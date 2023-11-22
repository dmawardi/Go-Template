package adminpanel

import (
	"fmt"
	"net/http"
)

type AdminUserController interface {
	AllUsers(w http.ResponseWriter, r *http.Request)
	UserDetail(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Edit(w http.ResponseWriter, r *http.Request)
}

type adminUserController struct {
}

func NewUserAdminController() AdminUserController {
	return &adminUserController{}
}

func (c adminUserController) AllUsers(w http.ResponseWriter, r *http.Request) {
	// Parse the template
	tmpl, err := parseAdminTemplates()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%+v\n", tmpl.DefinedTemplates())

	// Data to be injected into template
	data := struct {
		Title       string
		SidebarList []string
		// Form
		FormPage bool
		FormData FormData
	}{
		Title:       "Hello world",
		SidebarList: sidebarList,
		FormPage:    true,
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: "/admin/users",
				FormMethod: "POST",
			},
			FormFields: []FormField{
				{
					Label:       "Username",
					Name:        "username",
					Placeholder: "Cilandak 213",
					Value:       "",
					Type:        "text",
					Required:    true,
					Disabled:    false,
					Errors:      []ErrorMessage{{Message: "This is an error message"}},
				},
				{
					Label:       "Password",
					Name:        "password",
					Placeholder: "",
					Value:       "",
					Type:        "text",
					Required:    true,
					Disabled:    true,
					Errors:      []ErrorMessage{{Message: "This is an error message"}},
				},
			},
		},
	}

	// Execute the template with data and write to response
	err = tmpl.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c adminUserController) UserDetail(w http.ResponseWriter, r *http.Request) {
	// Parse the template
	tmpl, err := parseAdminTemplates()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%+v\n", tmpl.DefinedTemplates())

	// Data to be injected into template
	data := struct {
		Title       string
		SidebarList []string
		FormPage    bool
		FormFields  []FormField
	}{
		Title:       "Hello world",
		SidebarList: sidebarList,
		FormPage:    false,
		FormFields: []FormField{
			{
				Label:       "Username",
				Name:        "username",
				Placeholder: "Enter username",
				Value:       "",
				Type:        "text",
				Required:    true,
			},
		},
	}

	// Execute the template with data and write to response
	tmpl.ExecuteTemplate(w, "layout.tmpl", data)
}

func (c adminUserController) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the Create User page"))
}

func (c adminUserController) Edit(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the Edit User page"))
}
