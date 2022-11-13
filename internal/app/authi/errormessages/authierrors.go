package errormessages

import "fmt"

var (
	UserAlreadyExists = fmt.Errorf("user already exists")
)
