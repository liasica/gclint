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

func shortlistSemanticCandidates(target string, candidates []string, relatedCandidateLimit int) []string {
	if target == "" {
		return candidates
	}

	if len(candidates) > relatedCandidateLimit {
		return candidates[:relatedCandidateLimit]
	}

	return candidates
}

func loadFallbackCandidates() []string {
	return []string{"fallback"}
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

func goodDerivedSourceReuse(target string, candidates []string, relatedCandidateLimit int) []string {
	shortlist := shortlistSemanticCandidates(target, candidates, relatedCandidateLimit)
	shortlist = candidates

	if len(shortlist) > relatedCandidateLimit {
		shortlist = shortlist[:relatedCandidateLimit]
	}

	return shortlist
}

func badDerivedSourceReuse(target string, candidates []string, relatedCandidateLimit int) []string {
	shortlist := shortlistSemanticCandidates(target, candidates, relatedCandidateLimit)
	fallbackCandidates := loadFallbackCandidates()

	shortlist = fallbackCandidates // want "variable \"shortlist\" was already established"

	if len(shortlist) > relatedCandidateLimit {
		shortlist = shortlist[:relatedCandidateLimit]
	}

	return shortlist
}
