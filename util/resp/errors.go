package resp

const (
	E_AUTH_FAILED    = 10001 // 未登录
	E_INTERNAL_ERROR = 10002 // passport内部错误（zjuam或crypto）
	E_WRONG_REQUEST  = 10003 // 参数错误
	E_WRONG_USERNAME = 10009 // 用户名错误
	E_WRONG_PASSWORD = 10004 // 用户名错误
	E_DATABASE_ERROR = 10005 // db炸了
)
