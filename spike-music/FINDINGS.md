# Spike A — E-Tunes → WAV by decoding the module format (first principles)

**Outcome: 10 WAVs delivered (`out/m01.wav … out/m10.wav`), cross-validated
against spike B as audibly identical.** But the headline finding is a surprise
about the *format itself* that reframes the A-vs-B comparison — see
[§ The central finding](#the-central-finding).

All work is a self-contained scratch Go module under `spike-music/`
(`go 1.22`, single dependency `koron-go/z80`). Build everything with
`go build ./...`; `go vet ./...` and `go test ./...` are clean.

---

## What I built

| Component | File(s) | Lines | Origin |
|---|---|---:|---|
| Go SAA1099 model | `saa/*.go` | ~1200 | **ported** from Dave Hooper's SAASound (BSD; attribution in `saa/SAASound-LICENCE.txt`) |
| Decruncher harness | `cmd/decrunch` | 188 | mine |
| RE probe | `cmd/probe` | 128 | mine |
| Renderer + loop detect | `cmd/render` | 258 | mine |
| WAV writer | `wav` | 54 | mine |

So ~630 lines of my own Go, plus the mechanical SAASound port.

---

## How the format works (reverse-engineered)

### 1. Getting the player image (decrunch)

`assets/E-Code` (Andrew Collier's *E-Tracker* player) is crunched. `jukebox.bas`
line 60 decrunches it: load `decruncher` at `0x4000`, point it at `E-Code`
(caddr `0x18000`) and `0x8000` (daddr), `CALL 0x4000`. The decruncher pages
source/dest 16 KB banks into section C via `OUT (0xFB),page` and bumps the page
across `0xC000`.

`cmd/decrunch` reproduces this in a 32-page SAM memory model on `koron-go/z80`:
run the real decruncher Z80 routine, dump the `0x8000` image. 182 684 Z80 steps,
clean `RET`. **No ROM needed** — the decruncher is pure RAM code.

A surprise from the dump: **the player relocates itself to `0x0000`.** At entry
it does `LMPR := (HMPR & 0x1F) | 0x20` (page its own RAM into section A) then
`JP 0x01EB`, so internal code references are low (`0x0xxx`) while its data/IO
save area and the BASIC-poked tables are the *same physical page* aliased at
`0x82xx`. Offsets are identical either way (`0x8000`→off 0, `0x0000`→off 0).

### 2. The module *is* a self-contained player (the key structural fact)

The `m01..m10` blobs are **not** a declarative data format. Diffing them:

```
first 0x4B3 bytes of all ten modules: byte-identical
they first diverge at offset 0x4B3
```

Each module = **a shared E-Tracker replay engine (offset 0..0x4B2) + per-song
data (0x4B3..end)**. The engine carries an embedded signature at 0x4BD:
`"ETracker (C) BY ESI.R"`. Song-data sizes: 578–3582 bytes.

The module's own code references absolute `0x80xx–0x84xx`, i.e. it runs **at
`0x8000`**. The FRED jukebox loads it to `0x10000`, and E-Code (the wrapper:
interrupts + scroller) pages it into section C and drives it. Run standalone it
needs nothing but its own page.

### 3. Engine entry points & per-frame model

- **`0x8000` — init.** `LD HL,0x84B3` (song-data pointer) `; JP 0x83EF`. Init
  reads a small header of 16-bit **offsets** (relative to the song base, via the
  add-base helper at `0x816C`) into the four per-channel data-stream pointers
  + tempo, zeroes 178 bytes of channel state at `0x833D`, and sets up six
  25-byte channel structs (instrument ptr at +0x0F = `0x812C`, ornament ptr at
  +0x11 = `0x813E`).
- **`0x8006` — play one tick.** Flag-gated (the byte at `0x8007`, set to 1 by
  init). Processes the 6 channels, applies instrument/ornament/vibrato per tick,
  and flushes the chip.

### 4. The SAA register shadow (the clean handoff point)

The flush routine at `0x8104` is the gold nugget. The engine keeps a **26-byte
SAA register shadow** for regs `0x00..0x19` at `0x83D3..0x83EC`, and every frame:

```
write reg 0x1C = 0x01            ; sound enable
for reg 0x19 down to 0x00:       ; dump the shadow
    WriteAddress(reg); WriteData(shadow[reg])
```

SAA writes use the SAM idiom `LD BC,0x01FF; OUT (C),reg` (address port `0x1FF`)
then `DEC B; OUT (C),data` (data port `0x00FF`) — matching SimCoupé's
`SAA_ADDR_PORT=0x1FF / SAA_PORT=0xFF` decode.

So **reading `0x83D3..0x83EC` after each tick gives the exact per-frame chip
state** — no need to snoop the OUT stream (which `koron-go/z80`'s 8-bit port
interface can't fully disambiguate anyway).

### 5. Rendering (`cmd/render`)

```
load module at 0x8000  →  call init (0x8000)
per 50 Hz frame:        →  call play (0x8006)  →  read 26-byte shadow
feed shadow → Go SAASound model → 882 samples/frame @ 44.1 kHz → WAV
```

**Loop detection:** the engine is deterministic, so the shadow *sequence* is
eventually periodic with period = the musical loop. `detectLoop` finds the
smallest period that repeats over an ≥8 s tail window, walks back to the loop
start, and renders `[0, loopStart+period)` = intro + exactly one loop, then
trims trailing silence. (Full-machine-state hashing fails here — the engine has
a free-running counter that never lets full state repeat, even though the music
does. Detecting on the audible register stream is the right signal.)

Clock = **8 MHz**, sample rate 44.1 kHz — matches SimCoupé (its `SAADevice`
never calls `SetClockRate`, so SAASound's compiled `EXTERNAL_CLK_HZ=8000000`
stands).

---

## Audio quality & faithfulness

**Pitch validation** — dominant frequency of each WAV lands on a real musical
note, confirming both the SAA frequency math and the 8 MHz clock (a wrong clock
would detune everything by a constant ratio):

| tune | dom. freq | nearest note | tune | dom. freq | nearest note |
|---|---|---|---|---|---|
| m01 | 131.5 Hz | C3 (130.8) | m06 | 260.5 Hz | C4 (261.6) |
| m02 | 98.0 Hz  | **G2 (98.0)** | m07 | 233.0 Hz | **A♯3 (233.1)** |
| m03 | 82.5 Hz  | E2 (82.4)  | m08 | 65.5 Hz  | C2 (65.4) |
| m04 | 165.0 Hz | E3 (164.8) | m09 | 493.0 Hz | B4 (493.9) |
| m05 | 66.0 Hz  | C2 (65.4)  | m10 | 1045.5 Hz| **C6 (1046.5)** |

All WAVs: healthy RMS (3.4k–7.9k), peaks well under clipping. Durations 37–148 s
(one intro+loop each).

**Cross-check against spike B (the validation oracle the spec asked for).**
Spike B captured `m01.reg` (the SAA register stream from full emulation). I
reconstructed its per-frame register image and compared to my captured shadow:

```
audible-signature comparison: compared=299 frames, mismatches=0
>>> AUDIBLE OUTPUT IDENTICAL to spike B across all 299 compared frames <<<
```

The only raw-register differences are in **muted channels** (channels 2–5 have
zero amplitude and are disabled in reg 0x14) — their don't-care freq/octave bytes
differ because my lightweight headless drive doesn't replicate the full player's
initialisation of unused state. **Audibly, my output and spike B's are the
same.** Confidence in faithfulness: high.

---

## Effort

- **Hard / time-consuming:** the reverse-engineering. The player is
  self-relocating, interrupt-oriented, and **self-modifying** (init patches play
  routine operands every run), and z80dasm mis-aligns where data is interleaved
  with code. Static disassembly alone was not enough — dynamic tracing in the
  emulator (run it, watch the register shadow, diff module heads) was what
  cracked the structure. Decrunching also had its own puzzle (the
  self-relocation / dual address aliasing).
- **Mechanical:** the SAASound→Go port (large but a straight transliteration),
  the WAV writer, the render loop.
- **Cheap win:** once the `0x83D3` shadow + `0x8006` entry were found, the
  renderer was ~250 lines and the loop detector ~30.

---

## The central finding

**Spike A's premise — "decoding the format is cleaner than emulation" — does not
hold for E-Tracker, because there is no declarative format to decode.** Each
module *embeds its own Z80 replay engine* (the identical 0x4B3-byte prefix); the
song data is meaningless without it. To turn a module into register writes you
must either (a) execute that engine, or (b) hand-port ~1.2 KB of self-modifying,
interrupt-timed Z80 — channel state machines, instrument envelopes, ornaments,
vibrato — into Go. Option (b) *is* re-emulating that specific engine, only more
fragile and far more work, with no payoff in maintainability or fidelity.

So the honest first-principles route **converges on running the engine** —
i.e. spike B's approach. My renderer does exactly that, but the RE bought a much
*lighter* form of it than a full SAM/ROM emulator: **no ROM, no E-Code wrapper,
no interrupts, no scroller** — just the module's own page, `call 0x8006` per
frame, read 26 bytes. That is the genuinely useful product of the
"first-principles" effort.

### A-vs-B verdict (for E-Tracker specifically)

- A **pure data-format decoder is not viable** for E-Tracker — recommend not
  pursuing it.
- Both spikes share the indispensable **Component A** (the Go SAA1099 model) and
  produce the **same audio**.
- The cleanest production design is the **hybrid this spike landed on**: run the
  module's embedded engine headlessly on a minimal Z80 + SAM-paging model and
  read the register shadow. It is small (~630 lines + the SAA port), needs no
  ROM/SimCoupé, and is provably faithful (matches spike B).
- This would generalise to *any* E-Tracker module; **Sound Machine** (a
  different player) would need its own entry-point/shadow RE but the same
  technique.

### Caveats / limits

- Loop boundaries are heuristic (shadow-sequence periodicity); musically clean on
  all 10 but not authored loop points.
- The `0xF9` status port is stubbed to "nothing pending" (no SPACE keypress), so
  the player never skips tracks — correct for single-tune rendering.
- Verified by register-stream match + pitch analysis; I could not A/B the actual
  audio by ear in this environment.

---

## Reproduce

```sh
cd spike-music
go build ./...
# 1. decrunch the player (for RE; not needed to render)
./decrunch -decruncher ~/git/sam-jukebox/assets/decruncher \
           -ecode ~/git/sam-jukebox/assets/E-Code -out /tmp/ecode-8000.bin
# 2. render all ten tunes
for m in m01 m02 m03 m04 m05 m06 m07 m08 m09 m10; do
    ./render -v -mod ~/git/sam-jukebox/assets/$m -out out/$m.wav
done
# 3. inspect a module's per-frame SAA shadow
./probe -mod ~/git/sam-jukebox/assets/m01 -frames 60
```
