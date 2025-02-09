package order

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/url"
	pmongo "oms/pkg/mongo"
	"oms/service"
	"time"
)

type Controller struct {
	service service.Service
}

func NewController(s service.Service) *Controller {
	return &Controller{
		service: s,
	}
}

func (c *Controller) HandleGetOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		client := pmongo.GetMongoClient()

		// Get MongoDB collection
		collection := client.Database("orders").Collection("orders")

		// query parameters
		tenantID := ctx.Query("tenant_id")
		sellerID := ctx.Query("seller_id")
		status := ctx.Query("status")
		startDate := ctx.Query("start_date")
		endDate := ctx.Query("end_date")

		//MongoDB filter
		filter := bson.M{}

		if tenantID != "" {
			filter["tenantid"] = tenantID
		}
		if sellerID != "" {
			filter["sellerid"] = sellerID
		}
		if status != "" {
			filter["status"] = status
		}

		// Handle date range filter
		if startDate != "" && endDate != "" {
			startTime, err1 := time.Parse("2006-01-02", startDate)
			endTime, err2 := time.Parse("2006-01-02", endDate)

			if err1 == nil && err2 == nil {
				// Adjust time range to include full day
				startTime = startTime.UTC()
				endTime = endTime.Add(24 * time.Hour).UTC()

				filter["createdat"] = bson.M{"$gte": startTime, "$lt": endTime}
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
				return
			}
		}

		// Query MongoDB
		cursor, err := collection.Find(context.Background(), filter, options.Find().SetSort(bson.M{"createdat": -1}))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer cursor.Close(context.Background())

		// Decode results
		var orders []bson.M
		if err := cursor.All(context.Background(), &orders); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode orders"})
			return
		}

		ctx.JSON(http.StatusOK, orders)
	}
}

func (c *Controller) HandleOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		filePath := ctx.DefaultQuery("filePath", "")

		decodedPath, _ := url.QueryUnescape(filePath)

		if decodedPath == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "File path is required",
			})
			return
		}

		err := c.service.ProcessOrder(decodedPath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "File processed successfully",
		})
	}

}
