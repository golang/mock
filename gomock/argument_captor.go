package gomock

// ArgumentCaptor is a struct that composes a Matcher, and extends it by storing given arguments in the values slice.
type ArgumentCaptor struct {
	Matcher
	values []interface{}
}

// Matches method overrides the Matcher.Matches method
// First it appends any argument(s) used to the values slice.
// Then the parent Matches method is called.
func (ac *ArgumentCaptor) Matches(x interface{}) bool {
	ac.values = append(ac.values, x)
	return ac.Matcher.Matches(x)
}

// LastValue returns the last argument (or arguments) the matcher was called with as an interface{}.
// If the matcher was never called, nil is returned.
func (ac *ArgumentCaptor) LastValue() interface{} {
	if len(ac.values) < 1 {
		return nil
	}
	return ac.values[len(ac.values)-1]
}

// Values returns the all arguments the matcher was called with as a slice of interface{}.
// The values are ordered from first called to last called.
func (ac *ArgumentCaptor) Values() []interface{} {
	return ac.values
}

// Captor is a helper method that returns a new *ArgumentCaptor struct with Matcher set to the given matcher m
func Captor(m Matcher) *ArgumentCaptor {
	return &ArgumentCaptor{Matcher: m}
}

// AnyCaptor is a helper method that returns a new *ArgumentCaptor struct with the matcher set to an anyMatcher
func AnyCaptor() *ArgumentCaptor {
	return &ArgumentCaptor{Matcher: Any()}
}
