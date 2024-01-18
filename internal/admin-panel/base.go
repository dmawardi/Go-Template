package adminpanel

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dmawardi/Go-Template/internal/auth"
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
	ChangePasswordSuccess(w http.ResponseWriter, r *http.Request)
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
	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", PageRenderData{
		SidebarList: sidebar,
		PageType: PageType{
			HomePage: true,
		},
		// The section title is used on this page, to display login errors
		SectionTitle:  "Welcome to the admin panel",
		HeaderSection: header,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
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
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token", // Use the name of your auth cookie
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true, // Set to false if not using HTTPS
	})

	// Redirect to the login page, or return a success message
	http.Redirect(w, r, "/admin", http.StatusFound)
}
func (c adminBaseController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Init notification string
	notification := ""
	// Generate form
	passwordForm := c.generateChangePasswordForm()

	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Extract form data
		form, err := parseFormToMap(r)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		changePassword := models.ChangePassword{
			CurrentPassword:    form["currentPassword"],
			NewPassword:        form["newPassword"],
			ConfirmNewPassword: form["confirmNewPassword"],
		}
		// // Validate form data
		// // Validate struct
		pass, valErrors := helpers.GoValidateStruct(changePassword)

		// If validation passes
		if pass {
			// Perform password check
			// First, grab token from cookie
			// Validate the token
			tokenData, err := auth.ValidateAndParseToken(w, r)
			// If error detected
			if err != nil {
				http.Error(w, "Error parsing authentication token", http.StatusForbidden)
				return
			}
			// convert tokenData.UserID to int
			userID, err := strconv.Atoi(tokenData.UserID)
			if err != nil {
				// handle error
				fmt.Println("Error:", err)
			}

			// Find user using id found in token
			passMatch := c.service.CheckPasswordMatch(userID, []byte(changePassword.CurrentPassword))

			// If pasword match error is nil, and new password matches confirm new password
			if changePassword.NewPassword == changePassword.ConfirmNewPassword && passMatch {

				// Update the user's password
				_, err = c.service.Update(userID, &models.UpdateUser{Password: changePassword.ConfirmNewPassword})
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				// Redirect or render a success message
				http.Redirect(w, r, "/admin/change-password-success", http.StatusSeeOther)
				return
			} else {
				notification = "Old password does not match or new passwords do not match"
			}
		}

		// If validation fails
		if notification == "" {
			if changePassword.NewPassword != changePassword.ConfirmNewPassword {
				notification = "New passwords do not match"
			} else {
				// Else if validation fails, extract errors and manipulato for display
				newPasswordErrors := valErrors.Validation_errors["new_password"]
				errorSplit := strings.Split(newPasswordErrors[0], "does")
				errorString := fmt.Sprintf("Does%s", errorSplit[1])

				// Assign to notification
				notification = errorString
			}
		}
	}
	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", PageRenderData{
		SectionTitle: "Change Password",
		PageTitle:    "Change Password",
		// The section detail is used on this page, to display login errors
		SectionDetail: template.HTML("<p>" + notification + "</p>"),
		PageType: PageType{
			EditPage: true,
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: "/admin/change-password",
				FormMethod: "POST",
			},
			FormFields: passwordForm,
		},
		HeaderSection: header,
		SidebarList:   sidebar,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c adminBaseController) ChangePasswordSuccess(w http.ResponseWriter, r *http.Request) {
	serveAdminSuccess(w, "Change Password Success - Admin", "Change Password Success")
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
func (c adminBaseController) generateChangePasswordForm() []FormField {
	return []FormField{
		{DbLabel: "Password", Label: "Current Password", Name: "currentPassword", Placeholder: "", Value: "", Type: "password", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Password", Label: "New Password", Name: "newPassword", Placeholder: "", Value: "", Type: "password", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Password", Label: "Confirm new Password", Name: "confirmNewPassword", Placeholder: "", Value: "", Type: "password", Required: true, Disabled: false, Errors: []ErrorMessage{}},
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
