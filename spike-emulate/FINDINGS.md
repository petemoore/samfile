# Spike B — E-Tunes → WAV by emulation: FINDINGS

**Goal:** convert Pete's 10 FRED E-Tracker tunes (`m01`..`m10`) to WAV by running
Andrew Collier's *real* E-Code player headlessly on a Go Z80 + a Go SAA1099,
capturing the SAA register writes per 50 Hz frame. No understanding of the
E-Tracker module format required — the player is the spec.

**Result: it works.** All 10 tunes render to clean, structured, non-silent
stereo WAV (`out/m01.wav`..`out/m10.wav`, 60 s each). The captured per-frame SAA
register stream is dumped alongside (`out/m01.reg`..`out/m10.reg`) as the
validation oracle for spike A.

---

## What was built

A self-contained Go module (`spike-emulate/`):

| Piece | What it is |
|---|---|
| `saa/saa.go` | A faithful, class-by-class Go port of Dave Hooper's **SAASound** (SAA1099). `WriteAddress`/`WriteData`/`GenerateMany`. BSD-3 `LICENCE` preserved. |
| `sam/machine.go` | A minimal headless **SAM Coupé**: 512 KB paged RAM (LMPR/HMPR), the SAA `OUT` decode, VMPR/STATUS/keyboard ports, a 50 Hz IM1 frame interrupt, and SAM "long address" `LOAD CODE`/`POKE`/`CALL` helpers. Paging model adapted from `z80-test-harness-go`. Built on `koron-go/z80`. |
| `cmd/etunes/` | Orchestrator: reproduces `jukebox.bas` (decrunch → poke freq table → load module → run player), intercepts SAA writes → SAA → WAV, and dumps the register stream. |
| `cmd/dumpplayer/` | Runs just the decruncher and dumps the decrunched 32 KB player image (`out/player.bin`) for disassembly. |
| `scripts/gen-freqtable.py` | Regenerates `out/freqtable.bin` (the 169-byte note table) from `jukebox.bas`. |

Run: `go build ./cmd/etunes && ./etunes -m m07 -frames 3000 -wav out/m07.wav -reg out/m07.reg`

---

## How hard was the headless player setup? (medium-hard, but tractable)

The genuinely fiddly parts, in order of pain:

1. **SAM "long address" → physical page mapping.** `jukebox.bas` uses SAM BASIC
   long addresses (`LOAD CODE 98304`, `CALL 32768`). The decruncher's poke
   parameters revealed the convention: **physical page = `addr / 16384 - 1`**,
   and that page number is written straight to HMPR to map it at `&8000`. Once
   that clicked, every load/poke/call in the loader fell into place
   (decruncher→page 0, player→page 1, module→page 3, crunched E-Code→page 5).

2. **The SAA OUT decode under koron-go.** SAM selects the SAA *address* register
   at port `&1FF` and *data* at port `&FF` — the two differ only in **A8**.
   But koron-go's `IO.Out(port uint8, …)` only passes the low byte (`&FF` for
   both). Fix: recover A8 from the **B register** at OUT time (`OUT (C),r`
   idiom), which the player uses. Confirmed against SimCoupé `SAMIO.h`
   (`SAA_ADDR_PORT=0x1ff`, `SAA_MASK=0x1ff`).

3. **The player relocates itself.** `CALL 32768` enters at `&8000`, but the
   player immediately re-pages so its own page is *also* visible at `&0000`
   (sets `LMPR = (HMPR&0x1F)|0x20`) and `JP`s into the `&0000` mirror. All its
   real run-time addresses are `&0000–&3FFF`; `&0038` is its IM1 handler. You
   have to disassemble at `org 0`, not `org &8000`, to read it.

4. **Frame timing / interrupts.** The player is interrupt-driven: a 50 Hz frame
   interrupt (IM1 → `&0038`) runs the per-frame replay (calls the unrolled
   audio routines it builds in page 4, then walks the song), and **line**
   interrupts drive the palette raster for the scroller. koron-go has no
   T-state counter, so cycle-accurate raster timing isn't available. **The key
   realisation: audio is entirely on the *frame* interrupt; line interrupts are
   purely visual.** So the model raises *only* a frame interrupt per "frame"
   (STATUS port reports frame-pending, active-low) and never line interrupts.
   The scroller silently does nothing useful (no screen) — harmless.

What made it *tractable* rather than a multi-day slog: the emulator itself is
the oracle. Run the decruncher → dump → disassemble → see valid code; run the
player → histogram OUT ports → see exactly where the SAA is written. No need to
understand the E-Tracker format at all.

### Things that did **not** need solving
- The E-Tracker module format (B1's whole job). B2 treats it as opaque.
- Line-accurate raster / palette (visual only).
- ROM: the player and decruncher are self-contained; no SAM ROM is needed
  except an IM1 vector, which the player provides itself at `&0038`.

---

## Audio quality & faithfulness

- **Sample format:** 16-bit stereo, 44100 Hz (SimCoupé's `SAMPLE_FREQ`).
- **SAA clock:** 8 MHz — SimCoupé's SAASound default (`EXTERNAL_CLK_HZ`; it
  never calls `SetClockRate`), so pitch matches SimCoupé.
- **64× oversample + 5 Hz DC-blocking high-pass + 11.35× output gain** — all
  SAASound defaults, as SimCoupé uses them.
- **Levels:** peaks 7 k–24 k of the 32767 range, no clipping, RMS ~2–4 k.
  Headroom matches SAASound's design (6 channels × 480 × 11.35 ≈ 32 k max).
- **Stereo:** 4 of the 10 tunes use genuine L≠R panning; the rest are centred
  (the player writes symmetric L/R amplitude nibbles, e.g. `&44`, `&88`).
- **Tempo: faithful.** The player refreshes all 6 channels' frequency + 6
  amplitude registers **exactly once per frame** (300 freq writes/s = 6 × 50),
  with continuous per-frame pitch changes — exactly how E-Tracker applies
  ornaments at the 50 Hz tick. One tune-tick ↔ one frame interrupt ↔ 1/50 s of
  audio, so timing is correct by construction.

**Why this route is "guaranteed faithful":** it is the actual shipping player
code executing the actual module bytes. The only things that *can* differ from a
real SAM/SimCoupé are (a) the SAA1099 model — but it's a faithful port of the
same SAASound SimCoupé uses, and (b) sub-frame write timing (see Limitations).

### Residual faithfulness risk (recommended next check)
A direct PCM comparison against **SimCoupé's WAV recorder** (the spec's
"Fallback") has *not* been done in this time-box. That's the recommended final
validation: record one tune in SimCoupé and diff against `out/m0?.wav`. The
register-stream oracle below is the more important cross-check and is the whole
point of running B2.

---

## The register-stream oracle (for spike A)

`out/m0?.reg` logs **every** SAA write tagged by frame:

```
# frame   kind   reg   value
0         ADDR   28    -
0         DATA   28    2        <- reset/sync
1         DATA   28    1        <- enable
1         DATA   22    64       <- noise control
...
```

This is the ground truth for B1: run the same tune through B1's
first-principles decoder, dump its `(frame, reg, value)` stream, and `diff`.
They should match line-for-line (the SAA register state is a pure function of
song position, independent of the noise-LFSR/oversample internals). Format is
tab-separated and `diff`-friendly. The streams are distinct per tune (verified),
confirming each tune plays its own module data.

---

## Loop detection / duration

The tunes loop, but **the SAA register stream does not exactly repeat within
240 s** for `m01` (autocorrelation found no period ≤ 120 s; the song note
pointer climbs monotonically without wrapping in 320 s). These are long /
through-composed FRED tunes, and reliable automatic loop-trimming needs the
**order-list structure**, which is exactly what B1 decodes — so loop-point
extraction is properly spike-A territory.

Pragmatic choice for the deliverable: **fixed 60 s render**, leading silence
trimmed, 1 s fade-out for a clean end. The tool renders any length
(`-frames N`, 50 frames = 1 s), so full-length captures are one flag away once a
loop point is known.

---

## The real question: is B2's machinery worth it vs B1?

**Cruft / indirection that B2 carries** (and B1 would not):

- A Z80 CPU (koron-go) + a SAM hardware model (paging, ports, interrupts).
- Reproducing the FRED BASIC loader byte-for-byte (decrunch, pokes, paging).
- A self-relocating, interrupt-driven, scroller-entangled player whose audio you
  reach only by emulating frame interrupts and intercepting `OUT`.
- ~600 lines of SAM/Z80 glue on top of the (shared, unavoidable) ~750-line SAA.

**What that machinery buys you:**

- **Correctness for free.** Zero reverse-engineering of the module format. It
  played all 10 tunes correctly on the first complete run. The risk profile is
  "did I model the hardware right?" not "did I understand a 30-year-old tracker
  format right?"
- **A faithful oracle.** B2's register stream is precisely what lets B1 be
  *proven* correct rather than *believed* correct.

**Assessment.** For *shipping* a clean, fully-Go `samfile music-to-wav`, B1 is
the better destination: no Z80 emulator dependency, direct loop-point detection
from the order list, no scroller/raster baggage, and it generalises to other
E-Tracker modules without a matching BASIC loader. But B2 is **cheap insurance
and a strong validation tool**: it took roughly a day, the SAA port is shared
with B1 anyway, and it converts the "decode the format" risk into a one-line
`diff` against ground truth. **Recommendation: keep B2 as the test oracle; build
B1 as the product, validated frame-for-frame against these `.reg` dumps.**

The single biggest surprise: the *hard* part wasn't the SAA or the Z80 — it was
the SAM paging conventions and discovering that the player relocates to `&0000`
and hides its audio behind a frame interrupt. Once those were understood, the
SAASound port did the rest verbatim.

---

## Deliverables in this directory

- `out/m01.wav` .. `out/m10.wav` — the 10 tunes, 60 s, 16-bit/44.1 kHz stereo.
- `out/m01.reg` .. `out/m10.reg` — per-frame SAA register streams (oracle).
- `out/player.bin` — decrunched 32 KB E-Code player image (disassemble at
  `org 0`; useful for both spikes).
- `out/freqtable.bin` — the 169-byte note table (regenerate with
  `scripts/gen-freqtable.py`).
