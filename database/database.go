package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Dennis2340/book-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var bookCollection *mongo.Collection
var Collection = "Books"

func GetClient() *mongo.Client {

	// load the uri from env
	uri := os.Getenv("DATABASE_URL")

	// getting the context
	if client != nil {
		return client
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
	}

	return client
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	if bookCollection != nil {
		return bookCollection
	}

	bookCollection = client.Database("BookShop").Collection(collectionName)
	return bookCollection
}

func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if client == nil {
		return
	}

	err := client.Disconnect(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

func List_Books() ([]models.Book, error) {
	// get client
	client := GetClient()
	bookCollection := GetCollection(client, Collection)
	// mongo queries //
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var bookList []models.Book

	cursor, err := bookCollection.Find(ctx, bson.D{})
	defer cursor.Close(ctx)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	for cursor.Next(ctx) {
		var book models.Book
		err := cursor.Decode(&book)
		if err != nil {
			log.Println("Error decoding book:", err)
			continue
		}
		bookList = append(bookList, book)
	}

	return bookList, nil
}

func FindBook(name string) *models.Book {
	client := GetClient()
	bookCollection := GetCollection(client, Collection)

	// create the context //
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var book *models.Book

	filter := bson.D{{Key: "name", Value: name}}

	err := bookCollection.FindOne(ctx, filter).Decode(&book)
	if err != nil {
		return nil
	}
	return book
}

func CreateBook(book *models.Book) string {
	// get the client and the collection
	client := GetClient()
	bookCollection := GetCollection(client, Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := bookCollection.InsertOne(ctx, models.Book{
		Id:    primitive.NewObjectID(),
		Name:  book.Name,
		Price: book.Price,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return result.InsertedID.(primitive.ObjectID).Hex()
}

func DeleteBook(identifier string) error {
	client := GetClient()
	bookCollection := GetCollection(client, Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var filter primitive.D

	objectID, err := primitive.ObjectIDFromHex(identifier)
	if err == nil {
		// If identifier is a valid ObjectID, search by _id
		filter = bson.D{{Key: "id", Value: objectID}}
	} else {
		// If identifier is not a valid ObjectID, assume it's a name
		filter = bson.D{{Key: "name", Value: identifier}}
	}

	log.Printf("Deleting book with filter: %v", filter)

	result, err := bookCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting book: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no book found with identifier: %s", identifier)
	}

	log.Printf("Successfully deleted book with identifier: %s", identifier)
	return nil
}

func UpdateBook(id string, updatedBook *models.Book) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	client := GetClient()
	bookCollection := GetCollection(client, Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: objectID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: updatedBook.Name},
		{Key: "price", Value: updatedBook.Price},
	}}}

	result, err := bookCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating book: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no book found with id: %s", id)
	}

	return nil
}
