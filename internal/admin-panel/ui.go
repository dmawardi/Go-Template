package adminpanel

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
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
	TableHeaders   []string
	TableRows      []TableRow
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
// Function to build table data
func BuildTableData(listOfSchemaObjects []db.AdminPanelSchema, adminSchemaUrl string, tableHeaders []string) TableData {
	// Init table data
	tableData := TableData{
		AdminSchemaUrl: adminSchemaUrl,
		TableHeaders:   tableHeaders,
		TableRows:      []TableRow{},
	}

	// Loop through listOfSchemaObjects and build table rows
	for _, object := range listOfSchemaObjects {
		// Init table row
		row := TableRow{
			Data: []string{},
			// Fill in edit info
			Edit: EditInfo{
				EditUrl:   fmt.Sprintf("admin/%s/%s", adminSchemaUrl, object.GetID()),
				DeleteUrl: fmt.Sprintf("admin/%s/delete/%s", adminSchemaUrl, object.GetID()),
			},
		}

		// Append data based on the table headers
		for _, header := range tableHeaders {
			// Use header string values to get values from schema object
			row.Data = append(row.Data, object.ObtainValue(header))
		}

		// Append row to table data
		tableData.TableRows = append(tableData.TableRows, row)
	}

	return tableData
}
