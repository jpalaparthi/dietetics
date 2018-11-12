package main

import (
	"context"
	database "dietetics/databases"
	"dietetics/handler"
	proto "dietetics/proto/food/raw"

	"github.com/golang/glog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	mgo "gopkg.in/mgo.v2"

	"github.com/casbin/casbin"
	"github.com/casbin/casbin/persist"
	mongodbadapter "github.com/casbin/mongodb-adapter"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"github.com/spf13/viper"
)

var (
	adapter  persist.Adapter
	enforcer *casbin.Enforcer

	dbConnectionStr string
	dbName          string

	casbinDbConnection string
	casbindbName       string
	modelConf          string

	err error
)

// logWrapper implements the server.HandlerWrapper
func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {

		return fn(ctx, req, rsp)
	}
}

// Athenticator the request athenticator
func Athenticator(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {

		metadata, _ := metadata.FromContext(ctx)
		glog.Warningln(metadata)
		sub, act, secret, obj := metadata["Subject"], metadata["Act"], metadata["Secret"], req.Method()

		glog.Warning(metadata)
		if !ValidateUser(secret) {
			glog.Error("user is Unauthenticated to use this service method")
			return status.Error(codes.Unauthenticated, "user is Unauthenticated to use this service method")
		}
		if err := enforcer.LoadPolicy(); err != nil {
			glog.Fatal(err)
		}

		if !enforcer.Enforce(sub, obj, act) {
			glog.Error("user is Unauthenticated to use this service method")
			return status.Error(codes.Unauthenticated, "user is Unauthenticated to use this service method")
		}

		return fn(ctx, req, rsp)
	}

}

// ValidateUser is to validate user based on the password/secret provided.
// Yet to implement this method
func ValidateUser(secret string) bool {
	if secret == "5be1e0496b2e7deb114d58db" {
		return true
	}
	return false
}

func init() {

	viper.SetConfigName("app")      // no need to include file extension
	viper.AddConfigPath("./config") // set the path of your config file

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	} else {
		dbConnectionStr = viper.GetString("connections.DBConnection")
		dbName = viper.GetString("connections.DBName")
		casbinDbConnection = viper.GetString("casbin.DBConnection")
		casbindbName = viper.GetString("casbin.DBName")
		modelConf = viper.GetString("casbin.ModelConf")
	}

	adapter = mongodbadapter.NewAdapter(casbinDbConnection) // Your MongoDB URL.
	enforcer = casbin.NewEnforcer(modelConf, adapter)
	// Load the policy from DB.
}

func main() {
	s, err := database.GetConnection(dbConnectionStr, dbName)
	if err != nil {
		defer s.(*mgo.Session).Close()
		glog.Fatal("mongodb database is not connected")
	}
	defer s.(*mgo.Session).Close()

	if err != nil {
		glog.Fatal("mongodb database is not connected")
	}

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.dietetics"),
		micro.Version("0.1"),
		micro.WrapHandler(logWrapper),
		micro.WrapHandler(Athenticator),
	)

	// Initialise service
	service.Init()

	// create handler instance
	rf := new(handler.RawFood)
	rf.IRawFood = &database.RawFoodDB{Session: s, DBName: dbName}

	// Register Handler
	if err := proto.RegisterRawServiceHandler(service.Server(), rf); err != nil {
		glog.Error(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		glog.Fatal(err)
	}
}
