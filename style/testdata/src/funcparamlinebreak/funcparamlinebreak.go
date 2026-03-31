package funcparamlinebreak

type payload struct{}

// 通过：无参数
func noParams() {}

// 通过：1 个参数写在一行
func oneParamInline(a int) {}

// 报错：1 个参数但换行了
func oneParamMultiline( // want "function has 1 parameters \\(threshold 5\\), parameter list must not span multiple lines"
	a int,
) {
}

// 通过：2 个参数写在一行
func twoParamsInline(a int, b string) {}

// 报错：2 个参数但换行了
func twoParamsMultiline( // want "function has 2 parameters \\(threshold 5\\), parameter list must not span multiple lines"
	a int,
	b string,
) {
}

// 通过：3 个参数写在一行
func threeParamsInline(a int, b string, c bool) {}

// 报错：3 个参数但换行了
func threeParamsMultiline( // want "function has 3 parameters \\(threshold 5\\), parameter list must not span multiple lines"
	a int,
	b string,
	c bool,
) {
}

// 通过：4 个参数写在一行
func fourParamsInline(a int, b string, c bool, d float64) {}

// 报错：4 个参数但换行了
func fourParamsMultiline( // want "function has 4 parameters \\(threshold 5\\), parameter list must not span multiple lines"
	a int,
	b string,
	c bool,
	d float64,
) {
}

// 通过：5 个参数（等于阈值）换行
func fiveParamsMultiline(
	a int,
	b string,
	c bool,
	d float64,
	e []byte,
) {
}

// 通过：6 个参数（超过阈值）换行
func sixParamsMultiline(
	a int,
	b string,
	c bool,
	d float64,
	e []byte,
	f error,
) {
}

// 通过：5 个参数写在一行也行
func fiveParamsInline(a int, b string, c bool, d float64, e []byte) {}

// 通过：方法接收者不影响参数计数
type receiver struct{}

func (r *receiver) methodTwoParamsInline(a int, b string) {}

// 报错：方法接收者不影响参数计数，2 个参数换行了
func (r *receiver) methodTwoParamsMultiline( // want "function has 2 parameters \\(threshold 5\\), parameter list must not span multiple lines"
	a int,
	b string,
) {
}

// 通过：匿名函数，参数在一行
var inlineFuncLit = func(a int, b string) {}

// 报错：匿名函数，少于 5 个参数但换行
var multilineFuncLit = func( // want "function has 2 parameters \\(threshold 5\\), parameter list must not span multiple lines"
	a int,
	b string,
) {
}

// 通过：分组参数 a, b int 等于 2 个参数，在一行
func groupedParamsInline(a, b int) {}

// 报错：分组参数 a, b int 等于 2 个参数，但换行了
func groupedParamsMultiline( // want "function has 2 parameters \\(threshold 5\\), parameter list must not span multiple lines"
	a, b int,
) {
}

// 通过：匿名参数（接口方法签名中常见），在一行
type exampleInterface interface {
	Method(int, string, bool)
}

func usePayload(_ *payload) {}
