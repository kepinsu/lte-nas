package ie_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"lte-nas/ie"
)

func TestAuthenticationParameterRANDFields(t *testing.T) {
	t.Parallel()
	fields := &ie.AuthenticationParameterRANDFields{RandomChallenge: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}}
	value, err := ie.NewPositionalAuthenticationParameterRANDWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.AuthenticationParameterRAND()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("AuthenticationParameterRAND = %+v, want %+v", decoded, fields)
	}
}

func TestAuthenticationParameterAUTNFields(t *testing.T) {
	t.Parallel()
	fields := &ie.AuthenticationParameterAUTNFields{
		SequenceNumberXORAnonymityKey: [6]byte{1, 2, 3, 4, 5, 6},
		AuthenticationManagementField: [2]byte{7, 8},
		MessageAuthenticationCode:     [8]byte{9, 10, 11, 12, 13, 14, 15, 16},
	}
	payload, err := fields.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}; !bytes.Equal(payload, want) {
		t.Fatalf("Marshal = %x, want %x", payload, want)
	}
}

func TestAuthenticationFailureParameterFields(t *testing.T) {
	t.Parallel()
	fields := &ie.AuthenticationFailureParameterFields{
		SequenceNumberXORAnonymityKey: [6]byte{1, 2, 3, 4, 5, 6},
		MessageAuthenticationCode:     [8]byte{9, 10, 11, 12, 13, 14, 15, 16},
	}
	value, err := ie.NewAuthenticationFailureParameterWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.AuthenticationFailureParameter()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("AuthenticationFailureParameter = %+v, want %+v", decoded, fields)
	}
}

func TestAuthenticationResponseParameterFields(t *testing.T) {
	t.Parallel()
	fields := &ie.AuthenticationResponseParameterFields{Response: []byte{1, 2, 3, 4, 5, 6, 7, 8}}
	value, err := ie.NewPositionalAuthenticationResponseParameterWithFields(fields)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := value.AuthenticationResponseParameter()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, fields) {
		t.Fatalf("AuthenticationResponseParameter = %+v, want %+v", decoded, fields)
	}

	fields.Response[0] = 0xff
	if decoded.Response[0] == 0xff {
		t.Fatal("decoded response aliases input fields")
	}
}

func TestAuthenticationResponseParameterFields_Invalid(t *testing.T) {
	t.Parallel()
	if _, err := (&ie.AuthenticationResponseParameterFields{Response: []byte{1, 2, 3}}).Marshal(); !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("Marshal error = %v, want ErrInvalidLength", err)
	}
}
