# SAM BASIC v3.0 — Lexical Grammar Reference

A reference for implementing a Pike-style state-function lexer that converts plain-text SAM BASIC source into the on-disk tokenised body bytes consumed by `samfile basic-to-text` (and produced by the planned inverse `samfile text-to-basic`).

The deliverable is to emulate the SAM BASIC v3.0 editor's line-entry behaviour byte-for-byte. **Stock interpreter only** — the optional ROM extension vectors (`MTOKV`, `CMDV`, `EVALUV`, `MEPRO2`, etc.) are out of scope.

Citations are to:
- `~/git/sam-aarch64/docs/sam/sam-coupe_rom-v3.0_annotated-disassembly.txt` — `(ROM L#### / ROUTINE)`.
- `sam-coupe_tech-man_v3-0.txt` — `(TM p#)`.
- `sam-coupe_use-guide.txt` — `(UG p#)`.

Confidence prefix:
- **[ROM]** — proven by reading the disassembly.
- **[TM]** / **[UG]** — asserted by the manual; assumed unless ROM contradicts.
- **[inferred]** — best reading, but not directly confirmed.

---

## 1. Overview

**Character set.** SAM BASIC source is a stream of **SAM Coupé bytes**, one byte per glyph across the full `0x00`–`0xFF` range. Bytes `0x80`–`0xFF` are valid SAM characters (graphics symbols, accented letters) and appear in REM bodies and string literals as legitimate visible content. An implementing lexer **reads bytes, not Unicode runes** — UTF-8 decoding the input would mangle high bytes (`utf8.RuneError` → truncated to `0xFD` on cast) and is a category error against the SAM character set.

When the user presses Enter, the SAM BASIC editor passes the freshly-typed line — a buffer of raw ASCII bytes terminated by `0x0D` (CR) — to two routines in sequence:

1. **`TOKMAIN`** (ROM L13028 / `TOKMAIN` @ `0x3872`) walks the line left-to-right, identifies runs of alphabetic characters and the operator chars `<`/`>`, and replaces any complete keyword by its 1-byte token (`0x85`–`0xF6`) or 2-byte token (`0xFF nn`, `nn` in `0x3B`–`0x83`). Keyword matching is **case-insensitive** (`AND 0DFH` fold) and **strictly word-bounded** (`LOADX` is not `LOAD`+`X`; it is the procedure name `LOADX`). String literals (`"…"`) and everything after `REM` are skipped untouched. (ROM L15877–L15979)

2. **`LINESCAN`** (ROM L3457 / `LINESCAN` @ `0x0D13`) re-parses the now-tokenised line as a syntax check, and during that pass **inserts an "invisible" 6-byte form after every numeric literal**: marker `0x0E`, followed by 5 bytes encoding the value on the floating-point calculator stack format. (ROM L5261–L5269 / `INSERT5B`; TM p77.) During the same pass, the `LIF`/`LELSE` tokens get patched in place: the shared `SIF`/`LIF` handler at L6340 evaluates the IF's expression, then checks whether the byte immediately after that expression is `THEN` — if so it rewrites the `LIF` to `SIF` (`LD (HL),0D8H` at L6358) and records `IFTYPE = THENTOK` for the rest of the line. Later on the same line, `LELSE`'s handler (`FELSLP` at L6424) reads `IFTYPE`; if it is `THENTOK` it patches the `LELSE` byte to `ELSE` (`LD (HL),0DAH` at L6447, label `NLELS`). There is also a `LELSE LIF` → `LELSE SIF` patch at L6438–L6440 for the `ELSE IF cond THEN …` chain. (ROM L6340–L6364 / `SIF`/`LIF`/`IFL1`; L6424–L6450 / `FELSLP`/`NLELS`.)

The tokenised line is then prefixed with `MSB LSB LenLo LenHi` and the trailing `0x0D` is kept as the terminator. (TM p77 "FORMAT OF A BASIC PROGRAM".)

This document specifies the **input recogniser** — the subset of "what the editor would have accepted on Enter" — without modelling any of the run-time semantic errors that fire later.

---

## 2. Whitespace and line structure

### 2.1 Line terminator
- A SAM BASIC line ends at the first `0x0D` (CR). This byte is **preserved** in the tokenised body as the line's last byte. (ROM L3556 / `STMTLP3`; TM p77.)
- LF (`0x0A`) has no special meaning to the editor — but the lexer is consuming a text file from a Unix or Windows host, so the implementer should accept and discard a `\r\n` or `\n` line ending and emit a single `0x0D`.
- An input line consisting of only `0x0D` (no text, with or without a leading line number) is **valid but silently dropped**. (ROM L4131–L4133 / `INSLN3`.)

### 2.2 Statement separator
- `:` (`0x3A`) separates statements within a line. Stored verbatim as ASCII `0x3A`. (ROM L3563 / `STMTLP3`.)
- The parser also treats `THEN` (`0x8D`) as a "statement separator" in certain contexts — but only at run time / syntax-check time. The lexer treats `THEN` as just another keyword.

### 2.3 Leading whitespace and line-number prefix
- The editor accepts a *line number* in the range **1 to 0xFEFF (65279)** at the start of the buffer. The on-disk 16-bit field has the same maximum because `0xFF` is reserved as the program-end sentinel (TM p77 "the final line in the program is followed by FFH (so the maximum line number allowed is FEFFH)"; UG p50 "programs can have up to 65279 lines"). The cap is enforced inside `EVALLINO` at L4077–L4079: after `INTTOFP` + `FPTOBC` give BC = parsed value, `LD A,B; ADD A,1; RET C` sets carry (which the callers treat as a fatal `NONSENSE` error — see `MAINE2` L3779 `JP C,NONSENSE` and `LINESCAN` L3466 `JR NC,STMTLP1`) whenever B = `0xFF`, i.e. for any value in `0xFF00..0xFFFF`. (ROM L4068–L4083 / `EVALLINO`; ROM L3778–3779 / `MAINE2`; ROM L3465–3468 / `LINESCAN`.)
  - **SimCoupé empirical (2026-05-14):** line numbers `1`, `2`, `10`, `10000`, `65279` all accepted; line `0` rejected (reserved as the "no line number" signal); line `65280` (`0xFF00`) rejected (fails the `B = 0xFF` cap). The earlier note about `10000` being rejected was a mistaken observation — re-tested with a stock ROM and stands accepted. **Lexer policy: accept 1..0xFEFF inclusive; reject 0 and anything ≥ 0xFF00.**
  - Line 0 is reserved (`Z if BC=0` from `EVALLINO` is the "no line number" signal).
  - **The line number is parsed as a decimal integer via `INTTOFP` then `FPTOBC`** — so it can have leading zeros (`0010 PRINT …` is line 10) but cannot use `&hex` or `.frac` or `E` notation.
- **Leading whitespace before the line number is skipped.** `EVALLINO` enters via the character-fetch routines (`RST 18H` / `GTCH3`, ROM L394–L408) which silently skip every byte in `0x00`–`0x1F` *and* `0x20` (space) before delivering the first significant character to the line-number parser. So `    10 PRINT "hi"` parses as line 10; any quantity of leading spaces, tabs, or other control bytes is fine. **Confirmed by ROM (`GTCH3` skip + `EVALLINO` entry path) and by observation.**
- **After the line number, the editor performs a conditional one-space drop** before storing the body. The exact rule, from `INSERTLN` at ROM L4106–L4116, is:
  - Let `b1` = byte at CHAD (i.e. the first byte after the parsed digits).
  - If `b1` is not a space (`0x20`): no modification. CHAD is left where it is. Body = the bytes from CHAD up to (not including) `0x0D` verbatim. (Jump at L4108 `JR NZ,INSLN3`.)
  - If `b1` is a space, examine `b2` = byte at CHAD+1:
    - If `b2` is `0x0D`: **preserve the space.** CHAD is left at `b1`. Body = `0x20 0x0D` (single-space body). This is the special case at L4116 `DEC BC ;AVOID ANY ACTION IF E.G. 10 (space) CR` whose explicit purpose is to make the line look stored-with-its-space rather than as an empty-body delete.
    - If `b2` is anything else: **drop the space.** CHAD is advanced past `b1` (via `JR NZ,INSLN2` at L4113 with BC already INCed). Body = bytes from `b2` up to (not including) `0x0D`. The annotated ROM comment at L4113–L4115 calls this out explicitly: "INC CHAD AND DELETE FIRST SPACE IN A LINE LIKE: 10 test. PREVENTS SPACES ACCUMULATING WITH MULTI-EDIT/ENTERS."
  - Net effect: typing `10 PRINT "hi"` and typing `10PRINT "hi"` both store identical bodies (the editor's `LIST` formatter inserts its own space between line number and body at display time — see ROM L26009–L26010 / `OUTLN3`). Typing `10  PRINT "hi"` (two spaces) stores a body with one leading space, because only the first space is dropped.
  - **SimCoupé observation (confirmed):** `10 \n` stores body `0x20 0x0D` (matches the L4116 special case).
  - **SimCoupé observations (confirmed 2026-05-14):**
    - `20 PRINT "with space"` → body length 14 (PRINT keyword + `"with space"` + `0x0D`); no leading space byte — the one input space was dropped.
    - `20  X` (two spaces) → body length 9, body = `0x20 0x58 0x0E …(5 invisible bytes) 0x0D` — one of the two spaces dropped, the second preserved as the leading body byte.
    - These confirm the three-branch rule end-to-end.
  - **Lexer policy**: implement the one-space drop. Given an input line `<spaces><digits><body><0x0D>`, after parsing the digits, peek at the first body byte: if it is `0x20` AND the next byte is not `0x0D`, skip it. (Equivalently: take the verbatim bytes from CHAD-after-parsing-digits to `0x0D`, then if the body length is ≥ 2 and the body starts with `0x20`, remove that first byte.)
- A line typed without a leading line number is an *immediate-mode command* and never goes on disk. **The lexer's job is for stored programs only**, so the input format must always begin with a line number (matching what `samfile basic-to-text` emits — see `sambasic.go` L65: `"%5d "`).

### 2.4 Whitespace in the body
- The character fetch routines `GETCHAR` / `NEXTCHAR` (RST 18H / RST 20H) **transparently skip any byte in `0x00`–`0x1F` except `0x0D`**. (ROM L394–L408 / `GTCH1`.) That is, when LINESCAN walks the line it treats most control codes as if they were whitespace; only `0x0D` (line end) is significant, plus `0x0E` (number marker).
- Spaces (`0x20`) are **stored as-is** in the body — they are not collapsed. The tokeniser deliberately preserves spaces so that listings re-display them. The lone exception is the one trailing space immediately following a tokenised keyword, which is consumed (see §3.5): the test at ROM L15965–L15970 (`TOK55`/`TOK6`) checks the byte right after the matched keyword, and if it is `0x20` includes that one byte in the close-up region passed to `RECLAIM1`. Only one space is consumed; a second trailing space stays in the body.
- Spaces are significant **only as keyword/identifier boundaries**. They cannot appear inside a number literal (the lexer must reject `1 23` as a single number, but `1 23` may be valid as token+separator+token at line entry — the editor will read it as `1` (number), then either find a context where `23` is another expression or report a syntax error.)
- TABs (`0x09`), formfeeds, etc. are accepted by the *editor*, but on entry they are silently dropped. **Lexer policy**: treat ASCII TAB as a space.

### 2.5 Empty / no-text line
- After stripping the line-number, if all that remains is `0x0D`, the line is treated as a **delete** of that line number from the program (ROM L4131 / `INSLN3`, jumping past the `INSERTLN` action).
- For the lexer producing a fresh file, this case is irrelevant: a saved file never contains an empty line for an existing number.

---

## 3. Keyword tokenisation rules (TOKMAIN)

The keyword table is at upper-ROM `KEYWTAB` = `0xF8C8` (ROM L26722). Format:
- One leading separator byte (`0xA0` at the start; the high bit of `0xA0` is what `GETTOKEN`'s `BIT 7,(HL)` scan keys off).
- Each keyword is its uppercase ASCII letters with **bit 7 set on the last letter** (e.g. `PI` is `0x50 0xC9` = `'P' 'I'|0x80`). Multi-word keywords have an embedded space (`0x20`) between words — e.g. `GO TO` is `0x47 0x4F 0x20 0x54 0xCF` (`'G' 'O' ' ' 'T' 'O'|0x80`).
- Keywords are listed in the order `PI` (index 0x3B), `RND`, … `INK` (index 0xFF, reserved). The index plus `0x3B` is the function/operator code (3B–84) or, for commands, plus `0x3B` puts you in 0x85–0xFE. (ROM L15938–L15946 / `TOK42`.)

The full canonical-spelling table is already in `/Users/pmoore/git/samfile/sambasic/keywords.go`; the lexer should not duplicate it.

### 3.1 Where TOKMAIN looks
- TOKMAIN reads through the line **starting from the first character past any line-number** (line-number processing is in `INSERTLN` / `EVALLINO` and happens *before* tokenisation). (ROM L13028.)
- The main loop `TOKMLP` (ROM L15877) advances character by character. The cases below tell it whether to try a keyword match.

### 3.2 Candidate first character
TOKMAIN only attempts to match a keyword when CHAD points to:
- A letter (`A`–`Z` or `a`–`z`) — `ALPHA` carry-set (ROM L13621 / `ALPHA`, L15883 / `POSFIRST`).
- `<` (`0x3C`) — to catch `<=` and `<>` (ROM L15885).
- `>` (`0x3E`) — to catch `>=` (ROM L15888).

All other characters — digits, punctuation, `"`, `&`, `(`, `)`, `+`, `-`, `*`, `/`, `=`, `,`, `;`, `:`, `^`, `?`, control codes, high-bit bytes — are **never the start of a keyword match**.

### 3.3 The match itself (GETTOKEN — ROM L20029–L20087)
TOKMAIN copies up to 15 bytes from the line into a scratch buffer in CDBUFF and calls `GETTOKEN`. GETTOKEN walks the keyword table comparing:

- **Letter case is folded via `XOR (HL); AND 0DFH`** (ROM L20044, L20064). This is bit-5 masking. `'A' XOR 'a' = 0x20`, so AND with `0xDF` masks bit 5. For ASCII letters this is the standard case-fold. **Beware**: bit-5 masking also folds `@`↔`` ` ``, `[`↔`{`, `\`↔`|`, `]`↔`}`, `^`↔`~`, `_`↔`?`, plus all of `0x00`–`0x1F` with `0x20`–`0x3F`. But since the *candidate* must be a letter (per §3.2) and the table entries are letters/space, this can only produce false matches if the input contains low-control characters in mid-word position — and those are stripped by GTCH1 anyway. In practice: **assume strict ASCII case-fold on `A`–`Z`/`a`–`z`**.

- **Embedded spaces in multi-word keywords are optional in the input.** GTTOK3 (ROM L20051–L20061): "if the list has a space at this position, accept either a matching space in the input or skip past it". So `GO TO` and `GOTO` both match the keyword `GO TO`; but `GO  TO` (two spaces) does *not* — see the strict-one-space rule below.
  - **The match consumes exactly one input space, or none, per list-side space.** Multiple consecutive spaces in input where the list has one space will *fail*: GTTOK3 doesn't loop over extra spaces. ROM behaviour is "compare next list char vs next input char; if list has space, also try input no-space"; a second space stays in the input. **Verified on SimCoupé**: typing `GO  TO` (two spaces) into the SAM BASIC editor is rejected — `GO` then `TO` is not a valid statement. **Multi-word keyword spacing is strictly one space**: between the words of a multi-word keyword you must type either zero spaces (run them together, e.g. `GOTO`) or exactly one space (`GO TO`). Any other amount of whitespace is a syntax error from LINESCAN's perspective.
  - Single-word keywords cannot contain spaces. Input `LO AD` does *not* match `LOAD`; the partial-match fails at the space and falls through.

- **Word-boundary check after match (GTTOK6 — ROM L20077–L20082):** after the last letter has matched, the routine looks at the input character *immediately past* the matched word. If the keyword's last letter is a letter (i.e. the keyword does not end in `$`, `=` or `>`), it calls `ALDU` and requires the trailing input character to be **not** a letter, **not** `_`, **not** `$`. (ROM L13653 / `ALDU`.) Otherwise the candidate is rejected and search continues.
  - Consequence: `LOADX`, `LOAD_X`, `LOAD$` are not tokenised — they are stored verbatim as ASCII characters (and are syntactically interpreted as identifier `LOADX` / `LOAD_X` / variable name `LOAD$`).
  - Keywords ending in `$` (e.g. `INKEY$`, `CHR$`, `HEX$`, `STR$`) bypass the trailing-letter check — so `INKEY$abc` would tokenise as `INKEY$` + `abc`. **[inferred from GTTOK6 logic; corroborated by Tech Manual p55: "A keyword which is followed by a letter will not be recognised, so printx is assumed to be a procedure name, but print x and print1 become PRINT x and PRINT 1"]** (TM p55.)
  - The operators `<>`, `<=`, `>=` also bypass the trailing check (the `>`/`=` ending tells ALDU to skip — ROM L20077: `CP 7EH; CCF`).

- **Maximal munch is not used.** The table is walked **in order from PI (3B) downwards through the indices**, and the *first* matching entry wins. (ROM L20039 / `RET Z` after `DEC C`.) Practical consequences:
  - Both `IF` entries exist in the table (`0xD7` LIF and `0xD8` SIF). Per the explanatory comment ROM L6358–L6362: at line entry, `I F` always matches `LIF` (0xD7) because LIF appears first; LINESCAN later patches to `SIF` if `THEN` follows. **The lexer must emit `0xD7` for `IF` always** — let LINESCAN's THEN-patching logic happen after, or (simpler) handle the patch inside the lexer's higher-level grammar pass. The patch is *necessary* for correct round-tripping of corpus files.
  - Same situation for `ELSE`: the table has `LELSE` (0xD9) before `ELSE` (0xDA). The lexer emits `0xD9` always; LINESCAN patches to `0xDA` if the line contains a preceding `SIF`. (ROM L6447 / `NLELS`.)
  - **Special hack**: if `GETTOKEN` returns `INK` (entry `0xFF` in the table), TOKMAIN rewrites the token to `PEN` (`0xA1`). (ROM L15947–L15950 / `TOK42`.) This is the editor's compat shim for Spectrum users typing `INK`. Lexer must replicate: input `INK` (or `ink`, `Ink`, etc.) tokenises to `0xA1` (`PEN`).
  - The table also has two `IF`s and two `ELSE`s as noted; per ROM matching order, the **first** of each pair (`LIF` 0xD7, `LELSE` 0xD9) is the one TOKMAIN ever emits.

### 3.4 1-byte vs 2-byte token byte
- Table index `0x3B` through `0x84` → emitted as `0xFF nn` (2 bytes). (ROM L15938 / L15941.)
- Table index `0x85` through `0xFE` → emitted as 1 byte. (ROM L15946 / `TOK42`.)
- Mapping is mechanical: `token_byte = table_index + 0x3B`. The keywords.go file in this repo already encodes the canonical spellings for both ranges.

### 3.5 Result: where the bytes land in the line
- TOKMAIN overwrites the **first letter** of the matched word with the token byte (replacing 1 letter with 1 byte for the command range, or with the leading `0xFF` for the function range; the second byte of a 2-byte token is then written at +1). It then closes up the rest of the matched letters via `RECLAIM1`. (ROM L15962–L15972.)
- If the matched word had a *leading* space in the input, that leading space is overwritten by the token (no space byte remains). (ROM L15952–L15958 / `TOK43`, `TOK5`.)
- If the matched word had a *trailing* space, that single trailing space is **also closed up** (consumed). The test is a single `CP " "; JR NZ,TOK6; INC HL` at L15965–L15969: it looks at exactly the one byte immediately after the matched keyword and, if it is `0x20`, extends the close-up region by one byte. No look-ahead beyond that — the rule is "drop exactly one trailing space if present, regardless of what follows it." So `PRINT X` becomes `PRINTTOK X` (no intervening space), `PRINT  X` becomes `PRINTTOK 0x20 X` (one of the two spaces consumed), and `PRINT"hi"` is unaffected (no space to drop, no space ever introduced). (ROM L15965–L15969 / `TOK55`, `TOK6`.)
- **One *leading* space is also closed up around tokenised keywords.** Empirical (SimCoupé 2026-05-14): `10 PRINT "a":PRINT "b"` stores 10 bytes (`<PRINT>"a":<PRINT>"b" CR`); `20 PRINT "a" : PRINT "b"` stores 11 bytes (`<PRINT>"a" :<PRINT>"b" CR` — the space *before* `:` is kept, but the space *between* `:` and the second `PRINT` is gone); `30 :: PRINT "x"` stores 7 bytes (`::<PRINT>"x" CR`). If only trailing-space drop existed, those bodies would be 10 / 12 / 9 bytes. They're 10 / 11 / 7 — exactly one space removed on each side of every tokenised keyword whose adjacent byte is `0x20`. The corresponding ROM code is presumed to be a `DEC DE; LD A,(DE); CP " "; JR Z, …` test paired with the trailing-space test at L15965–L15970, but the precise location has not been re-derived; flagged for ROM citation. **Lexer policy**: when emitting a keyword token byte, consume one immediately-preceding `0x20` (if any) from the buffer in addition to the immediately-following one (if any). The combined rule is "drop up to one space on each side of every tokenised keyword". Both drops are independent: a leading space is dropped iff present; the same for trailing. Two leading spaces leave one, two trailing spaces leave one.

This is a small but real subtlety: **the byte stream after tokenisation cannot in general be reconstructed by spelling out keywords and joining with spaces**, because the original whitespace around them has been collapsed unpredictably. For a *forward* lexer this doesn't matter — the lexer just needs to know "after a keyword, if input has a space, drop it; otherwise emit no space".

### 3.6 Recognition is suppressed inside…
- **String literals (`"…"`)** — TOKMAIN's main loop has a special arm at L15896–L15905: when it sees `"`, it scans forward to the next `"` (or `0x0D`) and resumes from there. No keywords inside strings.
- **After `REM`** (token `0xB7`) — after a successful tokenisation, TOKMAIN checks the just-emitted byte; if it's `0xB7` it returns immediately, leaving the rest of the line as raw ASCII. (ROM L15977–L15981 / `TOK6`/`TOKFIN`.) So `REM PRINT "hi"` keeps `PRINT` as the five ASCII letters `0x50 0x52 0x49 0x4E 0x54`.
- **After an `FF` byte (FN leader)** — line L15891–L15893 / L15907 / `FNTS`. When TOKMAIN's character-scan encounters a `0xFF` (already-emitted FN leader from an earlier match), it skips the *next* byte (the FN code) before resuming. This handles re-tokenisation of already-tokenised text (e.g. when `MERGE`ing or when the editor processes a line that was reloaded for editing).
- **After a single-byte keyword (`0x85`–`0xF6`)** that isn't `REM` — TOKMAIN simply doesn't *recognise* those as candidates (per §3.2 they aren't letters or `<`/`>`) so it walks past them.

### 3.7 Stream-number suffix tokens (`PRINT #`, `INPUT #`, `OPEN #`, `CLOSE #`, `CLS #`)
- **`#` is never part of a keyword token**. `PRINT#3;"hi"` tokenises as `PRINTTOK '#' '3'…` — three independent items. Recognition of the stream syntax is in the statement-handler routines, not the tokeniser. (See for example ROM L20053–L20055 in `GETTOKEN` where `#` is not in the keyword table; and the run-time `IMINKEYS` at ROM L5173 doing its own `CP "#"` check.)

### 3.8 Multi-word keywords
The complete list of multi-word entries in `keywords.go` (i.e. keyword strings containing a space) is the **complete set of multi-word keywords**:

- `GO TO` (0xB4), `GO SUB` (0xB5)
- `DEF FN` (0xC8), `DEF KEYCODE` (0xC9), `DEF PROC` (0xCA), `END PROC` (0xCB)
- `LOOP IF` (0xD3), `EXIT IF` (0xD6)
- `END IF` (0xDB)
- `ON ERROR` (0xDD)

All accept their words run together: `GOTO`, `GOSUB`, `DEFFN`, `DEFKEYCODE`, `DEFPROC`, `ENDPROC`, `LOOPIF`, `EXITIF`, `ENDIF`, `ONERROR` are all valid input and tokenise the same.

### 3.9 Operator-token forms
The "operator" keyword entries (matched only when `<` or `>` is the candidate first char):
- `<>` → `0xFF 0x81` (NOTEQUAL)
- `<=` → `0xFF 0x82` (LESSEQ)
- `>=` → `0xFF 0x83` (GREATEREQ)

These are tokenised by TOKMAIN since `<` and `>` are first-char candidates. `=`, `<`, `>` on their own are **single ASCII bytes** (`0x3D`, `0x3C`, `0x3E`), not tokens. (See OPPRIORT table ROM L5417–L5438 — `<`/`=`/`>` get codes 0x12/0x13/0x14, and the run-time scanner looks for them as ASCII at ROM L5316.)

### 3.10 Stated non-features

These are surface forms users coming from other dialects might expect, which stock SAM BASIC v3.0 explicitly does **not** support:

- **No `?` shortcut for `PRINT`.** Spectrum BASIC and many other dialects accept `?` as a synonym for `PRINT`. SAM does not. Per §3.2, `?` (`0x3F`) is not a keyword first-char candidate (it's not a letter and not `<`/`>`), so TOKMAIN walks past it. `?` at statement start stores as a raw `0x3F` byte and is rejected by LINESCAN's statement dispatcher as not-a-statement. **[ROM L15877 candidate-char check; UG/TM list no `?` shortcut.]**

- **Reserved/dead keyword table slots are unreachable.** The keyword table has a number of slots whose text is `"-"` (a single byte `0xAD`, i.e. `'-' | 0x80` — a one-letter "keyword" whose only letter is `-`). These slots occupy table indexes 0x49 (INARRAY), 0x4E (NUMBER), 0x51 (USING), 0x52 (SHIFT), 0x68, 0x6A, 0x75, 0x77, 0x78, 0x79, 0x7D, 0x84, 0xD0 (EDIT), 0xF7–0xFE, and 0xFF (which is the special INK→PEN shim, see §3.3). Because per §3.2 only letters and `<`/`>` are candidate first characters, and `-` (0x2D) is none of those, GETTOKEN never even attempts to match these slots. They are **dead** — no input string can ever cause them to fire. (The 0xFF/INK slot is the lone exception, reached only because TOKMAIN's special-case in `TOK42` rewrites the *result* of a successful `INK` match, not because the table's "INK" text is matched against literal `-`.) **Lexer policy: no special handling required.**

---

## 4. Numeric literals

### 4.1 Where they appear
Numeric literals are *not* tokenised by TOKMAIN — they are processed by `LINESCAN` at syntax-check time. Any decimal digit, `&`, `BIN` token, or leading `.` that the expression scanner (`SCANNING`) encounters triggers the literal-number processor at `SDECIMAL` (ROM L5243). (BIN is itself a token, since `BIN` is a keyword; the editor tokenises `BIN` first, then LINESCAN's scanner recognises the `BIN`-token-leader as a binary-literal start.)

### 4.2 Accepted surface forms

| Form | Example | Grammar |
|---|---|---|
| Decimal integer | `0`, `123`, `65535` | `digit+` |
| Decimal with fraction | `1.5`, `0.5`, `.5`, `1.` | `digit+ '.' digit*` or `'.' digit+` |
| Scientific | `1E5`, `1.5e-3`, `1E+5`, `.1E5`, `1.E2` | `<mantissa> ('E'\|'e') ('+'\|'-')? digit+` |
| Hex | `&FF`, `&80FF`, `&12345`, `&0a`, `&0000FFFFFF` | `'&' hexdigit+` (any number of leading zeros; value must fit in 24 bits, `[0-9A-Fa-f]`) |
| Binary | `BIN 10101` | the keyword `BIN` followed by `[01]+` (max 16 bits) |

Notes:
- **Decimal point alone is not a number.** `.` followed by a non-digit is rejected (`NONSENSE` ROM L5637). So `.+` is an error, but `.5` works and so does `1.`.
- **`.E5` is rejected.** No leading digit, no fractional digit. SDECIMAL (ROM L5632) takes the `.` branch, calls `RST 20H` to skip `.`, calls `NUMERIC`; `E` is not numeric, so `JP NC,NONSENSE` fires. **[ROM]**
- **`1E` / `1E+` are rejected.** After the mantissa, `GEXSGN2` (ROM L5695) calls `NUMERIC`; if the next character (after any optional `+`/`-`) is not a digit, `JP NC,NONSENSE` fires. So an `E` must be followed by at least one digit (optionally preceded by a sign). **[ROM]**
- **Exponent underflow is a hard rejection, not silent zero-coercion.** `10 PRINT 1E-300` is rejected by the SAM BASIC editor at line-entry time (SimCoupé empirical 2026-05-14). The ROM's `POFTEN`/exponent handling fires `Number too large` when |exponent| exceeds the representable range, regardless of whether the underflow would round to 0 in IEEE-style FP. **Lexer policy**: reject any literal whose decoded exponent magnitude exceeds the SAM FP range (≈ ±38 decimal, more precisely the biased-exponent byte ≤ 127). Do not coerce to zero.
- **`&` must be followed by at least one hex digit.** Although `AMPDILP` (ROM L18684) itself contains no minimum-digit check (it would compute the value `0` from zero digits), the SAM BASIC editor **rejects `10 &` (ampersand with no hex digits)** at line-entry time — verified on SimCoupé. The rejection comes from LINESCAN treating a bare `&` as not-a-complete-expression (`NONSENSE`). **Lexer policy: reject `&` not immediately followed by `[0-9A-Fa-f]` as a syntax error.**
- **A numeric literal must be terminated by a non-alphabetic character.** Although the byte-level number-parsers (`AMPDILP` for hex at ROM L18684, the decimal-scientific parser at `DECIMAL` ROM 0x1778) walk digits until the first non-digit character and would, at byte level, parse `&FFG` as hex `&FF` followed by leftover `G` (or `1G` as decimal `1` followed by `G`), the SAM BASIC editor **rejects all of `10 PRINT &ffg`, `10 PRINT 1G`, `10 PRINT 1.5G`, `10 PRINT 1E5G`, and `10 PRINT 1G:PRINT "x"` at line-entry time** — verified on SimCoupé (2026-05-14). The rejection comes from LINESCAN's later syntax-check pass, which treats the alphabetic character immediately following the digits as a syntax error (likely `NONSENSE` after the unexpected token). **Lexer policy**: after consuming the digits of any numeric literal (decimal, hex `&…`, scientific `…E…`, float, and binary `BIN …`), peek the next character; if it is `[A-Za-z]` (or `_`), emit a syntax error (`bad number syntax: "1.5G"`). Operators (`+ - * / : , ;`), end-of-line, and end-of-input are all valid terminators. **[ROM L18691–L18695 / `AMPVALID`; SimCoupé empirical 2026-05-14.]**
- **Whitespace between a numeric literal and a following identifier is also rejected** at line-entry time in PRINT-like contexts: `10 PRINT 1 G` is rejected (SimCoupé empirical 2026-05-14). This is a **parser-level** rule (PRINT's expression-list parser requires `,` / `;` / `'` between expressions; an identifier directly after another expression has no valid grammar-rule). Our lexer is intentionally *not* a parser and will tokenise `1 G` as `[number 1, identifier G]` without error. The asymmetry is acceptable because (a) corpus files all passed the editor's parse, so they never contain `1 G` patterns; (b) the corpus round-trip test only validates "bytes we produce match bytes we read", not "we reject everything the editor rejects". A future tightening could add minimal expression-list checking, but is out of v1 scope.
- **Hex case-fold**: `&aB` is `&AB`. (ROM L18689 / `AMPVALID` does `OR 20H`.)
- **Hex range**: AMPERSAND accumulates 4 bits per digit into AHL. The check is **value-based**, not digit-count: only an overflow *into bit 24* triggers `Number too large` (ROM L18702–L18704 / `AMPERLP`). So `&FFFFFF` (16777215) is the maximum *value*, but any number of leading-zero hex digits is allowed: `&0FFFFFF`, `&0000FFFFFF`, etc. all accept and store value 0xFFFFFF (SimCoupé empirical 2026-05-14). Lexer policy: walk all hex digits regardless of count; reject only when the accumulated value would exceed 0xFFFFFF.
- **Binary range**: each `0` or `1` shifts BC (16-bit register pair) left by one. A `1` shifted off the top is *not* checked — wait, actually look at ROM L5621–L5627: `RL C; RL B; JR NC,NXBINDIG; RST 08; DB 28` — the carry test fires only when there's an overflow into bit 17. So up to 16 binary digits work; the 17th causes `Number too large`. **[ROM]**
- **Sign**: a leading `-` is **always a unary operator**, never part of the literal. (ROM L5517 / `UNARMIN` → `LD E,NEGATE`.) So `-5` is parsed as unary-minus applied to `5` (literal `5`). However, the 5-byte invisible form is inserted only for the *literal* `5`. The visible byte stream for `LET A=-5` becomes `LET_TOK ASCII 'A' '=' '-' '5' 0x0E 00 00 05 00 00`. **The lexer must emit no sign byte inside the 5-byte form for `-5`** — only the literal's value.
- **Unary `+`**: skipped silently (ROM L5514 / `UNARPLU`). Affects no bytes.

### 4.3 The 5-byte invisible form (`0x0E` + 5 bytes)

Inserted by LINESCAN *after* the visible textual form of every numeric literal. (ROM L5261–L5269 / `INSERT5B`; `MAKESIX` at L11904 opens the 6 bytes; copies 5 from FPCS.) The marker byte is `0x0E`. (TM p77 "The invisible forms are 0EH followed by 5 bytes to store the value.")

The 5 bytes are the floating-point-calculator (FPC) representation of the value (TM p49):

**Integer fast-path (signed value in `-65535..65535`):**
```
byte 0: 0x00              ; exponent = 0 ⇒ "small integer" sentinel
byte 1: 0x00 if value >= 0, 0xFF if value < 0
byte 2: value & 0xFF      ; LSB of |value| (or of 65536+value if negative — see below)
byte 3: (value >> 8) & 0xFF
byte 4: 0x00
```

The TM is explicit (p49) that "negative values are stored in negated form (65536 minus the number)" — e.g. `-1` is `00 FF FF FF 00` and `-0x80` is `00 FF 80 FF 00`. **[TM, corroborated by ROM `FPBCINT` L6951 onwards.]**

`Number(uint16)` in `tokens.go` already handles the positive-integer case. The lexer also needs the negative case if it ever materialises a negative literal — but per §4.2, the editor never sees a negative literal directly: `-N` is unary-minus over `N`. So the lexer's `Num` only ever needs the positive integer form for integers in `0..65535`.

**Floating-point form (value outside small-int range, or non-integer):**
```
byte 0: exponent + 0x80    ; biased exponent. Mantissa is in (0.5, 1)
byte 1: sign bit (0x80 = negative, 0x00 = positive) OR'd with top 7 bits of mantissa
byte 2-4: next three bytes of mantissa, MSB first
```
The mantissa always has an implicit leading `1` bit (TM p49 "the first bit is always 1, allowing it to be actually used as a SGN bit"). Range about `1E-39` to `1E38`, accuracy 9–10 decimal digits.

Implementing the FP encoder in Go: the simplest approach is to mirror the ROM's `INTTOFP` + `POFTEN` (multiply by 10^exp) routine in floating-point arithmetic, then encode the resulting `float64` into the SAM 5-byte form. For values that are integer-valued and in `-65535..65535`, **emit the small-integer form, not the FP form** — the ROM does this (e.g. `&80` = `00 00 80 00 00`, integer form). This matters for corpus byte-exactness.

#### Range / overflow

- `Number too large` (error 28) fires for:
  - decimal/scientific values whose exponent character is non-numeric or whose exponent's absolute value > 127 (ROM L5701 / `NTLERR`, `FPTOA` returning carry, `RLCA; JR NC` testing top bit);
  - `&` numbers with > 6 hex digits (ROM L18704);
  - `BIN` numbers with > 16 bits (ROM L5625).
- A line of code is limited to **0x3EFF** (≈16127) bytes after tokenisation (`OOMERR` ROM L4123).

#### The `Num.Display` field
The existing Go type `*sambasic.Num` stores both the textual `Display` (the ASCII characters as typed) and the 5-byte `Value`. The lexer must populate both. The `Display` is the original typed bytes including any leading `.`, the optional `E`/`e`, the optional `+`/`-` after `E`, etc. (The leading unary `+`/`-` is *not* part of `Display`; it's a separate `literal` token.)

### 4.4 What is *not* a numeric literal at this level
- Octal: no syntax. (No `&O…`, no leading `0` magic.)
- Other bases: none.
- Constants like `PI`: those are keywords (2-byte token `0xFF 0x3B`), not numeric literals — no 5-byte form is inserted for them.

---

## 5. String literals

(ROM L15896–L15910 / TOKMAIN's quote-scanner; L5529–L5594 / `SQUOTE` evaluator.)

- Delimited by `"` (`0x22`) on both sides.
- **Embedded `"` is escaped by doubling**: `"a""b"` is the 3-character string `a"b`. (ROM L5546–L5549 / `SQUOTE`: checks the char after the closing `"`; if another `"`, copies and continues.)
- **Strings are stored verbatim** between the two delimiting quote bytes in the tokenised body — *including* both `"` bytes. No escape transformation happens at tokenise time; the doubling rule is interpreted only at run time when the string value is needed. So for the lexer: copy the bytes between the opening `"` and the closing `"` (inclusive) directly to the output, including any doubled quotes.
- **Doubled-quote byte storage rule (confirmed).** A `""` doubled-quote sequence inside a string literal represents a single `"` *at run time* but **both `"` bytes are stored verbatim on disk**. Verified on SimCoupé: `PRINT "hello""pete"` outputs `hello"pete` at run time, and the saved program preserves all four quote bytes (opening `"`, `""` pair, closing `"`) in the tokenised body. The `""`→`"` collapse is performed by `SQUOTE` (ROM L5546–L5549) at run time, not at tokenisation. **Lexer policy**: when the source text contains `""` inside a string, emit both `"` bytes; do not collapse.
- A string is **terminated implicitly by `0x0D`** (line end) if no closing `"` is found. **The editor accepts unterminated strings at line-entry time** — the `NONSENSE` check at ROM L5540 is part of `SQUOTE`, the run-time evaluator, not the line-entry tokeniser. Empirical (SimCoupé 2026-05-14): `10 PRINT """` is accepted at entry; body stored is `bb 22 22 22 0d` (PRINT keyword + three literal `"` bytes + terminator). The run-time error "42 String too long" only fires when `PRINT` tries to consume the unbalanced quotes. **Lexer policy**: accept unterminated strings; copy bytes verbatim until the closing `"` *or* end-of-line / end-of-input, whichever comes first.
- **No backslash escapes**, no `\n`, no `\t`.
- **Max length 65520** runtime (TM p55) but at tokenise time the only limit is the 16127-byte line cap.
- **Embedded control chars** (`0x00`–`0x1F`) are not stripped by the editor's character fetch inside string literals because TOKMAIN's quote-scanner reads bytes directly via `LD A,(DE)` (ROM L15900), not via the control-skipping GTCH1 path. So if some byte < 0x20 is in the line's buffer between two quotes, it stays in the stored bytes. (How such a byte gets into the buffer is via the editor's `[CTRL]` / special-key sequences — beyond the lexer's scope.)

**Lexer policy on control bytes inside strings:** if the source text contains a `{N}` escape (the convention used by `samfile basic-to-text`, see §9) inside a string, render it as the corresponding byte `N`. Otherwise treat every byte literally.

---

## 6. Identifiers (variable, array, FN, PROC, LABEL names)

### 6.1 Allowed characters

(ROM `ALPHA` L13621, `ALPHANUM` L13637, `ALNUMUND` L13665, `ALDU` L13653.)

- First character: a letter `A`–`Z` or `a`–`z`. (ROM `ALPHA`, `GETALPH` L13563.)
- Subsequent characters: letter, digit, or underscore `_`. (`ALNUMUND` — bit 5 says "letter or digit or `_`".)
- Spaces are *allowed* inside numeric variable names but **do not count toward the length** and are not stored as part of the variable's identity. (TM p55 "Numbers, letters, underlines and spaces may follow the first character, to a total of up to 32 characters (spaces do not count)".)
- Optional trailing `$` denotes a string-type. (ROM L11772 / `VARNAME` checks for `$`; `FNNAME` L11663.)
- Optional trailing `(` (immediately, no space) denotes array — handled at expression-parsing time, not part of the identifier's stored bytes. (ROM L11793–L11798 / `VARAR`.)

### 6.2 Length limits

- Numeric variable name: **up to 32** characters, spaces not counted (`BC=0B00H` initial in `VARNAME` for 10-char limit is the **string/array** limit — see below; the numeric 32-char limit comes from TM p55 and is enforced at the *variable creation* time via the type/length byte 5-bit field).
- String / string-array / numeric-array variable name: **up to 10** characters, spaces not counted. (ROM L11765 / `VARNAME` initialises `B = 0x0B`, i.e. max name length 10; TM p55.)
- FN, PROC, LABEL names: identifier rules same as numeric variables — letters, digits, `_`. No `$` (well, `DEF FN x$(…)` exists with `$`); no spaces stored. **[ROM L11663 / `FNNAME` — same `GETALPH` + `ALNUMUND` loop.]**

These limits are run-time checks; the **lexer accepts longer identifiers and lets LINESCAN error out later** — that's what the editor itself does.

### 6.3 Storage form

- Identifiers are stored **verbatim** in the line body as raw ASCII bytes (case as typed; spaces as typed). Casing is preserved in the *program text*.
- Case-folding happens only when matching identifier strings to existing entries in the **variable storage** at run time (TM p76: stored "in lower case if they are letters"). That happens during execution, not tokenisation.

### 6.4 The `FN ` keyword prefix
The keyword `FN` is itself a 2-byte token (`0xFF 0x42`). The user types `FN name(args)` — `FN` tokenises, then the identifier `name` follows as literal bytes. (ROM L11663 / `FNNAME`.) Same pattern for `DEF FN name`.

### 6.5 Invisible bytes after a bare-identifier statement (procedure call)

When a statement begins with an identifier that is **not** a recognised command keyword (i.e. the tokenised byte is in the range `0x00–0x8F`), `LINESCAN`'s command dispatcher treats it as a procedure call and falls into `PROCS` (ROM L11809, dispatched from L3577 `JP C,PROCS` in `STMTLP3`). At *syntax-check* time `PROCS` jumps to `PROCSY` (L11868), which:

1. Skips the bare name with `ALNUMUND` (L11868–L11871).
2. Calls `MKCLBF` (L11889) with `A = 0xFD` (set at L11873).
3. `MKCLBF` calls `MAKESIX` (L11904) which opens **6 bytes** at the post-name position via `MAKEROOM` and pre-fills the first byte with `0x0E`; the remaining 5 bytes are **not zeroed** — they retain whatever shifted-down ELINE content sat there.
4. `MKCLBF` then writes `A` (= `0xFD`) into the next **three** bytes (L11893–L11898). After step 4 the buffer reads `0E FD FD FD ?? ??`, where `??` are MAKEROOM leftovers.

Parameter syntax-checking continues at `PCSYL` (L11878) via `SCANNING` for any actual arguments — those arguments produce the **normal** `0x0E + 5-byte FP` form for any numeric literals among them, separately from the buffer described above.

The buffer is created identically for `FN`-call statements with `A = 0xFE` instead (`0E FE FE FE ?? ??`). `LKCALL` (L11254) later distinguishes the two via the D register.

#### Compilation pass — what the last two bytes mean

The `?? ??` bytes are only filled in when the program is **executed** (or when `COMPILE` is otherwise forced — see below). `INSERTLN` sets `COMPFLG = 0xFF` (L4102 via `SCOMP`); the next `CALL COMPILE` at `MAINEXEC` (L3804) walks every procedure-call buffer with `LKCALL` (L11254) and patches each via `LOOKDP` (L11373):

- **`DEF PROC name` found** (label `LKDP3` at L11401):
  - `B = 0x80 | (page & 0x1F)` — page byte with bit 7 as "page valid" marker.
  - `DE = addr of the line-number bytes` of the matched `DEF PROC` line (set by `EX DE,HL` at L11399, after backing 4 bytes from the first text char to reach the line header).
  - Patch at L11414–L11420 overwrites bytes 3, 4, 5 of the calling buffer:
    - byte 3 (was the third `FD`) ← `B` (page byte, top bit set)
    - byte 4 ← `E` (addr low)
    - byte 5 ← `D` (addr high)
  - Stored form: **`0E FD FD <0x80|page> <addr-lo> <addr-hi>`**.
  - Worked example: `0E FD FD 80 F3 9C` decodes as page `0x80 & 0x1F = 0`, address `0x9CF3` in the SAM 8000-BFFFH program area — i.e. an absolute pointer into the program-store page-0 segment, pointing at the line-number-MSB of the matching `DEF PROC` line.
- **No matching `DEF PROC`** (label `LKDP4` at L11409):
  - `B = 0xFF` — the "no DEF PROC" marker.
  - **`D`, `E` are leftover register state** from `LKFC`'s last successful pass through L12088–L12090 (`LD E,(HL); INC HL; LD D,(HL)`) before it bailed at the program-end stopper (`RET C` at L12083). That last pass set DE to the **line-length word of the most recently scanned line** in the program (low byte then high byte; the SAM line header stores length as LSB/MSB per TM p77).
  - Stored form: **`0E FD FD FF <last-line-len-lo> <last-line-len-hi>`**.
  - Worked examples (Pete's tests, 2026-05-14):
    - Single-line program `60 X` (line body `58 0E FD FD FF 08 00 0D` = 8 bytes) → buffer `0E FD FD FF 08 00`. `08 00` matches that line's own length (it's also the last/only line scanned).
    - Single-line program `20  X` (one space kept in body → body length 9) → buffer `0E FD FD FF 09 00`.
    - Multi-line program calling undefined `Y` where `X` is defined via `DEF PROC X` → buffer `0E FD FD FF 02 00`. The `02 00` is the length of the **last** program line scanned by `LKFC` (typically a short `END PROC` line, which tokenises to roughly two bytes plus `0x0D`).

#### Determinism

The bytes are **deterministic given the full tokenised program text** at the moment `COMPILE` runs. Specifically:
- The defined-proc case depends only on the absolute program-store address of the matching `DEF PROC` line, which is determined by the cumulative sizes and order of preceding lines.
- The undefined-proc case depends on the line length of the last full line `LKFC` examined before bailing — which is the **last line in the program**, regardless of where in the program the call appears. (Because `LKFC` scans the *entire* program looking for a `DEF PROC` match each time `LKCALL` finds a calling buffer, and DE retains its value through the failed loop.)

They are **not** FPCS leakage: `INSERT5B` (L5261) is never reached on this code path. The bytes come from `MAKEROOM` (which doesn't zero) and are overwritten only by `LOOKDP`.

#### Pre-COMPILE form (MKCLBF placeholder)

The literal byte pattern that `MKCLBF` writes after the bare identifier is fully determined by ROM code:

- **Byte 0**: `0x0E` — hard-coded in `MAKESIX` (L11906: `LD (HL),0EH`).
- **Bytes 1, 2, 3**: copies of register `A`, which is set by the caller:
  - `PROC`-call (called from `PROCSY` L11873): `LD A,0FDH` → bytes are `FD FD FD`.
  - `FN`-call (called from `FNSYN` L11635): `LD A,0FEH` → bytes are `FE FE FE`.
  - (`DEF FN` parameter buffers go through `MAKESIX` directly at L15821, *not* `MKCLBF`, and are filled by the `DEF FN` body's own logic — out of scope here.)
- **Bytes 4, 5**: **not initialised by MKCLBF**. `MAKESIX` opens 6 bytes via `MAKEROOM` (L11904–L11908); `MAKEROOM` shifts existing memory outward but never zeros the gap, and `MKCLBF` only writes bytes 0–3 before `LD (CHAD),HL` at L11900. So bytes 4–5 hold whatever happened to be at that ELINE offset before the gap was opened.

Authoritative literal patterns immediately post-`MKCLBF`:

| Call kind | Pre-COMPILE buffer | Caller |
|---|---|---|
| `PROC`-call (bare identifier statement) | `0E FD FD FD ?? ??` | `PROCSY` L11868 |
| `FN`-call (`FN name(...)`) | `0E FE FE FE ?? ??` | `FNSYN` L11631 |

Note: the disasm author's comment at L11899 (`BUFFER= 0E FE FE FE ?? ??`) describes only the FN case; the PROC case differs in bytes 1–3.

#### COMPILE re-run gate

**Q: Does SAM BASIC re-run `COMPILE` every time the program runs, or does it trust the disk-stored buffer bytes?**

**A: It re-runs.** `COMPILE` (L12018) gates the full label/PROC/FN pass on `COMPFLG`:

```
33F2 3A405B    LD A,(COMPFLG)
33F5 A7        AND A
33F6 282C      JR Z,COMPILEL   ; skip full pass if COMPFLG=0
```

If `COMPFLG=0` only `ELCOMAL` runs (compiles the edit line). If `COMPFLG≠0` it does the labels pass, then `COMALL` (PROCs+FNs), then `XOR A; LD (COMPFLG),A` (L12057–L12058) to clear the flag.

**`COMPFLG` is set to `0xFF` by `SCOMP` (L11991–L11993).** The comment at L11989 lists SCOMP's callers: **LOAD**, DELETE, KEYIN, RENUM. Plus:

- `INSERTLN` (L4102) — every new program line entered from the editor.
- `DOCOMP` (L12013) — calls `SCOMP` then falls into `COMPILE`. `DOCOMP` is invoked by `CLEAR`/`RUN` (L13176) and by `LDPROG` (the BASIC-program LOAD finaliser, L22699: `DW DOCOMP` immediately after the loaded image is placed and `NVARS`/`NUMEND`/`SAVARS` are restored).

So the gate sequence is:

1. **LOAD a BASIC file from disk** → `LVMMAIN` → `LDFL` → `LDPRDT` → `LDPROG` → `CALL R1OFFCL ; DW DOCOMP` (L22697–L22699). `DOCOMP` forces `COMPFLG=0xFF` then runs the full `COMPILE` pass before LOAD returns.
2. **RUN** → `RUN` (L13143) → `CLR1` → `DOCOMP` (L13176). Forces a full re-COMPILE every time `RUN` is issued, regardless of prior state.
3. **First-RUN-after-LOAD** therefore runs the full COMPILE pass **twice** in quick succession (once at the end of LOAD, once at the start of RUN) — both pass-overs are idempotent because `LOOKDP` (L11373) re-scans the program for each `0E FD FD …` buffer and unconditionally overwrites bytes 3–5 of the buffer with `LD (HL),B / E / D` (L11416–L11420).

**`LKCALL`'s buffer locator does not care about the post-COMPILE state of bytes 3–5.** It searches with `CPIR` for the *second* `FD`/`FE` (byte index 2; preserved across COMPILE), confirms byte 0 is `0E` and byte 1 is `FD`/`FE` (also preserved), and then only sanity-checks that byte 3 has bit 7 set (L11283–L11286). Both the MKCLBF pre-COMPILE value (`FD` or `FE`, bit 7 set) and the post-COMPILE values (`0x80|page` or `0xFF`, bit 7 set) satisfy this — so LKCALL finds the buffer regardless, and LKDP3/LKDP4 then clobber bytes 3–5 with current addresses.

**Consequence**: the on-disk values of bytes 3–5 of any PROC/FN calling buffer are **irrelevant to runtime behaviour after LOAD**. Whatever was saved gets overwritten before the program ever executes.

#### When they're "stable"

A freshly-typed line that has not yet been executed contains `0E FD FD FD ?? ??` where the last two bytes are `MAKEROOM` garbage (whatever ELINE held at that offset before the room was opened). Running the program once causes `COMPILE` to patch every buffer. `INSERTLN` sets `COMPFLG = 0xFF` (L4102) so any subsequent insert/edit forces a full re-compile on next RUN, and `DOCOMP` (called from both `RUN` and `LDPROG`) forces it again unconditionally. **On-disk SAVE files reflect whatever state was last in memory** — typically the post-COMPILE form if the program was RUN before SAVE, or the raw `FD FD FD ?? ??` form if SAVEd immediately after entry. **Either form is functionally equivalent: SAM rebuilds bytes 3–5 on every LOAD via `LDPROG → DOCOMP`.**

#### Lexer policy

**Supersedes the prior "treat 3 bytes as don't-care" recommendation.** Given the COMPILE re-run gate above, a text-to-bytes lexer should emit a **fixed deterministic placeholder** for every bare-identifier statement and every `FN`-call:

| Call kind | Lexer-emitted bytes | Rationale |
|---|---|---|
| `PROC`-call (bare identifier as statement) | `0E FD FD FD 00 00` | Bytes 0–3 match MKCLBF's PROC output exactly; bytes 4–5 are guaranteed-overwritten by LKDP3/LKDP4 on LOAD, so any value works — `00 00` is the conventional "uninitialised" choice. |
| `FN`-call (`FN name(args)`) | `0E FE FE FE 00 00` | Same reasoning, with the FN marker bytes. |

Disks produced by this lexer will not be byte-identical to corpus disks for these 6-byte regions if the corpus disk was SAVEd after a RUN (the corpus will hold post-COMPILE address bytes; the lexer holds `FD FD 00 00`). However, **a lexer-produced disk LOADed into SAM will execute identically to a post-COMPILE corpus disk loaded into SAM**, because:

1. LOAD copies the on-disk bytes verbatim into memory.
2. `LDPROG → DOCOMP` runs the full COMPILE pass before LOAD returns.
3. COMPILE walks the program with `LKCALL`, finds every buffer (since byte 0 = `0E`, bytes 1–2 = `FD FD` or `FE FE`, byte 3 has bit 7 set in both pre- and post-COMPILE forms), and `LKDP3`/`LKDP4` unconditionally overwrites bytes 3–5 with addresses derived from the freshly loaded program layout.

**Verifier recommendations**:

- **Functional equivalence test (preferred for byte-exact corpus testing)**: after LOAD, dump the resident program image and compare to the lexer-emitted image *with bytes 3–5 of every `0E FD FD …` / `0E FE FE …` buffer normalised* (either zeroed or re-run through SAM's own COMPILE). This is the test that actually mirrors what SAM does at runtime.
- **Mask-based byte-exact test**: treat bytes 3–5 of any `0E FD FD ??` / `0E FE FE ??` calling buffer as don't-care when comparing against arbitrary corpus disks. Bytes 0–2 (`0E FD FD` or `0E FE FE`) are deterministic from text and must match.
- **Fix-up pass (optional)**: if a corpus disk's bytes 3–5 must match exactly (e.g. to detect tampering), implement a fix-up pass that walks the emitted program and patches bytes 3–5 the same way `LOOKDP` does. For the defined-PROC case this requires resolving the page+absolute address of the target `DEF PROC` line (computable from the laid-out program); for the undefined-PROC case it requires the LSB/MSB of the last full program line's length word (also computable). The fix-up pass exactly reproduces `LOOKDP`'s logic in L11373–L11422.

**Confidence**: high (≥95%) on the re-run gate itself — `DW DOCOMP` at L22699 is an unconditional call invoked unconditionally on every LOAD-program code path. The remaining 5% uncertainty is whether some edge-case LOAD variant (MERGE, LOAD CODE-as-program, autorun boot) bypasses `LDPROG`. From the disassembly:
- MERGE program goes through `MEPROG` (L22405 → L13007) which calls `MEPRO2` and ultimately the per-line `INSERTLN` path, which itself calls `SCOMP` — so MERGE also forces a re-COMPILE.
- AUTOLOAD on boot uses `ALHK` which routes through the same `LKTH`/`LDFL` chain as user-issued LOAD.

**Empirical confirmation (recommended on SimCoupé)**: build two identical-source disks, one SAVEd before any RUN (bytes 3–5 = `FD ?? ??`) and one SAVEd after RUN (bytes 3–5 = `<page> <lo> <hi>`). Load each, dump memory immediately after LOAD completes, and verify bytes 3–5 in memory are identical between the two cases. Then verify both programs RUN with identical output. If both checks pass, the re-run gate is empirically confirmed and the lexer can safely emit `0E FD FD FD 00 00` for all bare-identifier PROC calls.

---

## 7. Operators and punctuation

| Surface | Stored bytes | Notes |
|---|---|---|
| `+` | `0x2B` | ASCII literal. |
| `-` | `0x2D` | ASCII literal. May be unary or binary depending on context — irrelevant to lexer. |
| `*` | `0x2A` | ASCII. |
| `/` | `0x2F` | ASCII. |
| `^` | `0x5E` | Exponentiation; ASCII. (ROM L5331 / OPPRIORT row "power-of".) |
| `=` | `0x3D` | ASCII. |
| `<` | `0x3C` | ASCII (when not start of `<=`/`<>`). |
| `>` | `0x3E` | ASCII (when not start of `>=`). |
| `<>` | `0xFF 0x81` | Token. |
| `<=` | `0xFF 0x82` | Token. |
| `>=` | `0xFF 0x83` | Token. |
| `,` | `0x2C` | ASCII. |
| `;` | `0x3B` | ASCII. |
| `:` | `0x3A` | Statement separator; ASCII. |
| `(` | `0x28` | ASCII. |
| `)` | `0x29` | ASCII. |
| `#` | `0x23` | Stream prefix; ASCII; not part of any token. |
| `$` | `0x24` | String-suffix on identifier; ASCII. |
| `&` | `0x26` | Hex literal prefix; ASCII (followed immediately by hex digits before LINESCAN converts the *whole* `&xxx` plus 5-byte form). |
| `?` | `0x3F` | No special meaning at tokenise time; stored as ASCII. (Reserved/unused at run time.) |
| `.` | `0x2E` | ASCII (decimal point inside numeric literals; also other roles). |
| `_` | `0x5F` | ASCII; identifier char. |
| `"` | `0x22` | String delimiter. |
| `'` | `0x27` | ASCII; **not** an alternative comment marker — there is no `'`-comment syntax in stock SAM BASIC. **[ROM]** No comment in TOKMAIN about it; UG p30 mentions only `REM`. |
| `~` | `0x7E` | ASCII; no special role. |
| `**` | `0x2A 0x2A` | Not recognised — only `^` is the exponent operator. (ROM has no entry for `**`.) |
| `MOD`, `DIV`, `BOR`, `BAND`, `OR`, `AND`, `NOT` | 2-byte tokens — see keyword table (`0xFF 0x7A`, `…0x7B`, `…0x7C`, `…0x7E`, `…0x7F`, `…0x80`, `…0x76`). | Tokenised by TOKMAIN as normal alphabetic keywords. |

**Non-ASCII operator characters: none.** SAM BASIC v3.0 uses only ASCII operator characters; the only multi-byte operators are the three `<>`/`<=`/`>=` which take a 2-byte keyword form.

---

## 8. REM and comments

(ROM L15977–L15979; UG p30.)

- Only one comment form: the `REM` token (`0xB7`).
- TOKMAIN, after emitting any token, checks whether the byte it just wrote is `0xB7`. If so, **it stops tokenising the rest of the line.** (`CP 0B7H; JP NZ,TOKRST` ROM L15978.) The remainder of the line — up to but not including the line-terminator `0x0D` — is left as **raw ASCII as typed**.
  - That means: `REM Print "hello"` → `REM_TOK 0x20 'P' 'r' 'i' 'n' 't' 0x20 0x22 'h' 'e' 'l' 'l' 'o' 0x22 0x0D`. Keywords inside the comment are NOT tokenised; case is preserved.
- LINESCAN's statement loop, upon seeing the `REM` token at the start of a statement, jumps to `REMARK` (ROM L3651) which discards everything to the next `0x0D`. So no number-form-insertion happens inside a REM either. (`F1 POP AF; JR LINEEND`.)
- **No alternative comment forms**: `'`, `//`, `;`, `~`, `(*…*)` etc. are not comment syntax. **Empirical (SimCoupé 2026-05-14):** the editor **rejects** `10 'hello` (apostrophe at statement-start) and `10 PRINT "x":'comment` (apostrophe after `:`) at line-entry time. Although TOKMAIN walks past `'` (it's not a keyword candidate), LINESCAN's syntax dispatcher errors when `'` is used as a statement-starting token. **Lexer policy**: reject `'` outside string literals and REM bodies as a syntax error (`unexpected character "'"`); inside strings and REM bodies it stays as ASCII `0x27`.
- `REM` doesn't need any space before whatever follows: `REMhello` will be tokenised as `REM` + `hello` (because `REM` ends in a letter and `h` is a letter — wait, that breaks the §3.3 word-boundary rule). **Re-checking**: GTTOK6's trailing-letter check says "if the keyword's last letter is a letter, the next input char must be non-letter". So `REMhello` would *not* tokenise as `REM` + `hello`; it would fall through and tokenise nothing — `REMhello` would be stored as 7 ASCII bytes and at run time be an undefined procedure call. **`REM hello` (with a space) is the canonical form.**

---

## 9. Control-character escapes `{N}`

**The SAM editor has no `{N}` syntax.** This is purely a convention of `samfile basic-to-text` (see `sambasic.go` L93 `case b < 0x20: fmt.Printf("{%v}", int(b))`).

- The editor accepts control bytes 0x00–0x1F into a line buffer only via the keyboard driver (e.g. `[CTRL]+key` sequences, or via `KEYIN`). The character-fetch routine `GTCH1` silently *skips* them during line scanning (ROM L394–L408), so they never participate in tokenisation or syntax checking. Inside string literals they survive because the quote-scanner reads bytes directly (§5).
- **Lexer policy**: accept `{N}` (where `N` is a decimal byte value `0`–`255`) as a single-byte literal escape, emitting the byte directly into the output. Required for round-tripping any corpus file that contains, e.g., an `INVERSE 1` colour-control byte inside a `PRINT` string.
- Recommended scoping: accept `{N}` only inside string literals, after `REM`, and (debatably) anywhere — matching whatever `basic-to-text` emits. Note that `basic-to-text` currently emits `{N}` regardless of context (any non-string control byte will be rendered the same way), so the inverse must accept it anywhere.
- The brace characters `{` (`0x7B`) and `}` (`0x7D`) themselves have **no syntactic meaning** in stock SAM BASIC, so they are unambiguously available as the escape delimiters.

---

## 10. Line-entry validation: what the editor rejects

The editor rejects a line at entry time (i.e. LINESCAN fires an error) for any of the following — the lexer should mirror exactly:

| Rejected at entry | ROM citation |
|---|---|
| Unterminated string (`"` without closing `"` before `0x0D`) | L5540 / `SQUOTE` |
| `&` not followed by at least one hex digit | Verified on SimCoupé: typing `10 &` (ampersand, no hex digits) into the SAM BASIC editor is rejected. Although `AMPDILP` itself contains no minimum-digit check (ROM L18684), LINESCAN treats a bare `&` as an incomplete expression and fires `NONSENSE`. **Lexer policy: reject.** |
| `.` not followed by a digit | L5637 / `NONSENSE` after `JP NC,NONSENSE` (`INSIST ON E.G. .1 OR .8`) |
| `E` (after digits/`.`) not followed by an optional sign and a digit | L5696 / `NONSENSE` after second `JP NC,NONSENSE` |
| Hex literal with > 6 digits overflowing 24-bit accumulator | L18704 / `NTLERR` (error 28 "Number too large") |
| Binary literal with > 16 digits | L5625 (error 28) |
| Scientific-notation exponent magnitude > 127 | L5703 / `NTLERR` |
| Line number `< 1` or `> 0xFEFF` | L4079 / `EVALLINO` "JP C,NONSENSE" via MAINE2 |
| Tokenised line body > 0x3EFF bytes | L4123 / `OOMERR` |
| Single line containing > 127 statements (`SUBPPC` overflows to bit 7) | L3759–L3762 / `MAINE1` — emits error 33 "No room for line" |
| Any structural syntax error caught by `LINESCAN`'s statement dispatcher (wrong keyword for context, missing brackets, undefined token byte, etc.) | various |

The editor **does not** reject things like:
- `LOAD 1.5` — accepted at entry (LINESCAN sees `LOAD <number>`, which is structurally valid; run-time errors when the file can't be found).
- Reference to undefined variables / FN / PROC / LABEL.
- Out-of-range arguments to commands that are not range-checked syntactically.

**Lexer policy**: aim for "accept everything LINESCAN accepts at syntax check time; reject what LINESCAN rejects". Specifically the lexer needs to:
1. Recognise the lexical surface (keywords, numbers, strings, identifiers, operators) correctly.
2. Insert the `0x0E` + 5 bytes after every numeric literal.
3. Validate the few constraints in the table above (line-number range, line-length cap, number-range).

The lexer does **not** need to validate statement-level syntax (`LOAD` must be followed by `"name"`, etc.) — those errors are LINESCAN's job and the editor only rejects the line then. If the goal is "encode anything that was once successfully stored on a SAM disk", a permissive lexer that just performs lexical translation is correct, because anything stored *did* pass LINESCAN at the time it was entered.

---

## 11. Line editing semantics

The SAM BASIC editor treats a typed program as a sequence of **edit operations**, not as a single static program. Each line, on Enter, is dispatched through `INSLN3` (ROM L4106–L4133) and friends, which perform an in-order insert / replace / delete against the program-store linked-list keyed by line number. A lexer consuming source text that is intended to behave "like typing into the SAM editor" must reproduce the *final* state that the editor would have arrived at after replaying the same sequence of Enter-presses.

The behaviours below have been observed in the SAM BASIC editor (SimCoupé). Where a ROM citation is provided it is exact; where the routine is hard to pin down, the cite is marked "(citation TBD)".

### 11.1 Leading whitespace before the line number is ignored

Already covered in §2.3: the character-fetch path (`GTCH3` ROM L394–L408, called from `EVALLINO` ROM L4068+) silently skips bytes in `0x00`–`0x20` before the first significant character. So `    10 PRINT "hi"` parses as line 10. **Lexer policy**: strip any leading whitespace before parsing the line number.

### 11.2 Lines are stored sorted by line number

The editor maintains the program as a linked list **sorted by line number ascending**. `INSLN3` walks the list looking for the insertion point that preserves the sort. Re-entering lines out of typing order does *not* keep them in typing order. Input:

```
10 X
5 Y
15 Z
```

is stored on disk as:

```
5 Y
10 X
15 Z
```

(ROM L4106+ / `INSLN3`; see also TM p77 "lines are held in order of line number".)

**Lexer policy**: parse all lines first, then emit them sorted by line number. Do not preserve typing order.

### 11.3 Duplicate line numbers: last entry wins

If the typed input contains two lines with the same number, the second replaces the first. The editor's `INSLN3` flow detects an existing line with the matching number, unlinks/reclaims it, and inserts the new content. Input:

```
10 PRINT "first"
10 PRINT "second"
```

stores as line 10 with body `PRINT "second"`. (ROM L4106–L4133 / `INSLN3`; the "found existing line" branch falls through to `RECLAIM` then `INSERTLN`.)

**Lexer policy**: when collecting lines, key by line number; later occurrences overwrite earlier ones.

### 11.4 Bare line number deletes

A line consisting of **only** a line number followed by the line terminator (no body, no whitespace) is a **delete operation**: any existing line with that number is removed from the program, and no new line is inserted. After stripping the line-number prefix, if the only remaining byte is `0x0D`, `INSLN3` jumps past the `INSERTLN` action — equivalent to "delete-if-exists, then no-op". (ROM L4131–L4133 / `INSLN3` — already noted in §2.5.)

So:

```
10 PRINT "hi"
10
```

stores no line 10. The second line is a delete that removes the first.

**Lexer policy**: a bare-line-number input (digits + `0x0D` only, no body characters at all) means "drop any previously-collected line for this number; do not emit this one".

### 11.5 Post-line-number byte handling: the conditional one-space drop

The boundary between "delete this line number" and "store with this body" is decided in `INSERTLN` at ROM L4106–L4116, immediately before the `INSLN3` label. The routine inspects the first one or two bytes after the parsed line number and applies a conditional one-space drop — see §2.3 for the full citation. Restated as cases:

| Input | Stored? | Body bytes (before the `0x0D` terminator) | ROM path |
|---|---|---|---|
| `10\n` (line number, terminator only) | **No (delete)** | n/a — line text length is 1 (just `0x0D`), so `INSERTLN` reaches L4131–4133 `LD A,C; DEC A; OR B; RET Z` and returns without storing. | L4108 `JR NZ,INSLN3` (no space) → L4133 `RET Z` |
| `10 \n` (line number, one space, terminator) | **Yes (stored)** | `0x20 0x0D` (single space + terminator). **Empirically confirmed on SimCoupé.** | L4116 `DEC BC` (space-then-CR special case) → space preserved |
| `10   \n` (line number, three spaces, terminator) | **Yes (stored)** | `0x20 0x20 0x0D` (two spaces + terminator). The first space is dropped because the byte after it is another space (not CR), so the L4116 special case does not fire; one of the three input spaces is consumed. **Re-verify on SimCoupé.** | L4113 `JR NZ,INSLN2` (space-then-not-CR) → one-space drop |
| `10 X\n` | **Yes (stored)** | `0x58 0x0D` (`X` + terminator). The space between line number and `X` is dropped. The visual gap in a `LIST` of this line comes from the LIST formatter's own post-line-number space (ROM L26009–L26010), not from a stored body byte. **Confirmed on SimCoupé 2026-05-14** (line `40 PRINT "with space"` body length 14, no leading `0x20`). | L4113 `JR NZ,INSLN2` → one-space drop |
| `    10\n` (leading whitespace, line number, terminator) | **No (delete)** | leading whitespace stripped before line-number parse (via `GTCH3`); what remains is `10\n` — a bare-number delete. **Empirically confirmed on SimCoupé.** | L394–L408 `GTCH3` skip → then `RET Z` as `10\n` row above |

So the precise rule the lexer must implement is the §2.3 algorithm: strip leading whitespace, parse the line number, look at the first body byte. If body is empty → delete. If body starts `0x20 0x0D` → keep both bytes. If body starts `0x20 <not-0x0D>` → drop the leading `0x20`. Otherwise → keep verbatim.

**Lexer policy**: a bare-line-number input (digits + `0x0D` only, no body characters at all) means "drop any previously-collected line for this number; do not emit this one". A line-number-followed-by-content (including just a single space) stores. When storing, apply the conditional one-space drop above so that re-typing `10 X` after editing `10 X` does not accumulate extra spaces.

### 11.6 Implementation summary for the lexer

Given an input file as a sequence of source lines:

1. Parse each line as `<optional leading whitespace> <line number> [body] <terminator>`. Leading whitespace before the line number is discarded; nothing else is.
2. If the body is empty (i.e. zero bytes between the parsed line number and the `0x0D` terminator), this is a **delete**: remove any previously-collected entry for this line number; emit nothing.
3. Otherwise: apply the conditional one-space drop from §2.3 / §11.5 to the body (skip the first byte only if it is `0x20` and the next byte is not `0x0D`), then **upsert** the parsed line into a map keyed by line number; later entries overwrite earlier ones.
4. After consuming the whole input, emit the surviving lines in ascending order of line number.

This matches the final state the SAM BASIC editor would have arrived at after typing the same input as a sequence of Enter-presses. The pre-dispatch routine that decides delete-vs-store and applies the one-space drop is `INSERTLN` ROM L4096–L4133 (label `INSLN3` at L4120 is the merge point after the space-handling block; the empty-body delete short-circuit is at L4131–L4133).

---

## 12. Edge cases and gotchas

1. **`IF`/`ELSE` token patching by LINESCAN.** As detailed in §3.3, both `IF` and `ELSE` always tokenise to their "long" forms at TOKMAIN time, and LINESCAN later flips them in place based on `THEN`-presence within the line. The lexer **must** implement this patch step; otherwise corpus round-trip fails for every line containing `IF … THEN`. The ROM's exact condition (L6340–L6364) is "the character immediately after the IF's expression is `THEN`" — not just "`THEN` appears somewhere on the line". For the lexer this distinction normally collapses to the same answer because input that's well-formed enough to reach LINESCAN typically has THEN directly after the expression. The simple algorithm is:
   - After full tokenisation of a line, scan for `LIF` (0xD7). If a `THEN` (0x8D) appears later on the same line before any `0x0D` or `:` (`0x3A`) statement-separator, patch the `LIF` to `SIF` (0xD8). Iterate for multiple `IF…THEN`s on the same line.
   - After patching all `IF`s, scan for `LELSE` (0xD9). If a *preceding* `SIF` (0xD8) is on the same line, patch the `LELSE` to `ELSE` (0xDA). (ROM L6447 / `NLELS`.) A second patch (ROM L6438–L6440) rewrites the `LIF` immediately following such a `LELSE` to `SIF` as well, so `ELSE IF cond THEN …` always tokenises with both short forms.

2. **`INK` → `PEN` rewrite** (§3.3). `INK` (and `ink`, etc.) tokenises directly to the `PEN` token byte `0xA1`, not to the table's `INK` slot at `0xFF`.

3. **First statement of a line is special only at run time.** The tokeniser treats every statement-position identically. No "first token must be a command" check at tokenise time — that's a LINESCAN constraint that the editor enforces but for our purposes only matters when reproducing a corpus that already passed that check.

4. **Trailing-space consumption (§3.5)**: when emitting a token, the editor consumes one trailing space if present. Implementation detail of the lexer's source-to-token loop:
   - Buffer the matched keyword's input span (start..end).
   - Emit token bytes.
   - If `source[end+1] == ' '`, advance the source cursor past that one space.

5. **Leading-space consumption**: similarly, if the matched keyword had `' '` immediately before its first letter in the input, that space is *overwritten* by the token, not preserved. (ROM L15952–L15958 / `TOK43`/`TOK5`.) For a lexer this means: look back at the previously-emitted byte; if it was a space *and* the previous-previous byte was not a space, drop it before emitting the token. Or, more simply: when about to emit a token, if the just-emitted last byte is `0x20`, replace it with the token's first byte. (Behaviour exactly matches: keyword recognised at column N replaces the space at column N-1 with its first token byte, deleting it from the line.)

6. **`AUTO` line-number generation** — this is run-time only and not part of the lexer.

7. **`0x0E` byte inside a number's `Display` field**: the lexer must be careful never to include `0x0E` in the display bytes — but since `Display` is composed of ASCII characters from the user's typing, this should not arise.

8. **The `0xFF` end-of-program sentinel**: this is the file-level sentinel after the last line, not a per-line concern. (TM p77 "The final line in the program is followed by FFH".) Already handled by `File.ProgBytes()`.

9. **Tokenisation of an already-tokenised buffer**: TOKMAIN is also called by `KEYIN` (ROM L15398), `MERGE`, and `VAL`. In those cases the input buffer may already contain token bytes. TOKMAIN's main loop handles this gracefully — it skips past any byte that isn't a letter or `<`/`>` candidate, and the `0xFF` arm (L15891) walks past the FN-leader+code pair. For our purposes (lexing fresh ASCII from disk-text format), this doesn't matter; the input is always plain ASCII.

10. **`PI`, `RND`, etc. are 2-byte tokens, not numeric literals.** Their bytes (`0xFF 0x3B`, `0xFF 0x3C`, …) are emitted by TOKMAIN at tokenisation time; LINESCAN does **not** insert a 5-byte form after them. So `LET A=PI` becomes `LET_TOK ' ' 'A' '=' 0xFF 0x3B 0x0D` — five tokenised bytes plus terminator, no `0x0E`.

11. **Operator priorities are run-time concerns.** Lexer should not collapse `A AND B AND C` into anything special — emit each token as encountered.

12. **`SCREEN` vs `SCREEN$`**: these are *different* keywords. `SCREEN` is 1-byte token `0xE7`; `SCREEN$` is 2-byte token `0xFF 0x4C`. TOKMAIN's GETTOKEN handles this correctly because `SCREEN` (table index 0xE7) appears before `SCREEN$` (index 0x4C) in the table — wait, that's wrong order. Looking again: the table is ordered by **token byte ascending**, so `SCREEN$` (0x4C) comes before `SCREEN` (0xE7). The first match wins, so `SCREEN$abc` would try `SCREEN$` first, succeed, and not even look at `SCREEN`. For `SCREENabc`, `SCREEN$` fails (no `$`), then `SCREEN` matches but fails the trailing-letter rule (`a` follows), so neither tokenises. For `SCREEN 1`, `SCREEN$` fails (no `$` after `SCREEN`), then `SCREEN` matches (followed by space), succeeding. Behaviour as expected. Similarly for the many other `$`-suffix duplicates: `INKEY`/`INKEY$`, `MEM`/`MEM$`, `PATH$` (no non-$ form in stock), `STRING$`, `VAL`/`VAL$`, `TRUNC$`, `CHR$`, `STR$`, `BIN$`, `HEX$`, `USR`/`USR$`.

13. **`POINT` is 2-byte (`0xFF 0x3D`) and is a function**; **`PLOT` is 1-byte (`0x9B`)** and is a command. Both share an obvious surface relationship but have different tokenisations. No special handling — just keyword table lookup.

14. **Reserved/unused keyword table slots** (those marked `""` or `"-"` in `keywords.go`) — these indexes are *never matched* by GETTOKEN (the table entries `DC "-"` would match a literal `-`, but `-` is not a candidate first character, so they're effectively dead). No special handling needed.

---

## 13. Open questions / ambiguities

The remaining open questions — the ones not yet answered by ROM reading, SimCoupé testing, or this revision pass.

1. **`'` (apostrophe) at start of statement.** **Resolved (SimCoupé empirical 2026-05-14):** the editor rejects `10 'hello` and `10 PRINT "x":'comment` at line-entry. TOKMAIN walks past `'` (it's not a keyword candidate), but LINESCAN errors when `'` appears as a statement-starting token. See §8 for the lexer policy. **Remaining sub-question** (corpus only): do any extant corpus files have `'` at statement start anyway — e.g. from a dialect-extension ROM that *did* treat it as a remark? If yes, those will surface as `unexpected character "'"` errors during corpus testing and we'll have data to inform an opt-in dialect mode.

2. **`{N}` outside strings — corpus evidence.** The inverse-direction lexer accepts `{N}` as a single-byte literal escape (§9). `{N}` outside a quoted region would inject a raw byte into the post-tokenisation stream — which is fine in principle but unusual: TOKMAIN's main loop normally produces a body containing token bytes, ASCII printables, `0x0E` (number marker), `0x0D` (line terminator), `0xFF` (2-byte-token leader). **Open question**: do any real corpus files contain raw bytes in `0x00`–`0x1F` (other than `0x0D`, `0x0E`) outside strings? If not, `{N}` outside strings can be treated as a round-trip-only feature with no live corpus traffic, and the lexer can reject it as a sanity check. If yes, the lexer must support them as-typed.

---

## Appendix A: file-level format reminder

Per `sambasic/file.go`:

```
Program bytes = concat( Line[0], Line[1], …, Line[N-1] ) ++ 0xFF
Line bytes    = [ MSB(num) LSB(num) Lo(len) Hi(len) ] ++ body ++ 0x0D
body          = sequence of tokens emitted by the lexer
```

The 0xFF after the last line is the program-end sentinel. The line-length field is the byte-count of `body ++ 0x0D` (i.e. it includes the line's own terminator but not the 4-byte header).

After `0xFF`, the file continues with `NumericVars`, `Gap`, and `StringArrayVars` regions — these are run-time state and not produced by the lexer; the lexer's job ends at producing `body` bytes for each line.
