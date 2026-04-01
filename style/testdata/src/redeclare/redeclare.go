package redeclare

import (
	"encoding/json"
	"errors"
	"fmt"
)

type workerResponse struct{}

type result struct {
	resp *workerResponse
	err  error
}

type stdoutReader interface {
	ReadBytes(delimiter byte) ([]byte, error)
}

type workerClient struct {
	stdout stdoutReader
}

func readFirstValue() (string, error) {
	return "first", nil
}

func readSecondValue() (string, error) {
	return "second", nil
}

func readPair() (string, string) {
	return "left", "right"
}

func badRedeclare() error {
	firstValue, err := readFirstValue()
	if err != nil {
		return err
	}

	secondValue, firstValue := readPair() // want "existing variable \"firstValue\" must not be reused in short variable declaration"
	if secondValue == "" || firstValue == "" {
		return nil
	}

	thirdValue, err := readSecondValue() // want "existing variable \"err\" must not be reused in short variable declaration"
	if thirdValue == "" {
		return nil
	}

	_, err = readSecondValue()
	if err != nil {
		return err
	}

	return nil
}

func goodRedeclare() error {
	firstValue, err := readFirstValue()
	if err != nil {
		return err
	}

	var secondValue string
	secondValue, firstValue = readPair()

	if secondValue == "" || firstValue == "" {
		return nil
	}

	return nil
}

func innerScopeShadowingIsForbidden() error {
	firstValue, err := readFirstValue()
	if err != nil {
		return err
	}

	if firstValue != "" {
		firstValue, secondValue := readPair() // want "existing variable \"firstValue\" must not be reused in short variable declaration"
		if firstValue == "" || secondValue == "" {
			return nil
		}
	}

	return nil
}

func failValidation() error {
	return nil
}

func innerScopeErrShadowingIsForbidden() error {
	firstValue, err := readFirstValue()
	if err != nil {
		return err
	}

	if err := failValidation(); err != nil { // want "existing variable \"err\" must not be reused in short variable declaration"
		return err
	}

	_ = firstValue

	return nil
}

func rangeShadowingIsForbidden(values []string) []string {
	firstValue := "seed"

	for _, firstValue := range values { // want "existing variable \"firstValue\" must not be reused in short variable declaration"
		if firstValue == "" {
			continue
		}
	}

	_ = firstValue

	return values
}

func nestedFunctionShadowingIsAllowed(values []string) []string {
	firstValue := "seed"

	go func() {
		for _, firstValue := range values {
			if firstValue == "" {
				return
			}
		}
	}()

	_ = firstValue

	return values
}

func goroutineErrShadowingIsForbidden(c *workerClient) <-chan result {
	responseChannel := make(chan result, 1)

	go func() {
		line, err := c.stdout.ReadBytes('\n')
		if err != nil {
			responseChannel <- result{err: err}
			return
		}

		var response workerResponse
		if err := json.Unmarshal(line, &response); err != nil { // want "existing variable \"err\" must not be reused in short variable declaration"
			responseChannel <- result{err: fmt.Errorf("decode semantic worker response: %w", err)}
			return
		}

		responseChannel <- result{resp: &response}
	}()

	return responseChannel
}

func withTransaction(fn func() error) error {
	return fn()
}

// closureErrInSameShortVarDecl 复现 bug:
// 外层 err := f(func() { err := ... }) 闭包内的 err 不应触发 redeclare
func closureErrInSameShortVarDecl() {
	err := withTransaction(func() error {
		val, err := readFirstValue()
		if err != nil {
			return err
		}
		_ = val
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// closureErrInSameShortVarDeclMultiReturn 多返回值版本
func closureErrInSameShortVarDeclMultiReturn() {
	result, err := readFirstValue()
	if err != nil {
		panic(err)
	}

	_ = withTransaction(func() error {
		val, err := readSecondValue()
		if err != nil {
			return err
		}
		_ = val
		return nil
	})

	_ = result
}

// closureWithInnerRedeclare 闭包内部自身的 redeclare 仍应报错
func closureWithInnerRedeclare() {
	err := withTransaction(func() error {
		val, err := readFirstValue()
		if err != nil {
			return err
		}

		val2, err := readSecondValue() // want "existing variable \"err\" must not be reused in short variable declaration"
		_, _ = val, val2
		return nil
	})

	_ = errors.New("use errors import")

	if err != nil {
		panic(err)
	}
}
