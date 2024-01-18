package adminpanel

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Helper function to format time.Time fields for forms
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339) // or another suitable format
}

// Helper function to convert list of strings to list of ints
func convertStringSliceToIntSlice(stringSlice []string) ([]int, error) {
	intSlice := make([]int, 0, len(stringSlice)) // Create a slice of ints with the same length

	for _, str := range stringSlice {
		num, err := strconv.Atoi(str) // Convert string to int
		if err != nil {
			return nil, err
		}
		intSlice = append(intSlice, num) // Append the converted int to the slice
	}
	return intSlice, nil
}

// Parses all the template files in the templates directory
func ParseAdminTemplates() (*template.Template, error) {
	// Parse the base template
	tmpl := template.New("/internal/admin-panel/templates/layout.tmpl")

	// Walk through all files in the templates directory
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// If the file is not a directory and has the .html extension
		if !info.IsDir() && filepath.Ext(path) == ".tmpl" {
			// Parse the file
			_, err = tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})

	// Return error if there is filepath walk issue
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// Function to render the Admin error page to the response
func serveAdminError(w http.ResponseWriter, sectionTitle string) {
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Error - Admin",
		SectionTitle: sectionTitle,
		SidebarList:  sidebar,
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

// Function to render the Admin success page to the response
func serveAdminSuccess(w http.ResponseWriter, pageTitle string, sectionTitle string) {
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    pageTitle,
		SectionTitle: sectionTitle,
		SidebarList:  sidebar,
		PageType: PageType{
			EditPage:    false,
			ReadPage:    false,
			CreatePage:  false,
			DeletePage:  false,
			SuccessPage: true,
		},
		FormData:      FormData{},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// Checks if a string contains another string (Used to search for resource)
func containsString(s, searchTerm string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(searchTerm))
}

// Checks if array contains data
func arrayContainsString(arr []string, target string) bool {
	for _, s := range arr {
		if s == target {
			return true
		}
	}
	return false
}

// Function to sort permissions data from enforcer
func sortByRoleResourceAlphabetically(a, b map[string]interface{}) bool {
	resourceA, okA := a["resource"].(string)
	resourceB, okB := b["resource"].(string)

	// If either of the elements doesn't have a valid "resource" string, consider it greater (move it to the end)
	if !okA || !okB {
		return false
	}

	// Compare the "resource" strings alphabetically
	return resourceA < resourceB
}

func sortByKeyAlphabetically(a, b map[string]string, key string) bool {
	valueA, okA := a[key]
	valueB, okB := b[key]

	// If either of the elements doesn't have a valid string for the given key, consider it greater (move it to the end)
	if !okA || !okB {
		return false
	}

	// Compare the strings alphabetically
	return valueA < valueB
}

// Edit a slice of table rows to add row span to first column and remove <td> tags from
// subsequent rows
func editTableDataRowSpan(tableRows []TableRow) {
	var lastRecordedStart struct {
		resource string
		index    int
	}
	for i, row := range tableRows {
		// Extract row variables
		rowData := row.Data
		currentResource := rowData[0].Label

		// If current resource is different from last recorded resource, then must edit row span
		if currentResource != lastRecordedStart.resource {
			// Calculate row span
			rowSpan := i - lastRecordedStart.index

			// If the difference between the current index and the last recorded index is greater than 1, then must edit row span
			if i-lastRecordedStart.index > 1 {
				// Add row span to first cell of last recorded start row
				tableRows[lastRecordedStart.index].Data[0].RowSpan = rowSpan

				// Remove <td> tags from subsequent rows: count from last recorded index + 1 till before current index
				for j := lastRecordedStart.index + 1; j < i; j++ {
					// Chop first element off of data
					tableRows[j].Data = tableRows[j].Data[1:]
				}
			}
			// Reassign lastRecordedResource details
			lastRecordedStart.resource = currentResource
			lastRecordedStart.index = i

		}
	}
}

// Convert json string to map
func jsonToMap(jsonStr string) (map[string]string, error) {
	var result map[string]string
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Convert struct string to map. Struct string will be in format: "[key1]:value1|[key2]:value2"
func stringToMap(input string) (map[string]string, error) {
	result := make(map[string]string)

	// Split the input string by "|"
	parts := strings.Split(input, "|")

	for _, part := range parts {
		// Check if the part has "[]" to identify a key name
		keyValueSlice := strings.Split(part, ":")

		// If a key value pair is found
		if len(keyValueSlice) == 2 {
			// Grab the first item in slice as key, and remove the "[" and "]" characters
			key := keyValueSlice[0]
			// // Grab the second item in slice as value
			value := keyValueSlice[1]
			// Add key value pair to result map
			result[key[1:len(key)-1]] = value
		}

	}

	return result, nil
}

// parseFormToMap parses the form data and converts it into a map[string]string
func parseFormToMap(r *http.Request) (map[string]string, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	formMap := make(map[string]string)
	for key, values := range r.Form { // range over map
		// In form data, key can have multiple values,
		// we'll take the first one only
		formMap[key] = values[0]
	}

	return formMap, nil
}
