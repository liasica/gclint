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

// Pass: receive error return with err
func goodSingleError() {
	err := returnError()
	_ = err
}

// Pass: receive error from multi-return with err
func goodMultiReturn() {
	_, err := returnValueAndError()
	_ = err
}

// Pass: ignore error with _
func goodBlankIdentifier() {
	_ = returnError()
}

// Pass: ignore error from multi-return with _
func goodBlankMultiReturn() {
	_, _ = returnValueAndError()
}

// Pass: function does not return error
func goodNoError() {
	_ = returnNoError()
}

// Pass: multi-return but last value is not error
func goodTwoNoError() {
	_, _ = returnTwoNoError()
}

// Fail: receive error with e
func badSingleVarE() {
	e := returnError() // want `error return must be received by "err", not "e"`
	_ = e
}

// Fail: receive error with myErr
func badSingleVarMyErr() {
	myErr := returnError() // want `error return must be received by "err", not "myErr"`
	_ = myErr
}

// Fail: receive error from multi-return with fetchErr
func badMultiReturnCustomName() {
	_, fetchErr := returnValueAndError() // want `error return must be received by "err", not "fetchErr"`
	_ = fetchErr
}

// Fail: receive error from triple-return with loadErr
func badThreeReturnCustomName() {
	_, _, loadErr := returnThreeWithError() // want `error return must be received by "err", not "loadErr"`
	_ = loadErr
}

// Pass: assignment (=) receives error with err
func goodAssignExisting() {
	var err error

	err = returnError()
	_ = err
}

// Fail: assignment (=) receives error with non-err name
func badAssignExisting() {
	var e error

	e = returnError() // want `error return must be received by "err", not "e"`
	_ = e
}

// Pass: multi-return assignment receives error with err
func goodAssignMultiReturn() {
	var n int
	var err error

	n, err = returnValueAndError()
	_ = n
	_ = err
}

// Fail: multi-return assignment receives error with non-err name
func badAssignMultiReturn() {
	var n int
	var fetchErr error

	n, fetchErr = returnValueAndError() // want `error return must be received by "err", not "fetchErr"`
	_ = n
	_ = fetchErr
}
