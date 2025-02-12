package types

const (
	EmailAlreadyExists = "User with this email already exists"
)

type CreateUserError struct {
	Msg   string // Description of the error
	Field string // Which field in the table might have thrown the error (if any)
}

func (p CreateUserError) Error() string {
	return p.Msg
}

func GetCreateUserError(msg string, field string) error {
	return CreateUserError{
		Msg:   msg,
		Field: field,
	}
}
