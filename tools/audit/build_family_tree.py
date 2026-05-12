#!/usr/bin/env python3
"""Materialise a per-family directory tree under docs/dos-families/.

For every DOS family produced by family_table.py's clustering, write a
self-contained directory with:

  <head_sha16>[-label]/
    README.md                  — family metadata, variants, page layouts,
                                 source-binding notes, sample disks
    body.bin                   — exact slot-0 body of the family head
    body.hex                   — xxd-style hex dump for grepability
    variants/<sha>.bin         — body for every other variant (big families)
    variants/<sha>.md          — byte-diff summary against head
    src/                       — original assembly source (copied in for
                                 the SAMDOS-2 and MasterDOS families;
                                 stub for everything else)

Plus a top-level docs/dos-families/INDEX.md that cross-links into
docs/dos-families.md (the human-readable family table) and into each
per-family directory.

The point of this tree is to give future agents one place to grep,
diff and reason about a specific DOS without having to re-run the
ROM-contract extractor first. Once written, each directory is
self-describing: a fresh agent can `cat README.md`, `xxd body.bin`,
and (for samdos / masterdos) read the original commented assembly
right there.

Usage:
    build_family_tree.py [--threshold 1.5]

Idempotent: re-running overwrites body.bin / body.hex / READMEs but
leaves manually-added notes alone, as long as they live outside the
generated files (e.g. NOTES.md).
"""

from __future__ import annotations

import argparse
import hashlib
import shutil
import sys
from collections import defaultdict
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))
from family_table import (
    byte_diff_pct,
    cluster_into_families,
    collect_variants,
)
from survey_dos import (
    DISK_SIZE,
    decode_slot0_loadexec,
    is_bootable,
)

REPO_ROOT = Path.home() / "git" / "samfile"
OUT_ROOT = REPO_ROOT / "docs" / "dos-families"
SAMDOS_REPO = Path.home() / "git" / "samdos"
MASTERDOS_REPO = Path.home() / "git" / "masterdos"

# Hand-curated labels for known families. Keyed by family-head SHA16.
# Tail families with no known label use head_sha16 only.
#
# Rule: only label a family if the binding to a named DOS is verified by
# byte-for-byte SHA match against an upstream reference binary. Sample-
# disk names are not enough — they're often just disks happen to use that
# DOS, not name it. Speculative labels mislead future agents more than
# they help. Expand this list as more reference binaries are found.
FAMILY_LABELS: dict[str, str] = {
    "9bc0fb4b949109e8": "samdos2",
    "152b811ed65b651d": "masterdos-v2.3",
}

# Family head -> a "binary-of-record" SHA whose body matches an upstream
# source's assemble output. Recorded so the README can say "the source
# in src/ assembles to variants/<sha>.bin", not the family-head body.
SOURCE_BINDS: dict[str, dict] = {
    "9bc0fb4b949109e8": {
        "label": "samdos2",
        "src_repo": SAMDOS_REPO,
        "src_subdir": "src",
        "src_files_glob": "*.s",
        "bin_sha_full": "3cca541beb3f9fe93402a770997945b2be852e69f278d2b176ba0bbc4fbb6077",
        "bin_sha16": "3cca541beb3f9fe9",
        "bin_path": "res/samdos2.reference.bin",
        "upstream_zip": "https://ftp.nvg.ntnu.no/pub/sam-coupe/sources/SamDos2InCometFormatMasterv1.2.zip",
        "notes": [
            "Source from https://github.com/stefandrissen/samdos (Stefan",
            "Drissen). HEAD is samdos2, assembled byte-identical to",
            "variants/3cca541beb3f9fe9.bin. The git history of that repo",
            "carries the five upstream comp1..comp5 versions as separate",
            "commits.",
            "",
            "The upstream archive `SamDos2InCometFormatMasterv1.2.zip`",
            "contains a SAM .dsk image (also fetched into",
            "`upstream/SamDos2InCometFormatMasterv1.2.dsk` if available)",
            "with all 28 source files (a1.s..h2.s, comp1.s..comp5.s,",
            "gm1.s, gm2.s, ldit1.s..ldit3.s, ld1.s).",
        ],
    },
    "152b811ed65b651d": {
        "label": "masterdos-15750-v2.3",
        "src_repo": MASTERDOS_REPO,
        "src_subdir": "src",
        "src_files_glob": "*.asm",
        "bin_sha_full": "152b811ed65b651df25e29f49e15340bec84ef3deebfba4eaa6cd76bfbb31fae",
        "bin_sha16": "152b811ed65b651d",
        "bin_path": "res/MDOS23.bin",
        "notes": [
            "Source from https://github.com/dandoore/masterdos (Dan Doore).",
            "`src/masterdos23.asm` assembles to body.bin. The README",
            "notes v2.2 and v2.3 differ only at points labelled `Fix_*`",
            "in the source.",
        ],
    },
}


def safe_dirname(head: str) -> str:
    label = FAMILY_LABELS.get(head)
    return f"{head}-{label}" if label else head


def hex_dump(body: bytes, width: int = 16) -> str:
    """Return an xxd-style hex dump with ASCII gloss."""
    lines = []
    for off in range(0, len(body), width):
        chunk = body[off:off + width]
        hex_part = " ".join(f"{b:02x}" for b in chunk)
        ascii_part = "".join(chr(b) if 32 <= b < 127 else "." for b in chunk)
        lines.append(f"{off:06x}  {hex_part:<{width*3-1}}  |{ascii_part}|")
    return "\n".join(lines) + "\n"


def diff_summary(head_body: bytes, variant_body: bytes) -> tuple[float, list[tuple[int, int]]]:
    """(diff_pct, list of (start, end_inclusive) runs of differing bytes)."""
    if len(head_body) != len(variant_body):
        return float("nan"), []
    n = len(head_body)
    diff_positions = [i for i in range(n) if head_body[i] != variant_body[i]]
    pct = 100 * len(diff_positions) / n
    runs: list[tuple[int, int]] = []
    if not diff_positions:
        return pct, runs
    run_start = diff_positions[0]
    prev = diff_positions[0]
    for p in diff_positions[1:]:
        if p == prev + 1:
            prev = p
            continue
        runs.append((run_start, prev))
        run_start = p
        prev = p
    runs.append((run_start, prev))
    return pct, runs


def write_variant_diff(out_path: Path, head_sha: str, variant_sha: str,
                       head_body: bytes, variant_body: bytes,
                       disks: list[str]) -> None:
    pct, runs = diff_summary(head_body, variant_body)
    lines = [
        f"# Variant `{variant_sha}` vs family head `{head_sha[:16]}`",
        "",
        f"- **Body length:** {len(variant_body)} bytes",
        f"- **Disks with this body:** {len(disks)}",
        f"- **Byte-diff vs head:** {pct:.3f}%",
        f"- **Distinct differing runs:** {len(runs)}",
        "",
    ]
    if len(disks) <= 20:
        lines.append("## Disks")
        lines.append("")
        for d in sorted(disks):
            lines.append(f"- {d}")
        lines.append("")
    else:
        lines.append("## Sample disks (first 20)")
        lines.append("")
        for d in sorted(disks)[:20]:
            lines.append(f"- {d}")
        lines.append(f"")
        lines.append(f"... and {len(disks) - 20} more.")
        lines.append("")

    if runs:
        lines.append("## Differing byte runs")
        lines.append("")
        lines.append("| Start | End | Length | Head bytes | Variant bytes |")
        lines.append("|---:|---:|---:|---|---|")
        for (s, e) in runs[:50]:
            head_slc = head_body[s:e + 1]
            var_slc = variant_body[s:e + 1]
            head_hex = head_slc.hex(" ") if e - s + 1 <= 16 else (head_slc[:16].hex(" ") + " …")
            var_hex = var_slc.hex(" ") if e - s + 1 <= 16 else (var_slc[:16].hex(" ") + " …")
            lines.append(f"| 0x{s:04x} | 0x{e:04x} | {e - s + 1} | `{head_hex}` | `{var_hex}` |")
        if len(runs) > 50:
            lines.append(f"")
            lines.append(f"... and {len(runs) - 50} more runs.")
        lines.append("")
    out_path.write_text("\n".join(lines))


def page_label(load: int) -> str:
    page = (load >> 14)
    return f"p{page + 1}"


def copy_source_tree(bind: dict, dst: Path) -> bool:
    """Copy commented assembly source files for a family. Returns True
    if anything was copied. Idempotent: wipes dst and re-copies."""
    src_root = bind["src_repo"] / bind["src_subdir"]
    if not src_root.exists():
        return False
    if dst.exists():
        shutil.rmtree(dst)
    dst.mkdir(parents=True, exist_ok=True)
    glob = bind["src_files_glob"]
    n = 0
    for f in sorted(src_root.glob(glob)):
        if f.is_file():
            shutil.copy2(f, dst / f.name)
            n += 1
    # Always copy README if present (under src or the repo root).
    for cand in [src_root / "Readme.md", src_root / "README.md",
                 bind["src_repo"] / "README.md"]:
        if cand.exists():
            shutil.copy2(cand, dst / f"UPSTREAM-{cand.name}")
            break
    return n > 0


def write_family_readme(out_dir: Path, head: str, info: dict,
                        family_variants: list[str], variants: dict,
                        bind: dict | None) -> None:
    label = FAMILY_LABELS.get(head)
    title = f"# DOS family `{head}`" + (f" — {label}" if label else "")
    total_disks = sum(len(variants[v]["disks"]) for v in family_variants)
    lengths = sorted({variants[v]["length"] for v in family_variants})
    loads = sorted({variants[v]["load"] for v in family_variants})
    execs = sorted({variants[v]["exec"] for v in family_variants
                    if variants[v]["exec"] is not None})

    lines = [
        title,
        "",
        "Self-contained materialisation of one DOS family from the SAM",
        "Coupé corpus. The family is the equivalence class of slot-0 DOS",
        "bodies clustered at 1.5% byte-diff (see",
        "`docs/dos-families.md` for the full table).",
        "",
        "## Identity",
        "",
        f"- **Family-head SHA16:** `{head}`",
        f"- **Variants in family:** {len(family_variants)}",
        f"- **Disks in family:** {total_disks}",
        f"- **Body length(s):** {', '.join(str(L) for L in lengths)}",
        f"- **Load address(es):** {', '.join(f'0x{L:06x} ({page_label(L)})' for L in loads)}",
        f"- **Execution address(es):** {', '.join(f'0x{e:06x}' for e in execs) if execs else '(unset / 0xFF)'}",
        "",
        "## Files in this directory",
        "",
        f"- `body.bin` — exact slot-0 body of the family-head SHA",
        f"  (`{head}`). Header-decoded, so byte 0 is the first byte the",
        f"  ROM would copy to the body's load address.",
        f"- `body.hex` — xxd-style hex dump of `body.bin` (big families only).",
        f"- `variants/*.md` — byte-diff summary of each variant against the",
        f"  head, including the differing byte ranges in hex so the variant",
        f"  body can be reconstructed from head + diff. Only written for",
        f"  families with at least 5 disks; extract any variant body with",
        f"  `tools/audit/extract_dos.py <sha>`.",
        f"- `src/` — commented original assembly source (only present when",
        f"  upstream source is known to assemble to a binary in this family).",
        "",
    ]

    if bind:
        lines.extend([
            "## Source binding",
            "",
            f"- **Binary-of-record SHA16:** `{bind['bin_sha16']}` "
            f"(full SHA `{bind['bin_sha_full']}`)",
            f"- **Upstream source:** `{bind['src_repo']}` "
            f"({bind['src_subdir']}/{bind['src_files_glob']})",
            f"- **Reference binary in upstream:** "
            f"`{bind['src_repo']}/{bind['bin_path']}`",
            "",
        ])
        if "upstream_zip" in bind:
            lines.append(f"- **Upstream archive:** {bind['upstream_zip']}")
            lines.append("")
        if "notes" in bind:
            lines.append("### Notes from upstream README")
            lines.append("")
            for n in bind["notes"]:
                lines.append(f"> {n}" if n else ">")
            lines.append("")
    else:
        primary_load = loads[0]
        lines.extend([
            "## Source binding",
            "",
            "No upstream source identified for this family yet. The body",
            "must be disassembled directly from `body.bin` if you need to",
            "reason about its semantics. Disassemble with any z80",
            "disassembler, e.g.",
            "",
            "```",
            f"z80dasm body.bin -a -t -o 0x{primary_load:06x} > body.z80.s",
            "```",
            "",
            "If the family has multiple load addresses, pick the one that",
            "covers the binary you care about. See `references/README.md`",
            "for the SAM ROM v3.0 annotated disassembly, which is the",
            "single biggest aid when reading these bodies.",
            "",
        ])

    lines.append("## Variants in this family")
    lines.append("")
    lines.append("| Variant SHA16 | Disks | Length | Load | Exec | Within-fam diff vs head |")
    lines.append("|---|---:|---:|---|---|---:|")
    head_body = variants[head]["body"]
    sorted_vs = sorted(family_variants,
                       key=lambda v: -len(variants[v]["disks"]))
    for v in sorted_vs:
        info_v = variants[v]
        ex_str = f"0x{info_v['exec']:06x}" if info_v.get("exec") is not None else "—"
        if v == head:
            pct_str = "head"
        else:
            pct, _ = diff_summary(head_body, info_v["body"])
            pct_str = "n/a (length differs)" if pct != pct else f"{pct:.3f}%"
        marker = " ←source-of-record" if bind and v.startswith(bind["bin_sha16"]) else ""
        lines.append(
            f"| `{v}`{marker} | {len(info_v['disks'])} | {info_v['length']} | "
            f"0x{info_v['load']:06x} ({page_label(info_v['load'])}) | {ex_str} | {pct_str} |"
        )
    lines.append("")

    # Sample disks (first ~10).
    sample_disks = []
    for v in sorted_vs:
        sample_disks.extend(variants[v]["disks"])
    sample_disks.sort()
    lines.append("## Sample disks")
    lines.append("")
    for d in sample_disks[:10]:
        lines.append(f"- {d}")
    if len(sample_disks) > 10:
        lines.append(f"")
        lines.append(f"... and {len(sample_disks) - 10} more.")
    lines.append("")

    (out_dir / "README.md").write_text("\n".join(lines))


def main() -> None:
    ap = argparse.ArgumentParser()
    ap.add_argument("--threshold", type=float, default=1.5,
                    help="family-membership %% threshold (default 1.5)")
    ap.add_argument("--big-family-min-disks", type=int, default=5,
                    help="big families (>= this many disks) materialise all "
                         "variants; smaller families only emit the head")
    args = ap.parse_args()

    print(f"collecting slot-0 variants ...")
    variants = collect_variants()
    print(f"clustering at threshold {args.threshold}% ...")
    families = cluster_into_families(variants, args.threshold)
    print(f"{len(families)} families across {len(variants)} variants")

    OUT_ROOT.mkdir(parents=True, exist_ok=True)
    # Wipe stale family dirs that the generator owns; leave the top-level
    # INDEX.md and any human-authored files at OUT_ROOT level untouched.
    for child in OUT_ROOT.iterdir():
        if child.is_dir():
            shutil.rmtree(child)

    index_rows = []
    big_threshold = args.big_family_min_disks
    for head, members in families.items():
        out_dir = OUT_ROOT / safe_dirname(head)
        out_dir.mkdir(parents=True, exist_ok=True)
        head_info = variants[head]
        head_body = head_info["body"]
        # Head body + hex dump (hex only for big families to keep the
        # tree compact; regen any tail family's hex with
        # `xxd body.bin > body.hex` if needed).
        (out_dir / "body.bin").write_bytes(head_body)
        total_disks_pre = sum(len(variants[v]["disks"]) for v in members)
        if total_disks_pre >= big_threshold:
            (out_dir / "body.hex").write_text(hex_dump(head_body))

        bind = SOURCE_BINDS.get(head)
        if bind:
            copied = copy_source_tree(bind, out_dir / "src")
            if not copied:
                # Source repo missing — record a stub.
                stub = out_dir / "src" / "MISSING.md"
                stub.parent.mkdir(parents=True, exist_ok=True)
                stub.write_text(
                    f"# Upstream source not available locally\n\n"
                    f"Expected at `{bind['src_repo']}` but the path does"
                    f" not exist. Clone it manually if you need to grep"
                    f" the assembly source.\n"
                )

        total_disks = sum(len(variants[v]["disks"]) for v in members)
        if total_disks >= big_threshold:
            v_dir = out_dir / "variants"
            v_dir.mkdir(parents=True, exist_ok=True)
            for v in members:
                if v == head:
                    continue
                vinfo = variants[v]
                # Variant .md captures the diff to head — including the
                # differing byte ranges in hex, so the variant body can
                # be reconstructed from head + diff. Skip writing the
                # variant body to keep the tree small; extract any
                # variant body with `tools/audit/extract_dos.py <sha>`.
                write_variant_diff(v_dir / f"{v}.md", head, v,
                                   head_body, vinfo["body"], vinfo["disks"])

        write_family_readme(out_dir, head, head_info, members, variants, bind)

        index_rows.append({
            "head": head,
            "label": FAMILY_LABELS.get(head, ""),
            "variants": len(members),
            "disks": total_disks,
            "lengths": sorted({variants[v]["length"] for v in members}),
            "loads": sorted({variants[v]["load"] for v in members}),
            "has_source": bool(bind),
            "dir": safe_dirname(head),
        })

    index_rows.sort(key=lambda r: -r["disks"])

    # Write INDEX.md.
    n_disks = sum(r["disks"] for r in index_rows)
    n_with_source = sum(1 for r in index_rows if r["has_source"])
    idx = [
        "# DOS families — per-family directory tree",
        "",
        "One subdirectory per family in [`docs/dos-families.md`].",
        "Each contains the family-head body, a hex dump, READMEs",
        "describing the variants and sample disks, and (when",
        "available) the commented original assembly source.",
        "",
        "## Summary",
        "",
        f"- **Families:** {len(index_rows)}",
        f"- **Total disks across families:** {n_disks}",
        f"- **Families with upstream source attached:** {n_with_source}",
        "",
        "Big families (≥ 5 disks) materialise all member-variant",
        "bodies under `variants/`. Smaller families only emit the",
        "family-head body to keep the tree compact; if you need a",
        "non-head variant body, run `tools/audit/extract_dos.py`",
        "against any disk that contains it.",
        "",
        "## Index",
        "",
        "| Rank | Family | Variants | Disks | Length(s) | Source |",
        "|---:|---|---:|---:|---|:-:|",
    ]
    for i, r in enumerate(index_rows, 1):
        label = f" ({r['label']})" if r["label"] else ""
        lens = ", ".join(str(L) for L in r["lengths"])
        src = "yes" if r["has_source"] else "—"
        idx.append(
            f"| {i} | [`{r['head']}`]({r['dir']}/){label} | "
            f"{r['variants']} | {r['disks']} | {lens} | {src} |"
        )
    idx.append("")
    idx.append("Regenerate this tree with `tools/audit/build_family_tree.py`.")
    idx.append("")
    (OUT_ROOT / "INDEX.md").write_text("\n".join(idx))

    # Top-level references/ — pointers to canonical external materials
    # rather than copies. Keeps the samfile repo lean while giving
    # future agents one place to learn where to look.
    refs_dir = OUT_ROOT / "references"
    refs_dir.mkdir(parents=True, exist_ok=True)
    rom_disasm = Path.home() / "git" / "sam-aarch64" / "docs" / "sam" / "sam-coupe_rom-v3.0_annotated-disassembly.txt"
    samdos_repo = SAMDOS_REPO
    masterdos_repo = MASTERDOS_REPO
    (refs_dir / "README.md").write_text(
        "# External reference materials\n"
        "\n"
        "Per-family directories embed the commented assembly source for\n"
        "the DOSes whose source we have. The materials below live\n"
        "outside the samfile repo and are too large or too remote to\n"
        "ship inline — but every agent reasoning about a family should\n"
        "know they exist.\n"
        "\n"
        "## ROM v3.0 annotated disassembly\n"
        "\n"
        f"- **Path:** `{rom_disasm}`\n"
        "- **Size:** ~1.1 MB, 27353 lines\n"
        "- **Use for:** the LOAD-path semantics that every DOS plugs\n"
        "  into. Grep for `BOOTEX`, `BTCK`, `LDHD`, `GTFLE`, `HCONR`\n"
        "  to walk the SAVE / LOAD / boot-sector chain.\n"
        "- **Source:** captured in `~/git/sam-aarch64/docs/sam/` and in\n"
        "  `~/git/migrate-build-disk-to-go/docs/sam/` (same file).\n"
        "\n"
        "## SAMDOS source — upstream tokenised archive\n"
        "\n"
        "- **Upstream:** https://ftp.nvg.ntnu.no/pub/sam-coupe/sources/SamDos2InCometFormatMasterv1.2.zip\n"
        "- **Contains:** a SAM .dsk image with 28 source files\n"
        "  (`a1..h2.S`, `comp1..comp5.S`, `gm1`, `gm2`, `ldit1..3`,\n"
        "  `ld1`). These are in COMET tokenised format — not plain\n"
        "  text. Use `samfile extract -i <dsk>` to pull them out, then\n"
        "  a Comet detokeniser (or compare against Stefan Drissen's\n"
        "  git history at https://github.com/stefandrissen/samdos which\n"
        "  has the lowered / detokenised plain-text form).\n"
        f"- **Local clean copy:** `{samdos_repo}` (Drissen's repo).\n"
        "  HEAD assembles to `9bc0fb4b949109e8-samdos2/variants/3cca541beb3f9fe9.bin`.\n"
        "  Earlier comp1..comp5 versions are reachable via `git log`.\n"
        "\n"
        "## MasterDOS source\n"
        "\n"
        f"- **Upstream:** https://github.com/dandoore/masterdos (Dan Doore)\n"
        f"- **Local clean copy:** `{masterdos_repo}`\n"
        "- `src/masterdos23.asm` assembles to\n"
        "  `152b811ed65b651d-masterdos-v2.3/body.bin`. v2.2 and v2.3\n"
        "  differ only at points labelled `Fix_*` in the source.\n"
        "\n"
        "## How to disassemble a DOS body when no source is available\n"
        "\n"
        "Most non-SAMDOS / non-MasterDOS families have no upstream\n"
        "source. Disassemble the body directly:\n"
        "\n"
        "```\n"
        "# z80dasm: the load address is recorded in each family's README\n"
        "z80dasm -a -t -o0x008009 body.bin > body.z80.s\n"
        "```\n"
        "\n"
        "or use any equivalent z80 disassembler. The annotated ROM\n"
        "disassembly above is invaluable for cross-referencing CALL\n"
        "targets and hardware port writes.\n"
    )

    print(f"wrote {refs_dir}/README.md")
    print(f"wrote {OUT_ROOT}/INDEX.md ({len(index_rows)} families)")
    print(f"  families with source attached: {n_with_source}")
    print(f"  big families (>= {big_threshold} disks): "
          f"{sum(1 for r in index_rows if r['disks'] >= big_threshold)}")


if __name__ == "__main__":
    main()
