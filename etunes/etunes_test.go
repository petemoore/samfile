package etunes

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"
)

func read(t *testing.T, p string) []byte {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestIsETrackerModule(t *testing.T) {
	if !IsETrackerModule(read(t, "testdata/m01")) {
		t.Error("m01 should be detected")
	}
	if IsETrackerModule([]byte("not a module, just some text bytes here.....")) {
		t.Error("garbage should not be detected")
	}
}

func TestLoopSplit(t *testing.T) {
	_, intro01, loop01, err := Frames(read(t, "testdata/m01"))
	if err != nil {
		t.Fatal(err)
	}
	if intro01 != 0 || loop01 != 1760 {
		t.Errorf("m01: intro=%d loop=%d, want 0/1760", intro01, loop01)
	}
	_, intro03, loop03, err := Frames(read(t, "testdata/m03"))
	if err != nil {
		t.Fatal(err)
	}
	if intro03 != 288 || loop03 != 3672 {
		t.Errorf("m03: intro=%d loop=%d, want 288/3672", intro03, loop03)
	}
}

func TestOracle(t *testing.T) {
	shadows, _, _, err := Frames(read(t, "testdata/m01"))
	if err != nil {
		t.Fatal(err)
	}
	oracle := readReg(t, "testdata/m01.reg")
	// The oracle file records the SAA writes performed during each frame tick;
	// readReg builds out[N] as the chip state *before* frame N's writes, so
	// out[N+1] holds the state *after* frame N's writes.  Frame 0 in the oracle
	// is the init call (only touches reg 28, outside the 26-byte shadow), and
	// frame 1 is the first play tick — so out[2] is the first audible snapshot.
	// Our shadows[0] is also captured after the first play tick, hence the
	// oracle alignment is f+2 (not f+1 as one might first assume).
	const oracleOffset = 2
	n := len(shadows)
	if len(oracle)-oracleOffset < n {
		n = len(oracle) - oracleOffset
	}
	for f := 0; f < n; f++ {
		got, want := audibleSig(shadows[f]), audibleSig(oracle[f+oracleOffset])
		if got != want {
			t.Fatalf("audible mismatch at frame %d:\n got %v\nwant %v", f, got, want)
		}
	}
}

func audibleSig(r [shadowLen]byte) [19]int {
	var sig [19]int
	for ch := 0; ch < 6; ch++ {
		sig[ch] = int(r[ch])
	}
	for ch := 0; ch < 6; ch++ {
		if r[ch] != 0 {
			sig[6+2*ch] = int(r[8+ch])
			op := r[0x10+ch/2]
			if ch%2 == 1 {
				sig[6+2*ch+1] = int(op >> 4)
			} else {
				sig[6+2*ch+1] = int(op & 0x0F)
			}
		} else {
			sig[6+2*ch], sig[6+2*ch+1] = -1, -1
		}
	}
	sig[18] = int(r[0x14])<<16 | int(r[0x15])<<8 | int(r[0x16])
	return sig
}

func readReg(t *testing.T, p string) [][shadowLen]byte {
	t.Helper()
	f, err := os.Open(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var state [32]byte
	var out [][shadowLen]byte
	sc := bufio.NewScanner(f)
	cur := -1
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		fr, _ := strconv.Atoi(fields[0])
		for fr > cur {
			var snap [shadowLen]byte
			copy(snap[:], state[:shadowLen])
			out = append(out, snap)
			cur++
		}
		if fields[1] == "DATA" {
			reg, _ := strconv.Atoi(fields[2])
			val, _ := strconv.Atoi(fields[3])
			if reg < len(state) {
				state[reg] = byte(val)
			}
		}
	}
	var snap [shadowLen]byte
	copy(snap[:], state[:shadowLen])
	out = append(out, snap)
	return out
}
