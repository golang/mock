package gomock_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func iptr(x int) *int {
	return &x
}

func sptr(x string) *string {
	return &x
}

func TestImplementsInterface(t *testing.T) {
	type yes interface{}
	type no interface{}
	type testCase struct {
		matcher gomock.Matcher
		yes     []yes
		no      []no
	}
	gomock.RegisterType("io.Writer", reflect.ValueOf(new(io.Writer)).Type().Elem())
	tests := []testCase{
		testCase{
			gomock.ImplementsInterface("error"),
			[]yes{errors.New("foo"), fmt.Errorf("bar: %d", 4)},
			[]no{"xyz", nil, int64(46), iptr(6)},
		},
		testCase{
			gomock.ImplementsInterface(reflect.ValueOf(new(error)).Type().Elem()),
			[]yes{errors.New("foo"), fmt.Errorf("bar: %d", 4)},
			[]no{"xyz", nil, int64(46), iptr(6)},
		},
		testCase{
			gomock.ImplementsInterface("io.Reader"),
			[]yes{strings.NewReader("foo"), bytes.NewReader([]byte("bytes"))},
			[]no{"xyz", nil, int64(46), iptr(6)},
		},
		testCase{
			gomock.ImplementsInterface(reflect.ValueOf(new(io.Reader)).Type().Elem()),
			[]yes{strings.NewReader("foo"), bytes.NewReader([]byte("bytes"))},
			[]no{"xyz", nil, int64(46), iptr(6)},
		},
		testCase{
			gomock.ImplementsInterface("io.Writer"),
			[]yes{os.Stdout, os.Stderr},
			[]no{4, "abc", strings.NewReader("foo")},
		},
		testCase{
			gomock.ImplementsInterface(reflect.ValueOf(new(io.Writer)).Type().Elem()),
			[]yes{os.Stdout, os.Stderr},
			[]no{4, "abc", strings.NewReader("foo")},
		},
		testCase{
			gomock.ImplementsInterface("io.ReadWriter"),
			[]yes{os.Stdout, os.Stderr},
			[]no{4, "abc", strings.NewReader("foo")},
		},
		testCase{
			gomock.ImplementsInterface(reflect.ValueOf(new(io.ReadWriter)).Type().Elem()),
			[]yes{os.Stdout, os.Stderr},
			[]no{4, "abc", strings.NewReader("foo")},
		},
	}
	for i, test := range tests {
		for _, x := range test.yes {
			if !test.matcher.Matches(x) {
				t.Errorf(`test %d: "%v %s" should be true.`, i, x, test.matcher)
			}
		}
		for _, x := range test.no {
			if test.matcher.Matches(x) {
				t.Errorf(`test %d: "%v %s" should be false.`, i, x, test.matcher)
			}
		}
	}
}
