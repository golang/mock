package gomock

import (
	"reflect"
	"testing"
)

type mockTestReporter struct {
	errorCalls int
	fatalCalls int
}

func (o *mockTestReporter) Errorf(format string, args ...interface{}) {
	o.errorCalls++
}

func (o *mockTestReporter) Fatalf(format string, args ...interface{}) {
	o.fatalCalls++
}

func TestCall_After(t *testing.T) {
	t.Run("SelfPrereqCallsFatalf", func(t *testing.T) {
		tr1 := &mockTestReporter{}

		c := &Call{t: tr1}
		c.After(c)

		if tr1.fatalCalls != 1 {
			t.Errorf("number of fatal calls == %v, want 1", tr1.fatalCalls)
		}
	})

	t.Run("LoopInCallOrderCallsFatalf", func(t *testing.T) {
		tr1 := &mockTestReporter{}
		tr2 := &mockTestReporter{}

		c1 := &Call{t: tr1}
		c2 := &Call{t: tr2}
		c1.After(c2)
		c2.After(c1)

		if tr1.errorCalls != 0 || tr1.fatalCalls != 0 {
			t.Error("unexpected errors")
		}

		if tr2.fatalCalls != 1 {
			t.Errorf("number of fatal calls == %v, want 1", tr2.fatalCalls)
		}
	})
}

func TestCall_SetArg(t *testing.T) {
	t.Run("SetArgSlice", func(t *testing.T) {
		c := &Call{
			methodType: reflect.TypeOf(func([]byte) {}),
		}
		c.SetArg(0, []byte{1, 2, 3})

		in := []byte{4, 5, 6}
		c.call([]interface{}{in})

		if in[0] != 1 || in[1] != 2 || in[2] != 3 {
			t.Error("Expected SetArg() to modify input slice argument")
		}
	})

	t.Run("SetArgPointer", func(t *testing.T) {
		c := &Call{
			methodType: reflect.TypeOf(func(*int) {}),
		}
		c.SetArg(0, 42)

		in := 43
		c.call([]interface{}{&in})
		if in != 42 {
			t.Error("Expected SetArg() to modify value pointed to by argument")
		}
	})
}
