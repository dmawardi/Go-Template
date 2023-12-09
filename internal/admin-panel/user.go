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

// Table headers to show on find all page
var userTableHeaders = []TableHeader{
	{Label: "ID", ColumnSortLabel: "id", Pointer: false, DataType: "int"},
	{Label: "Username", ColumnSortLabel: "username", Pointer: false, DataType: "string"},
	{Label: "Email", ColumnSortLabel: "email", Pointer: false, DataType: "string"},
	{Label: "Verified", ColumnSortLabel: "verified", Pointer: true, DataType: "bool"},
}

func NewUserAdminController(service service.UserService) AdminUserController {
	return &adminUserController{
		service: service,
		// Use values from above
		adminHomeUrl: "/admin/users",
		schemaName:   "Users",
		tableHeaders: userTableHeaders,
		formSelectors: FormSelectors{
			"role": roleSelection,
		},
	}
}

type adminUserController struct {
	service service.UserService
	// For links
	adminHomeUrl string
	// For HTML text rendering
	schemaName string
	// Custom table headers
	tableHeaders  []TableHeader
	formSelectors FormSelectors
}

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

func (c adminUserController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab query parameters
	searchQuery := r.URL.Query().Get("search")
	// Grab basic query params
	baseQueryParams, err := helpers.ExtractBasicFindAllQueryParams(r)
	if err != nil {
		http.Error(w, "Error extracting query params", http.StatusBadRequest)
		return
	}

	// Generate query params to extract
	queryParamsToExtract := models.UserConditionQueryParams()
	// Extract query params
	extractedConditionParams, err := helpers.ExtractSearchAndConditionParams(r, queryParamsToExtract)
	if err != nil {
		fmt.Println("Error extracting conditions: ", err)
		http.Error(w, "Can't find conditions", http.StatusBadRequest)
		return
	}

	// Grab all users from database
	users, err := c.service.FindAll(baseQueryParams.Limit, baseQueryParams.Offset, baseQueryParams.Order, extractedConditionParams)
	if err != nil {
		http.Error(w, "Error finding data", http.StatusInternalServerError)
		return
	}
	// Convert data to AdminPanelSchema
	adminUserSlice := c.convertDataToAdminPanelSchema(*users.Data)

	// Build the table data
	tableData := BuildTableData(adminUserSlice, users.Meta, c.adminHomeUrl, c.tableHeaders)

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Admin: Users",
		SectionTitle: "Select a user to edit",
		SidebarList:  sidebarList,
		TableData:    tableData,
		SchemaHome:   c.adminHomeUrl,
		SearchTerm:   searchQuery,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   true,
			CreatePage: false,
			DeletePage: false,
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: c.adminHomeUrl,
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

		// Extract form submission from request and build into map[string]string
		formFieldMap, err := c.extractFormFromRequest(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		// Populate previouisly entered values (Avoids password)
		populateValuesWithForm(r, &createUserForm, formFieldMap)
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

		// Extract form submission from request and build into map[string]string
		fieldMap, err := c.extractFormFromRequest(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		// Populate previously entered values (Avoids password)
		err = populateValuesWithForm(r, &editUserForm, fieldMap)
		if err != nil {
			http.Error(w, "Error populating form", http.StatusInternalServerError)
			return
		}
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
	userData := c.getValuesUsingFieldMap(*foundUser)
	// Populate form field placeholders with data from database
	err = populatePlaceholdersWithDBData(&editUserForm, userData)
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
		SchemaHome:   c.adminHomeUrl,
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

// Form generation
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
		{DbLabel: "Role", Label: "Role", Name: "role", Placeholder: "Enter role", Value: "user", Type: "select", Required: false, Disabled: false, Errors: []ErrorMessage{}, Selectors: c.formSelectors["role"]()},
		{DbLabel: "Verified", Label: "Verified", Name: "verified", Placeholder: "", Value: "false", Type: "checkbox", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "VerificationCode", Label: "Verification Code", Name: "verification_code", Placeholder: "Enter verification code", Value: "", Type: "text", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "VerificationCodeExpiry", Label: "Verification Code Expiry", Name: "verification_code_expiry", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "CreatedAt", Label: "Created At", Name: "created_at", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "UpdatedAt", Label: "Updated At", Name: "updated_at", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
	}
}

// Extract forms
// Used to extract form submission from request and build into models.CreateUser
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

// Used to extract form submission from request and build into models.UpdateUser
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
// Used to extract form submission from request and build into map[string]string (Used in populateValuesWithForm)
func (c adminUserController) extractFormFromRequest(r *http.Request) (map[string]string, error) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing form: %s", err.Error()))
	}
	// Map of user fields
	fieldMap := map[string]string{
		"Name":     r.FormValue("name"),
		"Username": r.FormValue("username"),
		"Email":    r.FormValue("email"),
		"Role":     r.FormValue("role"),
		"Verified": r.FormValue("verified"),
	}
	return fieldMap, nil
}

// For dynamic data iteration: takes a user and returns a map for easier dynamic access
func (c adminUserController) getValuesUsingFieldMap(user db.User) map[string]string {
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
	return fieldMap
}

// Convert schema slice to AdminPanelSchema through appending to new slice of AdminPanelSchema for standardization
func (c adminUserController) convertDataToAdminPanelSchema(slice []db.User) []AdminPanelSchema {
	// Init AdminPanelSchema
	var schemaSlice []AdminPanelSchema
	// Loop through users and append to schemaSlice
	for _, user := range slice {
		// Append user to schemaSlice
		schemaSlice = append(schemaSlice, user)
	}

	return schemaSlice
}

// Form Selectors
// For role selection in form
func roleSelection() []FormFieldSelector {
	return []FormFieldSelector{
		{Value: "user", Label: "User", Selected: true},
		{Value: "admin", Label: "Admin", Selected: false},
		{Value: "moderator", Label: "Moderator", Selected: false},
	}
}
