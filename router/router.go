package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/http"
	"go.mongodb.org/mongo-driver/mongo"
	"oms/controller/order"
	"oms/repo"
	"oms/service"
)

func InternalRoutes(ctx context.Context, s *http.Server, client *mongo.Client) (err error) {
	rtr := s.Engine.Group("/api/v1")
	rtr.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "mst"})
	})

	newRepository := repo.NewRepository(client)
	newService := service.NewService(newRepository)
	controller := order.NewController(newService)

	rtr.POST("/order", controller.HandleOrders())
	rtr.GET("/order", controller.HandleGetOrders())

	return
}
