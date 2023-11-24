package adminpanel

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/go-chi/chi"
)

type AdminUserController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	// Edit is also used to view the record details
	Edit(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type adminUserController struct {
}

func NewUserAdminController() AdminUserController {
	return &adminUserController{}
}

func (c adminUserController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Parse the template
	tmpl, err := parseAdminTemplates()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%+v\n", tmpl.DefinedTemplates())

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Admin User Home",
		SectionTitle: "Select a user to edit",
		SidebarList:  sidebarList,
		TableData: TableData{
			TableHeaders: []string{"ID", "Username", "Email"},
			TableRows: []TableRow{
				{Data: []string{"1", "admin", "admin@bulba.com"}},
				{Data: []string{"2", "admin", "admin@bulba.com"}},
				{Data: []string{"3", "admin", "admin@bulba.com"}},
			},
		},
		PageType: PageType{
			EditPage:   false,
			ReadPage:   true,
			CreatePage: false,
			DeletePage: false,
		},
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

func (c adminUserController) Create(w http.ResponseWriter, r *http.Request) {
	// Parse the template
	tmpl, err := parseAdminTemplates()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%+v\n", tmpl.DefinedTemplates())

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:   "Hello world",
		SidebarList: sidebarList,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: true,
			DeletePage: false,
		},
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

func (c adminUserController) Edit(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	fmt.Println("id parameter from request: ", stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Init a new user struct
	foundUser := &db.User{}
	// Search for user by ID and store in foundUser
	app.DbClient.Find(foundUser, idParameter)

	// Parse the template
	tmpl, err := parseAdminTemplates()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%+v\n", tmpl.DefinedTemplates())

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:   "Hello world",
		SidebarList: sidebarList,
		PageType: PageType{
			EditPage:   true,
			ReadPage:   false,
			CreatePage: false,
			DeletePage: false,
		},
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

func (c adminUserController) Delete(w http.ResponseWriter, r *http.Request) {
	// Parse the template
	tmpl, err := parseAdminTemplates()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%+v\n", tmpl.DefinedTemplates())

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:   "Hello world",
		SidebarList: sidebarList,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: false,
			DeletePage: true,
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: "/admin/users",
				FormMethod: "POST",
			},
			FormFields: []FormField{},
		},
	}

	// Execute the template with data and write to response
	err = tmpl.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

type UserEditableFields struct {
	Username string `json:"username,omitempty" valid:"length(6|25)"`
	Password string `json:"password,omitempty" valid:"length(6|30)"`
	Name     string `json:"name,omitempty" valid:"length(6|80)"`
	Email    string `json:"email,omitempty" valid:"email"`
	Verified bool   `json:"verified,omitempty" valid:"bool"`
}
