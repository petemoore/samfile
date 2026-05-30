// Package sam is a minimal headless SAM Coupé machine built on koron-go/z80.
// It models exactly the subset needed to run the FRED E-music loader: 512 KB
// of paged RAM (LMPR/HMPR), the SAA1099 OUT decode, a 50 Hz frame interrupt,
// and helpers that mirror SAM BASIC's "long address" LOAD CODE / POKE / CALL.
//
// The paging model (resolveRead/resolveWritePage) is adapted from
// sam-aarch64/tools/z80-test-harness-go/harness.go, which sources it from the
// SAM Coupé Technical Manual v3.0 §6.10.
package sam

import (
	"github.com/koron-go/z80"
)

const (
	PageSize = 16384
	NumPages = 32

	portLMPR   = 0xFA
	portHMPR   = 0xFB
	portVMPR   = 0xFC
	portStatus = 0xF9 // read STATUS / write LINE
	portKbd    = 0xFE // keyboard + border

	statusIntLine  = 0x01
	statusIntFrame = 0x08
	statusKeyMask  = 0xE0

	// SAA1099 ports (SimCoupé Base/SAMIO.h): address register when the full
	// 9-bit port == 0x1FF (A8 high), data when port == 0x0FF. Low byte is
	// 0xFF for both, so we recover A8 from the B register at OUT time.
	saaLowByte = 0xFF
)

// SAAWrite records one write routed to the SAA, tagged with the frame in
// which it occurred. IsAddr distinguishes a register-select (OUT 511) from a
// data write (OUT 255). For data writes, Reg is the register that was
// selected by the most recent address write.
type SAAWrite struct {
	Frame int
	IsAddr bool
	Reg   uint8 // valid for data writes: the currently-selected register
	Val   uint8
}

// Machine is the headless SAM.
type Machine struct {
	ram  [NumPages][PageSize]byte
	lmpr uint8
	hmpr uint8

	cpu *z80.CPU

	// SAA routing
	saaWrite   func(isAddr bool, val uint8)
	curSAAReg  uint8
	frame      int
	Writes     []SAAWrite
	captureSAA bool

	// device port state
	vmpr         uint8
	framePending bool

	// diagnostics
	Steps      uint64
	unknownOut map[uint8]int
}

// New builds a machine with cleared RAM and a koron CPU wired to it.
func New() *Machine {
	m := &Machine{lmpr: 0x1F, hmpr: 0x00, unknownOut: map[uint8]int{}}
	cpu := &z80.CPU{Memory: m, IO: m}
	m.cpu = cpu
	return m
}

// CPU exposes the underlying koron CPU for register inspection / setup.
func (m *Machine) CPU() *z80.CPU { return m.cpu }

// --- memory: koron z80.Memory ------------------------------------------------

func (m *Machine) resolveReadPage(addr uint16) int {
	section := addr >> 14
	switch section {
	case 0:
		if m.lmpr&0x20 == 0 {
			return -1 // ROM0 — we have no ROM, read as 0xFF
		}
		return int(m.lmpr & 0x1F)
	case 1:
		return int((m.lmpr&0x1F + 1) & 0x1F)
	case 2:
		return int(m.hmpr & 0x1F)
	default: // 3
		if m.lmpr&0x40 != 0 {
			return -1 // ROM1
		}
		return int((m.hmpr&0x1F + 1) & 0x1F)
	}
}

func (m *Machine) Get(addr uint16) uint8 {
	pg := m.resolveReadPage(addr)
	if pg < 0 {
		return 0xFF
	}
	return m.ram[pg][addr&0x3FFF]
}

func (m *Machine) Set(addr uint16, value uint8) {
	pg := m.resolveReadPage(addr)
	if pg < 0 {
		return // ROM write dropped
	}
	m.ram[pg][addr&0x3FFF] = value
}

// --- IO: koron z80.IO --------------------------------------------------------

func (m *Machine) In(port uint8) uint8 {
	switch port {
	case portLMPR:
		return m.lmpr
	case portHMPR:
		return m.hmpr
	case portVMPR:
		return m.vmpr
	case portStatus:
		// Active-low interrupt bits | keyboard rows (high = no key).
		v := uint8(0xFF)
		if m.framePending {
			v &^= statusIntFrame // assert frame interrupt (bit clear)
			// reading the status acknowledges the (edge-triggered) frame int
			m.framePending = false
		}
		return v
	case portKbd:
		return 0xFF // no keys pressed (active low)
	}
	return 0xFF
}

func (m *Machine) Out(port uint8, value uint8) {
	switch port {
	case portLMPR:
		m.lmpr = value
	case portHMPR:
		// bits 5-7 are mode-3 CLUT; preserve.
		m.hmpr = (m.hmpr & 0xE0) | (value & 0x1F)
	case saaLowByte:
		// Recover A8 from B register (OUT (C),r idiom). A8 high => address.
		isAddr := (m.cpu.BC.Hi & 0x01) != 0
		if isAddr {
			m.curSAAReg = value & 31
		}
		if m.captureSAA {
			m.Writes = append(m.Writes, SAAWrite{
				Frame: m.frame, IsAddr: isAddr, Reg: m.curSAAReg, Val: value,
			})
		}
		if m.saaWrite != nil {
			m.saaWrite(isAddr, value)
		}
	case portVMPR:
		m.vmpr = value
	case portStatus:
		// LINE port write: sets next line-interrupt line. The frame-only
		// timing model does not raise line interrupts, so this is a no-op.
	case 0xF8:
		// CLUT / palette (and HPEN/LPEN) — visual only, ignore.
	case 0xFE:
		// border / speaker — ignore
	default:
		m.unknownOut[port]++
	}
}

// SetSAAHandler installs the callback invoked for every SAA write (after A8
// decode). isAddr true => WriteAddress(val); false => WriteData(val).
func (m *Machine) SetSAAHandler(fn func(isAddr bool, val uint8)) { m.saaWrite = fn }

// EnableSAACapture turns on the per-frame register-stream log in m.Writes.
func (m *Machine) EnableSAACapture(on bool) { m.captureSAA = on }

// Frame returns the current frame counter.
func (m *Machine) Frame() int { return m.frame }

// SetFrame sets the frame counter (used to tag captured writes).
func (m *Machine) SetFrame(f int) { m.frame = f }

// UnknownOuts returns ports written that the model does not handle, with counts.
func (m *Machine) UnknownOuts() map[uint8]int { return m.unknownOut }

// --- paging control ----------------------------------------------------------

func (m *Machine) SetLMPR(v uint8) { m.lmpr = v }
func (m *Machine) SetHMPR(v uint8) { m.hmpr = (m.hmpr & 0xE0) | (v & 0x1F) }
func (m *Machine) LMPR() uint8     { return m.lmpr }
func (m *Machine) HMPR() uint8     { return m.hmpr }

// --- raw page access ---------------------------------------------------------

// LoadPage copies data into a physical page starting at offset, spilling into
// consecutive pages on 16 KB boundaries.
func (m *Machine) LoadPage(page, offset int, data []byte) {
	for i, b := range data {
		p := page + (offset+i)/PageSize
		o := (offset + i) % PageSize
		m.ram[p&(NumPages-1)][o] = b
	}
}

// ReadPage reads length bytes from physical page storage starting at offset,
// spanning page boundaries.
func (m *Machine) ReadPage(page, offset, length int) []byte {
	out := make([]byte, length)
	for i := range out {
		p := page + (offset+i)/PageSize
		o := (offset + i) % PageSize
		out[i] = m.ram[p&(NumPages-1)][o]
	}
	return out
}

// PeekPhys reads a byte from a physical page at offset (paging-independent).
func (m *Machine) PeekPhys(page, off int) byte { return m.ram[page&(NumPages-1)][off&0x3FFF] }

// PeekPhysWord reads a little-endian 16-bit word from a physical page.
func (m *Machine) PeekPhysWord(page, off int) uint16 {
	return uint16(m.PeekPhys(page, off)) | uint16(m.PeekPhys(page, off+1))<<8
}

// --- SAM "long address" helpers ---------------------------------------------
//
// The FRED loader (jukebox.bas) uses SAM BASIC long addresses where a decrunch
// page byte is (addr DIV 16384) - 1 and that value is written straight to HMPR
// to map the page at &8000. So physical page = (addr DIV 16384) - 1.

// PhysFromLinear maps a BASIC long address to (physical page, offset).
func PhysFromLinear(linear int) (page, offset int) {
	return linear/PageSize - 1, linear % PageSize
}

// LoadLinear deposits data as `LOAD name CODE linear` would.
func (m *Machine) LoadLinear(linear int, data []byte) {
	pg, off := PhysFromLinear(linear)
	m.LoadPage(pg, off, data)
}

// Poke writes one byte at a long address (POKE linear,val).
func (m *Machine) Poke(linear int, val byte) {
	pg, off := PhysFromLinear(linear)
	m.LoadPage(pg, off, []byte{val})
}

// DPoke writes a 16-bit little-endian word at a long address (DPOKE).
func (m *Machine) DPoke(linear int, val uint16) {
	m.Poke(linear, byte(val&0xFF))
	m.Poke(linear+1, byte(val>>8))
}

// --- execution ---------------------------------------------------------------

// RunUntilPC single-steps until PC reaches target, HALT, or maxSteps. Returns
// the reason ("pc", "halt", "steps") and the step count consumed.
func (m *Machine) RunUntilPC(target uint16, maxSteps uint64) (string, uint64) {
	var n uint64
	for n = 0; n < maxSteps; n++ {
		if m.cpu.PC == target {
			return "pc", n
		}
		m.cpu.Step()
		m.Steps++
		if m.cpu.HALT {
			return "halt", n
		}
	}
	return "steps", n
}

// StepFrame raises a 50 Hz frame interrupt (IM1 → &0038) and runs `budget`
// Z80 instructions so the player's ISR services it and the main loop spins
// until the next frame. The frame counter (used to tag SAA writes) is
// advanced first. Returns true if the CPU halted.
func (m *Machine) StepFrame(budget uint64) bool {
	m.frame++
	m.framePending = true
	m.cpu.Interrupt = z80.IM1Interrupt()
	for i := uint64(0); i < budget; i++ {
		m.cpu.Step()
		m.Steps++
		if m.cpu.HALT {
			return true
		}
	}
	return false
}

// Run runs `n` raw instructions (no interrupt), for init sequences.
func (m *Machine) Run(n uint64) string {
	for i := uint64(0); i < n; i++ {
		m.cpu.Step()
		m.Steps++
		if m.cpu.HALT {
			return "halt"
		}
	}
	return "steps"
}
