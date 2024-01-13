package models

type Location struct {
	ID string
	WelcomeText string
	PhoneNumbers string
	Email string
	ImagePath string
}

func (l *Location) LoadInfo() {
	DB.Find(l)
}