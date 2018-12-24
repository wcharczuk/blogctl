package web

// ResultProvider is the provider interface for results.
type ResultProvider interface {
	InternalError(err error) Result
	BadRequest(err error) Result
	NotFound() Result
	NotAuthorized() Result
	Status(statusCode int, result ...interface{}) Result
}

// ResultOrDefault returns a result or a default.
func ResultOrDefault(defaultResult interface{}, result ...interface{}) interface{} {
	if len(result) > 0 {
		return result[0]
	}
	return defaultResult
}
