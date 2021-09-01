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


//AuthorId: "6124b2da54569e9a7f2d51e2",


func TestProject_Create(t *testing.T) {
	clt, err := Connect()
	if err !=nil{
		t.Errorf("Connection to Mongo database failed error: %s",err)
	}
	var projectDoc entity.Project
	err = faker.FakeData(&projectDoc)
	if err!=nil{
		t.Errorf("Error in generation fake data %s ",err)
	}
	projectDoc.Id=""
	projectDoc.AuthorId = "6124b2da54569e9a7f2d51e2"
	projectRepo := new(Project)
	projectRepo.InitCollection(clt)
	b, err := projectRepo.Create(context.TODO(), projectDoc)
	if err!=nil{
		fmt.Printf("error in create function %s\n",err)
	}
	projectDoc.Id = b.Id
	if assert.NoError(t, err) {
		assert.NotNil(t, b, "function create doesn't create object or doesn't return ")
		assert.NotNil(t, backlogDoc.Id, "Id of the inserted document not setted")
		testProject_Find(t,projectRepo, projectDoc)
		testProject_Update(t,projectRepo, projectDoc)
		testProject_All(t,projectRepo)
	} else {
		t.Errorf("Main function of Backlog_test throw error on creating document, error: %s", err)
	}
	clt.Disconnect(context.TODO())
}

func testProject_Find(t *testing.T, repository *Project, project entity.Project) {
	b, err := repository.Find(context.TODO(), project.Id)
	if assert.NoError(t, err) {
		//assert.True(t, cmp.Equal(b,project),"Projects are not equal")
		if diff := cmp.Diff(project,b);diff!=""{
			t.Log(diff)
		}
	} else {
		t.Errorf("Find function return error:%s", err)
	}
}

func testProject_Update(t *testing.T, repository *Project,project entity.Project) {
	project.Name = "Name test"
	project,err := repository.Update(context.TODO(), project)
	if err != nil {
		assert.Error(t, err)
	}
	b, err := repository.Find(context.TODO(), backlogDoc.Id)
	if assert.NoError(t, err) {
		assert.True(t, cmp.Equal(project,b),"Projects not equal ")
	} else {
		t.Errorf("Find function return error:%s", err)
	}
}

func testProject_All(t *testing.T, repository *Project) {
	collection, err := repository.All(context.TODO())
	assert.NoError(t, err, "List documents from mongodb return error")
	assert.GreaterOrEqual(t, len(collection), 1)
}

func testProject_Delete(t *testing.T, repository *Project) {
	err := repository.Delete(context.TODO(), backlogDoc.Id)
	assert.NoError(t, err, "Delete document from mongodb return error")
	doc, err := repository.Find(context.TODO(), backlogDoc.Id)
	assert.Nil(t, doc)
	assert.Error(t, err)
}

