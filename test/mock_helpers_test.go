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

func (o *mockTestObject) test(success bool) bool {
	args := o.Called(success)
	return args.Bool(0)
}

func (o *mockTestObject) test2(success bool) bool {
	return o.test(success)
}

func Test_FindMockCallByMethod_Success(t *testing.T) {
	testMock := mockTestObject{}
	testMock.On("test").Return(true)

	call := FindMockCallByMethod(t, &testMock.Mock, "test")
	require.NotNil(t, call)
	assert.Equal(t, "test", call.Method)
}

func Test_FindMockCallByMethod_Failure(t *testing.T) {
	testMock := mockTestObject{}
	testMock.On("test").Return(true)

	assert.Panics(t, func() {
		FindMockCallByMethod(t, &testMock.Mock, "test2")
	})
}

func Test_ReplaceMockReturnArgs_Success(t *testing.T) {
	testMock := mockTestObject{}
	testMock.On("test").Return(true)
	ReplaceMockReturnArgs(t, &testMock.Mock, "test", false)
}

func Test_ReplaceMockReturnArgs_MethodNotFound_Failure(t *testing.T) {
	testMock := mockTestObject{}
	testMock.On("test").Return(true)

	assert.Panics(t, func() {
		ReplaceMockReturnArgs(t, &testMock.Mock, "test2", false)
	})
}

func Test_ReplaceMockReturnArgs_ArgCountMismatch_Failure(t *testing.T) {
	testMock := mockTestObject{}
	testMock.On("test").Return(true)

	assert.Panics(t, func() {
		ReplaceMockReturnArgs(t, &testMock.Mock, "test", false, false, false)
	})
}

func Test_LimitMockCallCount_Success(t *testing.T) {
	testMock := mockTestObject{}
	testMock.On("test").Return(true)
	LimitMockCallCount(t, &testMock.Mock, "test", 5)
	call := FindMockCallByMethod(t, &testMock.Mock, "test")
	assert.Equal(t, 5, call.Repeatability)
}

func Test_LimitMockCallCount_MethodNotFound_Failure(t *testing.T) {
	testMock := mockTestObject{}
	testMock.On("test").Return(true)

	assert.Panics(t, func() {
		LimitMockCallCount(t, &testMock.Mock, "test2", 5)
	})
}
