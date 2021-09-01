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

type Task struct {
	coll          *mongo.Collection
	connectstring string
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	Title         string             `bson:"title"`
	Brief         string             `bson:"brief"`
	Content       string             `bson:"content"`
	State         int                `bson:"state"`
	ExecutorId    primitive.ObjectID `bson:"executor_id"`
	AuthorId      primitive.ObjectID `bson:"author_id"`
	CreatedAt     int64          `bson:"created_at"`
	UpdatedAt     int64          `bson:"updated_at"`
	Prioritet     int                `bson:"prioritet"`
	Action        []ActionStory      `bson:"action"`
	Achievement   []Achievements     `bson:"achievements"`
}

type ActionStory struct {
	Action      string             `bson:"action"`
	CreatedTime int64          `bson:"created_time"`
	UpdatedTime int64          `bson:"updated_time"`
	AuthorId    primitive.ObjectID `bson:"author_id,omitempty" faker:"-"`
}

type Achievements struct {
	Achievment  string             `bson:"achievment"`
	CreatedTime int64          `bson:"created_time"`
	UpdatedTime int64          `bson:"updated_time"`
	AuthorId    primitive.ObjectID `bson:"author_id,omitempty" faker:"-"`
}

func (d *Task) transformTo() entity.Task {
	idStrId := d.Id.Hex()
	authorStrId := d.AuthorId.Hex()
	executorStrId := d.ExecutorId.Hex()
	ent := entity.Task{
		Id:         idStrId,
		Title:      d.Title,
		Brief:      d.Brief,
		Content:    d.Content,
		State:      d.State,
		AuthorId:   authorStrId,
		ExecutorId: executorStrId,
		Prioritet:  d.Prioritet,
		CreatedAt:  time.Unix(0,d.CreatedAt * int64(time.Millisecond)),
		UpdatedAt:  time.Unix(0,d.UpdatedAt * int64(time.Millisecond)),
	}
	if (len(d.Achievement) == 0) && (len(d.Action) == 0) {
		return ent
	}
	if len(d.Achievement)>0 {
	achievements := make([]entity.Achievements, len(d.Achievement))
	for _, value := range d.Achievement {
		var strId string
		if !value.AuthorId.IsZero() {
			strId = value.AuthorId.Hex()
		}else{
			strId = ""
		}
			achieve := entity.Achievements{
				Achievment:  value.Achievment,
				CreatedTime: time.Unix(0,value.CreatedTime * int64(time.Millisecond)),
				UpdatedTime: time.Unix(0,value.UpdatedTime * int64(time.Millisecond)),
				AuthorId:    strId,
			}
			achievements = append(achievements, achieve)
		}
		ent.Achievement = achievements
	}

	if len(d.Action)>0 {
		actions := make([]entity.ActionStory, len(d.Action))
		for _, value := range d.Action {
			var strId string
			if !value.AuthorId.IsZero() {
				strId = value.AuthorId.Hex()
			} else {
				strId = ""
			}

			action := entity.ActionStory{
				Action:      value.Action,
				CreatedTime: time.Unix(0,value.CreatedTime * int64(time.Millisecond)),
				UpdatedTime: time.Unix(0,value.UpdatedTime * int64(time.Millisecond)),
				AuthorId:    strId,
			}
			actions = append(actions, action)
		}
		ent.Action = actions
	}
	return ent
}

func (d *Task) transformFrom(task entity.Task) {
	if len(task.Id) > 0 {
		objId, err := primitive.ObjectIDFromHex(task.Id)
		if err != nil {
			log.Fatalf("Task id  didn't recognized, error: %s, id: %s", err, objId)
		}
		d.Id = objId
	}
	if len(task.AuthorId) > 0 {
		objId, err := primitive.ObjectIDFromHex(task.AuthorId)
		if err != nil {
			log.Fatalf("Task Author id  didn't recognized, error: %s, id: %s", err, objId)
		}
		d.AuthorId = objId
	}
	if len(task.ExecutorId) > 0 {
		objId, err := primitive.ObjectIDFromHex(task.ExecutorId)
		if err != nil {
			log.Fatalf("Task Executor id  didn't recognized, error: %s, id: %s", err, objId)
		}
		d.ExecutorId = objId
	}
	d.UpdatedAt = task.UpdatedAt.UnixNano()/int64(time.Millisecond)
	d.CreatedAt = task.CreatedAt.UnixNano()/int64(time.Millisecond)
	d.Brief = task.Brief
	d.Title = task.Title
	d.Content = task.Content
	d.State = task.State
	d.Prioritet = task.Prioritet

	achievements := make([]Achievements, len(task.Achievement))
	for _, value := range task.Achievement {
		var achieve Achievements
		if len(value.AuthorId)>0{
			objId, err := primitive.ObjectIDFromHex(value.AuthorId)
			if err != nil {
				log.Fatalf("Task Author object id didn't recognized, error: %s", err)
			}
			achieve.AuthorId = objId
		}
		achieve.Achievment = value.Achievment
		achieve.CreatedTime = value.CreatedTime.UnixNano()/int64(time.Millisecond)
		achieve.UpdatedTime = value.UpdatedTime.UnixNano()/int64(time.Millisecond)
		achievements = append(achievements, achieve)
	}
	d.Achievement = achievements
	actions := make([]ActionStory, len(task.Action))
	for _, value := range task.Action {
		var action ActionStory
		if len(value.AuthorId)>0{
			objId, err := primitive.ObjectIDFromHex(value.AuthorId)
			if err != nil {
				log.Fatalf("Task Author object id didn't recognized, error: %s", err)
			}
			action.AuthorId = objId
		}
		action.Action = value.Action
		action.UpdatedTime = value.UpdatedTime.UnixNano()/int64(time.Millisecond)
		action.CreatedTime = value.CreatedTime.UnixNano()/int64(time.Millisecond)

		actions = append(actions, action)
	}
	d.Action = actions
}

func (d *Task) DatabaseAddr(db string) {
	d.connectstring = db
}

func (d *Task) Create(ctx context.Context, task entity.Task) (entity.Task, error) {
	d.transformFrom(task)
	insertedResult, err := d.coll.InsertOne(ctx, d)
	if err != nil {
		log.Fatalf("Cannot create document in collection backlog:%s", err)
	}
	if oid, ok := insertedResult.InsertedID.(primitive.ObjectID); ok {
		d.Id = oid
	} else {
		return entity.Task{}, errors.New("can't cast ObjectId to hex")
	}
	result := d.transformTo()
	return result, err
}

func (d *Task) Find(ctx context.Context, id string) (entity.Task, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.Task{}, err
	}
	opts := options.FindOne().SetSort(bson.D{{"age", 1}})
	document := d.coll.FindOne(ctx, bson.D{{"_id", objId}}, opts)
	err = document.Decode(&d)
	result := d.transformTo()
	return result, err
}

func (d *Task) All(ctx context.Context) ([]entity.Task, error) {
	opts := options.Find().SetSort(bson.D{{"age", 1}})
	cursor, err := d.coll.Find(ctx, bson.D{{}}, opts)
	defer cursor.Close(context.TODO())
	if err != nil {
		log.Fatalf("Can't load cursor with all documents, error: %s", err)
	}
	var tasks []Task
	err = cursor.All(ctx, &tasks)
	if err != nil {
		log.Fatalf("Can't decode all documents, error: %s", err)
	}
	results := make([]entity.Task, len(tasks))
	for _, value := range tasks {
		b := value.transformTo()
		results = append(results, b)
	}

	return results, nil
}

func (d *Task) Update(ctx context.Context, task entity.Task) (entity.Task, error) {
	d.transformFrom(task)
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": bson.M{"$eq": d.Id}}
	update, err := bson.Marshal(d)
	if err != nil {
		return entity.Task{}, err
	}
	var data bson.M
	err = bson.Unmarshal(update, &data)
	if err != nil {
		return entity.Task{}, err
	}
	result, err := d.coll.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: data}}, opts)

	if err != nil {
		return entity.Task{}, err
	}
	if result.MatchedCount == 0 {
		return entity.Task{}, errors.New("No matched documents")
	}
	if result.ModifiedCount == 0 {
		return entity.Task{}, errors.New("No modified document")
	}
	if err != nil {
		return entity.Task{}, err
	}
	entity := d.transformTo()
	return entity, err
}

func (d *Task) Delete(ctx context.Context, id string) error {
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

func (d *Task) InitCollection(clt *mongo.Client) {
	d.coll = clt.Database("scrumDocs").Collection("tasks")
}
