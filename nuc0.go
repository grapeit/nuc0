package main

import (
	"io/ioutil"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.GET("/status", ledStatus)
	//r.RunTLS("0.0.0.0:3333", "certificate.pem", "key.pem")
	r.Run("0.0.0.0:3333")
}

func ledStatus(c *gin.Context) {
	data, error := ioutil.ReadFile("/proc/acpi/nuc_led")
	if error == nil {
		c.JSON(200, gin.H{
			"status": string(data),
		})
	} else {
		c.String(501, error.Error())
	}
}
