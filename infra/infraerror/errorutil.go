package infraerror

import "errors"

func CountCauses(err error) int {
	i := 0
	for err != nil {
		unwrapErr := errors.Unwrap(err)
		if unwrapErr == nil {
			break
		}
		err = unwrapErr
		i++
	}
	return i
}
