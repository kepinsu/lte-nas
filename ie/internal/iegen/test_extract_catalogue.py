"""Tests of extraction rules, with small synthetic specification excerpts."""

import contextlib
import csv
import io
import json
import tempfile
import unittest
from pathlib import Path
from unittest import mock

import extract_catalogue as extractor


def specification(clauses, tables=""):
    return (
        "3GPP TS 24.301 version 18.9.0 Release 18\n"
        "8.2 EPS mobility management messages\n" + tables + "\n"
        "9.9 Other information elements\n"
        "9.9.2 Common information elements\n" + clauses + "\n"
        "10 List of system parameters\n"
    )


class TableEnvelopeTest(unittest.TestCase):
    def test_wire_overhead(self):
        cases = (
            ("V", "1/2", (1, 1, 1)),
            ("TV", "1", (1, 1, 1)),
            ("T", "1", (2, 0, 0)),
            ("V", "16", (3, 16, 16)),
            ("TV", "6", (3, 5, 5)),
            ("LV", "2-101", (4, 1, 100)),
            ("TLV", "3-102", (4, 1, 100)),
            ("TLV", "5-n", (4, 3, 255)),
            ("LV-E", "2-n", (6, 0, 65535)),
            ("TLV-E", "4-65538", (6, 1, 65535)),
        )
        for wire_format, length, expected in cases:
            with self.subTest(format=wire_format, length=length):
                self.assertEqual(
                    extractor.table_envelope(wire_format, length),
                    extractor.Envelope(*expected),
                )

    def test_unsupported_fraction(self):
        with self.assertRaisesRegex(ValueError, "fractional"):
            extractor.table_envelope("LV", "1/2")

    def test_unsupported_tag_only_length(self):
        with self.assertRaisesRegex(ValueError, "unsupported format/length"):
            extractor.table_envelope("T", "2")


class ExtractionTest(unittest.TestCase):
    def test_explicit_declarations(self):
        text = specification(
            "9.9.2.0 Additional information\n"
            "The Additional information is a type 4 information element with a minimum length of 3 octets.\n"
            "9.9.3.7 Detach type\n"
            "The Detach type is a type 1 information element.\n"
            "9.9.3.25 Nonce\n"
            "The Nonce is a type 3 information element with a length of 5 octets.\n"
            "9.9.3.56 Ciphering key data\n"
            "The Ciphering key data is a type 6 information element, with a minimum length of 35 octets "
            "and a maximum length of 2291 octets.\n"
            "9.9.4.30 Extended quality of service\n"
            "The Extended quality of service is a type 4 information element with a length of 12 octets.\n"
        )
        entries = extractor.extract_catalogue(text, {})
        self.assertEqual(
            [entry.envelope for entry in entries],
            [extractor.Envelope(*value) for value in (
                (4, 1, 255), (1, 1, 1), (3, 4, 4), (6, 32, 2288), (4, 10, 10),
            )],
        )
        self.assertEqual(entries[-1].kind, "ExtendedQualityOfService")
        self.assertIn("12 octets", entries[-1].evidence[0])

    def test_inherited_definition_and_union(self):
        text = specification(
            "9.9.4.1 Access point name\nSee clause 10.5.6.1 in TS 24.008.\n",
            "Table 8.3.6.1: Mandatory fields\nAccess point name M LV 2-101\n9.9.4.1\n"
            "Table 8.3.20.1: Optional fields\n28 Access point name O TLV 3-102\n9.9.4.1\n",
        )
        entry, = extractor.extract_catalogue(text, {})
        self.assertEqual(entry.envelope, extractor.Envelope(4, 1, 100))
        self.assertEqual(len(entry.evidence), 2)
        self.assertTrue(entry.rule.startswith("union"))

    def test_explicit_bounds_override_context_specific_tables(self):
        text = specification(
            "9.9.4.20 Remote UE context list\n"
            "The Remote UE context list is a type 6 information element with a minimum length of 5 octets "
            "and a maximum length of 65538 octets.\n",
            "Table 8.3.23.1: Remote UE report\n79 Remote UE Context Connected O TLV-E 3-65538\n9.9.4.20\n",
        )
        entry, = extractor.extract_catalogue(text, {})
        self.assertEqual(entry.envelope, extractor.Envelope(6, 2, 65535))
        # Contradicting context-specific evidence is still shown by --explain.
        self.assertIn("TLV-E 3-65538", "\n".join(entry.evidence))

    def test_declared_category_and_table_lengths(self):
        text = specification(
            "9.9.3.15 ESM message container\nThe ESM message container is a type 6 information element.\n",
            "Table 8.2.1.1: Attach accept\nESM message container M LV-E 5-n\n9.9.3.15\n"
            "Table 8.2.6.1: Detach request\n78 ESM message container O TLV-E 3-n\n9.9.3.15\n",
        )
        entry, = extractor.extract_catalogue(text, {})
        self.assertEqual(entry.envelope, extractor.Envelope(6, 0, 65535))

    def test_wrapped_reference_and_next_row_boundary(self):
        text = specification(
            "9.9.3.44 Voice domain preference and UE's usage setting\nSee TS 24.008.\n",
            "Table 8.2.4.1: Attach request\n5D Voice domain preference O TLV 3\n"
            "and UE's usage setting\n9.9.3.44\nD- Device properties O TV 1\n9.9.2.0A\n",
        )
        entry, = extractor.extract_catalogue(text, {})
        self.assertEqual(entry.kind, "VoiceDomainPreferenceAndUEUsageSetting")
        self.assertEqual(entry.envelope, extractor.Envelope(4, 1, 1))

    def test_missing_evidence_is_not_guessed(self):
        for table in (
            "",
            "Table 8.2.4.1: Attach request\nUnreferenced IE O TLV 3\nNext IE O TV 1\n9.9.2.0A\n",
        ):
            with self.subTest(table=table):
                text = specification("9.9.4.1 Access point name\nSee TS 24.008.\n", table)
                with self.assertRaisesRegex(ValueError, "missing or conflicting"):
                    extractor.extract_catalogue(text, {})

    def test_conflicting_references_require_documented_override(self):
        text = specification(
            "9.9.3.16 GPRS timer\nSee TS 24.008.\n",
            "Table 8.2.1.1: Attach accept\n17 T3402 GPRS timer O TV 2\n9.9.3.16\n"
            "Table 8.2.29.1: Tracking area update request\n6A T3324 GPRS timer 2 O TLV 3\n9.9.3.16\n",
        )
        with self.assertRaisesRegex(ValueError, "conflicting"):
            extractor.extract_catalogue(text, {})
        overrides = json.loads((extractor.DIRECTORY / "catalogue_overrides.json").read_text())
        entry, = extractor.extract_catalogue(text, {"9.9.3.16": overrides["9.9.3.16"]})
        self.assertEqual(entry.envelope, extractor.Envelope(3, 1, 1))
        self.assertEqual(entry.rule, "documented override")
        self.assertIn("T3324", "\n".join(entry.evidence))

    def test_overrides_require_source_and_reason(self):
        text = specification("9.9.2.3 Mobile identity\nSee TS 24.008.\n")
        with self.assertRaisesRegex(ValueError, "source and reason"):
            extractor.extract_catalogue(text, {"9.9.2.3": {"format": 4, "minLength": 1, "maxLength": 9}})
        with self.assertRaisesRegex(ValueError, "missing clauses"):
            extractor.extract_catalogue(text, {"9.9.3.16": {}})
        with self.assertRaisesRegex(ValueError, "source and reason"):
            extractor.extract_catalogue(text, {"9.9.2.3": None})

    def test_invalid_source_or_bounds(self):
        text = specification(
            "9.9.4.30 Extended quality of service\n"
            "The Extended quality of service is a type 4 information element with a length of 12 octets.\n"
        )
        for source, message in (
            (text.replace("18.9.0", "17.0.0"), "other versions"),
            (text.replace("9.9 Other information elements", "Contents only"), "body not found"),
            (text.replace("10 List of system parameters", "Truncated"), "end of clause"),
            (text.replace("12 octets", "300 octets"), "invalid value bounds"),
            (text.replace("9.9.4.30 Extended", "9.9.4.30 Duplicate\n9.9.4.30 Extended"), "duplicate clause"),
        ):
            with self.subTest(message=message):
                with self.assertRaisesRegex(ValueError, message):
                    extractor.extract_catalogue(source, {})


class TaggedUsagesTest(unittest.TestCase):
    def render(self, text, overrides=None):
        entries = extractor.extract_catalogue(text, {})
        start = extractor.re.search(r"(?m)^9\.9\s+Other information elements\s*$", text)
        return list(csv.DictReader(io.StringIO(extractor.render_tagged_usages(
            entries, extractor.message_tables(text[:start.start()]), overrides or {},
        ))))

    def test_only_compatible_tagged_rows(self):
        text = specification(
            "9.9.4.1 Access point name\nSee clause 10.5.6.1 in TS 24.008.\n"
            "9.9.4.5 ESM information transfer flag\n"
            "The ESM information transfer flag is a type 1 information element.\n",
            "Table 8.3.14.1: PDN connectivity request\n"
            "Request type M V 1/2\n9.9.4.14\n"
            "Access point name M LV 2-101\n9.9.4.1\n"
            "D- ESM information transfer flag O TV 1\n9.9.4.5\n"
            "28 Access point name O TLV 3-102\n9.9.4.1\n"
            "Table 8.2.33.1: Control plane service request\n"
            "78 ESM message container O TLV-E 3-n\n9.9.3.15\n",
        )
        rows = self.render(text)
        self.assertEqual(
            [(row["kind"], row["tag"], row["clause"]) for row in rows],
            [("AccessPointName", "0x28", "9.9.4.1"), ("ESMInformationTransferFlag", "0xd0", "9.9.4.5")],
        )

    def test_overrides_disambiguate_definitions_reused_under_multiple_tags(self):
        text = specification(
            "9.9.3.25 Nonce\nThe Nonce is a type 3 information element with a length of 5 octets.\n",
            "Table 8.2.20.1: Authentication request\n"
            "55 Nonce UE O TV 5\n9.9.3.25\n"
            "56 Nonce MME O TV 5\n9.9.3.25\n",
        )
        with self.assertRaisesRegex(ValueError, "duplicate tagged usage name Nonce"):
            self.render(text)
        rows = self.render(text, {
            "0x55:9.9.3.25": "NonceUE",
            "0x56:9.9.3.25": "NonceMME",
        })
        self.assertEqual([(row["kind"], row["tag"]) for row in rows], [
            ("NonceUE", "0x55"), ("NonceMME", "0x56"),
        ])

    def test_incompatible_clause_reference_does_not_create_a_usage(self):
        text = specification(
            "9.9.3.16 GPRS timer\nThe GPRS timer is a type 3 information element with a length of 2 octets.\n",
            "Table 8.2.29.1: Tracking area update request\n"
            "6A T3324 GPRS timer 2 O TLV 3\n9.9.3.16\n",
        )
        self.assertEqual(self.render(text), [])

    def test_tag_format_mismatch_is_rejected(self):
        text = specification(
            "9.9.3.16 GPRS timer\nThe GPRS timer is a type 3 information element with a length of 2 octets.\n",
            "Table 8.2.1.1: Invalid example\nD0 IE O TV 2\n9.9.3.16\n",
        )
        with self.assertRaisesRegex(ValueError, "tag/format mismatch"):
            self.render(text)


class CommandTest(unittest.TestCase):
    def setUp(self):
        directory = tempfile.TemporaryDirectory()
        self.addCleanup(directory.cleanup)
        self.directory = Path(directory.name)
        self.source = self.directory / "spec.txt"
        self.source.write_text(specification(
            "9.9.4.30 Extended quality of service\n"
            "The Extended quality of service is a type 4 information element with a length of 12 octets.\n"
        ))
        self.overrides = self.directory / "overrides.json"
        self.overrides.write_text("{}")
        self.output = self.directory / "catalogue.csv"
        self.tagged_usages_output = self.directory / "tagged_usages.csv"
        self.tagged_usage_overrides = self.directory / "tagged_usage_overrides.json"
        self.tagged_usage_overrides.write_text("{}")
        self.arguments = [
            str(self.source), "--overrides", str(self.overrides),
            "--output", str(self.output), "--tagged-usages-output", str(self.tagged_usages_output),
            "--tagged-usage-overrides", str(self.tagged_usage_overrides),
        ]

    def test_write_check_explain_and_difference(self):
        with mock.patch.object(extractor, "EXPECTED_COUNT", 1), contextlib.redirect_stdout(io.StringIO()) as log:
            self.assertEqual(extractor.main(self.arguments), 0)
            original = self.output.read_bytes()
            self.assertEqual(extractor.main(self.arguments + ["--check", "--explain", "ExtendedQualityOfService"]), 0)
            self.assertEqual(self.output.read_bytes(), original)
            self.assertIn("subtract IEI/length overhead", log.getvalue())
            self.assertEqual(list(csv.reader(original.decode().splitlines()))[-1][-2:], ["10", "10"])
            self.output.write_text("incorrect catalogue\n")
            self.assertEqual(extractor.main(self.arguments + ["--check"]), 1)
            self.assertEqual(self.output.read_text(), "incorrect catalogue\n")
            self.assertIn("--- ", log.getvalue())

    def test_partial_catalogue_does_not_overwrite_output(self):
        self.output.write_text("keep existing data\n")
        with contextlib.redirect_stderr(io.StringIO()), self.assertRaises(SystemExit) as error:
            extractor.main(self.arguments)
        self.assertEqual(error.exception.code, 1)
        self.assertEqual(self.output.read_text(), "keep existing data\n")
        self.assertFalse(self.tagged_usages_output.exists())

    def test_tagged_usages_are_also_checked_without_writing(self):
        with mock.patch.object(extractor, "EXPECTED_COUNT", 1), contextlib.redirect_stdout(io.StringIO()) as log:
            self.assertEqual(extractor.main(self.arguments), 0)
            original = self.output.read_bytes()
            self.tagged_usages_output.write_text("incorrect tagged usages\n")
            self.assertEqual(extractor.main(self.arguments + ["--check"]), 1)
            self.assertEqual(self.output.read_bytes(), original)
            self.assertEqual(self.tagged_usages_output.read_text(), "incorrect tagged usages\n")
            self.assertIn("tagged_usages.csv", log.getvalue())

    def test_unknown_explanation_does_not_write_output(self):
        with mock.patch.object(extractor, "EXPECTED_COUNT", 1), contextlib.redirect_stderr(io.StringIO()), self.assertRaises(SystemExit):
            extractor.main(self.arguments + ["--explain", "Unknown"])
        self.assertFalse(self.output.exists())


if __name__ == "__main__":
    unittest.main()
