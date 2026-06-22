package middleware

// Middlewares 中间件集合
type Middlewares struct {
	Trace         *TraceMiddleware
	RequestLogger *RequestLogger
}
