package main

import (
	"BlockDoc/entity"
)

type ProjectStatusInterface interface {
	OpenProject(project entity.Project) int
	InOrderProject(project entity.Project) int
	ClosedProject(project entity.Project) int
}

type BacklogStatusInterface interface {
	OpenBacklog (backlog entity.Backlog) int
	CloseBacklog(backlog entity.Backlog) int
	InOrderBacklog(backlog entity.Backlog) int
}

type StatusInterface interface {
	ToDo (task entity.Task) int
	DoIt (task entity.Task) int
	Done (task entity.Task) int
}

type Status struct {
	State int `json:"state" bson:"state"`
}

type ProjectStatus struct{
	State int `json:"state" bson:"state"`
}
