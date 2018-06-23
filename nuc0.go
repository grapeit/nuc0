package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.GET("/led", led)
	r.RunTLS(":3333", "certificate.pem", "key.pem")
}

func led(c *gin.Context) {
	data, error := ioutil.ReadFile("/proc/acpi/nuc_led")
	if error == nil {
		c.JSON(200, gin.H{
			"status": string(data),
		})
	} else {
		c.String(501, error.Error())
	}
}
