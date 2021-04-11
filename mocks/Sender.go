// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	twilio "github.com/kevinburke/twilio-go"
	mock "github.com/stretchr/testify/mock"

	url "net/url"
)

// Sender is an autogenerated mock type for the Sender type
type Sender struct {
	mock.Mock
}

// SendMessage provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Sender) SendMessage(_a0 string, _a1 string, _a2 string, _a3 []*url.URL) (*twilio.Message, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 *twilio.Message
	if rf, ok := ret.Get(0).(func(string, string, string, []*url.URL) *twilio.Message); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*twilio.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, []*url.URL) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
