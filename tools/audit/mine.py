#!/usr/bin/env python3
"""Generate audit reports from ~/sam-corpus/findings.db's `checks`
table.

Reports written under ~/sam-corpus/analyses/. Produced in fail-safe
order — coverage.md and disk-health.md (the fallback floor) come
first as plain pandas counts and never fail; the richer reports are
best-effort with graceful degradation.
"""

from __future__ import annotations

import sqlite3
import traceback
from pathlib import Path

import pandas as pd

CORPUS = Path.home() / "sam-corpus"
DB = CORPUS / "findings.db"
OUT = CORPUS / "analyses"


def load_checks() -> pd.DataFrame:
    conn = sqlite3.connect(DB)
    df = pd.read_sql_query("SELECT * FROM checks", conn)
    conn.close()
    return df


def report_coverage(checks: pd.DataFrame) -> None:
    """Per-rule applies/fails/fail-rate. The denominator-gap fix."""
    rows = []
    for rule_id, grp in checks.groupby("rule_id"):
        applicable = grp[grp["outcome"].isin(["pass", "fail"])]
        fails = grp[grp["outcome"] == "fail"]
        applies_n = len(applicable)
        fails_n = len(fails)
        passes_n = (grp["outcome"] == "pass").sum()
        na_n = (grp["outcome"] == "not_applicable").sum()
        disks = grp.loc[grp["outcome"] == "fail", "disk"].nunique()
        rate = (100.0 * fails_n / applies_n) if applies_n else float("nan")
        sev = fails["severity"].mode().iloc[0] if not fails.empty else ""
        scope = grp["scope"].mode().iloc[0] if not grp.empty else ""
        rows.append({
            "rule_id": rule_id,
            "severity": sev,
            "scope": scope,
            "applies": applies_n,
            "passes": passes_n,
            "fails": fails_n,
            "not_applicable": na_n,
            "fail_rate_pct": rate,
            "disks_affected": disks,
        })
    df = pd.DataFrame(rows).sort_values(
        ["fail_rate_pct", "fails"], ascending=[False, False], na_position="last"
    )
    md = [
        "# Per-rule coverage and failure rate",
        "",
        "Sorted by fail-rate (descending). `applies` = (pass + fail). `not_applicable` = subjects the rule deliberately skips.",
        "",
        "Rules where `scope == legacy` lack denominator metadata — only their fails are recorded; pass-rate is NaN.",
        "",
        "| Rule | Severity | Scope | Applies | Passes | Fails | N/A | Fail-rate | Disks |",
        "|---|---|---|---:|---:|---:|---:|---:|---:|",
    ]
    for _, r in df.iterrows():
        rate = "N/A" if pd.isna(r.fail_rate_pct) else f"{r.fail_rate_pct:.1f}%"
        md.append(
            f"| `{r.rule_id}` | {r.severity} | {r.scope} | {r.applies} | {r.passes} | {r.fails} | {r.not_applicable} | {rate} | {r.disks_affected} |"
        )
    (OUT / "coverage.md").write_text("\n".join(md) + "\n")
    print(f"wrote {OUT / 'coverage.md'} ({len(df)} rules)")


def report_disk_health(checks: pd.DataFrame) -> None:
    """Per-disk findings + structural-pass score."""
    rows = []
    for disk, grp in checks.groupby("disk"):
        fails = grp[grp["outcome"] == "fail"]
        total_findings = len(fails)
        fatal = (fails["severity"] == "fatal").sum()
        structural = (fails["severity"] == "structural").sum()
        # Pass rate over checks where the rule fired with structural/fatal anywhere.
        # Use the rule's modal severity from the corpus to bucket.
        rule_severity = (
            checks[checks["outcome"] == "fail"]
            .groupby("rule_id")["severity"].agg(lambda x: x.mode().iloc[0] if not x.mode().empty else "")
        ).to_dict()
        grp_with_sev = grp.copy()
        grp_with_sev["rule_severity"] = grp_with_sev["rule_id"].map(rule_severity).fillna("")
        struct = grp_with_sev[
            grp_with_sev["rule_severity"].isin(["structural", "fatal"]) &
            grp_with_sev["outcome"].isin(["pass", "fail"])
        ]
        sp = (struct["outcome"] == "pass").sum()
        sf = (struct["outcome"] == "fail").sum()
        struct_rate = sp / (sp + sf) if (sp + sf) else 1.0
        distinct_rules_fired = fails["rule_id"].nunique()
        dialect = grp["dialect"].dropna().mode().iloc[0] if not grp["dialect"].dropna().empty else ""
        rows.append({
            "disk": disk,
            "dialect": dialect,
            "total_findings": total_findings,
            "fatal": int(fatal),
            "structural": int(structural),
            "structural_pass_rate": struct_rate,
            "distinct_rules_fired": distinct_rules_fired,
        })
    df = pd.DataFrame(rows).sort_values(
        ["structural_pass_rate", "total_findings"], ascending=[True, False]
    )
    md = [
        "# Per-disk health",
        "",
        "Sorted by structural-pass-rate (ascending — worst first). Disks at the top with very low pass-rate and many distinct rules fired are candidate 'not really a disk' specimens that probably shouldn't dilute per-rule pattern analysis.",
        "",
        "| Disk | Dialect | Total findings | Fatal | Structural | Struct pass-rate | Distinct rules fired |",
        "|---|---|---:|---:|---:|---:|---:|",
    ]
    for _, r in df.iterrows():
        md.append(
            f"| {r.disk[:60]} | {r.dialect} | {r.total_findings} | {r.fatal} | {r.structural} | {r.structural_pass_rate:.2f} | {r.distinct_rules_fired} |"
        )
    (OUT / "disk-health.md").write_text("\n".join(md) + "\n")
    print(f"wrote {OUT / 'disk-health.md'} ({len(df)} disks)")


# Attribute columns considered for conditional / decision-tree analyses.
ATTR_COLS = [
    "dialect", "boot_signature_present", "used_slot_count",
    "file_type", "page_offset_form", "first_track", "first_sector",
    "first_side", "has_autorun_or_autoexec", "dir_mirror_populated",
    "slot_is_erased", "sectors_count", "file_length", "mgt_flags",
    "pages",
]


def report_conditional(checks: pd.DataFrame) -> None:
    """Per-rule conditional fail-rate per attribute value vs baseline."""
    rows = []
    for rule_id, grp in checks.groupby("rule_id"):
        applicable = grp[grp["outcome"].isin(["pass", "fail"])]
        if len(applicable) < 10:
            continue
        baseline = (applicable["outcome"] == "fail").mean()
        for col in ATTR_COLS:
            if col not in applicable.columns or applicable[col].isna().all():
                continue
            for val, sub in applicable.groupby(col):
                if len(sub) < 10:
                    continue
                cond = (sub["outcome"] == "fail").mean()
                support_disks = sub["disk"].nunique()
                if cond in (0.0, 1.0) or abs(cond - baseline) > 0.3:
                    rows.append({
                        "rule_id": rule_id,
                        "attribute": col,
                        "value": str(val),
                        "support": len(sub),
                        "support_disks": support_disks,
                        "baseline_fail_rate": baseline,
                        "conditional_fail_rate": cond,
                        "delta": cond - baseline,
                    })
    df = pd.DataFrame(rows).sort_values(
        ["rule_id", "conditional_fail_rate"], ascending=[True, False]
    )
    md = [
        "# Conditional failure rates per attribute",
        "",
        "Surfaces (rule, attribute, value) combinations where the conditional fail-rate differs substantially from the rule's baseline. |Δ| > 30pp OR conditional rate ∈ {0%, 100%}, support ≥ 10 events.",
        "",
        "| Rule | Attribute | Value | Support (events) | Support (disks) | Baseline | Conditional | Δ |",
        "|---|---|---|---:|---:|---:|---:|---:|",
    ]
    for _, r in df.iterrows():
        md.append(
            f"| `{r.rule_id}` | {r.attribute} | {r.value} | {r.support} | {r.support_disks} | {r.baseline_fail_rate*100:.1f}% | {r.conditional_fail_rate*100:.1f}% | {r.delta*100:+.1f}pp |"
        )
    (OUT / "conditional.md").write_text("\n".join(md) + "\n")
    print(f"wrote {OUT / 'conditional.md'} ({len(df)} rows)")


def report_disk_clusters(checks: pd.DataFrame) -> None:
    """Hierarchical clustering of disks by rule-fire vector."""
    from sklearn.cluster import AgglomerativeClustering

    pivot = (
        checks[checks["outcome"] == "fail"]
        .pivot_table(index="disk", columns="rule_id", values="ref", aggfunc="count", fill_value=0)
    )
    if pivot.empty:
        (OUT / "disk-clusters.md").write_text("# Disk clusters\n\nNo failure events to cluster.\n")
        return
    n_clusters = min(8, max(2, len(pivot) // 100))
    X = (pivot > 0).astype(int).values
    cl = AgglomerativeClustering(n_clusters=n_clusters)
    labels = cl.fit_predict(X)
    pivot["cluster"] = labels
    md = ["# Disk clusters by rule-fire pattern", "", f"{len(pivot)} disks partitioned into {n_clusters} clusters."]
    for cid, grp in pivot.groupby("cluster"):
        top_rules = (grp.drop(columns="cluster") > 0).sum().sort_values(ascending=False).head(8)
        md.append("")
        md.append(f"## Cluster {cid} — {len(grp)} disks")
        md.append("")
        md.append("Most-fired rules in this cluster (rule, disks-with-fail, share):")
        for rule, n in top_rules.items():
            md.append(f"- `{rule}` — {n}/{len(grp)} disks ({100*n/len(grp):.0f}%)")
        md.append("")
        md.append("Example disks:")
        for disk in list(grp.index)[:5]:
            md.append(f"- {disk[:60]}")
    (OUT / "disk-clusters.md").write_text("\n".join(md) + "\n")
    print(f"wrote {OUT / 'disk-clusters.md'} ({n_clusters} clusters)")


def report_patterns(checks: pd.DataFrame) -> None:
    """Per-rule decision-tree splits."""
    from sklearn.tree import DecisionTreeClassifier, export_text

    md = [
        "# Per-rule patterns (decision-tree splits)",
        "",
        "For each rule with ≥ 50 applicable checks and both pass + fail outcomes, a depth-4 decision tree is trained on the subject attributes to predict fail vs pass. The tree's splits show which attribute values drive failure.",
        "",
    ]
    feat_cols = [c for c in ATTR_COLS if c in checks.columns]
    for rule_id, grp in checks.groupby("rule_id"):
        applicable = grp[grp["outcome"].isin(["pass", "fail"])].copy()
        if len(applicable) < 50:
            continue
        y = (applicable["outcome"] == "fail").astype(int)
        if y.nunique() < 2:
            continue
        X = applicable[feat_cols].copy()
        for c in X.columns:
            if X[c].dtype == object:
                X[c] = X[c].astype("category").cat.codes
            else:
                X[c] = pd.to_numeric(X[c], errors="coerce").fillna(-1)
        try:
            tree = DecisionTreeClassifier(max_depth=4, min_samples_leaf=10).fit(X, y)
        except ValueError:
            continue
        md.append(f"## `{rule_id}` (applies={len(applicable)}, fail-rate={y.mean()*100:.1f}%)")
        md.append("")
        md.append("```")
        md.append(export_text(tree, feature_names=list(X.columns), max_depth=4).strip())
        md.append("```")
        md.append("")
    (OUT / "patterns.md").write_text("\n".join(md) + "\n")
    print(f"wrote {OUT / 'patterns.md'}")


def report_high_confidence(checks: pd.DataFrame) -> None:
    """Distill conditional fail-rates into the high-confidence action
    list: every (rule, attribute, value) where conditional fail-rate
    is 0% or 100% on support ≥ 10 distinct disks. These are the
    patterns the spec's autonomous fix loop is allowed to act on
    (after source-grounding)."""
    rows = []
    for rule_id, grp in checks.groupby("rule_id"):
        applicable = grp[grp["outcome"].isin(["pass", "fail"])]
        if len(applicable) < 10:
            continue
        baseline = (applicable["outcome"] == "fail").mean()
        for col in ATTR_COLS:
            if col not in applicable.columns or applicable[col].isna().all():
                continue
            for val, sub in applicable.groupby(col):
                support_disks = sub["disk"].nunique()
                if support_disks < 10:
                    continue
                cond = (sub["outcome"] == "fail").mean()
                if cond not in (0.0, 1.0):
                    continue
                rows.append({
                    "rule_id": rule_id,
                    "attribute": col,
                    "value": str(val),
                    "support_events": len(sub),
                    "support_disks": support_disks,
                    "conditional_fail_rate": cond,
                    "baseline_fail_rate": baseline,
                })
    if rows:
        df = pd.DataFrame(rows).sort_values(
            ["conditional_fail_rate", "support_disks"], ascending=[False, False]
        )
    else:
        df = pd.DataFrame(columns=[
            "rule_id", "attribute", "value", "support_events",
            "support_disks", "conditional_fail_rate", "baseline_fail_rate",
        ])
    md = [
        "# High-confidence patterns (act-on candidates)",
        "",
        "Each row is a (rule, attribute, value) slice where the conditional fail-rate is 0% or 100% on support ≥ 10 distinct disks. These meet the statistical half of the spec's act-on threshold; before acting on any one, ground it in source code (samdos / ROM disasm / samfile) per the design spec.",
        "",
        "| Rule | Attribute | Value | Conditional fail-rate | Baseline | Support (disks) | Support (events) |",
        "|---|---|---|---:|---:|---:|---:|",
    ]
    for _, r in df.iterrows():
        md.append(
            f"| `{r.rule_id}` | {r.attribute} | {r.value} | {r.conditional_fail_rate*100:.0f}% | {r.baseline_fail_rate*100:.1f}% | {r.support_disks} | {r.support_events} |"
        )
    if df.empty:
        md.append("| _no high-confidence patterns surfaced_ | | | | | | |")
    (OUT / "high-confidence-patterns.md").write_text("\n".join(md) + "\n")
    print(f"wrote {OUT / 'high-confidence-patterns.md'} ({len(df)} candidates)")


def main() -> None:
    OUT.mkdir(parents=True, exist_ok=True)
    checks = load_checks()
    print(f"loaded {len(checks)} CheckEvents")
    # Fallback floor — these must always succeed.
    report_coverage(checks)
    report_disk_health(checks)
    # Richer analyses — best-effort.
    for fn, name in [
        (report_conditional, "conditional.md"),
        (report_high_confidence, "high-confidence-patterns.md"),
        (report_disk_clusters, "disk-clusters.md"),
        (report_patterns, "patterns.md"),
    ]:
        try:
            fn(checks)
        except Exception:
            traceback.print_exc()
            print(f"({name} failed; check above)")


if __name__ == "__main__":
    main()
