package adminpanel

import "net/http"

// Admin base controller (non-schema related routes)
type AdminBaseController interface {
	Home(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
}

type adminBaseController struct {
}

// Constructor
func NewAdminBaseController() AdminBaseController {
	return &adminBaseController{}
}

// RECEIVER FUNCTIONS
func (c adminBaseController) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the admin main home page"))
}
func (c adminBaseController) Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the admin login page"))
}
func (c adminBaseController) Logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the admin logout page"))
}
func (c adminBaseController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the admin logout page"))
}
