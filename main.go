package main

import (
	"context"
	"os"

	"rest-projects/internal/metrics"
	"rest-projects/internal/tracer"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func main() {
	cleanup := tracer.InitTracer()
	defer cleanup(context.Background())

	server := gin.New()
	server.Use(otelgin.Middleware("GoAPI"))

	server.GET("/", func(c *gin.Context) {
		tracer := otel.Tracer("GoAPI") // Assuming your service name is GoAPI

		_, span := tracer.Start(c.Request.Context(), "get-hello-world")
		defer span.End()

		span.SetName("get-hello-world")
		span.SetAttributes(attribute.String("message", "Hello, World!"))

		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server.GET("/metrics", metrics.PrometheusHandler())

	server.Run(":" + port)
}
