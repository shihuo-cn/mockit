package iface

import (
	"context"
	"time"
)

// DemoInterface 全局接口，所有实现的中间件或者其他都需遵循这个接口实现mock
type DemoInterface interface {
	First(ctx context.Context, key int64) (int64, error)
	Put(ctx context.Context, key, val int64) error
	List(ctx context.Context, relationId int64, pageIndex, pageSize int) ([]*KV, error)
	Aggregate(ctx context.Context, dm *DemoInterfaceModel, time time.Time) ([]*DemoInterfaceModel, int, error)
}

type DemoType uint
type DemoInterfaceModel struct {
	Type   DemoType
	Date   time.Time
	Int    int
	Float  float64
	Str    string
	Slice  []string
	Slice2 []int
	Struct struct {
		Name  string
		Value int
	}
	Err             error
	InnerModel      *DemoInnerModel
	InnerModelSlice []*DemoInnerModel
}

type DemoInnerModel struct {
	One   int
	Two   string
	Three time.Time
	Four  error
}

type KV struct {
	Id   int64  `gorm:"column:id;primary_key" json:"id" mock:"id"`
	Key  int64  `gorm:"column:key" json:"key" mock:"key"`
	Val  int64  `gorm:"column:val" json:"val" mock:"val"`
	Name string `gorm:"column:name" json:"name" mock:"name"`
}
