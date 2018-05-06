package gomock

import (
	"fmt"
	"testing"
)

// Reporter wraps a *testing.T, and provides a more useful failure mode
// when interacting with Controller.
//
// For example, consider:
//   func TestMyThing(t *testing.T) {
//     mockCtrl := gomock.NewController(t)
//     defer mockCtrl.Finish()
//     mockObj := something.NewMockMyInterface(mockCtrl)
//     go func() {
//       mockObj.SomeMethod(4, "blah")
//     }
//   }
//
// It hangs without any indication that it's missing an EXPECT() on `mockObj`.
// Providing the Reporter to the gomock.Controller ctor avoids this, and terminates
// with useful feedback. i.e.
//   func TestMyThing(t *testing.T) {
//     mockCtrl := gomock.NewController(Reporter{t})
//     defer mockCtrl.Finish()
//     mockObj := something.NewMockMyInterface(mockCtrl)
//     go func() {
//       mockObj.SomeMethod(4, "blah") // crashes the test now
//     }
//   }
type Reporter struct {
	T *testing.T
}

// ensure Reporter implements gomock.TestReporter.
var _ TestReporter = Reporter{}

// Errorf is equivalent to testing.T.Errorf.
func (r Reporter) Errorf(format string, args ...interface{}) {
	r.T.Errorf(format, args...)
}

// Fatalf crashes the program with a panic to allow users to diagnose
// missing expects.
func (r Reporter) Fatalf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}
