package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Article struct {
	Title     string    `bson:"title"`
	Tags      []string  `bson:"tags"`
	Body      string    `bson:"body"`
	CreatedAt time.Time `bson:"created_at"`
}

var (
	client     *mongo.Client
	collection *mongo.Collection
)

func ConnectToMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get MongoDB URI from environment variable or use default
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27018"
	}

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	// Get collection
	collection = client.Database("golang").Collection("article")
	log.Println("Connected to MongoDB successfully!")
	return nil
}

func CloseMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.Disconnect(ctx)
}

func StoreArticle(article *Article) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if article.CreatedAt.IsZero() {
		article.CreatedAt = time.Now()
	}

	_, err := collection.InsertOne(ctx, article)
	if err != nil {
		log.Printf("Error storing article: %v", err)
		return err
	}

	log.Printf("Article stored successfully: %s", article.Title)
	return nil
}

func GetTopTags(n int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{"$unwind", "$tags"}},
		{{"$group", bson.D{
			{"_id", "$tags"},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$sort", bson.D{{"count", -1}}}},
		{{"$limit", n}},
		{{"$project", bson.D{
			{"_id", 1},
			{"count", 1},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Error aggregating tags: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Printf("Error reading aggregation results: %v", err)
		return nil, err
	}

	var topTags []string
	for _, result := range results {
		if tag, ok := result["_id"].(string); ok {
			topTags = append(topTags, tag)
		}
	}

	log.Printf("Retrieved top %d tags: %v", n, topTags)
	return topTags, nil
}

func GetAllArticles() ([]Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding articles: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var articles []Article
	if err = cursor.All(ctx, &articles); err != nil {
		log.Printf("Error reading articles: %v", err)
		return nil, err
	}

	log.Printf("Retrieved %d articles", len(articles))
	return articles, nil
}

func GetArticleByTitle(title string) (*Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var article Article
	err := collection.FindOne(ctx, bson.M{"title": title}).Decode(&article)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Article not found: %s", title)
			return nil, nil
		}
		log.Printf("Error finding article: %v", err)
		return nil, err
	}

	log.Printf("Retrieved article: %s", title)
	return &article, nil
}
