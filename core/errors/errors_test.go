package errors_test

import (
	gerr "errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wyy-go/wzo/core/errors"
)

func Test_NilError(t *testing.T) {
	var err *errors.Error

	err = err.TakeOption()
	require.Equal(t, err.Code, int32(0))
	require.Equal(t, err.Message, "")
	require.Equal(t, err.Metadata, map[string]string(nil))
	//require.Equal(t, err.Error(), "<nil>")
}

func Test_Error(t *testing.T) {
	err := errors.New(400, "请求参数错误", "")
	require.Error(t, err)
	//require.Equal(t, err.Error(), "请求参数错误")

	innerErr := gerr.New("内部错误1")
	err = err.TakeOption(errors.WithDetail(innerErr.Error()), errors.WithMetadata("k1", "v1"))

	require.Equal(t, err.Code, int32(400))
	require.Equal(t, err.Message, "请求参数错误")
	require.Equal(t, err.Metadata, map[string]string{"k1": "v1"})
	//require.Equal(t, err.Error(), "请求参数错误: 内部错误1")

	err = err.TakeOption(errors.WithDetail("内部错误2"))
	//require.Equal(t, err.Error(), "请求参数错误: 内部错误2")

	err = err.TakeOption(errors.WithDetail("内部错误3"))
	//require.Equal(t, err.Error(), "请求参数错误: 内部错误3")

	err = err.TakeOption(errors.WithMessage("另一个错误"))
	//require.Equal(t, err.Error(), "另一个错误: 内部错误3")

	err = err.TakeOption(errors.WithMessagef("另一个错误(%v)", "设备号111"))
	//require.Equal(t, err.Error(), "另一个错误(设备号111): 内部错误3")
}

type testError struct {
	s string
}

func (e *testError) Error() string { return e.s }

func newTestError(s string) *testError {
	return &testError{s: s}
}

func Test_Error_Unwrap(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		var err *errors.Error

		gotErr := new(testError)
		ok := gerr.As(err, &gotErr)
		require.False(t, ok)
	})
	t.Run("should not found", func(t *testing.T) {
		err1 := errors.Newf(400, "请求参数错误", "")

		gotErr := new(testError)
		ok := gerr.As(err1, &gotErr)
		require.False(t, ok)
	})
	t.Run("Unwrap", func(t *testing.T) {
		err := newTestError("内部错误")
		err1 := errors.Newf(400, "请求参数错误", "").TakeOption(errors.WithDetail(err.Error()))

		gotErr := new(testError)
		ok := gerr.As(err1, &gotErr)
		require.True(t, ok)
		require.Equal(t, gotErr.Error(), "内部错误")
	})
}

func Test_FromError(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		var err error
		gotErr := errors.ParseError(err)
		require.Nil(t, gotErr)
	})
	t.Run("not Error", func(t *testing.T) {
		err := newTestError("内部错误")
		gotErr := errors.ParseError(err)

		require.Equal(t, gotErr.Code, int32(500))
		require.Equal(t, gotErr.Message, "服务器错误")
		require.Equal(t, gotErr.Error(), "服务器错误: 内部错误")
	})
	t.Run("Error", func(t *testing.T) {
		err := errors.New(400, "请求参数错误", "")
		gotErr := errors.ParseError(err)

		require.Equal(t, gotErr.Code, int32(400))
		require.Equal(t, gotErr.Message, "请求参数错误")
		require.Equal(t, gotErr.Metadata, map[string]string(nil))
		require.Equal(t, gotErr.Error(), "请求参数错误")
	})
}

func Test_Error_EqualCode(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		var err1 error
		require.False(t, errors.EqualCode(err1, 400))
	})
	t.Run("nil Error", func(t *testing.T) {
		var err2 *errors.Error
		require.False(t, errors.EqualCode(err2, 400))
	})
	t.Run("not equal", func(t *testing.T) {
		err1 := new(testError)
		require.False(t, errors.EqualCode(err1, 400))

		err2 := errors.Newf(400, "请求参数错误1", "")
		require.False(t, errors.EqualCode(err2, 500))
	})
	t.Run("equal", func(t *testing.T) {
		err1 := errors.Newf(400, "请求参数错误1", "")
		require.True(t, errors.EqualCode(err1, 400))
	})
}
