package mongodb

import (
	"BlockDoc/entity"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"testing"
	"time"
)

var backlogDoc = entity.Backlog{
	BriefDescription: "test test test",
	Title:            "TEST",
	Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Praesent blandit elit congue, pulvinar turpis eu, interdum est.",
	User_Id: "6124b2da54569e9a7f2d51e2",
}

var retrospectives = []entity.Retrospective{{
	Information: "Test information",
	AuthorId: "6124b2da54569e9a7f2d51e1",
},{
	Information: "Test information 3",
	AuthorId: "6124b2da54569e9a7f2d51e5",
}}

var backlogDoc_v2 = entity.Backlog{
	BriefDescription: "test 2 test 2 test 2",
	Title:            "TEST 2",
	Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Praesent blandit elit congue, pulvinar turpis eu, interdum est.",
	User_Id: "6124b2da54569e9a7f2d51e8",
	Retro:   retrospectives,
}


func Connect() (*mongo.Client,error){
	clt, err := mongo.NewClient(options.Client().ApplyURI("mongodb://admin:12345@127.0.0.1:27017/"))
	if err!=nil{
		//log.Fatalf("Build connection client for mongo failed with error: %s",err)
		return nil,err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = clt.Connect(ctx)
	if err != nil {
		log.Printf("Connection not establish error: %s", err)
		return nil,err
	}
	err = clt.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Printf("Connection not establish error: %s", err)
		return nil,err
	}
	return clt,nil
}


func TestBacklog_Create(t *testing.T) {
	clt, err := Connect()
	if err !=nil{
		t.Errorf("Connection to Mongo database failed error: %s",err)
	}
	backlogRepo := new(Backlog)
	backlogRepo.InitCollection(clt)
	b, err := backlogRepo.Create(context.TODO(), backlogDoc)
	if err!=nil{
		fmt.Printf("error in create function %s\n",err)
	}
	backlogDoc.Id = b.Id
	if assert.NoError(t, err) {
		assert.NotNil(t, b, "function create doesn't create object or doesn't return ")
		assert.NotNil(t, backlogDoc.Id, "Id of the inserted document not setted")
		testBacklog_Find(t,backlogRepo)
		testBacklog_Update(t,backlogRepo)
		testBacklog_All(t,backlogRepo)
	} else {
		t.Errorf("Main function of Backlog_test throw error on creating document, error: %s", err)
	}
	clt.Disconnect(context.TODO())
}

func testBacklog_Find(t *testing.T, repository *Backlog) {
	b, err := repository.Find(context.TODO(), backlogDoc.Id)
	if assert.NoError(t, err) {
		assert.Equalf(t, backlogDoc.Title, b.Title, "backlog title is not the same after updated them expected %s result is %s", backlogDoc.Title, b.Title)
		assert.Equalf(t, backlogDoc.BriefDescription, b.BriefDescription, "backlog brief description is not changing by update function expected %s result is %s", backlogDoc.BriefDescription, b.BriefDescription)
		assert.Equalf(t, backlogDoc.User_Id, b.User_Id, "user id of document is not changed after update method run, expected %d result is %d ", backlogDoc.User_Id, b.User_Id)
		assert.Equalf(t, backlogDoc.Content, b.Content, "Content of document is not changed after update method run, expected %s result is %s ", backlogDoc.Content, b.Content)
	} else {
		t.Errorf("Find function return error:%s", err)
	}
}

func testBacklog_Update(t *testing.T, repository *Backlog) {
	backlogDoc.Title = "TITLE TEST 2"
	backlogDoc,err := repository.Update(context.TODO(), backlogDoc)
	if err != nil {
		assert.Error(t, err)
	}
	b, err := repository.Find(context.TODO(), backlogDoc.Id)
	if assert.NoError(t, err) {
		assert.Equalf(t, backlogDoc.Title, b.Title, "backlog title is not the same after updated them expected %s result is %s", backlogDoc.Title, b.Title)
		assert.Equalf(t, backlogDoc.BriefDescription, b.BriefDescription, "backlog brief description is not changing by update function expected %s result is %s", backlogDoc.BriefDescription, b.BriefDescription)
		assert.Equalf(t, backlogDoc.User_Id, b.User_Id, "user id of document is not changed after update method run, expected %d result is %d ", backlogDoc.User_Id, b.User_Id)
		assert.Equalf(t, backlogDoc.Content, b.Content, "Content of document is not changed after update method run, expected %s result is %s ", backlogDoc.Content, b.Content)
	} else {
		t.Errorf("Find function return error:%s", err)
	}
}

func testBacklog_All(t *testing.T, repository *Backlog) {
	collection, err := repository.All(context.TODO())
	assert.NoError(t, err, "List documents from mongodb return error")
	assert.GreaterOrEqual(t, len(collection), 1)
}

func testBacklog_Delete(t *testing.T, repository *Backlog) {
	err := repository.Delete(context.TODO(), backlogDoc.Id)
	assert.NoError(t, err, "Delete document from mongodb return error")
	doc, err := repository.Find(context.TODO(), backlogDoc.Id)
	assert.Nil(t, doc)
	assert.Error(t, err)
}
