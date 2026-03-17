package middleware

// Middlewares 中间键集合
type Middlewares struct {
	Before *BeforeMiddleware
	Cors   *CorsMiddleware
}
