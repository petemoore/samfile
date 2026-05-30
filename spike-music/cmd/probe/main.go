// Command probe runs a single E-Tracker module's own replay engine headlessly
// in a minimal SAM Coupé Z80 model and reports what it does each frame: the
// 26-byte SAA register shadow the engine maintains at 0x83D3-0x83EC.
//
// It is a reverse-engineering instrument (see FINDINGS.md), not the renderer.
// It confirms the "module is a self-contained player" finding: each module
// loads at 0x8000, its init entry is 0x8000 (-> jp 0x83EF), and its per-frame
// "play one tick" entry is 0x8006 (a flag-gated 6-channel processor that
// recomputes and flushes the SAA register shadow). cmd/render is the renderer.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/koron-go/z80"
)

const (
	pageSize = 16384
	numPages = 32
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
	// SAA writes (ports 0x1FF/0x00FF) are not snooped here; the engine's
	// 26-byte register shadow (read from memory below) is the authoritative
	// per-frame state, which sidesteps koron-go/z80's 8-bit port interface.
}

func main() {
	var (
		modPath = flag.String("mod", "", "module file (m01..m10)")
		frames  = flag.Int("frames", 60, "play frames to run")
	)
	flag.Parse()
	if *modPath == "" {
		fmt.Fprintln(os.Stderr, "usage: probe -mod FILE [-frames N]")
		os.Exit(2)
	}
	mod, err := os.ReadFile(*modPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	s := &sam{lmpr: 0x1F, hmpr: 0x02} // section C (0x8000) = physical page 2
	if len(mod) > pageSize {
		copy(s.ram[2][:], mod[:pageSize])
		copy(s.ram[3][:], mod[pageSize:])
	} else {
		copy(s.ram[2][:], mod)
	}
	cpu := &z80.CPU{Memory: s, IO: s}

	const sentinel = 0xFFFE
	run := func(entry uint16) {
		cpu.PC = entry
		cpu.SP = 0xDFFE
		s.Set(cpu.SP-1, byte(sentinel>>8))
		s.Set(cpu.SP-2, byte(sentinel&0xFF))
		cpu.SP -= 2
		for n := 0; n < 5_000_000 && cpu.PC != sentinel && !cpu.HALT; n++ {
			cpu.Step()
		}
	}

	run(0x8000) // init
	fmt.Printf("after-init shadow (regs 00..19): % 02X\n", readShadow(s))

	prev := readShadow(s)
	changes := 0
	for f := 0; f < *frames; f++ {
		run(0x8006) // play one tick
		cur := readShadow(s)
		if cur != prev {
			changes++
			fmt.Printf("frame %3d regs: % 02X\n", f, cur[:])
			prev = cur
		}
	}
	fmt.Printf("frames=%d changed=%d\n", *frames, changes)
}

func readShadow(s *sam) [26]byte {
	var sh [26]byte
	for i := range sh {
		sh[i] = s.Get(uint16(0x83D3 + i))
	}
	return sh
}
