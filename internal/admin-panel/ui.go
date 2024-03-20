package adminpanel

import (
	"html/template"
)

// Data used to render each page
// Contains state for the page
type PageRenderData struct {
	// In HEAD
	PageTitle string
	// In BODY
	SectionTitle  string
	SectionDetail template.HTML
	SidebarList   AdminSideBar
	// Schema home used to return to the schema home page from delete
	SchemaHome string // eg. /admin/users/
	// Page type (Used for content selection)
	PageType PageType
	// Form
	FormData  FormData
	TableData TableData
	// Search
	SearchTerm             string
	RecordsPerPageSelector []int
	// Special section data for policies
	PolicySection PolicySection
	HeaderSection HeaderSection
}

// Variables for header
type HeaderSection struct {
	HomeUrl           template.URL
	ViewSiteUrl       template.URL
	ChangePasswordUrl template.URL
	LogOutUrl         template.URL
}

// Variables for policy section
type PolicySection struct {
	FocusedPolicies []PolicyEditDataRow
	PolicyResource  string
	Selectors       PolicyEditSelectors
}

// Page type (Used for dynamic selective rendering)
type PageType struct {
	HomePage    bool
	EditPage    bool
	ReadPage    bool
	CreatePage  bool
	DeletePage  bool
	SuccessPage bool
	// Used for policy section
	PolicyMode string // eg. "policy" or "inheritance"
}
