package models

import (
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "10.5.0.5:6379",
		Password: "123", // no password set
		DB:       0,     // use default DB
	})

	if RedisClient == nil {
		panic("Error occurred while conncting Redis")
	}
}

type UserSession struct {
	UserID    int64
	State     string
	AdmMode string
	Location  string
	Breweries []string
	Styles    []string
}

func (session *UserSession) LoadInfo() {
	DB.Find(session)
}

func (session *UserSession) SetUserState(State string) {
	DB.Find(session)
	session.State = State
	DB.Save(session)
}

func (session *UserSession) SetAdminMode(Mode string) {
	DB.Find(session)
	session.AdmMode = Mode
	DB.Save(session)
}

func (session *UserSession) CleanUserFilters(UserID int64) {
	session.Breweries = nil
	session.Styles = nil
	DB.Save(session)
}
