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
	adminpanel "github.com/dmawardi/Go-Template/internal/helpers/adminPanel"
	"github.com/dmawardi/Go-Template/internal/helpers/request"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/go-chi/chi"
)

// Table headers to show on find all page
var postTableHeaders = []TableHeader{
	{Label: "ID", ColumnSortLabel: "id", Pointer: false, DataType: "int", Sortable: true},
	{Label: "Title", ColumnSortLabel: "title", Pointer: false, DataType: "string", Sortable: true},
	{Label: "Body", ColumnSortLabel: "body", Pointer: false, DataType: "string", Sortable: true},
	{Label: "User", ColumnSortLabel: "user", Pointer: false, DataType: "foreign", ForeignKeyRepKeyName: "Username"},
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
	baseQueryParams, err := request.ExtractBasicFindAllQueryParams(r)
	if err != nil {
		http.Error(w, "Error extracting query params", http.StatusBadRequest)
		return
	}

	// Generate query params to extract
	queryParamsToExtract := controller.PostConditionQueryParams()
	// Extract query params
	extractedConditionParams, err := request.ExtractSearchAndConditionParams(r, queryParamsToExtract)
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
		PageTitle:              "Admin: " + c.pluralSchemaName,
		SectionTitle:           fmt.Sprintf("Select a %s to edit", c.schemaName),
		SidebarList:            sidebar,
		TableData:              tableData,
		SchemaHome:             c.adminHomeUrl,
		SearchTerm:             searchQuery,
		RecordsPerPageSelector: recordsPerPage,
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
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminPostController) Create(w http.ResponseWriter, r *http.Request) {
	// Init new Create form
	createForm := c.generateCreateForm()

	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Extract form submission
		formFieldMap, err := adminpanel.ParseFormToMap(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		// Convert relationsip to int
		userId, err := strconv.Atoi(formFieldMap["user"])
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		toValidate := models.CreatePost{
			Title: formFieldMap["title"],
			Body:  formFieldMap["body"],
			User:  db.User{ID: uint(userId)},
		}

		// Validate struct
		pass, valErrors := request.GoValidateStruct(toValidate)
		// If failure detected
		// If validation passes
		if pass {
			// Create
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

		// Populate previously entered values (Avoids password inputs)
		err = populateFormValuesWithSubmittedFormMap(&createForm, formFieldMap)
		if err != nil {
			http.Error(w, "Error populating form", http.StatusInternalServerError)
			return
		}
	}

	// Render preparation
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    fmt.Sprintf("Create %s", c.schemaName),
		SectionTitle: fmt.Sprintf("Create a new %s", c.schemaName),
		SidebarList:  sidebar,
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
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminPostController) Edit(w http.ResponseWriter, r *http.Request) {
	// Init new form
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
		// Extract form submission
		formFieldMap, err := adminpanel.ParseFormToMap(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		// Convert relationsip to int
		userId, err := strconv.Atoi(formFieldMap["user"])
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		toValidate := models.UpdatePost{
			Title: formFieldMap["title"],
			Body:  formFieldMap["body"],
			User:  db.User{ID: uint(userId)},
		}

		// Validate struct
		pass, valErrors := request.GoValidateStruct(toValidate)
		// If failure detected
		// If validation passes
		if pass {
			// Update
			_, err = c.service.Update(idParameter, &toValidate)
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

		// Populate previously entered values (Avoids password)
		err = populateFormValuesWithSubmittedFormMap(&editForm, formFieldMap)
		if err != nil {
			fmt.Printf("Error populating form: %v\n", err)
			http.Error(w, "Error populating form", http.StatusInternalServerError)
			return
		}
	}

	// If not POST, ie. GET
	// Find current details to use as placeholder values
	// Search for by ID and store in found
	found, err := c.service.FindById(idParameter)
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
		SidebarList:  sidebar,
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
		HeaderSection: header,
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
		SectionTitle: fmt.Sprintf("Are you sure you wish to delete %s: %s?", c.schemaName, stringParameter),
		SidebarList:  sidebar,
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
		HeaderSection: header,
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

	// Decode request body as JSON and store
	err := json.NewDecoder(r.Body).Decode(&listOfIds)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Convert string slice to int slice
	intIdList, err := helpers.ConvertStringSliceToIntSlice(listOfIds.SelectedItems)
	if err != nil {
		bulkResponse.Errors = append(bulkResponse.Errors, err)
		bulkResponse.Success = false
		request.WriteAsJSON(w, bulkResponse)
		return
	}

	// Bulk Delete
	err = c.service.BulkDelete(intIdList)
	// If error detected send error response
	if err != nil {
		bulkResponse.Errors = append(bulkResponse.Errors, err)
		bulkResponse.Success = false
		request.WriteAsJSON(w, bulkResponse)
		return
	}
	// else if successful
	bulkResponse.Success = true
	request.WriteAsJSON(w, bulkResponse)
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
// Used to build Create form
func (c adminPostController) generateCreateForm() []FormField {
	return []FormField{
		{DbLabel: "Title", Label: "Title", Name: "title", Placeholder: "", Value: "", Type: "text", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Body", Label: "Body", Name: "body", Placeholder: "", Value: "", Type: "textarea", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "User", Label: "User", Name: "user", Placeholder: "", Value: "", Type: "select", Required: true, Disabled: false, Errors: []ErrorMessage{}, Selectors: c.formSelectors.UserSelection()},
	}
}

// Used to build Edit form
func (c adminPostController) generateEditForm() []FormField {
	return []FormField{
		{DbLabel: "Title", Label: "Title", Name: "title", Placeholder: "", Value: "", Type: "text", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Body", Label: "Body", Name: "body", Placeholder: "", Value: "", Type: "textarea", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "User", Label: "User", Name: "user", Placeholder: "", Value: "", Type: "select", Required: true, Disabled: false, Errors: []ErrorMessage{}, Selectors: c.formSelectors.UserSelection()},
	}
}

// Extract forms

// Basic helper functions
// For dynamic data iteration: takes a db struct and returns a map for easier dynamic access
// Used prior to populating form placeholders
func (c adminPostController) getValuesUsingFieldMap(post db.Post) map[string]string {
	// Map of user fields
	fieldMap := map[string]string{
		"ID":        fmt.Sprint(post.ID),
		"CreatedAt": post.CreatedAt.Format(time.RFC3339),
		"UpdatedAt": post.UpdatedAt.Format(time.RFC3339),
		"Title":     post.Title,
		"Body":      post.Body,
		// Foreign key: Return ID as string
		"User": fmt.Sprint(post.UserID),
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
