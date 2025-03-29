package persistence

func toAnySlice[T any](values []T) []any {
	docs := make([]any, len(values))
	for i, value := range values {
		docs[i] = value
	}
	return docs
}
