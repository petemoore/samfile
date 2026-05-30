// Command render converts a SAM Coupé E-Tracker module to a WAV file.
//
// Pipeline (see FINDINGS.md for the reverse-engineering that justifies it):
//
//	module loads at 0x8000  ->  run init entry 0x8000 (jp 0x83EF)
//	per frame (50 Hz):       ->  run play entry 0x8006 (6-channel tick)
//	                         ->  read the 26-byte SAA register shadow
//	                             the engine maintains at 0x83D3..0x83EC
//	feed shadow -> Go SAASound model -> 44.1 kHz stereo PCM -> WAV
//
// Loop detection: the replay engine is deterministic, so its full mutable
// state (physical page 2 = logical 0x8000..0xBFFF) repeating exactly means the
// song has looped. We render from frame 0 up to and including the frame whose
// post-state first repeats an earlier frame's state (= intro + one full loop),
// then trim trailing silence.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/koron-go/z80"
	"github.com/petemoore/samfile-spike-decode/spike-music/saa"
	"github.com/petemoore/samfile-spike-decode/spike-music/wav"
)

const (
	pageSize        = 16384
	numPages        = 32
	playEntry       = 0x8006
	initEntry       = 0x8000
	shadowBase      = 0x83D3 // regs 0x00..0x19 live here (26 bytes)
	shadowLen       = 26
	frameHz         = 50
	sampleRate      = 44100
	samplesPerFrame = sampleRate / frameHz // 882
)

type sam struct {
	ram        [numPages][pageSize]byte
	lmpr, hmpr uint8
}

func (s *sam) page(addr uint16) int {
	switch addr >> 14 {
	case 0:
		return int(s.lmpr & 0x1F)
	case 1:
		return int((s.lmpr&0x1F + 1) & 0x1F)
	case 2:
		return int(s.hmpr & 0x1F)
	default:
		return int((s.hmpr&0x1F + 1) & 0x1F)
	}
}

func (s *sam) Get(addr uint16) uint8    { return s.ram[s.page(addr)][addr&0x3FFF] }
func (s *sam) Set(addr uint16, v uint8) { s.ram[s.page(addr)][addr&0x3FFF] = v }
func (s *sam) In(port uint8) uint8 {
	switch port {
	case 0xFA:
		return s.lmpr
	case 0xFB:
		return s.hmpr
	case 0xF9:
		return 0x00 // status/interrupt port: nothing pending
	}
	return 0xFF
}
func (s *sam) Out(port uint8, v uint8) {
	switch port {
	case 0xFA:
		s.lmpr = v
	case 0xFB:
		s.hmpr = v
	}
}

func main() {
	var (
		modPath   = flag.String("mod", "", "module file (m01..m10)")
		outPath   = flag.String("out", "", "output WAV path")
		maxFrames = flag.Int("max-frames", frameHz*360, "hard cap on frames (default 360s; loop search covers up to half this)")
		verbose   = flag.Bool("v", false, "verbose")
		regDump   = flag.String("regdump", "", "also write per-frame SAA registers (regs 00..19) as hex CSV to this path")
	)
	flag.Parse()
	if *modPath == "" || *outPath == "" {
		fmt.Fprintln(os.Stderr, "usage: render -mod FILE -out FILE.wav [-max-frames N]")
		os.Exit(2)
	}
	mod, err := os.ReadFile(*modPath)
	must(err)

	s := &sam{lmpr: 0x1F, hmpr: 0x02}
	if len(mod) > pageSize {
		copy(s.ram[2][:], mod[:pageSize])
		copy(s.ram[3][:], mod[pageSize:])
	} else {
		copy(s.ram[2][:], mod)
	}
	cpu := &z80.CPU{Memory: s, IO: s}

	const sentinel = 0xFFFE
	call := func(entry uint16) {
		cpu.PC = entry
		cpu.SP = 0xDFFE
		s.Set(cpu.SP-1, byte(sentinel>>8))
		s.Set(cpu.SP-2, byte(sentinel&0xFF))
		cpu.SP -= 2
		for n := 0; n < 5_000_000; n++ {
			if cpu.PC == sentinel || cpu.HALT {
				return
			}
			cpu.Step()
		}
	}

	call(initEntry)

	// Run frames, capturing the SAA register shadow each frame.
	type frame [shadowLen]byte
	var shadows []frame
	for f := 0; f < *maxFrames; f++ {
		call(playEntry)
		var sh frame
		for i := 0; i < shadowLen; i++ {
			sh[i] = s.Get(uint16(shadowBase + i))
		}
		shadows = append(shadows, sh)
	}

	if *regDump != "" {
		var b []byte
		for f, sh := range shadows {
			b = append(b, []byte(fmt.Sprintf("%d", f))...)
			for _, r := range sh {
				b = append(b, []byte(fmt.Sprintf(",%02X", r))...)
			}
			b = append(b, '\n')
		}
		must(os.WriteFile(*regDump, b, 0o644))
	}

	// Loop detection on the SAA shadow *sequence*. The replay engine is
	// deterministic, so once past any intro the register stream is periodic
	// with period = the musical loop length. (We detect on the audible
	// register stream rather than full machine state because the engine keeps
	// a free-running frame/LFO counter that never lets full state repeat
	// exactly even though the music loops.) Find the smallest period whose
	// repeat holds over a long verification window at the tail, then find the
	// earliest frame from which that period holds continuously to the end, and
	// render [0, loopStart+period) = intro + exactly one loop.
	renderN := len(shadows)
	if p, loopStart := detectLoop(shadows); p > 0 {
		renderN = loopStart + p
		if *verbose {
			fmt.Printf("loop detected: intro %d frames, loop %d frames (%.2fs); rendering %d frames (%.2fs)\n",
				loopStart, p, float64(p)/frameHz, renderN, float64(renderN)/frameHz)
		}
	} else if *verbose {
		fmt.Printf("no loop within %d frames; rendering all + trimming silence\n", *maxFrames)
	}

	// Render: feed each frame's 26 registers to the SAA, enable sound once,
	// then generate one frame of samples.
	chip := saa.New(0, 0)             // SAM defaults: 8 MHz clock, 44100 Hz
	chip.WriteAddressData(0x1C, 0x01) // sound enable
	pcm := make([]int16, 0, renderN*samplesPerFrame*2)
	buf := make([]int16, samplesPerFrame*2)
	for f := 0; f < renderN; f++ {
		sh := shadows[f]
		for r := 0; r < shadowLen; r++ {
			chip.WriteAddressData(byte(r), sh[r])
		}
		chip.GenerateMany(buf, samplesPerFrame)
		pcm = append(pcm, buf...)
	}

	pcm = trimTrailingSilence(pcm)
	must(wav.WriteStereo16(*outPath, pcm, sampleRate))
	fmt.Printf("%s: %d frames -> %s (%.2fs)\n",
		*modPath, renderN, *outPath, float64(len(pcm)/2)/float64(sampleRate))
}

// detectLoop finds the musical loop in a shadow-register sequence. It returns
// (period, loopStart) such that shadows[f] == shadows[f+period] for all
// f in [loopStart, len-period), with the smallest such period — or (0,0) if no
// stable loop is found. A period is accepted only if it repeats over a
// verification window of at least minWindow frames at the tail (guards against
// short coincidental matches).
func detectLoop[T comparable](shadows []T) (period, loopStart int) {
	n := len(shadows)
	if n < 100 {
		return 0, 0
	}
	const minWindow = 400 // ≥8s of identical repetition to accept a period
	maxPeriod := n / 2
	for p := 1; p <= maxPeriod; p++ {
		// Verify p over a window at the tail.
		w := minWindow
		if w > n-p {
			w = n - p
		}
		ok := true
		for i := 0; i < w; i++ {
			if shadows[n-1-i] != shadows[n-1-i-p] {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		// Found a tail period p. Walk back to the earliest frame from which
		// the period holds continuously to the end.
		ls := n - p
		for ls > 0 && shadows[ls-1] == shadows[ls-1+p] {
			ls--
		}
		return p, ls
	}
	return 0, 0
}

// trimTrailingSilence removes trailing stereo frames whose absolute amplitude
// stays under a small threshold, leaving a short tail.
func trimTrailingSilence(pcm []int16) []int16 {
	const thresh = 64
	n := len(pcm) / 2
	last := 0
	for i := 0; i < n; i++ {
		l, r := pcm[2*i], pcm[2*i+1]
		if abs16(l) > thresh || abs16(r) > thresh {
			last = i
		}
	}
	end := (last + frameHz/2) * 2 // keep ~0.5s tail
	if end > len(pcm) {
		end = len(pcm)
	}
	return pcm[:end]
}

func abs16(v int16) int16 {
	if v < 0 {
		return -v
	}
	return v
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
