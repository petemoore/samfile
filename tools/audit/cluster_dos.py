#!/usr/bin/env python3
"""Within each (length, load_address) bucket of slot-0 DOS bodies,
identify which byte positions are stable across all variants (= code)
vs which vary (= data / config).

The hypothesis: most "unique" DOS SHAs in the corpus are the same
DOS with different baked-in config (cursor position, default drive,
menu strings, magazine table-of-contents, etc.). If two variants
share their entire code section and differ only in a small,
contiguous data section, they're the same DOS family.

Output (one section per (length, load) bucket with >= 2 variants):

  - Number of variants, number of disks
  - "Stable byte count" — positions where all variants agree
  - "Variable byte count" — positions where any pair disagrees
  - The variable byte ranges (start..end, length)
  - A hex preview of each variable run for each variant
"""

from __future__ import annotations

import hashlib
import sys
from collections import defaultdict
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))
from survey_dos import (
    DISK_SIZE,
    decode_slot0_loadexec,
    is_bootable,
)

CORPUS = Path.home() / "sam-corpus" / "disks"
OUT = Path.home() / "git" / "samfile" / "docs" / "dos-clusters.md"


def collect_bucket(min_variants: int = 2, min_disks: int = 10) -> dict:
    """Walk the corpus, group bootable-disk slot-0 bodies by
    (length, load_address), and return a dict keyed by that tuple
    with {variants: {sha: (body, [disks])}}."""
    buckets: dict[tuple[int, int], dict] = defaultdict(lambda: {"variants": {}})
    for p in sorted(CORPUS.glob("*.mgt")):
        data = p.read_bytes()
        if len(data) != DISK_SIZE or not is_bootable(data):
            continue
        info = decode_slot0_loadexec(data)
        if info is None:
            continue
        load_addr, exec_addr, length, body = info
        key = (length, load_addr)
        h = hashlib.sha256(body).hexdigest()[:16]
        if h not in buckets[key]["variants"]:
            buckets[key]["variants"][h] = (body, [])
        buckets[key]["variants"][h][1].append(p.stem)

    # Filter to interesting buckets only.
    filtered = {}
    for key, b in buckets.items():
        n_variants = len(b["variants"])
        n_disks = sum(len(ds) for (_, ds) in b["variants"].values())
        if n_variants >= min_variants and n_disks >= min_disks:
            filtered[key] = b
    return filtered


def variation_profile(variants: list[bytes]) -> tuple[list[bool], list[tuple[int, int]]]:
    """For a list of equal-length variants, return (stable_mask,
    variable_runs).
    stable_mask[i] = True iff variants[*][i] is the same for all variants.
    variable_runs = list of (start, end_inclusive) ranges where stability is False.
    """
    n = len(variants[0])
    stable: list[bool] = [True] * n
    for i in range(n):
        b0 = variants[0][i]
        for v in variants[1:]:
            if v[i] != b0:
                stable[i] = False
                break
    runs: list[tuple[int, int]] = []
    in_run = False
    start = 0
    for i in range(n):
        if not stable[i]:
            if not in_run:
                start = i
                in_run = True
        else:
            if in_run:
                runs.append((start, i - 1))
                in_run = False
    if in_run:
        runs.append((start, n - 1))
    return stable, runs


def main() -> None:
    buckets = collect_bucket(min_variants=2, min_disks=10)
    if not buckets:
        print("no qualifying buckets")
        return

    md = [
        "# DOS clusters: stable code vs variable data",
        "",
        "For each `(length, load_address)` bucket with multiple slot-0 body",
        "variants, this report identifies which byte positions are stable",
        "across every variant (likely code) and which vary (likely a data",
        "section — magazine table-of-contents, baked-in config, build",
        "stamps).",
        "",
        "A bucket where variation is small and concentrated in a few short",
        "runs almost certainly represents **one DOS family** with per-disk",
        "data patches; the unique SHAs collapse into a single family-level",
        "fingerprint for rule-scoping purposes. A bucket where variation",
        "is large and spread across the body means the variants are",
        "genuinely different DOSes that happen to share the same size.",
        "",
    ]

    summary_rows = []
    for (length, load), b in sorted(buckets.items(), key=lambda kv: -sum(len(ds) for (_, ds) in kv[1]["variants"].values())):
        variants = list(b["variants"].items())  # [(sha, (body, [disks]))]
        bodies = [v[1][0] for v in variants]
        stable, runs = variation_profile(bodies)
        stable_count = sum(1 for s in stable if s)
        variable_count = length - stable_count
        n_variants = len(variants)
        n_disks = sum(len(ds) for (_, ds) in b["variants"].values())
        summary_rows.append({
            "length": length,
            "load": load,
            "variants": n_variants,
            "disks": n_disks,
            "stable": stable_count,
            "variable": variable_count,
            "runs": runs,
            "data": variants,
        })

    md.append("## Buckets")
    md.append("")
    md.append("| Length | Load | Variants | Disks | Stable bytes | Variable bytes | Variation % | Distinct runs |")
    md.append("|---:|---:|---:|---:|---:|---:|---:|---:|")
    for r in summary_rows:
        pct = 100 * r["variable"] / r["length"]
        md.append(
            f"| {r['length']} | 0x{r['load']:05x} | {r['variants']} | {r['disks']} "
            f"| {r['stable']} | {r['variable']} | {pct:.2f}% | {len(r['runs'])} |"
        )
    md.append("")

    for r in summary_rows:
        length, load = r["length"], r["load"]
        md.append(f"## length={length}, load=0x{load:05x} — {r['variants']} variants across {r['disks']} disks")
        md.append("")
        pct = 100 * r["variable"] / length
        md.append(f"- Stable bytes:   **{r['stable']:5d}** ({100*r['stable']/length:.2f}%)")
        md.append(f"- Variable bytes: **{r['variable']:5d}** ({pct:.2f}%)")
        md.append(f"- Distinct variable runs: **{len(r['runs'])}**")
        md.append("")
        if not r["runs"]:
            md.append("All variants identical — this bucket is a single DOS, the multiple")
            md.append("SHA buckets must have been a survey artefact.")
            md.append("")
            continue

        # Show the variable runs (compact + with hex preview).
        md.append("### Variable runs")
        md.append("")
        md.append("| Start | End | Length |")
        md.append("|---:|---:|---:|")
        for (s, e) in r["runs"]:
            md.append(f"| 0x{s:04x} | 0x{e:04x} | {e-s+1} |")
        md.append("")

        # Show the first ~5 variants' hex content for the variable runs.
        md.append("### Hex preview of variable runs (top 5 variants by disk count)")
        md.append("")
        ordered = sorted(r["data"], key=lambda v: -len(v[1][1]))[:5]
        for (sha, (body, disks)) in ordered:
            md.append(f"**`{sha}` ({len(disks)} disks)** — example: `{disks[0][:60]}`")
            md.append("")
            md.append("```")
            for (s, e) in r["runs"]:
                slc = body[s:e+1]
                # Show as hex with up to 32 bytes per line, plus ASCII gloss.
                ascii_glyph = "".join(chr(b) if 32 <= b < 127 else "." for b in slc)
                md.append(f"  0x{s:04x}..0x{e:04x} ({e-s+1}b) {slc.hex(' ')}  |{ascii_glyph}|")
            md.append("```")
            md.append("")

    OUT.write_text("\n".join(md) + "\n")
    print(f"wrote {OUT}")
    print()
    for r in summary_rows:
        pct = 100 * r["variable"] / r["length"]
        print(
            f"  length={r['length']:5d} load=0x{r['load']:05x}  "
            f"variants={r['variants']:3d}  disks={r['disks']:4d}  "
            f"stable={r['stable']:5d}  variable={r['variable']:5d} ({pct:5.2f}%)  "
            f"runs={len(r['runs'])}"
        )


if __name__ == "__main__":
    main()
