package samfile

import "fmt"

// ChainStepSubject wraps one step in a file's sector chain. Track
// and Sector identify the sector itself; ChainIndex is the 0..N-1
// position within this slot's chain.
type ChainStepSubject struct {
	SlotIndex   int
	ChainIndex  int
	Track       uint8
	Sector      uint8
	NextTrack   uint8
	NextSector  uint8
	Position    string // first | intermediate | last | orphan
	OnSAMMap    bool
	OnDirSAMMap bool
	Disk        *DiskImage
	Journal     *DiskJournal
}

func (s *ChainStepSubject) Ref() string {
	return fmt.Sprintf("slot=%d,chain=%d,track=%d,sector=%d", s.SlotIndex, s.ChainIndex, s.Track, s.Sector)
}

func (s *ChainStepSubject) Attributes() map[string]any {
	side := 0
	if s.Track&0x80 != 0 {
		side = 1
	}
	// Distance from dir tracks: dir is tracks 0..3 of side 0.
	dirDistance := 999
	t := int(s.Track & 0x7F)
	if side == 0 {
		if t >= 4 {
			dirDistance = t - 3
		} else {
			dirDistance = 0
		}
	} else {
		dirDistance = t + 1
	}
	return map[string]any{
		"slot_index":               s.SlotIndex,
		"chain_index":              s.ChainIndex,
		"chain_position":           s.Position,
		"track":                    int(s.Track),
		"sector":                   int(s.Sector),
		"side":                     side,
		"next_track":               int(s.NextTrack),
		"next_sector":              int(s.NextSector),
		"on_sam_map":               s.OnSAMMap,
		"on_dir_sam_map":           s.OnDirSAMMap,
		"distance_from_dir_tracks": dirDistance,
	}
}
