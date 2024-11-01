package database


import (
    "time"
    "context"
    "fmt"
    "log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    
)




func DBInstance() *mongo.Client{
    MongoDB := "mongodb://localhost:27017"
    fmt.Println(MongoDB)

    ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
     defer cancel()


    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://foo:bar@localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }

   
    
    fmt.Println("Has been connected to mongoDB")
    return client
}



var Client *mongo.Client = DBInstance()


func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
    collection := client.Database("restaurant").Collection(collectionName)

    return collection
}



