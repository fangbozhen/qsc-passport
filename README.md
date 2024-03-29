# 求是潮通行证 Passport v4

## 产品简介

为求是潮web产品提供登录服务，并管理qsc内部人员数据的后台

## 功能需求

- 通过求是潮账号登陆
- 通过浙大统一认证登录
- 查询用户信息
- 修改潮人的人资数据

## 接口说明

### Model

```typescript
// 别问我为什么外面是大写里面是小写首字母，问就是兼容
interface User {
  Name: string          // 真实姓名
  ZjuId: string         // 学号
  LoginType: LoginType  // 用户类型
  
  // 对非潮人为null
  QscUser?: {
    zjuid: string       // 兼容性设计，与`User.ZjuId`相同
    name: string        // 兼容性设计，与`User.Name`相同
    qscid: string       // Qsc Id，注意区分大小写
    gender: string      // 性别 enum {"男", "女"}
    phone: string       // 电话号
    email: string       // 邮箱
    note: string        // 备注信息
    birthday: time.Time
    position: Position  // 身份
    department: string  // 部门
    direction: string   // 部门下分方向
    jointime: string    // RFC3339 注意读出是GMT
    status: string      // 状态【保留】
    privilege: map<string, string>   // 权限组【保留】
  }
}

enum LoginType {
	LT_ZJU = "zju",
	LT_QSC = "qsc",
}

enum Position {
	POS_INTERN     = "实习成员",
	POS_NORMAL     = "正式成员",
	POS_CONSULTANT = "顾问",
	POS_MANAGER    = "中管",
	POS_MASTER     = "高管",
	POS_ADVANCED   = "高级成员",
}

```

### 关于cookie

身份验证基于cookie实现，理论上在同一个域名下的服务可以共享登录状态
如果是在后台调用`/qsc/login`，则必须手动处理cookie缓存（至少golang默认不保存cookie）

### (返回值说明)
- code  错误代码，0表示成功
- err   错误提示信息，成功时为空
  - E_AUTH_FAILED    = 10001 // 未登录
  - E_INTERNAL_ERROR = 10002 // passport内部错误（zjuam或crypto）
  - E_WRONG_REQUEST  = 10003 // 参数错误
  - E_WRONG_USERNAME = 10009 // 用户名错误
  - E_WRONG_PASSWORD = 10004 // 用户名错误
  - E_DATABASE_ERROR = 10005 // db炸了

### 求是潮登录网页 /qsc/login [GET]

- 成功跳转url：?success=xxx
- 【注】跳转URL会添加SESSION_TOKEN=xxx，其值与Header中相同
  
### 求是潮登录API /qsc/login [POST]

- request body
```json
{
    "username": "321010xxxx",
    "password": "abcabc",
}
```
- response
```json
{
    "code": 0,
    "err": "",
    "data": User
}
```

### 浙大统一认证登录 /zju/login [GET]

- 认证流程
    - 浏览器端访问此API
    - 返回302重定向，进入浙大统一认证页面
    - 用户输入账号密码
    - 浙大统一认证平台核验密码，验证通过后继续流程
    - 浙大统一认证平台将用户自动重定向到`/zju/login_success`
    - 后端获取并记录用户信息，返回302重定向到`success_url` （或发送错误，则重定向`failed_url`）
- 注意事项
    - 统一认证登录采用oauth机制，客户端必须保证用户可以在302跳转后可以与统一认证页面交互
    - 用户名及密码错误不会返回；会卡在统一认证界面；需要客户端设计返回上一级的按钮
    - failed_url重定向时会携带query参数：错误码code和错误信息reason
- request query:    success=AAA&fail=BBB
- response [302]
  - 【注】跳转URL会添加SESSION_TOKEN=xxx，其值与Header中相同
  redirect to zjuam

### 浙大登录后重定向地址 /zju/login_success [GET;不需要主动调用]
- query `code=abcabc` （zjuam自动回传）
- response [302]
redirect to `success_url` or `failed_url`

### 登出账户 /logout [GET]
无参数，无返回；如果没登录也不会报错

### 查询用户资料 /profile [GET]
- response
```json
{
    "code": 0,
    "err": "",
    "data": {
        "logined": true,
        "user": User,
    }
}
```


### 求是潮设置密码页面 /qsc/set_password [GET]
- 成功跳转url：?success=xxx


### 求是潮设置密码API /qsc/set_password [POST]
要求在用户已登录的状态下调用（携带cookie）

- request body
```json
{
    "password": "abcabc",
}
```

- 成功无data返回，code=0


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


#### /user/updateMany [POST]

```json
{
    "qscid": $iduser, //qscId 
    "department": , //空即不修改该字段
    "direction": ,
    "postion": ,  
}
```

- example

```json
{
  "qscid": [
    "1238",
    "1239",
    "1240"
  ],
  "department": "摄影",
  "direction": "",
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



