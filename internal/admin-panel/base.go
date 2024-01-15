package adminpanel

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/service"
)

// Admin base controller (non-schema related routes)
type AdminBaseController interface {
	Home(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
}

type adminBaseController struct {
	service service.UserService
}

// Constructor
func NewAdminBaseController(userService service.UserService) AdminBaseController {
	return &adminBaseController{userService}
}

// RECEIVER FUNCTIONS
func (c adminBaseController) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the admin main home page"))
}

// Admin login page
func (c adminBaseController) Login(w http.ResponseWriter, r *http.Request) {
	// Init token string
	tokenString := ""
	loginErrorMsg := ""

	// Generate form
	loginForm := c.generateLoginForm()

	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Extract form data
		login, err := c.extractLoginFormSubmission(r)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// Validate form data
		// Validate struct
		pass, valErrors := helpers.GoValidateStruct(login)

		// If validation passes
		if pass {
			// Login user
			tokenString, err = c.service.LoginUser(&login)
			if err == nil {
				// Set token in cookie
				createAndSetHeaderCookie(w, tokenString)

				// Redirect or render a success message
				http.Redirect(w, r, "/admin/home", http.StatusSeeOther)
				return
			}

			// Else if login fails
			fmt.Printf("Error logging in for email: %s\n", login.Email)
			loginErrorMsg = "Invalid email or password"

		}

		// Else if validation fails
		// Populate form field errors
		SetValidationErrorsInForm(loginForm, *valErrors)

		// Extract form submission from request and build into map[string]string
		formFieldMap, err := extractLoginFormAsMap(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Populate previously entered values (Avoids password)
		err = populateValuesWithFormName(&loginForm, formFieldMap)
		if err != nil {
			http.Error(w, "Error populating form", http.StatusInternalServerError)
			return
		}
	}
	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "login.tmpl", PageRenderData{
		// The section title is used on this page, to display login errors
		SectionTitle: loginErrorMsg,
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: "/admin/login",
				FormMethod: "POST",
			},
			FormFields: loginForm,
		},
		HeaderSection: header,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminBaseController) Logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the admin logout page"))
}
func (c adminBaseController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the admin logout page"))
}

func (c adminBaseController) extractLoginFormSubmission(r *http.Request) (models.Login, error) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		return models.Login{}, errors.New("Error parsing form")
	}

	// Build struct for validation
	toValidate := models.Login{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	return toValidate, nil
}

func (c adminBaseController) generateLoginForm() []FormField {
	return []FormField{
		{DbLabel: "Email", Label: "Email", Name: "email", Placeholder: "", Value: "", Type: "text", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Password", Label: "Password", Name: "password", Placeholder: "", Value: "", Type: "password", Required: true, Disabled: false, Errors: []ErrorMessage{}},
	}
}

func createAndSetHeaderCookie(w http.ResponseWriter, tokenString string) {
	// Create the cookie
	expire := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:     "jwt_token",
		Value:    tokenString,
		Expires:  expire,
		HttpOnly: true,
		Secure:   true, // Set to false if not using HTTPS
		Path:     "/",
	}

	// Set the cookie in the response header
	http.SetCookie(w, &cookie)
}

func extractLoginFormAsMap(r *http.Request) (map[string]string, error) {
	// Extract form submission from request and build into map[string]string
	formFieldMap := make(map[string]string)
	for k, v := range r.Form {
		formFieldMap[k] = v[0]
	}
	return formFieldMap, nil
}
