package main

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Counter struct {
	Name string `bson:"name" json:"name"`
	Num  int    `bson:"num" json:"num"`
}

var client *mongo.Client
var collection *mongo.Collection
var counterCache = make(map[string]int)
var mu sync.Mutex

func pushCacheToDB() {
	mu.Lock()
	defer mu.Unlock()
	if len(counterCache) == 0 {
		return
	}
	LogInfo("Pushing cache to DB:", counterCache)
	var operations []mongo.WriteModel
	for name, num := range counterCache {
		model := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"name": name}).
			SetUpdate(bson.M{"$set": bson.M{"num": num}}).
			SetUpsert(true)
		operations = append(operations, model)
	}
	if len(operations) > 0 {
		_, err := collection.BulkWrite(context.Background(), operations)
		if err != nil {
			LogErrorf("Error pushing cache to DB: %v", err)
		} else {
			counterCache = make(map[string]int)
		}
	}
}
func getCountByName(name string, num int) (Counter, error) {
	if name == "demo" {
		return Counter{Name: "demo", Num: 1234567890}, nil
	}
	if num > 0 {
		return Counter{Name: name, Num: num}, nil
	}
	mu.Lock()
	defer mu.Unlock()
	currentNum, inCache := counterCache[name]
	if !inCache {
		var counterFromDB Counter
		err := collection.FindOne(context.Background(), bson.M{"name": name}).Decode(&counterFromDB)
		if err != nil && err != mongo.ErrNoDocuments {
			LogErrorf("Error getting count from DB for %s: %v", name, err)
			return Counter{}, err
		}
		currentNum = counterFromDB.Num
	}
	currentNum++
	counterCache[name] = currentNum
	return Counter{Name: name, Num: currentNum}, nil
}
