package adminpanel

import (
	"github.com/dmawardi/Go-Template/internal/db"
	"gorm.io/gorm"
)

type SelectorService interface {
	RoleSelection() []FormFieldSelector
	UserSelection() []FormFieldSelector
	ActionSelection() []FormFieldSelector
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

	return []FormFieldSelector{
		{Value: "user", Label: "User", Selected: true},
		{Value: "admin", Label: "Admin", Selected: false},
		{Value: "moderator", Label: "Moderator", Selected: false},
	}
}

func (c selectorService) ActionSelection() []FormFieldSelector {
	return []FormFieldSelector{
		{Value: "create", Label: "Create", Selected: true},
		{Value: "read", Label: "Read", Selected: false},
		{Value: "update", Label: "Update", Selected: false},
		{Value: "delete", Label: "Delete", Selected: false},
	}
}
