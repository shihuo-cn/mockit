package mockimpl

import (
	"context"
	"github.com/olivere/elastic/v7"
	"github.com/shihuo-cn/mockit/iface"
	"net/http"
	"reflect"
	"time"
)

var _ iface.DemoInterface = &MockEs{}

func NewMockEs(cli *http.Client) *MockEs {
	url := "http://mockes.com"
	client, err := elastic.NewSimpleClient(elastic.SetURL(url), elastic.SetHttpClient(cli))
	if err != nil {
		panic(err)
	}
	return &MockEs{cli: client}
}

type MockEs struct {
	cli *elastic.Client
}

func (m MockEs) First(ctx context.Context, key int64) (int64, error) {
	panic("implement me")
}

func (m MockEs) Put(ctx context.Context, key, val int64) error {
	panic("implement me")
}

func (m MockEs) List(ctx context.Context, relationId int64, pageIndex, pageSize int) ([]*iface.KV, error) {
	termQuery := elastic.NewTermQuery("rid", relationId)
	res, err := m.cli.Search("_all").
		Query(termQuery).
		Sort("@timestamp", false).
		Size(pageSize).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	msgNum := len(res.Hits.Hits)
	if msgNum == 0 {
		return nil, nil
	}
	list := make([]*iface.KV, msgNum, msgNum)
	for i, item := range res.Each(reflect.TypeOf(&iface.KV{})) {
		list[i] = item.(*iface.KV)
	}
	return list, nil
}

func (m MockEs) Aggregate(ctx context.Context, dm *iface.DemoInterfaceModel, time time.Time) ([]*iface.DemoInterfaceModel, int, error) {
	panic("implement me")
}
