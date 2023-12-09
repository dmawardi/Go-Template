package adminpanel

import (
	"fmt"

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
			fmt.Printf("header: %v \nField data: %v\n", header.Label, fieldData)

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

type BulkDeleteRequest struct {
	SelectedItems []string `json:"selected_items"`
}
