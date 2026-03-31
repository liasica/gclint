package errvarname

import "errors"

var errSentinel = errors.New("sentinel")

func returnError() error {
	return errSentinel
}

func returnValueAndError() (int, error) {
	return 0, nil
}

func returnThreeWithError() (string, int, error) {
	return "", 0, nil
}

func returnNoError() int {
	return 0
}

func returnTwoNoError() (int, string) {
	return 0, ""
}

// 通过：用 err 接收 error 返回值
func goodSingleError() {
	err := returnError()
	_ = err
}

// 通过：用 err 接收多返回值中的 error
func goodMultiReturn() {
	_, err := returnValueAndError()
	_ = err
}

// 通过：用 _ 忽略 error
func goodBlankIdentifier() {
	_ = returnError()
}

// 通过：用 _ 忽略多返回值中的 error
func goodBlankMultiReturn() {
	_, _ = returnValueAndError()
}

// 通过：函数不返回 error
func goodNoError() {
	_ = returnNoError()
}

// 通过：多返回值但最后一个不是 error
func goodTwoNoError() {
	_, _ = returnTwoNoError()
}

// 报错：用 e 接收 error
func badSingleVarE() {
	e := returnError() // want `error return must be received by "err", not "e"`
	_ = e
}

// 报错：用 myErr 接收 error
func badSingleVarMyErr() {
	myErr := returnError() // want `error return must be received by "err", not "myErr"`
	_ = myErr
}

// 报错：多返回值中用 fetchErr 接收 error
func badMultiReturnCustomName() {
	_, fetchErr := returnValueAndError() // want `error return must be received by "err", not "fetchErr"`
	_ = fetchErr
}

// 报错：三返回值中用 loadErr 接收 error
func badThreeReturnCustomName() {
	_, _, loadErr := returnThreeWithError() // want `error return must be received by "err", not "loadErr"`
	_ = loadErr
}

// 通过：赋值语句（=）用 err 接收
func goodAssignExisting() {
	var err error

	err = returnError()
	_ = err
}

// 报错：赋值语句（=）用其他名接收
func badAssignExisting() {
	var e error

	e = returnError() // want `error return must be received by "err", not "e"`
	_ = e
}

// 通过：多返回值赋值语句用 err 接收
func goodAssignMultiReturn() {
	var n int
	var err error

	n, err = returnValueAndError()
	_ = n
	_ = err
}

// 报错：多返回值赋值语句用其他名接收
func badAssignMultiReturn() {
	var n int
	var fetchErr error

	n, fetchErr = returnValueAndError() // want `error return must be received by "err", not "fetchErr"`
	_ = n
	_ = fetchErr
}
