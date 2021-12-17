# Mockit
![Go](https://github.com/shihuo-cn/mockit/workflows/Go/badge.svg)
## 目标：将mock变得简单，让代码维护变得容易

## 常见Mock难点
- 不同中间件，mock库设计模式不一致，学习代价高，差异化明显
- mock方案强依赖服务端，无法灵活解耦
- 单元测试适配各种中间件的方案后，依赖管理杂乱
- 综上所述不想写mock，也没有良好的可执行方案，放弃单测
## mockit做到了什么
- 统一简化语法
- 无需服务端
- 解耦依赖项
- testMain统一管理

## mockit 使用
### 目前支持
- `Redis`,`MySQL`,`Interface`,`HTTP`
- `GRPC` 可以使用proto生成interface使用`Interface`模拟
- `ElasticSearch`
    + 使用`HTTP`方式代理client,不过es的返回值比较复杂，请求路径没有普通HTTP直观
    + 使用`Interface`方式，将`dao`层抽象成接口方式，这种方式下，接口返回值模拟相对方便直观

> 理论上业务抽象使用`Interface`方式都可达成

### 准备 (具体参照kit_test.go)
1. interface生成 https://github.com/golang/mock
```golang
// interface生成方式
$ mockgen -source ./iface/interface.go -package mockimpl -destination ./mockimpl/interface.go
// 而后将以下new方法的返回值改成interface{}
before: func NewMockDemoInterface(ctrl *gomock.Controller) *MockDemoInterface {
after: func NewMockDemoInterface(ctrl *gomock.Controller) interface{}} {
    mock := &MockDemoInterface{ctrl: ctrl}
    mock.recorder = &MockDemoInterfaceMockRecorder{mock}
    return mock
}
``` 
2. sqlmock依赖replace
> 目前需要替换下sqlmock库，目前pr还在合并中，预计最近2周就能OK
```
replace github.com/DATA-DOG/go-sqlmock v1.5.0 => github.com/Rennbon/go-sqlmock v1.5.1-0.20211212104631-9c4a20760689
```
### mockit自身单测
- 当前目录下 `iface` 中有4个方法， `mockimpl`中分别为各个mock实例的实现
- kit_test.go中`mockSrv`引用了这些`mockimpl`实例，在`testMain`中列举了mock的启动方式，并且在之后所有的Test中介绍了如何使用

## 引用传送
- https://github.com/DATA-DOG/go-sqlmock
- https://github.com/jarcoal/httpmock
- https://github.com/golang/mock
- https://github.com/alicebob/miniredis
- https://github.com/go-redis/redis
- https://github.com/olivere/elastic
- https://github.com/jinzhu/gorm
- https://github.com/go-gorm/gorm

