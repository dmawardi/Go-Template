package adminpanel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

// For role selection in form
func roleSelection() []FormFieldSelector {
	return []FormFieldSelector{
		{Value: "user", Label: "User", Selected: true},
		{Value: "admin", Label: "Admin", Selected: false},
		{Value: "moderator", Label: "Moderator", Selected: false},
	}
}

var tableHeaders = []TableHeader{
	{Label: "ID", ColumnSortLabel: "id"},
	{Label: "Name", ColumnSortLabel: "name"},
	{Label: "Username", ColumnSortLabel: "username"},
}

// Schema home used to return to the schema home page from delete
var adminHomeUrl = "/admin/users"

type AdminUserController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	// Edit is also used to view the record details
	Edit(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	// Bulk delete (from table)
	BulkDelete(w http.ResponseWriter, r *http.Request)
	// Success pages
	CreateSuccess(w http.ResponseWriter, r *http.Request)
	EditSuccess(w http.ResponseWriter, r *http.Request)
	DeleteSuccess(w http.ResponseWriter, r *http.Request)
}

type adminUserController struct {
	service service.UserService
}

func NewUserAdminController(service service.UserService) AdminUserController {
	return &adminUserController{service: service}
}

func (c adminUserController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab query parameters
	searchQuery := r.URL.Query().Get("search")
	page, err := helpers.GrabIntQueryParamOrDefault(r, "page", 1)
	if err != nil {
		http.Error(w, "Error parsing page", http.StatusBadRequest)
		return
	}
	limit, err := helpers.GrabIntQueryParamOrDefault(r, "limit", 10)
	if err != nil {
		http.Error(w, "Error parsing limit", http.StatusBadRequest)
		return
	}
	// Grab order
	order := helpers.GrabQueryParamOrDefault(r, "order", "id")
	// Calculate offset using pages and limit
	offset := (page - 1) * limit

	// Generate query params to extract
	queryParamsToExtract := models.UserQueryParams()
	// Extract query params
	extractedQueryParams, err := helpers.ExtractConditionParams(r, queryParamsToExtract)
	if err != nil {
		fmt.Println("Error extracting conditions: ", err)
		http.Error(w, "Can't find conditions", http.StatusBadRequest)
		return
	}

	// Grab all users from database
	users, err := c.service.FindAll(limit, offset, order, extractedQueryParams)
	if err != nil {
		http.Error(w, "Error finding data", http.StatusInternalServerError)
		return
	}
	// Convert data to AdminPanelSchema
	adminUserSlice := c.convertDataToAdminPanelSchema(*users.Data)

	// Build the table data
	tableData := BuildTableData(adminUserSlice, users.Meta, adminHomeUrl, tableHeaders)

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Admin: Users",
		SectionTitle: "Select a user to edit",
		SidebarList:  sidebarList,
		TableData:    tableData,
		SchemaHome:   adminHomeUrl,
		SearchTerm:   searchQuery,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   true,
			CreatePage: false,
			DeletePage: false,
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: adminHomeUrl,
				FormMethod: "get",
			},
			FormFields: []FormField{},
		},
	}

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c adminUserController) Create(w http.ResponseWriter, r *http.Request) {
	// Init new User Create form
	createUserForm := c.generateCreateForm()

	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Extract user form submission
		userToValidate, err := c.extractCreateFormSubmission(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Validate struct
		pass, valErrors := helpers.GoValidateStruct(userToValidate)
		// If failure detected
		// If validation passes
		if pass {
			// Create user
			_, err = c.service.Create(&userToValidate)
			if err != nil {
				http.Error(w, "Error creating user", http.StatusInternalServerError)
				return
			}
			// Redirect or render a success message
			http.Redirect(w, r, "/admin/users/create/success", http.StatusSeeOther)
			return
		}

		// If validation fails
		// Populate form field errors
		SetValidationErrorsInForm(createUserForm, *valErrors)
		// Populate previouisly entered values (Avoids password)
		c.populateValuesWithForm(r, &createUserForm)
	}

	// Render preparation
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
				FormAction: "/admin/users/create",
				FormMethod: "post",
			},
			FormFields: createUserForm,
		},
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c adminUserController) Edit(w http.ResponseWriter, r *http.Request) {
	// Init new User Edit form
	editUserForm := c.generateEditForm()

	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Extract user form submission
		userToValidate, err := c.extractUpdateFormSubmission(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		fmt.Printf("Edit user form extracted: %+v\n", userToValidate)

		// Validate struct
		pass, valErrors := helpers.GoValidateStruct(userToValidate)
		// If failure detected
		// If validation passes
		if pass {
			// Update user
			_, err = c.service.Update(idParameter, &userToValidate)
			if err != nil {
				http.Error(w, "Error updating user", http.StatusInternalServerError)
				return
			}
			// Redirect or render a success message
			http.Redirect(w, r, "/admin/users/edit/success", http.StatusSeeOther)
			return
		}

		// If validation fails
		// Populate form field errors
		SetValidationErrorsInForm(editUserForm, *valErrors)
		// Populate previouisly entered values (Avoids password)
		c.populateValuesWithForm(r, &editUserForm)
	}

	// If not POST, ie. GET
	// Find current details of user to use as placeholder values
	// Init a new user struct
	foundUser := &db.User{}
	// Search for user by ID and store in foundUser
	foundUser, err = c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Populate form field placeholders with data from database
	err = c.populatePlaceholders(*foundUser, &editUserForm)
	if err != nil {
		http.Error(w, "Error generating form", http.StatusInternalServerError)
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
				FormAction: "/admin/users/" + stringParameter,
				FormMethod: "post",
			},
			FormFields: editUserForm,
		},
	}

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c adminUserController) Delete(w http.ResponseWriter, r *http.Request) {
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	if err != nil {
		serveAdminError(w, "Unable to interpret ID")
		return
	}
	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Delete user
		err = c.service.Delete(idParameter)
		if err != nil {
			http.Error(w, "Error deleting user", http.StatusInternalServerError)
			return
		}
		// Redirect to success page
		http.Redirect(w, r, "/admin/users/delete/success", http.StatusSeeOther)
		return
	}

	// GET request
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Delete User",
		SectionTitle: "Are you sure you wish to delete user: " + stringParameter + "?",
		SidebarList:  sidebarList,
		SchemaHome:   adminHomeUrl,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: false,
			DeletePage: true,
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: "/admin/users/delete/" + stringParameter,
				FormMethod: "post",
			},
			FormFields: []FormField{},
		},
	}

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c adminUserController) BulkDelete(w http.ResponseWriter, r *http.Request) {
	// Grab body of request
	// Init
	var listOfIds BulkDeleteRequest

	// Prepare response
	bulkResponse := models.BulkDeleteResponse{
		// Set deleted records to length of selected items
		DeletedRecords: len(listOfIds.SelectedItems),
		Errors:         []error{},
	}

	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&listOfIds)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Convert string slice to int slice
	intIdList, err := convertStringSliceToIntSlice(listOfIds.SelectedItems)
	if err != nil {
		bulkResponse.Errors = append(bulkResponse.Errors, err)
		bulkResponse.Success = false
		helpers.WriteAsJSON(w, bulkResponse)
		return
	}

	// Bulk Delete users
	err = c.service.BulkDelete(intIdList)
	// If error detected send error response
	if err != nil {
		bulkResponse.Errors = append(bulkResponse.Errors, err)
		bulkResponse.Success = false
		helpers.WriteAsJSON(w, bulkResponse)
		return
	}
	// else if successful
	bulkResponse.Success = true
	helpers.WriteAsJSON(w, bulkResponse)
}

// Success handlers
func (c adminUserController) CreateSuccess(w http.ResponseWriter, r *http.Request) {
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "User Creation form submitted",
		SectionTitle: "User Created Successfully!",
		SidebarList:  sidebarList,
		PageType: PageType{
			EditPage:    false,
			ReadPage:    false,
			CreatePage:  false,
			DeletePage:  false,
			SuccessPage: true,
		},
		FormData: FormData{},
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminUserController) EditSuccess(w http.ResponseWriter, r *http.Request) {
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "User Edit form submitted",
		SectionTitle: "User Updated Successfully!",
		SidebarList:  sidebarList,
		PageType: PageType{
			EditPage:    false,
			ReadPage:    false,
			CreatePage:  false,
			DeletePage:  false,
			SuccessPage: true,
		},
		FormData: FormData{},
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminUserController) DeleteSuccess(w http.ResponseWriter, r *http.Request) {
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "User Delete form submitted",
		SectionTitle: "User Deleted Successfully!",
		SidebarList:  sidebarList,
		PageType: PageType{
			EditPage:    false,
			ReadPage:    false,
			CreatePage:  false,
			DeletePage:  false,
			SuccessPage: true,
		},
		FormData: FormData{},
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// Form Helper functions
// Used to build Create user form
func (c adminUserController) generateCreateForm() []FormField {
	return []FormField{
		{DbLabel: "Name", Label: "Name", Name: "name", Placeholder: "Enter name", Value: "", Type: "text", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Username", Label: "Username", Name: "username", Placeholder: "Enter username", Value: "", Type: "text", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Email", Label: "Email", Name: "email", Placeholder: "Enter email", Value: "", Type: "email", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Password", Label: "Password", Name: "password", Placeholder: "Enter password", Value: "", Type: "password", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Role", Label: "Role", Name: "role", Placeholder: "Enter role", Value: "user", Type: "select", Required: false, Disabled: false, Errors: []ErrorMessage{}, Selectors: roleSelection()},
		{DbLabel: "Verified", Label: "Verified", Name: "verified", Placeholder: "", Value: "true", Type: "checkbox", Required: false, Disabled: false, Errors: []ErrorMessage{}},
	}
}

// Used to build Edit user form
func (c adminUserController) generateEditForm() []FormField {
	return []FormField{
		{DbLabel: "ID", Label: "ID", Name: "id", Placeholder: "", Value: "", Type: "number", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "Name", Label: "Name", Name: "name", Placeholder: "Enter name", Value: "", Type: "text", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Username", Label: "Username", Name: "username", Placeholder: "Enter username", Value: "", Type: "text", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Email", Label: "Email", Name: "email", Placeholder: "Enter email", Value: "", Type: "email", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		// {DbLabel: "Password", Label: "Password", Name: "password", Placeholder: "Enter password", Value: "", Type: "password", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Role", Label: "Role", Name: "role", Placeholder: "Enter role", Value: "user", Type: "select", Required: false, Disabled: false, Errors: []ErrorMessage{}, Selectors: roleSelection()},
		{DbLabel: "Verified", Label: "Verified", Name: "verified", Placeholder: "", Value: "false", Type: "checkbox", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "VerificationCode", Label: "Verification Code", Name: "verification_code", Placeholder: "Enter verification code", Value: "", Type: "text", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "VerificationCodeExpiry", Label: "Verification Code Expiry", Name: "verification_code_expiry", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "CreatedAt", Label: "Created At", Name: "created_at", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "UpdatedAt", Label: "Updated At", Name: "updated_at", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
	}
}

// Used to extract form submission from request
func (c adminUserController) extractCreateFormSubmission(r *http.Request) (models.CreateUser, error) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		return models.CreateUser{}, errors.New("Error parsing form")
	}

	// Preparation for validation
	// Parse the verified attribute from the form
	verified := false
	// If the verified attribute is present in the form (ie. true)
	if r.FormValue("verified") != "" {
		// Set verified to true
		verified = true
	}
	// Build struct for validation
	userToValidate := models.CreateUser{
		Name:     r.FormValue("name"),
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Role:     r.FormValue("role"),
		Verified: verified,
	}

	return userToValidate, nil
}

// Used to extract form submission from request
func (c adminUserController) extractUpdateFormSubmission(r *http.Request) (models.UpdateUser, error) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		return models.UpdateUser{}, errors.New("Error parsing form")
	}

	// Preparation for validation
	// Parse the verified attribute from the form
	verified := false
	// If the verified attribute is present in the form (ie. true)
	if r.FormValue("verified") != "" {
		// Set verified to true
		verified = true
	}
	// Build struct for validation
	userToValidate := models.UpdateUser{
		Name:     r.FormValue("name"),
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Role:     r.FormValue("role"),
		Verified: verified,
	}

	return userToValidate, nil
}

// Basic helper functions
// Used to populate form field placeholders with data from database
func (c adminUserController) populatePlaceholders(user db.User, form *[]FormField) error {
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
	for i := range *form {
		// Get pointer to field
		field := &(*form)[i]
		if field.Type == "select" {
			// Update selectors with default value
			field.Selectors = addDefaultSelectedToSelector(field.Selectors, fieldMap[field.DbLabel])
			// Else treat as ordinary input
		} else {
			// If the field exists in the map, populate the placeholder
			if val, ok := fieldMap[field.DbLabel]; ok {
				field.Placeholder = val
			} else {
				return fmt.Errorf("field: %s not found in map", field.DbLabel)
			}
		}
	}

	return nil
}

// Used to populate form field placeholders with data from database
func (c adminUserController) populateValuesWithForm(r *http.Request, form *[]FormField) error {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		return errors.New("Error parsing form")
	}

	// Map of user fields
	fieldMap := map[string]string{
		"Name":     r.FormValue("name"),
		"Username": r.FormValue("username"),
		"Email":    r.FormValue("email"),
		"Role":     r.FormValue("role"),
		"Verified": r.FormValue("verified"),
	}

	// Loop through fields and populate placeholders
	for i := range *form {
		// Get pointer to field
		field := &(*form)[i]
		// If the field exists in the map, populate the placeholder
		if val, ok := fieldMap[field.DbLabel]; ok {
			field.Value = val
		} else {
			return fmt.Errorf("field: %s not found in map", field.DbLabel)
		}
	}
	return nil
}

// Convert data to AdminPanelSchema
func (c adminUserController) convertDataToAdminPanelSchema(slice []db.User) []db.AdminPanelSchema {
	// Init AdminPanelSchema
	var schemaSlice []db.AdminPanelSchema
	// Loop through users and append to schemaSlice
	for _, user := range slice {
		// Append user to schemaSlice
		schemaSlice = append(schemaSlice, user)
	}

	return schemaSlice
}
