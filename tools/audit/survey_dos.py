#!/usr/bin/env python3
"""Survey every disk in ~/sam-corpus/disks/ for its on-disk DOS.

For each disk, slot 0's file body is (by convention) the DOS that
boots that disk. Hash the body, group by hash, count distinct disks.
Output a per-DOS catalog: hash, disk count, filename, body length,
sample disks. Cross-reference against known SAM-DOS variants via
ASCII fingerprinting (strings present in the body).

Output: writes docs/dos-catalog.md in the samfile repo and a brief
summary to stdout. Does NOT touch findings.db — this is a one-shot
survey.
"""

from __future__ import annotations

import hashlib
from pathlib import Path

CORPUS = Path.home() / "sam-corpus" / "disks"
REPO = Path.home() / "git" / "samfile"
OUT = REPO / "docs" / "dos-catalog.md"

DISK_SIZE = 819_200
DIR_TRACKS = [0, 10240, 20480, 30720]  # byte offsets of side-0 tracks 0..3
ENTRY_SIZE = 256


def sector_offset(track: int, sector: int) -> int:
    return ((track >> 7) * 5120
            + (sector - 1) * 512
            + (track & 0x7F) * 10240)


def walk_chain(disk: bytes, first_track: int, first_sector: int) -> bytes | None:
    """Return body bytes by walking the sector chain. Returns None on
    any malformed step (out-of-range track/sector, missing terminator
    after 1600 steps to defend against cycles)."""
    out = bytearray()
    t, s = first_track, first_sector
    visited = set()
    while True:
        if not (4 <= t <= 79 or 128 <= t <= 207) or not (1 <= s <= 10):
            return None
        if (t, s) in visited:
            return None  # cycle
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


def slot_zero_body(disk: bytes) -> bytes | None:
    """Return the body bytes of slot 0's file by walking its sector
    chain, or None if slot 0 is erased or the chain is malformed.
    Slot 0's body is the loaded boot code — the actual on-disk DOS.
    Filename / type byte are dir metadata, not part of the ROM
    contract; the survey deliberately doesn't look at them."""
    slot = disk[:ENTRY_SIZE]
    if slot[0] == 0:
        return None  # erased
    return walk_chain(disk, slot[0x0D], slot[0x0E])


# The only fingerprint is the SHA-256 of the on-disk slot-0 body.
# Two bodies with the same SHA are the same DOS; two bodies with
# different SHAs are different DOSes. Identification of which DOS
# each SHA corresponds to is a manual step (disassembly / source
# comparison), not something this survey tries to guess.


def main() -> None:
    by_hash: dict[str, dict] = {}
    no_dos: list[str] = []
    bad: list[str] = []
    total = 0
    disks = sorted(CORPUS.glob("*.mgt"))
    for path in disks:
        total += 1
        data = path.read_bytes()
        if len(data) != DISK_SIZE:
            bad.append(path.stem)
            continue
        body = slot_zero_body(data)
        if body is None:
            no_dos.append(path.stem)
            continue
        # Trim trailing 0x00 / 0xFF padding (last sector usually padded).
        trimmed = body.rstrip(b"\x00\xFF")
        h = hashlib.sha256(trimmed).hexdigest()[:16]
        if h not in by_hash:
            by_hash[h] = {
                "hash": h,
                "body_len": len(trimmed),
                "raw_body_len": len(body),
                "disks": [],
            }
        by_hash[h]["disks"].append(path.stem)

    rows = sorted(by_hash.values(), key=lambda r: -len(r["disks"]))

    md = [
        "# DOS catalog",
        "",
        "Empirical survey of the SAM Coupé corpus at `~/sam-corpus/disks/`.",
        "For each disk, slot 0's file body — the loaded boot code, i.e.",
        "the actual on-disk DOS — is read by walking its sector chain,",
        "then SHA-256-hashed after stripping trailing 0x00 / 0xFF padding.",
        "Two bodies with the same SHA are the same DOS; identifying which",
        "DOS each SHA corresponds to is a manual step (disassembly /",
        "source comparison), not something this survey tries to guess.",
        "",
        "The slot-0 filename and the dir-entry file-type byte are dir",
        "metadata, not part of the ROM ↔ DOS contract, so the survey",
        "deliberately ignores them.",
        "",
        f"- Total disks scanned: **{total}**",
        f"- Disks with slot-0 DOS: **{sum(len(r['disks']) for r in rows)}**",
        f"- Disks with no slot-0 file (non-bootable archives): **{len(no_dos)}**",
        f"- Unique DOS bodies (by SHA-256, truncated to 16 hex chars): **{len(rows)}**",
        f"- Disks with malformed images: **{len(bad)}**",
        "",
        "## Unique DOSes (most-used first)",
        "",
        "| SHA-256 (16) | Disks | Body length (trimmed) | Body length (raw) |",
        "|---|---:|---:|---:|",
    ]
    for r in rows:
        md.append(
            f"| `{r['hash']}` | {len(r['disks'])} | {r['body_len']} | {r['raw_body_len']} |"
        )
    md.append("")
    md.append("## Sample disks per DOS")
    md.append("")
    md.append("Sample list of corpus disks for each unique DOS (up to 5).")
    md.append("Use any of these to extract the body for disassembly:")
    md.append("")
    md.append("```bash")
    md.append("# Extract slot 0's body for disassembly:")
    md.append("python3 ~/git/samfile/tools/audit/extract_dos.py <hash-prefix>")
    md.append("```")
    md.append("")
    for r in rows:
        md.append(f"### `{r['hash']}` ({len(r['disks'])} disks, body={r['body_len']} bytes)")
        md.append("")
        for d in r["disks"][:5]:
            md.append(f"- {d}")
        md.append("")

    if no_dos:
        md.append(f"## Disks with no slot-0 file ({len(no_dos)})")
        md.append("")
        md.append("Non-bootable archive disks. Dir entries were written by")
        md.append("*some* DOS but the DOS isn't on the disk.")
        md.append("")
        md.append("Sample (first 20):")
        for d in no_dos[:20]:
            md.append(f"- {d}")
        md.append("")

    if bad:
        md.append(f"## Malformed images ({len(bad)})")
        md.append("")
        for d in bad:
            md.append(f"- {d}")
        md.append("")

    OUT.write_text("\n".join(md) + "\n")
    print(f"wrote {OUT}")
    print(f"scanned: {total} disks; unique DOSes: {len(rows)}; no-DOS: {len(no_dos)}; bad: {len(bad)}")
    print()
    print("top 15 by disk count:")
    for r in rows[:15]:
        print(f"  {len(r['disks']):4d}  {r['hash']}  body={r['body_len']} bytes")


if __name__ == "__main__":
    main()
