package mongodb

import (
	"BlockDoc/entity"
	"context"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"testing"
)




func TestUserStory_Create(t *testing.T) {
	clt, err := Connect()
	if err !=nil{
		t.Errorf("Connection to Mongo database failed error: %s",err)
	}
	userStoryRepo := new(UserStory)
	userStoryRepo.InitCollection(clt)
	var userStoryDoc entity.UserStory
	err = faker.FakeData(&userStoryDoc)
	userStoryDoc.Id = ""
	b, err := userStoryRepo.Create(context.TODO(), userStoryDoc)
	if err!=nil{
		fmt.Printf("error in create function %s\n",err)
	}
	userStoryDoc.Id = b.Id
	if assert.NoError(t, err) {
		assert.NotNil(t, b, "function create doesn't create object or doesn't return ")
		assert.NotNil(t, userStoryDoc.Id, "Id of the inserted document not setted")
		testUserStory_Find(t,userStoryRepo, userStoryDoc)
		testUserStory_Update(t,userStoryRepo, userStoryDoc)
		testUserStory_All(t,userStoryRepo, userStoryDoc)
	} else {
		t.Errorf("Main function of Backlog_test throw error on creating document, error: %s", err)
	}
	clt.Disconnect(context.TODO())
}

func testUserStory_Find(t *testing.T, repository *UserStory, storyDoc entity.UserStory) {
	b, err := repository.Find(context.TODO(), storyDoc.Id)
	if assert.NoError(t, err) {
		assert.True(t, cmp.Equal(storyDoc,b),"UserStory objects are not equal")
	} else {
		t.Errorf("Find function return error:%s", err)
	}
}

func testUserStory_Update(t *testing.T, repository *UserStory, storyDoc entity.UserStory) {
	storyDoc.Title = "TITLE TEST 2"
	storyDoc,err := repository.Update(context.TODO(), storyDoc)
	if err != nil {
		assert.Error(t, err)
	}
	b, err := repository.Find(context.TODO(), storyDoc.Id)
	if assert.NoError(t, err) {
		assert.True(t, cmp.Equal(storyDoc,b),"User objects are not equal")
	} else {
		t.Errorf("Find function return error:%s", err)
	}
}

func testUserStory_All(t *testing.T, repository *UserStory, storyDoc entity.UserStory) {
	collection, err := repository.All(context.TODO())
	assert.NoError(t, err, "List documents from mongodb return error")
	assert.GreaterOrEqual(t, len(collection), 1)
}

func testUserStory_Delete(t *testing.T, repository *UserStory,storyDoc entity.UserStory) {
	err := repository.Delete(context.TODO(), backlogDoc.Id)
	assert.NoError(t, err, "Delete document from mongodb return error")
	doc, err := repository.Find(context.TODO(), backlogDoc.Id)
	assert.Nil(t, doc)
	assert.Error(t, err)
}

