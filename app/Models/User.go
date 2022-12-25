package Models

import "time"

//enum position
const (
	POS_INTERN     = "实习成员"
	POS_NORMAL     = "正式成员"
	POS_CONSULTANT = "顾问"
	POS_MANAGER    = "中管"
	POS_MASTER     = "高管"
)

//enum status
const (
	STATUS_NORMAL = "在职"
	STATUS_RETIRE = "退休"
)

type (
	// 数据表表结构体类
	User struct {
		ZjuId      string    `json:"zjuid"`
		QscId      string    `json:"qscid"`
		Password   string    `json:"-"`
		Name       string    `json:"name"`
		Gender     string    `json:"gender"`
		Department string    `json:"department"`
		Direction  string    `json:"direction,omitempty"`
		Position   string    `json:"position"`
		Status     string    `json:"status"`
		Phone      string    `json:"phone"`
		Email      string    `json:"email"`
		Birthday   time.Time `json:"birthday"`
		JoinTime   time.Time `json:"jointime"`
	}

	// 请求结构体类
	UserReq struct {
		QscId    string `json:"qscid"`
		Password string `json:"password"`
	}
)
