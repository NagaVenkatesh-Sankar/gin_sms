package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type payload struct {
	DeviceOS  string    `json:"deviceOS" binding:"required"`
	SmsSchema smsSchema `json:"sms"`
	User      user      `json:"user"`
}
type user struct {
	UserName string `binding:"required"`
	Password string `binding:"required"`
}
type smsSchema struct {
	From    string `json:"from" binding:"required"`
	To      string `json:"to" binding:"required"`
	Message string `json:"message"`
}

var jsonData payload
var data []payload

func main() {
	r := gin.Default()
	r.Use(AuthRequired())
	r.POST("/sms", sms)

	r.Run(":8090") // listen and serve on 0.0.0.0:8090
}

// AuthRequired middleware
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		jsonData = payload{}
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if jsonData.User.UserName != "naga" || jsonData.User.Password != "123" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
			return
		}
	}
}
func sms(c *gin.Context) {
	// Create the device type
	currentDevice, err := CreateDevice(jsonData.DeviceOS)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Send the sms
	smsContext := jsonData.SmsSchema
	_, err = currentDevice.SendSms(smsContext)
	if err != nil {
		log.Printf("Issue with the SMS module, %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SMS Sent successfully."})

}
