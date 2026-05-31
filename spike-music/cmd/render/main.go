// Command render converts a SAM Coupé E-Tracker module to a WAV file (one or
// more loop iterations).
//
// Pipeline (see FINDINGS.md for the reverse-engineering that justifies it):
//
//	module loads at 0x8000  ->  run init entry 0x8000 (jp 0x83EF)
//	per frame (50 Hz):       ->  run play entry 0x8006 (6-channel tick)
//	                         ->  read the 26-byte SAA register shadow
//	                             the engine maintains at 0x83D3..0x83EC
//	feed shadow -> Go SAASound model -> 44.1 kHz stereo PCM -> WAV
//
// Loop detection (exact): the engine's master order-list pointer is the
// self-modified operand at 0x8462. It advances through the order list and, on
// the 0xFF end-marker, is reloaded from the loop point (0x84A5, set by the 0xFE
// marker which E-Tracker songs place at the very start) — so it jumps
// *backwards*. That backward jump is the loop boundary; songs loop from the
// start (no intro). We capture the SAA shadow each frame up to the wrap;
// [0, wrap) is exactly one seamless loop, repeated -loops times in the output.
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
	orderPtrAddr    = 0x8462 // self-modified operand: live master order-list pointer
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
func (s *sam) Get16(addr uint16) uint16 {
	return uint16(s.Get(addr)) | uint16(s.Get(addr+1))<<8
}
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
		loops     = flag.Int("loops", 4, "number of times to repeat the loop body in the output")
		loopsAbbr = flag.Int("l", 0, "alias for -loops (0 = use -loops)")
		maxFrames = flag.Int("max-frames", frameHz*600, "hard cap on frames if no loop wrap is found")
		verbose   = flag.Bool("v", false, "verbose")
		regDump   = flag.String("regdump", "", "also write per-frame SAA registers (regs 00..19) as hex CSV to this path")
	)
	flag.Parse()
	if *modPath == "" || *outPath == "" {
		fmt.Fprintln(os.Stderr, "usage: render -mod FILE -out FILE.wav [-max-frames N]")
		os.Exit(2)
	}
	if *loopsAbbr != 0 {
		*loops = *loopsAbbr // -l overrides -loops
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

	// Capture the intro + loop structure. The engine's master order-list pointer
	// (the self-modified operand at 0x8462) advances monotonically through the
	// order list; on the 0xFF end-marker it is reloaded from the loop point
	// (0x84A5, set by an 0xFE marker) and jumps *backwards*. We time the first
	// two backward jumps (wraps):
	//
	//	wrap1            = end of the first pass (intro + one loop body)
	//	loopFrames       = wrap2 - wrap1   (one loop-only pass)
	//	introFrames      = wrap1 - loopFrames (patterns before the loop point;
	//	                   0 when the 0xFE marker is at the start, as in 9/10 tunes)
	//
	// shadows[0:wrap1] is the whole first pass; shadows[introFrames:wrap1] is the
	// loop body. Output = intro once + loop body `*loops` times.
	type frame [shadowLen]byte
	var shadows []frame
	var wraps []int
	prevOrder := s.Get16(orderPtrAddr)
	for f := 0; f < *maxFrames && len(wraps) < 2; f++ {
		call(playEntry)
		var sh frame
		for i := 0; i < shadowLen; i++ {
			sh[i] = s.Get(uint16(shadowBase + i))
		}
		order := s.Get16(orderPtrAddr)
		if f > 0 && order < prevOrder {
			wraps = append(wraps, f)
		}
		prevOrder = order
		shadows = append(shadows, sh)
	}

	introFrames, loopFrames := 0, len(shadows)
	if len(wraps) >= 2 {
		wrap1, wrap2 := wraps[0], wraps[1]
		loopFrames = wrap2 - wrap1
		introFrames = wrap1 - loopFrames
		shadows = shadows[:wrap1] // intro + one loop body
	} else if *verbose {
		fmt.Printf("WARNING: <2 order-pointer wraps in %d frames; treating all as loop\n", *maxFrames)
	}
	if *verbose {
		fmt.Printf("intro=%d frames (%.2fs)  loop=%d frames (%.2fs)  x%d\n",
			introFrames, float64(introFrames)/frameHz,
			loopFrames, float64(loopFrames)/frameHz, *loops)
	}

	if *regDump != "" {
		var b []byte
		for f, sh := range shadows {
			b = fmt.Appendf(b, "%d", f)
			for _, r := range sh {
				b = fmt.Appendf(b, ",%02X", r)
			}
			b = append(b, '\n')
		}
		must(os.WriteFile(*regDump, b, 0o644))
	}

	// Render the intro once, then the loop body `*loops` times. The chip is fed
	// continuously (never reset), so repeats are phase-seamless.
	chip := saa.New(0, 0)             // SAM defaults: 8 MHz clock, 44100 Hz
	chip.WriteAddressData(0x1C, 0x01) // sound enable
	var pcm []int16
	buf := make([]int16, samplesPerFrame*2)
	emit := func(lo, hi int) {
		for f := lo; f < hi; f++ {
			sh := shadows[f]
			for r := 0; r < shadowLen; r++ {
				chip.WriteAddressData(byte(r), sh[r])
			}
			chip.GenerateMany(buf, samplesPerFrame)
			pcm = append(pcm, buf...)
		}
	}
	emit(0, introFrames) // intro, once
	for rep := 0; rep < *loops; rep++ {
		emit(introFrames, introFrames+loopFrames) // loop body
	}

	must(wav.WriteStereo16(*outPath, pcm, sampleRate))
	fmt.Printf("%s: intro=%d + loop=%d frames x%d -> %s (%.2fs)\n",
		*modPath, introFrames, loopFrames, *loops, *outPath, float64(len(pcm)/2)/float64(sampleRate))
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
