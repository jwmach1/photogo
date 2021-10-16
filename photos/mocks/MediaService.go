// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	data "velocitizer.com/photogo/data"
)

// MediaService is an autogenerated mock type for the MediaService type
type MediaService struct {
	mock.Mock
}

// Get provides a mock function with given fields: mediaItem
func (_m *MediaService) Get(mediaItem data.MediaItem) ([]byte, error) {
	ret := _m.Called(mediaItem)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(data.MediaItem) []byte); ok {
		r0 = rf(mediaItem)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(data.MediaItem) error); ok {
		r1 = rf(mediaItem)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: nextPageToken
func (_m *MediaService) List(nextPageToken string) (*data.MediaResponse, error) {
	ret := _m.Called(nextPageToken)

	var r0 *data.MediaResponse
	if rf, ok := ret.Get(0).(func(string) *data.MediaResponse); ok {
		r0 = rf(nextPageToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*data.MediaResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(nextPageToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
