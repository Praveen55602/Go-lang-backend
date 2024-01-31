package controllers

import (
	"auth/database"
	helper "auth/helpers"
	"auth/models"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(*database.Client, "user")
var validate = validator.New()

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

func Signup(ctx *gin.Context) {
	var newUser models.User

	//we are expecting ctx to have all the data required by user object and if it will satisfy all the condition it will bind it to the newUser object
	err := ctx.Bind(&newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//now that newUser object has been binded with json data from request
	// we check the data for validation that we've set in the struct
	validationErr := validate.Struct(newUser)
	if validationErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	var newCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	countEmail, err := userCollection.CountDocuments(newCtx, bson.M{"email": newUser.Email})
	defer cancel()

	if err != nil {
		log.Panic(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the email"})
		return
	}

	if countEmail > 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "account with this email already exists"})
		return
	}

	newUser.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	newUser.Updated_at = newUser.Created_at
	newUser.ID = primitive.NewObjectID()
	newUser.User_id = newUser.ID.Hex()
	//token, refreshToken, _ := helper.GenerateAllTokens(*newUser.Email, *newUser.First_name, *newUser.Last_name, *newUser.User_type, newUser.User_id)

	hashedPassword, hashErr := HashPassword(*newUser.Password)
	if hashErr != nil {
		log.Panic("error while trying to hash the password")
		return
	}

	newUser.Password = &hashedPassword
	resultInserNumber, insertErr := userCollection.InsertOne(newCtx, newUser)
	if insertErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	defer cancel()
	ctx.JSON(http.StatusOK, resultInserNumber)

}

func Login(ctx *gin.Context) {
	var newCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User
	var foundUser models.User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	searchErr := userCollection.FindOne(newCtx, bson.M{"email": user.Email}).Decode(&foundUser)

	if searchErr != nil || !VerifyPassword(*user.Password, *foundUser.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "wrong email or password"})
		return
	}

	if foundUser.Email == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user does not exists"})
		return
	}

	//token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
	ctx.JSON(http.StatusOK, foundUser)
	//helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

}

func GetUsers() {

}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		var err error
		err = helper.MatchUserTypeToUid(c, userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		//bson.M is used to provide the filter condition in the mongodB driver for go
		err = userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode((&user))
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
