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

func TestUser_Create(t *testing.T) {
	clt, err := Connect()
	if err !=nil{
		t.Errorf("Connection to Mongo database failed error: %s",err)
	}
	userRepo := new(User)
	userRepo.InitCollection(clt)
	var userDoc entity.User
	err = faker.FakeData(&userDoc)
	if err!=nil{
		t.Errorf("error in generation fake data")
	}
	b, err := userRepo.Create(context.TODO(), userDoc)
	if err!=nil{
		fmt.Printf("error in create function %s\n",err)
	}
	userDoc.Id = b.Id
	if assert.NoError(t, err) {
		assert.NotNil(t, b, "function create doesn't create object or doesn't return ")
		assert.NotNil(t, userDoc.Id, "Id of the inserted document not setted")
		testUser_Find(t,userRepo, userDoc)
		testUser_Update(t,userRepo,userDoc)
		testUser_All(t,userRepo)
	} else {
		t.Errorf("Main function of Backlog_test throw error on creating document, error: %s", err)
	}
	clt.Disconnect(context.TODO())
}

func testUser_Find(t *testing.T, repository *User, user entity.User) {
	b, err := repository.Find(context.TODO(), user.Id)
	if assert.NoError(t, err) {
		assert.Equal(t, user,b)
		assert.True(t, cmp.Equal(user,b),"User objects are not equal")
	} else {
		t.Errorf("Find function return error:%s", err)
	}
}

func testUser_Update(t *testing.T, repository *User,user entity.User) {
	user.UserName = "User_test"
	user,err := repository.Update(context.TODO(), user)
	if err != nil {
		assert.Error(t, err)
	}
	b, err := repository.Find(context.TODO(), user.Id)
	if assert.NoError(t, err) {
		assert.True(t, cmp.Equal(user,b),"User objects are not equal")
	} else {
		t.Errorf("Find function return error:%s", err)
	}
}

func testUser_All(t *testing.T, repository *User) {
	collection, err := repository.All(context.TODO())
	assert.NoError(t, err, "List documents from mongodb return error")
	assert.GreaterOrEqual(t, len(collection), 1)
}

func testUser_Delete(t *testing.T, repository *User, user entity.User) {
	err := repository.Delete(context.TODO(), user.Id)
	assert.NoError(t, err, "Delete document from mongodb return error")
	doc, err := repository.Find(context.TODO(), user.Id)
	assert.Nil(t, doc)
	assert.Error(t, err)
}


