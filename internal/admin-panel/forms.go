package adminpanel

// Go Template components

type FormData struct {
	FormFields  []FormField
	FormDetails FormDetails
}

// Data for each form field
type FormField struct {
	// Label above input
	Label   string
	DbLabel string
	// Used for form submission
	Name string
	// Is this field required?
	Required bool
	// Is this field disabled?
	Disabled bool

	// Current value
	Value string
	// Silhouette
	Placeholder string
	Type        string
	Errors      []ErrorMessage
	Selectors   []FormFieldSelector
}

type FormFieldSelector struct {
	Value string
	Label string
}

// Display of errors in form
type ErrorMessage struct {
	Message string
}

type UserEditForm struct {
	Title  string
	Fields []FormField
}

// Form Details
type FormDetails struct {
	FormAction string
	FormMethod string
}
