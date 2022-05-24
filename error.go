package melon

import (
	"errors"
	"fmt"
)

func IsErrorIn(err error, errorList ...error) bool {
	for _, e := range errorList {
		if errors.Is(e, err) {
			return true
		}
	}
	return false
}

func MustOK(ok bool, supplier func() error) {
	if ok {
		return
	}
	panic(supplier())
}

func MustOKStr(ok bool, str string) {
	if ok {
		return
	}
	panic(errors.New(str))
}

func MustIgnore(err error, ignoreErrors ...error) {
	if IsErrorIn(err, ignoreErrors...) {
		return
	}
	Must(err)
}

func MustWrapStr(err error, wrapStr string) error {
	if err == nil {
		return nil
	}
	panic(fmt.Errorf("%s: %w", wrapStr, err))
}

func MustWrap(err error, supplier func(err error) error) error {
	if err == nil {
		return nil
	}
	panic(supplier(err))
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Must2[T any](t T, err error) T {
	Must(err)
	return t
}

func Must3[T1, T2 any](t1 T1, t2 T2, err error) (T1, T2) {
	Must(err)
	return t1, t2
}

func Must4[T1, T2, T3 any](t1 T1, t2 T2, t3 T3, err error) (T1, T2, T3) {
	Must(err)
	return t1, t2, t3
}

func Must5[T1, T2, T3, T4 any](t1 T1, t2 T2, t3 T3, t4 T4, err error) (T1, T2, T3, T4) {
	Must(err)
	return t1, t2, t3, t4
}

func Must6[T1, T2, T3, T4, T5 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, err error) (T1, T2, T3, T4, T5) {
	Must(err)
	return t1, t2, t3, t4, t5
}
