package parser

import "github.com/BeanCodeDe/authi/pkg/adapter"

type (
	ParseTokenResponse struct {
		claim *adapter.Claims
		err   error
	}

	ParseTokenRecord struct {
		authorizationString string
	}

	ParserMock struct {
		parseTokenRecordArray   []*ParseTokenRecord
		parseTokenResponseArray []*ParseTokenResponse
	}
)

func (mock *ParserMock) ParseToken(authorizationString string) (*adapter.Claims, error) {
	parseTokenRecord := &ParseTokenRecord{authorizationString: authorizationString}
	mock.parseTokenRecordArray = append(mock.parseTokenRecordArray, parseTokenRecord)
	response := mock.parseTokenResponseArray[len(mock.parseTokenResponseArray)-1]
	return response.claim, response.err
}
