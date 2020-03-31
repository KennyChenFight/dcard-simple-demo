package middleware

import (
	"github.com/KennyChenFight/dcard-simple-demo/lib/auth"
	"github.com/KennyChenFight/dcard-simple-demo/lib/constant"
	"github.com/KennyChenFight/dcard-simple-demo/lib/lua"
	"github.com/KennyChenFight/dcard-simple-demo/lib/validate"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"time"
	"xorm.io/xorm"
)

const (
	IPLimitPeriod     = 3600
	IPLimitTimeFormat = "2006-01-02 15:04:05"
	IPLimitMaximum  = 1000
)

var (
	db          *xorm.Engine
	redisClient *redis.Client
)

func Init(database *xorm.Engine, client *redis.Client) {
	db = database
	redisClient = client
}

// do nothing and provide injection of database object only
// normally it is used by public endpoint
func Plain() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constant.Db, db)
		c.Set(constant.StatusCode, nil)
		c.Set(constant.Error, nil)
		c.Set(constant.Output, nil)
		c.Next()

		statusCode := c.GetInt(constant.StatusCode)
		err := c.MustGet(constant.Error)
		output := c.MustGet(constant.Output)
		if err != nil {
			if validationErr, ok := err.(validator.ValidationErrors); ok {
				sendResponse(c, statusCode, map[string]interface{}{"error": validationErr.Translate(validate.BindingTrans)})
			} else {
				sendResponse(c, statusCode, map[string]interface{}{"error": err.(error).Error()})
			}
		} else {
			sendResponse(c, statusCode, output)
		}
	}
}

// a middleware to handle user authorization
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := auth.Verify(c.GetHeader("Authorization"))
		if err != nil {
			sendResponse(c, http.StatusUnauthorized, map[string]interface{}{"error": err.Error()})
		} else {
			if newToken, err := auth.Sign(userId); err != nil {
				sendResponse(c, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			} else {
				c.Header("Authorization", newToken) // update JWT Token
				c.Set(constant.UserId, userId)
			}
		}
	}
}

// a middleware prepare a database session for the handler
func TX() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := db.NewSession()
		defer session.Close()

		err := session.Begin()
		if err != nil {
			sendResponse(c, http.StatusInternalServerError, map[string]interface{}{"error": err})
			return
		}
		userId := c.GetString(constant.UserId)
		c.Set(constant.Session, session)
		c.Set(constant.StatusCode, nil)
		c.Set(constant.Error, nil)
		c.Set(constant.Output, nil)
		c.Set(constant.Update, false)
		c.Set(constant.UserId, userId)
		c.Next()

		session = c.MustGet(constant.Session).(*xorm.Session)
		statusCode := c.GetInt(constant.StatusCode)
		err1 := c.MustGet(constant.Error)
		output := c.MustGet(constant.Output)
		update := c.GetBool(constant.Update)
		if err1 == nil {
			if err := session.Commit(); err != nil {
				session.Rollback()
				sendResponse(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			} else {
				sendResponse(c, statusCode, output)
			}
		} else {
			session.Rollback()
			if validationErr, ok := err1.(validator.ValidationErrors); ok {
				if update {
					sendResponse(c, statusCode, map[string]interface{}{"error": validationErr.Translate(validate.UpdateTrans)})
				} else {
					sendResponse(c, statusCode, map[string]interface{}{"error": validationErr.Translate(validate.BindingTrans)})
				}
			} else {
				sendResponse(c, statusCode, map[string]interface{}{"error": err1.(error).Error()})
			}
		}
	}
}

// a middleware for IP Limit checking
// use redis to verify.
func IPLimitIntercept() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now().Unix()
		key := c.Request.URL.Path + "-" + c.Request.Method + "-" + c.ClientIP()
		script := redis.NewScript(lua.SCRIPT)
		args := []interface{}{now, IPLimitMaximum, IPLimitPeriod}

		value, err := script.Run(redisClient, []string{key}, args...).Result()
		// only when redis is disconnected or lua runtime error, error will show up. and it will be rollback.
		// if script's any redis operations are wrong, it will not get error because it is recognized as logical error
		// for example: wrong key
		if err != nil {
			sendResponse(c, http.StatusInternalServerError, err)
			return
		}

		result := value.([]interface{})
		remaining := result[0].(int64)
		// in normal situation: 0~9
		// otherwise, "-1" means too much requests in period
		if remaining == -1 {
			sendResponse(c, http.StatusTooManyRequests, err)
			return
		}
		t := result[1].(int64)
		reset := time.Unix(t, 0).Format(IPLimitTimeFormat)

		c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
		c.Header("X-RateLimit-Reset", reset)
	}
}

// send a http response to the user with JSON format
func sendResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
	c.Abort()
}
