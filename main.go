package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var receipt_database = myMap{}

type myMap struct {
	sync.Map
}

func (m *myMap) Insert(key string, value *Receipt) {
	m.Store(key, value)
}

func (m *myMap) Select(key string) *Receipt {
	found_receipt, _ := m.Load(key)
	r, _ := found_receipt.(*Receipt)
	if r != nil {
		return r
	} else {
		return nil
	}
}

func tryGetIdScore(c *gin.Context) {
	id := c.Param("id")

	if resp := receipt_database.Select(id); resp != nil {
		c.JSON(http.StatusOK, receipt_query_success_response{Points: resp.Points})
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func tryInstallReceipt(c *gin.Context) {
	new_receipt := &ReceiptContent{}
	if err := c.BindJSON(&new_receipt); err == nil {
		new_uuid := uuid.NewString()
		// Validators validate input fields, but the receipt itself is
		// fallible since requires quantiative assessment of all fields
		// in order to uphold that points can be uint64 or something
		if valid_receipt, err := NewReceipt(new_receipt, new_uuid); err == nil {
			receipt_database.Insert(new_uuid, valid_receipt)
			c.JSON(http.StatusOK, receipt_install_success_response{ID: new_uuid})
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
		}
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
	router.GET("/receipts/:id/points", tryGetIdScore)

	// Bind
	router.Run("localhost:33824")
}
