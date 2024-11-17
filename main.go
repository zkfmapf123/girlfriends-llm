package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zkfmapf123/go-llm/routes"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	routes.GetVacationRouter(r)
	r.Run()
}
