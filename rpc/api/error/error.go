package api

type EmptyError struct{}

func (e *EmptyError) Error() string {
	return "empty user"
}

