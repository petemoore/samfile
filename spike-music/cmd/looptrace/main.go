// Command looptrace reports an E-Tracker module's intro/loop structure by
// watching the engine's master order-list pointer (the self-modified operand at
// 0x8462). The pointer advances monotonically through the order list and, on the
// 0xFF end-marker, is reloaded from the loop point (0x84A5, set by an 0xFE
// marker) — jumping backwards. We time the first two backward jumps (wraps):
//
//	loopLen  = wrap2 - wrap1   (one loop-only pass)
//	introLen = wrap1 - loopLen (whatever precedes the loop point; 0 if the song
//	                            loops from the start)
//
// This is the authoritative loop signal — exact and aligned to emission. (An
// earlier read-address-fingerprint approach produced bogus "intros" because
// first-pass channel phase hasn't settled; the order pointer doesn't have that
// problem.)
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
func (s *sam) Get16(addr uint16) uint16 {
	return uint16(s.ram[s.page(addr)][addr&0x3FFF]) | uint16(s.ram[s.page(addr+1)][(addr+1)&0x3FFF])<<8
}
func (s *sam) In(port uint8) uint8 {
	switch port {
	case 0xFA:
		return s.lmpr
	case 0xFB:
		return s.hmpr
	case 0xF9:
		return 0x00
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
		modPath   = flag.String("mod", "", "module file")
		maxFrames = flag.Int("max-frames", 50*600, "frame cap")
	)
	flag.Parse()
	if *modPath == "" {
		fmt.Fprintln(os.Stderr, "usage: looptrace -mod FILE")
		os.Exit(2)
	}
	mod, err := os.ReadFile(*modPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	s := &sam{lmpr: 0x1F, hmpr: 0x02}
	copy(s.ram[2][:], mod)
	cpu := &z80.CPU{Memory: s, IO: s}
	const sentinel = 0xFFFE
	call := func(entry uint16) {
		cpu.PC = entry
		cpu.SP = 0xDFFE
		s.Set(cpu.SP-1, byte(sentinel>>8))
		s.Set(cpu.SP-2, byte(sentinel&0xFF))
		cpu.SP -= 2
		for n := 0; n < 5_000_000 && cpu.PC != sentinel && !cpu.HALT; n++ {
			cpu.Step()
		}
	}

	call(0x8000) // init

	// Order pointer (master sequence position) is the self-modified operand at
	// 0x8462 (the LD HL,nnnn at 0x8461). It advances through the order list and
	// is reloaded from the loop point (0x84A5) when an 0xFF end-marker is hit.
	order0 := s.Get16(0x8462)
	loopPt0 := s.Get16(0x84A5)

	// Detect the first TWO order-pointer wraps (0xFF end-markers). The first
	// wrap ends the whole first pass (intro + loop body); the second ends the
	// first loop-only pass. So loopLen = wrap2-wrap1, and introLen = wrap1-loopLen.
	var wraps []int
	prevOrder := order0
	for f := 0; f < *maxFrames && len(wraps) < 2; f++ {
		call(0x8006)
		op := s.Get16(0x8462)
		if f > 0 && op < prevOrder {
			wraps = append(wraps, f)
		}
		prevOrder = op
	}

	fmt.Printf("%s:\n", *modPath)
	fmt.Printf("  order-list start=0x%04X  loop-point(0x84A5)=0x%04X\n", order0, loopPt0)
	if len(wraps) < 2 {
		fmt.Printf("  fewer than 2 wraps found (%v)\n", wraps)
		return
	}
	w1, w2 := wraps[0], wraps[1]
	loopLen := w2 - w1
	introLen := w1 - loopLen
	fmt.Printf("  wrap1=%d wrap2=%d => intro=%d frames (%.2fs)  loop=%d frames (%.2fs)\n",
		w1, w2, introLen, float64(introLen)/50, loopLen, float64(loopLen)/50)
}
