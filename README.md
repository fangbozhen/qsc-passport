# QSC passport

## 功能

**后端登录模块直接使用pta的passportv4的部分**

 - [x] 潮人登录
   - [ ] QSCid登录
   - [ ] 浙大统一认证登录
   - [ ] 修改密码
 - [ ] 人资管理系统
    - [ ] 用户组批量权限设置
    - [ ] 基础的筛选、检索功能
    - [ ] 操作日志记录
    - [ ] excel批量导入成员信息

## model

User类型
```go
type (
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
)
```

## Json格式

### 错误代码
```go
const (
	AuthFailedError    = 10001 // 未登录
	InternalError      = 10002 // passport内部错误
	WrongRequestError  = 10003 // 参数错误
	WrongPasswordError = 10004 // 密码错误
	WrongUsernameError = 10005
	DatabaseError      = 10006 // db炸了
)
```

## 路由

### 求是潮登录界面 /qsc/login [GET]

提供求是潮id登录和密码，以及浙大统一认证

 - 成功跳转url:`?success=qscid`

### 求是潮登录API /qsc/login [POST]

request body

```json
{
   "code": 0,
   "error": "",
   "data": {
      "qscid": "",
      "password": ""
   }
}
```

成功即返回`code=0`，无data

### 浙大统一认证登录 /zju/login [GET]

暂略

### 登出账号 /logout [GET]

暂略

### 用户界面 /users/:id [GET]

对所有潮人提供基本信息的查询

要求在登录状态，携带cookie，其他无要求

response

```json
{
   "code": 0,
   "err": "",
   "data": {
      "logined": true,
      "user": $User
   }
}
```

#### 修改密码API /qsc/reset-password [GET]

提供修改密码服务

输入原密码，新密码，重复输入新密码

#### 修改密码API /qsc/reset-password [POST]

request body
```json
{
   "code": 0,
   "err": "",
   "data": {
      "origin-password": "",
      "new-password": ""
   }
}
```

成功即返回`code=0`，无data

### 后台管理界面 /admin Group

默认跳转到/login，携带cookie跳转到/index

#### 登录 /login [GET]

登录界面，基本同正常登录

#### 登录API /login [POST]

同上，后端验证是否有权限登录后台

#### 主页 /index [GET]

以列表形式呈现所有潮人（分页），呈现内容即`User`类中的内容 

基本的筛选、搜索功能

response
```json
{
   "code": 0,
   "err": "",
   "data": $userarray, //返回一个User数组，包含要呈现的User信息
}
```

#### 修改用户信息 /user/:id [GET]

显示单个用户所有信息，并提供修改，以及删除用户按钮

Response即单个用户信息`User`

#### 修改用户信息 /user/:id/update [PUT]

request 发送最新的`User`信息
成功即返回`code=0`，无data

#### 删除用户 /user/:id/delete [DELETE]

request

```json
{
   "code": 0,
   "err": "",
   "data": {
      "qscid": ""
   }
}
```


#### 操作日志 /log [GET]

暂略