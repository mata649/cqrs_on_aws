package auth

import "github.com/mata649/cqrs_on_aws/platform/server/auth"

func SetupAuth(clientID string, clientSecret string, keySecret string) {

	auth.SetJWTSecretKey(keySecret)

}
