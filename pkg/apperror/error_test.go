package apperror_test

import (
	"errors"
	"fmt"
	"testing"

	"zero-backend/pkg/apperror"
)

// 测试用错误码
var (
	testUserNotFound = apperror.NewCode(404001, "USER.NOT_FOUND", "用户不存在")
	testFileTooLarge = apperror.NewCode(400001, "FILE.TOO_LARGE", "文件大小不能超过 %d MB，当前 %d MB")
	testDBError      = apperror.NewCode(500001, "DB.ERROR", "数据库操作失败")
)

func TestCode_Value(t *testing.T) {
	if got := testUserNotFound.Value(); got != 404001 {
		t.Errorf("Value() = %d, want 404001", got)
	}
}

func TestCode_String(t *testing.T) {
	if got := testUserNotFound.String(); got != "USER.NOT_FOUND" {
		t.Errorf("String() = %q, want %q", got, "USER.NOT_FOUND")
	}
}

func TestCode_Template(t *testing.T) {
	if got := testUserNotFound.Template(); got != "用户不存在" {
		t.Errorf("Template() = %q, want %q", got, "用户不存在")
	}
}

func TestNew_Basic(t *testing.T) {
	err := apperror.New(testUserNotFound)

	if err.Error() != "用户不存在" {
		t.Errorf("Error() = %q, want %q", err.Error(), "用户不存在")
	}
	if err.Code() != testUserNotFound {
		t.Errorf("Code() mismatch")
	}
	if err.Cause() != nil {
		t.Errorf("Cause() should be nil")
	}
}

func TestNew_WithCause(t *testing.T) {
	cause := errors.New("sql: no rows")
	err := apperror.New(testDBError, apperror.WithCause(cause))

	if err.Cause() != cause {
		t.Errorf("Cause() = %v, want %v", err.Cause(), cause)
	}
	if err.Unwrap() != cause {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestNew_WithArgs(t *testing.T) {
	err := apperror.New(testFileTooLarge, apperror.WithArgs(10, 25))

	want := "文件大小不能超过 10 MB，当前 25 MB"
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestNew_WithCauseAndArgs(t *testing.T) {
	cause := errors.New("file too large")
	err := apperror.New(testFileTooLarge, apperror.WithCause(cause), apperror.WithArgs(5, 100))

	if err.Cause() != cause {
		t.Errorf("Cause() mismatch")
	}
	want := "文件大小不能超过 5 MB，当前 100 MB"
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestError_Cached(t *testing.T) {
	err := apperror.New(testFileTooLarge, apperror.WithArgs(10, 20))

	msg1 := err.Error()
	msg2 := err.Error()
	if msg1 != msg2 {
		t.Errorf("Error() should return cached message")
	}
}

func TestWrap(t *testing.T) {
	cause := errors.New("connection refused")
	err := apperror.Wrap(testDBError, cause)

	if err.Error() != "数据库操作失败" {
		t.Errorf("Error() = %q, want %q", err.Error(), "数据库操作失败")
	}
	if err.Code() != testDBError {
		t.Errorf("Code() mismatch")
	}
	if err.Cause() != cause {
		t.Errorf("Cause() = %v, want %v", err.Cause(), cause)
	}
}

func TestWrap_WithArgs(t *testing.T) {
	cause := errors.New("file validation failed")
	err := apperror.Wrap(testFileTooLarge, cause, 10, 50)

	if err.Cause() != cause {
		t.Errorf("Cause() mismatch")
	}
	want := "文件大小不能超过 10 MB，当前 50 MB"
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestError_Is_SameCode(t *testing.T) {
	err1 := apperror.New(testUserNotFound)
	err2 := apperror.New(testUserNotFound)

	if !errors.Is(err1, err2) {
		t.Error("errors.Is should return true for same Code")
	}
}

func TestError_Is_DifferentCode(t *testing.T) {
	err1 := apperror.New(testUserNotFound)
	err2 := apperror.New(testDBError)

	if errors.Is(err1, err2) {
		t.Error("errors.Is should return false for different Code")
	}
}

func TestError_Is_NonAppError(t *testing.T) {
	err := apperror.New(testUserNotFound)
	plain := errors.New("plain error")

	if errors.Is(err, plain) {
		t.Error("errors.Is should return false for non-*Error target")
	}
}

func TestError_AsType(t *testing.T) {
	err := apperror.New(testUserNotFound)

	appErr, ok := errors.AsType[*apperror.Error](err)
	if !ok {
		t.Fatal("errors.AsType should return true")
	}
	if appErr.Code() != testUserNotFound {
		t.Errorf("Code() mismatch")
	}
}

func TestError_AsType_Wrapped(t *testing.T) {
	inner := apperror.New(testUserNotFound)
	wrapped := fmt.Errorf("service layer: %w", inner)

	appErr, ok := errors.AsType[*apperror.Error](wrapped)
	if !ok {
		t.Fatal("errors.AsType should unwrap and find *Error")
	}
	if appErr.Code() != testUserNotFound {
		t.Errorf("Code() mismatch")
	}
}

func TestError_Unwrap_Chain(t *testing.T) {
	root := errors.New("sql: duplicate key")
	dbErr := apperror.Wrap(testDBError, root)
	svcErr := apperror.New(testDBError, apperror.WithCause(dbErr))

	// 第一层 unwrap
	if errors.Unwrap(svcErr) != dbErr {
		t.Errorf("first Unwrap should return dbErr")
	}

	// 第二层 unwrap
	if errors.Unwrap(dbErr) != root {
		t.Errorf("second Unwrap should return root")
	}

	// errors.Is 应能穿透整个链
	if !errors.Is(svcErr, dbErr) {
		t.Error("errors.Is should find dbErr in chain")
	}
}

func TestError_Format_VerbV(t *testing.T) {
	err := apperror.New(testUserNotFound)
	got := fmt.Sprintf("%v", err)
	if got != "用户不存在" {
		t.Errorf("%%v = %q, want %q", got, "用户不存在")
	}
}

func TestError_Format_VerbPlusV(t *testing.T) {
	cause := errors.New("db timeout")
	err := apperror.New(testDBError, apperror.WithCause(cause))

	got := fmt.Sprintf("%+v", err)

	for _, want := range []string{"apperror.Error{", "DB.ERROR", "500001", "数据库操作失败", "db timeout"} {
		if !contains(got, want) {
			t.Errorf("%%+v output missing %q\nGot: %s", want, got)
		}
	}
}

func TestError_Format_VerbS(t *testing.T) {
	err := apperror.New(testUserNotFound)
	got := fmt.Sprintf("%s", err)
	if got != "用户不存在" {
		t.Errorf("%%s = %q, want %q", got, "用户不存在")
	}
}

func TestError_Format_VerbQ(t *testing.T) {
	err := apperror.New(testUserNotFound)
	got := fmt.Sprintf("%q", err)
	if got != `"用户不存在"` {
		t.Errorf("%%q = %q, want %q", got, `"用户不存在"`)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
