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




func TestTask_Create(t *testing.T) {
	clt, err := Connect()
	if err !=nil{
		t.Errorf("Connection to Mongo database failed error: %s",err)
	}
	var taskDoc entity.Task
	err = faker.FakeData(&taskDoc)
	taskDoc.AuthorId = "6124b2da54569e9a7f2d51e2"
	taskDoc.ExecutorId = "6124b2da54569e9a7f2d51e2"
	taskRepo := new(Task)
	taskRepo.InitCollection(clt)
	b, err := taskRepo.Create(context.TODO(), taskDoc)
	if err!=nil{
		fmt.Printf("error in create function %s\n",err)
	}
	taskDoc.Id = b.Id
	fmt.Println("backlogDoc id "+ taskDoc.Id)
	fmt.Println("b id "+b.Id)
	if assert.NoError(t, err) {
		assert.NotNil(t, b, "function create doesn't create object or doesn't return ")
		assert.NotNil(t, taskDoc.Id, "Id of the inserted document not setted")
		testTask_Find(t,taskRepo,taskDoc)
		testTask_Update(t,taskRepo,taskDoc)
		testTask_All(t,taskRepo)
	} else {
		t.Errorf("Main function of Backlog_test throw error on creating document, error: %s", err)
	}
	clt.Disconnect(context.TODO())
}

func testTask_Find(t *testing.T, repository *Task, task entity.Task) {
	b, err := repository.Find(context.TODO(), task.Id)
	if assert.NoError(t, err) {
		assert.True(t, cmp.Equal(task,b),"User objects are not equal")
	} else {
		t.Errorf("Find function return error:%s", err)
	}
}

func testTask_Update(t *testing.T, repository *Task, task entity.Task) {
	task.Title = "Task test 2"
	task,err := repository.Update(context.TODO(), task)
	if err != nil {
		assert.Error(t, err)
	}
	b, err := repository.Find(context.TODO(), task.Id)
	if assert.NoError(t, err) {
		assert.Equal(t, task,b)
		assert.True(t, cmp.Equal(task,b),"User objects are not equal")
	} else {
		t.Errorf("Find function return error:%s", err)
	}
}

func testTask_All(t *testing.T, repository *Task) {
	collection, err := repository.All(context.TODO())
	assert.NoError(t, err, "List documents from mongodb return error")
	assert.GreaterOrEqual(t, len(collection), 1)
}

func testTask_Delete(t *testing.T, repository *Project) {
	err := repository.Delete(context.TODO(), backlogDoc.Id)
	assert.NoError(t, err, "Delete document from mongodb return error")
	doc, err := repository.Find(context.TODO(), backlogDoc.Id)
	assert.Nil(t, doc)
	assert.Error(t, err)
}


