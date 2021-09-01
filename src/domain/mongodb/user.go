package mongodb

import (
	"BlockDoc/entity"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type User struct {
	coll             *mongo.Collection
	connectstring string
	Id               primitive.ObjectID `bson:"_id,omitempty"`
	UserName         string             `bson:"user_name"`
	FirstName        string             `bson:"first_name"`
	SecondName       string             `bson:"second_name"`
	Email            string             `bson:"email"`
	password         string             `bson:"password"`
}

func (d *User) transformTo() entity.User {
	idStrId := d.Id.Hex()
	ent := entity.User{
		Id:         idStrId,
		UserName:   d.UserName,
		FirstName:  d.FirstName,
		SecondName: d.SecondName,
		Email:      d.Email,
		Password:   d.password,
	}
	return ent
}

func (d *User) transformFrom(user entity.User) {
	if len(user.Id) > 0 {
		objId, err := primitive.ObjectIDFromHex(user.Id)
		if err != nil {
			log.Fatalf("User Id is not recognized error:%s, id:%s",err,user.Id)
		}
		d.Id = objId
	}

	d.UserName = user.UserName
	d.FirstName = user.FirstName
	d.SecondName = user.SecondName
	d.Email = user.Email
	d.password = user.Password
}

func (d *User) DatabaseAddr(db string) {
	d.connectstring = db
}


func (d *User) Create(ctx context.Context, user entity.User) (entity.User, error) {
	d.transformFrom(user)
	insertedResult, err := d.coll.InsertOne(ctx, d)
	if err != nil {
		log.Fatalf("Cannot create document in collection backlog:%s", err)
	}
	if oid, ok := insertedResult.InsertedID.(primitive.ObjectID); ok {
		d.Id = oid
	} else {
		return entity.User{}, errors.New("can't cast ObjectId to hex")
	}
	result := d.transformTo()
	return result, err
}

func (d *User) Find(ctx context.Context, id string) (entity.User, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.User{}, err
	}
	opts := options.FindOne().SetSort(bson.D{{"age", 1}})
	document := d.coll.FindOne(ctx, bson.D{{"_id", objId}}, opts)
	err = document.Decode(&d)
	result := d.transformTo()
	return result, err
}

func (d *User) All(ctx context.Context) ([]entity.User, error) {
	opts := options.Find().SetSort(bson.D{{"age", 1}})
	cursor, err := d.coll.Find(ctx, bson.D{{}}, opts)
	defer cursor.Close(context.TODO())
	if err != nil {
		log.Fatalf("Can't load cursor with all documents, error: %s", err)
	}
	var users []User
	err = cursor.All(ctx, &users)
	if err != nil {
		log.Fatalf("Can't decode all documents, error: %s", err)
	}
	results := make([]entity.User, len(users))
	for _, value := range users {
		b := value.transformTo()
		results = append(results, b)
	}

	return results, nil
}

func (d *User) Update(ctx context.Context, user entity.User) (entity.User, error) {
	d.transformFrom(user)
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": bson.M{"$eq": d.Id}}
	update, err := bson.Marshal(d)
	if err != nil {
		return entity.User{}, err
	}
	var data bson.M
	err = bson.Unmarshal(update, &data)
	if err != nil {
		return entity.User{}, err
	}
	result, err := d.coll.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: data}}, opts)

	if err != nil {
		return entity.User{}, err
	}
	if result.MatchedCount == 0 {
		return entity.User{}, errors.New("No matched documents")
	}
	if result.ModifiedCount == 0 {
		return entity.User{}, errors.New("No modified document")
	}
	if err != nil {
		return entity.User{}, err
	}
	entity := d.transformTo()
	return entity, err
}

func (d *User) Delete(ctx context.Context, id string) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    "en_US",
		Strength:  1,
		CaseLevel: false,
	})
	res, err := d.coll.DeleteOne(ctx, bson.D{{"_id", objId}}, opts)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("no deleted document")
	}
	return err
}

func (d *User) InitCollection(clt *mongo.Client) {
	d.coll = clt.Database("scrumDocs").Collection("users")
}
