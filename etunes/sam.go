// Package etunes decodes SAM Coupé E-Tracker music modules to PCM by running
// the module's own Z80 replay engine headlessly on a minimal SAM model and
// capturing the per-frame SAA1099 register shadow the engine maintains.
package etunes

import "github.com/koron-go/z80"

const (
	pageSize = 16384
	numPages = 32
)

// sam is a minimal SAM Coupé memory+IO model: 32 × 16 KB RAM pages with
// LMPR (port 0xFA) / HMPR (port 0xFB) paging. Enough to run a self-contained
// E-Tracker module loaded at logical 0x8000.
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
		return 0x00 // status/interrupt port: nothing pending (no SPACE keypress)
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

// newCPU returns a CPU wired to a sam with the module loaded at 0x8000.
func newCPU(mod []byte) (*z80.CPU, *sam) {
	s := &sam{lmpr: 0x1F, hmpr: 0x02}
	copy(s.ram[2][:], mod) // all known modules are < 16 KB
	cpu := &z80.CPU{Memory: s, IO: s}
	return cpu, s
}

// call runs the routine at entry until it RETs to a sentinel or halts.
func call(cpu *z80.CPU, s *sam, entry uint16) {
	const sentinel = 0xFFFE
	cpu.PC = entry
	cpu.SP = 0xDFFE
	s.Set(cpu.SP-1, byte(sentinel>>8))
	s.Set(cpu.SP-2, byte(sentinel&0xFF))
	cpu.SP -= 2
	for n := 0; n < 5_000_000 && cpu.PC != sentinel && !cpu.HALT; n++ {
		cpu.Step()
	}
}
