package funcparamlinebreak

type payload struct{}

// Pass: no parameters
func noParams() {}

// Pass: 1 parameter on a single line
func oneParamInline(a int) {}

// Fail: 1 parameter but with line break
func oneParamMultiline( // want "function has 1 parameters \\(threshold 5\\), parameter list must not span multiple lines"
	a int,
) {
}

// Pass: 2 parameters on a single line
func twoParamsInline(a int, b string) {}

// Fail: 2 parameters but with line break
func twoParamsMultiline( // want "function has 2 parameters \\(threshold 5\\), parameter list must not span multiple lines"
	a int,
	b string,
) {
}

// Pass: 3 parameters on a single line
func threeParamsInline(a int, b string, c bool) {}

// Fail: 3 parameters but with line break
func threeParamsMultiline( // want "function has 3 parameters \\(threshold 5\\), parameter list must not span multiple lines"
	a int,
	b string,
	c bool,
) {
}

// Pass: 4 parameters on a single line
func fourParamsInline(a int, b string, c bool, d float64) {}

// Fail: 4 parameters but with line break
func fourParamsMultiline( // want "function has 4 parameters \\(threshold 5\\), parameter list must not span multiple lines"
	a int,
	b string,
	c bool,
	d float64,
) {
}

// Pass: 5 parameters (equals threshold) with line break
func fiveParamsMultiline(
	a int,
	b string,
	c bool,
	d float64,
	e []byte,
) {
}

// Pass: 6 parameters (exceeds threshold) with line break
func sixParamsMultiline(
	a int,
	b string,
	c bool,
	d float64,
	e []byte,
	f error,
) {
}

// Pass: 5 parameters on a single line is also fine
func fiveParamsInline(a int, b string, c bool, d float64, e []byte) {}

// Pass: method receiver does not affect parameter count
type receiver struct{}

func (r *receiver) methodTwoParamsInline(a int, b string) {}

// Fail: method receiver does not affect parameter count, 2 parameters with line break
func (r *receiver) methodTwoParamsMultiline( // want "function has 2 parameters \\(threshold 5\\), parameter list must not span multiple lines"
	a int,
	b string,
) {
}

// Pass: anonymous function, parameters on a single line
var inlineFuncLit = func(a int, b string) {}

// Fail: anonymous function, fewer than 5 parameters but with line break
var multilineFuncLit = func( // want "function has 2 parameters \\(threshold 5\\), parameter list must not span multiple lines"
	a int,
	b string,
) {
}

// Pass: grouped parameters a, b int count as 2, on a single line
func groupedParamsInline(a, b int) {}

// Fail: grouped parameters a, b int count as 2, but with line break
func groupedParamsMultiline( // want "function has 2 parameters \\(threshold 5\\), parameter list must not span multiple lines"
	a, b int,
) {
}

// Pass: anonymous parameters (common in interface method signatures), on a single line
type exampleInterface interface {
	Method(int, string, bool)
}

func usePayload(_ *payload) {}
