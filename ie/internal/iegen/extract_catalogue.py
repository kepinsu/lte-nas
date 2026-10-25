#!/usr/bin/env python3
"""Extract the TS 24.301 V18.9.0 clause 9.9 IE catalogue.

Use a local ETSI PDF or layout-preserving text. No network requests are made.
Clause declarations take precedence over context-specific message-table bounds.
Inherited definitions fall back to the tables in clauses 8.2/8.3. Ambiguities
require a documented entry in catalogue_overrides.json; public tagged usage names
that cannot be derived uniquely require tagged_usage_overrides.json.
"""

import argparse
import csv
import difflib
import io
import json
import re
import sys
from dataclasses import dataclass
from pathlib import Path


SPECIFICATION = "3GPP TS 24.301 V18.9.0 (Release 18)"
SPECIFICATION_URL = (
    "https://www.etsi.org/deliver/etsi_ts/124300_124399/124301/"
    "18.09.00_60/ts_124301v180900p.pdf"
)
EXPECTED_COUNT = 129
DIRECTORY = Path(__file__).resolve().parent
FIELDS = ("kind", "name", "clause", "format", "minLength", "maxLength")
TAGGED_USAGE_FIELDS = ("kind", "tag", "clause")
CLAUSE = re.compile(r"(?m)^(9\.9\.[234]\.\d+[A-Z]?)\s+([^\n]+)")
DECLARATION = re.compile(
    r"\bis (?:an? )?type ([1346]) information element[^.]*\.", re.IGNORECASE
)
TABLE_ROW = re.compile(
    r"\b[MOC]\s+(TLV-E|LV-E|TLV|LV|TV|V|T)\s+"
    r"(\d+(?:/2)?(?:-(?:\d+|n))?)\s*$"
)
OVERHEAD = {1: 0, 2: 1, 3: 1, 4: 2, 6: 3}
MAX_LENGTH = {1: 1, 2: 0, 3: 65535, 4: 255, 6: 65535}


@dataclass(frozen=True)
class Envelope:
    format: int
    minimum: int
    maximum: int


@dataclass(frozen=True)
class TableEvidence:
    envelope: Envelope
    source: str
    tag: int | None
    table: str


@dataclass(frozen=True)
class Entry:
    kind: str
    name: str
    clause: str
    envelope: Envelope
    rule: str
    evidence: tuple[str, ...]

    def row(self):
        """Return the public CSV columns, without extraction-only metadata."""
        return (
            self.kind, self.name, self.clause, self.envelope.format,
            self.envelope.minimum, self.envelope.maximum,
        )


def read_source(path):
    """Read text as-is, or extract PDF text in pypdf's layout mode."""
    if path.suffix.lower() != ".pdf":
        return path.read_text(encoding="utf-8")
    try:
        from pypdf import PdfReader
    except ImportError as error:
        raise ValueError(
            "PDF input requires pypdf; install ie/internal/iegen/requirements.txt, "
            "or supply layout-preserving .txt input"
        ) from error
    reader = PdfReader(path)
    return "\n".join(page.extract_text(extraction_mode="layout") for page in reader.pages)


def kind_name(name):
    """PascalCase headings while preserving acronyms and internal capitals."""
    # UE's usage setting describes UE usage, not a separate 'S' Go identifier.
    words = re.findall(r"[A-Za-z0-9]+", name.replace("'s", ""))
    return "".join(word[0].upper() + word[1:] for word in words)


def table_envelope(wire_format, length):
    """Convert a message-table wire length to value-only octet bounds."""
    if (wire_format, length) in (("V", "1/2"), ("TV", "1")):
        # The CSV describes storage: one byte for a four-bit value.
        return Envelope(1, 1, 1)
    if wire_format == "T" and length == "1":
        return Envelope(2, 0, 0)
    if "/" in length:
        raise ValueError(f"unsupported fractional length: {wire_format} {length}")
    formats = {
        "V": (3, 0), "TV": (3, 1),
        "LV": (4, 1), "TLV": (4, 2),
        "LV-E": (6, 2), "TLV-E": (6, 3),
    }
    if wire_format not in formats:
        raise ValueError(f"unsupported format/length: {wire_format} {length}")
    category, overhead = formats[wire_format]
    limits = length.split("-", 1)
    minimum = int(limits[0]) - overhead
    maximum = MAX_LENGTH[category] if limits[-1] == "n" else int(limits[-1]) - overhead
    return Envelope(category, minimum, maximum)


def message_tables(text):
    """Associate each 8.2/8.3 table row with its following 9.9 clause reference."""
    result = {}
    lines = text.splitlines()
    table = None
    for number, line in enumerate(lines):
        heading = re.search(r"\bTable (8\.[23]\.[\dA-Za-z.]+):", line)
        if heading:
            table = heading[1]
        row = TABLE_ROW.search(line)
        if not row or table is None:
            continue
        # Names can wrap before the clause number. Stop at the next row so a
        # missing reference cannot accidentally bind to the next IE.
        end = number + 1
        while end < min(len(lines), number + 6) and not TABLE_ROW.search(lines[end]):
            end += 1
        references = re.findall(r"9\.9\.[234]\.\d+[A-Za-z]?", " ".join(lines[number:end]))
        if not references:
            continue
        clause = references[0].upper()
        tag_match = re.match(r"\s*([0-9A-F]{2}|[89A-F]-)\s+", line, re.IGNORECASE)
        tag = None
        if tag_match and row[1] in ("T", "TV", "TLV", "TLV-E"):
            tag = int(tag_match[1].replace("-", "0"), 16)
        evidence = TableEvidence(
            table_envelope(row[1], row[2]),
            f"Table {table}, {clause}: {line.strip()}",
            tag,
            table,
        )
        result.setdefault(clause, []).append(evidence)
    return result


def declared_envelope(declaration):
    """Convert an explicit clause declaration, whose lengths include the IEI."""
    category = int(declaration[1])
    if category == 1:
        return Envelope(1, 1, 1)
    lengths = [int(value) for value in re.findall(r"(\d+)\s+octets?\b", declaration[0])]
    if not lengths:
        return None
    overhead = OVERHEAD[category]
    minimum = lengths[0] - overhead
    if "minimum length" in declaration[0].lower():
        maximum = lengths[1] - overhead if len(lengths) > 1 else MAX_LENGTH[category]
    else:
        maximum = minimum
    return Envelope(category, minimum, maximum)


def select_envelope(clause, body, tables, overrides):
    """Use a documented override, an explicit clause, or unambiguous table rows."""
    declaration = DECLARATION.search(body)
    evidence = []
    if declaration:
        evidence.append("Clause declaration: " + " ".join(declaration[0].split()))
    evidence.extend(sorted({table.source for table in tables}))
    if clause in overrides:
        override = overrides[clause]
        required = {"format", "minLength", "maxLength", "source", "reason"}
        if (not isinstance(override, dict) or set(override) != required
                or not isinstance(override["source"], str) or not override["source"].strip()
                or not isinstance(override["reason"], str) or not override["reason"].strip()):
            raise ValueError(f"{clause}: override requires envelope, source and reason")
        evidence.extend(("Override source: " + override["source"], override["reason"]))
        return Envelope(override["format"], override["minLength"], override["maxLength"]), "documented override", tuple(evidence)
    if declaration:
        envelope = declared_envelope(declaration)
        if envelope is not None:
            return envelope, "explicit clause declaration; subtract IEI/length overhead", tuple(evidence)
        tables = [table for table in tables if table.envelope.format == int(declaration[1])]
    formats = {table.envelope.format for table in tables}
    if len(formats) != 1:
        raise ValueError(
            f"{clause}: missing or conflicting format evidence {sorted(formats)}; "
            "add a sourced override rather than guessing"
        )
    category = formats.pop()
    envelope = Envelope(
        category,
        min(table.envelope.minimum for table in tables),
        max(table.envelope.maximum for table in tables),
    )
    return envelope, "union of message-table value bounds; subtract wire overhead", tuple(evidence)


def extract_catalogue(text, overrides):
    """Return ordered catalogue entries with evidence; reject unsupported sources."""
    if not re.search(r"3GPP TS 24\.301 version 18\.9\.0 Release 18", text):
        raise ValueError(f"expected {SPECIFICATION}; other versions need a separate review")
    start = re.search(r"(?m)^9\.9\s+Other information elements\s*$", text)
    if start is None:
        raise ValueError("clause 9.9 body not found (a table of contents is not sufficient)")
    chapter = text[start.end():]
    end = re.search(r"(?m)^10\s+List of system parameters\s*$", chapter)
    if end is None:
        raise ValueError("end of clause 9.9 not found")
    chapter = chapter[:end.start()]
    tables = message_tables(text[:start.start()])
    sections = list(CLAUSE.finditer(chapter))
    clauses = {section[1] for section in sections}
    if len(clauses) != len(sections):
        raise ValueError("duplicate clause headings")
    if set(overrides) - clauses:
        raise ValueError(f"overrides for missing clauses: {sorted(set(overrides) - clauses)}")
    entries = []
    for number, section in enumerate(sections):
        clause, name = section[1], " ".join(section[2].split())
        next_start = sections[number + 1].start() if number + 1 < len(sections) else len(chapter)
        body = chapter[section.end():next_start]
        envelope, rule, evidence = select_envelope(clause, body, tables.get(clause, []), overrides)
        values = (envelope.format, envelope.minimum, envelope.maximum)
        if any(type(value) is not int for value in values) or envelope.format not in MAX_LENGTH:
            raise ValueError(f"{clause}: invalid envelope {envelope}")
        if not 0 <= envelope.minimum <= envelope.maximum <= MAX_LENGTH[envelope.format]:
            raise ValueError(f"{clause}: invalid value bounds {envelope}")
        if envelope.format in (1, 2, 3) and envelope.minimum != envelope.maximum:
            raise ValueError(f"{clause}: fixed-length format has variable bounds")
        entries.append(Entry(kind_name(name), name, clause, envelope, rule, evidence))
    if len({entry.kind for entry in entries}) != len(entries):
        raise ValueError("duplicate Go kind names")
    return entries


def render_csv(entries):
    """Use stable LF endings, column order and specification order."""
    output = io.StringIO(newline="")
    writer = csv.writer(output, lineterminator="\n")
    writer.writerow(FIELDS)
    writer.writerows(entry.row() for entry in entries)
    return output.getvalue()


def tagged_usage_key(tag, clause):
    """Return the stable override key for one public tagged usage."""
    return f"0x{tag:02x}:{clause}"


def tagged_evidence(entry, evidence):
    """Return whether one table row is a valid tagged usage of entry.

    Message tables occasionally contain an incorrect clause reference. The table
    still describes the wire layout, but it must not create a constructor for an
    incompatible catalogue definition.
    """
    if evidence.tag is None:
        return False
    envelope = evidence.envelope
    fixed = envelope.minimum if envelope.format in (1, 2, 3) else 0
    if envelope.format in (1, 2, 3) and envelope.minimum != envelope.maximum:
        raise ValueError(f"{evidence.source}: variable length in fixed format")
    if not 0 <= fixed <= 65535 or evidence.tag == 0 or (envelope.format == 3 and fixed == 0):
        raise ValueError(f"{evidence.source}: invalid tagged envelope")
    if (envelope.format == 1 and (evidence.tag < 0x80 or evidence.tag & 0x0f)) or (
            envelope.format != 1 and evidence.tag >= 0x80):
        raise ValueError(f"{evidence.source}: tag/format mismatch")
    if entry.envelope.format != envelope.format:
        return False
    if envelope.format in (1, 2, 3) and (
            entry.envelope.minimum != fixed or entry.envelope.maximum != fixed):
        return False
    return True


def render_tagged_usages(entries, tables, overrides):
    """Render stable public tagged usages from the message tables.

    The catalogue provides the internal envelope. Each compatible tagged table
    row produces one (TAG, clause) binding. `overrides` only assigns a public
    usage name where a definition is used under several TAGs; it never changes a
    wire value or a format.
    """
    usages = {}
    for entry in entries:
        for evidence in tables.get(entry.clause, []):
            if not tagged_evidence(entry, evidence):
                continue
            key = evidence.tag, entry.clause
            usages[key] = (entry, evidence.envelope)
    if not usages:
        return TAGGED_USAGE_FIELDS[0] + "," + TAGGED_USAGE_FIELDS[1] + "," + TAGGED_USAGE_FIELDS[2] + "\n"

    available_overrides = {tagged_usage_key(tag, clause) for tag, clause in usages}
    if set(overrides) - available_overrides:
        raise ValueError("tagged usage overrides for missing bindings: " + ", ".join(sorted(set(overrides) - available_overrides)))

    names = {}
    rows = []
    for (tag, clause), (entry, envelope) in usages.items():
        key = tagged_usage_key(tag, clause)
        name = overrides.get(key, entry.kind)
        if not isinstance(name, str) or not name or kind_name(name) != name:
            raise ValueError(f"{key}: tagged usage override must be a Go identifier")
        if name in names:
            raise ValueError(
                f"duplicate tagged usage name {name}: {names[name]} and {key}; add a tagged usage override"
            )
        names[name] = key
        fixed = envelope.minimum if envelope.format in (1, 2, 3) else 0
        rows.append((tag, envelope.format, fixed, clause, name))

    output = io.StringIO(newline="")
    writer = csv.writer(output, lineterminator="\n")
    writer.writerow(TAGGED_USAGE_FIELDS)
    for tag, _format, fixed, clause, name in sorted(rows):
        writer.writerow((name, f"0x{tag:02x}", clause))
    return output.getvalue()


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__, epilog="Source: " + SPECIFICATION_URL)
    parser.add_argument("source", type=Path, help="local TS 24.301 V18.9.0 .pdf or layout-preserving .txt")
    parser.add_argument("--output", type=Path, default=DIRECTORY / "catalogue.csv")
    parser.add_argument("--tagged-usages-output", type=Path, default=DIRECTORY / "tagged_usages.csv")
    parser.add_argument("--overrides", type=Path, default=DIRECTORY / "catalogue_overrides.json")
    parser.add_argument("--tagged-usage-overrides", type=Path, default=DIRECTORY / "tagged_usage_overrides.json")
    parser.add_argument("--check", action="store_true", help="compare without writing; exit 1 on any difference")
    parser.add_argument("--explain", metavar="CLAUSE_OR_KIND", help="show the derivation and source evidence for one IE")
    arguments = parser.parse_args(argv)
    try:
        overrides = json.loads(arguments.overrides.read_text(encoding="utf-8"))
        if not isinstance(overrides, dict):
            raise ValueError("overrides must be a JSON object indexed by clause")
        tagged_usage_overrides = json.loads(arguments.tagged_usage_overrides.read_text(encoding="utf-8"))
        if not isinstance(tagged_usage_overrides, dict):
            raise ValueError("tagged usage overrides must be a JSON object indexed by tag and clause")
        text = read_source(arguments.source)
        entries = extract_catalogue(text, overrides)
        if len(entries) != EXPECTED_COUNT:
            raise ValueError(f"expected {EXPECTED_COUNT} clauses, found {len(entries)}; refusing partial output")
        if arguments.explain:
            matches = [entry for entry in entries if arguments.explain in (entry.kind, entry.clause)]
            if not matches:
                raise ValueError(f"no such clause or kind: {arguments.explain}")
            entry = matches[0]
            print(f"{entry.clause} {entry.kind}: {entry.envelope}")
            print("Rule: " + entry.rule)
            for evidence in entry.evidence:
                print("  " + evidence)
        start = re.search(r"(?m)^9\.9\s+Other information elements\s*$", text)
        if start is None:
            raise ValueError("clause 9.9 body not found")
        outputs = (
            (arguments.output, render_csv(entries)),
            (arguments.tagged_usages_output, render_tagged_usages(
                entries, message_tables(text[:start.start()]), tagged_usage_overrides)),
        )
        different = False
        for path, result in outputs:
            if arguments.check:
                existing = path.read_text(encoding="utf-8")
                if existing != result:
                    different = True
                    sys.stdout.writelines(difflib.unified_diff(
                        existing.splitlines(keepends=True), result.splitlines(keepends=True),
                        fromfile=str(path), tofile="extracted " + path.name,
                    ))
            else:
                path.write_text(result, encoding="utf-8", newline="")
        if different:
            return 1
        print(f"{'Checked' if arguments.check else 'Wrote'} {len(entries)} IEs: {arguments.output}")
        print(f"{'Checked' if arguments.check else 'Wrote'} tagged usages: {arguments.tagged_usages_output}")
        return 0
    except (OSError, ValueError) as error:
        parser.exit(1, f"error: {error}\n")


if __name__ == "__main__":
    sys.exit(main())
