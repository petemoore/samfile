#!/usr/bin/env python3
"""Ingest JSONL CheckEvents from samfile --format jsonl into the
`checks` table of ~/sam-corpus/findings.db.

Reads: ~/sam-corpus/outputs-jsonl/<disk>.jsonl
Writes: ~/sam-corpus/findings.db (creates / replaces `checks` table).
Existing `disks` and `findings` tables are not touched.
"""

from __future__ import annotations

import json
import sqlite3
from pathlib import Path

CORPUS = Path.home() / "sam-corpus"
JSONL_DIR = CORPUS / "outputs-jsonl"
DB = CORPUS / "findings.db"

# Column schema for the checks table. Order matters — used to build
# the INSERT placeholders.
COLUMNS = [
    # event header
    ("disk", "TEXT"),
    ("rule_id", "TEXT"),
    ("scope", "TEXT"),
    ("ref", "TEXT"),
    ("outcome", "TEXT"),
    # disk-scope attrs
    ("dialect", "TEXT"),
    ("boot_signature_present", "INTEGER"),
    ("used_slot_count", "INTEGER"),
    # slot-scope attrs
    ("slot_index", "INTEGER"),
    ("filename", "TEXT"),
    ("file_type", "TEXT"),
    ("file_type_byte", "INTEGER"),
    ("file_length", "INTEGER"),
    ("page_offset_form", "TEXT"),
    ("pages", "INTEGER"),
    ("mgt_flags", "INTEGER"),
    ("first_track", "INTEGER"),
    ("first_sector", "INTEGER"),
    ("first_side", "INTEGER"),
    ("has_autorun_or_autoexec", "INTEGER"),
    ("dir_mirror_populated", "INTEGER"),
    ("slot_is_erased", "INTEGER"),
    ("file_type_info_hex", "TEXT"),
    ("sectors_count", "INTEGER"),
    # chain-step attrs (currently unused — reserved)
    ("chain_position", "TEXT"),
    ("chain_index", "INTEGER"),
    ("track", "INTEGER"),
    ("sector", "INTEGER"),
    ("side", "INTEGER"),
    ("next_track", "INTEGER"),
    ("next_sector", "INTEGER"),
    ("on_sam_map", "INTEGER"),
    ("on_dir_sam_map", "INTEGER"),
    ("distance_from_dir_tracks", "INTEGER"),
    # finding payload (NULL on pass / not_applicable)
    ("severity", "TEXT"),
    ("message", "TEXT"),
    ("citation", "TEXT"),
]

COL_NAMES = [c[0] for c in COLUMNS]


def column_value(event: dict, col: str):
    if col in ("disk", "rule_id", "scope", "ref", "outcome"):
        return event.get(col)
    attrs = event.get("attrs") or {}
    if col in attrs:
        v = attrs[col]
        if isinstance(v, bool):
            return int(v)
        return v
    finding = event.get("finding") or {}
    if col == "severity":
        return finding.get("Severity")
    if col == "message":
        return finding.get("Message")
    if col == "citation":
        return finding.get("Citation")
    return None


def main() -> None:
    conn = sqlite3.connect(DB)
    c = conn.cursor()
    c.execute("DROP TABLE IF EXISTS checks")
    c.execute(
        "CREATE TABLE checks ("
        + ", ".join(f"{name} {typ}" for name, typ in COLUMNS)
        + ")"
    )
    c.execute("CREATE INDEX idx_checks_rule ON checks(rule_id)")
    c.execute("CREATE INDEX idx_checks_outcome ON checks(outcome)")
    c.execute("CREATE INDEX idx_checks_disk ON checks(disk)")
    c.execute("CREATE INDEX idx_checks_type ON checks(file_type)")

    placeholders = ", ".join("?" for _ in COL_NAMES)
    insert = f"INSERT INTO checks ({', '.join(COL_NAMES)}) VALUES ({placeholders})"

    n = 0
    batch = []
    for path in sorted(JSONL_DIR.glob("*.jsonl")):
        for line in path.read_text(errors="replace").splitlines():
            if not line.strip():
                continue
            try:
                event = json.loads(line)
            except json.JSONDecodeError:
                continue
            event.setdefault("disk", path.stem)
            batch.append(tuple(column_value(event, col) for col in COL_NAMES))
            n += 1
            if len(batch) >= 5000:
                c.executemany(insert, batch)
                batch = []
    if batch:
        c.executemany(insert, batch)
    conn.commit()
    print(f"ingested {n} CheckEvents into {DB}")
    conn.close()


if __name__ == "__main__":
    main()
