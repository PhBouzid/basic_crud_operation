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
	"time"
)

type Backlog struct {
	coll             *mongo.Collection
	connectstring    string
	Id               primitive.ObjectID `bson:"_"`
	Title            string             `bson:"title"`
	BriefDescription string             `bson:"brief"`
	Content          string             `bson:"content"`
	TimeCreated      int64              `bson:"time_created"`
	User_Id          primitive.ObjectID `bson:"user_id"`
	Retro            []Retrospective    `bson:"retrospectives"`
}

type Retrospective struct {
	TaskId      primitive.ObjectID `bson:"task_id"`
	Information string             `bson:"information"`
	AuthorId    primitive.ObjectID `bson:"author"`
	TimeCreated int64              `bson:"time_created"`
}

func (d *Backlog) Create(ctx context.Context, backlog entity.Backlog) (entity.Backlog, error) {
	d.transformFrom(backlog)
	insertedResult, err := d.coll.InsertOne(ctx, d)
	if err != nil {
		log.Fatalf("Cannot create document in collection backlog:%s", err)
	}
	if oid, ok := insertedResult.InsertedID.(primitive.ObjectID); ok {
		d.Id = oid
	} else {
		return entity.Backlog{}, errors.New("can't cast ObjectId to hex")
	}
	result, err := d.transformTo()
	return result, err
}

func (d *Backlog) transformTo() (entity.Backlog, error) {
	useStrId := d.User_Id.Hex()
	idStrId := d.Id.Hex()
	ent := entity.Backlog{
		Id:               idStrId,
		Title:            d.Title,
		BriefDescription: d.BriefDescription,
		Content:          d.Content,
		TimeCreated:      time.Unix(0,d.TimeCreated*int64(time.Millisecond)),
		User_Id:          useStrId,
	}
	if len(d.Retro) == 0 {
		return ent, nil
	}
	retros := make([]entity.Retrospective, len(d.Retro))
	for _, value := range d.Retro {
		taskStrId := value.TaskId.Hex()
		authStrId := value.AuthorId.Hex()
		ret := entity.Retrospective{
			TaskId:      taskStrId,
			Information: value.Information,
			AuthorId:    authStrId,
			TimeCreated: time.Unix(0,d.TimeCreated*int64(time.Millisecond)),
		}
		retros = append(retros, ret)
	}
	ent.Retro = retros
	return ent, nil
}

func (d *Backlog) transformFrom(backlog entity.Backlog) {
	if len(backlog.Id) > 0 {
		objId, err := primitive.ObjectIDFromHex(backlog.Id)
		if err != nil {
			panic(err)
		}
		d.Id = objId
	}
	if len(backlog.User_Id) > 0 {
		userObjId, err := primitive.ObjectIDFromHex(backlog.User_Id)
		if err != nil {
			panic(err)
		}
		d.User_Id = userObjId
	}
	d.Title = backlog.Title
	d.BriefDescription = backlog.BriefDescription
	d.Content = backlog.Content
	d.TimeCreated = backlog.TimeCreated.UnixNano()/int64(time.Millisecond)
	if len(backlog.Retro) == 0 {
		return
	}
	retros := make([]Retrospective, len(backlog.Retro))
	for _, value := range backlog.Retro {
		taskObjId, err := primitive.ObjectIDFromHex(value.TaskId)
		if err != nil {
			log.Fatalf("Task object id didn't recognized, error: %s", err)
		}
		authObjId, err := primitive.ObjectIDFromHex(value.AuthorId)
		if err != nil {
			log.Fatalf("Author object id didn't recognized, error: %s", err)
		}

		r := Retrospective{
			Information: value.Information,
			TimeCreated: value.TimeCreated.UnixNano()/int64(time.Millisecond),
			TaskId:      taskObjId,
			AuthorId:    authObjId,
		}
		retros = append(retros, r)
	}
	d.Retro = retros
}

func (d *Backlog) DatabaseAddr(db string) {
	d.connectstring = db
}

func (d *Backlog) Find(ctx context.Context, id string) (entity.Backlog, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.Backlog{}, err
	}
	opts := options.FindOne().SetSort(bson.D{{"age", 1}})
	document := d.coll.FindOne(ctx, bson.D{{"_id", objId}}, opts)
	err = document.Decode(&d)
	result, err := d.transformTo()
	return result, err
}

func (d *Backlog) All(ctx context.Context) ([]entity.Backlog, error) {
	opts := options.Find().SetSort(bson.D{{"age", 1}})
	cursor, err := d.coll.Find(ctx, bson.D{{}}, opts)
	defer cursor.Close(context.TODO())
	if err != nil {
		log.Fatalf("Can't load cursor with all documents, error: %s", err)
	}
	var backlogs []Backlog
	err = cursor.All(ctx, &backlogs)
	if err != nil {
		log.Fatalf("Can't decode all documents, error: %s", err)
	}
	results := make([]entity.Backlog, len(backlogs))
	for _, value := range backlogs {
		b, err := value.transformTo()
		if err != nil {
			log.Fatalf("one or more documents can't be decode, error: %s", err)
		}
		results = append(results, b)
	}

	return results, nil
}

func (d *Backlog) Update(ctx context.Context, backlog entity.Backlog) (entity.Backlog, error) {
	d.transformFrom(backlog)
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": bson.M{"$eq": d.Id}}
	update, err := bson.Marshal(d)
	if err != nil {
		return entity.Backlog{}, err
	}
	var data bson.M
	err = bson.Unmarshal(update, &data)
	if err != nil {
		return entity.Backlog{}, err
	}
	result, err := d.coll.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: data}}, opts)

	if err != nil {
		return entity.Backlog{}, err
	}
	if result.MatchedCount == 0 {
		return entity.Backlog{}, errors.New("No matched documents")
	}
	if result.ModifiedCount == 0 {
		return entity.Backlog{}, errors.New("No modified document")
	}
	if err != nil {
		return entity.Backlog{}, err
	}
	entity, err := d.transformTo()
	return entity, err
}

func (d *Backlog) Delete(ctx context.Context, id string) error {
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

func (d *Backlog) InitCollection(clt *mongo.Client) {
	d.coll = clt.Database("scrumDocs").Collection("backlog")
}
