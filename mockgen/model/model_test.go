package model

import (
	"fmt"
	"testing"
)

func TestImpPath(t *testing.T) {
	nonVendor := "github.com/foo/bar"
	if nonVendor != impPath(nonVendor) {
		t.Errorf("")

	}
	testCases := []struct {
		input string
		want  string
	}{
		{"foo/bar", "foo/bar"},
		{"vendor/foo/bar", "foo/bar"},
		{"vendor/foo/vendor/bar", "bar"},
		{"/vendor/foo/bar", "foo/bar"},
		{"qux/vendor/foo/bar", "foo/bar"},
		{"qux/vendor/foo/vendor/bar", "bar"},
		{"govendor/foo", "govendor/foo"},
		{"foo/govendor/bar", "foo/govendor/bar"},
		{"vendors/foo", "vendors/foo"},
		{"foo/vendors/bar", "foo/vendors/bar"},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %s", tc.input), func(t *testing.T) {
			if got := impPath(tc.input); got != tc.want {
				t.Errorf("got %s; want %s", got, tc.want)
			}
		})
	}
}

func TestTypeString(t *testing.T) {
	pkgMap := map[string]string{"pkg": "pkg2"}
	testCases := []struct {
		input Type
		want  string
	}{
		{&EmptyLength{}, "[]"},
		{&LiteralLength{Len: 1}, "[1]"},
		{&ConstLength{Name: "C"}, "[C]"},
		{&ConstLength{Package: "pkg", Name: "C"}, "[pkg2.C]"},
		{&ArrayType{
			Len:  &ConstLength{Package: "pkg", Name: "C"},
			Type: &NamedType{Package: "pkg2", Type: "T"},
		}, "[pkg2.C]T"},
	}
	for _, tc := range testCases {
		if got := tc.input.String(pkgMap, ""); got != tc.want {
			t.Errorf("got %s; want %s", got, tc.want)
		}
	}
}

func TestTypeAddImports(t *testing.T) {
	testCases := []struct {
		input Type
		want  map[string]bool
	}{
		{&EmptyLength{}, map[string]bool{}},
		{&LiteralLength{Len: 1}, map[string]bool{}},
		{&ConstLength{Name: "C"}, map[string]bool{}},
		{&ConstLength{Package: "pkg", Name: "C"}, map[string]bool{"pkg": true}},
		{&ArrayType{
			Len:  &ConstLength{Package: "pkg", Name: "C"},
			Type: &NamedType{Package: "pkg2", Type: "T"},
		}, map[string]bool{"pkg": true, "pkg2": true}},
	}
	for _, tc := range testCases {
		m := map[string]bool{}
		tc.input.addImports(m)
		assertMapEqual(t, m, tc.want)
	}
}

func assertMapEqual(t *testing.T, actual, expected map[string]bool) {
	if len(actual) != len(expected) {
		t.Errorf("expected %v to have the same length as %v", expected, actual)
	}
	for k, v := range actual {
		if actual[k] != v {
			t.Errorf("expected %v to have key %s with value %t", expected, k, v)
		}
	}
}
