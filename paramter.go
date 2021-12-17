package mockit

import "net/http"

// ExpectParam mockit ioc/控制反转 参数
type ExpectParam interface {
	// WithPath 路径类
	// interface : mock client name
	// http:  url "http://mock.com/get"
	WithPath(path string) ExpectParam
	// WithMethod 动作/方法类
	// http :  POST GET DELETE等
	// mysql: "SELECT (.+) FROM DEMO"  "INSERT INTO$"等
	WithMethod(method string) ExpectParam
	// WithArgs 入参
	WithArgs(args ...interface{}) ExpectParam
	// WithArgsByIndex 下标入参
	// NOTE: 目前针对interface封装,不覆盖的会走Any（忽略匹配）逻辑
	WithArgsByIndex(index int, value interface{}) ExpectParam
	// WithHttpResponseFunc 设置Http response
	// http 设置可变response用
	WithHttpResponseFunc(f func(req *http.Request) (*http.Response, error)) ExpectParam
	// WithReturns 设置预期返回值
	// http: 限制1个，默认匹配interface类型，匹配string时，response 返回为string,其他情况response使用json
	// mysql: 支持 struct/slice/array 默认使用json tag,如果mockit设置tag为其他，如"mock",可参见单元测试
	// 其他: 全量
	WithReturns(m ...interface{}) ExpectParam
	// WithKeyValReturn 设置单个预期返回值
	// mysql：如查询计数时，可使用 "COUNT", 20   => 代表返回20，注意这个`COUNT`必须和查询语句的大小写相同
	// http: 支持 string json 2种key,区分response返回的类型
	WithKeyValReturn(key ExpectKey, val interface{}) ExpectParam
}

// ExpectKey custom key
// In the case of mysql： `json` `string` `func`  choose excepted response
type ExpectKey = string

const (
	// for http response encode
	HttpResponseFunc   ExpectKey = "func"   // first
	HttpResponseString ExpectKey = "string" // second
	HttpResponseJson   ExpectKey = "json"   // third

)

type expectParam struct {
	// 动作/方法
	// http :  POST GET DELETE等
	// mysql: "SELECT (.+) FROM DEMO"  "INSERT INTO$"等
	method string
	// interface : mock client name
	// http:  url "http://mock.com/get"
	path string

	// default:
	// 0: not set
	// 1: WithArgs
	// 2: WithArgsByIndex
	argsKind argsKind
	args     []interface{}
	idxArgs  map[int]interface{}

	// return
	returns         []interface{}
	key             ExpectKey
	val             interface{}
	responseHandler func(req *http.Request) (*http.Response, error)
}

func NewExpectParam() ExpectParam {
	return new(expectParam)
}
func (p *expectParam) WithPath(path string) ExpectParam {
	p.path = path
	return p
}
func (p *expectParam) WithMethod(method string) ExpectParam {
	p.method = method
	return p
}
func (p *expectParam) WithArgs(args ...interface{}) ExpectParam {
	if !(p.argsKind == normal || p.argsKind == notSet) {
		return p
	}
	p.args = args
	p.argsKind = normal
	return p
}

func (p *expectParam) WithArgsByIndex(index int, value interface{}) ExpectParam {
	if !(p.argsKind == byIndex || p.argsKind == notSet) {
		return p
	}
	if len(p.idxArgs) == 0 {
		p.idxArgs = make(map[int]interface{})
	}
	p.argsKind = byIndex
	p.idxArgs[index] = value
	return p
}

func (p *expectParam) WithHttpResponseFunc(f func(req *http.Request) (*http.Response, error)) ExpectParam {
	p.responseHandler = f
	p.key = HttpResponseFunc
	return p
}
func (p *expectParam) WithReturns(m ...interface{}) ExpectParam {
	if p == nil {
		p = new(expectParam)
	}
	p.returns = m
	return p
}
func (p *expectParam) WithKeyValReturn(key ExpectKey, val interface{}) ExpectParam {
	p.key = key
	p.val = val
	return p
}
func (p *expectParam) WithHttpResponseHandler(f func(req *http.Request) (*http.Response, error)) ExpectParam {
	p.key = HttpResponseFunc
	p.responseHandler = f
	return p
}
