package main

import (
	"context"
	"log"
	"time"

	"github.com/Chethu16/Chethu/src/cmd"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var echoLambda *echoadapter.EchoLambda
var mongoClient *mongo.Client

func init() {
	// ✅ MongoDB Atlas URI
	mongoURI := 

	log.Println("🔧 Connecting to MongoDB...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	clientOpts := options.Client().ApplyURI(mongoURI)
	var err error
	mongoClient, err = mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}

	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ MongoDB ping failed: %v", err)
	}

	log.Println("✅ MongoDB connection successful")

	// ✅ Setup Echo server
	log.Println("🚀 Bootstrapping Echo server...")
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ✅ Initialize validator
	validate := validator.New()

	// ✅ Pass DB + Validator into route setup
	db := mongoClient.Database("testing")
	cmd.SetupRoutes(e, db, validate)

	// ✅ Wrap Echo into Lambda adapter
	echoLambda = echoadapter.New(e)
	log.Println("✅ Echo Lambda initialized")
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("📥 Lambda received request: %s", req.Path)
	return echoLambda.ProxyWithContext(ctx, req)
}

func main() {
	log.Println("🟢 Lambda starting...")
	lambda.Start(handler)
}
