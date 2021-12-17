package mockimpl

import (
	"context"
	"github.com/shihuo-cn/mockit/iface"
	"gorm.io/gorm"
	"time"
)

var _ iface.DemoInterface = &mockMySQL{}

type mockMySQL struct {
	db *gorm.DB
}

func NewSqlDao(db *gorm.DB) iface.DemoInterface {
	return &mockMySQL{
		db: db,
	}
}
func (s *mockMySQL) tableDemo(ctx context.Context) *gorm.DB {
	return s.db.Table("demo")
}

func (s *mockMySQL) First(ctx context.Context, key int64) (int64, error) {
	var count int64
	err := s.tableDemo(ctx).
		Where("key = ?", key).
		Select("COUNT(1)").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *mockMySQL) Put(ctx context.Context, key, val int64) error {
	m := &iface.KV{
		Key:  key,
		Val:  val,
		Name: "name",
	}
	err := s.tableDemo(ctx).Create(m).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *mockMySQL) List(ctx context.Context, relationId int64, pageIndex, pageSize int) ([]*iface.KV, error) {
	var arr []*iface.KV
	err := s.tableDemo(ctx).
		Where("relation_id = ? AND is_deleted = ?",
			relationId, false).
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Scan(&arr).Error
	if err != nil {
		return nil, err
	}
	return arr, nil
}

func (s *mockMySQL) Aggregate(ctx context.Context, m *iface.DemoInterfaceModel, time time.Time) ([]*iface.DemoInterfaceModel, int, error) {
	panic("implement me")
}
