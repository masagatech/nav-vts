package main

import (
	"context"
	"fmt"

	"github.com/masagatech/nav-vts/app/connectors"
	"github.com/masagatech/nav-vts/app/handlers"
	"github.com/masagatech/nav-vts/app/models"
	"github.com/masagatech/nav-vts/app/queue"
	"github.com/masagatech/nav-vts/app/servers"
	"github.com/masagatech/nav-vts/app/services"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

var ctx = context.Background()

func main() {
	fmt.Println("test")
	// mongoTest()

	// load config
	c := services.Confiuration{}
	config := c.LoadConfig()

	// connect mongo /// config
	db := connectors.NewDb(config)
	// close the cconnection on application exit

	// connect to redis
	redis := connectors.NewRedis(config)

	// connect rmq /// config
	rmqClient := connectors.NewRabbtMq(config)
	rmq := queue.New(rmqClient)
	rmq.Listner("", "test")

	tcpServer := servers.NewTCPServer(config)
	handlers.NewInit(tcpServer, rmq)

	go tcpServer.Start()

	rest := servers.RESTServer{
		App: &models.App{
			DB:    db,
			Redis: redis,
		},
	}
	rest.Start(config)

}

// connect tile38 /// config

// tile38 := connectors.NewTile38(config)

// if err := tile38.Keys.Set("fleet", "truck1").Point(33.5123, -112.2693).Do(); err != nil {
// 	panic(err)
// }

// if err := tile38.Keys.Set("fleet", "truck2").Point(33.4626, -112.1695).
// 	// optional params
// 	Field("speed", 20).
// 	Expiration(20).
// 	Do(); err != nil {
// 	panic(err)
// }

// response, err := tile38.Search.Nearby("fleet", 33.462, -112.268, 6000).
// 	Where("speed", 0, 100).
// 	Match("truck*").
// 	Format(t38c.FormatPoints).Do()
// if err != nil {
// 	panic(err)
// }

// // truck1 {33.5123 -112.2693}
// fmt.Println(response.Points[0].ID, response.Points[0].Point)
// start tcp server /// config

// start rest service /// config

// var ctx = context.Background()

// func ExampleClient() {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "vts_redis:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})

// 	err := rdb.Set(ctx, "key", "value", 0).Err()
// 	if err != nil {
// 		panic(err)
// 	}

// 	val, err := rdb.Get(ctx, "key").Result()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("key", val)

// 	val2, err := rdb.Get(ctx, "key2").Result()
// 	if err == redis.Nil {
// 		fmt.Println("key2 does not exist")
// 	} else if err != nil {
// 		panic(err)
// 	} else {
// 		fmt.Println("key2", val2)
// 	}
// 	// Output: key value
// 	// key2 does not exist
// }

// func mongoTest() {
// 	clientOptions := options.Client().ApplyURI("mongodb://vts_mongo:27017")
// 	client, err := mongo.Connect(context.TODO(), clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Check the connection
// 	err = client.Ping(context.TODO(), nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Connected to MongoDB!")

// 	// Get a handle for your collection
// 	collection := client.Database("test").Collection("trainers")

// 	// Some dummy data to add to the Database`
// 	ash := Trainer{"Ash", 10, "Pallet Town"}
// 	misty := Trainer{"Misty", 10, "Cerulean City"}
// 	brock := Trainer{"Brock", 15, "Pewter City"}

// 	// Insert a single document
// 	insertResult, err := collection.InsertOne(context.TODO(), ash)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

// 	// Insert multiple documents
// 	trainers := []interface{}{misty, brock}

// 	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

// 	// Update a document
// 	filter := bson.D{{"name", "Ash"}}

// 	update := bson.D{
// 		{"$inc", bson.D{
// 			{"age", 1},
// 		}},
// 	}

// 	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

// 	// Find a single document
// 	var result Trainer

// 	err = collection.FindOne(context.TODO(), filter).Decode(&result)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("Found a single document: %+v\n", result)

// 	findOptions := options.Find()
// 	findOptions.SetLimit(2)

// 	var results []*Trainer

// 	// Finding multiple documents returns a cursor
// 	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Iterate through the cursor
// 	for cur.Next(context.TODO()) {
// 		var elem Trainer
// 		err := cur.Decode(&elem)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		results = append(results, &elem)
// 	}

// 	if err := cur.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// 	// Close the cursor once finished
// 	cur.Close(context.TODO())

// 	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)

// 	// Delete all the documents in the collection
// 	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)

// 	// Close the connection once no longer needed
// 	err = client.Disconnect(context.TODO())

// 	if err != nil {
// 		log.Fatal(err)
// 	} else {
// 		fmt.Println("Connection to MongoDB closed.")
// 	}
// }
