package lambdalabs

import (
	"github.com/unweave/unweave/api/types"
)

const apiURL = "https://cloud.lambdalabs.com/api/v1/"

// err400 can happen when ll doesn't have enough capacity to create the instance
func err400(msg string, err error) *types.Error {
	return &types.Error{
		Code:     400,
		Provider: types.LambdaLabsProvider,
		Message:  msg,
		Err:      err,
	}
}

func err401(msg string, err error) *types.Error {
	return &types.Error{
		Code:       401,
		Provider:   types.LambdaLabsProvider,
		Message:    msg,
		Suggestion: "Make sure your LambdaLabs credentials are up to date",
		Err:        err,
	}
}

func err403(msg string, err error) *types.Error {
	return &types.Error{
		Code:       403,
		Provider:   types.LambdaLabsProvider,
		Message:    msg,
		Suggestion: "Make sure your LambdaLabs credentials are up to date",
		Err:        err,
	}
}

func err404(msg string, err error) *types.Error {
	return &types.Error{
		Code:       404,
		Provider:   types.LambdaLabsProvider,
		Message:    msg,
		Suggestion: "",
		Err:        err,
	}
}

func err500(msg string, err error) *types.Error {
	if msg == "" {
		msg = "Unknown error"
	}
	return &types.Error{
		Code:       500,
		Message:    msg,
		Suggestion: "LambdaLabs might be experiencing issues. Check the service status page at https://status.lambdalabs.com/",
		Provider:   types.LambdaLabsProvider,
		Err:        err,
	}
}

// We return this when LambdaLabs doesn't have enough capacity to create the instance.
func err503(msg string, err error) *types.Error {
	return &types.Error{
		Code:     503,
		Provider: types.LambdaLabsProvider,
		Message:  msg,
		Err:      err,
	}
}

func errUnknown(code int, err error) *types.Error {
	return &types.Error{
		Code:       code,
		Message:    "Unknown error",
		Suggestion: "",
		Provider:   types.LambdaLabsProvider,
		Err:        err,
	}
}
