# 用户模块

## 1. 文件目录

```shell
user/
├── cmd             // 启动器
└── internal    
    ├── repository  // 存储仓库
    │   └── db      // 存储db操作
    │       └── dao
    └── service     // 具体实现的微服务 
```

## 2. 表设计

```go
type User struct {
    UserID         int64  `gorm:"primarykey"`
	UserName       string `gorm:"unique"`
	NickName       string
	PasswordDigest string
}
```

user表设计的比较简单，因为这不是项目重点，只有一个主键`user_id`，用户名`user_name`，昵称`nick_name`，密码`PasswordDigest`。

密码是加密存储的，所以我们需要对这个`user model`的对象进行一些加解密操作。

加密密码，我们对这个user对象进行`SetPassword`操作，这个 password 就是需要加密的原始password

```go
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), consts.PassWordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}
```

校验密码，我们对user对象进行`CheckPassword`密码操作，password就是需要校验的password，与这个查询出来的user对象进行密码校验。

```go
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}
```

详细用法请看 `app/user/internal/service/user.go` 对dao的调用.
