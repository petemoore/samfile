package samfile

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type (
	FileType uint8

	DiskImage [819200]byte

	SectorData [512]byte

	FileEntry struct {
		Type                       FileType
		Name                       Filename
		Sectors                    uint16
		FirstSector                Sector
		SectorAddressMap           SectorAddressMap
		FileTypeInfo               [11]byte
		StartAddressPage           uint8
		StartAddressPageOffset     uint16
		Pages                      uint8
		LengthMod16K               uint16
		ExecutionAddressPage       uint8
		ExecutionAddressPageOffset uint16
		SAMBasicStartLine          uint16
	}

	Filename         [10]byte
	SectorAddressMap [195]byte

	DirectoryListing [80]*FileEntry

	FilePart struct {
		Data       [510]byte
		NextSector Sector
	}

	FileHeader struct {
		Type         FileType
		LengthMod16K uint16
		PageOffset   uint16
		Pages        uint8
		StartPage    uint8
	}

	File struct {
		Header FileHeader
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

func (sam SectorAddressMap) String() string {
	out := ""
	sector := Sector{
		Track:  4,
		Sector: 1,
	}
	for _, b := range sam {
		for j := 0; j < 8; j++ {
			if b&0x1 == 0x1 {
				out += "    " + sector.String() + "\n"
			}
			sector.Sector++
			if sector.Sector == 11 {
				sector.Sector = 1
				sector.Track++
				if sector.Track == 80 {
					sector.Track = 128
				}
			}
			b >>= 1
		}
	}
	h := make([]byte, hex.EncodedLen(len(sam[:])))
	hex.Encode(h, sam[:])
	return out + string(h)
}

func (ft FileType) String() string {
	switch ft {
	case FT_ERASED:
		return "Erased"
	case FT_ZX_SNAPSHOT:
		return "ZX Snapshot"
	case FT_SAM_BASIC:
		return "SAM Basic"
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

func (sector Sector) String() string {
	return fmt.Sprintf("Track %v / Sector %v", sector.Track, sector.Sector)
}

func Load(file string) (*DiskImage, error) {
	image, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Can't load disk2.mgt: %v", err)
	}
	d := DiskImage{}
	copy(d[:], image)
	return &d, nil
}

func (i *DiskImage) SectorData(sector Sector) (*SectorData, error) {
	if sector.Sector < 1 || sector.Sector > 10 {
		return nil, fmt.Errorf("Sector out of range: %v", sector.Sector)
	}
	if (sector.Track >= 80 && sector.Track < 128) || sector.Track >= 208 {
		return nil, fmt.Errorf("Track out of range: %v (should be 0-79 or 128-207)", sector.Track)
	}
	start := 512 * (20*uint(sector.Track) + uint(sector.Sector) - 1)
	if sector.Track >= 128 {
		start -= 255 * 10 * 512
	}
	data := SectorData{}
	copy(data[:], i[start:])
	return &data, nil
}

func (data *SectorData) FilePart() *FilePart {
	fp := FilePart{
		NextSector: Sector{
			Track:  data[510],
			Sector: data[511],
		},
	}
	copy(fp.Data[:], data[:])
	return &fp
}

func (i *DiskImage) DirectoryListing() *DirectoryListing {
	dl := DirectoryListing{}
	index := 0
	for track := uint8(0); track < 4; track++ {
		for sector := uint8(1); sector <= 10; sector++ {
			sectorData, err := i.SectorData(
				Sector{
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

func FileEntryFrom(data [256]byte) *FileEntry {
	fe := FileEntry{
		Type:    FileType(data[0]),
		Sectors: uint16(data[11])<<8 | uint16(data[12]),
		FirstSector: Sector{
			Track:  data[13],
			Sector: data[14],
		},
		StartAddressPage:           data[236],
		StartAddressPageOffset:     uint16(data[237]) | uint16(data[238])<<8,
		Pages:                      data[239],
		LengthMod16K:               uint16(data[240]) | uint16(data[241])<<8,
		ExecutionAddressPage:       data[242],
		ExecutionAddressPageOffset: uint16(data[243]) | uint16(data[244])<<8,
		SAMBasicStartLine:          uint16(data[243]) | uint16(data[244])<<8,
	}
	copy(fe.Name[:], data[1:])
	copy(fe.SectorAddressMap[:], data[15:])
	copy(fe.FileTypeInfo[:], data[221:])
	return &fe
}

func (dl *DirectoryListing) Output() {
	for _, fe := range dl {
		err := fe.Output()
		if err != nil {
			log.Printf("ERROR: %v", err)
		}
	}
}

func (fe *FileEntry) Free() bool {
	if strings.HasPrefix(fe.Type.String(), "UNKNOWN") {
		return true
	}
	if fe.FirstSector.Track == 0 {
		return true
	}
	return false
}

func (fe *FileEntry) Output() error {
	if fe.Free() {
		return nil
	}
	if fe.FirstSector.Track < 4 {
		return fmt.Errorf("First sector has track < 4: %v", fe.FirstSector)
	}
	defer fmt.Println("")
	fmt.Printf("%q\n", fe.Name)
	fmt.Printf("  Type:                              %v\n", fe.Type)
	// fmt.Printf("  Sectors:                           %v\n", fe.Sectors)
	if fe.Type == FT_ERASED {
		return nil
	}
	// fmt.Printf("  First Sector:                      %v\n", fe.FirstSector)
	// fmt.Printf("  Sector Address Map:\n%v\n", fe.SectorAddressMap)
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
		fmt.Printf("  Start Line:                        %v\n", fe.SAMBasicStartLine)
	case FT_CODE:
		if fe.ExecutionAddressPage != 255 {
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
	return uint32(fe.ExecutionAddressPageOffset&0x3fff) | uint32(fe.ExecutionAddressPage&0x1f)<<14
}

func (fe *FileEntry) StartAddress() uint32 {
	return uint32(fe.StartAddressPageOffset&0x3fff) | uint32(fe.StartAddressPage&0x1f+1)<<14
}

func (fe *FileEntry) Length() uint32 {
	return uint32(fe.LengthMod16K&0x3fff) | uint32(fe.Pages)<<14
}

func (di *DiskImage) File(filename string) (*File, error) {

	for _, fe := range di.DirectoryListing() {
		if fe.Name.String() == filename {
			fileLength := fe.Length()
			raw := make([]byte, fileLength+9, fileLength+9)
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
				Header: FileHeader{
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

func (fe *FileEntry) File() *File {
	f := File{}
	return &f
}
