#!/usr/bin/env python3
"""Run the *current branch's* samfile verify against only the disks
whose slot-0 body SHA matches a reference binary exactly. Aggregate
the JSONL findings into a per-rule report listing every fire.

This is the canonical "scientific test" for whether a rule fires on
the canonical SAMDOS-2 cohort itself. The strict-SHA filter on the
existing `findings.db` (family_report.py --strict-sha) only re-slices
data captured at an earlier rule version. This script:

  1. Picks all corpus disks whose slot-0 body SHA == reference SHA.
  2. Rebuilds samfile from the current working tree.
  3. Runs `samfile verify --format jsonl` against each.
  4. Aggregates fails per rule, lists every disk/ref/severity/message,
     and emits a markdown report.

Usage:
    strict_sha_scan.py [--ref-bin PATH] [--sha SHA256] [--out PATH]

  --ref-bin   Path to a binary; its sha256 is the match target.
              Defaults to ~/git/samdos/res/samdos2.reference.bin
              (samdos2, source-of-record for the SAMDOS-2 family).
  --sha       Override: an exact sha256 hex to match. If given, --ref-bin
              is ignored.
  --out       Output markdown path. Defaults to
              docs/strict-sha-scan-<sha16>.md
"""

from __future__ import annotations

import argparse
import hashlib
import json
import shutil
import subprocess
import sys
import tempfile
from collections import defaultdict
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))
from survey_dos import (
    DISK_SIZE,
    decode_slot0_loadexec,
    is_bootable,
)

REPO_ROOT = Path.home() / "git" / "samfile"
CORPUS = Path.home() / "sam-corpus" / "disks"
REPO_DOCS = REPO_ROOT / "docs"
DEFAULT_REF = Path.home() / "git" / "samdos" / "res" / "samdos2.reference.bin"


def sha256_full(b: bytes) -> str:
    return hashlib.sha256(b).hexdigest()


def find_matching_disks(target_sha: str) -> list[Path]:
    """Walk CORPUS, return paths to disks whose slot-0 body SHA == target."""
    matches: list[Path] = []
    for p in sorted(CORPUS.glob("*.mgt")):
        data = p.read_bytes()
        if len(data) != DISK_SIZE or not is_bootable(data):
            continue
        info = decode_slot0_loadexec(data)
        if info is None:
            continue
        _load, _exec, _length, body = info
        if sha256_full(body) == target_sha:
            matches.append(p)
    return matches


def build_samfile() -> Path:
    """Build the samfile binary from REPO_ROOT, return its path."""
    out = Path(tempfile.gettempdir()) / "samfile-strict-sha-scan"
    print(f"==> building samfile -> {out}")
    subprocess.run(
        ["go", "build", "-o", str(out), "./cmd/samfile"],
        cwd=REPO_ROOT, check=True,
    )
    return out


def run_verify(samfile: Path, disk: Path) -> list[dict]:
    """Run samfile verify --format jsonl on one disk. Return parsed events."""
    proc = subprocess.run(
        [str(samfile), "verify", "--format", "jsonl", "-i", str(disk)],
        check=False, capture_output=True, text=True,
    )
    events = []
    for line in proc.stdout.splitlines():
        line = line.strip()
        if not line:
            continue
        try:
            events.append(json.loads(line))
        except json.JSONDecodeError as e:
            print(f"warn: parse fail on {disk.name}: {e}", file=sys.stderr)
    return events


def main() -> None:
    ap = argparse.ArgumentParser()
    ap.add_argument("--ref-bin", type=Path, default=DEFAULT_REF,
                    help=f"reference binary (default {DEFAULT_REF})")
    ap.add_argument("--sha", type=str, default=None,
                    help="explicit sha256; overrides --ref-bin")
    ap.add_argument("--out", type=Path, default=None,
                    help="output markdown path")
    args = ap.parse_args()

    if args.sha:
        target_sha = args.sha.lower()
        if len(target_sha) != 64:
            sys.exit("--sha must be a full 64-char sha256 hex")
        source_label = f"explicit sha256 `{target_sha}`"
    else:
        if not args.ref_bin.exists():
            sys.exit(f"reference binary not found: {args.ref_bin}")
        target_sha = sha256_full(args.ref_bin.read_bytes())
        source_label = f"`{args.ref_bin}` (sha256 `{target_sha}`)"

    sha16 = target_sha[:16]
    print(f"target sha = {target_sha}")
    print(f"source     = {source_label}")
    print(f"finding matching disks in {CORPUS} ...")
    disks = find_matching_disks(target_sha)
    if not disks:
        sys.exit(f"no corpus disks have slot-0 body matching {target_sha}")
    print(f"found {len(disks)} disks")

    samfile = build_samfile()

    # Run verify on each disk, collect events. Message and severity
    # for fail events live under `ev["finding"]`, not at the top level.
    all_events: dict[str, list[dict]] = {}
    rule_counts: dict[str, dict[str, int]] = defaultdict(lambda: defaultdict(int))
    rule_severity: dict[str, str] = {}
    rule_citation: dict[str, str] = {}
    rule_fails: dict[str, list[dict]] = defaultdict(list)
    rule_message: dict[str, set] = defaultdict(set)
    fired_rules: set[str] = set()
    print("==> running verify on each disk")
    for disk in disks:
        events = run_verify(samfile, disk)
        all_events[disk.stem] = events
        for ev in events:
            rid = ev.get("rule_id", "?")
            outcome = ev.get("outcome", "?")
            rule_counts[rid][outcome] += 1
            finding = ev.get("finding") or {}
            if finding.get("Severity"):
                rule_severity[rid] = finding["Severity"]
            if finding.get("Citation"):
                rule_citation[rid] = finding["Citation"]
            if outcome == "fail":
                fired_rules.add(rid)
                rule_fails[rid].append({"disk": disk.stem, **ev})
                msg = finding.get("Message")
                if msg:
                    rule_message[rid].add(msg)
    print(f"  collected {sum(len(v) for v in all_events.values())} events "
          f"across {len(disks)} disks")
    print(f"  rules that fired (fail): {len(fired_rules)}")

    # Build the report.
    out_path = args.out or (REPO_DOCS / f"strict-sha-scan-{sha16}.md")
    md = [
        f"# Strict-SHA scan: every rule fire on `{sha16}` disks",
        "",
        f"Reference binary: {source_label}.",
        "",
        f"Each of the **{len(disks)} disks** in the corpus has a slot-0",
        f"body whose sha256 matches the reference binary exactly. The",
        f"current branch's `samfile verify` was run against each. The",
        f"report below lists every rule that fired at least once on this",
        f"cohort, with the specific findings.",
        "",
        f"Because the cohort is one exact SAMDOS-2 build, any fire is",
        f"either:",
        f"",
        f"- a real corruption / writer bug on the specific disk, or",
        f"- a rule documenting a writer convention that SAMDOS-2 itself",
        f"  does not enforce when SAVE'ing files (i.e. a false-positive",
        f"  rule that should be rewritten or removed).",
        "",
        "## Summary",
        "",
        f"- **Cohort size:** {len(disks)} disks",
        f"- **Total events emitted:** {sum(len(v) for v in all_events.values())}",
        f"- **Rules with at least one fire (fail):** {len(fired_rules)}",
        "",
        "## Disks in cohort",
        "",
    ]
    for d in sorted(disks):
        md.append(f"- `{d.stem}`")
    md.append("")
    md.append("## Rules that fired")
    md.append("")
    if not fired_rules:
        md.append("None. Every rule passed or was not applicable on every")
        md.append("disk in the cohort. This is the desired outcome for a")
        md.append("rule set aligned with the canonical SAMDOS-2 SAVE path.")
        md.append("")

    # Order by fail count desc.
    fired_sorted = sorted(
        fired_rules,
        key=lambda r: -rule_counts[r].get("fail", 0),
    )
    for rid in fired_sorted:
        counts = rule_counts[rid]
        passes = counts.get("pass", 0)
        fails = counts.get("fail", 0)
        na = counts.get("not_applicable", 0)
        sev = rule_severity.get(rid, "?")
        citation = rule_citation.get(rid, "")
        total_applicable = passes + fails
        rate = 100 * fails / total_applicable if total_applicable else 0.0
        md.append(f"### `{rid}` — {fails}/{total_applicable} fails "
                  f"({rate:.1f}%), severity `{sev}`")
        md.append("")
        if citation:
            md.append(f"Source citation: `{citation}`")
            md.append("")
        if rule_message[rid]:
            md.append("**Distinct failure messages:**")
            md.append("")
            for m in sorted(rule_message[rid]):
                md.append(f"- `{m}`")
            md.append("")
        md.append("**Every fire (disk, slot/ref, filename, message):**")
        md.append("")
        md.append("| Disk | Ref | Filename | Message |")
        md.append("|---|---|---|---|")
        for f in sorted(rule_fails[rid], key=lambda f: (f.get("disk", ""), f.get("ref", ""))):
            disk = f.get("disk", "?")
            ref = f.get("ref", "")
            finding = f.get("finding") or {}
            msg = finding.get("Message", "")
            loc = finding.get("Location") or {}
            filename = loc.get("Filename", "") or ""
            if len(msg) > 100:
                msg = msg[:97] + "..."
            md.append(f"| `{disk}` | `{ref}` | `{filename}` | {msg} |")
        md.append("")
        if na:
            md.append(f"_({na} not-applicable events for this rule on this cohort)_")
            md.append("")

    # Also include rules that fired only as "pass" — for completeness.
    pass_only_rules = sorted(
        rid for rid in rule_counts
        if rule_counts[rid].get("fail", 0) == 0
        and rule_counts[rid].get("pass", 0) > 0
    )
    if pass_only_rules:
        md.append("## Rules that only passed (never fired)")
        md.append("")
        md.append(f"({len(pass_only_rules)} rules — listed for completeness.)")
        md.append("")
        md.append("| Rule | passes | not-applicable |")
        md.append("|---|---:|---:|")
        for rid in pass_only_rules:
            c = rule_counts[rid]
            md.append(f"| `{rid}` | {c.get('pass', 0)} | {c.get('not_applicable', 0)} |")
        md.append("")

    # And rules that were never applicable to any disk in the cohort.
    not_app_only = sorted(
        rid for rid in rule_counts
        if rule_counts[rid].get("pass", 0) == 0
        and rule_counts[rid].get("fail", 0) == 0
        and rule_counts[rid].get("not_applicable", 0) > 0
    )
    if not_app_only:
        md.append("## Rules that were never applicable")
        md.append("")
        md.append(f"({len(not_app_only)} rules — never observed a subject they apply to.)")
        md.append("")
        for rid in not_app_only:
            md.append(f"- `{rid}`")
        md.append("")

    out_path.parent.mkdir(parents=True, exist_ok=True)
    out_path.write_text("\n".join(md) + "\n")
    print(f"==> wrote {out_path}")

    # Companion JSONL: every fail event from the cohort, one per line.
    # Use this for ad-hoc jq / grep on attrs that the markdown summary
    # doesn't surface.
    jsonl_path = out_path.with_suffix(".fails.jsonl")
    with jsonl_path.open("w") as fp:
        for rid in fired_sorted:
            for ev in rule_fails[rid]:
                # ev already has disk-stem prepended; remove the duplicate
                # disk key from the wrapper layer to keep one event = one line.
                clean = {k: v for k, v in ev.items() if k != "disk" or ev.get("disk") == ev.get("disk")}
                fp.write(json.dumps(clean) + "\n")
    print(f"==> wrote {jsonl_path}")


if __name__ == "__main__":
    main()
