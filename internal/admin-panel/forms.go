package adminpanel

// Go Template components
type FormField struct {
	Label string
	// Current value
	Placeholder string
	InputType   string
	FieldType   string
}

type UserEditForm struct {
	Title  string
	Fields []FormField
}
