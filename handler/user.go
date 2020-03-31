package handler

import (
	"errors"
	"github.com/KennyChenFight/dcard-simple-demo/lib/auth"
	"github.com/KennyChenFight/dcard-simple-demo/lib/constant"
	"github.com/KennyChenFight/dcard-simple-demo/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"xorm.io/xorm"
)

func UserCreate(c *gin.Context) {
	var user struct {
		model.User `xorm:"extends"`
		Password   string `xorm:"-" json:"password" binding:"required"`
	}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.Set(constant.StatusCode, http.StatusBadRequest)
		c.Set(constant.Error, err)
		return
	}

	uid, err := uuid.NewRandom()
	if err != nil {
		c.Set(constant.StatusCode, http.StatusInternalServerError)
		c.Set(constant.Error, err)
		return
	}
	user.Id = uid.String()
	if digest, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
		c.Set(constant.StatusCode, http.StatusInternalServerError)
		c.Set(constant.Error, err)
		return
	} else {
		user.PasswordDigest = string(digest)
	}

	db := c.MustGet(constant.Db).(*xorm.Engine)
	session := db.NewSession()
	defer session.Close()

	q := `insert into users(id, email, password_digest, name)
			select ?, ?, ?, ?
			where not exists (select 1 from users where email = ?)`

	result, err := db.Exec(q, user.Id, user.Email, user.PasswordDigest, user.Name, user.Email)
	if err != nil {
		c.Set(constant.StatusCode, http.StatusInternalServerError)
		c.Set(constant.Error, err)
		return
	}

	affected, err := result.RowsAffected()
	if affected == 0 {
		c.Set(constant.StatusCode, http.StatusForbidden)
		c.Set(constant.Error, errors.New("the email is already used"))
		return
	}

	if err := session.Commit(); err != nil {
		c.Set(constant.StatusCode, http.StatusInternalServerError)
		c.Set(constant.Error, err)
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
		c.Set(constant.StatusCode, http.StatusCreated)
		c.Set(constant.Output, map[string]interface{}{"userId": user.Id})
	}
}
