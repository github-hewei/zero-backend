package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-web/errcode"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) List(ctx context.Context, req *ListRequest) (*ListResult, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListResult), args.Error(1)
}

func (m *mockUserService) Create(ctx context.Context, req *CreateRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockUserService) Update(ctx context.Context, req *UpdateRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockUserService) Delete(ctx context.Context, req *DeleteRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockUserService) Detail(ctx context.Context, id uint32) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *mockUserService) GetPointsLogs(ctx context.Context, req *PointsLogListRequest) (*ListResult, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ListResult), args.Error(1)
}

func (m *mockUserService) ChangePoints(ctx context.Context, req *PointsChangeRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func setupGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func newTestBinder() *bind.Binder {
	validate := bind.NewValidate()
	trans := bind.MustNewTrans(validate)
	return bind.New(validate, trans, errcode.InvalidInput)
}

func TestHandler_List(t *testing.T) {
	mockSvc := new(mockUserService)
	mockSvc.On("List", mock.Anything, mock.Anything).
		Return(&ListResult{List: []*User{}, Total: 0}, nil)

	h := newHandler(newTestBinder(), mockSvc)
	r := setupGin()
	r.POST("/list", h.List)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/list",
		strings.NewReader(`{"page":1,"limit":10}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "请求成功")
	mockSvc.AssertExpectations(t)
}

func TestHandler_Create(t *testing.T) {
	mockSvc := new(mockUserService)
	mockSvc.On("Create", mock.Anything, mock.Anything).Return(nil)

	h := newHandler(newTestBinder(), mockSvc)
	r := setupGin()
	r.POST("/create", h.Create)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/create",
		strings.NewReader(`{"username":"testuser","password":"123456","mobile":"13800138000","nick_name":"测试","gender":1,"status":1}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "创建成功")
	mockSvc.AssertExpectations(t)
}

func TestHandler_Detail(t *testing.T) {
	mockSvc := new(mockUserService)
	mockSvc.On("Detail", mock.Anything, uint32(1)).
		Return(&User{Username: "testuser", Mobile: "13800138000"}, nil)

	h := newHandler(newTestBinder(), mockSvc)
	r := setupGin()
	r.POST("/detail", h.Detail)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/detail",
		strings.NewReader(`{"id":1}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "请求成功")
	mockSvc.AssertExpectations(t)
}

func TestHandler_Detail_NotFound(t *testing.T) {
	mockSvc := new(mockUserService)
	mockSvc.On("Detail", mock.Anything, uint32(999)).
		Return(nil, assert.AnError)

	h := newHandler(newTestBinder(), mockSvc)
	r := setupGin()
	r.POST("/detail", h.Detail)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/detail",
		strings.NewReader(`{"id":999}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotContains(t, w.Body.String(), "请求成功")
	mockSvc.AssertExpectations(t)
}
