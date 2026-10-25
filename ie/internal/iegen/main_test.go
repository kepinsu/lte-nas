package main

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
)

const testCatalogue = `kind,name,clause,format,minLength,maxLength
GPRSTimer,GPRS timer,9.9.3.16,3,1,1
	GPRSTimer2,GPRS timer 2,9.9.3.16A,4,1,1
	ESMCause,ESM cause,9.9.4.4,3,1,1
	UENetworkCapability,UE network capability,9.9.3.34,4,2,13
	DeviceProperties,Device properties,9.9.2.0A,1,1,1
`

const testTaggedUsages = `kind,tag,clause
T3402Value,0x17,9.9.3.16
T3346Value,0x5f,9.9.3.16A
ESMCause,0x58,9.9.4.4
UENetworkCapability,0x58,9.9.3.34
DeviceProperties,0xc0,9.9.2.0A
`

func TestParseTaggedUsages(t *testing.T) {
	t.Parallel()
	definitions := parseTestCSV(t, testCatalogue)
	got, err := parseTaggedUsages(definitions, parseTestCSV(t, testTaggedUsages))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 5 || got[0].kind != "T3402Value" || got[0].tag != 0x17 || got[0].definition[0] != "GPRSTimer" {
		t.Fatalf("usages = %+v", got)
	}
}

func TestParseTaggedUsagesRejectsInvalidInput(t *testing.T) {
	t.Parallel()
	definitions := parseTestCSV(t, testCatalogue)
	for _, tc := range []struct {
		name   string
		usages string
	}{
		{
			name: "unknown clause",
			usages: `kind,tag,clause
Unknown,0x17,9.9.3.999
`,
		},
		{
			name: "duplicate binding",
			usages: `kind,tag,clause
T3402Value,0x17,9.9.3.16
AnotherT3402Value,0x17,9.9.3.16
`,
		},
		{
			name: "duplicate name",
			usages: `kind,tag,clause
T3402Value,0x17,9.9.3.16
T3402Value,0x5f,9.9.3.16A
`,
		},
		{
			name: "tag format mismatch",
			usages: `kind,tag,clause
Invalid,0x80,9.9.3.16A
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if _, err := parseTaggedUsages(definitions, parseTestCSV(t, tc.usages)); err == nil {
				t.Fatal("parseTaggedUsages accepted invalid input")
			}
		})
	}
}

func TestGenerateTaggedFormats(t *testing.T) {
	t.Parallel()
	usages, err := parseTaggedUsages(parseTestCSV(t, testCatalogue), parseTestCSV(t, testTaggedUsages))
	if err != nil {
		t.Fatal(err)
	}
	generated, err := generateTaggedFormats(usages)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(generated, []byte("0x17: {format: IEFormatThree, length: 1}")) ||
		!bytes.Contains(generated, []byte("0x5f: {format: IEFormatFour, length: 0}")) ||
		!bytes.Contains(generated, []byte("0x58: {}")) {
		t.Fatalf("generated formats:\n%s", generated)
	}
	again, err := generateTaggedFormats(usages)
	if err != nil || !bytes.Equal(generated, again) {
		t.Fatalf("generation is not reproducible: %v", err)
	}
}

func TestGenerateTaggedFormatsRejectsEmptyUsages(t *testing.T) {
	t.Parallel()
	if generated, err := generateTaggedFormats(nil); err == nil || generated != nil {
		t.Fatalf("accepted empty usages: %s, %v", generated, err)
	}
}

func parseTestCSV(t *testing.T, text string) [][]string {
	t.Helper()
	rows, err := csv.NewReader(strings.NewReader(text)).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	return rows
}
