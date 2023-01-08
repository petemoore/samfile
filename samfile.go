// See https://www.worldofsam.org/products/samdos and https://sam.speccy.cz/systech/sam-coupe_tech-man_v3-0.pdf
// for details of the disk format
package samfile

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
)

type (
	FileType uint8

	DiskImage [819200]byte

	SectorData [512]byte

	FileEntry struct {
		Type                   FileType
		Name                   Filename
		Sectors                uint16
		FirstSector            *Sector
		SectorAddressMap       *SectorAddressMap
		FileTypeInfo           [11]byte
		StartAddressPage       uint8
		StartAddressPageOffset uint16
		Pages                  uint8
		LengthMod16K           uint16
		ExecutionAddressDiv16K uint8
		ExecutionAddressMod16K uint16
		SAMBASICStartLine      uint16
		MGTFlags               uint8
		MGTFutureAndPast       [10]byte
		ReservedA              [4]byte
		ReservedB              [11]byte
	}

	Filename         [10]byte
	SectorAddressMap [195]byte

	DiskJournal [80]*FileEntry

	FilePart struct {
		Data       [510]byte
		NextSector *Sector
	}

	FileHeader struct {
		Type         FileType
		LengthMod16K uint16
		PageOffset   uint16
		Pages        uint8
		StartPage    uint8
	}

	File struct {
		Header *FileHeader
		Body   []byte
	}

	Sector struct {
		// 0-79 or 128-207
		Track uint8
		// 1-10
		Sector uint8
	}
)

const (
	FT_ERASED      = FileType(0)
	FT_ZX_SNAPSHOT = FileType(5)
	FT_SAM_BASIC   = FileType(16)
	FT_NUM_ARRAY   = FileType(17)
	FT_STR_ARRAY   = FileType(18)
	FT_CODE        = FileType(19)
	FT_SCREEN      = FileType(20)
)

func (file *File) Output() {
	fmt.Printf("Type:                %v\n", file.Header.Type)
	fmt.Printf("Start:               %v\n", file.Header.Start())
	fmt.Printf("Length:              %v\n", file.Header.Length())
	fmt.Printf("Body:\n%v\n", file.Body)
}

func (fileHeader *FileHeader) Start() uint32 {
	return uint32(fileHeader.PageOffset&0x3fff) | uint32(fileHeader.StartPage&0x1f+1)<<14
}

func (fileHeader *FileHeader) Length() uint32 {
	return uint32(fileHeader.LengthMod16K&0x3fff) | uint32(fileHeader.Pages)<<14
}

func (sam *SectorAddressMap) String() string {
	out := ""
	h := make([]byte, hex.EncodedLen(len(sam[:])))
	hex.Encode(h, sam[:])
	return out + string(h)
}

func (sam *SectorAddressMap) filterSectors(used bool) []*Sector {
	sectors := []*Sector{}
	track := uint8(4)
	sector := uint8(1)
	for _, b := range sam {
		for j := 0; j < 8; j++ {
			if (b&0x1 == 1) == used {
				sectors = append(sectors,
					&Sector{
						Track:  track,
						Sector: sector,
					},
				)
			}
			sector++
			if sector == 11 {
				sector = 1
				track++
				if track == 80 {
					track = 128
				}
			}
			b >>= 1
		}
	}
	return sectors
}

func (sam *SectorAddressMap) UsedSectors() []*Sector {
	return sam.filterSectors(true)
}

func (sam *SectorAddressMap) FreeSectors() []*Sector {
	return sam.filterSectors(false)
}

func (sam *SectorAddressMap) Merge(s *SectorAddressMap) {
	for i := range sam {
		sam[i] = sam[i] | s[i]
	}
}

func (ft FileType) String() string {
	switch ft {
	case FT_ERASED:
		return "Erased"
	case FT_ZX_SNAPSHOT:
		return "ZX Snapshot"
	case FT_SAM_BASIC:
		return "SAM BASIC"
	case FT_NUM_ARRAY:
		return "Number Array"
	case FT_STR_ARRAY:
		return "String Array"
	case FT_CODE:
		return "Code"
	case FT_SCREEN:
		return "Screen"
	default:
		return fmt.Sprintf("UNKNOWN (%v)", uint8(ft))
	}
}

func (sector *Sector) String() string {
	return fmt.Sprintf("Track %v / Sector %v", sector.Track, sector.Sector)
}

func Load(filename string) (*DiskImage, error) {
	image, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Can't load disk image %q: %v", filename, err)
	}
	d := DiskImage{}
	copy(d[:], image)
	return &d, nil
}

func (di *DiskImage) Save(filename string) error {
	err := os.WriteFile(filename, di[:], 0400)
	if err != nil {
		return fmt.Errorf("ERROR: Can't write disk image %q: %v", filename, err)
	}
	return nil
}

func (i *DiskImage) SectorData(sector *Sector) (*SectorData, error) {
	if sector.Sector < 1 || sector.Sector > 10 {
		debug.PrintStack()
		return nil, fmt.Errorf("Sector out of range: %v", sector.Sector)
	}
	if (sector.Track >= 80 && sector.Track < 128) || sector.Track >= 208 {
		return nil, fmt.Errorf("Track out of range: %v (should be 0-79 or 128-207)", sector.Track)
	}
	start := sector.Offset()
	data := SectorData{}
	copy(data[:], i[start:])
	return &data, nil
}

func (data *SectorData) FilePart() *FilePart {
	fp := FilePart{
		NextSector: &Sector{
			Track:  data[510],
			Sector: data[511],
		},
	}
	copy(fp.Data[:], data[:])
	return &fp
}

func (i *DiskImage) DiskJournal() *DiskJournal {
	dl := DiskJournal{}
	index := 0
	for track := uint8(0); track < 4; track++ {
		for sector := uint8(1); sector <= 10; sector++ {
			sectorData, err := i.SectorData(
				&Sector{
					Track:  track,
					Sector: sector,
				},
			)
			if err != nil {
				log.Printf("Error reading directory listing: %v", err)
				continue
			}
			for offset := 0; offset < 512; offset += 256 {
				var raw [256]byte
				copy(raw[:], sectorData[offset:])
				dl[index] = FileEntryFrom(raw)
				index++
			}
		}
	}
	return &dl
}

func FileEntryFrom(data [0x100]byte) *FileEntry {
	fe := FileEntry{
		Type:    FileType(data[0x00]),
		Sectors: uint16(data[0x0b])<<8 | uint16(data[0x0c]), // big endian!
		FirstSector: &Sector{
			Track:  data[0x0d],
			Sector: data[0x0e],
		},
		MGTFlags:               data[0xdc],
		StartAddressPage:       data[0xec],
		StartAddressPageOffset: uint16(data[0xed]) | uint16(data[0xee])<<8,
		Pages:                  data[0xef],
		LengthMod16K:           uint16(data[0xf0]) | uint16(data[0xf1])<<8,
		ExecutionAddressDiv16K: data[0xf2],
		ExecutionAddressMod16K: uint16(data[0xf3]) | uint16(data[0xf4])<<8,
		SAMBASICStartLine:      uint16(data[0xf3]) | uint16(data[0xf4])<<8,
		SectorAddressMap:       &SectorAddressMap{},
		Name:                   Filename{},
	}
	copy(fe.Name[:], data[0x01:])
	copy(fe.SectorAddressMap[:], data[0x0f:])
	copy(fe.MGTFutureAndPast[:], data[0xd2:])
	copy(fe.FileTypeInfo[:], data[0xdd:])
	copy(fe.ReservedA[:], data[0xe8:])
	copy(fe.ReservedB[:], data[0xf5:])
	return &fe
}

func (fe *FileEntry) Raw() [0x100]byte {
	raw := [0x100]byte{}
	raw[0x00] = byte(fe.Type)
	// raw[1]..raw[10]
	for i := 0; i < 0x0a; i++ {
		raw[i+1] = fe.Name[i]
	}
	raw[0x0b] = byte((fe.Sectors >> 8) & 0xff) // big
	raw[0x0c] = byte(fe.Sectors & 0xff)        // endian !!
	raw[0x0d] = fe.FirstSector.Track
	raw[0x0e] = fe.FirstSector.Sector
	// raw[0x0f]..raw[0xd1]
	for i := 0; i < 0xc3; i++ {
		raw[i+0x0f] = fe.SectorAddressMap[i]
	}
	// raw[0xd2]..raw[0xdb]
	for i := 0; i < 0x0a; i++ {
		raw[i+0xd2] = fe.MGTFutureAndPast[i]
	}
	raw[0xdc] = fe.MGTFlags
	// raw[0xdd]..raw[0xe7]
	for i := 0; i < 0x0b; i++ {
		raw[i+0xdd] = fe.FileTypeInfo[i]
	}
	// raw[0xe8]..raw[0xeb]
	for i := 0; i < 0x04; i++ {
		raw[i+0xe8] = fe.ReservedA[i]
	}
	raw[0xec] = fe.StartAddressPage
	raw[0xed] = byte(fe.StartAddressPageOffset & 0xff)
	raw[0xee] = byte((fe.StartAddressPageOffset >> 8) & 0xff)
	raw[0xef] = fe.Pages
	raw[0xf0] = byte(fe.LengthMod16K & 0xff)
	raw[0xf1] = byte((fe.LengthMod16K) >> 8 & 0xff)
	raw[0xf2] = fe.ExecutionAddressDiv16K
	raw[0xf3] = byte(fe.ExecutionAddressMod16K & 0xff)
	raw[0xf4] = byte((fe.ExecutionAddressMod16K) >> 8 & 0xff)
	// raw[0xf5]..raw[0xff]
	for i := 0; i < 0x0b; i++ {
		raw[i+0xf5] = fe.ReservedB[i]
	}
	return raw
}

func (dj *DiskJournal) Output() {
	for _, fe := range dj {
		err := fe.Output()
		if err != nil {
			log.Printf("ERROR: %v", err)
		}
	}
}

func (dj *DiskJournal) CombinedSectorMap() *SectorAddressMap {
	sam := new(SectorAddressMap)
	for _, fe := range dj {
		sam.Merge(fe.SectorAddressMap)
	}
	return sam
}

func (dj *DiskJournal) filterFileEntries(used bool) []int {
	entries := []int{}
	for i, fe := range dj {
		if fe.Used() == used {
			entries = append(entries, i)
		}
	}
	return entries
}

func (dj *DiskJournal) UsedFileEntries() []int {
	return dj.filterFileEntries(true)
}

func (dj *DiskJournal) FreeFileEntries() []int {
	return dj.filterFileEntries(false)
}

func (fe *FileEntry) Used() bool {
	if strings.HasPrefix(fe.Type.String(), "UNKNOWN") {
		return false
	}
	if fe.FirstSector.Track == 0 {
		return false
	}
	return true
}

func (fe *FileEntry) Output() error {
	if !fe.Used() {
		return nil
	}
	if fe.FirstSector.Track < 4 {
		return fmt.Errorf("First sector has track < 4: %v", fe.FirstSector)
	}
	defer fmt.Println("")
	fmt.Printf("%q\n", fe.Name)
	fmt.Printf("  Type:                              %v\n", fe.Type)
	if fe.Type == FT_ERASED {
		return nil
	}
	switch fe.Type {
	case FT_NUM_ARRAY:
		fmt.Printf("  Number Array Info:                 %v\n", fe.FileTypeInfo)
	case FT_STR_ARRAY:
		fmt.Printf("  String Array Info:                 %v\n", fe.FileTypeInfo)
	case FT_SCREEN:
		fmt.Printf("  Screen Mode:                       %v\n", fe.FileTypeInfo[0])
	case FT_SAM_BASIC:
		fmt.Printf("  Program length:                    %v\n", fe.ProgramLength())
		fmt.Printf("  Numeric variables offset:          %v\n", fe.NumericVariableOffset())
		fmt.Printf("  String/array variables offset:     %v\n", fe.StringArrayVariableOffset())
	}
	fmt.Printf("  Start:                             %v\n", fe.StartAddress())
	fmt.Printf("  Length:                            %v\n", fe.Length())
	switch fe.Type {
	case FT_SAM_BASIC:
		fmt.Printf("  Start Line:                        %v\n", fe.SAMBASICStartLine)
	case FT_CODE:
		if fe.ExecutionAddressDiv16K != 255 {
			fmt.Printf("  Execution Address:                 %v\n", fe.ExecutionAddress())
		}
	}
	return nil
}

func (fe *FileEntry) ProgramLength() uint32 {
	return uint32(fe.FileTypeInfo[0])<<16 | uint32(fe.FileTypeInfo[1]) | uint32(fe.FileTypeInfo[2])<<8
}

func (fe *FileEntry) NumericVariableOffset() uint32 {
	return uint32(fe.FileTypeInfo[3])<<16 | uint32(fe.FileTypeInfo[4]) | uint32(fe.FileTypeInfo[5])<<8
}

func (fe *FileEntry) StringArrayVariableOffset() uint32 {
	return uint32(fe.FileTypeInfo[6])<<16 | uint32(fe.FileTypeInfo[7]) | uint32(fe.FileTypeInfo[8])<<8
}

func (fe *FileEntry) ExecutionAddress() uint32 {
	return uint32(fe.ExecutionAddressMod16K&0x3fff) | uint32(fe.ExecutionAddressDiv16K&0x1f)<<14
}

func (fe *FileEntry) StartAddress() uint32 {
	return uint32(fe.StartAddressPageOffset&0x3fff) | uint32(fe.StartAddressPage&0x1f+1)<<14
}

func (fe *FileEntry) Length() uint32 {
	return uint32(fe.LengthMod16K&0x3fff) | uint32(fe.Pages)<<14
}

func (di *DiskImage) File(filename string) (*File, error) {

	for _, fe := range di.DiskJournal() {
		if fe.Name.String() == filename {
			fileLength := fe.Length()
			raw := make([]byte, fileLength+9)
			sectorData, err := di.SectorData(fe.FirstSector)
			if err != nil {
				return nil, err
			}
			filepart := sectorData.FilePart()
			i := uint16(0)
			for {
				copy(raw[510*i:], filepart.Data[:])
				i++
				if i == fe.Sectors {
					break
				}
				sectorData, err = di.SectorData(filepart.NextSector)
				if err != nil {
					return nil, err
				}
				filepart = sectorData.FilePart()
			}
			file := &File{
				Header: &FileHeader{
					Type:         FileType(raw[0]),
					LengthMod16K: uint16(raw[1]) | uint16(raw[2])<<8,
					PageOffset:   uint16(raw[3]) | uint16(raw[4])<<8,
					Pages:        raw[7],
					StartPage:    raw[8] & 0x1f,
				},
				Body: raw[9:],
			}
			return file, nil
		}
	}
	return nil, fmt.Errorf("File %v not found", filename)
}

func (filename Filename) String() string {
	b := make([]byte, 0, 10)
	for _, k := range filename {
		if k == 0 {
			break
		}
		b = append(b, k)
	}
	return strings.TrimRight(string(b), " ")
}

func (di *DiskImage) AddCodeFile(name string, data []byte, loadAddress, executionAddress uint32) error {
	if loadAddress < 1<<14 {
		return fmt.Errorf("Load address %v of %q is in ROM but must be %v of higher to be loaded into RAM", loadAddress, name, 1<<14)
	}
	if int(loadAddress) > 1<<19-len(data) {
		return fmt.Errorf("Load address %v of %v byte file %q higher than maximum allowed %v", loadAddress, len(data), name, 1<<19-len(data))
	}
	if executionAddress > 0 && executionAddress < loadAddress {
		return fmt.Errorf("Execution address %v of %q lower than load address %v", executionAddress, name, loadAddress)
	}
	if int(executionAddress) >= int(loadAddress)+len(data) {
		return fmt.Errorf("Execution address %v of %q is higher than the memory region it is loaded to (%v to %v)", executionAddress, name, loadAddress, int(loadAddress)+len(data)-1)
	}
	fe := &FileEntry{
		Type:                   FT_CODE,
		StartAddressPage:       uint8(loadAddress>>14) - 1,
		StartAddressPageOffset: uint16(loadAddress & 0x3fff),
		ExecutionAddressDiv16K: 0xff,
		ExecutionAddressMod16K: 0xffff,
	}
	if executionAddress > 0 {
		fe.ExecutionAddressDiv16K = uint8(executionAddress>>14) - 1
		fe.ExecutionAddressMod16K = uint16((executionAddress & 0x3fff) | 0x8000)
	}
	return di.addFile(
		name,
		fe,
		data,
	)
}

func (fe *FileEntry) CreateHeader() *FileHeader {
	return &FileHeader{
		Type:         fe.Type,
		LengthMod16K: fe.LengthMod16K,
		PageOffset:   fe.StartAddressPageOffset,
		Pages:        fe.Pages,
		StartPage:    fe.StartAddressPage,
	}
}

func (di *DiskImage) addFile(name string, fe *FileEntry, data []byte) error {
	dj := di.DiskJournal()
	freeFileEntries := dj.FreeFileEntries()
	if len(freeFileEntries) < 1 {
		return fmt.Errorf("Cannot add file %q to disk; disk already contains maximum number of files (80).", name)
	}
	requiredSectorCount := (len(data) + 9 + 509) / 510
	freeSectors := dj.CombinedSectorMap().FreeSectors()
	if len(freeSectors) < requiredSectorCount {
		return fmt.Errorf("Cannot add file %q to disk; not enough space (%v free sectors required but only %v sectors available).", name, requiredSectorCount, len(freeSectors))
	}
	fe.Name = *(*[10]byte)([]byte(name + "          "))
	fe.Sectors = uint16(requiredSectorCount)
	fe.FirstSector = freeSectors[0]
	fe.Pages = uint8(len(data) >> 14)
	fe.LengthMod16K = uint16(len(data) & 0x3fff)
	fe.SectorAddressMap = &SectorAddressMap{}

	f := &File{
		Header: fe.CreateHeader(),
		Body:   data,
	}
	raw := f.Raw()

	sd := &SectorData{}
	for i := 0; i < requiredSectorCount; i++ {
		if i < requiredSectorCount-1 {
			copy(sd[:], raw[i*510:(i+1)*510])
			sd[510] = freeSectors[i+1].Track
			sd[511] = freeSectors[i+1].Sector
		} else {
			sd = &SectorData{} // otherwise sd has non-zero values
			copy(sd[:], raw[i*510:])
		}
		offset, mask := freeSectors[i].SAMMask()
		fe.SectorAddressMap[offset] |= byte(mask)
		di.WriteSector(freeSectors[i], sd)
	}
	dj[freeFileEntries[0]] = fe
	di.WriteFileEntry(dj, freeFileEntries[0])
	return nil
}

func (di *DiskImage) WriteFileEntry(dj *DiskJournal, index int) {
	offset := index << 8
	rawFileEntry := dj[index].Raw()
	for i, b := range rawFileEntry {
		di[i+offset] = b
	}
}

func (sector *Sector) SAMMask() (offset uint8, mask uint8) {
	bitOffset := (int(sector.Track)&0x7f)*10 + int(sector.Sector) - 1 + ((int(sector.Track)&0x80)>>7)*800 - 40
	return uint8(bitOffset >> 3), 1 << bitOffset & 0x07
}

func (sector *Sector) Offset() int {
	return int(sector.Track>>7)*5120 + (int(sector.Sector)-1)*512 + int(sector.Track&0x7f)*10240
}

func (di *DiskImage) WriteSector(sector *Sector, sd *SectorData) {
	offset := sector.Offset()
	for i, b := range sd {
		di[i+offset] = b
	}
}

func (fh *FileHeader) Raw() [9]byte {
	return [9]byte{
		byte(fh.Type),
		byte(fh.LengthMod16K),
		byte(fh.LengthMod16K >> 8),
		byte(fh.PageOffset),
		byte(fh.PageOffset >> 8),
		0,
		0,
		fh.Pages,
		fh.StartPage,
	}
}

func (file *File) Raw() []byte {
	h := file.Header.Raw()
	return append(h[:], file.Body...)
}
