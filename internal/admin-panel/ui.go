package adminpanel

import (
	"fmt"
	"net/http"

	"github.com/dmawardi/Go-Template/internal/models"
)

// Data used to render each page
// Contains state for the page
type PageRenderData struct {
	// In HEAD
	PageTitle string
	// In BODY
	SectionTitle string
	SidebarList  []sidebarItem
	// Schema home used to return to the schema home page from delete
	SchemaHome string // eg. /admin/users/
	// Page type (Used for content selection)
	PageType PageType
	// Form
	FormData  FormData
	TableData TableData
	// Search
	SearchTerm string
}

// Page type (Used for content selection)
type PageType struct {
	EditPage    bool
	ReadPage    bool
	CreatePage  bool
	DeletePage  bool
	SuccessPage bool
}

// Data table
type TableData struct {
	AdminSchemaUrl string // eg. /users/
	TableHeaders   []TableHeader
	TableRows      []TableRow
	MetaData       models.ExtendedSchemaMetaData
}

// Used for table header information. Also holds information for sorting and pointer + data type
type TableHeader struct {
	Label string
	// label used in db
	ColumnSortLabel string
	// Is data type a pointer
	Pointer bool
	// Used for pointer to string extraction
	DataType string
}

// Data to complete a table row
type TableRow struct {
	Data []string
	Edit EditInfo
}

// Edit info for the Edit column in the table
type EditInfo struct {
	EditUrl   string // eg. admin/users/1
	DeleteUrl string // eg. admin/users/delete/1
}

// Form data
// Function to build table data from slice of adminpanel schema objects, admin schema url (eg. /admin/users) and table headers
func BuildTableData(listOfSchemaObjects []AdminPanelSchema, metaData models.SchemaMetaData, adminSchemaBaseUrl string, tableHeaders []TableHeader) TableData {
	// Calculate currently showing records and total pages
	currentlyShowing := metaData.CalculateCurrentlyShowingRecords()
	// Init table data
	tableData := TableData{
		AdminSchemaUrl: adminSchemaBaseUrl,
		TableHeaders:   tableHeaders,
		TableRows:      []TableRow{},
		// Build extended metadata
		MetaData: models.NewExtendedSchemaMetaData(metaData, currentlyShowing),
	}

	// Loop through listOfSchemaObjects and build table rows
	for _, object := range listOfSchemaObjects {
		// Init table row
		row := TableRow{
			Data: []string{},
			// Fill in edit info
			Edit: EditInfo{
				EditUrl:   fmt.Sprintf("%s/%s", adminSchemaBaseUrl, object.GetID()),
				DeleteUrl: fmt.Sprintf("%s/delete/%s", adminSchemaBaseUrl, object.GetID()),
			},
		}

		// Iterate through tableheaders
		for _, header := range tableHeaders {
			// Grab data from the schema object
			fieldData := object.ObtainValue(header.Label)

			// convert fieldData to string
			stringFieldData := fmt.Sprint(fieldData)

			// Use header string values to get values from schema object and append
			row.Data = append(row.Data, stringFieldData)
		}

		// Append row to table data
		tableData.TableRows = append(tableData.TableRows, row)
	}

	return tableData
}

// Function to render the Admin error page to the response
func serveAdminError(w http.ResponseWriter, sectionTitle string) {
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Error - Admin",
		SectionTitle: sectionTitle,
		SidebarList:  sidebarList,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: false,
			DeletePage: true,
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
