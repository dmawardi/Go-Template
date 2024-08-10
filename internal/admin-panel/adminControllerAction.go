package adminpanel

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dmawardi/Go-Template/internal/controller/core"
	"github.com/dmawardi/Go-Template/internal/db"
	webapi "github.com/dmawardi/Go-Template/internal/helpers/webApi"

	"github.com/dmawardi/Go-Template/internal/helpers/request"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/go-chi/chi"
)

// Table headers to show on find all page
var actionTableHeaders = []TableHeader{
	{Label: "ID", ColumnSortLabel: "id", Pointer: false, DataType: "int", Sortable: true},
	{Label: "Admin", ColumnSortLabel: "admin_id", Pointer: false, DataType: "int", Sortable: false},
	{Label: "Description", ColumnSortLabel: "description", Pointer: false, DataType: "string", Sortable: true},
	{Label: "EntityType", ColumnSortLabel: "entity_type", Pointer: false, DataType: "string"},
	{Label: "EntityID", ColumnSortLabel: "entity_id", Pointer: true, DataType: "string", Sortable: true},
}

func NewAdminActionController(service webapi.ActionService) AdminActionController {
	return &adminActionController{
		service: service,
		// Use values from above
		adminHomeUrl:     "/admin/actions",
		schemaName:       "Action",
		pluralSchemaName: "Actions",
		tableHeaders:     actionTableHeaders,
	}
}

type adminActionController struct {
	service webapi.ActionService
	// For links
	adminHomeUrl string
	// For HTML text rendering
	schemaName       string
	pluralSchemaName string
	// Custom table headers
	tableHeaders  []TableHeader
}

type AdminActionController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	View(w http.ResponseWriter, r *http.Request)
}

func (c adminActionController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab query parameters
	searchQuery := r.URL.Query().Get("search")
	// Grab basic query params
	baseQueryParams, err := request.ExtractBasicFindAllQueryParams(r)
	if err != nil {
		http.Error(w, "Error extracting query params", http.StatusBadRequest)
		return
	}

	// Generate query params to extract
	queryParamsToExtract := core.ActionConditionQueryParams()
	// Extract query params
	extractedConditionParams, err := request.ExtractSearchAndConditionParams(r, queryParamsToExtract)
	if err != nil {
		fmt.Println("Error extracting conditions: ", err)
		http.Error(w, "Can't find conditions", http.StatusBadRequest)
		return
	}

	// Grab all items
	found, err := c.service.FindAll(baseQueryParams.Limit, baseQueryParams.Offset, baseQueryParams.Order, extractedConditionParams)
	if err != nil {
		http.Error(w, "Error finding data", http.StatusInternalServerError)
		return
	}
	// Convert data to AdminPanelSchema
	schemaSlice := *found.Data
	var adminSchemaSlice []models.AdminPanelSchema
	for _, item := range schemaSlice {
		// Append to schemaSlice
		adminSchemaSlice = append(adminSchemaSlice, item)
	}

	// fmt.Printf("Found: %v\n", *found.Data)
	fmt.Printf("SchemaSlice: %+v\n\n", schemaSlice)
	fmt.Printf("AdminSchemaSlice: %+v\n\n", adminSchemaSlice)
	fmt.Printf("TableHeaders: %v\n", c.tableHeaders)
	// Build the table data
	tableData := BuildTableData(adminSchemaSlice, found.Meta, c.adminHomeUrl, c.tableHeaders)

	// Generate Find All page render data
	data := GenerateFindAllRenderData(tableData, c.schemaName, c.pluralSchemaName, c.adminHomeUrl, searchQuery)

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}


func (c adminActionController) View(w http.ResponseWriter, r *http.Request) {
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

	// Find Current record
	found := &db.Action{}
	// Search for by ID and store in found
	found, err = c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s not found", c.schemaName), http.StatusNotFound)
		return
	}

	// Convert db struct to map for placeholder population
	currentData := getValuesUsingFieldMap(*found)
	// Populate form field placeholders with data from database
	err = populateValuessWithDBData(&editForm, currentData)
	if err != nil {
		http.Error(w, "Error generating form", http.StatusInternalServerError)
		return
	}

	data := GenerateEditRenderData(editForm, c.schemaName, c.pluralSchemaName, c.adminHomeUrl, stringParameter)	

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}



// Form generation
// Used to build Edit form
func (c adminActionController) generateEditForm() []FormField {
	return []FormField{
		{DbLabel: "Name", Label: "Name", Name: "name", Placeholder: "Enter name", Value: "", Type: "text", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Username", Label: "Username", Name: "username", Placeholder: "Enter username", Value: "", Type: "text", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Email", Label: "Email", Name: "email", Placeholder: "Enter email", Value: "", Type: "email", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Password", Label: "Password", Name: "password", Placeholder: "Enter password", Value: "", Type: "password", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Role", Label: "Role", Name: "role", Placeholder: "Enter role", Value: "user", Type: "select", Required: false, Disabled: false, Errors: []ErrorMessage{}, Selectors: RoleSelection()},
		{DbLabel: "Verified", Label: "Verified", Name: "verified", Placeholder: "", Value: "false", Type: "checkbox", Required: false, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "VerificationCode", Label: "Verification Code", Name: "verification_code", Placeholder: "Enter verification code", Value: "", Type: "text", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "VerificationCodeExpiry", Label: "Verification Code Expiry", Name: "verification_code_expiry", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "CreatedAt", Label: "Created At", Name: "created_at", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
		{DbLabel: "UpdatedAt", Label: "Updated At", Name: "updated_at", Placeholder: "", Value: "", Type: "datetime-local", Required: false, Disabled: true, Errors: []ErrorMessage{}},
	}
}