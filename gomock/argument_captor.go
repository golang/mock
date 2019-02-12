package gomock

type ArgumentCaptor interface {
	Matches(x interface{}) bool
	Value() interface{}
	AllValues() []interface{}
}

type argumentCaptor struct {
	Matcher
	values []interface{}
}

func (ac *argumentCaptor) Matches(x interface{}) bool {
	ac.values = append(ac.values, x)
	return ac.Matcher.Matches(x)
}

// If called multiple times, return the last argument value used
func (ac *argumentCaptor) Value() interface{} {
	if len(ac.values) < 1 {
		return nil
	}
	return ac.values[len(ac.values)-1]
}

func (ac *argumentCaptor) AllValues() []interface{} {
	return ac.values
}

func Captor(m Matcher) ArgumentCaptor {
	return &argumentCaptor{Matcher: m}
}

func AnyCaptor() ArgumentCaptor {
	return &argumentCaptor{Matcher: Any()}
}
