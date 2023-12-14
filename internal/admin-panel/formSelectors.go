package adminpanel

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
	"gorm.io/gorm"
)

type SelectorService interface {
	RoleSelection() []FormFieldSelector
	UserSelection() []FormFieldSelector
}
type selectorService struct {
	DB *gorm.DB
}

func NewSelectorService(db *gorm.DB) SelectorService {
	return &selectorService{DB: db}
}

// Form Selectors
// For role selection in form
func (c selectorService) RoleSelection() []FormFieldSelector {
	return []FormFieldSelector{
		{Value: "user", Label: "User", Selected: true},
		{Value: "admin", Label: "Admin", Selected: false},
		{Value: "moderator", Label: "Moderator", Selected: false},
	}
}

func (c selectorService) UserSelection() []FormFieldSelector {
	var users []db.User
	err := c.DB.Select("id, username").Find(&users)
	if err != nil {
		return nil
	}
	fmt.Printf("users in selection: %+v", users)

	return []FormFieldSelector{
		{Value: "user", Label: "User", Selected: true},
		{Value: "admin", Label: "Admin", Selected: false},
		{Value: "moderator", Label: "Moderator", Selected: false},
	}
}
