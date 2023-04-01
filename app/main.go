package main

import (
	"github.com/jpdel518/go-ent/infrastructure/file"
	"github.com/jpdel518/go-ent/infrastructure/file/s3"
	"github.com/jpdel518/go-ent/infrastructure/rdb"
	"github.com/jpdel518/go-ent/infrastructure/rdb/mysql"
	"github.com/jpdel518/go-ent/presentation/handler"
	"github.com/jpdel518/go-ent/usecase"
	"github.com/jpdel518/go-ent/utils"
	"os"
	"time"
)

func init() {
	utils.LoadEnv()
	utils.LoggingSettings(os.Getenv("LOG_FILE"))
	mysql.InitDatabase()
}

func main() {
	// Dependency Injection
	client := mysql.NewClient()
	userRepository := rdb.NewUserRepository(client)
	carRepository := rdb.NewCarRepository(client)
	session := s3.NewS3Session()
	userFileRepository := file.NewUserFileRepository(session)
	userUsecase := usecase.NewUserUsecase(userRepository, carRepository, userFileRepository, 30*time.Second)
	handler.NewHandler(userUsecase)
}
