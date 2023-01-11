package resp

const (
	AuthFailedError    = 10001 // 未登录
	InternalError      = 10002 // passport内部错误
	WrongRequestError  = 10003 // 参数错误
	WrongPasswordError = 10004 // 密码错误
	WrongUsernameError = 10005
	DatabaseError      = 10006 // db炸了
)
