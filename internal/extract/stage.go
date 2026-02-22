package extract

// DetectStage identifies the lifecycle stage from tokenized input.
func (e *Extractor) DetectStage(tokens []string) string {
	for _, t := range tokens {
		if stage, ok := e.stageKeywords[t]; ok {
			return stage
		}
	}
	return ""
}
