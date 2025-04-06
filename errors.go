package iterate

import (
	"fmt"
)

type ErrCallback struct {
	Coordinates []int
	Err         error
}

func (e ErrCallback) Error() string {
	return fmt.Sprintf("callback returned error at %v: %v", e.Coordinates, e.Err)
}

func (e ErrCallback) Unwrap() error {
	return e.Err
}
