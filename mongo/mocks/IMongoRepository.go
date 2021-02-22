// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	errors "github.com/skiprco/go-utils/v2/errors"
	mock "github.com/stretchr/testify/mock"

	mongo "github.com/skiprco/go-utils/v2/mongo"
)

// IMongoRepository is an autogenerated mock type for the IMongoRepository type
type IMongoRepository struct {
	mock.Mock
}

// Count provides a mock function with given fields: ctx, collectionName, query, methodName
func (_m *IMongoRepository) Count(ctx context.Context, collectionName string, query map[string]interface{}, methodName string) (int64, *errors.GenericError) {
	ret := _m.Called(ctx, collectionName, query, methodName)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}, string) int64); ok {
		r0 = rf(ctx, collectionName, query, methodName)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 *errors.GenericError
	if rf, ok := ret.Get(1).(func(context.Context, string, map[string]interface{}, string) *errors.GenericError); ok {
		r1 = rf(ctx, collectionName, query, methodName)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*errors.GenericError)
		}
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, collectionName, entityId, methodName
func (_m *IMongoRepository) Delete(ctx context.Context, collectionName string, entityId string, methodName string) *errors.GenericError {
	ret := _m.Called(ctx, collectionName, entityId, methodName)

	var r0 *errors.GenericError
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *errors.GenericError); ok {
		r0 = rf(ctx, collectionName, entityId, methodName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*errors.GenericError)
		}
	}

	return r0
}

// GetMultiple provides a mock function with given fields: ctx, collectionName, query, responses, methodName, opts
func (_m *IMongoRepository) GetMultiple(ctx context.Context, collectionName string, query map[string]interface{}, responses interface{}, methodName string, opts ...*mongo.GetMultipleOption) *errors.GenericError {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, collectionName, query, responses, methodName)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *errors.GenericError
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}, interface{}, string, ...*mongo.GetMultipleOption) *errors.GenericError); ok {
		r0 = rf(ctx, collectionName, query, responses, methodName, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*errors.GenericError)
		}
	}

	return r0
}

// GetOne provides a mock function with given fields: ctx, collectionName, query, acceptsEmptyResult, response, methodName
func (_m *IMongoRepository) GetOne(ctx context.Context, collectionName string, query map[string]interface{}, acceptsEmptyResult bool, response interface{}, methodName string) *errors.GenericError {
	ret := _m.Called(ctx, collectionName, query, acceptsEmptyResult, response, methodName)

	var r0 *errors.GenericError
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}, bool, interface{}, string) *errors.GenericError); ok {
		r0 = rf(ctx, collectionName, query, acceptsEmptyResult, response, methodName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*errors.GenericError)
		}
	}

	return r0
}

// Save provides a mock function with given fields: ctx, collectionName, entity, entityId, methodName
func (_m *IMongoRepository) Save(ctx context.Context, collectionName string, entity interface{}, entityId interface{}, methodName string) *errors.GenericError {
	ret := _m.Called(ctx, collectionName, entity, entityId, methodName)

	var r0 *errors.GenericError
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, interface{}, string) *errors.GenericError); ok {
		r0 = rf(ctx, collectionName, entity, entityId, methodName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*errors.GenericError)
		}
	}

	return r0
}
