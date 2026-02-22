package extract

// DetectAudience identifies the target audience from tokenized input.
func (e *Extractor) DetectAudience(tokens []string) string {
	for _, t := range tokens {
		if val, ok := e.audienceKeywords[t]; ok {
			return val
		}
	}
	return ""
}

// DetectDepth identifies the response depth from tokenized input.
func (e *Extractor) DetectDepth(tokens []string) string {
	for _, t := range tokens {
		if val, ok := e.depthKeywords[t]; ok {
			return val
		}
	}
	return ""
}

// DetectStyle identifies the response style from tokenized input.
func (e *Extractor) DetectStyle(tokens []string) string {
	for _, t := range tokens {
		if val, ok := e.styleKeywords[t]; ok {
			return val
		}
	}
	return ""
}

// DetectFormat identifies the desired output format from tokenized input.
func (e *Extractor) DetectFormat(tokens []string) string {
	for _, t := range tokens {
		if val, ok := e.formatKeywords[t]; ok {
			return val
		}
	}
	return ""
}
