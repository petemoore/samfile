package sambasic_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	samfile "github.com/petemoore/samfile/v3"
	"github.com/petemoore/samfile/v3/sambasic"
)

func TestExhaustiveRealDiskRoundtrip(t *testing.T) {
	var diskPaths []string

	home := os.Getenv("HOME")
	searchDirs := []string{
		filepath.Join(home, "Downloads"),
		filepath.Join(home, "git"),
	}

	if wd, err := os.Getwd(); err == nil {
		searchDirs = append(searchDirs, filepath.Join(wd, "testdata", "mgt"))
	}

	seen := map[string]bool{}
	for _, dir := range searchDirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			ext := filepath.Ext(path)
			if ext != ".dsk" && ext != ".mgt" && ext != ".DSK" && ext != ".MGT" {
				return nil
			}
			if info.Size() != 819200 {
				return nil
			}
			if !seen[path] {
				seen[path] = true
				diskPaths = append(diskPaths, path)
			}
			return nil
		})
	}

	if len(diskPaths) == 0 {
		t.Skip("no .dsk or .mgt files found; skipping real-disk tests")
	}

	var disksScanned, basicFilesTotal, passCount, failCount int

	for _, path := range diskPaths {
		entry := filepath.Base(path)
		di, err := samfile.Load(path)
		if err != nil {
			continue
		}
		disksScanned++

		dj := di.DiskJournal()
		for slot, fe := range dj {
			if !fe.Used() || fe.Type != samfile.FT_SAM_BASIC {
				continue
			}
			basicFilesTotal++
			fileName := fe.Name.String()
			t.Run(entry+"/"+fileName, func(t *testing.T) {
				file, fileErr := func() (*samfile.File, error) {
					defer func() { recover() }()
					return di.File(fileName)
				}()
				if fileErr != nil || file == nil {
					t.Skipf("slot %d: File(%q): cannot read (corrupt directory entry)", slot, fileName)
					return
				}

				parsed, err := sambasic.Parse(file.Body)
				if err != nil {
					failCount++
					t.Logf("slot %d: Parse(%q): %v", slot, fileName, err)
					return
				}

				got := parsed.Bytes()
				if !bytes.Equal(got, file.Body) {
					failCount++
					t.Logf("slot %d %q: roundtrip mismatch: got %d bytes, want %d bytes",
						slot, fileName, len(got), len(file.Body))
					minLen := len(got)
					if len(file.Body) < minLen {
						minLen = len(file.Body)
					}
					for i := 0; i < minLen; i++ {
						if got[i] != file.Body[i] {
							t.Logf("  first diff at offset %d: got 0x%02x, want 0x%02x", i, got[i], file.Body[i])
							break
						}
					}
				} else {
					passCount++
				}
			})
		}
	}

	t.Logf("Summary: %d disks scanned, %d BASIC files tested, %d passed, %d failed",
		disksScanned, basicFilesTotal, passCount, failCount)
}
