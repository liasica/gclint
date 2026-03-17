// Package redeclare contains analyzer fixtures for same-scope short redeclaration checks.
package redeclare

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

func innerScopeShadowingIsAllowed() error {
	firstValue, err := readFirstValue()
	if err != nil {
		return err
	}

	if firstValue != "" {
		firstValue, secondValue := readPair()
		if firstValue == "" || secondValue == "" {
			return nil
		}
	}

	return nil
}
