package structs_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tenntenn/structs"
)

func TestMerge(t *testing.T) {
	t.Parallel()

	vs := func(vs ...any) []any { return vs }
	of := structs.Of
	F := structs.F
	Tag := structs.Tag[string]

	type A struct{ N int }
	type B struct{ S string }
	type C struct {
		N int `json:"n"`
	}
	type D struct {
		S string `json:"s"`
	}
	type E struct{ n int }

	cases := []struct {
		name string
		vs   []any
		want any
	}{
		{"zero", vs(), of()},
		{"single", vs(A{N: 100}), of(F("N", 100))},
		{"two", vs(A{N: 100}, B{S: "test"}), of(F("N", 100), F("S", "test"))},
		{"tag", vs(C{N: 100}, D{S: "test"}), of(F("N", 100, Tag(`json:"n"`)), F("S", "test", Tag(`json:"s"`)))},
		{"unexported", vs(E{n: 100}, B{S: "test"}), of(F("S", "test"))},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := structs.Merge(tt.vs...)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Error("(got, want):", diff)
			}
		})
	}
}
