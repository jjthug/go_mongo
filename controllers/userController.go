package controllers

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"go_mongo/database"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()
func HashPassword(password string)

func VerifyPassword(password string)

func SignUp

func Login

func GetUsers


func GetUser()
