package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	ipHandler "main/internal/ip/delivery"
	ipRepo "main/internal/ip/repository"
	ipUse "main/internal/ip/usecase"
	mware "main/internal/middleware"
	"os"
	"strconv"
)

func ParseFlags(server *echo.Echo) (int, int, string) {
	limit, _ := strconv.Atoi(os.Getenv("LIMIT"))
	maskLength := os.Getenv("MASK")
	coolDown, _ := strconv.Atoi(os.Getenv("COOLDOWN"))

	return limit, coolDown, maskLength
}

func InitDataBase(server *echo.Echo) *sql.DB {
	err := godotenv.Load("project.env")
	usernameDB := os.Getenv("POSTGRES_USERNAME")
	passwordDB := os.Getenv("POSTGRES_PASSWORD")
	nameDB := os.Getenv("POSTGRES_NAME")

	connectString := "user=" + usernameDB + " password=" + passwordDB + " dbname=" + nameDB + " sslmode=disable"
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		server.Logger.Fatal("NO CONNECTION TO BD", err.Error())
	}

	return db
}

func main() {

	server := echo.New()
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer prLogger.Sync()

	db := InitDataBase(server)
	lim, cool, mask := ParseFlags(server)

	repo := ipRepo.NewIpRepository(db)
	use := ipUse.NewIpUseCase(repo, mask, lim, cool)
	handler := ipHandler.NewIpHandler(use, logger)
	middleware := mware.NewMiddlewareHandler(logger)
	middleware.SetMiddleware(server)

	handler.InitHandlers(server)
	port := os.Getenv("PORT")
	server.Logger.Fatal(server.Start(port))
}
