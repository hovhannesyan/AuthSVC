package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hovhannesyan/AuthSVC/pkg/config"
	"github.com/hovhannesyan/AuthSVC/pkg/db"
	"github.com/hovhannesyan/AuthSVC/pkg/pb"
	"github.com/hovhannesyan/AuthSVC/pkg/services"
	"github.com/hovhannesyan/AuthSVC/pkg/utils"
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
}

func main() {
	defer f.Close()

	dbHandler := db.Init(db.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	jwt := utils.JwtWrapper{
		SecretKey:       os.Getenv("JWT_SECRET_KEY"),
		Issuer:          "AuthSVC",
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
