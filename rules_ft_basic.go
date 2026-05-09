package samfile

import (
	"github.com/petemoore/samfile/v3/sambasic"
)

// §7 FT_SAM_BASIC rules (catalog docs/disk-validity-rules.md §7).
// Rules in this file check FT_SAM_BASIC invariants: FileTypeInfo triplets,
// VARS/gap sizes, program sentinel byte, line-number encoding, auto-RUN
// start-line validity, and MGTFlags convention. They apply to all dialects
// (BASIC-VARS-GAP-INVARIANT consults ctx.Dialect internally).

// bodyData reads the file body (excluding the 9-byte header) by
// walking fe's sector chain. Mirrors the chain-walk loop in
// (*DiskImage).File but without the filename-lookup wrapper, so
// callers that already have a *FileEntry don't re-iterate the
// directory. Returns ("body bytes", nil) on success or
// (nil, err) when a SectorData call fails — rules treat the error
// as "no finding" because Phase 3's §1/§3 rules already report the
// underlying chain problem.
//
// The returned slice is fe.Length() bytes long; it does NOT include
// the body-header bytes 0..8, matching the convention of samfile.File's
// Body field.
func bodyData(di *DiskImage, fe *FileEntry) ([]byte, error) {
	fileLength := fe.Length()
	raw := make([]byte, fileLength+9)
	sd, err := di.SectorData(fe.FirstSector)
	if err != nil {
		return nil, err
	}
	fp := sd.FilePart()
	i := uint16(0)
	for {
		copy(raw[510*i:], fp.Data[:])
		i++
		if i == fe.Sectors {
			break
		}
		sd, err = di.SectorData(fp.NextSector)
		if err != nil {
			return nil, err
		}
		fp = sd.FilePart()
	}
	return raw[9:], nil
}

// Ensure the sambasic package is imported by this file so it is
// available when the §7 rules are added in the next commit.
var _ = sambasic.REM
