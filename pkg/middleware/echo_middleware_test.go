package middleware

import (
	"testing"

	"github.com/BeanCodeDe/authi/pkg/adapter"
)

type (
	parseTokenReturn struct {
		claim *adapter.Claims
		err   error
	}

	tokenParserMock struct {
		parseTokenArray  []string
		parseTokenReturn []*parseTokenReturn
	}
)

func TestCheckToken_Successfully(t *testing.T) {
	tokenParserMock := &tokenParserMock{parseTokenReturn: []*parseTokenReturn{{claim: &adapter.Claims{}, err: nil}}}
	middleware := &EchoMiddleware{tokenParser: tokenParserMock}

	middleware.CheckToken(nil)
	//TODO Implement testing when test classes are implemented by echo
}

func (tokenParserMock *tokenParserMock) ParseToken(authorizationString string) (*adapter.Claims, error) {
	tokenParserMock.parseTokenArray = append(tokenParserMock.parseTokenArray, authorizationString)
	parseTokenReturn := tokenParserMock.parseTokenReturn[len(tokenParserMock.parseTokenArray)-1]
	return parseTokenReturn.claim, parseTokenReturn.err
}
