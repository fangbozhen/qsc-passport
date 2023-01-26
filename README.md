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
    "Qscid": "321010xxxx",
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
    "pageNumber":"", // 从1开始
    "pageSize":"" ,  
    "filter":{
      "selector1": "",
      "selector2": "",
      // ...返回的是所有的selector的交集
    }, 
    "sortBy": {
    	"col":"",
    	"isDescend": boolean //0升序，1降序
	}
}
```
- 注意：写filter的时候每一个键值对的Key都要和 UserProfileQsc(api文档最前面)里面定义的一样，不然会检索不到
- 比如要检索还健在的潮人，应该输入 "Status":"alive"，而不是status
- 在比如检索qscid的时候，是QscId 而不是 Qscid

response

```json
{
   "code": 0,
   "err": "",
   "data": $userarray, //返回一个User数组，包含要呈现的User信息
}
```
- example: 返回是这样的
```json
{
  "code": "",
  "err": "",
  "data":{
    "users": {
      [{user1},{user2},{user3}]
    }
  }
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



