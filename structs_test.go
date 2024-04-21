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

	type A struct {
		N int
	}

	type B struct {
		S string
	}

	type C struct {
		N int `json:"n"`
	}

	type D struct {
		S string `json:"s"`
	}

	cases := []struct {
		name string
		vs   []any
		want any
	}{
		{"zero", vs(), of()},
		{"single", vs(A{N: 100}), of("N", 100)},
		{"two", vs(A{N: 100}, B{S: "test"}), of("N", 100, "S", "test")},
		{"tag", vs(C{N: 100}, D{S: "test"}), of("N", 100, `json:"n"`, "S", "test", `json:"s"`)},
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
