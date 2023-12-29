package adminpanel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

// Table headers to show on find all page
var postTableHeaders = []TableHeader{
	{Label: "ID", ColumnSortLabel: "id", Pointer: false, DataType: "int"},
	{Label: "Title", ColumnSortLabel: "title", Pointer: false, DataType: "string"},
	{Label: "Content", ColumnSortLabel: "content", Pointer: false, DataType: "string"},
	{Label: "User", ColumnSortLabel: "user", Pointer: false, DataType: "foreign"},
}

func NewAdminPostController(service service.PostService, selectorService SelectorService) AdminPostController {
	return &adminPostController{
		service: service,
		// Use values from above
		adminHomeUrl:     "/admin/posts",
		schemaName:       "Post",
		pluralSchemaName: "Posts",
		tableHeaders:     postTableHeaders,
		formSelectors:    selectorService,
	}
}

type adminPostController struct {
	service service.PostService
	// For links
	adminHomeUrl string
	// For HTML text rendering
	schemaName       string
	pluralSchemaName string
	// Custom table headers
	tableHeaders  []TableHeader
	formSelectors SelectorService
}

type AdminPostController interface {
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
	// For sidebar
	ObtainFields() BasicAdminController
}

// CRUD handlers
func (c adminPostController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab query parameters
	searchQuery := r.URL.Query().Get("search")
	// Grab basic query params
	baseQueryParams, err := helpers.ExtractBasicFindAllQueryParams(r)
	if err != nil {
		http.Error(w, "Error extracting query params", http.StatusBadRequest)
		return
	}

	// Generate query params to extract
	queryParamsToExtract := controller.PostConditionQueryParams()
	// Extract query params
	extractedConditionParams, err := helpers.ExtractSearchAndConditionParams(r, queryParamsToExtract)
	if err != nil {
		fmt.Println("Error extracting conditions: ", err)
		http.Error(w, "Can't find conditions", http.StatusBadRequest)
		return
	}

	// Find all with options from database
	found, err := c.service.FindAll(baseQueryParams.Limit, baseQueryParams.Offset, baseQueryParams.Order, extractedConditionParams)
	if err != nil {
		http.Error(w, "Error finding data", http.StatusInternalServerError)
		return
	}

	// Convert data to AdminPanelSchema
	schemaSlice := *found.Data
	var adminSchemaSlice []AdminPanelSchema
	for _, post := range schemaSlice {
		// Append to schemaSlice
		adminSchemaSlice = append(adminSchemaSlice, post)
	}

	// Build the table data
	tableData := BuildTableData(adminSchemaSlice, found.Meta, c.adminHomeUrl, c.tableHeaders)

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Admin: " + c.pluralSchemaName,
		SectionTitle: fmt.Sprintf("Select a %s to edit", c.schemaName),
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
func (c adminPostController) Create(w http.ResponseWriter, r *http.Request) {
	// Init new User Create form
	createForm := c.generateCreateForm()

	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Extract user form submission
		toValidate, err := c.extractCreateFormSubmission(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Validate struct
		pass, valErrors := helpers.GoValidateStruct(toValidate)
		// If failure detected
		// If validation passes
		if pass {
			// Create user
			_, err = c.service.Create(&toValidate)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error creating %s", c.schemaName), http.StatusInternalServerError)
				return
			}
			// Redirect or render a success message
			http.Redirect(w, r, fmt.Sprintf("%s/create/success", c.adminHomeUrl), http.StatusSeeOther)
			return
		}

		// If validation fails
		// Populate form field errors
		SetValidationErrorsInForm(createForm, *valErrors)

		// Extract form submission from request and build into map[string]string
		formFieldMap, err := c.extractFormFromRequest(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		// Populate previously entered values (Avoids password)
		populateValuesWithForm(r, &createForm, formFieldMap)
	}

	// Render preparation
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    fmt.Sprintf("Create %s", c.schemaName),
		SectionTitle: fmt.Sprintf("Create a new %s", c.schemaName),
		SidebarList:  sidebarList,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: true,
			DeletePage: false,
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: fmt.Sprintf("%s/create", c.adminHomeUrl),
				FormMethod: "post",
			},
			FormFields: createForm,
		},
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminPostController) Edit(w http.ResponseWriter, r *http.Request) {
	// Init new User Edit form
	editForm := c.generateEditForm()

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
				http.Error(w, fmt.Sprintf("Error updating %s", c.schemaName), http.StatusInternalServerError)
				return
			}
			// Redirect or render a success message
			http.Redirect(w, r, fmt.Sprintf("%s/edit/success", c.adminHomeUrl), http.StatusSeeOther)
			return
		}

		// If validation fails
		// Populate form field errors
		SetValidationErrorsInForm(editForm, *valErrors)

		// Extract form submission from request and build into map[string]string
		fieldMap, err := c.extractFormFromRequest(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		// Populate previously entered values (Avoids password)
		err = populateValuesWithForm(r, &editForm, fieldMap)
		if err != nil {
			http.Error(w, "Error populating form", http.StatusInternalServerError)
			return
		}
	}

	// If not POST, ie. GET
	// Find current details to use as placeholder values
	// Init a new db struct
	found := &db.Post{}
	// Search for by ID and store in found
	found, err = c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s not found", c.schemaName), http.StatusNotFound)
		return
	}

	// Populate form field placeholders with data from database
	currentData := c.getValuesUsingFieldMap(*found)
	// Populate form field placeholders with data from database
	err = populatePlaceholdersWithDBData(&editForm, currentData)
	if err != nil {
		http.Error(w, "Error generating form", http.StatusInternalServerError)
		return
	}

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    fmt.Sprintf("Edit %s: %s", c.schemaName, stringParameter),
		SectionTitle: fmt.Sprintf("Edit %s: %s", c.schemaName, stringParameter),
		SidebarList:  sidebarList,
		PageType: PageType{
			EditPage:   true,
			ReadPage:   false,
			CreatePage: false,
			DeletePage: false,
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: fmt.Sprintf("%s/%s", c.adminHomeUrl, stringParameter),
				FormMethod: "post",
			},
			FormFields: editForm,
		},
	}

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminPostController) Delete(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, fmt.Sprintf("Error deleting %s", c.schemaName), http.StatusInternalServerError)
			return
		}
		// Redirect to success page
		http.Redirect(w, r, fmt.Sprintf("%s/delete/success", c.adminHomeUrl), http.StatusSeeOther)
		return
	}

	// GET request
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    fmt.Sprintf("Delete %s", c.schemaName),
		SectionTitle: fmt.Sprintf("Are you sure you wish to delete user: %s?", stringParameter),
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
				FormAction: fmt.Sprintf("%s/delete/%s", c.adminHomeUrl, stringParameter),
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
func (c adminPostController) BulkDelete(w http.ResponseWriter, r *http.Request) {
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
func (c adminPostController) CreateSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Create %s", c.schemaName), fmt.Sprintf("%s Created Successfully!", c.schemaName))
}
func (c adminPostController) EditSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Edit %s", c.schemaName), fmt.Sprintf("%s Updated Successfully!", c.schemaName))
}
func (c adminPostController) DeleteSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Delete %s", c.schemaName), fmt.Sprintf("%s Deleted Successfully!", c.schemaName))
}

// Form generation
// Used to build Create user form
func (c adminPostController) generateCreateForm() []FormField {
	return []FormField{
		{DbLabel: "Title", Label: "Title", Name: "title", Placeholder: "", Value: "", Type: "text", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Body", Label: "Body", Name: "body", Placeholder: "", Value: "", Type: "text", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "User", Label: "User", Name: "user", Placeholder: "", Value: "", Type: "select", Required: true, Disabled: false, Errors: []ErrorMessage{}, Selectors: c.formSelectors.UserSelection()},
	}
}

// Used to build Edit user form
func (c adminPostController) generateEditForm() []FormField {
	return []FormField{
		{DbLabel: "Title", Label: "Title", Name: "title", Placeholder: "", Value: "", Type: "text", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Body", Label: "Body", Name: "body", Placeholder: "", Value: "", Type: "text", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "User", Label: "User", Name: "user", Placeholder: "", Value: "", Type: "select", Required: true, Disabled: false, Errors: []ErrorMessage{}, Selectors: c.formSelectors.UserSelection()},
	}
}

// Extract forms
// Used to extract form submission from request and build into models.CreatePost
func (c adminPostController) extractCreateFormSubmission(r *http.Request) (models.CreatePost, error) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		return models.CreatePost{}, errors.New("Error parsing form")
	}

	user := r.FormValue("user")
	// Convert to int
	userId, err := strconv.Atoi(user)
	if err != nil {
		return models.CreatePost{}, errors.New("Error parsing form")
	}

	// Build struct for validation
	toValidate := models.CreatePost{
		Title: r.FormValue("title"),
		Body:  r.FormValue("body"),
		User:  db.User{ID: uint(userId)},
	}

	return toValidate, nil
}

// Used to extract form submission from request and build into models.UpdateUser
func (c adminPostController) extractUpdateFormSubmission(r *http.Request) (models.UpdatePost, error) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		return models.UpdatePost{}, errors.New("Error parsing form")
	}

	// Extract user
	user := r.FormValue("user")
	// Convert to int
	userId, err := strconv.Atoi(user)
	if err != nil {
		return models.UpdatePost{}, errors.New("Error parsing form")
	}
	// Build struct for validation
	toValidate := models.UpdatePost{
		Title: r.FormValue("title"),
		Body:  r.FormValue("body"),
		User:  db.User{ID: uint(userId)},
	}

	return toValidate, nil
}

// Basic helper functions
// Used to extract form submission from request and build into map[string]string (Used in populateValuesWithForm)
func (c adminPostController) extractFormFromRequest(r *http.Request) (map[string]string, error) {
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
func (c adminPostController) getValuesUsingFieldMap(post db.Post) map[string]string {
	// Map of user fields
	fieldMap := map[string]string{
		"ID":        fmt.Sprint(post.ID),
		"CreatedAt": post.CreatedAt.Format(time.RFC3339),
		"UpdatedAt": post.UpdatedAt.Format(time.RFC3339),
		"Title":     post.Title,
		"Body":      post.Body,
		// Foreign key: uses username
		"UserID": fmt.Sprint(post.UserID),
	}
	return fieldMap
}

// Used to build standardize controller fields for admin panel sidebar generation
func (c adminPostController) ObtainFields() BasicAdminController {
	return basicAdminController{
		AdminHomeUrl:     c.adminHomeUrl,
		SchemaName:       c.schemaName,
		PluralSchemaName: c.pluralSchemaName,
	}
}