package main

import (
	"log"

	"github.com/Dennis2340/book-api/database"
	"github.com/Dennis2340/book-api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func HandleCreation(c *gin.Context) {
	name := c.Param("name")
	book := database.FindBook(name)
	c.JSON(200, book)
}
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	// initialize database connection
	client := database.GetClient()
	defer database.Disconnect()

	//check database connection
	err = client.Ping(nil, nil)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to database")

	//set up gin router
	router := gin.Default()
	router.SetTrustedProxies([]string{"192.168.1.2"})
	router.POST("/books", routes.HandleCreation)
	router.GET("/books", routes.FindAll)
	router.GET("/books/:name", routes.HandleFindOneBook)
	router.DELETE("/books/:identifier", routes.DeleteBook)

	router.Run("localhost:8082")
}
