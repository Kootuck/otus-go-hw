package hw09structvalidator

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11|regexp:^[0-9]+$"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"length:5"`
	}

	RCase struct {
		Field string `validate:"regexp:^\\w+("` // typo -> open (
	}
	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func (r UserRole) String() string {
	return string(r)
}

var tests = []struct {
	name        string
	in          interface{}
	expectedErr error
}{
	{
		name: "positive case 1-> all values are correct and satisfy validation conditions",
		in: User{
			ID:     "123456789012345678901234567890123456",
			Name:   "Test User",
			Age:    25,
			Email:  "testuser@example.com",
			Role:   "admin",
			Phones: []string{"12345678901"},
		},
		expectedErr: nil,
	},
	{
		name: "positive case 2 -> all fields violate specified validation criteria",
		in: User{
			ID:     "123", // length not 36
			Name:   "Invalid User",
			Age:    15,                                      // Age less than 18
			Email:  "invalid",                               // has value, but doesn't match a regexp
			Role:   "invalid",                               // has value, but is not in a list of allowed values
			Phones: []string{"12345678901", "123456789012"}, // Length of [1] is > 11
		},
		expectedErr: ValidationErrors{
			ValidationError{Field: "ID", Error: NewStrictStringLengthError(3, 36)},
			ValidationError{Field: "Age", Error: NewIntMustBeLargerThanError(18, 15)},
			ValidationError{Field: "Email", Error: NewStringRegExpError("^\\w+@\\w+\\.\\w+$", "invalid")},
			ValidationError{Field: "Role", Error: NewStringNotAllowedError("invalid", "admin, stuff")},
			ValidationError{Field: "Phones", Error: NewStrictStringLengthError(12, 11)},
		},
	},
	{
		name: "positive case 3 -> all fields violate specified validation criteria by beign empty",
		in: User{ //
			ID:     "", // length not 36
			Name:   "Invalid User",
			Age:    51,         // Age is more than 50
			Email:  "",         // Email is empty -> has to satisfy regexp
			Role:   "",         // Role is empty -> Has to be one of the specificed
			Phones: []string{}, // Phone is empty
		},
		expectedErr: ValidationErrors{
			ValidationError{Field: "ID", Error: NewStrictStringLengthError(0, 36)},
			ValidationError{Field: "Age", Error: NewIntMustBeLowerThanError(50, 51)},
			ValidationError{Field: "Email", Error: NewStringRegExpError("^\\w+@\\w+\\.\\w+$", "")},
			ValidationError{Field: "Role", Error: NewStringNotAllowedError("", "admin, stuff")},
		},
	},
	{
		name: "positive case 4 -> validation is not defined and therefore skipped",
		in: Token{
			Header:    []byte(base64.StdEncoding.EncodeToString([]byte("header information"))),
			Payload:   []byte(base64.StdEncoding.EncodeToString([]byte("payload information"))),
			Signature: []byte(base64.StdEncoding.EncodeToString([]byte("signature information"))),
		},
		expectedErr: nil,
	},
	{
		name: "postivie case 5 -> \"Code\" is of allowed values",
		in: Response{
			Code: 404,
			Body: "",
		},
		expectedErr: nil,
	},
	{
		name: "postivie case 6 -> \"Code\" is not among allowed values",
		in: Response{
			Code: 999,
			Body: "",
		},
		expectedErr: ValidationErrors{
			ValidationError{Field: "Code", Error: NewIntNotAllowedError(9999, []string{"200", "404", "500"})},
		},
	},
	{
		name: "negative case 1 -> incorrect validation tag",
		in: App{
			Version: "1.21.",
		},
		expectedErr: ProgramError{}, // unsupported tag ("length" instead of "len")
	},
	{
		name: "negative case 2 -> incorrect validation tag (regexp)",
		in: RCase{
			Field: "Hello world",
		},
		expectedErr: ProgramError{},
	},
}

func TestValidate(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wantPE ProgramError
			var wantVE ValidationErrors
			tt := tt
			// t.Parallel()

			err := Validate(tt.in)

			// 1. check that positive cases are indeed positive
			if tt.expectedErr == nil && err != nil {
				t.Fatalf("Expected no error, but received %v", err)
				return
			}

			// 2. All good
			if tt.expectedErr == nil && err == nil {
				return
			}

			// 3. check that returned error is of correct type}
			if errors.As(tt.expectedErr, &wantPE) {
				var got ProgramError
				if ok := errors.As(err, &got); !ok {
					t.Fatalf("error type mismatch: expected \"ProgramError\"")
				}
			}

			if errors.As(tt.expectedErr, &wantVE) {
				var got ValidationErrors
				if ok := errors.As(err, &got); !ok {
					t.Fatalf("error type mismatch: expected \"ValidationErrors\", got %T", err)
				}
				// 3. check that received errors are exactly what were expected
				checkGotEqualsWant(t, got, wantVE)
			}
		})
	}
}

func checkGotEqualsWant(t *testing.T, got ValidationErrors, want ValidationErrors) {
	t.Helper()
	matchCount := 0
	for _, g := range got {
		for _, e := range want {
			if errors.Is(g.Error, e.Error) {
				matchCount++
				break
			}
		}
	}
	if len(want) != matchCount {
		t.Fatalf("Not all expected errors were received")
	}
	if len(got) != matchCount {
		t.Fatalf("Not all received errors were expected")
	}
}
