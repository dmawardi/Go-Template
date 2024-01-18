package adminpanel

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

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
