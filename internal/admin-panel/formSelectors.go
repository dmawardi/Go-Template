package adminpanel

import (
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/service"
	"gorm.io/gorm"
)

type SelectorService interface {
	RoleSelection() []FormFieldSelector
	UserSelection() []FormFieldSelector
	ActionSelection() []FormFieldSelector
}
type selectorService struct {
	DB   *gorm.DB
	Auth service.AuthPolicyService
}

func NewSelectorService(db *gorm.DB, auth service.AuthPolicyService) SelectorService {
	return &selectorService{DB: db, Auth: auth}
}

// Form Selectors
// For role selection in form
func (c selectorService) RoleSelection() []FormFieldSelector {
	roles, err := c.Auth.FindAllRoles()
	if err != nil {
		// Return default selector
		return []FormFieldSelector{
			{Value: "user", Label: "User", Selected: true},
			{Value: "admin", Label: "Admin", Selected: false},
			{Value: "moderator", Label: "Moderator", Selected: false},
		}
	}
	// Init form field selector
	var roleSelector []FormFieldSelector
	// Build []FormFieldSelector from []string
	for _, r := range roles {
		roleSelector = append(roleSelector, FormFieldSelector{Value: r, Label: helpers.CapitalizeFirstLetter(r)})
	}

	// Set basic default as user
	c.setDefaultSelected(roleSelector, "user")

	return roleSelector
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

// Takes a slice of FormFieldSelector and sets the Selected field to true for the selector with the specified default value
func (c selectorService) setDefaultSelected(selector []FormFieldSelector, defaultValue string) {
	for i, s := range selector {
		if s.Value == defaultValue {
			selector[i].Selected = true
		} else {
			selector[i].Selected = false
		}
	}
}
