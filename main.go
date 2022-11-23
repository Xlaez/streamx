package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	config "streamx/configs"
	"streamx/controllers"
	"streamx/routes"
	"streamx/services"
	"streamx/token"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server          *gin.Engine
	ctx             context.Context
	mongoClient     *mongo.Client
	redisClient     *redis.Client
	userCol         *mongo.Collection
	musicCol        *mongo.Collection
	userController  controllers.UserController
	musicController controllers.MusicController
	tokenMaker      token.Maker
	configg         config.Config
)

func initTokenMaker(config config.Config) error {
	var err error
	tokenMaker, err = token.NewPasetoMaker(config.TokenKey)

	if err != nil {
		return fmt.Errorf("cannot create the token maker: %w", err)
	}
	return nil
}

func init() {
	// load environmental vaiables
	configg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("could not load enviromental variables", err)
	}

	// create a context
	ctx = context.TODO()

	// init token maker
	initTokenMaker(configg)

	// connect to mongodb
	mongoconn := options.Client().ApplyURI(configg.DBUri)

	// context, _ := context.WithTimeout(context.Background(), 10*time.Second)

	mongoClient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		log.Panic((err.Error()))
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Panic((err.Error()))
	}

	fmt.Println("MongoDb connected successfully!")

	// connect to redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: configg.RedisUri,
	})

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Panic(err.Error())
	}

	err = redisClient.Set(ctx, "test", "Redis on!", 0).Err()

	if err != nil {
		log.Panic(err.Error())
	}

	fmt.Println("Redis client connected successfully!")

	GetCollections(mongoClient, configg)
	// create gin instance
	server = gin.Default()
	routes.UserRoutes(server, userController, tokenMaker)
	routes.MusicRoutes(server, musicController, tokenMaker)
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoClient.Disconnect(ctx)

	value, err := redisClient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		log.Panic(err.Error())
	}

	router := server.Group("/")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	log.Fatal(server.Run(":" + config.Port))
}

func GetCollections(client *mongo.Client, config config.Config) {
	userCol := client.Database("streamx").Collection("users")
	musicCol := client.Database("streamx").Collection("Music")

	userService := services.NewUserService(userCol, ctx)
	musicService := services.NewMusicService(musicCol, ctx)
	userController = controllers.NewUserController(userService, tokenMaker, config, redisClient)
	musicController = controllers.NewMusicController(musicService, tokenMaker, config, redisClient)

}
