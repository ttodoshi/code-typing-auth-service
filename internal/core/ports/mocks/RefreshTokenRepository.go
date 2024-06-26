// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	domain "github.com/ttodoshi/code-typing-auth-service/internal/core/domain"

	mock "github.com/stretchr/testify/mock"
)

// RefreshTokenRepository is an autogenerated mock type for the RefreshTokenRepository type
type RefreshTokenRepository struct {
	mock.Mock
}

// CreateRefreshToken provides a mock function with given fields: refreshToken
func (_m *RefreshTokenRepository) CreateRefreshToken(refreshToken domain.RefreshToken) (string, error) {
	ret := _m.Called(refreshToken)

	if len(ret) == 0 {
		panic("no return value specified for CreateRefreshToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.RefreshToken) (string, error)); ok {
		return rf(refreshToken)
	}
	if rf, ok := ret.Get(0).(func(domain.RefreshToken) string); ok {
		r0 = rf(refreshToken)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(domain.RefreshToken) error); ok {
		r1 = rf(refreshToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteRefreshToken provides a mock function with given fields: refreshToken
func (_m *RefreshTokenRepository) DeleteRefreshToken(refreshToken string) error {
	ret := _m.Called(refreshToken)

	if len(ret) == 0 {
		panic("no return value specified for DeleteRefreshToken")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(refreshToken)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetRefreshToken provides a mock function with given fields: refreshToken
func (_m *RefreshTokenRepository) GetRefreshToken(refreshToken string) (domain.RefreshToken, error) {
	ret := _m.Called(refreshToken)

	if len(ret) == 0 {
		panic("no return value specified for GetRefreshToken")
	}

	var r0 domain.RefreshToken
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (domain.RefreshToken, error)); ok {
		return rf(refreshToken)
	}
	if rf, ok := ret.Get(0).(func(string) domain.RefreshToken); ok {
		r0 = rf(refreshToken)
	} else {
		r0 = ret.Get(0).(domain.RefreshToken)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(refreshToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateRefreshToken provides a mock function with given fields: oldRefreshToken, newRefreshToken
func (_m *RefreshTokenRepository) UpdateRefreshToken(oldRefreshToken string, newRefreshToken string) (domain.RefreshToken, error) {
	ret := _m.Called(oldRefreshToken, newRefreshToken)

	if len(ret) == 0 {
		panic("no return value specified for UpdateRefreshToken")
	}

	var r0 domain.RefreshToken
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (domain.RefreshToken, error)); ok {
		return rf(oldRefreshToken, newRefreshToken)
	}
	if rf, ok := ret.Get(0).(func(string, string) domain.RefreshToken); ok {
		r0 = rf(oldRefreshToken, newRefreshToken)
	} else {
		r0 = ret.Get(0).(domain.RefreshToken)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(oldRefreshToken, newRefreshToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRefreshTokenRepository creates a new instance of RefreshTokenRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRefreshTokenRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *RefreshTokenRepository {
	mock := &RefreshTokenRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
