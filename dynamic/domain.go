package dynamic

import "gorm.io/gorm"

type Participant struct {
	gorm.Model
	Surname      string `json:"surname"`
	Name         string `json:"name"`
	Organizacion string `json:"organizacion"`
	Position     string `json:"position"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Type         string `json:"type"` // Speaker/Publication/Listaner
	Title        string `json:"title"`
}

type Administrator struct {
	gorm.Model
	Secret string `json:"secret"`
}