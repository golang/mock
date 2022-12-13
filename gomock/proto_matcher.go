package gomock

import (
	"fmt"

	protov1 "github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/proto"
)

func ProtoEq(p proto.Message) *protoMatcher {
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

func ProtoV1Eq(p protov1.Message) *protoV1Matcher {
	return &protoV1Matcher{
		Message: p,
	}
}

type protoV1Matcher struct {
	protov1.Message
}

func (m protoV1Matcher) Matches(v interface{}) bool {
	vp, ok := v.(protov1.Message)
	if !ok {
		return false
	}
	return protov1.Equal(vp, m.Message)
}

func (m protoV1Matcher) String() string {
	return fmt.Sprintf("is equal to %v (%T)", m.Message, m.Message)
}
