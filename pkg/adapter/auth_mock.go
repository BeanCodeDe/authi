package adapter

type (
	RefreshTokenRecord struct {
		UserId       string
		Token        string
		RefreshToken string
	}

	GetTokenRecord struct {
		UserId   string
		Password string
	}

	TokenResponse struct {
		TokenResponseDTO *TokenResponseDTO
		Err              error
	}

	AdapterMock struct {
		RefreshTokenRecordArray   []*RefreshTokenRecord
		RefreshTokenResponseArray []*TokenResponse

		GetTokenRecordArray   []*GetTokenRecord
		GetTokenResponseArray []*TokenResponse
	}
)

func (mock *AdapterMock) RefreshToken(userId string, token string, refreshToken string) (*TokenResponseDTO, error) {
	refreshTokenRecord := &RefreshTokenRecord{UserId: userId, Token: token, RefreshToken: refreshToken}
	mock.RefreshTokenRecordArray = append(mock.RefreshTokenRecordArray, refreshTokenRecord)

	response := mock.RefreshTokenResponseArray[len(mock.RefreshTokenRecordArray)-1]
	return response.TokenResponseDTO, response.Err
}
func (mock *AdapterMock) GetToken(userId string, password string) (*TokenResponseDTO, error) {
	getTokenRecord := &GetTokenRecord{UserId: userId, Password: password}
	mock.GetTokenRecordArray = append(mock.GetTokenRecordArray, getTokenRecord)

	response := mock.GetTokenResponseArray[len(mock.GetTokenRecordArray)-1]
	return response.TokenResponseDTO, response.Err
}
