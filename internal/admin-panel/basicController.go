package adminpanel

import (
	"net/http"
	"reflect"

	"github.com/dmawardi/Go-Template/internal/service"
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
}

// RECEIVER FUNCTIONS
func (b basicAdminController) ObtainUrlDetails() URLDetails {
	return URLDetails{
		AdminHomeUrl:     b.AdminHomeUrl,
		SchemaName:       b.SchemaName,
		PluralSchemaName: b.PluralSchemaName,
	}
}

// ADMIN SIDEBAR CREATION
//
// Uses the ObtainUrlDetails method to get the sidebar details of any Basic Admin Controller type
func ObtainUrlDetailsForBasicAdminController(input interface{}) URLDetails {
	// Use reflection to call ObtainUrlDetails method if it exists.
	value := reflect.ValueOf(input)
	// ObtainUrlDetails method
	method := value.MethodByName("ObtainUrlDetails")
	if !method.IsValid() {
		return URLDetails{}
	}

	// Call ObtainUrlDetails method
	result := method.Call(nil)

	// Check if result is valid (if it is a BasicAdminController)
	// If it has a result
	if len(result) == 1 {
		// Assign the result as an interface to resultFields
		interfaceFields := result[0].Interface()
		// Assign the fields of the resultFields to sidebarDetails
		sidebarDetails := URLDetails{
			AdminHomeUrl:     interfaceFields.(URLDetails).AdminHomeUrl,
			SchemaName:       interfaceFields.(URLDetails).SchemaName,
			PluralSchemaName: interfaceFields.(URLDetails).PluralSchemaName,
		}
		return sidebarDetails
	}

	return URLDetails{}
}

type URLDetails struct {
	AdminHomeUrl string
	SchemaName   string
	PluralSchemaName string
}


// func (c basicAdminController) FindAll(w http.ResponseWriter, r *http.Request) {
// 	// Grab query parameters
// 	searchQuery := r.URL.Query().Get("search")
// 	// Grab basic query params
// 	baseQueryParams, err := request.ExtractBasicFindAllQueryParams(r)
// 	if err != nil {
// 		http.Error(w, "Error extracting query params", http.StatusBadRequest)
// 		return
// 	}

// 	// Generate query params to extract
// 	queryParamsToExtract := c.ConditionQueryParams
// 	// Extract query params
// 	extractedConditionParams, err := request.ExtractSearchAndConditionParams(r, queryParamsToExtract)
// 	if err != nil {
// 		fmt.Println("Error extracting conditions: ", err)
// 		http.Error(w, "Can't find conditions", http.StatusBadRequest)
// 		return
// 	}

// 	// Find all with options from database
// 	found, err := c.service.FindAll(baseQueryParams.Limit, baseQueryParams.Offset, baseQueryParams.Order, extractedConditionParams)
// 	if err != nil {
// 		http.Error(w, "Error finding data", http.StatusInternalServerError)
// 		return
// 	}

// 	// Convert data to AdminPanelSchema
// 	schemaSlice := *found.Data
// 	var adminSchemaSlice []models.AdminPanelSchema
// 	for _, schema := range schemaSlice {
// 		// Append to schemaSlice
// 		adminSchemaSlice = append(adminSchemaSlice, schema)
// 	}

// 	// Build the table data
// 	tableData := BuildTableData(adminSchemaSlice, found.Meta, c.adminHomeUrl, c.tableHeaders)

// 	// Data to be injected into template
// 	data := PageRenderData{
// 		PageTitle:              "Admin: " + c.pluralSchemaName,
// 		SectionTitle:           fmt.Sprintf("Select a %s to edit", c.schemaName),
// 		SidebarList:            sidebar,
// 		TableData:              tableData,
// 		SchemaHome:             c.adminHomeUrl,
// 		SearchTerm:             searchQuery,
// 		RecordsPerPageSelector: recordsPerPage,
// 		PageType: PageType{
// 			EditPage:   false,
// 			ReadPage:   true,
// 			CreatePage: false,
// 			DeletePage: false,
// 		},
// 		FormData: FormData{
// 			FormDetails: FormDetails{
// 				FormAction: c.adminHomeUrl,
// 				FormMethod: "get",
// 			},
// 			FormFields: []FormField{},
// 		},
// 		HeaderSection: header,
// 	}

// 	// Execute the template with data and write to response
// 	err = app.AdminTemplates.ExecuteTemplate(w, "layout.go.tmpl", data)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// }