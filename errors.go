package gotoken

import "github.com/872409/gatom/util"

const TokenErrorCode = 100

var (
	ErrorGoTokenInvalid       = util.NewCodeError("Invalid goToken", TokenErrorCode+1)
	ErrorGoTokenExpired       = util.NewCodeError("GoToken expired", TokenErrorCode+2)
	ErrorGoTokenGen           = util.NewCodeError("GoToken gen error", TokenErrorCode+3)
	ErrorGoTokenHeaderParams  = util.NewCodeError("GoToken header params error", TokenErrorCode+4)
	ErrorGoTokenHeaderEncoded = util.NewCodeError("GoToken header params encoded error", TokenErrorCode+5)
)
