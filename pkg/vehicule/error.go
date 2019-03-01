package vehicule

// A NotFoundError is an error that represents that no vehicule was found.
type NotFoundError struct {
	msg string
}

func (e NotFoundError) Error() string {
	return e.msg
}

// A WrongUserError is an error that represents that a vehicule is being
// added to another user then the one authenticated
type WrongUserError struct {
	msg string
}

func (e WrongUserError) Error() string {
	return e.msg
}
