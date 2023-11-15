package adminpanel

import (
	"fmt"
	"html/template"
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
		Content     template.HTML
		SidebarList []string
	}{
		Title:       "Hello world",
		Content:     template.HTML("<h1>This is the content</h1>"),
		SidebarList: sidebarList,
	}

	// Execute the template with data and write to response
	tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (c adminUserController) UserDetail(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the User Detail page"))
}

func (c adminUserController) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the Create User page"))
}

func (c adminUserController) Edit(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the Edit User page"))
}
