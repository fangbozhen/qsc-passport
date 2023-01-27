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
    ZjuId      string    `json:"zjuid" bson:"ZjuId"`
    QscId      string    `json:"qscid" bson:"QscId"`
    Password   string    `json:"-" bson:"Password"`
    Name       string    `json:"name" bson:"Name"`
    Gender     string    `json:"gender" bson:"Gender"`
    Department string    `json:"department" bson:"Department"`
    Position   string    `json:"position" bson:"Position"`
    Status     string    `json:"status" bson:"Status"`
    Phone      string    `json:"phone" bson:"Phone"`
    Email      string    `json:"email" bson:"Email"`
    Note       string    `json:"note" bson:"Note"`
    Birthday   time.Time `json:"birthday,omitempty" bson:"Birthday"`
    JoinTime   time.Time `json:"jointime,omitempty" bson:"JoinTime"`
}

```
**时间的格式前后端统一使用RFC3339**



# JSON 格式

```json
{
    "err": "",
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
    "qscid": "pta",
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
    "password": "abcabc"
}
```

### /admin GROUP

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
    	"isDescend": boolean //true升序，false降序
	}
}
```
- 注意：写filter的时候每一个键值对的Key都要和 UserProfileQsc(api文档最前面)里面bson中定义的一样，不然会检索不到
- 比如要检索还健在的潮人，应该输入 "Status":"在职"，而不是status
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
    "qscid": "",
    "user": $user, //后端记得检查Password字段
}
```

- example

```json
{
    "qscid": "pta",
    "user": {
        "zjuid": "3210101234",
        "qscid": "pta",
        "name": "陈岩",
        "department": "产研技术",
        "gender": "男",
        "position": "中管",
        "status": "在职",
        "phone": "123",
        "email": "cy@qq.com",
        "note": ""
    }
}
```

**HTML表单**
`<input name="file">`

csv文件格式：
* 只支持csv-utf8格式（excel导出时选定该格式） 
* 从左至右列内容依次为（无需标题行）：浙大学号，求是潮id，姓名，性别，部门，职位，状态，电话，邮箱，生日
* 生日格式为：2003/10/13


#### /user/updateMany [POST]

```json
{
    "qscid": $iduser, //qscId 
    "department": , //空即不修改该字段
    "postion": ,  
}
```

- example

```json
{
    "qscid": ["1238", "1239", "1240"],
    "department": "摄影",
    "position": ""
}
```



#### /user/delete [POST]

request

```json
{
    "qscid": "1242"
}
```



