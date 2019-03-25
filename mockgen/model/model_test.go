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
		{"/vendor/foo/bar", "foo/bar"},
		{"qux/vendor/foo/bar", "foo/bar"},
		{"qux/vendor/foo/vendor/bar", "bar"},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input %s", tc.input), func(t *testing.T) {
			if got := impPath(tc.input); got != tc.want {
				t.Errorf("got %s; want %s", got, tc.want)
			}
		})
	}

}
