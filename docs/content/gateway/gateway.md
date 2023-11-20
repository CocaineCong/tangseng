# 网关模块

网关一般对接外部http请求，进行一系列的登陆，鉴权，降级，熔断，转发http，调用rpc请求操作等等...

## 项目结构

```shell
.gateway/
├── cmd              // 启动
├── http             // http入口，controller层
├── middleware       // 存放各种中间件
├── routes           // 路由相关
└── rpc              // 各种下游rpc调用
```
