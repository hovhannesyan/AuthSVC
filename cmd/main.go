package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hovhannesyan/RiskIndex-AuthSVC/pkg/config"
	"github.com/hovhannesyan/RiskIndex-AuthSVC/pkg/db"
	"github.com/hovhannesyan/RiskIndex-AuthSVC/pkg/pb"
	"github.com/hovhannesyan/RiskIndex-AuthSVC/pkg/services"
	"github.com/hovhannesyan/RiskIndex-AuthSVC/pkg/utils"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var f *os.File

func init() {
	var err error
	t := time.Now().Format("2006-01-02")

	logFile := fmt.Sprintf("./log/log_%s.txt", t)
	f, err = os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + logFile)
		panic(err)
	}
	logrus.SetOutput(f)

	logrus.SetLevel(logrus.DebugLevel)

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := config.LoadConfig(); err != nil {
		logrus.Fatalln(err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalln("error loading.env file: %s", err)
	}
}

func main() {
	defer f.Close()

	dbHandler := db.Init(db.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	jwt := utils.JwtWrapper{
		SecretKey:       os.Getenv("JWT_SECRET_KEY"),
		Issuer:          "RiskIndex-AuthSVC",
		ExpirationHours: viper.GetInt64("ExpirationHours"),
	}

	lis, err := net.Listen("tcp", viper.GetString("port"))

	if err != nil {
		logrus.Fatalln("failed to listing", err.Error())
	}

	logrus.Println("AuthSVC on", viper.GetString("port"))

	s := services.Server{
		DbHandler: dbHandler,
		Jwt:       jwt,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalln("failed to serve", err.Error())
	}
}
