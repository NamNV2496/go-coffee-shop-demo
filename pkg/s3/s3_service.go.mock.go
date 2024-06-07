// Code generated by MockGen. DO NOT EDIT.
// Source: s3_service.go
//
// Generated by this command:
//
//	mockgen -source=s3_service.go -destination=s3_service.go.mock.go -package=s3
//

// Package s3 is a generated GoMock package.
package s3

import (
	context "context"
	multipart "mime/multipart"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockS3Client is a mock of S3Client interface.
type MockS3Client struct {
	ctrl     *gomock.Controller
	recorder *MockS3ClientMockRecorder
}

// MockS3ClientMockRecorder is the mock recorder for MockS3Client.
type MockS3ClientMockRecorder struct {
	mock *MockS3Client
}

// NewMockS3Client creates a new mock instance.
func NewMockS3Client(ctrl *gomock.Controller) *MockS3Client {
	mock := &MockS3Client{ctrl: ctrl}
	mock.recorder = &MockS3ClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockS3Client) EXPECT() *MockS3ClientMockRecorder {
	return m.recorder
}

// PreviewImage mocks base method.
func (m *MockS3Client) PreviewImage(ctx context.Context, filename, bucketName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PreviewImage", ctx, filename, bucketName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PreviewImage indicates an expected call of PreviewImage.
func (mr *MockS3ClientMockRecorder) PreviewImage(ctx, filename, bucketName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PreviewImage", reflect.TypeOf((*MockS3Client)(nil).PreviewImage), ctx, filename, bucketName)
}

// Write mocks base method.
func (m *MockS3Client) Write(ctx context.Context, fileName, bucket string, img multipart.File, size int64, contentType string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", ctx, fileName, bucket, img, size, contentType)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockS3ClientMockRecorder) Write(ctx, fileName, bucket, img, size, contentType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockS3Client)(nil).Write), ctx, fileName, bucket, img, size, contentType)
}