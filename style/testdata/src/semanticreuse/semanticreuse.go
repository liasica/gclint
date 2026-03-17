// Package semanticreuse contains analyzer fixtures for semantic variable reuse checks.
package semanticreuse

func marshalUser() ([]byte, error) {
	return []byte("user"), nil
}

func marshalAudit() ([]byte, error) {
	return []byte("audit"), nil
}

func normalizeUserJSON(userJSON []byte) ([]byte, error) {
	return userJSON, nil
}

func badSemanticReuse() error {
	var userJSON []byte
	var err error

	userJSON, err = marshalUser()
	if err != nil {
		return err
	}

	userJSON, err = marshalAudit() // want "variable \"userJSON\" was already established"
	if err != nil {
		return err
	}

	_ = userJSON

	return nil
}

func goodSemanticReuse() error {
	var userJSON []byte
	var err error

	userJSON, err = marshalUser()
	if err != nil {
		return err
	}

	userJSON, err = normalizeUserJSON(userJSON)
	if err != nil {
		return err
	}

	_ = userJSON

	return nil
}
