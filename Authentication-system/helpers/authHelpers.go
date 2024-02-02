package helper

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CheckUserType(c *gin.Context, role string) error {
	userType := c.GetString("user_type")
	var err error = nil
	if userType != role {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	return err
}

func MatchUserTypeToUid(c *gin.Context, userId string) error {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	var err error = nil

	if userType == "USER" && uid != userId {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	err = CheckUserType(c, userType)
	return err
}

func HashPassword(password string) (string, error) {
	//will return a byte array and and error
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(actualPassword string, claimedPassword string) bool {
	//actual password is stored in the hashed form in the db so this function will compare the hash with the claimed password
	err := bcrypt.CompareHashAndPassword([]byte(actualPassword), []byte(claimedPassword))
	return err == nil
}
