package errshortdecl

func readFirst() (string, error) {
	return "first", nil
}

func readSecond() (string, error) {
	return "second", nil
}

func badShortDeclaration() error {
	firstValue, err := readFirst()
	if err != nil {
		return err
	}

	secondValue, err := readSecond() // want "existing err must not be reused in short variable declaration"
	if err != nil {
		return err
	}

	_ = firstValue
	_ = secondValue

	return nil
}

func innerScopeShortDeclarationIsAllowed() error {
	firstValue, err := readFirst()
	if err != nil {
		return err
	}

	if _, err := readSecond(); err != nil {
		return err
	}

	_ = firstValue

	return nil
}
