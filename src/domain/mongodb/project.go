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

type Project struct {
	connectstring string
	coll          *mongo.Collection
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `bson:"name"`
	Weeks         int                `bson:"weeks"`
	StartDate     int64              `bson:"start_date"`
	WorkDuration  int64              `bson:"work_duration"`
	AuthorId      primitive.ObjectID `bson:"author"`
}

func (d *Project) transformTo() entity.Project {
	idStrId := d.Id.Hex()
	authorStrId := d.AuthorId.Hex()
	ent := entity.Project{
		Id:           idStrId,
		Name:         d.Name,
		Weeks:        d.Weeks,
		StartDate:    time.Unix(0,d.StartDate*int64(time.Millisecond)),
		WorkDuration: d.WorkDuration,
		AuthorId:     authorStrId,
	}
	return ent
}

func (d *Project) transformFrom(project entity.Project) {
	if len(project.Id) > 0 {
		objId, err := primitive.ObjectIDFromHex(project.Id)
		if err != nil {
			panic(err)
		}
		d.Id = objId
	}
	if len(project.AuthorId) > 0 {
		userObjId, err := primitive.ObjectIDFromHex(project.AuthorId)
		if err != nil {
			panic(err)
		}
		d.AuthorId = userObjId
	}
	d.Name = project.Name
	d.Weeks = project.Weeks
	d.StartDate = project.StartDate.UnixNano()/int64(time.Millisecond)
	d.WorkDuration = project.WorkDuration
}

func (d *Project) DatabaseAddr(db string) {
	d.connectstring = db
}

func (d *Project) Create(ctx context.Context, project entity.Project) (entity.Project, error) {
	d.transformFrom(project)
	insertedResult, err := d.coll.InsertOne(ctx, d)
	if err != nil {
		log.Fatalf("Cannot create document in collection backlog:%s", err)
	}
	if oid, ok := insertedResult.InsertedID.(primitive.ObjectID); ok {
		d.Id = oid
	} else {
		return entity.Project{}, errors.New("can't cast ObjectId to hex")
	}
	result := d.transformTo()
	return result, err
}

func (d *Project) Find(ctx context.Context, id string) (entity.Project, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.Project{}, err
	}
	opts := options.FindOne().SetSort(bson.D{{"age", 1}})
	document := d.coll.FindOne(ctx, bson.D{{"_id", objId}}, opts)
	err = document.Decode(&d)
	result := d.transformTo()
	return result, err
}

func (d *Project) All(ctx context.Context) ([]entity.Project, error) {
	opts := options.Find().SetSort(bson.D{{"age", 1}})
	cursor, err := d.coll.Find(ctx, bson.D{{}}, opts)
	defer cursor.Close(context.TODO())
	if err != nil {
		log.Fatalf("Can't load cursor with all documents, error: %s", err)
	}
	var backlogs []Project
	err = cursor.All(ctx, &backlogs)
	if err != nil {
		log.Fatalf("Can't decode all documents, error: %s", err)
	}
	results := make([]entity.Project, len(backlogs))
	for _, value := range backlogs {
		b := value.transformTo()
		results = append(results, b)
	}

	return results, nil
}

func (d *Project) Update(ctx context.Context, project entity.Project) (entity.Project, error) {
	d.transformFrom(project)
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": bson.M{"$eq": d.Id}}
	update, err := bson.Marshal(d)
	if err != nil {
		return entity.Project{}, err
	}
	var data bson.M
	err = bson.Unmarshal(update, &data)
	if err != nil {
		return entity.Project{}, err
	}
	result, err := d.coll.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: data}}, opts)

	if err != nil {
		return entity.Project{}, err
	}
	if result.MatchedCount == 0 {
		return entity.Project{}, errors.New("No matched documents")
	}
	if result.ModifiedCount == 0 {
		return entity.Project{}, errors.New("No modified document")
	}
	if err != nil {
		return entity.Project{}, err
	}
	entity := d.transformTo()
	return entity, err
}

func (d *Project) Delete(ctx context.Context, id string) error {
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

func (d *Project) InitCollection(clt *mongo.Client) {
	d.coll = clt.Database("scrumDocs").Collection("projects")
}
