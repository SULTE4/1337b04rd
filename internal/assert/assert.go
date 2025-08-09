package assert

func Equal[T comparable](a, b T) bool {
	return a == b
}
