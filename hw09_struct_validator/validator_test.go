package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID    string `json:"id" validate:"len:36"`
		Name  string
		Age   int      `validate:"min:18|max:50"`
		Email string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role  UserRole `validate:"in:admin,stuff"`

		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
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

	OnlyMax struct {
		Age int `validate:"max:70"`
	}

	MinMax struct {
		Age int `validate:"min:18|max:60"`
	}
)

func TestValidateProgErr(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          "loko",
			expectedErr: ErrorIsNotStruct,
		},
		{
			in:          1254,
			expectedErr: ErrorIsNotStruct,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			i++
			myErr := Validate(tt.in)
			require.Error(t, myErr)
			require.Truef(t, errors.Is(myErr, tt.expectedErr),
				"expected error is <%v>, but actual error is <%v>", tt.expectedErr, myErr)
		})
	}
}

func TestValidateSucs(t *testing.T) {
	tests := []interface{}{
		App{
			Version: "12.12",
		},
		Response{
			Code: 200,
			Body: "body1",
		},
		MinMax{
			Age: 45,
		},
		OnlyMax{
			Age: 20,
		},

		User{
			ID:    "123456789012345678901234567890abcdef", // "len:36"`
			Name:  "Vasya",
			Age:   30,                // "min:18|max:50"
			Email: "palpalych@bk.ru", // "regexp:^\\w+@\\w+\\.\\w+$"`
			Role:  "admin",           // "in:admin,stuff"`

			Phones: []string{"79031234567", "12345678900"}, // `validate:"len:11"`
			meta:   nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			i++
			require.NoError(t, Validate(tt))
		})
	}
}

func TestValidateValidateErr(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr ValidationErrors
	}{
		{
			in:          App{Version: "v1"},
			expectedErr: ValidationErrors{{Field: "Version", Err: ErrorFIsNotVal}},
		},
		{
			in: Response{
				Code: 100,
				Body: "body1",
			},
			expectedErr: ValidationErrors{{"Code", ErrorFIsNotVal}},
		},
		{
			in:          OnlyMax{Age: 100},
			expectedErr: ValidationErrors{{Field: "Age", Err: ErrorFIsNotVal}},
		},
		{
			in:          MinMax{Age: 11},
			expectedErr: ValidationErrors{{Field: "Age", Err: ErrorFIsNotVal}},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890abcde", // "len:36"`
				Name:   "Vasya",
				Age:    70,                                         // "min:18|max:50"
				Email:  "palpalychbk.ru",                           // "regexp:^\\w+@\\w+\\.\\w+$"`
				Role:   "manager",                                  // "in:admin,stuff"`
				Phones: []string{"79031234567", "123456789003232"}, // `validate:"len:11"`
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrorFIsNotVal},
				{Field: "Age", Err: ErrorFIsNotVal},
				{Field: "Email", Err: ErrorFIsNotVal},
				{Field: "Role", Err: ErrorFIsNotVal},
				{Field: "Phones", Err: ErrorFIsNotVal},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			i++
			myErr := Validate(tt.in)
			require.Error(t, myErr)
			myValErr, res := myErr.(ValidationErrors) //nolint:errorlint
			if res {
				myStr := "wrong len off ValidationErrors, expected %d, but real %d"
				require.Equal(t, len(tt.expectedErr), len(myValErr), myStr, len(tt.expectedErr), len(myValErr))
				for j := 0; j < len(myValErr); j++ {
					myStr = "expected error is <%v>, but actual error is <%v>"
					require.Equal(t, tt.expectedErr[j], myValErr[j], myStr, tt.expectedErr[j], myValErr[j])
				}
			} else {
				t.Error("fatal Error")
			}
		})
	}
}
