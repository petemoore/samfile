# spike-emulate — E-Tunes → WAV by emulation (Spike B)

Time-boxed spike: convert the 10 FRED E-Tracker tunes to WAV by running Andrew
Collier's real **E-Code** player headlessly on a Go Z80 (`koron-go/z80`) + a Go
port of **SAASound** (SAA1099), capturing the chip's register writes.

See **[FINDINGS.md](FINDINGS.md)** for the write-up and the B1-vs-B2 comparison.

## Build & run

```sh
go build ./...

# render one tune (60 s) to WAV + register-stream oracle
go build -o etunes ./cmd/etunes
./etunes -m m07 -frames 3000 -wav out/m07.wav -reg out/m07.reg

# render all 10
for t in m01 m02 m03 m04 m05 m06 m07 m08 m09 m10; do
  ./etunes -m $t -frames 3000 -wav out/$t.wav -reg out/$t.reg
done
```

Flags: `-m <module>`, `-frames N` (50 = 1 s), `-budget` (Z80 insns/frame),
`-trim` (silence trim + fade, default on), `-hist` (OUT port histogram),
`-ptr` (trace the song pointer). `-assets` defaults to
`~/git/sam-jukebox/assets`.

## How it works (one paragraph)

`cmd/etunes` reproduces `jukebox.bas`: it runs the `decruncher` to decrunch
`E-Code` to the player image, pokes the 169-byte frequency table, loads a module,
sets up SAM paging exactly as `CALL 32768` would, then runs the player. The
player relocates itself to `&0000` and is driven by a 50 Hz IM1 frame interrupt;
its `OUT`s to the SAA (port `&1FF`/`&FF`, A8 recovered from the B register) are
routed to `saa.SAA`, and one frame of PCM is pulled per frame → WAV. Every SAA
write is logged per frame as the validation oracle for spike A.

## Layout

- `saa/` — Go port of SAASound (BSD-3, attribution in `saa/LICENCE`).
- `sam/` — headless SAM machine (paging, ports, frame interrupt) on koron-go/z80.
- `cmd/etunes/` — the converter.
- `cmd/dumpplayer/` — decrunch + dump the player image for disassembly.
- `scripts/gen-freqtable.py` — regenerate `out/freqtable.bin` from `jukebox.bas`.
- `out/` — deliverables (10 WAVs, 10 register streams, player image, freq table).

## Attribution

SAA1099 emulation ported from Dave Hooper's **SAASound**
(https://github.com/stripwax/SAASound), BSD-3-Clause — see `saa/LICENCE`.
E-Code player by Andrew Collier (ILLUSION); tunes by Pete Moore (FRED 51–56).
