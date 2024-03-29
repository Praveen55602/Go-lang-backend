package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id"`
	First_name *string            `json:"first_name" validate:"required,min=2,max=50"`
	Last_name  *string            `json:"last_name"`
	Password   *string            `json:"password" validate:"required"`
	Email      *string            `json:"email" validate:"required"`
	Phone      *string            `json:"phone" validate:"required"`
	//Token         *string            `json:"token"`   storing token and refresh token in the database makes the system a little bit more complex for simple authentication synarios it's best to not store them in the db
	User_type *string `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	//Refresh_token *string            `json:"refresh_token"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	User_id    string    `json:"user_id"`
}
