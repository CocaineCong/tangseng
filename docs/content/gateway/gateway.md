# 网关模块

网关一般对接外部http请求，进行一系列的鉴权，降级，熔断，转发http，调用rpc请求操作等等...

## 项目结构

```shell
.gateway/
├── cmd              // 启动
├── http             // http入口，controller层
├── middleware       // 存放各种中间件
├── routes           // 路由相关
└── rpc              // 各种下游rpc调用
```

tangseng 的网关相对来说就简单很多了，只有承接http请求，鉴权，rpc调用，后续再加上降级，熔断。

## rpc调用

定义客户端RPC调用器

```go
var (
    Register   *discovery.Resolver
    ctx        context.Context
    CancelFunc context.CancelFunc

    UserClient          user.UserServiceClient
    FavoriteClient      favorite.FavoritesServiceClient
    SearchEngineClient  search_engine.SearchEngineServiceClient
    IndexPlatformClient index_platform.IndexPlatformServiceClient
    SearchVectorClient  search_vector.SearchVectorServiceClient
)
```

初始化所有的rpc请求

```go
func Init() {
    Register = discovery.NewResolver([]string{config.Conf.Etcd.Address}, logrus.New())
    resolver.Register(Register)
    ctx, CancelFunc = context.WithTimeout(context.Background(), 3*time.Second)

    defer Register.Close()
    initClient(config.Conf.Domain[consts.UserServiceName].Name, &UserClient)
    initClient(config.Conf.Domain[consts.FavoriteServiceName].Name, &FavoriteClient)
    initClient(config.Conf.Domain[consts.SearchServiceName].Name, &SearchEngineClient)initClient(config.Conf.Domain[consts.IndexPlatformName].Name, &IndexPlatformClient)
    initClient(config.Conf.Domain[consts.SearchVectorName].Name, &SearchVectorClient)
}
```

传入对应的`service name`并连接相应的rpc请求，如果有负载均衡则配置负载均衡即可。

```go
func connectServer(serviceName string) (conn *grpc.ClientConn, err error) {
    opts := []grpc.DialOption{
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    }
    addr := fmt.Sprintf("%s:///%s", Register.Scheme(), serviceName)

    // Load balance
    if config.Conf.Services[serviceName].LoadBalance {
        log.Printf("load balance enabled for %s\n", serviceName)
        opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
    }

    conn, err = grpc.DialContext(ctx, addr, opts...)
    return
}
```
