// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/handler/handler.go

// Package handler_test is a generated GoMock package.
package handler_test

import (
	context "context"
	entity "dmmak/simple-rest-crud/internal/entity"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockRedditPostsRepo is a mock of RedditPostsRepo interface.
type MockRedditPostsRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRedditPostsRepoMockRecorder
}

// MockRedditPostsRepoMockRecorder is the mock recorder for MockRedditPostsRepo.
type MockRedditPostsRepoMockRecorder struct {
	mock *MockRedditPostsRepo
}

// NewMockRedditPostsRepo creates a new mock instance.
func NewMockRedditPostsRepo(ctrl *gomock.Controller) *MockRedditPostsRepo {
	mock := &MockRedditPostsRepo{ctrl: ctrl}
	mock.recorder = &MockRedditPostsRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedditPostsRepo) EXPECT() *MockRedditPostsRepoMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockRedditPostsRepo) Delete(ctx context.Context, postID string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, postID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockRedditPostsRepoMockRecorder) Delete(ctx, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRedditPostsRepo)(nil).Delete), ctx, postID)
}

// Get mocks base method.
func (m *MockRedditPostsRepo) Get(ctx context.Context, postID string) (*entity.RedditPost, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, postID)
	ret0, _ := ret[0].(*entity.RedditPost)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockRedditPostsRepoMockRecorder) Get(ctx, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRedditPostsRepo)(nil).Get), ctx, postID)
}

// Save mocks base method.
func (m *MockRedditPostsRepo) Save(ctx context.Context, post *entity.RedditPost) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, post)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockRedditPostsRepoMockRecorder) Save(ctx, post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockRedditPostsRepo)(nil).Save), ctx, post)
}