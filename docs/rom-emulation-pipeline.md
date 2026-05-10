# ROM-emulation pipeline — design notes

An alternative-and-complementary architectural angle for samfile: rather
than reimplement SAM ROM / DOS behaviour in Go, drive the real ROM
(and optionally the real DOS) under a Z80 emulator and let SAM's own
binary do the work. We get byte-identical output for free, and any
behaviour we haven't documented (yet) just falls out of the emulator
correctly.

This document captures four architectural insights from a 2026-05-14
spike session and the one concrete bug it surfaced in samfile's
current behaviour. Each insight has the same shape: **don't
reimplement what already exists on SAM in binary form — emulate the
SAM, let its own routines do the work, capture the outputs**.

**Relationship to the Pike-style lexer plan.** An earlier exploration
(the hand-written pure-Go lexer specified in
[`docs/superpowers/specs/2026-05-14-sambasic-text-lexer-design.md`](superpowers/specs/2026-05-14-sambasic-text-lexer-design.md)
and the lexical grammar reference in
[`docs/sambasic-grammar.md`](sambasic-grammar.md)) approached the same
problem by reimplementing `TOKMAIN`'s behaviour byte-for-byte in Go.
The ROM-emulation approach below is the **preferred direction** going
forward — it gets correctness for free, supports ROM extensions and
non-SAMDOS DOSes naturally, and avoids the maintenance cost of mirroring
ROM behaviour in two implementations. The Pike-style design documents
are kept on record as the alternative that was explored; their
lexical-grammar notes remain valuable as documentation of *how* the SAM
tokeniser works (independent of the implementation strategy).

A pure-Go lexer can still be useful as a **second opinion** for
validation: feed both implementations the same source text, diff their
outputs. If/when that's wanted, the ROM-emulation pipeline is the
oracle.

---

## 1. text → tokenised BASIC (the spike — working today)

**Idea.** Boot the SAM ROM under a Z80 emulator. Reach the BASIC
editor's main loop. Inject the source text one keystroke at a time
through the same `FLAGS` bit-5 / `LASTK` channel the ROM's `KYIP2`
(ROM L1786 / `KYIP2` @ `0x050A`) reads. The ROM's own `TOKMAIN`
(ROM L13028 / `TOKMAIN` @ `0x3872`) tokenises the line as it would
on real hardware. Extract the tokenised bytes from the ELINE buffer
on the post-CR unwind, sort lines by line number, lay them out as a
`sambasic.File`.

**Status.** Working end-to-end. Verified byte-identical to a real
SAM SAVE for a non-trivial program (`PRINT "HELLO SAM"`, `FOR i=1 TO 5`,
`PRINT i*i`, `NEXT i`, `PRINT "DONE"`).

**Prototype location.**

  - Branch: `worktree-spike-basic-rom-emulation` in `~/git/sam-aarch64`
  - Code: `tools/basic-emulator-spike/main.go` (~800 lines)
  - Head commit (working end-to-end): `38ba060`
  - Dependencies: `github.com/koron-go/z80` (pure-Go Z80 core) +
    `github.com/petemoore/samfile/v3`

**Key design decisions captured during the spike.**

  - **Banner skip via PC hijack at MAINER3.** Cold boot reaches the
    "MILES GORDON TECHNOLOGY plc" banner at WTFK (ROM L3886 /
    `WTFK` @ `0x0FA2`). Instead of dismissing it via matrix keypress,
    the spike hijacks PC at the `CALL ERRHAND1` inside MAINER3
    (`0x0F75 → 0x0F78`) so MAINER3's state setup still runs
    (CLSLOWER / SET 5,(TVFLAG) / RES 7,(FLAGS)) but the banner draw
    and key wait are skipped. Boot reaches MAINELP at step 878k
    (~20 ms wall on M1).
  - **No keyboard matrix emulation needed.** Once the boot reaches
    the editor, key input goes through FLAGS bit 5 / LASTK; we hook
    `Memory.Get` for those addresses and synthesise the bit / byte.
    No need to model the SAM keyboard matrix or the WD1772
    debouncer.
  - **Sidestepped INSERTLN.** After CR is consumed, `RST 8` with
    error code 0 fires (the "OK unwind" path) before INSERTLN runs.
    The spike doesn't fight this — it extracts the tokenised bytes
    directly from ELINE on RST 8 entry and assembles PROG in Go after
    sorting. Architecturally cleaner is to let INSERTLN run; that's
    follow-up work.
  - **Snapshot/restore amortises the boot.** The spike snapshots the
    full emulator state (RAM + paging + CPU registers) at MAINELP
    entry, then restores it before each line. Restore is ~50 µs vs
    ~20 ms cold boot. Per-line cost is ~1.2 ms steady-state.

**Use as samfile oracle.** The spike's MGT output is the closest
thing to "what real SAMDOS+SAM would have written" without using
real hardware. For testing the planned `samfile text-to-basic`
lexer, the spike can act as a ground-truth oracle: feed both the
hand-written lexer and the spike the same source text, diff the
PROG bytes. Disagreement = at least one of them is wrong (and the
spike has the higher prior of being right, because it uses the
actual SAM tokeniser).

## 2. tokenised BASIC → text (already partially solved differently)

**Idea.** Symmetric inverse of (1). Place the tokenised PROG bytes
into RAM, inject `LLIST<CR>` via the same FLAGS/LASTK channel,
intercept the printer stream, capture bytes.

**Status.** Not implemented in the spike, but partially solved
already: `~/git/sam-aarch64/tools/llist-capture.sh` builds a
one-shot test disk and uses SimCoupé's parallel-port-to-file
mechanism plus a `DI; HALT` stub for clean exit to capture the
ROM's LLIST output. See
[`docs/sambasic-roundtrip-caveats.md`](sambasic-roundtrip-caveats.md)
for the existing workflow.

**Why both might still be useful.** SimCoupé-based capture relies
on having SimCoupé installed and on the parallel-port hook. The
spike-style equivalent would be a pure-Go in-process variant —
faster, more deterministic, no external binary. For test suites
that already validate against `llist-capture.sh`, that's redundant;
for new contexts (CI without SimCoupé, integration into other Go
tools) the in-process variant could pay off.

**Where to trap printer output.**

  - `LLIST` (ROM L2013 / `LLIST` @ `0x0637`) and `LPRINT`
    (ROM L2255 / `LPRINT` @ `0x07A3`) both route through stream 3.
  - `IOPENT` (ROM L21215 / `IOPENT` @ `0xDC20`) stages each char in
    `OPCHAR = 0x5A72` (sysvar).
  - Cleanest hook is probably the printer channel's output function
    via `CURCHL` when stream 3 is selected; alternatively trap the
    actual Centronics data port (TBD — needs research).

## 3. Arbitrary file → DOS-format disk via emulated SAM (future)

**Idea.** samfile currently emits MGT bytes by reimplementing the
SAMDOS on-disk format in Go. That works for SAMDOS but doesn't
extend to MasterDOS, B-DOS, ProDOS, SCADS-enhanced disks, or any
DOS variant with different directory layouts or file metadata.

If samfile instead uses the emulated SAM, and lets the **real DOS**
running on that SAM do the disk work, every DOS that ever existed
for the SAM is automatically supported.

**Pipeline for `samfile add <file> --load-address <addr> -o <disk.mgt>`:**

  1. Load `disk.mgt` into a virtual MGT via WD1772 FDC emulation so
     the emulated SAM sees it as a real floppy.
  2. Boot the SAM ROM. The disk auto-boots → its DOS loads → DOS
     hooks into BASIC.
  3. Inject `<file>`'s bytes directly into SAM RAM at `<addr>`.
  4. Inject the DOS's `SAVE "<name>" CODE <addr>,<len>` (or whatever
     command that DOS provides) via the FLAGS/LASTK channel from (1).
  5. DOS runs its actual SAVE logic — writes sectors, updates
     directory, handles whatever format quirks that DOS uses.
  6. Extract the resulting MGT.

**Why this is elegant.** Future-proof: any DOS we have a binary for
is supported. Authentic: bytes on disk are byte-identical to what
real SAM running that DOS would write. Composable: combined with
(1), entire disks can be built from human-readable source (BASIC
text + asset files + a DOS image to use).

**Prerequisite.** FDC (WD1772) emulation. SimCoupé's
`~/git/simcoupe/Base/` has a working reference implementation. This
is the load-bearing new piece.

## 4. Extended-BASIC tokenisation via autobooted extension (future)

**Idea.** SAM BASIC extensions (SCADS, others) work by an
autobooting CODE file that loads itself into RAM and hooks into
BASIC's tokeniser via ROM extension vectors (`MTOKV`, `CMDV`,
`EVALUV`, `MEPRO2`, etc. — explicitly out of scope for the planned
pure-Go lexer per
[`docs/superpowers/specs/2026-05-14-sambasic-text-lexer-design.md`](superpowers/specs/2026-05-14-sambasic-text-lexer-design.md)).

If the spike's input disk has such an extension autobooted, the
spike's tokeniser pipeline gets the extended dialect **for free**:

  - Cold boot → extension's autoboot runs → BASIC is now extended.
  - Snapshot at the extended MAINELP.
  - Inject keystrokes — the editor uses the extension's patched
    `TOKMAIN`, so ELINE ends up with the extension's custom tokens
    alongside the standard SAM BASIC ones.
  - Extract from ELINE as before.

The output `.mgt` only runs correctly on a SAM with the same
extension loaded — that's the right semantics.

**Prerequisite.** Same FDC emulation as (3). Both (3) and (4)
unlock at once.

---

## Empirical evidence: the `NumericVars` default-init bug

The spike surfaced a concrete bug in samfile's current
`sambasic.File` defaults that's worth its own ticket.

**Symptom.** A disk built via `samfile.DiskImage.AddBasicFile`
crashes at auto-run on real SAM (and on SimCoupé) when the program
uses any variable — e.g. `FOR i=1 TO 5`. `samfile basic-to-text`
round-trip on the same disk produces the correct source text, so
the tokenisation is right.

**Cause.** `sambasic.File.numericVars()` defaults to
`make([]byte, 92)` — 92 zeros. Real SAM ROM SAVE produces the
NumericVars area as **46 bytes of `0xFF`** (the 23 letter-pointer
pairs marking "no variable defined") followed by **46 bytes of
PSVTAB content** (the pre-saved-variables table copied from ROM at
NEW / cold-boot — see ROM L13209–L13230 / `CLRSR` @ `0x396B`).

Auto-running BASIC that touches variables walks into corrupt
letter-pointers in the zero-filled area and crashes. Auto BASIC
that only does `CLEAR`/`LOAD CODE`/`CALL` (as build-disk does) never
exercises variables, so the zero default was untested.

**Fix.** Replace the zero default in `sambasic.File.numericVars()`
with the canonical init. The 92 bytes are stable and can be
hardcoded; the spike does this. A cleaner fix is to read them from
the ROM (`PSVTAB` @ ROM `0x39E3`, copied to NVARS area by `CLRSR`)
but that ties samfile to having a SAM ROM available, which it
currently isn't.

**Workaround used by the spike** (canonical bytes verified
byte-identical against Pete's hand-typed-and-SAVE'd ground-truth
disk):

```go
canonicalNumericVars := []byte{
    // 46 bytes of 0xFF — letter-pointer "no var defined" sentinels
    0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
    0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
    0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
    0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
    0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
    0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
    // 46 bytes of PSVTAB content (extracted from a real SAM SAVE)
    0x19, 0x00, 0x03, 0x00, 0xFF, 0xFF, 0x02, 0x08,
    0x00, 0x6F, 0x73, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x02, 0xFF, 0xFF, 0x72, 0x67, 0x00, 0x00, 0xC0,
    0x00, 0x00, 0x02, 0x08, 0x00, 0x6F, 0x73, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x02, 0xFF, 0xFF, 0x72,
    0x67, 0x00, 0x00, 0x00, 0x01, 0x00,
}
```

**Architectural observation.** This is exactly the kind of bug the
ROM-emulation pipeline catches automatically. If samfile's `AddBasicFile`
were routed through `SAMDOS HSAVE` under emulation (see (3)), the
canonical vars init would have been emitted by the real SAMDOS SAVE
as a side-effect. Every layer of "real" SAM software we keep in the
pipeline is one less class of subtle encoding bugs to discover by
hand.

---

## How these relate to each other

| # | Direction                              | Mechanism                                    | Status                                                   |
|---|----------------------------------------|----------------------------------------------|----------------------------------------------------------|
| 1 | text → tokenised BASIC                 | inject keystrokes, capture ELINE             | **working** in spike (`worktree-spike-basic-rom-emulation`, head `38ba060`) |
| 2 | tokenised BASIC → text                 | inject LLIST, capture printer stream         | partially solved via SimCoupé+`llist-capture.sh`         |
| 3 | arbitrary file → DOS-format disk       | inject SAVE under target DOS                 | future, needs FDC                                        |
| 4 | text → extended-dialect tokenised BASIC | autoboot SCADS-like extension first         | future, needs FDC (shares (3)'s prereq)                  |

(3) and (4) unlock together once the WD1772 FDC is wired in. (1)
already runs on `koron-go/z80` with just RAM + paging emulation —
no FDC needed because the spike doesn't boot from a disk; it boots
the ROM image directly with the banner-skip hijack.
