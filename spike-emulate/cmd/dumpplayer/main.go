// dumpplayer runs the FRED decruncher headlessly to produce the decrunched
// E-Code player image, and writes pages 1+2 (32 KB) to a file for disassembly.
package main

import (
	"flag"
	"fmt"
	"os"

	"spike-emulate/sam"
)

func main() {
	assets := flag.String("assets", "/Users/pmoore/git/sam-jukebox/assets", "assets dir")
	out := flag.String("o", "out/player.bin", "output file (decrunched pages 1+2)")
	flag.Parse()

	decr, err := os.ReadFile(*assets + "/decruncher")
	must(err)
	ecode, err := os.ReadFile(*assets + "/E-Code")
	must(err)

	m := sam.New()
	// decruncher: LOAD "decruncher" CODE -> linear 16384 (page 0)
	m.LoadLinear(16384, decr)
	// crunched E-Code: LOAD "E-Code" CODE 98304 (page 5)
	m.LoadLinear(98304, ecode)

	// PROC decrunch "E-Code",daddr=32768,caddr=98304,41,44,53 pokes:
	const mc = 16384
	m.Poke(mc+3, byte(32768/16384-1))    // 1  (dest HMPR page)
	m.DPoke(mc+4, uint16(32768%16384+32768)) // 0x8000 (dest offset)
	m.Poke(mc+6, byte(98304/16384-1))    // 5  (src HMPR page)
	m.DPoke(mc+7, uint16(98304%16384+32768)) // 0x8000 (src offset)
	m.Poke(mc+9, 41)
	m.Poke(mc+10, 44)
	m.Poke(mc+11, 53)

	// Run setup: LMPR=0x1F -> section B (&4000) = page 0 (decruncher code).
	m.SetLMPR(0x1F)
	// Sentinel return address on the stack; detect completion when popped.
	m.LoadPage(0, 0x3EFE, []byte{0xFF, 0xFF}) // logical &7EFE = 0xFFFF
	cpu := m.CPU()
	cpu.SP = 0x7EFE
	cpu.PC = 0x4000

	reason, steps := m.RunUntilPC(0xFFFF, 20_000_000)
	fmt.Printf("decrunch: reason=%s steps=%d totalSteps=%d HMPR=%02X LMPR=%02X PC=%04X\n",
		reason, steps, m.Steps, m.HMPR(), m.LMPR(), cpu.PC)
	if uo := m.UnknownOuts(); len(uo) > 0 {
		fmt.Printf("unknown OUT ports: %v\n", uo)
	}

	// Decrunched player landed at dest HMPR page 1 (&8000) spilling to page 2.
	img := m.ReadPage(1, 0, 2*sam.PageSize)
	must(os.WriteFile(*out, img, 0o644))
	fmt.Printf("wrote %d bytes to %s\n", len(img), *out)

	// Show a little of the head so we can sanity-check it's code.
	fmt.Printf("first 32 bytes @page1 &8000: % 02X\n", img[:32])
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
