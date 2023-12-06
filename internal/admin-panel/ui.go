package adminpanel

import (
	"fmt"

	"github.com/dmawardi/Go-Template/internal/db"
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
	MetaData       models.ExtendedSchemaMetaData
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
func BuildTableData(listOfSchemaObjects []db.AdminPanelSchema, metaData models.SchemaMetaData, adminSchemaBaseUrl string, tableHeaders []string) TableData {
	// Init table data
	tableData := TableData{
		AdminSchemaUrl: adminSchemaBaseUrl,
		TableHeaders:   tableHeaders,
		TableRows:      []TableRow{},
		// Build extended metadata
		MetaData: models.NewExtendedSchemaMetaData(metaData, metaData.CalculateTotalPages(), metaData.CalculateCurrentlyShowingRecords()),
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
			// Use header string values to get values from schema object and append
			row.Data = append(row.Data, object.ObtainValue(header))
		}

		// Append row to table data
		tableData.TableRows = append(tableData.TableRows, row)
	}

	return tableData
}
