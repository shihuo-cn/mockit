package mockit

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis"
	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	jinzhugorm "github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	iogorm "gorm.io/gorm"
	"net/http"
	"reflect"
	"testing"
	"unsafe"
)

// Mockit mock工具集，集成如下
// 1. http （elasticSearch）
// 2. MYSQL (支持多类sql，目前支持了mysql，如需添加请联系)
// 3. Redis
// 4. interface (GRPC)
//  NOTE: 需结合https://github.com/golang/mock使用，生成后需要将实例化返回值改为interface{}，如： NewMockXXX(ctrl *gomock.Controller) interface{}
type Mockit interface {
	// GetInterfaceClient 获取interface的mockClient,获取后 m.(*MockXX)获得实例
	GetInterfaceClient(name string) interface{}
	// MysqlExecExpect mysql 增删改使用
	MysqlExecExpect(ep ExpectParam, tb testing.TB)
	// MysqlQueryExpect mysql 查询使用
	MysqlQueryExpect(ep ExpectParam, tb testing.TB)
	// InterfaceExpect interface mockgen 生成使用
	// NOTE: grpc 推荐mockgen生成后使用此方案作为mock代理
	InterfaceExpect(ep ExpectParam, tb testing.TB)
	// HttpExpect http client 拦截使用
	// NOTE: ES也是此方案
	HttpExpect(ep ExpectParam, tb testing.TB, httpStatus ...int)
	// BeforeTest NOTE:清洁单元测试环境
	BeforeTest()
	// AfterTest NOTE:配合BeforeTest
	AfterTest()

	// InterceptHttpClient 拦截http client
	InterceptHttpClient(client *http.Client)
	// Gorm2DB 获取伪grom2.db做mock用
	// NOTEL: gorm.io/gorm
	Gorm2DB() *iogorm.DB
	// GormDB 获取伪grom.db做mock用
	// NOTE: github.com/jinzhu/gorm
	GormDB() *jinzhugorm.DB
	// RedisAddr 获取伪redis server addr
	RedisAddr() string
}

type mockit struct {
	sqlEx        sqlmock.Sqlmock      // sql excepted
	redisSrv     *miniredis.Miniredis // redis server
	gorm2DB      *iogorm.DB           // gorm.io conn
	gormDB       *jinzhugorm.DB       // jinzhu conn
	ormTag       string               // 实体解析的tag
	clientInitFs []MockClientInitFunc
	iFaceClients map[string]interface{}
	iFaceValue   map[string]reflect.Value
}

type MockClientInitFunc func(ctrl *gomock.Controller) interface{}

// NewWithMockTag 创建mock工具包
// @ormTag: mysql 映射的实体tag,如下可使用"mock"作为tag数据库映射，
// NOTE: gorm中可以带primary key;varchar等，请自定义tag使用
//  type Kit struct {
//    Name string "mock:"name"
//  }
// @fs: 使用mockgen生成的interface mock client
// NOTE: 一般默认生成的是  NewMockDemoInterface(ctrl *gomock.Controller) *MockDemoInterface
// TODO: 需要修改成 NewMockDemoInterface(ctrl *gomock.Controller) interface{}后传入，方可拦截代理
func NewWithMockTag(tag string, fs ...MockClientInitFunc) (Mockit, error) {
	kit := new(mockit)
	kit.iFaceClients = make(map[string]interface{})
	kit.iFaceValue = make(map[string]reflect.Value)
	kit.clientInitFs = fs
	kit.ormTag = tag
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	kit.sqlEx = mock
	// gorm2 io
	gormDB, err := iogorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &iogorm.Config{})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// gorm  jinzhu
	jinZhuDB, err := jinzhugorm.Open("mysql", sqlDB)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	kit.gormDB = jinZhuDB
	kit.gorm2DB = gormDB

	// redis
	kit.redisSrv, err = miniredis.Run()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// 初始化interface mock 并通过内部反射解耦原库使用

	t := log.New()
	t.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	kit.initGoMockInterface(t)
	// 启用httpmock拦截
	httpmock.Activate()
	return kit, nil
}

// New 创建mock工具包
// @fs: 使用mockgen生成的interface mock client
// NOTE: 一般默认生成的是  NewMockDemoInterface(ctrl *gomock.Controller) *MockDemoInterface
// TODO: 需要修改成 NewMockDemoInterface(ctrl *gomock.Controller) interface{}后传入，方可拦截代理
func New(fs ...MockClientInitFunc) (Mockit, error) {
	return NewWithMockTag("mock", fs...)
}

type argsKind uint8

const (
	notSet  argsKind = 0
	normal  argsKind = 1
	byIndex argsKind = 2
)

func (mk *mockit) GetInterfaceClient(name string) interface{} {
	if mk == nil {
		return nil
	}
	return mk.iFaceClients[name]
}

func (mk *mockit) MysqlExecExpect(ep ExpectParam, tb testing.TB) {
	p := ep.(*expectParam)
	if len(p.method) == 0 {
		tb.Fatal("there is no SQL statements in the method filed")
	}
	if len(p.returns) != 2 {
		tb.Fatal("returns requires the addition of lastInsertId and rowsAffected, and the type is int64")
	}
	var num1, num2 int64
	for k, v := range p.returns {
		var numTmp int64
		switch tmp := v.(type) {
		case uint:
			numTmp = int64(tmp)
		case int:
			numTmp = int64(tmp)
		case int64:
			numTmp = int64(tmp)
		case int32:
			numTmp = int64(tmp)
		case uint32:
			numTmp = int64(tmp)
		case uint64:
			numTmp = int64(tmp)
		default:
			tb.Fatal("the returns type must be int of the func mysqlExecExpect")
		}
		if k == 0 {
			num1 = numTmp
		} else {
			num2 = numTmp
		}
	}
	var args []driver.Value
	for _, v := range p.args {
		args = append(args, v)
	}
	mk.sqlEx.ExpectBegin()
	mk.sqlEx.ExpectExec(p.method).WithArgs(args...).WillReturnResult(sqlmock.NewResult(num1, num2))
	mk.sqlEx.ExpectCommit()
}

func (mk *mockit) MysqlQueryExpect(ep ExpectParam, tb testing.TB) {
	var (
		rows *sqlmock.Rows
		err  error
	)
	p := ep.(*expectParam)
	if len(p.method) == 0 {
		tb.Fatal("there is no SQL statements in the method filed")
	}
	if len(p.key) > 0 {
		rows = sqlmock.NewRows([]string{p.key}).AddRow(p.val)
	} else {
		if len(p.returns) != 1 {
			tb.Fatal("the query must return one result")

		}
		rows, err = sqlmock.NewRowsFromInterface(p.returns[0], mk.ormTag)
	}
	if err != nil {
		tb.Fatalf("new rows failed:%s", err)
	}
	var args []driver.Value
	for _, v := range p.args {
		args = append(args, v)
	}
	mk.sqlEx.ExpectQuery(p.method).WithArgs(args...).WillReturnRows(rows)
}

func (mk *mockit) InterfaceExpect(ep ExpectParam, tb testing.TB) {
	p := ep.(*expectParam)
	if len(p.path) == 0 || len(p.method) == 0 {
		tb.Fatalf("both path and method are required")
	}
	cli, exists := mk.iFaceValue[p.path]
	if !exists {
		tb.Fatalf("interface mock client %s not exists", p.path)
	}
	method := cli.Elem().MethodByName(p.method)
	if !method.IsValid() {
		tb.Fatalf("path:%s method:%s not exists", p.path, p.method)
	}
	numIn := method.Type().NumIn()
	inputs := make([]reflect.Value, numIn)

	switch p.argsKind {
	case normal:
		for i, arg := range p.args {
			inputs[i] = reflect.ValueOf(arg)
		}
	case byIndex:
		for i := 0; i < numIn; i++ {
			if arg, exists := p.idxArgs[i]; exists {
				inputs[i] = reflect.ValueOf(gomock.Eq(arg))
			} else {
				inputs[i] = reflect.ValueOf(gomock.Any())
			}
		}
	default:
		for i := 0; i < numIn; i++ {
			inputs[i] = reflect.ValueOf(gomock.Any())
		}
	}
	outputs := method.Call(inputs)
	if len(outputs) != 1 {
		tb.Fatal("method returns not match")
	}
	methodReturn := outputs[0].MethodByName("Return")
	originalMethod, exists := reflect.TypeOf(mk.iFaceClients[p.path]).MethodByName(p.method)
	if !exists {
		tb.Fatal("mock client is not formatted correctly")
	}
	var types []reflect.Type
	for i := 0; i < originalMethod.Type.NumOut(); i++ {
		types = append(types, originalMethod.Type.Out(i))
	}
	rs := reflect.New(reflect.TypeOf([]interface{}{})).Elem()
	for i, v := range p.returns {
		val := reflect.ValueOf(v)
		if val.IsValid() {
			rs = reflect.Append(rs, val)
		} else {
			rs = reflect.Append(rs, reflect.Zero(types[i]))
		}
	}
	methodReturn.CallSlice([]reflect.Value{rs})
}

func (mk *mockit) HttpExpect(ep ExpectParam, tb testing.TB, httpStatus ...int) {
	var (
		resp httpmock.Responder
		err  error
	)
	p := ep.(*expectParam)
	if len(p.path) == 0 || len(p.method) == 0 {
		tb.Fatalf("both path and method are required")
	}
	status := http.StatusOK
	if len(httpStatus) > 0 {
		status = httpStatus[0]
	}
	switch p.key {
	case HttpResponseFunc:
		if p.responseHandler == nil {
			tb.Fatal("responseHandler is nil")
		}
		resp = p.responseHandler
	case HttpResponseString:
		str, ok := p.val.(string)
		if !ok {
			tb.Fatal("return val must be string type")
		}
		resp = httpmock.NewStringResponder(status, str)
	case HttpResponseJson:
		resp, err = httpmock.NewJsonResponder(status, p.val)
		if err != nil {
			tb.Fatalf("NewJsonResponder failed:%s", err)
		}
	default:
		if len(p.returns) != 1 {
			tb.Fatal("response nil, please set response via any one of WithKeyValReturn/WithHttpResponseHandler/WithReturns")
		}
		arg := p.returns[0]
		switch argTmp := arg.(type) {
		case string:
			resp = httpmock.NewStringResponder(status, argTmp)
		default:
			resp, err = httpmock.NewJsonResponder(status, arg)
			if err != nil {
				tb.Fatalf("NewJsonResponder failed:%s", err)
			}
		}
	}
	httpmock.RegisterResponder(
		p.method,
		p.path,
		resp,
	)
}

// 初始化interface mock 并通过内部反射解耦原库使用
func (mk *mockit) initGoMockInterface(t gomock.TestReporter) {
	ctrl := gomock.NewController(t)
	for _, f := range mk.clientInitFs {
		cli := f(ctrl)
		val := reflect.ValueOf(cli)
		typ := val.Elem().Type()
		if typ.Kind() != reflect.Struct {
			panic("mock client must gen form mockgen")
		}
		cliName := typ.Name()
		mk.iFaceClients[cliName] = cli
		expect := val.Elem().FieldByName("recorder")
		unsafeExpect := reflect.NewAt(expect.Type(), unsafe.Pointer(expect.UnsafeAddr()))
		mk.iFaceValue[cliName] = unsafeExpect
	}
}

func (mk *mockit) BeforeTest() {
	httpmock.Activate()
}

func (mk *mockit) AfterTest() {
	httpmock.DeactivateAndReset()
}

func (mk *mockit) InterceptHttpClient(client *http.Client) {
	httpmock.ActivateNonDefault(client)
}

func (mk *mockit) Gorm2DB() *iogorm.DB {
	if mk == nil {
		return nil
	}
	return mk.gorm2DB
}

func (mk *mockit) RedisAddr() string {
	if mk == nil {
		return ""
	}
	return mk.redisSrv.Addr()
}

func (mk *mockit) GormDB() *jinzhugorm.DB {
	if mk == nil {
		return nil
	}
	return mk.gormDB
}
