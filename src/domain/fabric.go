package domain

import (
	"BlockDoc/domain/mongodb"
	"strings"
)

type RepositoryFactory interface {
	MakeProject() ProjectRepository
	MakeBacklog() BacklogRepository
	MakeUserStory() UserStoryRepository
	MakeTask() TaskRepository
	MakeUser() UserRepository
}

func GetRepositoryFactory(dbType, repository string) RepositoryFactory {
	switch dbType {
	case strings.ToLower("mongodb"):
		return MongoRepositoryFactory{}
	case "postgre":
	default:
		return MongoRepositoryFactory{}
	}
	return nil
}


type MongoRepositoryFactory struct{
	databaseAddr string
}



func (m MongoRepositoryFactory) MakeProject() ProjectRepository{
	repo := new(mongodb.Project)
	repo.DatabaseAddr("mongodb://admin:12345@127.0.0.1:27017")
	return repo
}

func (m MongoRepositoryFactory) MakeBacklog() BacklogRepository{
	repo := new(mongodb.Backlog)
	repo.DatabaseAddr("mongodb://admin:12345@127.0.0.1:27017")
	return repo
}

func (m MongoRepositoryFactory) MakeUserStory() UserStoryRepository{
	repo := new(mongodb.UserStory)
	repo.DatabaseAddr("mongodb://admin:12345@127.0.0.1:27017")
	return repo
}

func (m MongoRepositoryFactory) MakeTask() TaskRepository{
	repo := new(mongodb.Task)
	repo.DatabaseAddr("mongodb://admin:12345@127.0.0.1:27017")
	return repo
}

func (m MongoRepositoryFactory) MakeUser() UserRepository{
	repo := new(mongodb.User)
	repo.DatabaseAddr("mongodb://admin:12345@127.0.0.1:27017")
	return repo
}


