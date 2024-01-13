package models

import (
	pq "github.com/lib/pq"
)

type UserSession struct {
	UserID      int64 `gorm:"primaryKey"`
	State       string
	AdmMode     string
	Location    string
	Breweries   pq.StringArray `gorm:"type:text"`
	Styles      pq.StringArray `gorm:"type:text"`
	CurrentPage int
}

func (session *UserSession) NewSession(ID int64) {
	session.UserID = ID
	DB.Create(session)
}

func (session *UserSession) LoadInfo() {
	DB.Find(session)
}

func (session *UserSession) SetUserState(State string) {
	session.State = State
	session.Breweries.Value()
	DB.Save(session)
}

func (session *UserSession) SetAdminMode(Mode string) {
	session.AdmMode = Mode
	DB.Save(session)
}

func (session *UserSession) SetLocation(Location string) {
	session.Location = Location
	DB.Save(session)
}

func (session *UserSession) CleanUserFilters() {
	session.Breweries = nil
	session.Styles = nil
	DB.Save(session)
}

func (session *UserSession) AppendBrewery(newBrewery string) {
	if session.Breweries == nil {
		session.Breweries = make([]string, 0, 10)
	}
	session.Breweries = append(session.Breweries, newBrewery)
	DB.Save(session)
}

func (session *UserSession) RemoveBrewery(delBrewery string) {
	for i, b := range session.Breweries {
		if b == delBrewery {
			session.Breweries = append(session.Breweries[:i], session.Breweries[i+1:]...)
		}
	}
	DB.Save(session)
}

func (session *UserSession) AppendStyle(newStyle string) {
	if session.Styles == nil {
		session.Styles = make([]string, 0, 10)
	}
	session.Styles = append(session.Styles, newStyle)
	DB.Save(&session)
}

func (session *UserSession) RemoveStyle(delStyle string) {
	for i, s := range session.Breweries {
		if s == delStyle {
			session.Breweries = append(session.Styles[:i], session.Styles[i+1:]...)
		}
	}
	DB.Save(&session)
}

func (session *UserSession) SetCurrentPage(page int) {
	session.CurrentPage = page
	DB.Save(session)
}
