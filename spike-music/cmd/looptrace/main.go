// Command looptrace finds the EXACT musical loop of an E-Tracker module by
// watching which song-data addresses the replay engine reads each frame.
//
// Idea (Pete's): the engine reads sequentially through the song-data region
// (>= 0x84B3, the order/pattern/instrument streams that follow the shared
// engine code) and, at the end of the song, jumps its read pointers back to the
// start. So the per-frame *set of song-data addresses read* is a fingerprint of
// the song position, and that fingerprint sequence is exactly periodic with the
// musical loop — cleanly, because it ignores the engine's free-running frame
// counter (which lives below 0x84B3 and never lets full machine state repeat).
//
// We hash, per frame, the distinct addresses read in [0x84B3, 0xC000), then
// find the smallest period P that the fingerprint sequence repeats over a long
// tail window, and the earliest frame from which P holds to the end. That gives
// (intro length, loop length) precisely.
package main

import (
	"flag"
	"fmt"
	"hash/fnv"
	"os"
	"sort"

	"github.com/koron-go/z80"
)

const (
	pageSize     = 16384
	numPages     = 32
	songDataLow  = 0x84B3
	songDataHigh = 0xC000 // section C only; stack/section D is >= 0xC000
)

type sam struct {
	ram        [numPages][pageSize]byte
	lmpr, hmpr uint8

	tracing  bool
	frameSet map[uint16]struct{} // distinct song-data reads this frame
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

func (s *sam) Get(addr uint16) uint8 {
	if s.tracing && addr >= songDataLow && addr < songDataHigh {
		s.frameSet[addr] = struct{}{}
	}
	return s.ram[s.page(addr)][addr&0x3FFF]
}
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

	s := &sam{lmpr: 0x1F, hmpr: 0x02, frameSet: map[uint16]struct{}{}}
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

	fps := []uint64{} // per-frame read fingerprints
	orders := []uint16{}
	s.tracing = true
	wrapFrame := -1
	var prevOrder uint16 = order0
	for f := 0; f < *maxFrames; f++ {
		clear(s.frameSet)
		call(0x8006)
		fps = append(fps, fingerprint(s.frameSet))
		op := s.Get16(0x8462)
		orders = append(orders, op)
		if wrapFrame < 0 && f > 0 && op < prevOrder {
			wrapFrame = f // order pointer jumped backwards = song wrapped
		}
		prevOrder = op
	}

	p, ls := detectLoop(fps)
	fmt.Printf("%s:\n", *modPath)
	fmt.Printf("  order-list start=0x%04X  loop-point(0x84A5) after init=0x%04X  (equal => loops from start)\n",
		order0, loopPt0)
	loopPtFinal := s.Get16(0x84A5)
	fmt.Printf("  loop-point after playing=0x%04X\n", loopPtFinal)
	if wrapFrame >= 0 {
		fmt.Printf("  ORDER-POINTER wrap at frame %d => loop length %d frames (%.2fs)\n",
			wrapFrame, wrapFrame, float64(wrapFrame)/50)
	} else {
		fmt.Printf("  no order-pointer wrap in %d frames\n", *maxFrames)
	}
	if p > 0 {
		fmt.Printf("  read-fingerprint: intro=%d loop=%d (%.2fs)\n", ls, p, float64(p)/50)
	}
}

func fingerprint(set map[uint16]struct{}) uint64 {
	addrs := make([]int, 0, len(set))
	for a := range set {
		addrs = append(addrs, int(a))
	}
	sort.Ints(addrs)
	h := fnv.New64a()
	var b [2]byte
	for _, a := range addrs {
		b[0], b[1] = byte(a), byte(a>>8)
		h.Write(b[:])
	}
	return h.Sum64()
}

func detectLoop(fp []uint64) (period, loopStart int) {
	n := len(fp)
	if n < 50 {
		return 0, 0
	}
	const minWindow = 200
	for p := 1; p <= n/2; p++ {
		w := minWindow
		if w > n-p {
			w = n - p
		}
		ok := true
		for i := 0; i < w; i++ {
			if fp[n-1-i] != fp[n-1-i-p] {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		ls := n - p
		for ls > 0 && fp[ls-1] == fp[ls-1+p] {
			ls--
		}
		return p, ls
	}
	return 0, 0
}
