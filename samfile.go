// Package samfile reads, inspects and modifies individual files inside a
// SAM Coupé MGT floppy disk image (the 819200-byte .mgt format written by
// SAMDOS and used by the SAM emulator ecosystem). For whole-disk
// operations and format conversions, use samdisk
// (https://simonowen.com/samdisk/) instead — samfile only touches the
// contents of an existing image.
//
// # MGT image layout
//
// An MGT image is 80 cylinders × 2 sides × 10 sectors × 512 bytes,
// stored cylinder-interleaved (side 0 of each cylinder followed by side
// 1). Sides are encoded in bit 7 of the track byte: tracks 0–79 are
// side 0, tracks 128–207 are side 1. Tracks 0–3 of side 0 hold the
// SAMDOS directory; tracks 4–79 and 128–207 hold file bodies.
//
// The directory is 80 fixed-position slots, each 256 bytes, packed two
// per sector. Every slot is a [FileEntry] describing one file (or sits
// erased with a zero Type byte). The data area is a pool of 1560
// 512-byte sectors; a file occupies a chain of sectors whose payload is
// the first 510 bytes of each sector, with bytes 510–511 holding the
// (track, sector) link to the next sector ((0, 0) marks end-of-file).
//
// # API model
//
//   - [Load] reads an .mgt file into a [*DiskImage] (and rejects EDSK
//     format, which must be converted with samdisk first).
//   - [DiskImage.DiskJournal] parses the 80-slot directory into a
//     [*DiskJournal] of [*FileEntry].
//   - [DiskImage.File] walks the sector chain for a named file and
//     returns its assembled [*File] (9-byte [FileHeader] + body bytes).
//   - [DiskImage.AddCodeFile] writes a new code/data file to a free
//     slot and free sectors, updating both the directory and the
//     sector chain.
//   - [DiskImage.Save] writes the (possibly modified) image back to
//     disk.
//
// SAM BASIC programs are stored tokenised; [SAMBasic.Output]
// detokenises a body into a plain-text listing.
//
// Authoritative format references: the SAM Coupé Technical Manual v3.0
// (https://sam.speccy.cz/systech/sam-coupe_tech-man_v3-0.pdf), the
// annotated SAM ROM v3.0 disassembly, and the SAMDOS source.
package samfile

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/petemoore/samfile/v3/sambasic"
)

type (
	// FileType is the value of the status / file-type byte at offset
	// 0x00 of a directory entry (and offset 0 of the file body). A
	// value of 0 marks the slot erased; valid file types are listed
	// in the FT_* constants. SAMDOS also uses bit 7 of this byte for
	// HIDDEN and bit 6 for PROTECTED — those attribute bits are not
	// modelled by FileType and must be masked off before comparing.
	FileType uint8

	// DiskImage is the raw byte image of one MGT floppy: 80
	// cylinders × 2 sides × 10 sectors × 512 bytes, stored
	// cylinder-interleaved (see the package overview). Use Load and
	// Save to read and write it; use SectorData, DiskJournal, File
	// and AddCodeFile to read or modify its contents.
	DiskImage [819200]byte

	// SectorData is the 512 bytes of one disk sector. For data
	// sectors, the first 510 bytes are file payload and bytes
	// 510–511 are the (track, sector) link to the next sector in the
	// chain (see FilePart). For directory sectors, all 512 bytes
	// belong to two packed FileEntry slots.
	SectorData [512]byte

	// FileEntry is a parsed 256-byte SAMDOS directory entry. The
	// disk holds 80 such entries — packed two per sector across
	// tracks 0–3 of side 0 — and an entry describes either an
	// existing file or a free slot (Type == 0). The struct fields
	// mirror the on-disk byte layout one-for-one; see FileEntryFrom
	// for the byte map.
	//
	// StartAddressPage, StartAddressPageOffset, Pages and
	// LengthMod16K are duplicates of the first 5 of the 9 bytes of
	// the file body's FileHeader — written into the directory at
	// SAVE time so that LS-style operations don't have to read the
	// file body. They should always match the body header.
	FileEntry struct {
		Type        FileType
		Name        Filename
		Sectors     uint16 // big-endian on disk: count of 512-byte sectors the file occupies
		FirstSector *Sector
		// SectorAddressMap is the per-file 1560-bit bitmap of which
		// data sectors this file occupies (see [SectorAddressMap]).
		SectorAddressMap *SectorAddressMap
		// FileTypeInfo is the 11-byte type-dependent metadata block
		// at directory bytes 0xDD–0xE7. For FT_SAM_BASIC it holds
		// three 3-byte PAGEFORM length triplets (program / +numeric
		// vars / +gap); for FT_SCREEN, byte 0 is the screen MODE;
		// for FT_NUM_ARRAY and FT_STR_ARRAY it holds the array's
		// type/length byte plus name; for FT_CODE it is unused.
		FileTypeInfo           [11]byte
		StartAddressPage       uint8  // mirror of FileHeader.StartPage
		StartAddressPageOffset uint16 // mirror of FileHeader.PageOffset
		Pages                  uint8  // mirror of FileHeader.Pages
		LengthMod16K           uint16 // mirror of FileHeader.LengthMod16K
		// ExecutionAddressDiv16K and ExecutionAddressMod16K together
		// encode the auto-execution address of an FT_CODE file in
		// SAM's REL PAGE FORM. Set ExecutionAddressDiv16K = 0xFF to
		// disable auto-exec; see ExecutionAddress for the decoded
		// linear address.
		ExecutionAddressDiv16K uint8
		ExecutionAddressMod16K uint16
		// SAMBASICStartLine occupies the same two bytes as
		// ExecutionAddressMod16K: for FT_SAM_BASIC it is the
		// auto-RUN line number (or 0xFFFF / ExecutionAddressDiv16K =
		// 0xFF to opt out of auto-RUN).
		SAMBASICStartLine uint16
		MGTFlags          uint8 // "MGT use only" per the Tech Manual
		// MGTFutureAndPast occupies directory bytes 0xD2–0xDB. The
		// Tech Manual labels these bytes unused, but SAMDOS in fact
		// caches the file body's 9-byte FileHeader at bytes 1–9
		// (i.e. 0xD3–0xDB) as an in-RAM scratchpad. Byte 0 (0xD2)
		// is unused.
		MGTFutureAndPast [10]byte
		ReservedA        [4]byte  // spare 4 bytes at 0xE8–0xEB
		ReservedB        [11]byte // spare 11 bytes at 0xF5–0xFF
	}

	// Filename is the 10-byte, space-padded ASCII filename stored at
	// directory bytes 0x01–0x0A. Use String to obtain a trimmed,
	// printable form; SAMDOS matches filenames case-insensitively.
	Filename [10]byte

	// SectorAddressMap is the 195-byte / 1560-bit allocation bitmap
	// at directory bytes 0x0F–0xD1. Each bit corresponds to one of
	// the disk's 1560 data sectors. Bit 0 of byte 0 is (track 4,
	// sector 1); bits then proceed sector-then-track within side 0
	// up to (track 79, sector 10) at bit 759, and continue with the
	// side-1 data sectors at bit 760 onwards up to (track 207,
	// sector 10) at bit 1559. See Sector.SAMMask for the
	// per-sector bit-position computation.
	//
	// In a FileEntry, the map records the sectors that file
	// occupies. The disk-wide free map is not stored on disk; it is
	// computed at allocation time as the bitwise OR of every entry's
	// map (see DiskJournal.CombinedSectorMap).
	SectorAddressMap [195]byte

	// DiskJournal is the parsed SAMDOS directory: 80 *FileEntry
	// values in their on-disk slot order. Slot 0 is the first half
	// of (track 0, sector 1); slot 79 is the second half of (track
	// 3, sector 10). An unused slot has FileEntry.Used returning
	// false (Type byte == 0).
	DiskJournal [80]*FileEntry

	// FilePart is one 510-byte payload chunk of a file body
	// together with the (track, sector) link to the next chunk.
	// NextSector.Track == 0 marks the last chunk of a file.
	FilePart struct {
		Data       [510]byte
		NextSector *Sector
	}

	// FileHeader is the 9-byte header that prefixes every file
	// body on disk (file body bytes 0–8). The directory entry
	// duplicates these fields into its StartAddressPage /
	// StartAddressPageOffset / Pages / LengthMod16K members, but
	// the body header is what the ROM consumes when LOADing.
	//
	// LengthMod16K plus Pages gives the body length excluding the
	// header (see Length); PageOffset and StartPage together
	// encode the SAM address the file should be loaded to in REL
	// PAGE FORM (see Start).
	//
	// ExecutionAddressDiv16K and ExecutionAddressMod16KLo mirror
	// the directory entry's auto-execution-address gate at body-
	// header bytes 5 and 6. Setting both to 0xFF signals "no
	// auto-exec" — the convention the SAM ROM's LOAD-CODE path
	// at rom-disasm:22471-22484 checks to return cleanly after a
	// load instead of jumping to the file's start address. The
	// directory entry holds the full 16-bit ExecutionAddressMod16K
	// at offsets 0xF3-0xF4; only the low byte fits in the body
	// header, which is what these fields represent.
	FileHeader struct {
		Type                     FileType
		LengthMod16K             uint16
		PageOffset               uint16
		ExecutionAddressDiv16K   uint8
		ExecutionAddressMod16KLo uint8
		Pages                    uint8
		StartPage                uint8
	}

	// File is a complete file as read from disk: the 9-byte
	// FileHeader followed by Body. Body excludes the header bytes.
	File struct {
		Header *FileHeader
		Body   []byte
	}

	// Sector identifies a (Track, Sector) location on disk.
	// Track uses SAMDOS's side-encoding (bit 7 = side bit): values
	// 0–79 are side 0, values 128–207 are side 1; values 80–127
	// and 208+ are invalid. Sector is 1-based and runs 1–10.
	// Tracks 0–3 of side 0 hold the directory; tracks 4–79 and
	// 128–207 hold file data. Use Offset to map a Sector to a
	// byte offset within a DiskImage.
	Sector struct {
		Track  uint8 // 0–79 (side 0) or 128–207 (side 1)
		Sector uint8 // 1–10
	}
)

// SAMDOS file types — values of the status / file-type byte at offset
// 0x00 of a directory entry (and byte 0 of the file body). FT_ERASED
// (0) is the "free slot" sentinel; the others are the public SAM
// types defined by the Tech Manual. Bit 7 of the byte marks the file
// HIDDEN and bit 6 marks it PROTECTED — mask them off before comparing
// against these constants.
const (
	FT_ERASED      = FileType(0)  // slot is unused
	FT_ZX_SNAPSHOT = FileType(5)  // 48K ZX Spectrum snapshot (SAMDOS extension)
	FT_SAM_BASIC   = FileType(16) // tokenised SAM BASIC program; body decoded by SAMBasic.Output
	FT_NUM_ARRAY   = FileType(17) // saved numeric array
	FT_STR_ARRAY   = FileType(18) // saved string array
	FT_CODE        = FileType(19) // arbitrary code/data blob with load (and optional execution) address
	FT_SCREEN      = FileType(20) // SCREEN$ — display memory dump; mode stored at FileTypeInfo[0]
)

// Output prints a debug summary of file to stdout: type, start address,
// body length, and the raw body bytes.
func (file *File) Output() {
	fmt.Printf("Type:                %v\n", file.Header.Type)
	fmt.Printf("Start:               %v\n", file.Header.Start())
	fmt.Printf("Length:              %v\n", file.Header.Length())
	fmt.Printf("Body:\n%v\n", file.Body)
}

// Start decodes the file's load address from the header's REL PAGE FORM
// encoding (StartPage's low 5 bits give the 16K-page index, PageOffset's
// low 14 bits give the offset within that page). The returned address
// is a linear offset into SAM's 512K address space.
func (fileHeader *FileHeader) Start() uint32 {
	return uint32(fileHeader.PageOffset&0x3fff) | uint32((fileHeader.StartPage&0x1f)+1)<<14
}

// Length is the size in bytes of the file body, excluding the 9-byte
// header itself. Decoded as Pages × 16384 + (LengthMod16K & 0x3FFF).
func (fileHeader *FileHeader) Length() uint32 {
	return uint32(fileHeader.LengthMod16K&0x3fff) | uint32(fileHeader.Pages)<<14
}

// String returns the map as a 390-character lowercase hex dump (two
// chars per byte, no spacing).
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

// UsedSectors returns the Sector locations whose bits are set in sam,
// in disk order. For a FileEntry's map this is the file's sector
// chain; for a CombinedSectorMap it is every allocated sector on the
// disk.
func (sam *SectorAddressMap) UsedSectors() []*Sector {
	return sam.filterSectors(true)
}

// FreeSectors returns the Sector locations whose bits are clear in
// sam, in disk order. Most useful on a CombinedSectorMap, where it
// enumerates the disk's available data sectors.
func (sam *SectorAddressMap) FreeSectors() []*Sector {
	return sam.filterSectors(false)
}

// Merge sets in sam every bit that is set in s (bitwise OR, in place).
// Used to accumulate per-file maps into the disk-wide allocation map
// (see DiskJournal.CombinedSectorMap).
func (sam *SectorAddressMap) Merge(s *SectorAddressMap) {
	for i := range sam {
		sam[i] = sam[i] | s[i]
	}
}

// String returns the file type's mnemonic ("Code", "SAM BASIC",
// "Screen", etc.). Values outside the documented FT_* set render as
// "UNKNOWN (n)".
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

// String formats sector as "Track N / Sector M" with the raw
// side-encoded Track value (so side-1 locations show as 128–207).
func (sector *Sector) String() string {
	return fmt.Sprintf("Track %v / Sector %v", sector.Track, sector.Sector)
}

// edskMagic is the prefix every Extended CPC DSK image starts with. The
// full on-disk magic is "EXTENDED CPC DSK File\r\nDisk-Info\r\n" (34
// bytes); the 21-byte prefix below is unambiguous and matches both EDSK
// and any future EDSK-derived variants that share the literal.
// Reference: https://www.cpcwiki.eu/index.php/Format:DSK_disk_image_file_format
var edskMagic = []byte("EXTENDED CPC DSK File")

// Load reads filename as a raw MGT disk image. Files smaller than
// 819200 bytes are zero-padded on the right; files larger than that
// are truncated. Extended CPC DSK images ("EDSK" — magic bytes
// "EXTENDED CPC DSK File") are detected and rejected with a message
// pointing at samdisk, which can convert them to MGT.
func Load(filename string) (*DiskImage, error) {
	image, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error: can't load disk image %q: %v", filename, err)
	}
	if len(image) >= len(edskMagic) && bytes.Equal(image[:len(edskMagic)], edskMagic) {
		return nil, fmt.Errorf("error: EDSK format not supported; convert to MGT with samdisk (https://simonowen.com/samdisk/): samdisk %s OUTPUT.mgt", filename)
	}
	d := DiskImage{}
	copy(d[:], image)
	return &d, nil
}

// Save writes the whole 819200-byte image to filename, with mode 0400
// (read-only for the owner). The destination is overwritten if it
// already exists.
func (di *DiskImage) Save(filename string) error {
	err := os.WriteFile(filename, di[:], 0400)
	if err != nil {
		return fmt.Errorf("error: can't write disk image %q: %v", filename, err)
	}
	return nil
}

// SectorData returns a copy of the 512 bytes at sector's location.
// Returns an error if sector.Sector is outside 1–10 or sector.Track
// falls in the invalid gap 80–127 or above 207.
func (i *DiskImage) SectorData(sector *Sector) (*SectorData, error) {
	if sector.Sector < 1 || sector.Sector > 10 {
		return nil, fmt.Errorf("sector out of range: %v", sector.Sector)
	}
	if (sector.Track >= 80 && sector.Track < 128) || sector.Track >= 208 {
		return nil, fmt.Errorf("track out of range: %v (should be 0-79 or 128-207)", sector.Track)
	}
	start := sector.Offset()
	data := SectorData{}
	copy(data[:], i[start:])
	return &data, nil
}

// FilePart reinterprets data as one chunk of a file's sector chain:
// the first 510 bytes become the payload and the trailing 2 bytes
// become the (track, sector) link to the next chunk. A returned
// NextSector with Track == 0 marks end-of-file.
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

// DiskJournal parses the 80-slot SAMDOS directory from tracks 0–3 of
// side 0. Slots are returned in on-disk order, with each
// FileEntry's fields populated from its 256 directory bytes; use
// FileEntry.Used to tell occupied slots from free ones. The returned
// directory is a snapshot — mutating it does not change the
// underlying image until WriteFileEntry is called.
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
				log.Printf("error reading directory listing: %v", err)
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

// FileEntryFrom parses the 256 bytes of a single SAMDOS directory
// slot into a FileEntry. The on-disk byte map is:
//
//	0x00      Type / status byte
//	0x01–0x0A Filename (10 bytes, space-padded)
//	0x0B–0x0C Sector count (big-endian!)
//	0x0D–0x0E First-sector track and sector
//	0x0F–0xD1 SectorAddressMap (195 bytes)
//	0xD2–0xDB MGTFutureAndPast (SAMDOS scratchpad — see FileEntry doc)
//	0xDC      MGTFlags
//	0xDD–0xE7 FileTypeInfo (11 bytes, type-dependent)
//	0xE8–0xEB ReservedA
//	0xEC      StartAddressPage
//	0xED–0xEE StartAddressPageOffset (little-endian)
//	0xEF      Pages
//	0xF0–0xF1 LengthMod16K (little-endian)
//	0xF2      ExecutionAddressDiv16K
//	0xF3–0xF4 ExecutionAddressMod16K / SAMBASICStartLine (little-endian)
//	0xF5–0xFF ReservedB
//
// Inverse of FileEntry.Raw.
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

// Raw encodes fe back into the 256 raw bytes of a SAMDOS directory
// entry. Inverse of FileEntryFrom.
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
	raw[0xf1] = byte((fe.LengthMod16K >> 8) & 0xff)
	raw[0xf2] = fe.ExecutionAddressDiv16K
	raw[0xf3] = byte(fe.ExecutionAddressMod16K & 0xff)
	raw[0xf4] = byte((fe.ExecutionAddressMod16K >> 8) & 0xff)
	// raw[0xf5]..raw[0xff]
	for i := 0; i < 0x0b; i++ {
		raw[i+0xf5] = fe.ReservedB[i]
	}
	return raw
}

// Output prints a per-file summary of every occupied directory slot
// to stdout — the format used by the `samfile ls` command. Per-entry
// errors are logged but do not stop the walk.
func (dj *DiskJournal) Output() {
	for _, fe := range dj {
		err := fe.Output()
		if err != nil {
			log.Printf("error: %v", err)
		}
	}
}

// CombinedSectorMap returns the bitwise OR of every entry's
// SectorAddressMap. The result is the disk-wide allocation bitmap:
// a set bit means "some file owns this sector", a clear bit means
// "this sector is free". SAMDOS doesn't persist this bitmap on disk;
// it reconstructs it on demand at allocation time, and so does
// samfile (see AddCodeFile).
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

// UsedFileEntries returns the slot indices (0–79) whose entries are
// occupied (FileEntry.Used == true).
func (dj *DiskJournal) UsedFileEntries() []int {
	return dj.filterFileEntries(true)
}

// FreeFileEntries returns the slot indices (0–79) that are available
// for new files. AddCodeFile populates the lowest free slot.
func (dj *DiskJournal) FreeFileEntries() []int {
	return dj.filterFileEntries(false)
}

// Used reports whether the directory slot is occupied. SAMDOS itself
// only treats slots with Type == 0 as erased; samfile additionally
// rejects slots whose Type byte is not one of the documented FT_*
// values or whose FirstSector.Track is 0 (the directory tracks, which
// can never be a file's first sector).
func (fe *FileEntry) Used() bool {
	if strings.HasPrefix(fe.Type.String(), "UNKNOWN") {
		return false
	}
	if fe.FirstSector.Track == 0 {
		return false
	}
	return true
}

// Output prints a per-field human-readable summary of fe to stdout
// (the per-entry block of `samfile ls` output). Unused slots are
// silently skipped. Returns an error if FirstSector.Track is in the
// directory area (0–3, which is structurally invalid for a file).
func (fe *FileEntry) Output() error {
	if !fe.Used() {
		return nil
	}
	if fe.FirstSector.Track < 4 {
		return fmt.Errorf("first sector has track < 4: %v", fe.FirstSector)
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
		fmt.Printf("  Numeric variables size:            %v\n", fe.NumericVariablesSize())
		fmt.Printf("  Gap size:                          %v\n", fe.GapSize())
		fmt.Printf("  String/array variables size:       %v\n", fe.StringArrayVariablesSize())
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

// pageFormLength decodes a 19-bit length stored in SAM Coupé "PAGEFORM":
// byte 0 is a page count (16384 bytes per page); bytes 1-2 are a
// little-endian 16-bit address in section C (0x8000-0xBFFF) whose low
// 14 bits carry the in-page offset (bit 15 is always 1 by SCF;RR H,
// bit 14 always 0). The linear length is therefore
// page * 16384 + (raw_addr & 0x3fff). See ROM disasm RDTHREE
// (sam-coupe_rom-v3.0_annotated-disassembly.txt:7654-7659) and PAGEFORM
// (sam-coupe_rom-v3.0_annotated-disassembly.txt:7578-7589).
func pageFormLength(b0, b1, b2 byte) uint32 {
	return uint32(b0)*16384 + uint32(uint16(b1)|uint16(b2)<<8)&0x3fff
}

// SAM BASIC file layout (program area, in order):
//   [program text] [numeric variables] [gap] [string/array variables]
// The three FileTypeInfo length fields encode the cumulative offsets of the
// section boundaries; the four section sizes below are derived from them
// plus the file's total Length, and are easier for callers to reason
// about than the raw cumulative offsets. Per Tech Manual L4370-4382 and
// ROM disasm L16005-16012 ("CDE=LEN OF PROG ALONE" / "CDE=LEN OF
// PROG+NVARS+GAP").

// ProgramLength is the size in bytes of the tokenised SAM BASIC
// program text, excluding the trailing numeric-variables / gap /
// string-array sections. Only meaningful for FT_SAM_BASIC entries.
func (fe *FileEntry) ProgramLength() uint32 {
	return pageFormLength(fe.FileTypeInfo[0], fe.FileTypeInfo[1], fe.FileTypeInfo[2])
}

// NumericVariablesSize is the size in bytes of the numeric-variables
// section that immediately follows the program text in a saved SAM
// BASIC file. Only meaningful for FT_SAM_BASIC entries.
func (fe *FileEntry) NumericVariablesSize() uint32 {
	return pageFormLength(fe.FileTypeInfo[3], fe.FileTypeInfo[4], fe.FileTypeInfo[5]) -
		pageFormLength(fe.FileTypeInfo[0], fe.FileTypeInfo[1], fe.FileTypeInfo[2])
}

// GapSize is the size in bytes of the empty gap SAM BASIC leaves
// between the numeric variables and the string/array variables. On
// canonical SAVEs this is usually 512 (the ROM's MAKEROOM
// pre-allocation). Only meaningful for FT_SAM_BASIC entries.
func (fe *FileEntry) GapSize() uint32 {
	return pageFormLength(fe.FileTypeInfo[6], fe.FileTypeInfo[7], fe.FileTypeInfo[8]) -
		pageFormLength(fe.FileTypeInfo[3], fe.FileTypeInfo[4], fe.FileTypeInfo[5])
}

// StringArrayVariablesSize is the size in bytes of the string and
// array variables section, which occupies the remainder of the body
// after the gap. Only meaningful for FT_SAM_BASIC entries.
func (fe *FileEntry) StringArrayVariablesSize() uint32 {
	return fe.Length() - pageFormLength(fe.FileTypeInfo[6], fe.FileTypeInfo[7], fe.FileTypeInfo[8])
}

// ExecutionAddress decodes the auto-execution address of an FT_CODE
// entry from its REL PAGE FORM encoding. The caller should first
// check that ExecutionAddressDiv16K is not 0xFF (the sentinel for
// "no auto-execution address set").
func (fe *FileEntry) ExecutionAddress() uint32 {
	return uint32(fe.ExecutionAddressMod16K&0x3fff) | uint32(fe.ExecutionAddressDiv16K&0x1f)<<14
}

// StartAddress decodes the linear SAM address the file body should
// be loaded to, from the REL PAGE FORM encoding mirrored from the
// FileHeader. Equivalent to calling Start on the corresponding
// FileHeader.
func (fe *FileEntry) StartAddress() uint32 {
	return uint32(fe.StartAddressPageOffset&0x3fff) | uint32((fe.StartAddressPage&0x1f)+1)<<14
}

// Length is the size in bytes of the file body, excluding the 9-byte
// header. Decoded as Pages × 16384 + (LengthMod16K & 0x3FFF) —
// equivalent to FileHeader.Length on the same file.
func (fe *FileEntry) Length() uint32 {
	return uint32(fe.LengthMod16K&0x3fff) | uint32(fe.Pages)<<14
}

// File reads the named file out of the disk image, walking its
// sector chain from FirstSector and assembling the body. The match
// against filename is exact against the trimmed Filename.String() of
// each occupied directory entry — there is no wildcard or
// case-folding. The returned File.Header is reconstructed from the
// first 9 bytes of the body; File.Body is the remainder.
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
					Type:                     FileType(raw[0]),
					LengthMod16K:             uint16(raw[1]) | uint16(raw[2])<<8,
					PageOffset:               uint16(raw[3]) | uint16(raw[4])<<8,
					ExecutionAddressDiv16K:   raw[5],
					ExecutionAddressMod16KLo: raw[6],
					Pages:                    raw[7],
					StartPage:                raw[8] & 0x1f,
				},
				Body: raw[9:],
			}
			return file, nil
		}
	}
	return nil, fmt.Errorf("file %v not found", filename)
}

// String returns filename with any trailing NULs and spaces removed,
// suitable for matching against user input.
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

// AddCodeFile writes data to the disk image as a new SAMDOS type-19
// (CODE) file named name, allocating a free directory slot and the
// required free sectors. loadAddress is the SAM address the file
// will be loaded to (must be ≥ 0x4000 — anything below is in ROM —
// and the file must fit within SAM's 512K address space).
// executionAddress is optional: 0 records "no auto-exec", any other
// value is the address the loader will JP to after loading and must
// lie within the loaded region.
//
// Returns an error if the address validations fail, if the disk has
// no free directory slots (max 80 files), or if there are not enough
// free sectors to hold the data plus the 9-byte file header.
func (di *DiskImage) AddCodeFile(name string, data []byte, loadAddress, executionAddress uint32) error {
	if loadAddress < 1<<14 {
		return fmt.Errorf("load address %v of %q is in ROM but must be %v of higher to be loaded into RAM", loadAddress, name, 1<<14)
	}
	if int(loadAddress) > 1<<19-len(data) {
		return fmt.Errorf("load address %v of %v byte file %q higher than maximum allowed %v", loadAddress, len(data), name, 1<<19-len(data))
	}
	if executionAddress > 0 && executionAddress < loadAddress {
		return fmt.Errorf("execution address %v of %q lower than load address %v", executionAddress, name, loadAddress)
	}
	if int(executionAddress) >= int(loadAddress)+len(data) {
		return fmt.Errorf("execution address %v of %q is higher than the memory region it is loaded to (%v to %v)", executionAddress, name, loadAddress, int(loadAddress)+len(data)-1)
	}
	fe := &FileEntry{
		Type:                   FT_CODE,
		StartAddressPage:       uint8(loadAddress>>14) - 1,
		StartAddressPageOffset: uint16((loadAddress & 0x3fff) | 0x8000),
		ExecutionAddressDiv16K: 0xff,
		ExecutionAddressMod16K: 0xffff,
	}
	if executionAddress > 0 {
		fe.ExecutionAddressDiv16K = uint8(executionAddress >> 14)
		fe.ExecutionAddressMod16K = uint16((executionAddress & 0x3fff) | 0x8000)
	}
	return di.addFile(
		name,
		fe,
		data,
	)
}

func NewDiskImage() *DiskImage {
	return &DiskImage{}
}

// SetStartAddressPageUnusedBits sets the upper 3 bits (bits 7..5) of
// the StartAddressPage byte for the named file, preserving the low 5
// bits (the physical page index). The change is applied to both the
// directory entry (raw[0xEC]) and the matching body-header byte 8 in
// the file's first sector. Returns an error if bits > 7 or if the
// file is not present on disk.
//
// The low 5 bits of StartAddressPage are the page index (0–31) and
// are derived by AddCodeFile from the load address — this method
// intentionally cannot disturb them. The high 3 bits are unused: no
// known consumer in SAM ROM v3, SAMDOS 2, or MasterDOS reads them
// (ROM's TSURPG paging routine XORs them away via the
// `XOR / AND 0xE0 / XOR` idiom at rom-disasm:14852-14859; SAMDOS's
// hsave masks with `AND 0x1F` at h.s:140-143). Real-SAVE output on
// FRED 02 / Defender disks nonetheless records samdos2's
// StartAddressPage as 0x7D (= 3<<5 | 0x1D = page 29 with the top
// two bits set) — those bits are an unmasked leak of the HMPR
// video-mode flags at the moment SAM ROM's BASIC SAVE wrote the
// byte (see rom-disasm:22866-22869). This method exists to
// reproduce that historical accident byte-for-byte when comparing
// against canonical reference images; if you don't have that
// constraint, you don't need to call this.
func (di *DiskImage) SetStartAddressPageUnusedBits(name string, bits uint8) error {
	if bits > 7 {
		return fmt.Errorf("StartAddressPage unused bits value %d out of range (0..7)", bits)
	}
	dj := di.DiskJournal()
	for slot, fe := range dj {
		if !fe.Used() || fe.Name.String() != name {
			continue
		}
		value := (fe.StartAddressPage & 0x1F) | (bits << 5)
		fe.StartAddressPage = value
		fe.MGTFutureAndPast[9] = value
		di.WriteFileEntry(dj, slot)
		di[fe.FirstSector.Offset()+8] = value
		return nil
	}
	return fmt.Errorf("file %v not found", name)
}

func pageForm3Byte(value uint32) [3]byte {
	page := byte(value / 16384)
	offset := uint16(value%16384) | 0x8000
	return [3]byte{page, byte(offset & 0xFF), byte(offset >> 8)}
}

func (di *DiskImage) AddBasicFile(name string, file *sambasic.File) error {
	body := file.Bytes()

	fe := &FileEntry{
		Type:                   FT_SAM_BASIC,
		StartAddressPage:       0,
		StartAddressPageOffset: 0x9CD5,
		MGTFlags:               0x20,
	}

	if file.StartLine == 0xFFFF {
		fe.ExecutionAddressDiv16K = 0xFF
		fe.ExecutionAddressMod16K = 0xFFFF
		fe.SAMBASICStartLine = 0xFFFF
	} else {
		fe.ExecutionAddressDiv16K = 0x00
		fe.ExecutionAddressMod16K = file.StartLine
		fe.SAMBASICStartLine = file.StartLine
	}

	nvars := pageForm3Byte(file.NVARSOffset())
	numend := pageForm3Byte(file.NUMENDOffset())
	savars := pageForm3Byte(file.SAVARSOffset())
	copy(fe.FileTypeInfo[0:3], nvars[:])
	copy(fe.FileTypeInfo[3:6], numend[:])
	copy(fe.FileTypeInfo[6:9], savars[:])

	// addFile sets fe.Pages, fe.LengthMod16K, and mirrors the body
	// header into MGTFutureAndPast — no need to populate either here.
	return di.addFile(name, fe, body)
}

// CreateHeader synthesises the 9-byte FileHeader that should prefix
// the file body, by copying the directory entry's mirrored
// metadata fields (Type, length, start-page/offset). Used by
// addFile to construct the on-disk header for a new file.
//
// Body header bytes 5-6 carry the CODE execution address gate;
// for non-CODE files they are always 0xFF — BASIC auto-RUN is
// signalled via the directory entry's 0xF2-0xF4 bytes only
// (see sam-file-header.md §3 and test-mgt-byte-layout.md).
func (fe *FileEntry) CreateHeader() *FileHeader {
	execDiv := uint8(0xFF)
	execModLo := uint8(0xFF)
	if fe.Type == FT_CODE {
		execDiv = fe.ExecutionAddressDiv16K
		execModLo = byte(fe.ExecutionAddressMod16K & 0xff)
	}
	return &FileHeader{
		Type:                     fe.Type,
		LengthMod16K:             fe.LengthMod16K,
		PageOffset:               fe.StartAddressPageOffset,
		ExecutionAddressDiv16K:   execDiv,
		ExecutionAddressMod16KLo: execModLo,
		Pages:                    fe.Pages,
		StartPage:                fe.StartAddressPage,
	}
}

func (di *DiskImage) addFile(name string, fe *FileEntry, data []byte) error {
	dj := di.DiskJournal()
	freeFileEntries := dj.FreeFileEntries()
	if len(freeFileEntries) < 1 {
		return fmt.Errorf("cannot add file %q to disk; disk already contains maximum number of files (80).", name)
	}
	requiredSectorCount := (len(data) + 9 + 509) / 510
	freeSectors := dj.CombinedSectorMap().FreeSectors()
	if len(freeSectors) < requiredSectorCount {
		return fmt.Errorf("cannot add file %q to disk; not enough space (%v free sectors required but only %v sectors available).", name, requiredSectorCount, len(freeSectors))
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

	// Mirror the 9-byte body header into MGTFutureAndPast[1..9] so
	// the directory entry's MGT "future and past" region matches the
	// canonical real-SAVE convention: every byte of the body header
	// is duplicated in the dir entry. (MGTFutureAndPast[0] is
	// reserved and stays zero.) Without this mirror an inspector
	// reading just the dir entry would see all zeros for the body
	// header bytes that are otherwise authoritatively held there;
	// real disks saved by ROM SAVE populate this region.
	header := f.Header.Raw()
	copy(fe.MGTFutureAndPast[1:10], header[:])

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

// WriteFileEntry encodes dj[index] back into the 256 bytes of
// directory slot index. Call this after mutating an entry to commit
// the change to the disk image. No bounds checking on index.
func (di *DiskImage) WriteFileEntry(dj *DiskJournal, index int) {
	offset := index << 8
	rawFileEntry := dj[index].Raw()
	for i, b := range rawFileEntry {
		di[i+offset] = b
	}
}

// SAMMask returns the (byte offset, bit mask) within a
// 195-byte SectorAddressMap that corresponds to sector. The bit
// position is computed as ((Track & 0x7f) × 10) + (Sector − 1) +
// (sideBit × 800) − 40 — the −40 accounts for the directory tracks
// being outside the map's domain. Used by AddCodeFile to update the
// per-file map.
func (sector *Sector) SAMMask() (offset uint8, mask uint8) {
	bitOffset := (int(sector.Track)&0x7f)*10 + int(sector.Sector) - 1 + ((int(sector.Track)&0x80)>>7)*800 - 40
	return uint8(bitOffset >> 3), 1 << (bitOffset & 0x07)
}

// Offset returns the byte offset into a DiskImage at which sector's
// 512 bytes begin. The MGT layout is cylinder-interleaved: within
// each 10240-byte cylinder, side 0's 5120 bytes precede side 1's
// 5120 bytes.
func (sector *Sector) Offset() int {
	return int(sector.Track>>7)*5120 + (int(sector.Sector)-1)*512 + int(sector.Track&0x7f)*10240
}

// WriteSector copies sd's 512 bytes into the disk image at sector's
// location. Unlike SectorData, this performs no validation — a
// malformed Sector silently corrupts the image.
func (di *DiskImage) WriteSector(sector *Sector, sd *SectorData) {
	offset := sector.Offset()
	for i, b := range sd {
		di[i+offset] = b
	}
}

// Raw encodes fh into its 9-byte on-disk form (the bytes that
// prefix the file body). Bytes 5–6 are emitted as 0x00,0x00; a
// real ROM SAVE writes 0xFF,0xFF there, but both are spec-compliant
// ("unused" per the Tech Manual).
func (fh *FileHeader) Raw() [9]byte {
	return [9]byte{
		byte(fh.Type),
		byte(fh.LengthMod16K),
		byte(fh.LengthMod16K >> 8),
		byte(fh.PageOffset),
		byte(fh.PageOffset >> 8),
		fh.ExecutionAddressDiv16K,
		fh.ExecutionAddressMod16KLo,
		fh.Pages,
		fh.StartPage,
	}
}

// Raw returns the bytes that go to disk: the 9-byte FileHeader
// followed by the body.
func (file *File) Raw() []byte {
	h := file.Header.Raw()
	return append(h[:], file.Body...)
}
