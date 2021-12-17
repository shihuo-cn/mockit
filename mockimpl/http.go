package mockimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shihuo-cn/mockit/iface"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var _ iface.DemoInterface = &MockHttp{}

func NewMockHttp(cli *http.Client) *MockHttp {
	m := &MockHttp{
		cli: cli,
		url: "http://mockhttp.com",
	}
	return m
}

type MockHttp struct {
	cli *http.Client
	url string
}

func (m *MockHttp) First(ctx context.Context, key int64) (int64, error) {
	resp, err := m.cli.Get(m.url + "/first?" + strconv.FormatInt(key, 10))
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, errors.New(resp.Status)
	}
	buff, err := ioutil.ReadAll(resp.Body)
	num, err := strconv.ParseInt(string(buff), 10, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (m *MockHttp) Put(ctx context.Context, key, val int64) error {
	panic("implement me")
}

func (m *MockHttp) List(ctx context.Context, relationId int64, pageIndex, pageSize int) ([]*iface.KV, error) {
	resp, err := m.cli.Get(fmt.Sprintf("%s/list?relationId=%d&pageIndex=%d&pageSize=%d", m.url, relationId, pageIndex, pageSize))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	buff, err := ioutil.ReadAll(resp.Body)
	var arr []*iface.KV
	err = json.Unmarshal(buff, &arr)
	if err != nil {
		return nil, err
	}
	return arr, nil
}

func (m *MockHttp) Aggregate(ctx context.Context, dm *iface.DemoInterfaceModel, time time.Time) ([]*iface.DemoInterfaceModel, int, error) {
	panic("implement me")
}
