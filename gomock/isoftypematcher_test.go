package gomock_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
)

type dog struct {
	breed, name string
}

type cat struct {
	breed, name string
}

func iptr(x int) *int {
	return &x
}

func sptr(x string) *string {
	return &x
}

func TestIsOfType(t *testing.T) {
	type yes interface{}
	type no interface{}
	type testCase struct {
		matcher gomock.Matcher
		yes     []yes
		no      []no
	}
	str := "str"
	sptrVar := &str
	sptr2Var := &sptrVar
	luckyDog := dog{breed: "Labrador", name: "Lucky"}
	fluffyCat := cat{breed: "Siamese", name: "Fluffy"}
	tests := []testCase{
		testCase{
			gomock.IsOfType("bool"),
			[]yes{true, false},
			[]no{"xyz", nil, int64(46), iptr(6)},
		},
		testCase{
			gomock.IsOfType("*errorString"),
			[]yes{errors.New("foo"), fmt.Errorf("bar: %d", 4)},
			[]no{"xyz", nil, int64(46), iptr(6)},
		},
		testCase{
			gomock.IsOfType("int"),
			[]yes{3, 4, 5},
			[]no{"xyz", nil, int64(46), iptr(6)},
		},
		testCase{
			gomock.IsOfType("*int"),
			[]yes{iptr(3), iptr(0), iptr(-1)},
			[]no{"xyz", nil, int64(46), 1, 0, -1, nil},
		},
		testCase{
			gomock.IsOfType("int64"),
			[]yes{int64(3), int64(4), int64(5)},
			[]no{"xyz", nil, 46},
		},
		testCase{
			gomock.IsOfType("float64"),
			[]yes{1.2, 2.3, 0.0, -1.0},
			[]no{"xyz", nil, 46, float32(1.2)},
		},
		testCase{
			gomock.IsOfType("float32"),
			[]yes{float32(1.2), float32(2.3), float32(0.0), float32(-1.0)},
			[]no{"xyz", nil, 46, float64(3.4), 1.2},
		},
		testCase{
			gomock.IsOfType("string"),
			[]yes{"abc", "3", ""},
			[]no{nil, int64(46), iptr(6), 7, sptr("def"), sptr},
		},
		testCase{
			gomock.IsOfType("*string"),
			[]yes{sptr("abc"), sptr("3"), sptr(""), sptrVar},
			[]no{nil, int64(46), iptr(6), 7, &sptrVar, sptr2Var, &sptr2Var},
		},
		testCase{
			gomock.IsOfType("**string"),
			[]yes{&sptrVar, sptr2Var},
			[]no{nil, int64(46), iptr(6), 7, sptrVar, &sptr2Var},
		},
		testCase{
			gomock.IsOfType("***string"),
			[]yes{&sptr2Var},
			[]no{nil, int64(46), iptr(6), 7, sptrVar, sptr2Var},
		},
		testCase{
			gomock.IsOfType("chan int"),
			[]yes{make(chan int)},
			[]no{nil, int64(46), iptr(6), 7, sptrVar, sptr2Var, make(chan string), make(chan struct{}), luckyDog},
		},
		testCase{
			gomock.IsOfType("dog"),
			[]yes{luckyDog},
			[]no{nil, &luckyDog, fluffyCat},
		},
		testCase{
			gomock.IsOfType("*dog"),
			[]yes{&luckyDog},
			[]no{nil, luckyDog, &fluffyCat},
		},
		testCase{
			gomock.IsOfType("cat"),
			[]yes{fluffyCat},
			[]no{nil, &luckyDog, luckyDog},
		},
		testCase{
			gomock.IsOfType("*cat"),
			[]yes{&fluffyCat},
			[]no{nil, luckyDog, &luckyDog},
		},
		testCase{
			gomock.IsOfType("nil"),
			[]yes{nil},
			[]no{"xyz", int64(46), iptr(6)},
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
