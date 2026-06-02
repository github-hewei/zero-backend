package response

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"zero-backend/internal/apperror"
	"zero-backend/internal/ctxkeys"
	"zero-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Response 响应数据结构体
type Response struct {
	ErrCode int    `json:"errcode"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Cost    string `json:"cost"`
	TraceId string `json:"traceId"`
}

// Error 输出错误信息
func Error(c *gin.Context, err error) {
	var coded apperror.Coded
	if errors.As(err, &coded) {
		code, msg := coded.CodePair()
		output(c, errors.Unwrap(err), Response{
			ErrCode: int(code),
			Message: msg,
		})
		return
	}

	// 未识别的错误统一包装为 SystemError
	systemError := apperror.NewSystemError(err, "系统异常")
	output(c, systemError.Err, Response{
		ErrCode: int(systemError.Code),
		Message: systemError.Message,
	})
}

// Success 输出成功信息
func Success(c *gin.Context, message string, data any) {
	if message == "" {
		message = "success"
	}

	output(c, nil, Response{
		ErrCode: int(apperror.ErrorCodeNone),
		Message: message,
		Data:    data,
	})
}

// SetCookie 设置cookie
func SetCookie(c *gin.Context, name, value string, maxAge int, path string) {
	// 设置 SameSite 允许Cookie跨站
	c.SetSameSite(http.SameSiteNoneMode)
	// 设置 Secure 强制Https请求（本地localhost除外）
	c.SetCookie(name, value, maxAge, path, "", true, true)
	// 接口调试模式
	//c.SetCookie(name, value, maxAge, "", "", false, false)
}

// output 输出响应JSON
func output(c *gin.Context, err error, resp Response) {
	traceId := ctxkeys.TraceID(c.Request.Context())

	beginTime, _ := ctxkeys.BeginTime(c.Request.Context())
	var cost time.Duration
	if !beginTime.IsZero() {
		cost = time.Since(beginTime)
	}

	resp.Cost = fmt.Sprintf("%.4f", cost.Seconds())
	resp.TraceId = traceId
	l := logger.Ctx(c.Request.Context())

	if err != nil {
		l.Err(err, "Error")
	}

	l.Info("Response",
		"msg", resp.Message,
		"cost", resp.Cost,
		"data", resp.Data,
		"errCode", resp.ErrCode)

	c.JSON(http.StatusOK, resp)
}
