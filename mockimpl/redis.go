package mockimpl

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/shihuo-cn/mockit/iface"
	"time"
)

var _ iface.DemoInterface = &MockRedis{}

type MockRedis struct {
	cli *redis.Client
	kf  keyFormatter
}
type keyFormatter struct {
}

func (keyFormatter) formatKey(key int64) string {
	return fmt.Sprintf("SHOHUO:%d", key)
}
func (m MockRedis) First(ctx context.Context, key int64) (int64, error) {
	k := m.kf.formatKey(key)
	num, err := m.cli.Get(ctx, k).Int64()
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (m MockRedis) Put(ctx context.Context, key, val int64) error {
	k := m.kf.formatKey(key)
	err := m.cli.Set(ctx, k, val, 0).Err()
	return err
}

func (m MockRedis) List(ctx context.Context, relationId int64, pageIndex, pageSize int) ([]*iface.KV, error) {
	panic("implement me")
}

func (m MockRedis) Aggregate(ctx context.Context, dm *iface.DemoInterfaceModel, time time.Time) ([]*iface.DemoInterfaceModel, int, error) {
	panic("implement me")
}

func NewMockRedis(addr string) *MockRedis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	return &MockRedis{
		cli: rdb,
	}
}
