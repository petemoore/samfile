package samfile

// DiskSubject wraps a DiskJournal + DiskImage for disk-scope rule
// evaluation. Single instance per Verify() run.
type DiskSubject struct {
	Journal *DiskJournal
	Disk    *DiskImage
	Dialect Dialect
}

func (s *DiskSubject) Ref() string { return "disk" }

func (s *DiskSubject) Attributes() map[string]any {
	used := len(s.Journal.UsedFileEntries())
	bootSig := false
	if sd, err := s.Disk.SectorData(&Sector{Track: 4, Sector: 1}); err == nil {
		// Match the ROM's masked compare (rom-disasm:20582-20598): bit 5
		// (case) and bit 7 ignored.
		expected := [4]byte{'B', 'O', 'O', 'T'}
		bootSig = true
		for i := 0; i < 4; i++ {
			if (sd[256+i]^expected[i])&0x5F != 0 {
				bootSig = false
				break
			}
		}
	}
	return map[string]any{
		"dialect":                s.Dialect.String(),
		"boot_signature_present": bootSig,
		"used_slot_count":        used,
	}
}
