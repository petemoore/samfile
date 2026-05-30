// etunes converts FRED E-Tracker tunes (m01..m10) to WAV by running Andrew
// Collier's real E-Code player headlessly on a Go Z80 + Go SAA1099, capturing
// the SAA register writes per 50 Hz frame.
package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"sort"

	"spike-emulate/saa"
	"spike-emulate/sam"
)

const (
	sampleRate     = 44100
	framesPerSec   = 50
	samplesPerFrame = sampleRate / framesPerSec // 882
)

func main() {
	assets := flag.String("assets", "/Users/pmoore/git/sam-jukebox/assets", "assets dir")
	module := flag.String("m", "m01", "module name (m01..m10)")
	frames := flag.Int("frames", 3000, "number of 50Hz frames to render (3000 = 60s)")
	budget := flag.Uint64("budget", 12000, "Z80 instructions per frame")
	trim := flag.Bool("trim", true, "trim leading silence and fade out the tail")
	wavOut := flag.String("wav", "", "WAV output path (empty = skip)")
	regOut := flag.String("reg", "", "register-stream dump path (empty = skip)")
	histogram := flag.Bool("hist", false, "print OUT port histogram and exit")
	ptrTrace := flag.Bool("ptr", false, "trace song pointer &02B2 per frame and exit")
	flag.Parse()

	m := buildMachine(*assets, *module)

	// Wire the Go SAA to the captured writes.
	chip := saa.New()
	m.SetSAAHandler(func(isAddr bool, val uint8) {
		if isAddr {
			chip.WriteAddress(val)
		} else {
			chip.WriteData(val)
		}
	})
	m.EnableSAACapture(true)

	// Run init (no interrupts) until the player reaches its main wait loop,
	// then drive frames. The player's main loop sits around &026F..&0286.
	// We give init a generous instruction budget to relocate, build tables and
	// call the module init.
	runPlayerInit(m)

	if *ptrTrace {
		prev := uint16(0)
		for f := 0; f < *frames; f++ {
			m.StepFrame(*budget)
			p := m.PeekPhysWord(1, 0x2B2) // song note pointer
			if p < prev {
				fmt.Printf("frame %4d ptr=%04X  <-- wrap\n", f, p)
			}
			if f%200 == 0 {
				fmt.Printf("frame %4d ptr(2B2)=%04X  2B5=%04X 2B4=%02X 2B7=%02X\n",
					f, p, m.PeekPhysWord(1, 0x2B5), m.PeekPhys(1, 0x2B4), m.PeekPhys(1, 0x2B7))
			}
			prev = p
		}
		return
	}

	var pcm []byte
	for f := 0; f < *frames; f++ {
		halted := m.StepFrame(*budget)
		if *wavOut != "" {
			// render this frame's audio AFTER its writes have been applied
			buf := make([]byte, samplesPerFrame*4)
			chip.GenerateMany(buf, samplesPerFrame)
			pcm = append(pcm, buf...)
		}
		if halted {
			fmt.Fprintf(os.Stderr, "player HALTed at frame %d\n", f)
			break
		}
	}

	if *histogram {
		printHistogram(m)
		fmt.Printf("total SAA writes: %d\n", len(m.Writes))
		return
	}

	fmt.Printf("module=%s frames=%d SAA writes=%d steps=%d\n",
		*module, *frames, len(m.Writes), m.Steps)
	printHistogram(m)

	if *regOut != "" {
		writeRegStream(*regOut, m.Writes)
		fmt.Printf("wrote register stream: %s\n", *regOut)
	}
	if *wavOut != "" {
		if *trim {
			pcm = trimAndFade(pcm)
		}
		writeWAV(*wavOut, pcm)
		fmt.Printf("wrote WAV: %s (%d samples, %.1fs)\n", *wavOut, len(pcm)/4, float64(len(pcm)/4)/sampleRate)
	}
}

// trimAndFade removes leading near-silence and applies a 1s linear fade-out so
// the (non-loop-trimmed) excerpt ends cleanly.
func trimAndFade(pcm []byte) []byte {
	const thresh = 200
	nSamp := len(pcm) / 4
	sample := func(i int) (int16, int16) {
		l := int16(uint16(pcm[i*4]) | uint16(pcm[i*4+1])<<8)
		r := int16(uint16(pcm[i*4+2]) | uint16(pcm[i*4+3])<<8)
		return l, r
	}
	start := 0
	for ; start < nSamp; start++ {
		l, r := sample(start)
		if abs16(l) > thresh || abs16(r) > thresh {
			break
		}
	}
	// keep 20ms pre-roll
	preroll := sampleRate / 50
	if start > preroll {
		start -= preroll
	} else {
		start = 0
	}
	out := append([]byte(nil), pcm[start*4:]...)

	// 1s fade-out
	nOut := len(out) / 4
	fade := sampleRate
	if fade > nOut {
		fade = nOut
	}
	for i := 0; i < fade; i++ {
		idx := nOut - fade + i
		g := float64(fade-i) / float64(fade)
		for c := 0; c < 2; c++ {
			off := idx*4 + c*2
			v := int16(uint16(out[off]) | uint16(out[off+1])<<8)
			v = int16(float64(v) * g)
			out[off] = byte(v)
			out[off+1] = byte(v >> 8)
		}
	}
	return out
}

func abs16(v int16) int16 {
	if v < 0 {
		return -v
	}
	return v
}

// buildMachine reproduces jukebox.bas lines 60..330 for a single tune.
func buildMachine(assets, module string) *sam.Machine {
	decr := mustRead(assets + "/decruncher")
	ecode := mustRead(assets + "/E-Code")
	freqtab := mustRead("out/freqtable.bin")
	etext := readOpt(assets + "/E-text")
	mod := mustRead(assets + "/" + module)

	m := sam.New()

	// --- decrunch E-Code (jukebox.bas line 60) ---
	m.LoadLinear(16384, decr)  // LOAD "decruncher" CODE  -> page 0
	m.LoadLinear(98304, ecode) // LOAD "E-Code" CODE 98304 -> page 5
	const mc = 16384
	m.Poke(mc+3, byte(32768/16384-1))
	m.DPoke(mc+4, uint16(32768%16384+32768))
	m.Poke(mc+6, byte(98304/16384-1))
	m.DPoke(mc+7, uint16(98304%16384+32768))
	m.Poke(mc+9, 41)
	m.Poke(mc+10, 44)
	m.Poke(mc+11, 53)
	m.SetLMPR(0x1F)
	m.LoadPage(0, 0x3EFE, []byte{0xFF, 0xFF})
	cpu := m.CPU()
	cpu.SP = 0x7EFE
	cpu.PC = 0x4000
	if r, _ := m.RunUntilPC(0xFFFF, 20_000_000); r != "pc" {
		fmt.Fprintf(os.Stderr, "warning: decrunch ended with %q\n", r)
	}

	// --- freq table poke (line 70): linear 33486 -> page1 +718 ---
	for i, b := range freqtab {
		m.Poke(33486+i, b)
	}

	// --- scroll text (line 280): E-text -> linear 33666 (page1 +898) ---
	if etext != nil {
		m.LoadLinear(33666, etext)
	}

	// --- module (line 320): LOAD a$ CODE 65536 -> page 3 ---
	m.LoadLinear(65536, mod)
	// scroller pointer DPOKE 33458 = 49152 (page1 +690)
	m.DPoke(33458, 49152)

	return m
}

// runPlayerInit sets up registers as BASIC's CALL 32768 would and runs the
// player's init path up to (but not entering) its frame loop.
func runPlayerInit(m *sam.Machine) {
	cpu := m.CPU()
	// CALL 32768: HMPR = (32768/16384 - 1) = 1, PC = &8000.
	m.SetHMPR(1)
	m.SetLMPR(0x1F)
	cpu.SP = 0xBF00
	cpu.PC = 0x8000
	cpu.IFF1 = false
	cpu.IM = 1
	// Run init until the player reaches its main wait loop (&026F). It sets
	// up IM1, builds tables, calls the module init, then EI and loops.
	r, n := m.RunUntilPC(0x026F, 5_000_000)
	fmt.Fprintf(os.Stderr, "init: reason=%s steps=%d PC=%04X HMPR=%02X LMPR=%02X IFF1=%v\n",
		r, n, cpu.PC, m.HMPR(), m.LMPR(), cpu.IFF1)
}

func printHistogram(m *sam.Machine) {
	uo := m.UnknownOuts()
	if len(uo) == 0 {
		fmt.Println("unknown OUT ports: none")
		return
	}
	type kv struct {
		port  uint8
		count int
	}
	var list []kv
	for p, c := range uo {
		list = append(list, kv{p, c})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].count > list[j].count })
	fmt.Println("unknown OUT ports (port: count):")
	for _, e := range list {
		fmt.Printf("  %02X: %d\n", e.port, e.count)
	}
}

func writeRegStream(path string, w []sam.SAAWrite) {
	f, err := os.Create(path)
	must(err)
	defer f.Close()
	fmt.Fprintln(f, "# frame\tkind\treg\tvalue   (E-Tracker SAA register stream — validation oracle for spike A)")
	for _, e := range w {
		if e.IsAddr {
			fmt.Fprintf(f, "%d\tADDR\t%d\t-\n", e.Frame, e.Val)
		} else {
			fmt.Fprintf(f, "%d\tDATA\t%d\t%d\n", e.Frame, e.Reg, e.Val)
		}
	}
}

func writeWAV(path string, pcm []byte) {
	f, err := os.Create(path)
	must(err)
	defer f.Close()
	const channels, bits = 2, 16
	byteRate := sampleRate * channels * bits / 8
	blockAlign := channels * bits / 8
	dataLen := len(pcm)
	w := func(v any) { binary.Write(f, binary.LittleEndian, v) }
	f.WriteString("RIFF")
	w(uint32(36 + dataLen))
	f.WriteString("WAVE")
	f.WriteString("fmt ")
	w(uint32(16))
	w(uint16(1)) // PCM
	w(uint16(channels))
	w(uint32(sampleRate))
	w(uint32(byteRate))
	w(uint16(blockAlign))
	w(uint16(bits))
	f.WriteString("data")
	w(uint32(dataLen))
	f.Write(pcm)
}

func mustRead(p string) []byte {
	b, err := os.ReadFile(p)
	must(err)
	return b
}

func readOpt(p string) []byte {
	b, err := os.ReadFile(p)
	if err != nil {
		return nil
	}
	return b
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
