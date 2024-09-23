package routes

import (
	"log"
	"sync"

	"github.com/Dennis2340/book-api/database"
	"github.com/Dennis2340/book-api/models"
	"github.com/gin-gonic/gin"
)

func HandleFindOneBook(c *gin.Context) {
	name := c.Param("name")
	book := database.FindBook(name)
	c.JSON(200, book)
}

func HandleCreation(c *gin.Context) {
	var newBook models.Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(400, gin.H{"ERROR": err.Error()})
		return
	}
	createBook := database.CreateBook(&newBook)
	c.JSON(201, createBook)
}

func FindAll(c *gin.Context) {
	booklist, err := database.List_Books()
	if err != nil {
		c.JSON(500, gin.H{"ERROR": err.Error()})
	}
	if len(booklist) == 0 {
		c.JSON(204, gin.H{"message": "No books found"})
		return
	}

	c.JSON(200, gin.H{"books": booklist})

}

func DeleteBook(c *gin.Context) {
	id := c.Param("identifier")
	log.Println("this is the id: ", id)
	err := database.DeleteBook(id)
	if err != nil {
		c.JSON(404, gin.H{"ERROR": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Book Successfully Deleted"})
}

func HandleConcurrentOperations(c *gin.Context) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		booklist, err := database.List_Books()
		if err != nil {
			log.Printf("Error fetching books: %v", err)
			return
		}
		log.Printf("Fetched books: %v", booklist)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		identifier := c.Param("identifier")
		err := database.DeleteBook(identifier)
		if err != nil {
			log.Printf("Error deleting book: %v", err)
			return
		}
		log.Printf("Book successfully deleted")
	}()

	wg.Wait()
	c.JSON(200, gin.H{"message": "Concurrent operations completed"})

}
