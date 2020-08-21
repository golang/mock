package gomock_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mock_sample "github.com/golang/mock/sample/mock_user"
)

func ExampleCall_DoAndReturn_latency() {
	t := &testing.T{} // provided by test
	ctrl := gomock.NewController(t)
	mockIndex := mock_sample.NewMockIndex(ctrl)

	mockIndex.EXPECT().Get(gomock.Any()).DoAndReturn(
		// signature of anonymous function must have the same number of input and output arguments as the mocked method.
		func(arg string) string {
			time.Sleep(1 * time.Millisecond)
			return "I'm sleepy"
		},
	)

	r := mockIndex.Get("foo")
	fmt.Println(r)
	// Output: I'm sleepy
}

func ExampleCall_DoAndReturn_captureArguments() {
	t := &testing.T{} // provided by test
	ctrl := gomock.NewController(t)
	mockIndex := mock_sample.NewMockIndex(ctrl)
	var s string

	mockIndex.EXPECT().Get(gomock.AssignableToTypeOf(s)).DoAndReturn(
		// signature of anonymous function must have the same number of input and output arguments as the mocked method.
		func(arg string) interface{} {
			s = arg
			return "I'm sleepy"
		},
	)

	r := mockIndex.Get("foo")
	fmt.Printf("%s %s", r, s)
	// Output: I'm sleepy foo
}

func ExampleCall_Do_latency() {
	t := &testing.T{} // provided by test
	ctrl := gomock.NewController(t)
	mockIndex := mock_sample.NewMockIndex(ctrl)

	mockIndex.EXPECT().Anon(gomock.Any()).Do(
		// signature of anonymous function must have the same number of input and output arguments as the mocked method.
		func(_ string) {
			fmt.Println("sleeping")
			time.Sleep(1 * time.Millisecond)
		},
	)

	mockIndex.Anon("foo")
	// Output: sleeping
}

func ExampleCall_Do_captureArguments() {
	t := &testing.T{} // provided by test
	ctrl := gomock.NewController(t)
	mockIndex := mock_sample.NewMockIndex(ctrl)

	var s string
	mockIndex.EXPECT().Anon(gomock.AssignableToTypeOf(s)).Do(
		// signature of anonymous function must have the same number of input and output arguments as the mocked method.
		func(arg string) {
			s = arg
		},
	)

	mockIndex.Anon("foo")
	fmt.Println(s)
	// Output: foo
}
