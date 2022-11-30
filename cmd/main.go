package main

import (
	"github.com/hovhannesyan/RiskIndex-AuthSVC/pkg/config"
	"github.com/hovhannesyan/RiskIndex-AuthSVC/pkg/db"
	"github.com/hovhannesyan/RiskIndex-AuthSVC/pkg/pb"
	"github.com/hovhannesyan/RiskIndex-AuthSVC/pkg/services"
	"github.com/hovhannesyan/RiskIndex-AuthSVC/pkg/utils"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := config.LoadConfig(); err != nil {
		logrus.Fatalln(err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading.env file: %s", err.Error())
	}

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
