package adminpanel

import (
	"fmt"
	"reflect"
)

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
	EditUrl   string
	DeleteUrl string
}

// Goes through a list of structs and returns a list of strings based on input slice
func GetStructFieldValues(listOfData []interface{}, listOfTableHeaders []string) [][]string {
	// Init a nested array to hold values
	var values [][]string
	// Get values for each row
	for _, rowItem := range listOfData {
		var row []string
		// Get values for each header
		row = getDynamicFieldValues(rowItem, listOfTableHeaders)

		// Add row to values slice
		values = append(values, row)
	}
	return values
}

// Function to get dynamic field values in the specified order
func getDynamicFieldValues(obj interface{}, fieldNames []string) []string {
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Struct {
		return nil
	}

	numFields := len(fieldNames)
	values := make([]string, numFields)

	for i := 0; i < numFields; i++ {
		fieldName := fieldNames[i]
		field := value.FieldByName(fieldName)

		if !field.IsValid() {
			values[i] = "" // Use an empty string for missing fields
		} else {
			values[i] = fmt.Sprint(field.Interface()) // Convert to string
		}
	}

	return values
}
