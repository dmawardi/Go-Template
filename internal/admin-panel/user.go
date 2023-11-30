package adminpanel

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/go-chi/chi"
)

var roleSelection = []FormFieldSelector{
	{Value: "user", Label: "User"},
	{Value: "admin", Label: "Admin"},
	{Value: "moderator", Label: "Moderator"},
}

type AdminUserController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	// Edit is also used to view the record details
	Edit(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type adminUserController struct {
	service service.UserService
}

func NewUserAdminController(service service.UserService) AdminUserController {
	return &adminUserController{service: service}
}

func (c adminUserController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Parse the template
	tmpl, err := parseAdminTemplates()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%+v\n", tmpl.DefinedTemplates())

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Admin: Users",
		SectionTitle: "Select a user to edit",
		SidebarList:  sidebarList,
		TableData: TableData{
			TableHeaders: []string{"ID", "Username", "Email"},
			TableRows: []TableRow{
				{Data: []string{"1", "admin", "admin@bulba.com"}},
				{Data: []string{"2", "admin", "admin@bulba.com"}},
				{Data: []string{"3", "admin", "admin@bulba.com"}},
			},
		},
		PageType: PageType{
			EditPage:   false,
			ReadPage:   true,
			CreatePage: false,
			DeletePage: false,
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: "/admin/users",
				FormMethod: "POST",
			},
			FormFields: []FormField{
				{
					Label:       "Username",
					Name:        "username",
					Placeholder: "Cilandak 213",
					Value:       "",
					Type:        "text",
					Required:    true,
					Disabled:    false,
					Errors:      []ErrorMessage{{Message: "This is an error message"}},
				},
				{
					Label:       "Password",
					Name:        "password",
					Placeholder: "",
					Value:       "",
					Type:        "text",
					Required:    true,
					Disabled:    true,
					Errors:      []ErrorMessage{{Message: "This is an error message"}},
				},
			},
		},
	}

	// Execute the template with data and write to response
	err = tmpl.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c adminUserController) Create(w http.ResponseWriter, r *http.Request) {
	// Parse the template
	tmpl, err := parseAdminTemplates()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	createUserForm := GenerateCreateUserForm()

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Create User",
		SectionTitle: "Create a new user",
		SidebarList:  sidebarList,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: true,
			DeletePage: false,
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: "/admin/users",
				FormMethod: "POST",
			},
			FormFields: createUserForm,
		},
	}

	// Execute the template with data and write to response
	err = tmpl.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c adminUserController) Edit(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Init a new user struct
	foundUser := &db.User{}
	// Search for user by ID and store in foundUser
	foundUser, err = c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Init new User Edit form
	editUserDataSchema := GenerateEditUserForm()

	// Populate form field placeholders with data from database
	err = PopulateUserPlaceholdersWithMap(*foundUser, &editUserDataSchema)
	if err != nil {
		http.Error(w, "Error generating form", http.StatusInternalServerError)
		fmt.Printf("Error generating form: %s\n", err.Error())
		return
	}

	// Parse the template
	tmpl, err := parseAdminTemplates()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Edit User: " + stringParameter,
		SectionTitle: "Edit User: " + stringParameter,
		SidebarList:  sidebarList,
		PageType: PageType{
			EditPage:   true,
			ReadPage:   false,
			CreatePage: false,
			DeletePage: false,
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: "/admin/users",
				FormMethod: "POST",
			},
			FormFields: editUserDataSchema,
		},
	}

	// Execute the template with data and write to response
	err = tmpl.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c adminUserController) Delete(w http.ResponseWriter, r *http.Request) {
	// Parse the template
	tmpl, err := parseAdminTemplates()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%+v\n", tmpl.DefinedTemplates())

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:   "Hello world",
		SidebarList: sidebarList,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: false,
			DeletePage: true,
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: "/admin/users",
				FormMethod: "POST",
			},
			FormFields: []FormField{},
		},
	}

	// Execute the template with data and write to response
	err = tmpl.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// Used to build Create user form
func GenerateCreateUserForm() []FormField {
	return []FormField{
		{DbLabel: "Name", Label: "Name", Name: "name", Placeholder: "Enter name", Value: "", Type: "text", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Username", Label: "Username", Name: "username", Placeholder: "Enter username", Value: "", Type: "text", Required: true, Disabled: false, Errors: []ErrorMessage{{Message: "This is an error message"}}},
		{DbLabel: "Email", Label: "Email", Name: "email", Placeholder: "Enter email", Value: "", Type: "email", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Password", Label: "Password", Name: "password", Placeholder: "Enter password", Value: "", Type: "password", Required: true, Disabled: false, Errors: []ErrorMessage{{Message: "This is an error message"}}},
		{DbLabel: "Role", Label: "Role", Name: "role", Placeholder: "Enter role", Value: "user", Type: "select", Required: false, Disabled: false, Errors: []ErrorMessage{}, Selectors: roleSelection},
		{DbLabel: "Verified", Label: "Verified", Name: "verified", Placeholder: "", Value: "false", Type: "checkbox", Required: false, Disabled: false, Errors: []ErrorMessage{}},
	}
}

// Used to build Edit user form
func GenerateEditUserForm() []FormField {
	return []FormField{
		{DbLabel: "ID", Label: "ID", Name: "id", Placeholder: "", Value: "", Type: "number", Required: false, Disabled: true, Errors: []ErrorMessage{{Message: "This is an error message"}}},
		{DbLabel: "Name", Label: "Name", Name: "name", Placeholder: "Enter name", Value: "", Type: "text", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Username", Label: "Username", Name: "username", Placeholder: "Enter username", Value: "", Type: "text", Required: false, Disabled: false, Errors: []ErrorMessage{{Message: "This is an error message"}}},
		{DbLabel: "Email", Label: "Email", Name: "email", Placeholder: "Enter email", Value: "", Type: "email", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		// {DbLabel: "Password", Label: "Password", Name: "password", Placeholder: "Enter password", Value: "", Type: "password", Required: false, Disabled: false, Errors: []ErrorMessage{{Message: "This is an error message"}}},
		{DbLabel: "Role", Label: "Role", Name: "role", Placeholder: "Enter role", Value: "user", Type: "select", Required: false, Disabled: false, Errors: []ErrorMessage{}, Selectors: roleSelection},
		{DbLabel: "Verified", Label: "Verified", Name: "verified", Placeholder: "", Value: "false", Type: "checkbox", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "VerificationCode", Label: "Verification Code", Name: "verification_code", Placeholder: "Enter verification code", Value: "", Type: "text", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "VerificationCodeExpiry", Label: "Verification Code Expiry", Name: "verification_code_expiry", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "CreatedAt", Label: "Created At", Name: "created_at", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "UpdatedAt", Label: "Updated At", Name: "updated_at", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
	}
}

// Used to populate form field placeholders with data from database
func PopulateUserPlaceholdersWithMap(user db.User, fields *[]FormField) error {
	// Map of user fields
	fieldMap := map[string]string{
		"ID":                     fmt.Sprint(user.ID),
		"CreatedAt":              user.CreatedAt.Format(time.RFC3339),
		"UpdatedAt":              user.UpdatedAt.Format(time.RFC3339),
		"Name":                   user.Name,
		"Username":               user.Username,
		"Email":                  user.Email,
		"Role":                   user.Role,
		"Verified":               fmt.Sprint(user.Verified),
		"VerificationCode":       user.VerificationCode,
		"VerificationCodeExpiry": user.VerificationCodeExpiry.Format(time.RFC3339),
	}

	// Loop through fields and populate placeholders
	for i := range *fields {
		// Get pointer to field
		field := &(*fields)[i]
		// If the field exists in the map, populate the placeholder
		if val, ok := fieldMap[field.DbLabel]; ok {
			field.Placeholder = val
		} else {
			return fmt.Errorf("field: %s not found in map", field.DbLabel)
		}
	}
	return nil
}
