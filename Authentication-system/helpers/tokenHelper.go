package helper

import (
	//"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

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

	//this is our payload
	accessTokenClaims := &Claims{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		User_type:  userType,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Second).Unix(),
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
			ExpiresAt: time.Now().Add(50 * time.Second).Unix(), //refresh token expires in 7 days
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

func UpdateAllTokens(refreshToken string) (string, string, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	claims, validationErr := ValidateToken(refreshToken)
	if validationErr != nil {
		return "", "", errors.New("expired refresh token")
	}

	//sending all three if there'll be any error while generating also that also will be send
	return GenerateAllTokens(claims.Email, claims.First_name, claims.Last_name, claims.User_type, claims.Uid)

}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	//this method does all the work compares the new generated sign with the already present sign in the token and also checks if th token is already expired, and then set all the payload info that is available in the token into the claims object that is created above if the token is valid
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	return claims, err
}
