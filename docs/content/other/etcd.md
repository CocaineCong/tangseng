# etcd 操作

> 本项目的服务注册，服务发现都是基于etcd,具体代码在`pkg/discovery/`下

# instance 实例

注册的微服务，将实例注册到etcd上。

```go
type Server struct {
	Name    string `json:"name"`
	Addr    string `json:"addr"`    // 地址
	Version string `json:"version"` // 版本
	Weight  int64  `json:"weight"`  // 权重
}
```
