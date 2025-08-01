package resultext

import (
	"errors"
	"io"
	"testing"

	. "github.com/pchchv/go-assert"
)

type Struct struct{}

func TestResult(t *testing.T) {
	result := returnOk()
	Equal(t, true, result.IsOk())
	Equal(t, false, result.IsErr())
	Equal(t, true, result.Err() == nil)
	Equal(t, Struct{}, result.Unwrap())

	result = returnErr()
	Equal(t, false, result.IsOk())
	Equal(t, true, result.IsErr())
	Equal(t, false, result.Err() == nil)
	PanicMatches(t, func() {
		result.Unwrap()
	}, "Result.Unwrap(): result is Err")
}

func TestUnwrap(t *testing.T) {
	er := Err[int, error](io.EOF)
	PanicMatches(t, func() { er.Unwrap() }, "Result.Unwrap(): result is Err")

	v := er.UnwrapOr(3)
	Equal(t, 3, v)

	v = er.UnwrapOrElse(func() int { return 2 })
	Equal(t, 2, v)

	v = er.UnwrapOrDefault()
	Equal(t, 0, v)
}

func TestAndXXX(t *testing.T) {
	ok := Ok[int, error](1)
	Equal(t, Ok[int, error](3), ok.And(func(int) int { return 3 }))
	Equal(t, Ok[int, error](3), ok.AndThen(func(int) Result[int, error] { return Ok[int, error](3) }))
	Equal(t, Err[int, error](io.EOF), ok.AndThen(func(int) Result[int, error] { return Err[int, error](io.EOF) }))

	err := Err[int, error](io.EOF)
	Equal(t, Err[int, error](io.EOF), err.And(func(int) int { return 3 }))
	Equal(t, Err[int, error](io.EOF), err.AndThen(func(int) Result[int, error] { return Ok[int, error](3) }))
	Equal(t, Err[int, error](io.EOF), err.AndThen(func(int) Result[int, error] { return Err[int, error](io.ErrUnexpectedEOF) }))
	Equal(t, Err[int, error](io.ErrUnexpectedEOF), ok.AndThen(func(int) Result[int, error] { return Err[int, error](io.ErrUnexpectedEOF) }))
}

func BenchmarkResultOk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if res := returnOk(); res.IsOk() {
			_ = res.Unwrap()
		}
	}
}

func BenchmarkResultErr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if res := returnErr(); res.IsOk() {
			_ = res.Unwrap()
		}
	}
}

func BenchmarkNoResultOk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if res, err := returnOkNoResult(); err != nil {
			_ = res
		}
	}
}

func BenchmarkNoResultErr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if res, err := returnErrNoResult(); err != nil {
			_ = res
		}
	}
}

func returnOk() Result[Struct, error] {
	return Ok[Struct, error](Struct{})
}

func returnErr() Result[Struct, error] {
	return Err[Struct, error](errors.New("bad"))
}

func returnOkNoResult() (Struct, error) {
	return Struct{}, nil
}

func returnErrNoResult() (Struct, error) {
	return Struct{}, errors.New("bad")
}
