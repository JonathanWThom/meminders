// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	twilio "github.com/kevinburke/twilio-go"
	mock "github.com/stretchr/testify/mock"

	url "net/url"
)

// Caller is an autogenerated mock type for the Caller type
type Caller struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *Caller) Create(_a0 context.Context, _a1 url.Values) (*twilio.Call, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *twilio.Call
	if rf, ok := ret.Get(0).(func(context.Context, url.Values) *twilio.Call); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*twilio.Call)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, url.Values) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
