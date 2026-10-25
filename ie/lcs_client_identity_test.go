package ie_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/kepinsu/lte-nas/ie"
)

func TestLCSClientIdentityFields(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		fields   *ie.LCSClientIdentityFields
		wantText string
	}{
		{
			name:     "GSM 7-bit",
			fields:   &ie.LCSClientIdentityFields{Text: "SUBSCRIBER", Encoding: ie.NetworkNameEncodingGSM7, AddCountryInitials: false},
			wantText: "SUBSCRIBER",
		},
		{
			name:     "GSM 7-bit with extension",
			fields:   &ie.LCSClientIdentityFields{Text: "€", Encoding: ie.NetworkNameEncodingGSM7, AddCountryInitials: true},
			wantText: "€",
		},
		{
			name:     "UCS2",
			fields:   &ie.LCSClientIdentityFields{Text: "1A2B", Encoding: ie.NetworkNameEncodingUCS2, AddCountryInitials: false},
			wantText: "1A2B",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := tc.fields.Marshal()
			if err != nil {
				t.Fatal(err)
			}

			// Verify minimum size (header + at least one encoded char)
			if len(payload) < 2 {
				t.Fatalf("payload too short: %d", len(payload))
			}

			// Round-trip via IE accessor (tagged)
			tagged, err := ie.NewLCSClientIdentityWithFields(tc.fields)
			if err != nil {
				t.Fatal(err)
			}
			if tagged.Kind != ie.LCSClientIdentity {
				t.Fatalf("tagged IE kind = %d, want %d", tagged.Kind, ie.LCSClientIdentity)
			}
			decoded, err := tagged.LCSClientIdentity()
			if err != nil {
				t.Fatal(err)
			}
			if decoded.Text != tc.wantText {
				t.Fatalf("LCSClientIdentity() text = %q, want %q", decoded.Text, tc.wantText)
			}
			if decoded.Encoding != tc.fields.Encoding {
				t.Fatalf("LCSClientIdentity() encoding = %d, want %d", decoded.Encoding, tc.fields.Encoding)
			}
			if decoded.AddCountryInitials != tc.fields.AddCountryInitials {
				t.Fatalf("LCSClientIdentity() addCountryInitials = %v, want %v", decoded.AddCountryInitials, tc.fields.AddCountryInitials)
			}

			// Round-trip via positional IE
			pos, err := ie.NewPositionalLCSClientIdentityWithFields(tc.fields)
			if err != nil {
				t.Fatal(err)
			}
			decoded2, err := pos.LCSClientIdentity()
			if err != nil {
				t.Fatal(err)
			}
			if decoded2.Text != tc.wantText {
				t.Fatalf("LCSClientIdentity() via positional text = %q, want %q", decoded2.Text, tc.wantText)
			}

			// Verify ParseLCSClientIdentityFields also round-trips
			parsed, err := ie.ParseLCSClientIdentityFields(payload)
			if err != nil {
				t.Fatal(err)
			}
			if parsed.Text != tc.wantText {
				t.Fatalf("ParseLCSClientIdentityFields text = %q, want %q", parsed.Text, tc.wantText)
			}
		})
	}
}

func TestLCSClientIdentityFields_Invalid(t *testing.T) {
	t.Parallel()

	t.Run("marshal empty text", func(t *testing.T) {
		fields := &ie.LCSClientIdentityFields{Text: "", Encoding: ie.NetworkNameEncodingGSM7}
		if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
			t.Fatalf("Marshal error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("marshal text too long", func(t *testing.T) {
		fields := &ie.LCSClientIdentityFields{Text: string(make([]byte, 300)), Encoding: ie.NetworkNameEncodingUCS2}
		if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidLength) {
			t.Fatalf("Marshal error = %v, want ErrInvalidLength", err)
		}
	})

	t.Run("parse missing extension bit", func(t *testing.T) {
		if _, err := ie.ParseLCSClientIdentityFields([]byte{0x00}); !errors.Is(err, ie.ErrInvalidValue) {
			t.Fatalf("Parse error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("parse invalid length zero", func(t *testing.T) {
		if _, err := ie.ParseLCSClientIdentityFields(nil); !errors.Is(err, ie.ErrInvalidLength) {
			t.Fatalf("Parse error = %v, want ErrInvalidLength", err)
		}
	})

	t.Run("parse reserved coding scheme", func(t *testing.T) {
		if _, err := ie.ParseLCSClientIdentityFields([]byte{0xA0}); !errors.Is(err, ie.ErrInvalidValue) {
			t.Fatalf("Parse error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("validate nil fields", func(t *testing.T) {
		var fields *ie.LCSClientIdentityFields
		if err := fields.Validate(); !errors.Is(err, ie.ErrInvalidValue) {
			t.Fatalf("Validate error = %v, want ErrInvalidValue", err)
		}
	})

	t.Run("parse too long", func(t *testing.T) {
		buf := make([]byte, 256)
		buf[0] = 0x80 // extension bit set, GSM7
		if _, err := ie.ParseLCSClientIdentityFields(buf); !errors.Is(err, ie.ErrInvalidLength) {
			t.Fatalf("Parse error = %v, want ErrInvalidLength", err)
		}
	})

	t.Run("unmarshal leaves unchanged on error", func(t *testing.T) {
		fields := &ie.LCSClientIdentityFields{Text: "unchanged", Encoding: ie.NetworkNameEncodingUCS2}
		err := fields.UnmarshalBinary([]byte{0x00}) // bad header
		if err == nil {
			t.Fatal("expected error")
		}
		if fields.Text != "unchanged" {
			t.Fatalf("fields modified after error: Text=%q", fields.Text)
		}
	})
}

func TestLCSClientIdentityMarshalLen(t *testing.T) {
	t.Parallel()
	f1 := &ie.LCSClientIdentityFields{Text: "test", Encoding: ie.NetworkNameEncodingGSM7}
	got := f1.MarshalLen()
	if got < 2 {
		t.Fatalf("MarshalLen = %d, want at least 2", got)
	}
	var f2 *ie.LCSClientIdentityFields
	if got := f2.MarshalLen(); got != 0 {
		t.Fatalf("MarshalLen nil = %d, want 0", got)
	}
}

func TestLCSClientIdentityMarshalTo(t *testing.T) {
	t.Parallel()
	fields := &ie.LCSClientIdentityFields{Text: "123", Encoding: ie.NetworkNameEncodingGSM7}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 16)
	err = fields.MarshalTo(buf)
	if err != nil {
		t.Fatal(err)
	}
	if len(buf) != 16 {
		t.Fatalf("buffer length changed from 16 to %d", len(buf))
	}
	if !bytes.Equal(buf[:len(payload)], payload) {
		t.Fatalf("MarshalTo copy = %x, want %x", buf[:len(payload)], payload)
	}

	// Short buffer
	short := make([]byte, len(payload)-1) // one byte too short
	err = fields.MarshalTo(short)
	if err == nil {
		t.Fatal("expected error for short buffer")
	}
	if !errors.Is(err, io.ErrShortBuffer) {
		t.Fatalf("MarshalTo error = %v, want ErrShortBuffer", err)
	}
}

func TestLCSClientIdentityAccessor(t *testing.T) {
	t.Parallel()
	fields := &ie.LCSClientIdentityFields{Text: "SUBSCRIBER", Encoding: ie.NetworkNameEncodingGSM7}
	ie, err := ie.NewLCSClientIdentityWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	result, err := ie.LCSClientIdentity()
	if err != nil {
		t.Fatal(err)
	}
	if result.Text != "SUBSCRIBER" {
		t.Fatalf("LCSClientIdentity text = %q, want %q", result.Text, "SUBSCRIBER")
	}
}

func TestLCSClientIdentityAccessorWrongIE(t *testing.T) {
	t.Parallel()
	wrongIE := ie.NewSSCode(0x42)
	_, err := wrongIE.LCSClientIdentity()
	if err == nil {
		t.Fatal("expected error for wrong IE kind")
	}
	if !errors.Is(err, ie.ErrUnexpectedIE) {
		t.Fatalf("error = %v, want ErrUnexpectedIE", err)
	}
}
