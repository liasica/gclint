// Package namedreturn contains analyzer fixtures for named return checks.
package namedreturn

type payload struct{}

func loadPayload() (*payload, error) {
	return &payload{}, nil
}

func badNamedReturn() (payloadValue *payload, err error) {
	payloadValue, err = loadPayload()
	if err != nil {
		return nil, err // want "named return values were assigned before this explicit return"
	}

	return payloadValue, nil // want "named return values were assigned before this explicit return"
}

func goodNamedReturn() (payloadValue *payload, err error) {
	payloadValue, err = loadPayload()
	if err != nil {
		return
	}

	return
}
