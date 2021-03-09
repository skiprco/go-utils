package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockTestObject struct {
	mock.Mock
}

func (o *mockTestObject) methodTest(success bool) bool {
	args := o.Called(success)
	return args.Bool(0)
}

func (o *mockTestObject) anotherMethodTest(success bool) bool {
	return o.methodTest(success)
}

func Test_FindMockCallByMethod_Success(t *testing.T) {
	testMock := mockTestObject{}
	testMock.On("methodTest").Return(true)

	call := FindMockCallByMethod(t, &testMock.Mock, "methodTest")
	require.NotNil(t, call)
	assert.Equal(t, "methodTest", call.Method)
}

func Test_FindMockCallByMethod_Failure(t *testing.T) {
	testMock := mockTestObject{}
	testMock.On("methodTest").Return(true)

	assert.Panics(t, func() {
		FindMockCallByMethod(t, &testMock.Mock, "anotherMethodTest")
	})
}

func Test_ReplaceMockReturnArgs_Success(t *testing.T) {
	testMock := mockTestObject{}
	testMock.On("methodTest").Return(true)
	ReplaceMockReturnArgs(t, &testMock.Mock, "methodTest", false)
}

func Test_ReplaceMockReturnArgs_MethodNotFound_Failure(t *testing.T) {
	testMock := mockTestObject{}
	testMock.On("methodTest").Return(true)

	assert.Panics(t, func() {
		ReplaceMockReturnArgs(t, &testMock.Mock, "anotherMethodTest", false)
	})
}

func Test_ReplaceMockReturnArgs_ArgCountMismatch_Failure(t *testing.T) {
	testMock := mockTestObject{}
	testMock.On("methodTest").Return(true)

	assert.Panics(t, func() {
		ReplaceMockReturnArgs(t, &testMock.Mock, "methodTest", false, false, false)
	})
}

func Test_LimitMockCallCount_Success(t *testing.T) {
	testMock := mockTestObject{}
	testMock.On("methodTest").Return(true)
	LimitMockCallCount(t, &testMock.Mock, "methodTest", 5)
	call := FindMockCallByMethod(t, &testMock.Mock, "methodTest")
	assert.Equal(t, 5, call.Repeatability)
}

func Test_LimitMockCallCount_MethodNotFound_Failure(t *testing.T) {
	testMock := mockTestObject{}
	testMock.On("methodTest").Return(true)

	assert.Panics(t, func() {
		LimitMockCallCount(t, &testMock.Mock, "anotherMethodTest", 5)
	})
}
