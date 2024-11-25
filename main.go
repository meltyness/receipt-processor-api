package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

var receipt_database = make(map[string]*Receipt)

func getItems(c *gin.Context) {
	id := c.Param("id")

	if resp := receipt_database[id]; resp != nil {
		c.JSON(http.StatusOK, receipt_query_success_response{Points: resp.Points})
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func tryInstallReceipt(c *gin.Context) {
	new_receipt := &ReceiptContent{}
	if err := c.BindJSON(&new_receipt); err == nil {
		new_uuid := uuid.NewString()		
		receipt_database[new_uuid] = NewReceipt(new_receipt, new_uuid)
		c.JSON(http.StatusOK, receipt_install_success_response{ID: new_uuid})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func main() {
	fmt.Println("Launching receipt solution, Beta v1.0.0")
	fmt.Println(" - Make sure to check entropy / strength on deployment")

	// Defined in api.go
	RegisterValidators()

	// Setup gin library
	gin.SetMode(gin.ReleaseMode)
	gin.EnableJsonDecoderDisallowUnknownFields()

	// Setup gin router
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// Register specified endpoints
	router.POST("/receipts/process", tryInstallReceipt)
	router.GET("/receipts/:id/points", getItems)

	// Bind
	router.Run("localhost:33824")
}
