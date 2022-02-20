package controllers

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/soltani-ard/echo-mongo-api/configs"
	"github.com/soltani-ard/echo-mongo-api/models"
	"github.com/soltani-ard/echo-mongo-api/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

// create collection
var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

// create validator instance
var validate = validator.New()

// CreateUser => Add User to MongoDB
func CreateUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	// validate request body
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	// validator err
	if validatorError := validate.Struct(&user); validatorError != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"data": validatorError.Error()},
		})
	}

	newUser := models.User{
		ID:       primitive.NewObjectID(),
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	return c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &echo.Map{"data": result},
	})
}

// GetUser => get user with id
func GetUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	uid := c.Param("userId")
	var user models.User
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(uid)

	err := userCollection.FindOne(ctx, bson.M{"id": objID}).Decode(&user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusBadRequest,
		Message: "success",
		Data:    &echo.Map{"data": user},
	})
}

// EditUser => update user
func EditUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	// validate request body
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	// validator err
	if validatorError := validate.Struct(&user); validatorError != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &echo.Map{"data": validatorError.Error()},
		})
	}

	update := bson.M{
		"name":     user.Name,
		"location": user.Location,
		"title":    user.Title,
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	var updateUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updateUser)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    &echo.Map{"data": err.Error()},
			})
		}
	}

	return c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &echo.Map{"data": result},
	})
}

// DeleteUser => delete user
func DeleteUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	if result.DeletedCount < 1 {
		return c.JSON(http.StatusNotFound, responses.UserResponse{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    &echo.Map{"data": "user with specified id not found"},
		})
	}

	return c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &echo.Map{
			"data": "user successfully deleted.",
		},
	})
}

// GetAllUsers => get all users list
func GetAllUsers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err := results.Decode(&singleUser); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    &echo.Map{"data": err.Error()},
			})
		}

		users = append(users, singleUser)
	}

	return c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &echo.Map{
			"data": users,
		},
	})
}
