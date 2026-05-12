package samfile

import (
	"encoding/hex"
	"fmt"
)

// SlotSubject wraps one directory slot's FileEntry for slot-scope
// rule evaluation. SlotIndex is the 0..79 position; FileEntry is
// the parsed dir entry. For erased slots (Type==0) the FileEntry
// still exists but most fields are zero.
type SlotSubject struct {
	SlotIndex int
	FileEntry *FileEntry
	Disk      *DiskImage
	Journal   *DiskJournal
}

func (s *SlotSubject) Ref() string { return fmt.Sprintf("slot=%d", s.SlotIndex) }

func (s *SlotSubject) Attributes() map[string]any {
	fe := s.FileEntry
	if fe == nil {
		return map[string]any{
			"slot_index":     s.SlotIndex,
			"slot_is_erased": true,
		}
	}
	pageOffsetForm := "other"
	switch fe.StartAddressPageOffset & 0xC000 {
	case 0x8000:
		pageOffsetForm = "0x8000"
	case 0x4000:
		pageOffsetForm = "0x4000"
	case 0x0000:
		pageOffsetForm = "0x0000"
	case 0xC000:
		pageOffsetForm = "0xC000"
	}
	dirMirrorPopulated := false
	for _, b := range fe.MGTFutureAndPast[1:10] {
		if b != 0 {
			dirMirrorPopulated = true
			break
		}
	}
	hasAutoRunOrAutoExec := fe.ExecutionAddressDiv16K != 0xFF
	ftrack, fsector, fside := 0, 0, 0
	if fe.FirstSector != nil {
		ftrack = int(fe.FirstSector.Track)
		fsector = int(fe.FirstSector.Sector)
		fside = int(fe.FirstSector.Track >> 7)
	}
	return map[string]any{
		"slot_index":              s.SlotIndex,
		"filename":                fe.Name.String(),
		"file_type":               fe.Type.String(),
		"file_type_byte":          int(fe.Type),
		"file_length":             int(fe.LengthMod16K),
		"page_offset_form":        pageOffsetForm,
		"pages":                   int(fe.Pages),
		"mgt_flags":               int(fe.MGTFlags),
		"first_track":             ftrack,
		"first_sector":            fsector,
		"first_side":              fside,
		"has_autorun_or_autoexec": hasAutoRunOrAutoExec,
		"dir_mirror_populated":    dirMirrorPopulated,
		"slot_is_erased":          fe.Type == 0,
		"file_type_info_hex":      hex.EncodeToString(fe.FileTypeInfo[:]),
		"sectors_count":           int(fe.Sectors),
	}
}
