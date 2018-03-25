package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.GET("/led", led)
	r.Run(":3333")
}

func led(c *gin.Context) {
	data, error := ioutil.ReadFile("/proc/acpi/nuc_led")
	if error == nil {
		c.JSON(200, gin.H{
			"status": string(data),
		})
	} else {
		c.String(500, error.Error())
	}
}
