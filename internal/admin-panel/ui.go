package adminpanel

import (
	"fmt"
	"strings"

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
	// Used for customized forms for roles and users
	Mode string // eg. "roles", "users", or "general" (all other shcmeas)
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
	Data []TableCell
	Edit EditInfo
}

// Data for a single cell
type TableCell struct {
	Label   string
	RowSpan int
	// Primarily used for the policy table
	EditLink string
}

// Edit info for the Edit column in the table
type EditInfo struct {
	EditUrl   string // eg. admin/users/1
	DeleteUrl string // eg. admin/users/delete/1
}

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
			Data: []TableCell{},
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
			row.Data = append(row.Data, TableCell{Label: stringFieldData})
		}

		// Append row to table data
		tableData.TableRows = append(tableData.TableRows, row)
	}

	return tableData
}

// Function build table data for Permissions
func BuildRolesTableData(policySlice []map[string]interface{}, adminSchemaBaseUrl string, tableHeaders []TableHeader) TableData {
	var tableRows []TableRow

	// Loop through policy slice to build table rows
	for _, policy := range policySlice {
		var rowData []TableCell

		// Iterate through tableheaders
		for _, header := range tableHeaders {
			// Grab data from the schema object
			value, found := policy[header.Label]

			// If the key is found, append the value to the row data
			if found {
				// Append with edit link if it's the first column (resource)
				if header.Label == "resource" {
					// Create the edit link from the label value
					editLink := strings.ReplaceAll(value.(string), "/", "-")
					// Append to row data with edit link
					rowData = append(rowData, TableCell{Label: fmt.Sprintf("%v", value), EditLink: editLink})

					// else if other column
				} else {
					rowData = append(rowData, TableCell{Label: fmt.Sprintf("%v", value), EditLink: ""})
				}

				// If the key is not found, append an empty string
			} else {
				rowData = append(rowData, TableCell{Label: ""}) // Add an empty string if the key is not found
			}
		}

		// Append to table rows
		tableRows = append(tableRows, TableRow{Data: rowData})
	}
	return TableData{
		AdminSchemaUrl: adminSchemaBaseUrl, // You can set this value as needed
		TableHeaders:   tableHeaders,
		TableRows:      tableRows,
	}
}
