package constants

const (
	// 管理员登录信息缓存前缀
	RedisAdminLoginKey = "ZAG:ADMIN:LOGIN"
	// 管理员刷新token缓存前缀
	RedisAdminRefreshTokenKey = "ZAG:ADMIN:REFRESH:TOKEN"
	// 用户登录信息缓存前缀
	RedisUserLoginKey = "ZAG:USER:LOGIN"
	// 用户刷新token缓存前缀
	RedisUserRefreshTokenKey = "ZAG:USER:REFRESH:TOKEN"
)
