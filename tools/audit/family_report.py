#!/usr/bin/env python3
"""Filter the audit's `checks` table to a single DOS family and
re-compute per-rule pass/fail rates restricted to that cohort.

Use this to answer: "if we only look at SAMDOS-2 disks (which our
rules were calibrated against), which rules still misfire?" Anything
that fires significantly inside a clean cohort is either genuinely
broken on that DOS or catching real corruption / writer bugs.

Usage:
    family_report.py [--rank N | --head SHA] [--threshold 1.5]

  --rank N      Family rank from family_table.py output (1 = top, default).
  --head SHA    Specific family-head SHA (overrides --rank).
  --threshold   Family-membership threshold percent (default 1.5).

Output: docs/family-coverage-<head>.md plus a stdout summary.
"""

from __future__ import annotations

import argparse
import hashlib
import sqlite3
import sys
from collections import Counter, defaultdict
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))
from family_table import (
    cluster_into_families,
    collect_variants,
)
from survey_dos import (
    decode_slot0_loadexec,
    is_bootable,
    DISK_SIZE,
)

REPO_DOCS = Path.home() / "git" / "samfile" / "docs"
DB = Path.home() / "sam-corpus" / "findings.db"


def main() -> None:
    ap = argparse.ArgumentParser()
    ap.add_argument("--rank", type=int, default=1,
                    help="family rank to inspect (1 = top, default). Ignored if --head or --strict-sha is set.")
    ap.add_argument("--head", type=str, default=None,
                    help="family-head SHA prefix to inspect (uses --threshold clustering).")
    ap.add_argument("--strict-sha", type=str, default=None,
                    help="exact slot-0 SHA prefix; report on disks with EXACTLY this SHA, no clustering. "
                         "Use this to test whether convention rules fire even within a single-build cohort.")
    ap.add_argument("--threshold", type=float, default=1.5,
                    help="family-membership %% threshold (only used without --strict-sha).")
    args = ap.parse_args()

    print(f"collecting slot-0 variants ...")
    variants = collect_variants()

    if args.strict_sha:
        # No clustering: take disks whose slot-0 SHA exactly matches.
        matches = [h for h in variants if h.startswith(args.strict_sha)]
        if not matches:
            sys.exit(f"no slot-0 SHA matches prefix {args.strict_sha!r}")
        if len(matches) > 1:
            print("multiple matching SHAs — pick a longer prefix:")
            for h in matches[:10]:
                print(f"  {h}  {len(variants[h]['disks'])} disks")
            sys.exit(1)
        head = matches[0]
        family_disks = set(variants[head]["disks"])
        mode_label = f"strict-SHA cohort `{head}`"
        print(f"strict-SHA mode: head={head}, {len(family_disks)} disks (no clustering)")
    else:
        print(f"clustering corpus into families at threshold {args.threshold}% ...")
        families = cluster_into_families(variants, args.threshold)
        family_list = []
        for h_head, vs in families.items():
            disks: list[str] = []
            for h in vs:
                disks.extend(variants[h]["disks"])
            family_list.append({"head": h_head, "disks": disks})
        family_list.sort(key=lambda f: -len(f["disks"]))
        if args.head:
            target = next((f for f in family_list if f["head"].startswith(args.head)), None)
            if target is None:
                sys.exit(f"no family with head SHA prefix {args.head!r}")
        else:
            if args.rank < 1 or args.rank > len(family_list):
                sys.exit(f"--rank out of range; only {len(family_list)} families")
            target = family_list[args.rank - 1]
        head = target["head"]
        family_disks = set(target["disks"])
        mode_label = f"family `{head}` (threshold {args.threshold}%)"
        print(f"family mode: head={head}, {len(family_disks)} disks")
    print()

    if not DB.exists():
        sys.exit(f"findings.db not found at {DB}; run tools/audit/run_audit.sh first")

    conn = sqlite3.connect(DB)
    c = conn.cursor()

    # Per-rule (applies, fails) over the family's disks vs the whole corpus.
    placeholders = ",".join("?" * len(family_disks))
    fam_q = f"""
        SELECT rule_id, outcome, COUNT(*)
        FROM checks
        WHERE disk IN ({placeholders})
        GROUP BY rule_id, outcome
    """
    all_q = """
        SELECT rule_id, outcome, COUNT(*) FROM checks GROUP BY rule_id, outcome
    """

    fam_counts: dict[str, dict[str, int]] = defaultdict(lambda: defaultdict(int))
    for r, o, n in c.execute(fam_q, tuple(family_disks)):
        fam_counts[r][o] = n
    all_counts: dict[str, dict[str, int]] = defaultdict(lambda: defaultdict(int))
    for r, o, n in c.execute(all_q):
        all_counts[r][o] = n

    # Per-rule severity (from a fail row's severity).
    sev_q = """
        SELECT rule_id, severity, COUNT(*) FROM checks
        WHERE outcome = 'fail' GROUP BY rule_id, severity
    """
    sev_by_rule: dict[str, str] = {}
    for r, s, n in c.execute(sev_q):
        if not sev_by_rule.get(r) or n > sev_by_rule.get(r, ("", 0))[1]:
            sev_by_rule[r] = (s, n)
    sev_by_rule = {k: v[0] for k, v in sev_by_rule.items()}

    conn.close()

    rows = []
    for rule_id in sorted(set(fam_counts) | set(all_counts)):
        fp = fam_counts[rule_id].get("pass", 0)
        ff = fam_counts[rule_id].get("fail", 0)
        ap = all_counts[rule_id].get("pass", 0)
        af = all_counts[rule_id].get("fail", 0)
        fam_rate = (100 * ff / (fp + ff)) if (fp + ff) else None
        all_rate = (100 * af / (ap + af)) if (ap + af) else None
        rows.append({
            "rule": rule_id,
            "severity": sev_by_rule.get(rule_id, ""),
            "fam_applies": fp + ff,
            "fam_fails": ff,
            "fam_fail_rate": fam_rate,
            "all_applies": ap + af,
            "all_fails": af,
            "all_fail_rate": all_rate,
        })

    # Sort by family fail-rate desc, treating None as -1.
    rows.sort(key=lambda r: -(r["fam_fail_rate"] or -1))

    md = [
        f"# Per-rule coverage on {mode_label}",
        "",
        f"This report restricts the audit pipeline's `checks` table to",
        f"the {len(family_disks)} disks in this cohort and recomputes",
        f"per-rule pass/fail rates. Compare the **cohort fail-rate**",
        f"column to the **all-corpus fail-rate**: if a rule was calibrated",
        f"against this DOS, the cohort fail-rate should be at or near 0%.",
        f"",
        f"Rules where the cohort fail-rate is materially > 0% are either:",
        f"",
        f"- genuinely buggy / over-strict on this DOS itself, or",
        f"- catching real corruption / writer bugs on the specific disks",
        f"  in this cohort.",
        "",
        f"With `--strict-sha`, the cohort contains only disks whose slot-0",
        f"body matches a single exact SHA — no clustering. This is the",
        f"falsifiable test for whether a 'convention' rule is enforced by",
        f"a single SAMDOS-2 build: if it fires in this cohort, the build",
        f"itself doesn't enforce it.",
        "",
        f"## Summary",
        "",
        f"- **Cohort:** {mode_label}",
        f"- **Disks in cohort:** {len(family_disks)}",
        f"- **Rules with ≥ 1 cohort fail:** {sum(1 for r in rows if r['fam_fails'] > 0)}",
        f"- **Rules with 0 cohort fails (clean):** {sum(1 for r in rows if r['fam_applies'] > 0 and r['fam_fails'] == 0)}",
        "",
        "## Table (sorted by cohort fail-rate, desc)",
        "",
        "| Rule | Severity | Cohort applies | Cohort fails | Cohort % | All-corpus % | Δ |",
        "|---|---|---:|---:|---:|---:|---:|",
    ]
    for r in rows:
        fam = "N/A" if r["fam_fail_rate"] is None else f"{r['fam_fail_rate']:.1f}%"
        alc = "N/A" if r["all_fail_rate"] is None else f"{r['all_fail_rate']:.1f}%"
        if r["fam_fail_rate"] is not None and r["all_fail_rate"] is not None:
            delta = f"{r['fam_fail_rate'] - r['all_fail_rate']:+.1f}pp"
        else:
            delta = "—"
        md.append(
            f"| `{r['rule']}` | {r['severity']} | {r['fam_applies']} | {r['fam_fails']} | {fam} | {alc} | {delta} |"
        )

    suffix = f"-strict-{head}" if args.strict_sha else f"-{head}"
    out = REPO_DOCS / f"family-coverage{suffix}.md"
    out.write_text("\n".join(md) + "\n")
    print(f"wrote {out}")
    print()
    print(f"{'Rule':40s}  {'Sev':>13s}  {'FamAppl':>7s}  {'FamFail':>7s}  {'Fam%':>7s}  {'All%':>7s}  Δ")
    print("=" * 110)
    for r in rows:
        fam = "    -" if r["fam_fail_rate"] is None else f"{r['fam_fail_rate']:6.2f}%"
        alc = "    -" if r["all_fail_rate"] is None else f"{r['all_fail_rate']:6.2f}%"
        if r["fam_fail_rate"] is not None and r["all_fail_rate"] is not None:
            d = f"{r['fam_fail_rate'] - r['all_fail_rate']:+5.1f}pp"
        else:
            d = "    -"
        print(f"{r['rule']:40s}  {r['severity']:>13s}  {r['fam_applies']:>7d}  {r['fam_fails']:>7d}  {fam}  {alc}  {d}")


if __name__ == "__main__":
    main()
