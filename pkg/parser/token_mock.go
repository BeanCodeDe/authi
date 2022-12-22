package parser

import "github.com/BeanCodeDe/authi/pkg/adapter"

type (
	ParseTokenResponse struct {
		Claim *adapter.Claims
		Err   error
	}

	ParseTokenRecord struct {
		AuthorizationString string
	}

	ParserMock struct {
		ParseTokenRecordArray   []*ParseTokenRecord
		ParseTokenResponseArray []*ParseTokenResponse
	}
)

func (mock *ParserMock) ParseToken(authorizationString string) (*adapter.Claims, error) {
	parseTokenRecord := &ParseTokenRecord{AuthorizationString: authorizationString}
	mock.ParseTokenRecordArray = append(mock.ParseTokenRecordArray, parseTokenRecord)
	response := mock.ParseTokenResponseArray[len(mock.ParseTokenResponseArray)-1]
	return response.Claim, response.Err
}
