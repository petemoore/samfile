#!/usr/bin/env python3
"""Extract slot 0's file body from a corpus disk, for disassembly.

Usage:
    extract_dos.py <hash-prefix>     # find first disk with that SHA
                                     # prefix, write its slot-0 body
                                     # to /tmp/<hash>.bin
    extract_dos.py --disk <path>     # extract slot-0 body from a
                                     # specific disk image

Both write the trimmed body bytes (trailing 0x00 / 0xFF stripped)
suitable for disassembly with any z80 disassembler.
"""

from __future__ import annotations

import argparse
import hashlib
from pathlib import Path

CORPUS = Path.home() / "sam-corpus" / "disks"
DISK_SIZE = 819_200
ENTRY_SIZE = 256


def sector_offset(track: int, sector: int) -> int:
    return ((track >> 7) * 5120
            + (sector - 1) * 512
            + (track & 0x7F) * 10240)


def walk_chain(disk: bytes, first_track: int, first_sector: int) -> bytes | None:
    out = bytearray()
    t, s = first_track, first_sector
    visited = set()
    while True:
        if not (4 <= t <= 79 or 128 <= t <= 207) or not (1 <= s <= 10):
            return None
        if (t, s) in visited:
            return None
        visited.add((t, s))
        off = sector_offset(t, s)
        if off + 512 > len(disk):
            return None
        sd = disk[off:off + 512]
        out.extend(sd[:510])
        nt, ns = sd[510], sd[511]
        if nt == 0 and ns == 0:
            break
        t, s = nt, ns
        if len(visited) > 1600:
            return None
    return bytes(out)


def extract_one(disk_path: Path) -> tuple[str, bytes]:
    data = disk_path.read_bytes()
    if len(data) != DISK_SIZE:
        raise SystemExit(f"{disk_path}: not a 819200-byte MGT image")
    slot = data[:ENTRY_SIZE]
    if slot[0] == 0:
        raise SystemExit(f"{disk_path}: slot 0 is erased — no DOS")
    body = walk_chain(data, slot[0x0D], slot[0x0E])
    if body is None:
        raise SystemExit(f"{disk_path}: slot 0 chain is malformed")
    trimmed = body.rstrip(b"\x00\xFF")
    return hashlib.sha256(trimmed).hexdigest()[:16], trimmed


def main() -> None:
    ap = argparse.ArgumentParser()
    g = ap.add_mutually_exclusive_group(required=True)
    g.add_argument("hash_prefix", nargs="?", help="SHA prefix to find in corpus")
    g.add_argument("--disk", type=Path, help="specific disk image to extract from")
    args = ap.parse_args()

    if args.disk:
        h, body = extract_one(args.disk)
        out = Path(f"/tmp/dos-{h}.bin")
        out.write_bytes(body)
        print(f"{args.disk.name}: sha={h}, {len(body)} bytes → {out}")
        return

    prefix = args.hash_prefix.lower()
    for disk in sorted(CORPUS.glob("*.mgt")):
        try:
            h, body = extract_one(disk)
        except SystemExit:
            continue
        if h.startswith(prefix):
            out = Path(f"/tmp/dos-{h}.bin")
            out.write_bytes(body)
            print(f"{disk.name}: sha={h}, {len(body)} bytes → {out}")
            return
    raise SystemExit(f"no corpus disk has slot-0 body matching SHA prefix {prefix!r}")


if __name__ == "__main__":
    main()
