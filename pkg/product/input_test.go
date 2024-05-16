package product

import (
	"errors"
	"testing"
)

func TestInput_Validate(t *testing.T) {
	tests := map[string]struct {
		in   *Input
		want error
	}{
		"it must err for a nil input": {
			in:   nil,
			want: ErrNilInput,
		},
		"it must err for an empty name": {
			in: &Input{
				Name:        "",
				Description: "lorem ipsum",
			},
			want: ErrInvalidName,
		},
		"it must err for an empty description": {
			in: &Input{
				Name:        "lorem ipsum",
				Description: "",
			},
			want: ErrInvalidDescription,
		},
		"it must not err for a valid input": {
			in: &Input{
				Name:        "lorem ipsum",
				Description: "dolor sit amet",
			},
			want: nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			got := test.in.Validate()

			if !errors.Is(got, test.want) {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}
