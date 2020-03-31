package handler

import (
	"errors"
	"github.com/KennyChenFight/dcard-simple-demo/lib/auth"
	"github.com/KennyChenFight/dcard-simple-demo/lib/constant"
	"github.com/KennyChenFight/dcard-simple-demo/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"xorm.io/xorm"
)

func Login(c *gin.Context) {
	// handle the input
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.Set(constant.StatusCode, http.StatusBadRequest)
		c.Set(constant.Error, err)
		return
	}

	db := c.MustGet(constant.Db).(*xorm.Engine)
	var user model.User
	found, err := db.Where("email = ?", input.Email).Get(&user)
	if err != nil {
		c.Set(constant.StatusCode, http.StatusInternalServerError)
		c.Set(constant.Error, err)
		return
	}

	if !found || bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(input.Password)) != nil {
		c.Set(constant.StatusCode, http.StatusUnauthorized)
		c.Set(constant.Error, errors.New("incorrect email or password"))
		return
	}

	if newToken, err := auth.Sign(user.Id); err != nil {
		c.Set(constant.StatusCode, http.StatusInternalServerError)
		c.Set(constant.Error, err)
	} else {
		// update JWT Token
		c.Header("Authorization", newToken)
		// allow CORS
		c.Header("Access-Control-Expose-Headers", "Authorization")
		c.Set(constant.StatusCode, http.StatusOK)
		c.Set(constant.Output, map[string]interface{}{"userId": user.Id})
	}
}
