package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// SetUp
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	// List
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "recipe_list.html", gin.H{
			// Template items can be added here...
		})
	})
	// Detail
	router.GET("/recipes/:recipe_id/", func(c *gin.Context) {
		recipe_id := c.Param("recipe_id")
		c.String(http.StatusOK, "We will fetch %s", recipe_id)
	})
	// New
	router.GET("/recipes/new/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Recipe Detail Will Go Here",
		})
	})
	
	router.Run(":8080")
}
