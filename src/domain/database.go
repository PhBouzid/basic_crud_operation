package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)


type Database interface {
	DatabaseAddress(addr string)
	Connect() error
	Disconnect() error
}

type MongoDatabase struct{
	Clt *mongo.Client
	address string
}

func (m *MongoDatabase) DatabaseAddress(addr string){
	m.address = addr
}

func (m *MongoDatabase) Disconnect() error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := m.Clt.Disconnect(ctx)
	return err
}


func(m *MongoDatabase) Connect() error{
	var err error
	m.Clt, err = mongo.NewClient(options.Client().ApplyURI(m.address))
	if err!=nil{
		//log.Fatalf("Build connection client for mongo failed with error: %s",err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = m.Clt.Connect(ctx)
	if err != nil {
		log.Printf("Connection not establish error: %s", err)
		return err
	}
	err = m.Clt.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Printf("Connection not establish error: %s", err)
		return err
	}
	return nil
}
