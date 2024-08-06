package adminpanel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	adminpanel "github.com/dmawardi/Go-Template/internal/helpers/adminPanel"
	"github.com/dmawardi/Go-Template/internal/helpers/data"
	"github.com/dmawardi/Go-Template/internal/helpers/request"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/go-chi/chi"
)

// Interface for all basic admin controllers (used for Admin panel to dynamically generate sidebar)
type BasicAdminController interface {
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
	// Obtain URL details for sidebar
	ObtainUrlDetails() URLDetails
}
type basicAdminController struct {
	service service.BasicModuleService
	// For links
	AdminHomeUrl string
	// For HTML text rendering
	SchemaName       string
	PluralSchemaName string
	// Custom table headers
	tableHeaders  []TableHeader
	formSelectors SelectorService
	// Conditional query params
	ConditionQueryParams map[string]string

	// Input functions for forms
	// Form creators
	generateCreateForm func() []FormField 
	generateEditForm   func() []FormField 
	// Submission preparation
	prepareSubmittedFormForCreation func(formFieldMap map[string]string) (struct{}, error)
}

// Helper functions for sidebar url details
func (c basicAdminController) ObtainUrlDetails() URLDetails {
	return URLDetails{
		AdminHomeUrl:     c.AdminHomeUrl,
		SchemaName:       c.SchemaName,
		PluralSchemaName: c.PluralSchemaName,
	}
}

// ADMIN SIDEBAR CREATION
//
func (c basicAdminController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab query parameters
	searchQuery := r.URL.Query().Get("search")
	// Grab basic query params
	baseQueryParams, err := request.ExtractBasicFindAllQueryParams(r)
	if err != nil {
		http.Error(w, "Error extracting query params", http.StatusBadRequest)
		return
	}

	// Generate query params to extract
	queryParamsToExtract := c.ConditionQueryParams
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
	var adminSchemaSlice []models.AdminPanelSchema
	for _, schema := range schemaSlice {
		// Append to schemaSlice
		adminSchemaSlice = append(adminSchemaSlice, schema)
	}

	// Build the table data
	tableData := BuildTableData(adminSchemaSlice, found.Meta, c.AdminHomeUrl, c.tableHeaders)

	// Generate Find All render data using input data
	data := GenerateFindAllRenderData(tableData, c.SchemaName, c.PluralSchemaName, c.AdminHomeUrl, searchQuery)

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c basicAdminController) Create(w http.ResponseWriter, r *http.Request) {
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

		// Prepare submitted form for creation
		toValidate, err := c.prepareSubmittedFormForCreation(formFieldMap)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Validate struct
		pass, valErrors := request.GoValidateStruct(toValidate)
		// If failure detected
		// If validation passes
		if pass {
			// Create
			_, err = c.service.Create(&toValidate)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error creating %s", c.SchemaName), http.StatusInternalServerError)
				return
			}
			// Redirect or render a success message
			http.Redirect(w, r, fmt.Sprintf("%s/create/success", c.AdminHomeUrl), http.StatusSeeOther)
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

	// Render page data
	data := GenerateCreateRenderData(createForm, c.SchemaName, c.PluralSchemaName, c.AdminHomeUrl)

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c basicAdminController) Edit(w http.ResponseWriter, r *http.Request) {
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
		// Prepare submitted form for creation
		toValidate, err := c.prepareSubmittedFormForCreation(formFieldMap)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Validate struct
		pass, valErrors := request.GoValidateStruct(toValidate)
		// If failure detected
		// If validation passes
		if pass {
			// Update
			_, err = c.service.Update(idParameter, &toValidate)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error updating %s", c.SchemaName), http.StatusInternalServerError)
				return
			}
			// Redirect or render a success message
			http.Redirect(w, r, fmt.Sprintf("%s/edit/success", c.AdminHomeUrl), http.StatusSeeOther)
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
		http.Error(w, fmt.Sprintf("%s not found", c.SchemaName), http.StatusNotFound)
		return
	}

	// Populate form field placeholders with data from database
	currentData := getValuesUsingFieldMap(*found)
	// Populate form field placeholders with data from database
	err = populatePlaceholdersWithDBData(&editForm, currentData)
	if err != nil {
		http.Error(w, "Error generating form", http.StatusInternalServerError)
		return
	}

	data := GenerateEditRenderData(editForm, c.SchemaName, c.PluralSchemaName, c.AdminHomeUrl, stringParameter)	

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c basicAdminController) Delete(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, fmt.Sprintf("Error deleting %s", c.SchemaName), http.StatusInternalServerError)
			return
		}
		// Redirect to success page
		http.Redirect(w, r, fmt.Sprintf("%s/delete/success", c.AdminHomeUrl), http.StatusSeeOther)
		return
	}

	data := GenerateDeleteRenderData(c.SchemaName, c.AdminHomeUrl, c.PluralSchemaName, stringParameter)

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c basicAdminController) BulkDelete(w http.ResponseWriter, r *http.Request) {
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
	intIdList, err := data.ConvertStringSliceToIntSlice(listOfIds.SelectedItems)
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
func (c basicAdminController) CreateSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Create %s", c.SchemaName), fmt.Sprintf("%s Created Successfully!", c.SchemaName))
}
func (c basicAdminController) EditSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Edit %s", c.SchemaName), fmt.Sprintf("%s Updated Successfully!", c.SchemaName))
}
func (c basicAdminController) DeleteSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Delete %s", c.SchemaName), fmt.Sprintf("%s Deleted Successfully!", c.SchemaName))
}