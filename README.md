# 求是潮通行证 Ver.3

## 一、产品简介

为求是潮web产品提供登录服务

## 二、功能需求

- 通过求是潮账号登陆
- 通过浙大统一认证登录
- 查询用户信息

## 三、接口说明

### 关于cookie
身份验证基于cookie实现，理论上在同一个域名下的服务可以共享登录状态
如果是在后台调用`/qsc/login`，则必须手动处理cookie缓存（至少golang默认不保存cookie）

### (返回值说明)
- code  错误代码，0表示成功
- err   错误提示信息，成功时为空

### 求是潮登录网页 /qsc/login [GET]
- 成功跳转url：?success=xxx

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
    "data": $user
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
        "user": $user
    }
}
```

### 求是潮设置密码 /qsc/set_password [GET]
- 成功跳转url：?success=xxx


### 求是潮登录API /qsc/login [POST]
要求在用户已登录的状态下调用（携带cookie）

- request body
```json
{
    "password": "abcabc",
}
- 成功无data返回，code=0

## 四、版本说明

## 五、流程规划
