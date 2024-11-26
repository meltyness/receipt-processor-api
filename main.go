/*
receipt_solution is a server program built using gin which:
- stores receipts in a JSON format defined in `receipt.go`
  - in doing so, computes a score
  - returns a UUID generated for the receipt-install transaction
  - no effort to deduplicate, each "process" installs a new receipt
  - invalid entries are rejected at install time

- can respond to queries by the UUIDs dispensed

Simply use `go run .` in order to launch the server

 $ go install golang.org/x/pkgsite/cmd/pkgsite@latest
 $ cd myproject
 $ pkgsite -open .
*/
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var receipt_database = myMap{}

// This implements the /process/receipts endpoint
//
// For a given input the API expects a Body with
// fields one would expect on a Receipt, see `receipt.go` for details.
//
// An error is returned when the UUID passed is unknown
func tryGetIdScore(c *gin.Context) {
	id := c.Param("id")

	if resp := receipt_database.Select(id); resp != nil {
		c.JSON(http.StatusOK, receipt_query_success_response{Points: resp.Points})
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

// This implements the /process/receipts endpoint
//
// For a given input the API expects a Body with
// fields one would expect on a Receipt, see `receipt.go` for details.
// Successful response assigns an ID to the receipt
// An error is returned when:
//   - Fields failed validation
//   - The receipt itself was scored in excess of uint64.max
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

// This sets up a gin Engine and listens on 33824
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
