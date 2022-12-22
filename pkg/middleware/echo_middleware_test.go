package middleware

import (
	"testing"

	"github.com/BeanCodeDe/authi/pkg/adapter"
	"github.com/BeanCodeDe/authi/pkg/parser"
)

func TestCheckToken_Successfully(t *testing.T) {
	tokenParserMock := &parser.ParserMock{parseTokenResponseArray: []*parser.ParseTokenResponse{{claim: &adapter.Claims{}, err: nil}}}
	middleware := &EchoMiddleware{tokenParser: tokenParserMock}

	middleware.CheckToken(nil)
	//TODO Implement testing when test classes are implemented by echo
}
