package gomock_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mock_sample "github.com/golang/mock/sample/mock_user"
)

func ExampleCall_DoAndReturn() {
	t := &testing.T{} // provided by test
	ctrl := gomock.NewController(t)
	mockIndex := mock_sample.NewMockIndex(ctrl)

	mockIndex.EXPECT().Get(gomock.Any()).DoAndReturn(
		func() string {
			time.Sleep(1 * time.Second)
			return "I'm sleepy"
		},
	)
}

func ExampleCall_DoAndReturn_captureArguments() {
	t := &testing.T{} // provided by test
	ctrl := gomock.NewController(t)
	mockIndex := mock_sample.NewMockIndex(ctrl)
	var s string

	mockIndex.EXPECT().Get(gomock.AssignableToTypeOf(s)).DoAndReturn(
		// When capturing arguments the anonymous function should have the same signature as the mocked method.
		func(arg string) interface{} {
			time.Sleep(1 * time.Second)
			fmt.Println(arg)
			return "I'm sleepy"
		},
	)
}

func ExampleCall_Do() {
	t := &testing.T{} // provided by test
	ctrl := gomock.NewController(t)
	mockIndex := mock_sample.NewMockIndex(ctrl)

	mockIndex.EXPECT().Anon(gomock.Any()).Do(
		func() {
			time.Sleep(1 * time.Second)
		},
	)
}

func ExampleCall_Do_captureArguments() {
	t := &testing.T{} // provided by test
	ctrl := gomock.NewController(t)
	mockIndex := mock_sample.NewMockIndex(ctrl)

	var s string
	mockIndex.EXPECT().Anon(gomock.AssignableToTypeOf(s)).Do(
		// When capturing arguments the anonymous function should have the same signature as the mocked method.
		func(arg string) {
			time.Sleep(1 * time.Second)
			fmt.Println(arg)
		},
	)
}
