#!/usr/bin/env python3
"""Extract the 169-byte E-Code frequency table from jukebox.bas.

jukebox.bas line 70 does `FOR n=33486 TO 33654: READ a: POKE n,a` — reading the
first 169 values from the DATA statements (lines 80..240) and poking them into
the decrunched player at linear &82CE (page 1 + 718). This script reproduces
that table as out/freqtable.bin so the headless player gets the same notes.

Usage: gen-freqtable.py [path/to/jukebox.bas] [out/freqtable.bin]
"""
import re
import sys

src = sys.argv[1] if len(sys.argv) > 1 else "/Users/pmoore/git/sam-jukebox/src/jukebox.bas"
dst = sys.argv[2] if len(sys.argv) > 2 else "out/freqtable.bin"

vals = []
for ln in open(src):
    m = re.match(r"\s*(\d+)\s+DATA\s+(.*)", ln)
    if not m:
        continue
    if 80 <= int(m.group(1)) <= 240:
        for tok in m.group(2).split(","):
            tok = tok.strip()
            if tok.lstrip("-").isdigit():
                vals.append(int(tok) & 0xFF)

table = bytes(vals[:169])  # 33486..33654 inclusive
assert len(table) == 169, f"expected 169 values, got {len(table)}"
open(dst, "wb").write(table)
print(f"wrote {len(table)} bytes to {dst}")
