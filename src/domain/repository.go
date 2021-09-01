package domain

import (
	"BlockDoc/entity"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	DatabaseAddr(db string)
}

type ProjectRepository interface {
	Find(ctx context.Context, id string) (entity.Project,error)
	All(ctx context.Context)([]entity.Project, error)
	Create(ctx context.Context, project entity.Project)(entity.Project,error)
	Update(ctx context.Context, project entity.Project)(entity.Project,error)
	Delete(ctx context.Context, id string)(error)
}

type BacklogRepository interface {
	Find(ctx context.Context, id string) (entity.Backlog,error)
	All(ctx context.Context)([]entity.Backlog, error)
	Create(ctx context.Context,backlog entity.Backlog)(entity.Backlog,error)
	Update(ctx context.Context,backlog entity.Backlog)(entity.Backlog,error)
	Delete(ctx context.Context, id string)(error)
}

type UserStoryRepository interface {
	Find(ctx context.Context, id string) (entity.UserStory,error)
	All(ctx context.Context)([]entity.UserStory, error)
	Create(ctx context.Context,story entity.UserStory)(entity.UserStory,error)
	Update(ctx context.Context, story entity.UserStory)(entity.UserStory,error)
	Delete(ctx context.Context, id string)(error)
}

type TaskRepository interface {
	Find(ctx context.Context, id string) (entity.Task,error)
	All(ctx context.Context)([]entity.Task, error)
	Create(ctx context.Context, task entity.Task)(entity.Task,error)
	Update(ctx context.Context, task entity.Task)(entity.Task,error)
	Delete(ctx context.Context, id string)(error)
}

type UserRepository interface {
	Find(ctx context.Context, id string) (entity.User,error)
	All(ctx context.Context)([]entity.User, error)
	Create(ctx context.Context, user entity.User)(entity.User,error)
	Update(ctx context.Context, user entity.User)(entity.User,error)
	Delete(ctx context.Context, id string)(error)
}

type MongoRepository interface {
	Repository
	InitCollection(client *mongo.Client)
}

type MongoProjectRepository interface {
	MongoRepository
	ProjectRepository
}

type MongoBacklogRepository interface {
	MongoRepository
	BacklogRepository
}

type MongoUserStoryRepository interface {
	MongoRepository
	UserStoryRepository
}

type MongoTaskRepository interface {
	MongoRepository
	TaskRepository
}

type MongoUserRepository interface {
	MongoRepository
	UserRepository
}
