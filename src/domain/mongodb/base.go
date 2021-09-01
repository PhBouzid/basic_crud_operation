package mongodb

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type BaseMongoCRUD struct {
	addr    string
	session *mongo.Client
	coll    *mongo.Collection
	Id      primitive.ObjectID `bson:"-"`
}


func (d *BaseMongoCRUD) connect(db string) error {
	d.addr = db
	fmt.Println("database "+d.addr)
	ss, err := mongo.NewClient(options.Client().ApplyURI(d.addr))
	if err != nil {
		return err
	}
	d.session = ss
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = d.session.Connect(ctx)
	if err != nil {
		log.Printf("Connection not establish error: %s", err)
		return err
	}
	err = d.session.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Printf("Connection not establish error: %s", err)
		return err
	}
	return nil
}

func (d BaseMongoCRUD) disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := d.session.Disconnect(ctx)
	return err
}

func (d BaseMongoCRUD) find(ctx context.Context, id primitive.ObjectID) *mongo.SingleResult {
	opts := options.FindOne().SetSort(bson.D{{"age", 1}})
	result := d.coll.FindOne(ctx, bson.D{{"_id", id}}, opts)
	return result
}

func (d BaseMongoCRUD) all(ctx context.Context) (*mongo.Cursor, error) {
	opts := options.Find().SetSort(bson.D{{"age", 1}})
	cursor, err := d.coll.Find(ctx, bson.D{{}}, opts)
	return cursor, err
}

func (d BaseMongoCRUD) create(ctx context.Context) error {
	result, err := d.coll.InsertOne(ctx, d)
	if err != nil {
		log.Fatalf("Cannot create document in collection backlog:%s", err)
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		d.Id = oid
	} else {
		return errors.New("can't cast ObjectId to hex")
	}
	return err
}

func (d BaseMongoCRUD) update(ctx context.Context) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": bson.M{"$eq": d.Id}}
	update, err := bson.Marshal(d)
	if err != nil {
		return err
	}
	var data bson.M
	err = bson.Unmarshal(update, &data)
	if err != nil {
		return err
	}
	result, err := d.coll.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: data}}, opts)

	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("No matched documents")
	}
	if result.ModifiedCount == 0 {
		return errors.New("No modified document")
	}
	return nil
}

func (d BaseMongoCRUD) delete(ctx context.Context, id primitive.ObjectID) error {
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    "en_US",
		Strength:  1,
		CaseLevel: false,
	})
	res, err := d.coll.DeleteOne(ctx, bson.D{{"_id", d.Id}}, opts)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("no deleted document")
	}
	return nil
}
