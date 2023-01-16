# API 接口文档

## user类

```go
type User struct {
	Name      string          `json:"Name"`
	ZjuId     string          `json:"ZjuId"`
	LoginType string          `json:"LoginType"` // @see enum LoginType
	QscUser   *UserProfileQsc `json:"QscUser,omitempty"`
}
type UserProfileQsc struct {
    ZjuId      string    `json:"zjuid"`
    QscId      string    `json:"qscid"`
    Password   string    `json:"-"`
    Name       string    `json:"name"`
    Gender     string    `json:"gender"`
    Department string    `json:"department"`
    Position   string    `json:"position"`
    Status     string    `json:"status"`
    Phone      string    `json:"phone"`
    Email      string    `json:"email"`
    Note       string    `json:"note"`
    Birthday   time.Time `json:"birthday"`
    JoinTime   time.Time `json:"jointime"`
}
```



# JSON 格式

```json
{
    "err": ,
    "code":0,
    "data":{}, 
}
```

出错

## 路由

### /qsc/login [GET]

前端网页

### /qsc/login [POST]

request

```json
{
    "username": "321010xxxx",
    "password": "abcabc",
}
```

### /zju/login [GET]

提供统一认证按钮

### /logout [GET]

### /user [GET]

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

### /qsc/reset-password [GET]

改密码界面

### /qsc/reset-password [POST]

request

```json
{
    "password": "abcabc",
}
```

### /admin GROUP

#### /login [GET]

#### /login [POST]

Cookie记得改

#### /user/list [GET]

request

```json
{
    "pageNumber": , // 从0开始
    "pageSize" ,  
    "filter": $user, 
    "sortby": {
    	"col":,
    	"isDescend": boolean //0升序，1降序
	}
}
```

reponse

```json
{
   "code": 0,
   "err": "",
   "data": $userarray, //返回一个User数组，包含要呈现的User信息
}
```

#### /user/updateOne [POST]

request

```json
{
    "qscid": ,
    "user": $user, //后端记得检查Password字段
}
```

**HTML表单**

#### /user/updateMany [POST]

```json
{
    "id": $iduser, //qscId 
    "department": , //空即不修改该字段
    "postion": ,  
}
```



#### /user/delete [POST]

request

```json
{
    "qscid": ,
}
```



