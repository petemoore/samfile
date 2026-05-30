// Command decrunch runs Andrew Collier's E-Tracker "decruncher" routine inside
// a minimal SAM Coupé Z80 model (koron-go/z80) to produce the decrunched
// E-Code player image at logical 0x8000, then writes it to a file.
//
// This mirrors jukebox.bas line 60:
//
//	CLEAR 32767: LOAD "decruncher" CODE : LET mc_start=16384
//	decrunch "E-Code",32768,98304,41,44,53
//
// where PROC decrunch (line 360) does:
//
//	LOAD name$ CODE caddr                       ; E-Code -> caddr (98304)
//	POKE mc_start+3, daddr DIV 16384 - 1        ; = 1   (daddr page)
//	DPOKE mc_start+4, daddr MOD 16384 + 32768   ; = 0x8000 (daddr addr)
//	POKE mc_start+6, caddr DIV 16384 - 1        ; = 5   (caddr page)
//	DPOKE mc_start+7, caddr MOD 16384 + 32768   ; = 0x8000 (caddr addr)
//	POKE mc_start+9, 41,44,53                   ; decode codes
//	CALL mc_start                               ; run decruncher at 0x4000
//
// The decruncher pages source (caddr) and destination (daddr) banks into
// section C (0x8000) via OUT (0xFB),page (HMPR), reads/writes at 0x8000, and
// bumps the page when it crosses 0xC000. We model 32 physical 16K pages with
// LMPR/HMPR paging (same scheme as the sam-aarch64 z80-test-harness), deposit
// E-Code into the caddr page, run from 0x4000 until RET to a sentinel, then
// dump the daddr page(s) as the 0x8000 player image.
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

// sam is a minimal SAM Coupé memory+IO model: 32 physical RAM pages plus
// LMPR (port 0xFA) / HMPR (port 0xFB) paging. No ROM is needed — the
// decruncher runs entirely from RAM (section B) and pages RAM into section C.
type sam struct {
	ram        [numPages][pageSize]byte
	lmpr, hmpr uint8
}

// page resolves the physical page index backing a logical address.
//
//	A (0x0000-0x3FFF): LMPR bit5? RAM(lmpr&0x1F) : ROM  -> we use RAM
//	B (0x4000-0x7FFF): RAM (lmpr&0x1F + 1)
//	C (0x8000-0xBFFF): RAM (hmpr&0x1F)
//	D (0xC000-0xFFFF): RAM (hmpr&0x1F + 1)
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
	}
	return 0xFF
}
func (s *sam) Out(port uint8, v uint8) {
	switch port {
	case 0xFA:
		s.lmpr = v
	case 0xFB:
		s.hmpr = v // low 5 bits = section-C physical page
	}
}

func main() {
	var (
		decruncherPath = flag.String("decruncher", "", "path to decruncher CODE file")
		ecodePath      = flag.String("ecode", "", "path to crunched E-Code CODE file")
		outPath        = flag.String("out", "", "output path for decrunched 0x8000 image")
		dumpLen        = flag.Int("len", pageSize, "bytes to dump from 0x8000")
	)
	flag.Parse()
	if *decruncherPath == "" || *ecodePath == "" || *outPath == "" {
		fmt.Fprintln(os.Stderr, "usage: decrunch -decruncher F -ecode F -out F [-len N]")
		os.Exit(2)
	}

	decr, err := os.ReadFile(*decruncherPath)
	must(err)
	ecode, err := os.ReadFile(*ecodePath)
	must(err)

	s := &sam{}
	// Default paging at BASIC CALL time: LMPR=0x1F (page0 in B), HMPR=... we
	// only care that section B holds the decruncher and we can reach the
	// caddr/daddr pages by their physical index, which the decruncher selects
	// explicitly via OUT (0xFB). Start HMPR at 0 (harmless).
	s.lmpr = 0x1F
	s.hmpr = 0x00

	// The decruncher CODE file is org'd at 0x4000 (first bytes: C3 0C 40 =
	// JP 0x400C). 0x4000 is section B = physical page (lmpr&0x1F + 1). With
	// lmpr=0x1F, section B = page 0. Deposit the decruncher into page 0.
	bPage := int((s.lmpr&0x1F + 1) & 0x1F)
	copy(s.ram[bPage][:], decr)

	// Patch the decrunch parameters exactly as PROC decrunch pokes them, at
	// logical 0x4000+offset (which lands in page bPage at offset 0+).
	poke := func(off int, b ...byte) { copy(s.ram[bPage][off:], b) }
	poke(3, 1)          // daddr page = 32768/16384 - 1 = 1
	poke(4, 0x00, 0x80) // daddr addr = 0x8000 (little-endian)
	poke(6, 5)          // caddr page = 98304/16384 - 1 = 5
	poke(7, 0x00, 0x80) // caddr addr = 0x8000
	poke(9, 41, 44, 53) // decode codes

	// Deposit E-Code into the caddr physical page (5). The decruncher does
	// OUT (0xFB),5 then reads 0x8000 (section C), i.e. physical page 5 off 0.
	const caddrPage = 5
	if len(ecode) > pageSize {
		copy(s.ram[caddrPage][:], ecode[:pageSize])
		copy(s.ram[caddrPage+1][:], ecode[pageSize:])
	} else {
		copy(s.ram[caddrPage][:], ecode)
	}

	// Run from 0x4000 with a sentinel return address so we detect the final
	// RET (the decruncher ends ...22 04 40 C9 = RET).
	const sentinel = 0xFFFE
	cpu := &z80.CPU{Memory: s, IO: s}
	cpu.PC = 0x4000
	cpu.SP = 0x6000 // stack in section B RAM, away from code/data
	// Push sentinel: when the routine RETs, PC := sentinel and we stop.
	s.Set(cpu.SP-1, byte(sentinel>>8))
	s.Set(cpu.SP-2, byte(sentinel&0xFF))
	cpu.SP -= 2

	steps := 0
	const maxSteps = 50_000_000
	for steps < maxSteps {
		if cpu.PC == sentinel {
			break
		}
		cpu.Step()
		steps++
		if cpu.HALT {
			fmt.Fprintf(os.Stderr, "HALT at PC=%04X after %d steps\n", cpu.PC, steps)
			break
		}
	}
	if steps >= maxSteps {
		fmt.Fprintf(os.Stderr, "WARNING: hit step cap (%d), PC=%04X\n", maxSteps, cpu.PC)
	}

	// The decrunched image lives in the daddr pages. daddr page started at 1
	// (physical page 1 = logical 0x8000-0xBFFF); page 2 = 0xC000-0xFFFF if it
	// spilled. Concatenate pages 1.. for the requested length.
	out := make([]byte, 0, *dumpLen)
	for p := 1; len(out) < *dumpLen; p++ {
		need := *dumpLen - len(out)
		if need > pageSize {
			need = pageSize
		}
		out = append(out, s.ram[p][:need]...)
	}
	must(os.WriteFile(*outPath, out, 0o644))
	fmt.Fprintf(os.Stderr, "decrunched %d bytes -> %s (%d steps, final PC=%04X HMPR=%02X)\n",
		len(out), *outPath, steps, cpu.PC, s.hmpr)
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
