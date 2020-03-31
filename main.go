package main

import (
	"fmt"
	"github.com/KennyChenFight/dcard-simple-demo/handler"
	"github.com/KennyChenFight/dcard-simple-demo/lib/auth"
	"github.com/KennyChenFight/dcard-simple-demo/lib/config"
	"github.com/KennyChenFight/dcard-simple-demo/lib/httputil"
	"github.com/KennyChenFight/dcard-simple-demo/lib/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	// autoload .env file
	_ "github.com/joho/godotenv/autoload"
	// register postgres driver
	_ "github.com/lib/pq"
	"log"
	"time"
	"xorm.io/core"
	"xorm.io/xorm"
)

// init the various object and inject the database object to the modules
func init() {
	// postgres connection
	connectStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		config.GetStr("DB_HOST"), config.GetStr("DB_PORT"), config.GetStr("DB_NAME"),
		config.GetStr("DB_USERNAME"), config.GetStr("DB_PASSWORD"), config.GetStr("DB_SSL_MODE"))
	db, _ := xorm.NewEngine("postgres", connectStr)
	if err := db.Ping(); err != nil {
		log.Panic("DB connection initialization failed:", err)
	}

	// redis connection
	redisOptions := redis.Options{
		Network:  "tcp",
		Addr:     config.GetStr("REDIS_ENDPOINT"),
		Password: config.GetStr("REDIS_PASSWORD"),
		PoolSize: config.GetInt("REDIS_POOL_SIZE"),
	}
	redisClient := redis.NewClient(&redisOptions)
	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Panic("Redis connection initialization failed:", err)
	}

	// jwt setting
	secretKey := config.GetBytes("SECRET_KEY")
	tokenLifeTime := time.Duration(config.GetInt("TOKEN_LIFETIME")) * time.Minute

	auth.Init(secretKey, tokenLifeTime)
	middleware.Init(db, redisClient)
	httputil.Init(core.SnakeMapper{})
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Static("/api-docs", "./swagger/dist")

	v1Router := router.Group("/v1/")
	{
		authRouter := v1Router.Group("/auth/")
		{
			authRouter.POST("/", middleware.Plain(), handler.Login)
		}

		userRouter := v1Router.Group("/users/")
		{
			userRouter.POST("/", middleware.Plain(), handler.UserCreate)
		}

		pairRouter := v1Router.Group("/pairs/")
		pairRouter.Use(middleware.Auth())
		{
			pairRouter.POST("/", middleware.IPLimitIntercept(), middleware.TX(), handler.PairCreate)
			pairRouter.GET("/", middleware.IPLimitIntercept(), middleware.TX(), handler.PairGetOne)
		}

	}
	return router
}

func main() {
	router := setupRouter()
	router.Run()
}
