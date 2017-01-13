package user

import (
	"bufio"
	"bytes"
	"code.google.com/p/gomock/sample"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/golang/mock/sample/imp1"
	mock_user "github.com/golang/mock/sample/mock_user"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestIndexSuite struct {
	*gomock.Controller
	*mock_user.MockIndex
}

func SetUpTest(t *testing.T) *TestIndexSuite {
	suite := &TestIndexSuite{}
	suite.Controller = gomock.NewController(t)
	suite.MockIndex = mock_user.NewMockIndex(suite.Controller)
	return suite
}

func (suite *TestIndexSuite) TearDownTest(t *testing.T) {
	suite.Finish()
}

func TestInvoke1(t *testing.T) {

	suite := SetUpTest(t)
	defer suite.TearDownTest(t)
	underlying := map[string]interface{}{}
	putFunc := func(key string, value interface{}) {
		underlying[key] = value
	}

	getFunc := func(key string) interface{} {
		if value, ok := underlying[key]; ok {
			return value
		} else {
			return nil
		}
	}

	getTwoFunc := func(key1, key2 string) (interface{}, interface{}) {
		return getFunc(key1), getFunc(key2)
	}

	suite.EXPECT().Put(gomock.Any(), gomock.Any()).Invoke(putFunc).AnyTimes()
	suite.EXPECT().Get(gomock.Any()).Invoke(getFunc).AnyTimes()
	suite.EXPECT().GetTwo(gomock.Any(), gomock.Any()).Invoke(getTwoFunc).AnyTimes()

	var value, value2 interface{}
	value = suite.Get("one")
	assert.Equal(t, value, nil)
	suite.Put("one", 1)
	value = suite.Get("one")
	assert.NotEqual(t, value, nil)
	assert.Equal(t, value.(int), 1)

	value, value2 = suite.GetTwo("one", "two")
	assert.Equal(t, value.(int), 1)
	assert.Equal(t, value2, nil)

	suite.Put("two", 2)
	value, value2 = suite.GetTwo("one", "two")
	assert.Equal(t, value.(int), 1)
	assert.Equal(t, value2, 2)

}

func TestInvokeNillableRet(t *testing.T) {
	suite := SetUpTest(t)
	defer suite.TearDownTest(t)
	suite.EXPECT().NillableRet().Invoke(
		func() error { return nil },
	)
	suite.EXPECT().NillableRet().Invoke(
		func() error { return fmt.Errorf("Error") },
	)
	var err error
	err = suite.NillableRet()
	assert.Equal(t, err, nil)
	err = suite.NillableRet()
	assert.NotEqual(t, err, nil)
}

func TestConcreteRet(t *testing.T) {
	suite := SetUpTest(t)
	defer suite.TearDownTest(t)
	suite.EXPECT().ConcreteRet().Invoke(
		func() chan<- bool { return nil },
	)
	suite.EXPECT().ConcreteRet().Invoke(
		func() chan<- bool { return make(chan<- bool, 1) },
	)
	ch := suite.ConcreteRet()
	//t.Errorf("ch=%+v\n",ch)
	assert.Equal(t, true, ch == nil)
	ch = suite.ConcreteRet()
	assert.NotEqual(t, ch, nil)
}

func TestEllip(t *testing.T) {
	suite := SetUpTest(t)
	defer suite.TearDownTest(t)
	var result string
	processor := func(s string, args ...interface{}) {
		sum := 0
		for i := 0; i < len(args); i++ {
			sum += args[i].(int)
		}
		result = fmt.Sprintf(s, sum)
	}
	processor2 := func(args ...string) {
		if len(args) == 0 {
			result = "none"
		} else {
			result = ""
			for i := 0; i < len(args); i++ {
				result += args[i]
			}
		}
	}
	suite.EXPECT().Ellip("%d", 0, 1, 1, 2, 3).Invoke(processor)
	tri := []interface{}{1, 3, 6, 10, 15}
	suite.EXPECT().Ellip("%d", tri...).Invoke(processor)

	suite.EXPECT().EllipOnly("%d", "5", "6", "7").Invoke(processor2)
	suite.EXPECT().EllipOnly().Invoke(processor2)

	suite.Ellip("%d", 0, 1, 1, 2, 3)
	assert.Equal(t, result, "7")
	suite.Ellip("%d", 1, 3, 6, 10, 15)
	assert.Equal(t, result, "35")

	suite.EllipOnly("%d", "5", "6", "7")
	assert.Equal(t, result, "%d567")
	suite.EllipOnly()
	assert.Equal(t, result, "none")
}

func TestRememberInvoke(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIndex := mock_user.NewMockIndex(ctrl)

	kv := map[string]interface{}{}
	putFunc := func(key string, value interface{}) {
		kv[key] = value
	}

	//getFunc := func(key string)interface{} {
	//	if value, ok := kv[key];ok {
	//		return value
	//	}else {
	//		return nil
	//	}
	//}

	//getTwoFunc := func(key1, key2 string)(interface{},interface{}) {
	//	return getFunc(key1), getFunc(key2)
	//}
	mockIndex.EXPECT().Put("a", 1).Invoke(putFunc)            // literals work
	mockIndex.EXPECT().Put("b", gomock.Eq(2)).Invoke(putFunc) // matchers work too

	// NillableRet returns error. Not declaring it should result in a nil return.
	mockIndex.EXPECT().NillableRet().Invoke(func() error { return nil })
	// Calls that returns something assignable to the return type.
	boolc := make(chan bool)
	// In this case, "chan bool" is assignable to "chan<- bool".
	mockIndex.EXPECT().ConcreteRet().Invoke(func() chan<- bool { return boolc })
	// In this case, nil is assignable to "chan<- bool".
	mockIndex.EXPECT().ConcreteRet().Invoke(func() chan<- bool { return nil })

	// Should be able to place expectations on variadic methods.
	mockIndex.EXPECT().Ellip("%d", 0, 1, 1, 2, 3).Invoke(
		func(fmt string, args ...interface{}) {
		},
	) // direct args
	tri := []interface{}{1, 3, 6, 10, 15}
	mockIndex.EXPECT().Ellip("%d", tri...).Invoke(
		func(fmt string, args ...interface{}) {
		},
	) // args from slice
	mockIndex.EXPECT().EllipOnly(gomock.Eq("arg")).Invoke(
		func(...string) {},
	)

	user.Remember(mockIndex, []string{"a", "b"}, []interface{}{1, 2})
	// Check the ConcreteRet calls.
	if c := mockIndex.ConcreteRet(); c != boolc {
		t.Errorf("ConcreteRet: got %v, want %v", c, boolc)
	}
	if c := mockIndex.ConcreteRet(); c != nil {
		t.Errorf("ConcreteRet: got %v, want nil", c)
	}

	// Try one with an action.
	calledString := ""
	mockIndex.EXPECT().Put(gomock.Any(), gomock.Any()).Invoke(func(key string, _ interface{}) {
		calledString = key
	})
	mockIndex.EXPECT().NillableRet().Invoke(func() error { return nil })
	user.Remember(mockIndex, []string{"blah"}, []interface{}{7})
	if calledString != "blah" {
		t.Fatalf(`Uh oh. %q != "blah"`, calledString)
	}

	// Use Do with a nil arg.
	mockIndex.EXPECT().Put("nil-key", gomock.Any()).Invoke(func(key string, value interface{}) {
		if value != nil {
			t.Errorf("Put did not pass through nil; got %v", value)
		}
	})
	mockIndex.EXPECT().NillableRet().Invoke(func() error { return nil })
	user.Remember(mockIndex, []string{"nil-key"}, []interface{}{nil})
}

func TestGrabPointerInvoke(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIndex := mock_user.NewMockIndex(ctrl)
	mockIndex.EXPECT().Ptr(gomock.Any()).Invoke(func(arg *int) { *arg = 7 })

	i := user.GrabPointer(mockIndex)
	if i != 7 {
		t.Errorf("Expected 7, got %d", i)
	}
}
func TestEmbeddedInterface(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmbed := mock_user.NewMockEmbed(ctrl)
	var msg string
	mockEmbed.EXPECT().RegularMethod().Invoke(func() { msg = "RegularMethod" })
	mockEmbed.EXPECT().EmbeddedMethod().Invoke(func() { msg = "EmbededMethod" })
	mockEmbed.EXPECT().ForeignEmbeddedMethod().Invoke(
		func() *bufio.Reader {
			buf := bytes.NewBuffer([]byte{})
			writer := bufio.NewWriter(buf)
			writer.WriteString("ForeignEmbeddedMethod\n")
			writer.Flush()
			return bufio.NewReader(buf)
		},
	)

	mockEmbed.RegularMethod()
	assert.Equal(t, msg, "RegularMethod")
	mockEmbed.EmbeddedMethod()
	assert.Equal(t, msg, "EmbededMethod")
	var emb imp1.ForeignEmbedded = mockEmbed // also does interface check
	reader := emb.ForeignEmbeddedMethod()
	s, err := reader.ReadString('\n')
	assert.Equal(t, err, nil)
	assert.Equal(t, s, "ForeignEmbeddedMethod\n")
}
func TestExpectTrueNilInvoke(t *testing.T) {
	// Make sure that passing "nil" to EXPECT (thus as a nil interface value),
	// will correctly match a nil concrete type.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIndex := mock_user.NewMockIndex(ctrl)
	mockIndex.EXPECT().Ptr(nil).Invoke(func(*int) {}) // this nil is a nil interface{}
	mockIndex.Ptr(nil)                                // this nil is a nil *int
}
