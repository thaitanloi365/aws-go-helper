package main

import (
	"aws-go-helper/config"
	"aws-go-helper/handlers"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo"
)

var echoLambda *echoadapter.EchoLambda

func init() {
	config.SetupEnv()
	e := echo.New()

	e.GET("/config", handlers.GetConfigHandler)
	e.GET("/stats", handlers.StatsHandler)
	e.GET("/", handlers.GetPublicInfoHandler)
	e.GET("/resize/:name", handlers.ResizeImageHandler)
	echoLambda = echoadapter.New(e)
}

// Handler root handler
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echoLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
