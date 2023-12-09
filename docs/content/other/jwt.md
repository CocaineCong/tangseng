# Tangseng



## Token

![image-20231209165044934](https://raw.githubusercontent.com/CremeU/cloud-img/main/image-20231209165044934.png)

token的意思是`令牌`，是服务端生成的一串字符串，作为客户端进行请求的一个标识。当用户第一次登录后，服务器将生成一个token并将此token返回给客户端，以后客户端只需带上这个token前来请求数据即可，无需再次带上用户名和密码。Token的组成:uid(用户唯一的身份标识)、time(当前时间的时间戳)、sign（签名，token的前几位以哈希算法压缩成的一定长度的十六进制字符串。为防止token泄露）。

### 结构

#### JWT头

>`Header`典型的由两部分组成：`token`的类型（“JWT”）和算法名称（比如：`HMAC SHA256`或者`RSA`等等）。

```json
{
"alg": "HS256",
"typ": "JWT"
}
```

- alg属性表示签名使用的算法，默认为HMAC SHA256（写为HS256）；
- typ属性表示令牌的类型，JWT令牌统一写为JWT。



#### 有效载荷

> `Preload`是JWT的主体内容部分，也是一个`JSON对象`，包含需要传递的数据。 JWT指定七个默认字段供选择。

`Registered claims `: 这里有一组预定义的声明，它们不是强制的，但是推荐

- iss：发行人
- exp：到期时间
- sub：主题
- aud：用户
- nbf：在此之前不可用
- iat：发布时间
- jti：JWT ID用于标识该JWT

`Public claims` : 可以随意定义。
`Private claims` : 用于在同意使用它们的各方之间共享信息，并且不是注册的或公开的声明。例如：

```go
{
    "sub": '10086',
    "name": 'FanOne',
    "admin":true
}
```

#### 签名哈希

>`签名哈希`是对上面两部分数据进行签名，通过指定的算法生成哈希，以确保数据不会被篡改。

签名是用于验证消息在传递过程中有没有被更改，并且，对于使用私钥签名的token，它还可以验证JWT的发送方是否为它所称的发送方。



## Double Token

用户正在app或者应用中操作 token突然过期，此时用户不得不返回登陆界面，重新进行一次登录，这种体验性不好，于是引入双token校验机制，首次登陆时服务端返回两个token ，accessToken和refreshToken，accessToken过期时间比较短，refreshToken时间较长，且每次使用后会刷新，每次刷新后的refreshToken都是不同。

accessToken的存在，保证了登录态的正常验证，因其过期时间的短暂也保证了帐号的安全性refreshToekn的存在，保证了用户无需在短时间内进行反复的登陆操作来保证登录态的有效性，同时也保证了活跃用户的登录态可以一直存续而不需要进行重新登录，反复刷新也防止某些不怀好意的人获取refreshToken后对用户帐号进行动手动脚的操作。

![image-20231209170311461](https://raw.githubusercontent.com/CremeU/cloud-img/main/image-20231209170311461.png)

在进行服务器请求的时候，先将Token发送验证，如果accessToken有效，则正常返回请求结果；如果accessToken无效，则验证refreshToken。

此时如果refreshToken有效则返回请求结果和新的accessToken和新的refreshToken。如果refreshToken无效，则提示用户进行重新登陆操作。

## 具体实现



```go
type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}
```



签发进行token

```go
func GenerateToken(id int64, username string) (accessToken, refreshToken string, err error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(consts.AccessTokenExpireDuration)
	rtExpireTime := nowTime.Add(consts.RefreshTokenExpireDuration)
	claims := Claims{
		ID:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "mall",
		},
	}
	// 加密并获得完整的编码后的字符串token
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: rtExpireTime.Unix(),
		Issuer:    "search-engine",
	}).SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
}
```



验证用户的token

```go
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
```



ParseRefreshToken 验证用户token

```go
func ParseRefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	accessClaim, err := ParseToken(aToken)
	if err != nil {
		return
	}

	refreshClaim, err := ParseToken(rToken)
	if err != nil {
		return
	}

	if accessClaim.ExpiresAt > time.Now().Unix() {
		// 如果 access_token 没过期,每一次请求都刷新 refresh_token 和 access_token
		return GenerateToken(accessClaim.ID, accessClaim.Username)
	}

	if refreshClaim.ExpiresAt > time.Now().Unix() {
		// 如果 access_token 过期了,但是 refresh_token 没过期, 刷新 refresh_token 和 access_token
		return GenerateToken(accessClaim.ID, accessClaim.Username)
	}

	// 如果两者都过期了,重新登陆
	return "", "", errors.New("身份过期，重新登陆")
}
```
