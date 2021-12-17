module github.com/shihuo-cn/mockit

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a // indirect
	github.com/alicebob/miniredis v2.5.0+incompatible
	github.com/go-redis/redis/v8 v8.11.4
	github.com/golang/mock v1.6.0
	github.com/gomodule/redigo v1.8.6 // indirect
	github.com/jarcoal/httpmock v1.0.8
	github.com/olivere/elastic/v7 v7.0.29
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.6.1
	github.com/yuin/gopher-lua v0.0.0-20210529063254-f4c35e4016d9 // indirect
	gorm.io/driver/mysql v1.2.1
	gorm.io/gorm v1.22.4
)

replace github.com/DATA-DOG/go-sqlmock v1.5.0 => github.com/Rennbon/go-sqlmock v1.5.1-0.20211212104631-9c4a20760689
