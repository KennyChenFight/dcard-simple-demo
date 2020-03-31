package main

import (
	"github.com/KennyChenFight/dcard-simple-demo/lib/auth"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestAuthStatusOk(t *testing.T) {
	r := gofight.New()

	r.POST("/v1/auth/").
		SetDebug(true).
		SetJSON(gofight.D{
			"email": "kenny@example.com",
			"password": "0000",
		}).
		Run(setupRouter(), func(response gofight.HTTPResponse, request gofight.HTTPRequest) {
			data := []byte(response.Body.String())

			userId, _ := jsonparser.GetString(data, "userId")
			token := response.HeaderMap.Get("Authorization")

			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "97327413-6b65-486f-b299-91be0871f898", userId)
			assert.NotEmpty(t, token)
			assert.Equal(t, "application/json; charset=utf-8", response.HeaderMap.Get("Content-Type"))
		})
}

func TestAuthBadRequestForMissingEmail(t *testing.T) {
	r := gofight.New()

	r.POST("/v1/auth/").
		SetDebug(true).
		SetJSON(gofight.D{
			"password": "0000",
		}).
		Run(setupRouter(), func(response gofight.HTTPResponse, request gofight.HTTPRequest) {
			data := []byte(response.Body.String())

			err, _ := jsonparser.GetString(data, "error", "email")

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Equal(t, "email is a required field", err)
			assert.Equal(t, "application/json; charset=utf-8", response.HeaderMap.Get("Content-Type"))
		})
}

func TestAuthBadRequestForMissingPassword(t *testing.T) {
	r := gofight.New()

	r.POST("/v1/auth/").
		SetDebug(true).
		SetJSON(gofight.D{
			"email": "kenny@example.com",
		}).
		Run(setupRouter(), func(response gofight.HTTPResponse, request gofight.HTTPRequest) {
			data := []byte(response.Body.String())

			err, _ := jsonparser.GetString(data, "error", "password")

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Equal(t, "password is a required field", err)
			assert.Equal(t, "application/json; charset=utf-8", response.HeaderMap.Get("Content-Type"))
		})
}

// todo
// should replace with: https://github.com/DATA-DOG/go-sqlmock
func TestUserStatusCreated(t *testing.T) {
	r := gofight.New()

	r.POST("/v1/users/").
		SetDebug(true).
		SetJSON(gofight.D{
			"email": "test@example.com",
			"password": "0000",
			"name": "tester",
		}).
		Run(setupRouter(), func(response gofight.HTTPResponse, request gofight.HTTPRequest) {
			data := []byte(response.Body.String())

			userId, _ := jsonparser.GetString(data, "userId")

			assert.Equal(t, http.StatusCreated, response.Code)
			assert.NotEmpty(t, userId)
			assert.Equal(t, "application/json; charset=utf-8", response.HeaderMap.Get("Content-Type"))
		})
}

func TestUserBadRequestForMissingEmail(t *testing.T) {
	r := gofight.New()

	r.POST("/v1/users/").
		SetDebug(true).
		SetJSON(gofight.D{
			"password": "0000",
			"name": "tester",
		}).
		Run(setupRouter(), func(response gofight.HTTPResponse, request gofight.HTTPRequest) {
			data := []byte(response.Body.String())

			err, _ := jsonparser.GetString(data, "error", "User.email")

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Equal(t, "email is a required field", err)
			assert.Equal(t, "application/json; charset=utf-8", response.HeaderMap.Get("Content-Type"))
		})
}

func TestUserBadRequestForMissingPassword(t *testing.T) {
	r := gofight.New()

	r.POST("/v1/users/").
		SetDebug(true).
		SetJSON(gofight.D{
			"email": "test@example.com",
			"name": "tester",
		}).
		Run(setupRouter(), func(response gofight.HTTPResponse, request gofight.HTTPRequest) {
			data := []byte(response.Body.String())

			err, _ := jsonparser.GetString(data, "error", "password")

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Equal(t, "password is a required field", err)
			assert.Equal(t, "application/json; charset=utf-8", response.HeaderMap.Get("Content-Type"))
		})
}

func TestUserBadRequestForMissingName(t *testing.T) {
	r := gofight.New()

	r.POST("/v1/users/").
		SetDebug(true).
		SetJSON(gofight.D{
			"email": "test@example.com",
			"password": "0000",
		}).
		Run(setupRouter(), func(response gofight.HTTPResponse, request gofight.HTTPRequest) {
			data := []byte(response.Body.String())

			err, _ := jsonparser.GetString(data, "error", "User.name")

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Equal(t, "name is a required field", err)
			assert.Equal(t, "application/json; charset=utf-8", response.HeaderMap.Get("Content-Type"))
		})
}

func TestUserForbiddenForRepeatedEmail(t *testing.T) {
	r := gofight.New()

	r.POST("/v1/users/").
		SetDebug(true).
		SetJSON(gofight.D{
			"email": "test@example.com",
			"password": "0000",
			"name": "tester",
		}).
		Run(setupRouter(), func(response gofight.HTTPResponse, request gofight.HTTPRequest) {
			data := []byte(response.Body.String())

			err, _ := jsonparser.GetString(data, "error")

			assert.Equal(t, http.StatusForbidden, response.Code)
			assert.Equal(t, "the email is already used", err)
			assert.Equal(t, "application/json; charset=utf-8", response.HeaderMap.Get("Content-Type"))
		})
}

func getDevAuth() (string, error) {
	userId := "97327413-6b65-486f-b299-91be0871f898"
	claims := auth.Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60 * 24 * 30).Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString([]byte("secret"))
}

// todo
// should replace with: https://github.com/elliotchance/redismock
func TestPairStatusOK(t *testing.T) {
	token, _ := getDevAuth()

	r := gofight.New()

	r.GET("/v1/pairs/").
		SetDebug(true).
		SetHeader(gofight.H{
			"Authorization": token,
		}).
		Run(setupRouter(), func(response gofight.HTTPResponse, request gofight.HTTPRequest) {
			data := []byte(response.Body.String())

			userIdOne, _ := jsonparser.GetString(data, "userIdOne")
			userIdTwo, _ := jsonparser.GetString(data, "userIdTwo")

			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "97327413-6b65-486f-b299-91be0871f898", userIdOne)
			assert.Equal(t, "80695811-0bf2-44fd-980d-1635de7734a8", userIdTwo)
			assert.Equal(t, "application/json; charset=utf-8", response.HeaderMap.Get("Content-Type"))
		})
}