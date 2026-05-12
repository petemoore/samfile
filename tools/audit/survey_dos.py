#!/usr/bin/env python3
"""Survey every disk in ~/sam-corpus/disks/ for its on-disk DOS.

The fingerprint follows the ROM `BOOTEX` routine
(rom-disasm:20473-20598, address 0xD8E5):

  ROM reads exactly track-4-sector-1 (512 bytes) into 0x8000-0x81FF
  via the INI loop at RSA5 (line 20554), validates 4 bytes at
  0x8100..0x8103 against the masked literal "BOOT" (BTWD, mask 0x5F,
  line 20586-20596), then jumps to 0x8009 (line 20598). Everything
  beyond is up to the DOS bootstrap.

So the strict ROM-contract fingerprint of a disk's DOS is the
SHA-256 of its T4S1 sector. By convention, the bootstrap then
loads slot 0's file body (the rest of the DOS), so we also record
that as a secondary fingerprint — its scope is a DOS-side
convention, not a ROM-enforced fact.

Output: writes docs/dos-catalog.md and a brief summary to stdout.
Does NOT touch findings.db.
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


def t4s1(disk: bytes) -> bytes:
    """Return the 512 bytes ROM loads at boot: track 4 sector 1.

    This is the strict ROM-contract identifier — exactly what ROM
    BOOTEX copies into 0x8000..0x81FF before jumping to 0x8009. Any
    disk whose T4S1 SHA matches another's is running the same
    boot-loaded code as far as ROM is concerned. What the
    bootstrap then chooses to do (typically: load slot 0's file
    body) is per-DOS, not ROM-mandated.
    """
    off = sector_offset(4, 1)
    return disk[off:off + 512]


# ROM BOOTEX checks T4S1 bytes 256..259 against the literal "BOOT"
# (BTWD at rom-disasm:26919) using mask 0x5F (BTCK at rom-disasm:
# 20586-20596 — ignores bit 5 (ASCII case) and bit 7 (high-bit) so
# "BOOT", "boot", "BOOt", `42 4F 4F D4` etc. all match). Any
# mismatch → "NO DOS" error and ROM refuses to boot. This is the
# exact ROM bootability gate.
BOOT_LITERAL = b"BOOT"
BOOT_OFFSET = 256


def is_bootable(disk: bytes) -> bool:
    sector = t4s1(disk)
    return all(
        (sector[BOOT_OFFSET + i] ^ BOOT_LITERAL[i]) & 0x5F == 0
        for i in range(4)
    )


def slot_zero_dos_code(disk: bytes) -> bytes | None:
    """Return the byte range slot 0's bootstrap convention would
    load as the rest of the DOS: walk slot 0's chain, decode the
    9-byte file body header, take exactly LengthMod16K|(Pages<<14)
    bytes after the header. Returns None if slot 0 is erased or
    the chain is malformed.

    This is *not* part of the ROM contract — the ROM only loads
    T4S1. The bootstrap at 0x8009 is per-DOS code and is free to
    fetch the rest of the DOS from anywhere on disk. Slot 0 is
    just the convention that every known SAM DOS happens to use.
    """
    slot = disk[:ENTRY_SIZE]
    if slot[0] == 0:
        return None  # erased
    chain = walk_chain(disk, slot[0x0D], slot[0x0E])
    if chain is None or len(chain) < 9:
        return None
    length_mod16k = chain[1] | (chain[2] << 8)
    pages = chain[7]
    length = (length_mod16k & 0x3FFF) | (pages << 14)
    if length == 0 or 9 + length > len(chain):
        return None
    return chain[9 : 9 + length]


# The only fingerprint is the SHA-256 of the on-disk slot-0 body.
# Two bodies with the same SHA are the same DOS; two bodies with
# different SHAs are different DOSes. Identification of which DOS
# each SHA corresponds to is a manual step (disassembly / source
# comparison), not something this survey tries to guess.


# ROM constants for the BOOT path (rom-disasm:20473-20598).
ROM_T4S1_LENGTH = 512
ROM_T4S1_LOAD_ADDRESS = 0x8000   # T4S1 → 0x8000..0x81FF, line 20550 (LD HL,0x8000)
ROM_T4S1_EXEC_ADDRESS = 0x8009   # JP 0x8009, line 20598


def decode_slot0_loadexec(disk: bytes) -> tuple[int, int | None, int, bytes] | None:
    """Decode slot 0's body-header to recover the load address, the
    execution address (or None if disabled), the body byte count,
    and the body bytes themselves. Returns None if slot 0 is
    erased / malformed / shorter than its header claims.

    This is the *bootstrap-convention* fingerprint: what slot 0's
    file header tells whoever is consuming it (typically the
    bootstrap loaded by ROM into 0x8000) about where this code
    wants to live and run. ROM itself never reads this header —
    that's part of the per-DOS convention layered on top of the
    ROM contract.

    Layout (matches samfile.FileHeader):
      [0]   Type
      [1-2] LengthMod16K   LE
      [3-4] PageOffset     LE   — low 14 bits = offset in page; high
                                  2 bits = SAM REL PAGE FORM page tag
      [5]   ExecutionAddressDiv16K
      [6]   ExecutionAddressMod16KLo
      [7]   Pages
      [8]   StartPage      — low 5 bits + 1 = SAM page (1..32)
    """
    slot = disk[:ENTRY_SIZE]
    if slot[0] == 0:
        return None
    chain = walk_chain(disk, slot[0x0D], slot[0x0E])
    if chain is None or len(chain) < 9:
        return None
    length_mod16k = chain[1] | (chain[2] << 8)
    page_offset = chain[3] | (chain[4] << 8)
    exec_div16k = chain[5]
    exec_mod16k_lo = chain[6]
    pages = chain[7]
    start_page = chain[8]
    length = (length_mod16k & 0x3FFF) | (pages << 14)
    if length == 0 or 9 + length > len(chain):
        return None
    body = chain[9 : 9 + length]
    # Linear load address: (StartPage low-5 bits + 1) is the page
    # 1..32, PageOffset low 14 bits is the offset within that page.
    load_address = ((start_page & 0x1F) + 1) * 0x4000 + (page_offset & 0x3FFF)
    # Execution: ExecDiv16K == 0xFF means "no auto-exec from this
    # header" — the body's entry point is implicit (= load_address)
    # or comes from a different mechanism. Only emit a distinct exec
    # address when the body header genuinely declares one.
    exec_address: int | None = None
    if exec_div16k != 0xFF:
        exec_address = exec_div16k * 0x4000 + exec_mod16k_lo
    return load_address, exec_address, length, body


def main() -> None:
    rom_by_hash: dict[str, dict] = {}
    boot_by_hash: dict[str, dict] = {}
    not_bootable: list[str] = []   # ROM would print "NO DOS" — skip
    no_slot0: list[str] = []       # bootable per ROM but slot 0 isn't usable
    bad: list[str] = []
    total = 0
    disks = sorted(CORPUS.glob("*.mgt"))
    for path in disks:
        total += 1
        data = path.read_bytes()
        if len(data) != DISK_SIZE:
            bad.append(path.stem)
            continue

        # ROM bootability gate (BTCK loop at rom-disasm:20586-20596).
        # If a disk fails this check, ROM refuses to boot it and the
        # T4S1 / slot-0 contents are irrelevant for DOS classification.
        if not is_bootable(data):
            not_bootable.append(path.stem)
            continue

        # ROM-contract fingerprint: T4S1 (512 bytes ROM loads to 0x8000,
        # jumps to 0x8009). Load and exec addresses are ROM constants —
        # part of the fingerprint so it's fully self-describing.
        sector = t4s1(data)
        rom_h = hashlib.sha256(sector).hexdigest()[:16]
        rom_row = rom_by_hash.setdefault(rom_h, {
            "hash": rom_h,
            "length": ROM_T4S1_LENGTH,
            "load_address": ROM_T4S1_LOAD_ADDRESS,
            "exec_address": ROM_T4S1_EXEC_ADDRESS,
            "disks": [],
        })
        rom_row["disks"].append(path.stem)

        # Bootstrap-convention fingerprint: slot 0's file-header-declared
        # load / exec / length, hashing exactly that many body bytes.
        info = decode_slot0_loadexec(data)
        if info is None:
            no_slot0.append(path.stem)
            continue
        load_addr, exec_addr, length, body = info
        boot_h = hashlib.sha256(body).hexdigest()[:16]
        boot_row = boot_by_hash.setdefault(boot_h, {
            "hash": boot_h,
            "length": length,
            "load_address": load_addr,
            "exec_address": exec_addr,
            "disks": [],
        })
        boot_row["disks"].append(path.stem)

    rom_rows = sorted(rom_by_hash.values(), key=lambda r: -len(r["disks"]))
    boot_rows = sorted(boot_by_hash.values(), key=lambda r: -len(r["disks"]))

    bootable_disks = sum(len(r["disks"]) for r in rom_rows)
    md = [
        "# DOS catalog",
        "",
        "Empirical survey of the SAM Coupé corpus at `~/sam-corpus/disks/`,",
        "fingerprinting each disk's DOS by following the actual ROM ↔ DOS",
        "load contract (rom-disasm:20473-20598, the BOOTEX routine).",
        "",
        "## ROM bootability gate",
        "",
        "Before recording any fingerprint, the survey applies ROM BOOTEX's",
        "own bootability check (BTCK loop at rom-disasm:20586-20596):",
        "T4S1 bytes 256..259 are XORed with the literal `\"BOOT\"` and",
        "AND-masked with 0x5F. If any of the four bytes mismatches, ROM",
        "prints `NO DOS` and refuses to boot — the disk's T4S1 and slot 0",
        "are irrelevant for DOS classification, so the survey skips them.",
        "",
        "## ROM-contract fingerprint (bootable disks only)",
        "",
        "ROM's BOOT path reads **exactly track-4-sector-1 (512 bytes)** into",
        "`0x8000..0x81FF` (rom-disasm:20550, `LD HL,0x8000` + `INI` loop at",
        "RSA5), validates the 4-byte signature, then jumps to `0x8009`.",
        "**Nothing else** is part of the ROM contract — the bootstrap at",
        "0x8009 is per-DOS code that chooses where to find the rest of",
        "itself.",
        "",
        f"Each disk's strict ROM-contract fingerprint is therefore:",
        f"",
        f"- **Length:** {ROM_T4S1_LENGTH} bytes (constant)",
        f"- **Load address:** 0x{ROM_T4S1_LOAD_ADDRESS:04x} (constant)",
        f"- **Execution address:** 0x{ROM_T4S1_EXEC_ADDRESS:04x} (constant)",
        "- **Content SHA-256:** varies per disk (see table)",
        "",
        f"- Total disks scanned: **{total}**",
        f"- Disks ROM would refuse (`NO DOS`): **{len(not_bootable)}**",
        f"- Disks ROM would boot: **{bootable_disks}**",
        f"- Unique T4S1 contents (among bootable): **{len(rom_rows)}**",
        f"- Disks with malformed images: **{len(bad)}**",
        "",
        "### Unique T4S1 contents (most-used first)",
        "",
        "| SHA-256 (16) | Disks | Length | Load | Exec |",
        "|---|---:|---:|---:|---:|",
    ]
    for r in rom_rows:
        md.append(
            f"| `{r['hash']}` | {len(r['disks'])} | {r['length']} "
            f"| 0x{r['load_address']:04x} | 0x{r['exec_address']:04x} |"
        )
    md.append("")
    md.append("## Bootstrap-convention fingerprint (slot 0's file body)")
    md.append("")
    md.append("Every known SAM DOS's bootstrap chooses to fetch the rest of")
    md.append("its code from slot 0's file (it doesn't have to — the ROM")
    md.append("contract is silent on this — but they all do). Decoding")
    md.append("slot 0's 9-byte file-header gives the precise byte count,")
    md.append("load address, and execution address declared by the file.")
    md.append("This is the *bootstrap-convention* fingerprint, not part of")
    md.append("the ROM contract; the actual bootstrap code at 0x8009 might")
    md.append("ignore these values and load anything from anywhere.")
    md.append("")
    md.append(f"- Bootable disks with a usable slot-0 file: **{sum(len(r['disks']) for r in boot_rows)}**")
    md.append(f"- Bootable disks with no usable slot-0 file: **{len(no_slot0)}**")
    md.append(f"- Unique slot-0 bodies: **{len(boot_rows)}**")
    md.append("")
    md.append("### Unique slot-0 bodies (most-used first)")
    md.append("")
    md.append("| SHA-256 (16) | Disks | Length | Load | Exec |")
    md.append("|---|---:|---:|---:|---:|")
    for r in boot_rows:
        exec_disp = f"0x{r['exec_address']:04x}" if r["exec_address"] is not None else "(none)"
        md.append(
            f"| `{r['hash']}` | {len(r['disks'])} | {r['length']} "
            f"| 0x{r['load_address']:04x} | {exec_disp} |"
        )
    md.append("")
    md.append("## Sample disks per slot-0 body")
    md.append("")
    md.append("Extract any of these for disassembly:")
    md.append("```bash")
    md.append("python3 ~/git/samfile/tools/audit/extract_dos.py <hash-prefix>")
    md.append("```")
    md.append("")
    for r in boot_rows:
        exec_disp = f"0x{r['exec_address']:04x}" if r["exec_address"] is not None else "(none)"
        md.append(
            f"### `{r['hash']}` ({len(r['disks'])} disks, "
            f"length={r['length']} bytes, load=0x{r['load_address']:04x}, exec={exec_disp})"
        )
        md.append("")
        for d in r["disks"][:5]:
            md.append(f"- {d}")
        md.append("")

    if not_bootable:
        md.append(f"## Disks ROM would refuse with `NO DOS` ({len(not_bootable)})")
        md.append("")
        md.append("T4S1[256..259] doesn't match `BOOT` (masked 0x5F).")
        md.append("ROM BOOTEX prints `NO DOS` and refuses to load.")
        md.append("These disks have no bootable DOS; skipped from")
        md.append("classification.")
        md.append("")
        md.append("Sample (first 20):")
        for d in not_bootable[:20]:
            md.append(f"- {d}")
        md.append("")

    if no_slot0:
        md.append(f"## Bootable but no slot-0 file ({len(no_slot0)})")
        md.append("")
        md.append("Disks where ROM would boot (BOOT signature present)")
        md.append("but slot 0 is erased / its chain malformed, so the")
        md.append("slot-0 bootstrap-convention fingerprint doesn't apply.")
        md.append("The T4S1 fingerprint above still classifies them.")
        md.append("")
        for d in no_slot0[:20]:
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
    print(f"scanned: {total} disks; "
          f"ROM-bootable: {bootable_disks}; "
          f"ROM-refused (NO DOS): {len(not_bootable)}; "
          f"bad: {len(bad)}")
    print(f"unique T4S1: {len(rom_rows)}; "
          f"unique slot-0 bodies: {len(boot_rows)}; "
          f"bootable-but-no-slot0: {len(no_slot0)}")
    print()
    print("top 15 T4S1 fingerprints by disk count:")
    for r in rom_rows[:15]:
        print(f"  {len(r['disks']):4d}  {r['hash']}")
    print()
    print("top 15 slot-0 body fingerprints by disk count:")
    for r in boot_rows[:15]:
        exec_disp = f"0x{r['exec_address']:04x}" if r["exec_address"] is not None else "(none)"
        print(f"  {len(r['disks']):4d}  {r['hash']}  length={r['length']:5d} load=0x{r['load_address']:04x} exec={exec_disp}")


if __name__ == "__main__":
    main()
