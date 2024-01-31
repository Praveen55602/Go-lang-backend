package helper

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// a token is made of 3 things
// Header.Payload.Sign
// the object of this struct will act as a payload for the token
type Claims struct {
	// add the neccessary info you want to add about user
	//but make sure not to add sensitve info
	Email              string
	First_name         string
	Last_name          string
	Uid                string
	User_type          string
	jwt.StandardClaims //this needs to be there will contain info about token
}

func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string) (string, string, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	//this is our payload
	accessTokenClaims := &Claims{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		User_type:  userType,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	}
	//here we are specifing what algo to use and what is the payload
	//the algo which is HS256 will be used in th header
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	//here we're signing the token i.e is adding the third part of the token. the accessToken object already contains both header and payload in encoded form now we've to add the third part which SignedString function will do by taking secret key parameter. and also convert all the component into string and will concat them together in the jwt format.
	tokenString, err := accessToken.SignedString(secretKey)

	if err != nil {
		log.Panic(err)
		return "", "", err
	}

	refreshTokenClaims := &Claims{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		User_type:  userType,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(), //refresh token expires in 7 days
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(secretKey)

	if err != nil {
		log.Panic(err)
		return "", "", err
	}

	return tokenString, refreshTokenString, nil

}
