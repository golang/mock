package gomock

import (
	"fmt"

	gogo "github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/proto" // TODO: change to "google.golang.org/protobuf/proto"
)

func protoEq(p proto.Message) *protoMatcher {
	return &protoMatcher{
		Message: p,
	}
}

type protoMatcher struct {
	proto.Message
}

func (m protoMatcher) Matches(v interface{}) bool {
	vp, ok := v.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(vp, m.Message)
}

func (m protoMatcher) String() string {
	return fmt.Sprintf("is equal to %v (%T)", m.Message, m.Message)
}

func gogoEq(p proto.Message) *gogoMatcher {
	return &gogoMatcher{
		Message: p,
	}
}

type gogoMatcher struct {
	gogo.Message
}

func (m gogoMatcher) Matches(v interface{}) bool {
	vp, ok := v.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(vp, m.Message)
}

func (m gogoMatcher) String() string {
	return fmt.Sprintf("is equal to %v (%T)", m.Message, m.Message)
}
