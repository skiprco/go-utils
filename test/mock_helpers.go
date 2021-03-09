package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
)

// FindMockCallByMethod searches a call inside a mock (defined by ".On(...)")
// and returns a pointer to it. If possible, please use "ReplaceMockReturnArgs"
// or "LimitMockCallCount" instead. Since calling ".On(...)" on the returned
// call won't likely have the expected result.
func FindMockCallByMethod(t *testing.T, mockObject *mock.Mock, methodName string) *mock.Call {
	for _, methodMock := range mockObject.ExpectedCalls {
		if methodMock.Method == methodName {
			// Call found
			return methodMock
		}
	}

	// Call not found
	msg := fmt.Sprintf(`Call with method "%s" not found in mock`, methodName)
	panic(msg)
}

// ReplaceMockReturnArgs replaces the returned arguments for an existing mock.
// This is useful for overwriting a happy mocked call.
//
// Example
//
// usecase, _, _, _, _, userMock, _ := testGetHappyUsecaseAndMocks()
//
// errHas := errors.NewGenericError(418, "test-domain", "test-subdomain", "during-has-operator-read", nil)
//
// test.ReplaceMockReturnArgs(t, &userMock.Mock, "HasOperatorRead", false, errHas)
func ReplaceMockReturnArgs(t *testing.T, mockObject *mock.Mock, methodName string, newArgs ...interface{}) {
	// Get method mock
	methodMock := FindMockCallByMethod(t, mockObject, methodName)

	// Assert arguments count
	if len(methodMock.ReturnArguments) != len(newArgs) {
		msg := fmt.Sprintf(`Unable to replace mock return args: Args count mismatch (%d vs %d)`, len(methodMock.ReturnArguments), len(newArgs))
		panic(msg)
	}

	// Replace arguments
	methodMock.ReturnArguments = newArgs
}

// LimitMockCallCount limits the number of calls the mock will respond to (See Once(), Times(), ...)
func LimitMockCallCount(t *testing.T, mockObject *mock.Mock, methodName string, callCount int) {
	methodMock := FindMockCallByMethod(t, mockObject, methodName)
	methodMock.Repeatability = callCount
}
