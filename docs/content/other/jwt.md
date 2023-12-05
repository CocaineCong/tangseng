# JWT 鉴权

tangseng中的登陆采用的是双token的验证方式。assessToken 和 refreshToken 进行验证。

通常来说：

assessToken:验证用户身份并且过期时间较短，一般过期时间是2～3个小时。

refreshToken:刷新assessToken时间，一般过期时间为2～3天。
