package parser

import (
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

const (
	validPubKey = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAvWn5rsSAy+uZb1E1Bxit
V8RJK9f1LjBil7iET81eTQUqBKq3sCm/8J42QmfE/WOHoAUQntFFuDIrP88eLzaj
XbGxjGWZQfdY0ltEHwoOBRsDBkcmiZdECHedYTRdAgsq1KLlItZgGs6M+qpy7cT9
IpQnOjpOCK9vlDCzQLlmikLTo/IjVRbUYBnsFjLj54Z1iM02LfgOE6RTACMmy07Y
byhWfdF9xFh0ODshZlaMC/wsdr+P0abn0oAiaSHXkl6Xlbs/osuSjavC6TaU8cFJ
UKBxx61oG+Goa9kaP0nMjz01apPH7fSGSRRV51ku0soq07Yhg8olw7uk40KqPkPR
WR5ylghQQGvMwD366dgEWA8LhBLx+9d6OrhWjJDlU75iISqIpriNzAIZFkrv9yp9
m24m2W+F/OkG5byG1M+6Q8yCmPsexJ0b4qz4Q7LtKWRdkz2UQgNr/Y9GSSkmOde6
6sNeFTcCK1v1SHGx4XhkTqOKsaiBWS4KVRA7bsWYzwsjqtfX28+vUQJUWWR+BXpz
4PJI8NPQR+qaGLSf97wCbKgyGQRH183MkoMrq6PHq1LT3IRYxnl/aa5turc1k2Kd
TsWxzjc68YMOp9sEX8W00aXoPDWXSnSoly96CTxiUoQI4oJ/WHfvqPtmU4s1Qq2h
0buEZdqExx8v5UmRxhDBkmMCAwEAAQ==
-----END PUBLIC KEY-----`

	expiredToken     = `Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYTlkYzY1MzktYjAyYy00ZTcwLTg0NjYtZTQwZTI4NTNjMTgwIiwiZXhwIjoxNjcxNTUyMjYxfQ.FFB_HEJsaczB5z9twn3zghsBON4cE--uCIqZ3hbrwi_aNTEFprP492aIzIZfDJ3ObXfn8N4od2cA_jmInc0G29zzWtNEaoHX73lIoV5wbAlREQeiQmY8rJR7zvi7e6qWESjTC_KUwWLOlNTFigWxONtzzNAikBuxEbNROekH6TO--vngM6pSQCoi72PAIEU3981zyiTjrQK2EGCxlIZaOcFbumHyh4dU90D6fVUUlcHrwh5jsga48RmBwLrbwgOyDJJxrYVmE8dg8rW_bSBHy79V_DUjWRncicKqZ2XOpbGJxYAuJPx4wOHHCYupMCaPxkGyb7ScBRGM01PPAdOi7oS1-8PmqXsOxfjvW-nF1DqiWoekkZGr_BAfvzTyxaqyPAaYQJPsPqZmM-3zQ-St4wG0iuATYHMcxnaEgXGFQslZWd23uQKP8NKmCr9-uJbmuTQ1gRnUay3z4wRS5-jCf33We9k_zwvxxnciJw1Rq95c2_3zeV3-KiOLdGR5Z_s3S_B1avUj0Z8DBEErUSo29DoopkoXf3Rc2Z8atSW0dpbUR27Cy1fG2HGm1rASVQFFIvYKrEfMM_m7GMAH8Rao2qFwj2kkBbvm3IsBA7sPwbWXTxMTwLBW9H8XXhcLvMdgbdw0ZWaBF_zWh2JxKIL7GFccA9ape8qBeiXaJFN6sz8`
	validToken       = `Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYTlkYzY1MzktYjAyYy00ZTcwLTg0NjYtZTQwZTI4NTNjMTgwIiwiZXhwIjo5OTk5OTk5OTk5fQ.qHgfMGKi14D9boWOds7Et57cFBF4ZlGaZ2fT4p5UTR8pnoLIp-wt0t6i3wzULb266BsOh1-rVVfEKTEMyCU7cWCefIGbBoVqnUFDBNqqolsLQjv5hzqHUwwr4tnfnbIkGGdg80TmF53DDBp8yGO0eLSXY_9VqS1qWy6RYvWSJAu5b2bmDYZdpMkrc8N2STtW5VPQrYki_YhU3YFraCyArH8C0F2QEVDoUYGL2ZxxUz5-cVTGaMk-3jDPSGpLNRaQzZZrab5W52fdEVJX0T1llPGR7agFDO6lipSvR4bCRyMH_z35IkNUSFDozM2OtnBeUMf8FUdmUQIKRSezWiwfF5rtaglDaVDYiGgSO9dvGqJQzpychEhRhbwEHDzbTVcMIenTK8NqqrpdsP7gDAqUhg5SCHASpPVdn_b9U1y3LfbtEj7rJY0mUTsbkyObolOr3Xb9pCwI_Ip7KjtY_BhmmlF0ZTAH32b-_LXA8bWbxZg2Mt4-MZQRhzkhJola--5vvxr--_Rc9Z_M0WEfaDkHSarETVBF7Dgwb2JslCkI6fWkcCnWzn2oFCzup2ubCCV2jRibcNNoy2I98OgBDIwH5Y7knumiZxkL6G09M0IxrXhzhgyBKD4l9Kvb35vSdrz9BjG6moYsqtbM-MhO41CK38DHGqcuXNyHYAXPdHYQgdQ`
	wrongFormatToken = `someToken`

	userId = "a9dc6539-b02c-4e70-8466-e40e2853c180"
)

func TestParseToken_Successfully(t *testing.T) {
	jwtParser := loadJWTParser(t, validPubKey)
	claim, err := jwtParser.ParseToken(validToken)

	assert.Nil(t, err)
	assert.NotNil(t, claim)
	assert.Equal(t, userId, claim.UserId.String())
}

func TestParseToken_ErrorWhenSplittingToken(t *testing.T) {
	jwtParser := loadJWTParser(t, validPubKey)
	claim, err := jwtParser.ParseToken(wrongFormatToken)

	assert.NotNil(t, err)
	assert.Nil(t, claim)
	assert.Equal(t, ErrTokenNotFound, err)
}

func TestParseToken_ErrorWhenParsingExpiredKey(t *testing.T) {
	jwtParser := loadJWTParser(t, validPubKey)
	claim, err := jwtParser.ParseToken(expiredToken)

	assert.NotNil(t, err)
	assert.Nil(t, claim)
	assert.ErrorIs(t, err, ErrClaimCouldNotBeParsed)
	assert.ErrorContains(t, err, "token is expired by")
}

func loadJWTParser(t *testing.T, pubKey string) *JWTParser {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubKey))
	if err != nil {
		assert.Fail(t, "Token couldn't be loaded: %w", err)
	}
	return &JWTParser{verifyKey: publicKey}
}
