package mongodb

import (
	"BlockDoc/entity"
	"context"
	"errors"
	_ "github.com/bxcodec/faker/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type UserStory struct {
	coll *mongo.Collection
	connectstring string
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Brief     string             `bson:"brief"`
	Content   string             `bson:"content"`
	Tasks     []TaskShort             `bson:"tasks"`
	CreatedAt int64          `bson:"created_at"`
	UpdatedAt int64          `bson:"updated_at"`
}

type TaskShort struct {
	TaskId     primitive.ObjectID `bson:"taskId"`
	Brief      string             `bson:"brief"`
	State      int                `bson:"state"`
	ExecutorId primitive.ObjectID `bson:"executor_id"`
}

func (d *UserStory) transformTo() entity.UserStory {
	idStrId := d.Id.Hex()
	ent := entity.UserStory{
		Id:           idStrId,
		Title:         d.Title,
		Brief:        d.Brief,
		Content:    d.Content,
		CreatedAt: time.Unix(0,d.CreatedAt * int64(time.Millisecond)),
		UpdatedAt:  time.Unix(0,d.UpdatedAt * int64(time.Millisecond)),
	}
	if len(d.Tasks) == 0 {
		return ent
	}
	tasks := make([]entity.TaskShort, len(d.Tasks))
	for _, value := range d.Tasks {
		taskStrId := value.TaskId.Hex()
		executorStrId := value.ExecutorId.Hex()
		task := entity.TaskShort{
			TaskId:      taskStrId,
			Brief: value.Brief,
			ExecutorId:    executorStrId,
			State: value.State,
		}
		tasks = append(tasks, task)
	}
	ent.Tasks = tasks
	return ent
}

func (d *UserStory) transformFrom(story entity.UserStory) {
	if len(story.Id) > 0 {
		objId, err := primitive.ObjectIDFromHex(story.Id)
		if err != nil {
			panic(err)
		}
		d.Id = objId
	}

	d.UpdatedAt = story.UpdatedAt.UnixNano()/int64(time.Millisecond)
	d.CreatedAt = story.CreatedAt.UnixNano()/int64(time.Millisecond)
	d.Brief = story.Brief
	d.Title = story.Title
	d.Content = story.Content
	tasks := make([]TaskShort, len(d.Tasks))
	for _, value := range story.Tasks {
		task := TaskShort{}
		if len(value.TaskId)>0 {
			taskObjId, err := primitive.ObjectIDFromHex(value.TaskId)
			if err != nil {
				log.Fatalf("User story documet has task object id that didn't recognized, error: %s, id: %s", err, value.TaskId)
			}
			task.TaskId = taskObjId
		}
		if len(value.ExecutorId)>0 {
			executorObjId, err := primitive.ObjectIDFromHex(value.ExecutorId)
			if err != nil {
				log.Fatalf("User story document has author object id that didn't recognized, error: %s, id: %s", err, value.ExecutorId)
			}
			task.ExecutorId = executorObjId
		}
		task.State = value.State
		task.Brief = value.Brief
		tasks = append(tasks, task)
	}
	d.Tasks = tasks
}

func (d *UserStory) DatabaseAddr(db string) {
	d.connectstring = db
}



func (d *UserStory) Create(ctx context.Context, story entity.UserStory) (entity.UserStory, error) {
	d.transformFrom(story)
	insertedResult, err := d.coll.InsertOne(ctx, d)
	if err != nil {
		log.Fatalf("Cannot create document in collection backlog:%s", err)
	}
	if oid, ok := insertedResult.InsertedID.(primitive.ObjectID); ok {
		d.Id = oid
	} else {
		return entity.UserStory{}, errors.New("can't cast ObjectId to hex")
	}
	result := d.transformTo()
	return result, err
}

func (d *UserStory) Find(ctx context.Context, id string) (entity.UserStory, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.UserStory{}, err
	}
	opts := options.FindOne().SetSort(bson.D{{"age", 1}})
	document := d.coll.FindOne(ctx, bson.D{{"_id", objId}}, opts)
	err = document.Decode(&d)
	result := d.transformTo()
	return result, err
}

func (d *UserStory) All(ctx context.Context) ([]entity.UserStory, error) {
	opts := options.Find().SetSort(bson.D{{"age", 1}})
	cursor, err := d.coll.Find(ctx, bson.D{{}}, opts)
	defer cursor.Close(context.TODO())
	if err != nil {
		log.Fatalf("Can't load cursor with all documents, error: %s", err)
	}
	var stories []UserStory
	err = cursor.All(ctx, &stories)
	if err != nil {
		log.Fatalf("Can't decode all documents, error: %s", err)
	}
	results := make([]entity.UserStory, len(stories))
	for _, value := range stories {
		b := value.transformTo()
		results = append(results, b)
	}

	return results, nil
}

func (d *UserStory) Update(ctx context.Context, story entity.UserStory) (entity.UserStory, error) {
	d.transformFrom(story)
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": bson.M{"$eq": d.Id}}
	update, err := bson.Marshal(d)
	if err != nil {
		return entity.UserStory{}, err
	}
	var data bson.M
	err = bson.Unmarshal(update, &data)
	if err != nil {
		return entity.UserStory{}, err
	}
	result, err := d.coll.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: data}}, opts)

	if err != nil {
		return entity.UserStory{}, err
	}
	if result.MatchedCount == 0 {
		return entity.UserStory{}, errors.New("No matched documents")
	}
	if result.ModifiedCount == 0 {
		return entity.UserStory{}, errors.New("No modified document")
	}
	if err != nil {
		return entity.UserStory{}, err
	}
	entity := d.transformTo()
	return entity, err
}

func (d *UserStory) Delete(ctx context.Context, id string) error {
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

func (d *UserStory) InitCollection(clt *mongo.Client) {
	d.coll = clt.Database("scrumDocs").Collection("stories")
}
