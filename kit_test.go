package mockit

import (
	"context"
	"errors"
	"github.com/jarcoal/httpmock"
	"github.com/shihuo-cn/mockit/iface"
	"github.com/shihuo-cn/mockit/mockimpl"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"os"
	"testing"
	"time"
)

// 以下所有方式都实现了 mockglobal.DemoInterface
// NOTE: 这里对于使用Interface或者struct无强制要求，2者都支持
type mockSrv struct {
	es      *mockimpl.MockEs    // struct 方式
	httpCli *mockimpl.MockHttp  // struct 方式
	iface   iface.DemoInterface // interface方式
	redis   *mockimpl.MockRedis // struct 方式
	mysql   iface.DemoInterface // interface方式
}

var (
	kit Mockit
	ctx = context.Background()
	srv *mockSrv
)

func TestMain(m *testing.M) {
	var err error
	// interface intercept agent
	// 初始化kit
	kit, err = New(mockimpl.NewMockDemoInterface)
	if err != nil {
		panic(err)
	}

	// http client 拦截
	{ // 初始化mock service
		cli := http.DefaultClient
		kit.InterceptHttpClient(cli)
		mockInterface := kit.GetInterfaceClient("MockDemoInterface")
		iface := mockInterface.(*mockimpl.MockDemoInterface)
		// gorm2 io
		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      kit.SqlDB(),
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		// 拦截
		srv = &mockSrv{
			es:      mockimpl.NewMockEs(cli),
			httpCli: mockimpl.NewMockHttp(cli),
			iface:   iface,
			redis:   mockimpl.NewMockRedis(kit.RedisAddr()),
			mysql:   mockimpl.NewSqlDao(gormDB),
		}
	}
	os.Exit(m.Run())
}

func TestMockInterface(t *testing.T) {
	t.Run("easy test", func(t *testing.T) {
		var (
			arg         = int64(1)
			outExpected = int64(20)
			errExpected = errors.New("this is error")
		)
		param := NewExpectParam().
			WithMethod("First").
			WithPath("MockDemoInterface").
			WithArgs(ctx, arg).
			WithReturns(outExpected, errExpected)
		kit.InterfaceExpect(param, t)
		outTmp, errTmp := srv.iface.First(ctx, arg)
		assert.EqualValues(t, errExpected, errTmp)
		assert.EqualValues(t, outExpected, outTmp)
	})
	t.Run("complex test", func(t *testing.T) {
		var (
			aggTime     = time.Now()
			aggExcepted = &iface.DemoInterfaceModel{
				Type:   1,
				Date:   time.Now(),
				Int:    2,
				Float:  3.75,
				Str:    "this is string",
				Slice:  []string{"str1", "str2"},
				Slice2: []int{1, 2},
				Struct: struct {
					Name  string
					Value int
				}{
					Name:  "name",
					Value: 100,
				},
				Err: errors.New("this is error"),
				InnerModel: &iface.DemoInnerModel{
					One: 0,
				},
				InnerModelSlice: []*iface.DemoInnerModel{
					{
						One:   1,
						Two:   "two",
						Three: time.Now(),
						Four:  errors.New("four"),
					},
					{
						One:   11,
						Two:   "two too",
						Three: time.Now().AddDate(1, 0, 0),
						Four:  errors.New("four too"),
					},
				},
			}
			return1Excepted = []*iface.DemoInterfaceModel{aggExcepted, aggExcepted, aggExcepted}
			return2Excepted = 2222
			return3Excepted = errors.New("complex test")
		)

		param2 := NewExpectParam().
			WithMethod("Aggregate").
			WithPath("MockDemoInterface").
			WithArgsByIndex(1, aggExcepted).
			WithArgsByIndex(2, aggTime).
			WithReturns(return1Excepted, return2Excepted, return3Excepted)
		kit.InterfaceExpect(param2, t)

		re1, re2, re3 := srv.iface.Aggregate(ctx, aggExcepted, aggTime)
		assert.EqualValues(t, return1Excepted, re1)
		assert.EqualValues(t, return2Excepted, re2)
		assert.EqualValues(t, return3Excepted, re3)
	})
}

func TestMockMySql(t *testing.T) {
	t.Run("insert", func(t *testing.T) {
		p := NewExpectParam().
			WithMethod("INSERT").
			WithReturns(1, 1)
		kit.MysqlExecExpect(p, t)
		err := srv.mysql.Put(ctx, 1, 2)
		assert.Nil(t, err)
	})
	t.Run("insert with args", func(t *testing.T) {
		// INSERT INTO `demo` (`key`,`val`,`name`) VALUES (?,?,?)
		arg1 := int64(5)
		arg2 := int64(10)
		p := NewExpectParam().
			WithMethod("^INSERT INTO `demo` \\(`key`,`val`,`name`\\) VALUES (.+)").
			WithArgs(arg1, arg2, "name").
			WithReturns(1, 1)
		kit.MysqlExecExpect(p, t)
		err := srv.mysql.Put(ctx, arg1, arg2)
		assert.Nil(t, err)
	})
	t.Run("select count", func(t *testing.T) {
		expectCount := int64(20)
		arg := int64(10)
		p := NewExpectParam().
			WithMethod("SELECT COUNT\\(1\\) FROM `demo` WHERE key = (.+)").
			WithArgs(arg).
			WithKeyValReturn("COUNT", expectCount)
		kit.MysqlQueryExpect(p, t)
		count, err := srv.mysql.First(ctx, arg)
		assert.Nil(t, err)
		assert.Equal(t, expectCount, count)
	})

	t.Run("select list", func(t *testing.T) {
		relationId := int64(20)
		pageIndex := 2
		pageSize := 20
		expectedReturns := []*iface.KV{
			{
				Id:   1,
				Key:  1,
				Val:  100,
				Name: "k1",
			},
			{
				Id:   2,
				Key:  2,
				Val:  200,
				Name: "k2",
			},
		}
		p := NewExpectParam().
			WithMethod("SELECT (.+) FROM `demo` WHERE relation_id = (.+)").
			WithArgs(relationId, false).
			WithReturns(expectedReturns)
		kit.MysqlQueryExpect(p, t)
		res, err := srv.mysql.List(ctx, relationId, pageIndex, pageSize)
		assert.Nil(t, err)
		assert.EqualValues(t, expectedReturns, res)
	})
	t.Run("error return", func(t *testing.T) {
		relationId := int64(200)
		pageIndex := 2
		pageSize := 20
		errExpected := gorm.ErrRecordNotFound
		p := NewExpectParam().
			WithMethod("SELECT (.+) FROM `demo` WHERE relation_id = (.+)").
			WithArgs(relationId, false).
			WithReturns(errExpected)
		kit.MysqlQueryExpect(p, t)
		res, err := srv.mysql.List(ctx, relationId, pageIndex, pageSize)
		assert.Equal(t, errExpected, err)
		assert.Nil(t, res)
	})
}

func TestMockRedis(t *testing.T) {
	keyExcepted := int64(10)
	valExcepted := int64(20)
	err := srv.redis.Put(ctx, keyExcepted, valExcepted)
	assert.Nil(t, err)
	val, err := srv.redis.First(ctx, keyExcepted)
	assert.Nil(t, err)
	assert.Equal(t, valExcepted, val)
}

func TestMockHttp(t *testing.T) {
	t.Run("get string response with kv", func(t *testing.T) {
		p := NewExpectParam().
			WithMethod("GET").
			WithPath("http://mockhttp.com/first").
			WithKeyValReturn(HttpResponseString, "20")
		kit.HttpExpect(p, t)
		num, err := srv.httpCli.First(ctx, 10)
		assert.Nil(t, err)
		assert.Equal(t, int64(20), num)
	})
	t.Run("get string response with returns", func(t *testing.T) {
		p := NewExpectParam().
			WithMethod("GET").
			WithPath("http://mockhttp.com/first").
			WithReturns("20")
		kit.HttpExpect(p, t)
		num, err := srv.httpCli.First(ctx, 10)
		assert.Nil(t, err)
		assert.Equal(t, int64(20), num)
	})
	t.Run("list json response with kv", func(t *testing.T) {
		expectResponse := []*iface.KV{
			{
				Id:   1,
				Key:  10,
				Val:  100,
				Name: "k1",
			},
			{
				Id:   2,
				Key:  20,
				Val:  200,
				Name: "k2",
			},
		}
		p := NewExpectParam().
			WithMethod("GET").
			WithPath("http://mockhttp.com/list").
			WithKeyValReturn(HttpResponseJson, expectResponse)
		kit.HttpExpect(p, t)
		list, err := srv.httpCli.List(ctx, 10, 20, 40)
		assert.Nil(t, err)
		assert.EqualValues(t, expectResponse, list)
	})
	t.Run("list json response with returns", func(t *testing.T) {
		expectResponse := []*iface.KV{
			{
				Id:   1,
				Key:  10,
				Val:  100,
				Name: "k1",
			},
			{
				Id:   2,
				Key:  20,
				Val:  200,
				Name: "k2",
			},
		}
		p := NewExpectParam().
			WithMethod("GET").
			WithPath("http://mockhttp.com/list").
			WithReturns(expectResponse)
		kit.HttpExpect(p, t)
		list, err := srv.httpCli.List(ctx, 10, 20, 40)
		assert.Nil(t, err)
		assert.EqualValues(t, expectResponse, list)
	})

	t.Run("list json response with func", func(t *testing.T) {
		expectResponse1 := []*iface.KV{
			{
				Id:   1,
				Key:  10,
				Val:  100,
				Name: "k1",
			},
			{
				Id:   2,
				Key:  20,
				Val:  200,
				Name: "k2",
			},
		}
		expectResponse2 := []*iface.KV{
			{
				Id:   11,
				Key:  110,
				Val:  1100,
				Name: "k11",
			},
			{
				Id:   12,
				Key:  120,
				Val:  1200,
				Name: "k12",
			},
		}
		p := NewExpectParam().
			WithMethod("GET").
			WithPath("http://mockhttp.com/list").
			WithHttpResponseFunc(func(req *http.Request) (*http.Response, error) {
				rid := req.URL.Query().Get("relationId")
				var resp []*iface.KV
				switch rid {
				case "1":
					resp = expectResponse1
				case "2":
					resp = expectResponse2
				default:
					resp = make([]*iface.KV, 0)
				}
				return httpmock.NewJsonResponse(http.StatusOK, resp)
			})
		kit.HttpExpect(p, t)
		list1, err := srv.httpCli.List(ctx, 1, 20, 40)
		assert.Nil(t, err)
		assert.EqualValues(t, list1, expectResponse1)
		list2, err := srv.httpCli.List(ctx, 2, 20, 40)
		assert.Nil(t, err)
		assert.EqualValues(t, list2, expectResponse2)
	})

}

func TestMockEs(t *testing.T) {
	expectResponse := []*iface.KV{
		{
			Id:   1,
			Key:  10,
			Val:  100,
			Name: "k1",
		},
		{
			Id:   2,
			Key:  20,
			Val:  200,
			Name: "k2",
		},
		{
			Id:   3,
			Key:  30,
			Val:  300,
			Name: "k3",
		},
	}
	str := `
{
  "hits": {
    "hits": [
      {
        "_source": {
          "id": 1,
          "key": 10,
          "val": 100,
          "name": "k1",
          "@version": "1",
          "@timestamp": "2021-12-12T22:39:55.760Z"
        }
      },
      {
        "_source": {
          "id": 2,
          "key": 20,
          "val": 200,
          "name": "k2",
          "@version": "1",
          "@timestamp": "2021-12-13T22:39:55.760Z"
        }
      },
      {
        "_source": {
          "id": 3,
          "key": 30,
          "val": 300,
          "name": "k3",
          "@version": "1",
          "@timestamp": "2021-12-14T22:39:55.760Z"
        }
      }
    ]
  }
}`
	p := NewExpectParam().
		WithMethod("POST").
		WithPath("http://mockes.com/_all/_search").
		WithReturns(str)
	kit.HttpExpect(p, t)
	list, err := srv.es.List(ctx, 1, 2, 3)
	assert.Nil(t, err)
	assert.Equal(t, expectResponse, list)
}
