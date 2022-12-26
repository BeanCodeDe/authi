package examples

import (
	"fmt"

	"github.com/BeanCodeDe/authi/pkg/adapter"
	"github.com/google/uuid"
)

func AdapterExample() {

	//User id that were previously created over REST
	userId := "693227c8-4178-4e72-b3b7-a8b8bae36f1b"

	//Generate UUID to trace your calls in the logs
	correlationId := uuid.NewString()

	//Initialize authi adapter
	authiAdapter := adapter.NewAuthiAdapter(correlationId)

	//Logging in previously created user with password `mySecretUserPassword`
	token, err := authiAdapter.GetToken(userId, "mySecretUserPassword")

	//Checking if an error occurred while loading user token
	if err != nil {
		panic(err)
	}

	//Printing access token for example
	fmt.Println(token.AccessToken)

	//Printing refresh token for example
	fmt.Println(token.RefreshToken)

	//Refreshing tokens to avoid outdated tokens
	refreshedToken, err := authiAdapter.RefreshToken(userId, token.AccessToken, token.RefreshToken)

	//Checking if an error occurred while loading refreshed token
	if err != nil {
		panic(err)
	}

	//Printing refreshed access token for example
	fmt.Println(refreshedToken.AccessToken)

	//Printing refreshed refresh token for example
	fmt.Println(refreshedToken.RefreshToken)
}
