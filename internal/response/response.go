package response

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"zero-backend/internal/apperror"
	"zero-backend/internal/ctxkeys"
	"zero-backend/internal/logger"

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
	var userError *apperror.UserError
	if errors.As(err, &userError) {
		output(c, Response{
			ErrCode: int(userError.Code),
			Message: userError.Message,
		})
		return
	}

	var systemError *apperror.SystemError
	if errors.As(err, &systemError) {
		output(c, Response{
			ErrCode: int(systemError.Code),
			Message: systemError.Message,
		})
		return
	}

	var unauthorizedError *apperror.UnauthorizedError
	if errors.As(err, &unauthorizedError) {
		output(c, Response{
			ErrCode: int(unauthorizedError.Code),
			Message: unauthorizedError.Message,
		})
		return
	}

	systemError = apperror.NewSystemError(err, "系统异常")
	output(c, Response{
		ErrCode: int(systemError.Code),
		Message: systemError.Message,
	})
}

// Success 输出成功信息
func Success(c *gin.Context, message string, data any) {
	if message == "" {
		message = "success"
	}

	output(c, Response{
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
func output(c *gin.Context, resp Response) {
	var traceId string
	if val := c.Request.Context().Value(ctxkeys.TraceIDKey{}); val != nil {
		traceId = val.(string)
	}

	var cost time.Duration
	if val := c.Request.Context().Value(ctxkeys.BeginTimeKey{}); val != nil {
		cost = time.Since(val.(time.Time))
	}

	resp.Cost = fmt.Sprintf("%.4f", cost.Seconds())
	resp.TraceId = traceId

	logger.Ctx(c.Request.Context()).Info("Response",
		"message", resp.Message,
		"cost", resp.Cost,
		"data", resp.Data,
		"errCode", resp.ErrCode)

	c.JSON(http.StatusOK, resp)
}
