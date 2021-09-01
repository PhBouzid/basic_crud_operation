//Packages entity makes data conversation between controllers and domain
package entity

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Backlog struct {
	Title            string          `json:"title"`
	BriefDescription string          `json:"brief"`
	Content          string          `json:"content"`
	TimeCreated      time.Time       `json:"time_created"`
	User_Id          string          `json:"user_id"`
	Id               string          `json:"id"`
	Retro            []Retrospective `json:"retrospectives,omitempty"`
}

type Retrospective struct {
	TaskId      string    `json:"task_id"`
	Information string    `json:"information"`
	AuthorId    string    `json:"author"`
	TimeCreated time.Time `json:"time_created"`
}

type Project struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Weeks        int       `json:"weeks"`
	StartDate    time.Time `json:"start_date"`
	WorkDuration int64     `json:"work_duration"`
	AuthorId     string    `json:"authorId" faker:"-"`
}

type UserStory struct {
	Id        string      `json:"id" faker:"-"`
	Title     string      `json:"title"`
	Brief     string      `json:"brief"`
	Content   string      `json:"content"`
	Tasks     []TaskShort `json:"tasks,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type TaskShort struct {
	TaskId     string `json:"task_id" faker:"-"`
	Brief      string `json:"brief"`
	State      int    `json:"state"`
	ExecutorId string `json:"executor_id" faker:"-"`
}

type Task struct {
	Id          string         `json:"id" faker:"-"`
	Title       string         `json:"title"`
	Brief       string         `json:"brief"`
	Content     string         `json:"content"`
	State       int            `json:"state"`
	ExecutorId  string         `json:"executor_id" faker:"-"`
	AuthorId    string         `json:"author_id" faker:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Prioritet   int            `json:"prioritet"`
	Action      []ActionStory  `json:"action"`
	Achievement []Achievements `json:"achievements"`
}

type ActionStory struct {
	Action      string    `json:"action"`
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time""`
	AuthorId    string    `json:"author_id" faker:"-"`
}

type Achievements struct {
	Achievment  string    `json:"achievment"`
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
	AuthorId    string    `json:"author_id" faker:"-"`
}

type User struct {
	UserName   string `json:"user_name"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Id         string `json:"id" faker:"-"`
	Email      string `json:"email"`
	Password   string `json:"-"`
}

type Token struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
	Email  string `json:"email"`
}
