package melon

import "errors"

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicOnErrorIgnore(err error, ignoreErrors ...error) {
	for _, ignoreErr := range ignoreErrors {
		if errors.Is(err, ignoreErr) {
			return
		}
	}
	PanicOnError(err)
}

func PanicOnFalse(ok bool, supplier func() error) {
	if ok {
		return
	}
	panic(supplier())
}

func IsErrorIn(err error, errorList ...error) bool {
	for _, e := range errorList {
		if errors.Is(e, err) {
			return true
		}
	}
	return false
}
