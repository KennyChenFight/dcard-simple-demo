package handler

import (
	"errors"
	"github.com/KennyChenFight/dcard-simple-demo/lib/constant"
	"github.com/KennyChenFight/dcard-simple-demo/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"xorm.io/xorm"
)

func PairCreate(c *gin.Context) {
	userId := c.GetString(constant.UserId)

	session := c.MustGet(constant.Session).(*xorm.Session)
	queryStr := `SELECT * FROM users WHERE id != ? and not exists (SELECT 1 FROM pairs WHERE user_id_one = ?) ORDER BY random() LIMIT 1;`

	var user model.User
	found, err := session.SQL(queryStr, userId, userId).Get(&user)
	if err != nil {
		c.Set(constant.StatusCode, http.StatusInternalServerError)
		c.Set(constant.Error, err)
		return
	}
	if !found {
		c.Set(constant.StatusCode, http.StatusForbidden)
		c.Set(constant.Error, errors.New("already have pair"))
		return
	}

	var pair model.Pair
	pair.UserIdOne = userId
	pair.UserIdTwo = user.Id
	affected, err := session.Insert(&pair)
	if err != nil {
		c.Set(constant.StatusCode, http.StatusInternalServerError)
		c.Set(constant.Error, err)
		return
	}

	if affected == 0 {
		c.Set(constant.StatusCode, http.StatusForbidden)
		c.Set(constant.Error, errors.New("already have pair"))
		return
	}

	c.Set(constant.StatusCode, http.StatusCreated)
	c.Set(constant.Output, map[string]interface{}{"user_id_two": pair.UserIdTwo})
}

func PairGetOne(c *gin.Context) {
	var pair model.Pair
	userId := c.GetString(constant.UserId)
	session := c.MustGet(constant.Session).(*xorm.Session)

	found, err := session.Where("user_id_one = ?", userId).Get(&pair)
	if err != nil {
		c.Set(constant.StatusCode, http.StatusInternalServerError)
		c.Set(constant.Error, err)
		return
	}

	if !found {
		c.Set(constant.StatusCode, http.StatusNotFound)
		c.Set(constant.Error, errNotFound)
		return
	}

	c.Set(constant.StatusCode, http.StatusOK)
	c.Set(constant.Output, pair)
}
