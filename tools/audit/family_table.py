#!/usr/bin/env python3
"""Collapse the corpus's slot-0 DOS variants into DOS families.

Two slot-0 bodies are in the same family if they have the same
length and their byte-wise diff is below `--threshold` percent.
This handles three real cases:

1. **Per-magazine launcher data.** Same DOS, magazine-specific
   embedded auto-launch program in the slot-0 data section. Pure
   data, not code. Two variants from this case typically differ by
   well under 1%.
2. **Memory-config rebase.** Same DOS code reassembled for a
   different SAM RAM page (e.g. page 14 for 256 KB, page 30 for
   512 KB). The instructions are identical; only the page-selector
   constants and a few sysvar pointers change. Two such variants
   typically differ by ~0.6%.
3. **Build / patch differences.** Same DOS, minor patches —
   bugfixes, branding strings, build stamps. Sub-percent typically.

A bucket of variants where these three causes combine still stays
well under a few percent. A pair of bodies whose diff exceeds the
threshold is, by definition, NOT in the same family at this
threshold.

Cross-length variants (e.g. MasterDOS 15700 vs 15750) are kept as
separate families even at the same load address, because a length
change implies an inserted region whose semantics need verifying
before assuming compatibility.

Output: writes docs/dos-families.md and a brief summary to stdout.

Usage:
    family_table.py [--threshold 1.5]
"""

from __future__ import annotations

import argparse
import hashlib
import sys
from collections import Counter, defaultdict
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))
from survey_dos import (
    DISK_SIZE,
    decode_slot0_loadexec,
    is_bootable,
)

CORPUS = Path.home() / "sam-corpus" / "disks"
OUT = Path.home() / "git" / "samfile" / "docs" / "dos-families.md"


def byte_diff_pct(b1: bytes, b2: bytes) -> float | None:
    """Return percentage of differing bytes, or None for different lengths."""
    if len(b1) != len(b2):
        return None
    d = sum(1 for a, b in zip(b1, b2) if a != b)
    return 100 * d / len(b1)


def collect_variants() -> dict:
    """sha16 -> {body, length, load, exec, disks: [...]}"""
    variants: dict[str, dict] = {}
    for p in sorted(CORPUS.glob("*.mgt")):
        data = p.read_bytes()
        if len(data) != DISK_SIZE or not is_bootable(data):
            continue
        info = decode_slot0_loadexec(data)
        if info is None:
            continue
        load_addr, exec_addr, length, body = info
        h = hashlib.sha256(body).hexdigest()[:16]
        if h not in variants:
            variants[h] = {
                "body": body,
                "length": length,
                "load": load_addr,
                "exec": exec_addr,
                "disks": [],
            }
        variants[h]["disks"].append(p.stem)
    return variants


def cluster_into_families(variants: dict, threshold: float) -> dict:
    """Union-find over variants. Two variants are in the same family
    iff their byte_diff_pct is below threshold. Returns a dict
    {family_head_sha: [sha, ...]}."""
    parent = {h: h for h in variants}

    def find(x: str) -> str:
        while parent[x] != x:
            parent[x] = parent[parent[x]]
            x = parent[x]
        return x

    def union(a: str, b: str) -> None:
        ra, rb = find(a), find(b)
        if ra == rb:
            return
        # Prefer larger disk count as the family head.
        ns = {ra: 0, rb: 0}
        for h in variants:
            r = find(h)
            if r in ns:
                ns[r] += len(variants[h]["disks"])
        if ns[ra] >= ns[rb]:
            parent[rb] = ra
        else:
            parent[ra] = rb

    # Anchor = top SHA in each (length, load) bucket.
    bucket_top: dict[tuple[int, int], str] = {}
    for h, v in variants.items():
        key = (v["length"], v["load"])
        if (key not in bucket_top
                or len(v["disks"]) > len(variants[bucket_top[key]]["disks"])):
            bucket_top[key] = h
    anchors = sorted(bucket_top.values(), key=lambda h: -len(variants[h]["disks"]))

    # Pass 1: every variant joins its closest anchor under threshold.
    for h, v in variants.items():
        best, best_pct = None, None
        for a in anchors:
            pct = byte_diff_pct(v["body"], variants[a]["body"])
            if pct is None or pct >= threshold:
                continue
            if best_pct is None or pct < best_pct:
                best, best_pct = a, pct
        if best is not None:
            union(h, best)

    # Pass 2: merge anchors that are within threshold of each other.
    for i, a in enumerate(anchors):
        for b in anchors[i + 1:]:
            pct = byte_diff_pct(variants[a]["body"], variants[b]["body"])
            if pct is not None and pct < threshold:
                union(a, b)

    families: dict[str, list[str]] = defaultdict(list)
    for h in variants:
        families[find(h)].append(h)
    return families


def family_summary(variants: dict, families: dict) -> list[dict]:
    rows = []
    for head, vs in families.items():
        disks = sum(len(variants[h]["disks"]) for h in vs)
        lengths = sorted({variants[h]["length"] for h in vs})
        loads = sorted({variants[h]["load"] for h in vs})
        # Per-length max within-variance, since cross-length comparison
        # isn't meaningful byte-by-byte.
        bodies_by_len: dict[int, list[bytes]] = defaultdict(list)
        for h in vs:
            bodies_by_len[variants[h]["length"]].append(variants[h]["body"])
        per_len: list[tuple[int, float]] = []
        for L, bodies in bodies_by_len.items():
            if len(bodies) < 2:
                per_len.append((L, 0.0))
                continue
            diff_count = 0
            for i in range(L):
                ref = bodies[0][i]
                if any(b[i] != ref for b in bodies[1:]):
                    diff_count += 1
            per_len.append((L, 100 * diff_count / L))
        max_var = max((v for _, v in per_len), default=0.0)
        rows.append({
            "head": head,
            "variants": len(vs),
            "disks": disks,
            "lengths": lengths,
            "loads": loads,
            "max_within_var_pct": max_var,
            "sample_disk": variants[head]["disks"][0] if variants[head]["disks"] else "",
        })
    rows.sort(key=lambda r: -r["disks"])
    return rows


def page_label(load: int) -> str:
    page = load >> 14
    return f"p{page + 1}"


def write_markdown(rows: list[dict], threshold: float, total_variants: int) -> None:
    n_disks = sum(r["disks"] for r in rows)
    n_big = sum(1 for r in rows if r["disks"] >= 5)
    n_one = sum(1 for r in rows if r["disks"] == 1)
    md = [
        "# DOS families",
        "",
        f"Slot-0 DOS body variants clustered into families. Two variants",
        f"are in the same family iff they have the same body length and",
        f"their byte-wise diff is below **{threshold:.1f}%** of the body",
        f"length.",
        "",
        "## Rationale",
        "",
        "Three real causes of small-percentage variation between same-",
        "DOS slot-0 bodies are all captured under one threshold:",
        "",
        "1. **Per-magazine / per-disk launcher data.** Magazine-specific",
        "   embedded auto-launch programs in the slot-0 data section.",
        "   Pure data, not code. Typical diff: well under 1%.",
        "2. **Memory-config rebase.** Same DOS code reassembled for a",
        "   different SAM RAM page (e.g. page 14 vs page 30). Identical",
        "   instructions, different page-selector constants and sysvar",
        "   pointers. Typical diff: ~0.6%.",
        "3. **Build / patch differences.** Bugfixes, branding, build",
        "   stamps. Sub-percent typically.",
        "",
        "A pair of bodies whose diff exceeds the threshold is *not* in",
        "the same family at this threshold. Cross-length variants (e.g.",
        "MasterDOS 15700 vs 15750) are kept as separate families — a",
        "length change implies an inserted region whose semantics need",
        "verifying before assuming compatibility.",
        "",
        "## Summary",
        "",
        f"- **Threshold:** {threshold:.1f}% byte-diff",
        f"- **Total bootable disks:** {n_disks}",
        f"- **Total unique slot-0 SHAs:** {total_variants}",
        f"- **Total families:** {len(rows)}",
        f"- **Families with ≥ 5 disks:** {n_big}",
        f"- **Families with 1 disk only (long-tail customs):** {n_one}",
        "",
        "## Table",
        "",
        "| Rank | Family head SHA | Variants | Disks | Lengths | Load addresses (page) | Max within-variance | Sample disk |",
        "|---:|---|---:|---:|---|---|---:|---|",
    ]
    for i, r in enumerate(rows, 1):
        lens = ",".join(str(L) for L in r["lengths"])
        loads = ", ".join(f"0x{L:06x} ({page_label(L)})" for L in r["loads"])
        md.append(
            f"| {i} | `{r['head']}` | {r['variants']} | {r['disks']} "
            f"| {lens} | {loads} | {r['max_within_var_pct']:.2f}% | {r['sample_disk'][:50]} |"
        )
    OUT.write_text("\n".join(md) + "\n")


def main() -> None:
    ap = argparse.ArgumentParser()
    ap.add_argument("--threshold", type=float, default=1.5,
                    help="byte-diff %% threshold for family membership (default 1.5)")
    args = ap.parse_args()

    print(f"collecting slot-0 variants from {CORPUS}...")
    variants = collect_variants()
    print(f"{len(variants)} unique slot-0 SHAs across {sum(len(v['disks']) for v in variants.values())} bootable disks")
    print(f"clustering at threshold {args.threshold:.1f}%...")
    families = cluster_into_families(variants, args.threshold)
    rows = family_summary(variants, families)
    write_markdown(rows, args.threshold, len(variants))
    print(f"wrote {OUT}")
    print()
    print(f"{'Rank':>4s}  {'Head SHA':17s}  {'Vars':>4s}  {'Disks':>5s}  {'Length(s)':17s}  {'MaxVar%':>7s}  Sample")
    print("=" * 110)
    for i, r in enumerate(rows, 1):
        lens = ",".join(str(L) for L in r["lengths"])
        print(
            f"  {i:>3d}  {r['head']:17s}  {r['variants']:>4d}  {r['disks']:>5d}  "
            f"{lens:17s}  {r['max_within_var_pct']:>6.2f}%  {r['sample_disk'][:50]}"
        )


if __name__ == "__main__":
    main()
