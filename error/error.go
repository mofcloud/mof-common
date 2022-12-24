package perr

import "github.com/pkg/errors"

func Wrap(err error) error {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	_, ok := err.(stackTracer)

	if !ok {
		err = errors.WithStack(err)
	}

	return err
}
