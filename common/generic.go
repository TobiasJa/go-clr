package common

func POINTER[T any](v T) *T {
	return &v
}
