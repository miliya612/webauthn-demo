package errUtil

type ErrTodoNotFound struct{}

func (e ErrTodoNotFound) Error() string {
	return "todo is not found"
}
