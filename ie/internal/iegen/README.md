# LTE NAS IE catalogue generation

This directory contains every input and tool used for generation. Generated Go
files remain in the public `ie` package because they define its constants and
constructors.

| File | Purpose |
| --- | --- |
| `extract_catalogue.py` | Extracts the 129 IEs from TS 24.301 V18.9.0. |
| `catalogue_overrides.json` | Documented manual envelope exceptions. |
| `tagged_usage_overrides.json` | Explicit usage names when one definition is used under multiple tags. |
| `catalogue.csv` | Extracted catalogue consumed by the Go generator. |
| `tagged_usages.csv` | Public tagged usage names, their tags, and source clauses. |
| `requirements.txt` | `pypdf` version used for PDF input. |
| `main.go` | Generates definitions, constants, constructors, tests, and the private Go registry. |
| `test_extract_catalogue.py` | Tests extraction rules and the script interface. |

## Source and full regeneration

Source: [ETSI TS 124 301 V18.9.0, Release 18, January 2025](https://www.etsi.org/deliver/etsi_ts/124300_124399/124301/18.09.00_60/ts_124301v180900p.pdf).

From the repository root, using Python **3.10 or later**:

```sh
WORK_DIR="$(mktemp -d)"
python3 -m venv "$WORK_DIR/venv"
"$WORK_DIR/venv/bin/python" -m pip install -r ie/internal/iegen/requirements.txt

curl --fail --location \
  https://www.etsi.org/deliver/etsi_ts/124300_124399/124301/18.09.00_60/ts_124301v180900p.pdf \
  --output "$WORK_DIR/ts_124301v180900p.pdf"

"$WORK_DIR/venv/bin/python" ie/internal/iegen/extract_catalogue.py \
  "$WORK_DIR/ts_124301v180900p.pdf"

go generate ./ie
```

The PDF is supplied explicitly to the script: neither extraction nor
`go generate` downloads it, and it is not kept in the repository. Default CSV
and override paths are relative to the script directory, not the current
directory. `go generate ./ie` runs from the repository root, while Go starts the
generator in the `ie` directory.

The script writes `catalogue.csv` and `tagged_usages.csv`. For a given release
and unchanged override files, its output is deterministic: `--check` must emit
no diff. When changing release, review both override files and regenerate the
CSVs. `go generate` consumes only those two CSV files and generates the public
tag constants, tagged constructors, `NewPositional…` constructors, and a private
registry derived from the usages. Positional IEs do not appear in that registry.

To use an existing text extraction, `pypdf` is not needed:

```sh
# Optional, if pdftotext is installed: preserve table layout.
pdftotext -layout ts_124301v180900p.pdf ts_124301v180900p.txt

python3 ie/internal/iegen/extract_catalogue.py ts_124301v180900p.txt
go generate ./ie
```

PDF mode uses `pypdf==6.21.1` with `extraction_mode="layout"`. Text without
layout preservation, another TS version, or an incomplete extraction is not an
equivalent input. The script fails when it cannot recognise the document, its
bounds, its 129 clauses, or the information required for one IE. It never
replaces CSV files with partial output.

## Check or explain without rewriting

`--check` compares the extraction with both CSV files and prints a diff on a
mismatch. It exits with `0` when they match and `1` otherwise.

```sh
python3 ie/internal/iegen/extract_catalogue.py ts_124301v180900p.txt --check

# A Go identifier or clause number is accepted.
python3 ie/internal/iegen/extract_catalogue.py ts_124301v180900p.txt \
  --check --explain ExtendedQualityOfService
python3 ie/internal/iegen/extract_catalogue.py ts_124301v180900p.txt \
  --check --explain 9.9.3.16
```

For `ExtendedQualityOfService`, clause 9.9.4.30 declares a type-4 IE of
**12 octets in total**. The CSV therefore records **10 value octets** after
subtracting the IEI and length field. `--explain` displays that declaration and
the relevant message-table rows.

`--output` selects the catalogue, `--tagged-usages-output` selects public usage
output, and `--overrides` and `--tagged-usage-overrides` select alternative
override files. `--help` describes every option. None of these modes reads an
existing CSV to derive values: CSV files are only comparison targets for
`--check`.

## Derivation rules

1. Clause headings 9.9.2, 9.9.3, and 9.9.4 provide the inventory and order.
   Go names are converted to PascalCase while preserving acronyms; possessive
   `UE's` becomes `UE` in the identifier. The display name retains TS casing.
2. An explicit declaration in a clause takes precedence over bounds specific to
   a message occurrence. The evidence remains visible through `--explain`,
   including differing table rows.
3. Without explicit bounds, rows from tables 8.2/8.3 are associated with their
   clause 9.9 reference. Formats must be compatible; the catalogue records the
   union of value bounds, not message-specific presence constraints.
4. Value lengths exclude the wire header:

   | Encoding | Octets subtracted from total length |
   | --- | ---: |
   | V | 0 |
   | TV, type 3 | 1 |
   | LV | 1 |
   | TLV | 2 |
   | LV-E | 2 |
   | TLV-E | 3 |

   `V 1/2` and `TV 1` denote a type-1 IE, stored in one CSV octet. For an `n`
   bound, the maximum representable value is 255 octets for type 4 or 65535
   octets for type 6.
5. Missing or conflicting format evidence is an error unless a documented
   `catalogue_overrides.json` exception applies.

### Public tagged usages

`tagged_usages.csv` is extracted from T/TV/TLV/TLV-E rows in tables 8.2/8.3
whose envelope is compatible with the corresponding clause 9.9 definition.
Positional V/LV/LV-E rows are excluded. Type-1 tags are left aligned (`D-`
becomes `0xd0`). Its columns are:

- `kind`: public tag and constructor name;
- `tag`: on-wire IEI;
- `clause`: internal definition whose envelope is used.

By default, `kind` uses the clause identifier. When that name would be repeated
because one definition is reused under several tags, the script requires a name
in `tagged_usage_overrides.json`. That file cannot add a usage absent from the
tables or modify a tag or format.

`catalogue.csv` remains the source of internal definitions used by constructors,
validation, and positional IEs. The Go generator derives the private registry
from the two CSV files: identical formats for a tag are grouped, while differing
formats mark that tag as ambiguous.

Constructors remain distinct and store the real tag in `IE.Kind`. When formats
or fixed lengths differ, the private registry marks the tag as ambiguous:
`Parse` and `ParseMultiIEs` return `ErrAmbiguousIE` for tags `32`, `58`, `5B`,
and `5C` rather than picking an envelope.

An inconsistent clause reference in a table does not create a usage. For
example, the T3324 TLV row that references 9.9.3.16 does not turn the GPRS timer
into TV; only the compatible association with clause 9.9.3.16A produces
`T3324Value`.

Unknown-tag skipping is not represented by registry entries. It is implemented
in `sequence.go` according to TS 24.007 V18.1.0 clauses 11.2.4–11.2.5 and
TS 24.301 clauses 7.5.1/7.6.1. Unknown IEIs `00..0F` and `7E..7F` are
comprehension-required. Other tags are skipped according to their inferred EPS
format. In Release 18, `7C` and `7D` are no longer comprehension-required,
unlike in earlier TS versions.

### Current manual exceptions

- **9.9.2.3 — Mobile identity**: LTE tables cover only some usages. The generic
  `1..9` bound comes from the referenced TS 24.008 clause 10.5.1.4 definition.
  The script does **not** extract TS 24.008, so this decision is maintained
  explicitly.
- **9.9.3.16 — GPRS timer**: table 8.2.29.1 presents T3324 as a `GPRS timer 2`
  TLV but references 9.9.3.16 rather than 9.9.3.16A. The override keeps
  `GPRS timer` as type 3 with one value octet.

These exceptions do not replace extraction of the remaining rows. Adding a
release requires reviewing sources and exceptions, rather than only changing the
expected count or copying a CSV.

## Tests

```sh
python3 -m unittest discover -s ie/internal/iegen -p 'test_*.py' -v
go generate ./ie
go test -race ./...
go vet ./...
```

Python tests cover headings, bounds, half-octets, multiline references,
conflicts, documented exceptions, tagged usage generation, and the `--check` /
`--explain` modes. Verification against the official PDF is kept separate using
`--check` so tests do not download external data.
