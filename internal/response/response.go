package response

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"zero-backend/internal/ctxkeys"
	"zero-backend/internal/errcode"

	"github.com/241x/zero-kit/logger"

	"github.com/241x/zero-kit/apperror"

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

// Error 输出错误响应
// 优先解析 *apperror.Error，提取错误码和消息；未识别的错误兜底为 Internal
func Error(c *gin.Context, err error) {
	if appErr, ok := errors.AsType[*apperror.Error](err); ok {
		output(c, appErr.Cause(), Response{
			ErrCode: appErr.Code().Value(),
			Message: appErr.Error(),
		})
		return
	}

	// 非 apperror 错误兜底
	output(c, err, Response{
		ErrCode: errcode.Internal.Value(),
		Message: errcode.Internal.Template(),
	})
}

// Success 输出成功响应
func Success(c *gin.Context, message string, data any) {
	if message == "" {
		message = "success"
	}

	output(c, nil, Response{
		ErrCode: 0,
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
