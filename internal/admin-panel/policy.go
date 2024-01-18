package adminpanel

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/go-chi/chi"
)

// Table headers to show on find all pages
var authPolicyTableHeaders = []TableHeader{
	{Label: "resource", ColumnSortLabel: "resource", Pointer: false, DataType: "string"},
	{Label: "role", ColumnSortLabel: "role", Pointer: false, DataType: "string"},
	{Label: "action", ColumnSortLabel: "action", Pointer: false, DataType: "string"},
}

var inheritanceTableHeaders = []TableHeader{
	{Label: "role", ColumnSortLabel: "role", Pointer: false, DataType: "string", Sortable: false},
	{Label: "inherits_from", ColumnSortLabel: "inherits_from", Pointer: false, DataType: "string", Sortable: false},
}

var roleTableHeaders = []TableHeader{
	{Label: "role", ColumnSortLabel: "role", Pointer: false, DataType: "string", Sortable: false},
}

// Constructor
func NewAdminAuthPolicyController(service service.AuthPolicyService, selectorService SelectorService) AdminAuthPolicyController {
	return &adminAuthPolicyController{
		service: service,
		// Use values from above
		adminHomeUrl:            "/admin/policy",
		schemaName:              "Policy",
		pluralSchemaName:        "Policies",
		tableHeaders:            authPolicyTableHeaders,
		inheritanceTableHeaders: inheritanceTableHeaders,
		roleTableheaders:        roleTableHeaders,
		formSelectors:           selectorService,
	}
}

type AdminAuthPolicyController interface {
	// Policy
	FindAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	CreateSuccess(w http.ResponseWriter, r *http.Request)
	// Edit is also used to view the record details
	Edit(w http.ResponseWriter, r *http.Request)
	// Roles
	FindAllRoles(w http.ResponseWriter, r *http.Request)
	CreateRole(w http.ResponseWriter, r *http.Request)
	CreateRoleSuccess(w http.ResponseWriter, r *http.Request)
	// Inheritance
	FindAllRoleInheritance(w http.ResponseWriter, r *http.Request)
	CreateInheritance(w http.ResponseWriter, r *http.Request)
	DeleteInheritance(w http.ResponseWriter, r *http.Request)
	CreateInheritanceSuccess(w http.ResponseWriter, r *http.Request)
	DeleteInheritanceSuccess(w http.ResponseWriter, r *http.Request)
	// For sidebar
	ObtainFields() BasicAdminController
}
type adminAuthPolicyController struct {
	service service.AuthPolicyService
	// For links
	adminHomeUrl string
	// For HTML text rendering
	schemaName       string
	pluralSchemaName string
	// Custom table headers
	tableHeaders            []TableHeader
	inheritanceTableHeaders []TableHeader
	roleTableheaders        []TableHeader

	// Form selectors
	formSelectors SelectorService
}

// CRUD handlers
func (c adminAuthPolicyController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab query parameters
	searchQuery := r.URL.Query().Get("search")

	// Find all with options from database
	groupsSlice, err := c.service.FindAll()
	if err != nil {
		http.Error(w, "Error finding data", http.StatusInternalServerError)
		return
	}
	// Filter by search query
	groupsSlice = searchPoliciesByResource(groupsSlice, searchQuery)

	// Sort by resource alphabetically
	sort.Slice(groupsSlice, func(i, j int) bool {
		// Give two items to compare to role resource alpha sorter
		return sortByRoleResourceAlphabetically(groupsSlice[i], groupsSlice[j])
	})

	// Build the roles table data
	tableData := BuildPolicyTableData(groupsSlice, c.adminHomeUrl, c.tableHeaders)
	// Add the row span attribute to the table based on resource grouping
	editTableDataRowSpan(tableData.TableRows)

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Admin: " + c.pluralSchemaName,
		SectionTitle: fmt.Sprintf("Select a %s to edit", c.schemaName),
		SidebarList:  sidebar,
		TableData:    tableData,
		SchemaHome:   c.adminHomeUrl,
		SearchTerm:   searchQuery,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   true,
			CreatePage: false,
			DeletePage: false,
			PolicyMode: "policy",
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: c.adminHomeUrl,
				FormMethod: "get",
			},
			FormFields: []FormField{},
		},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "policy.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminAuthPolicyController) FindAllRoleInheritance(w http.ResponseWriter, r *http.Request) {
	// Grab query parameters
	searchQuery := r.URL.Query().Get("search")

	// Find all with options from database
	inheritanceSlice, err := c.service.FindAllRoleInheritance()
	if err != nil {
		http.Error(w, "Error finding data", http.StatusInternalServerError)
		return
	}
	// // Filter by search query
	inheritanceSlice = searchMapUsingKeys(inheritanceSlice, []string{"inherits_from", "role"}, searchQuery)

	// Sort by resource alphabetically
	sort.Slice(inheritanceSlice, func(i, j int) bool {
		// Give two items to compare to alphabetic sorter
		return sortByKeyAlphabetically(inheritanceSlice[i], inheritanceSlice[j], "role")
	})

	// // Build the roles table data
	tableData := BuildRoleInheritanceTableData(inheritanceSlice, c.adminHomeUrl, c.inheritanceTableHeaders)

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Admin: " + c.pluralSchemaName,
		SectionTitle: fmt.Sprintf("Select a %s to edit", c.schemaName),
		SidebarList:  sidebar,
		TableData:    tableData,
		SchemaHome:   c.adminHomeUrl,
		SearchTerm:   searchQuery,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   true,
			CreatePage: false,
			DeletePage: false,
			PolicyMode: "inheritance",
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: c.adminHomeUrl + "/inheritance",
				FormMethod: "get",
			},
			FormFields: []FormField{},
		},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "policy.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func (c adminAuthPolicyController) FindAllRoles(w http.ResponseWriter, r *http.Request) {
	// Grab query parameters
	searchQuery := r.URL.Query().Get("search")

	// Find all with options from database
	rolesSlice, err := c.service.FindAllRoles()
	if err != nil {
		http.Error(w, "Error finding data", http.StatusInternalServerError)
		return
	}

	filteredSlice := []string{}
	// Iterate through roles slice and remove items that do not match search query
	for _, role := range rolesSlice {
		if containsString(role, searchQuery) {
			filteredSlice = append(filteredSlice, role)
		}
	}

	// // Build the roles table data
	tableData := BuildRoleTableData(filteredSlice, c.adminHomeUrl, c.roleTableheaders)

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    "Admin: " + c.pluralSchemaName,
		SectionTitle: fmt.Sprintf("Select a %s to edit", c.schemaName),
		SidebarList:  sidebar,
		TableData:    tableData,
		SchemaHome:   c.adminHomeUrl,
		SearchTerm:   searchQuery,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   true,
			CreatePage: false,
			DeletePage: false,
			PolicyMode: "policy",
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: c.adminHomeUrl + "/roles",
				FormMethod: "get",
			},
			FormFields: []FormField{},
		},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "policy.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c adminAuthPolicyController) Create(w http.ResponseWriter, r *http.Request) {
	// Init new form
	createForm := c.generateCreateForm()

	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Extract user form submission
		toValidate, err := c.extractCreateFormSubmission(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Validate struct
		pass, valErrors := helpers.GoValidateStruct(toValidate)
		// If failure detected
		// If validation passes
		if pass {
			// Create
			err = c.service.Create(toValidate)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error creating %s", c.schemaName), http.StatusInternalServerError)
				return
			}
			// Redirect or render a success message
			http.Redirect(w, r, fmt.Sprintf("%s/create/success", c.adminHomeUrl), http.StatusSeeOther)
			return
		}

		// If validation fails
		// Populate form field errors
		SetValidationErrorsInForm(createForm, *valErrors)

		// Extract form submission from request and build into map[string]string
		formFieldMap, err := c.extractFormFromRequest(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		fmt.Printf("formFieldMap: %+v\n", formFieldMap)
	}

	// Render preparation
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    fmt.Sprintf("Create %s", c.schemaName),
		SectionTitle: fmt.Sprintf("Create a new %s", c.schemaName),
		SidebarList:  sidebar,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: true,
			DeletePage: false,
			PolicyMode: "policy",
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: fmt.Sprintf("%s/create", c.adminHomeUrl),
				FormMethod: "post",
			},
			FormFields: createForm,
		},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "policy.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c adminAuthPolicyController) CreateRole(w http.ResponseWriter, r *http.Request) {
	// Init new form
	createForm := c.generateCreateRoleForm()

	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Extract user form submission
		submittedForm, err := c.extractCreateRoleFormSubmission(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Validate struct
		pass, valErrors := helpers.GoValidateStruct(submittedForm)
		// If failure detected
		// If validation passes
		if pass {
			// Create
			success, err := c.service.AssignUserRole(submittedForm.UserId, submittedForm.Role)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error assigning role %s", c.schemaName), http.StatusInternalServerError)
				return
			}
			if !*success {
				http.Error(w, fmt.Sprintf("Error assigning role %s", c.schemaName), http.StatusInternalServerError)
				return
			}
			// Redirect or render a success message
			http.Redirect(w, r, fmt.Sprintf("%s/create-role/success", c.adminHomeUrl), http.StatusSeeOther)
			return
		}

		// If validation fails
		// Populate form field errors
		SetValidationErrorsInForm(createForm, *valErrors)

		// Extract form submission from request and build into map[string]string
		formFieldMap, err := c.extractFormFromRequest(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		fmt.Printf("formFieldMap: %+v\n", formFieldMap)
	}

	// Render preparation
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    fmt.Sprintf("Create %s Role", c.schemaName),
		SectionTitle: fmt.Sprintf("Create a new %s Role", c.schemaName),
		SidebarList:  sidebar,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: true,
			DeletePage: false,
			PolicyMode: "policy",
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: fmt.Sprintf("%s/create-role", c.adminHomeUrl),
				FormMethod: "post",
			},
			FormFields: createForm,
		},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "policy.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c adminAuthPolicyController) CreateInheritance(w http.ResponseWriter, r *http.Request) {
	// Init new form
	createForm := c.generateCreateInheritanceForm()

	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Extract user form submission
		submittedForm, err := c.extractCreateInheritanceFormSubmission(r)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		fmt.Printf("submittedForm: %+v\n", submittedForm)

		// Validate struct
		pass, valErrors := helpers.GoValidateStruct(submittedForm)
		// If failure detected
		// If validation passes
		if pass {
			// Create
			err := c.service.CreateInheritance(models.G2Record{Role: submittedForm.Role, InheritsFrom: submittedForm.InheritsFrom})
			if err != nil {
				http.Error(w, fmt.Sprintf("Error assigning role %s", c.schemaName), http.StatusInternalServerError)
				return
			}

			// Redirect or render a success message
			http.Redirect(w, r, fmt.Sprintf("%s/create-inheritance/success", c.adminHomeUrl), http.StatusSeeOther)
			return
		}

		// If validation fails
		// Populate form field errors
		SetValidationErrorsInForm(createForm, *valErrors)
	}

	// Render preparation
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    fmt.Sprintf("Create %s Inheritance", c.schemaName),
		SectionTitle: fmt.Sprintf("Create a new %s Inheritance", c.schemaName),
		SidebarList:  sidebar,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: true,
			DeletePage: false,
			PolicyMode: "policy",
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: fmt.Sprintf("%s/create-inheritance", c.adminHomeUrl),
				FormMethod: "post",
			},
			FormFields: createForm,
		},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (c adminAuthPolicyController) DeleteInheritance(w http.ResponseWriter, r *http.Request) {
	// Grab params from URL
	inheritSlug := chi.URLParam(r, "inherit-slug")
	// Split by comma to separate the two params
	inheritArray := strings.Split(inheritSlug, ",")
	// Assign individually
	role := inheritArray[0]
	inherits := inheritArray[1]

	// If form is being submitted (method = POST)
	if r.Method == "POST" {
		// Delete user
		err := c.service.DeleteInheritance(models.G2Record{Role: role, InheritsFrom: inherits})
		if err != nil {
			http.Error(w, fmt.Sprintf("Error deleting %s", c.schemaName), http.StatusInternalServerError)
			return
		}
		// Redirect to success page
		http.Redirect(w, r, fmt.Sprintf("%s/delete-inheritance/success", c.adminHomeUrl), http.StatusSeeOther)
		return
	}

	// GET request
	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    fmt.Sprintf("Delete %s", c.schemaName),
		SectionTitle: fmt.Sprintf("Are you sure you wish to delete: %s?", fmt.Sprintf("%s inherits from %s", role, inherits)),
		SidebarList:  sidebar,
		SchemaHome:   c.adminHomeUrl,
		PageType: PageType{
			EditPage:   false,
			ReadPage:   false,
			CreatePage: false,
			DeletePage: true,
			PolicyMode: "policy",
		},
		FormData: FormData{
			FormDetails: FormDetails{
				FormAction: fmt.Sprintf("%s/delete-inheritance/%s", c.adminHomeUrl, inheritSlug),
				FormMethod: "post",
			},
			FormFields: []FormField{},
		},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err := app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func (c adminAuthPolicyController) Edit(w http.ResponseWriter, r *http.Request) {
	// Grab slug from URL
	policySlug := chi.URLParam(r, "id")
	// Unslug
	policyUnslug := UnslugifyResourceName(policySlug)
	// Detect request method
	method := r.Method

	// If form is being submitted (method = POST)
	if method == "POST" || method == "DELETE" {
		// Extract from request body, json policy
		pol := &models.PolicyRule{}

		// Decode request body as JSON and store in login
		err := json.NewDecoder(r.Body).Decode(&pol)
		if err != nil {
			http.Error(w, "Invalid policy", http.StatusBadRequest)
			return
		}

		// Validate the incoming DTO
		pass, _ := helpers.GoValidateStruct(pol)

		// If passes
		if pass {
			if method == "POST" {
				// Create policy
				err = c.service.Create(*pol)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error creating %s", c.schemaName), http.StatusInternalServerError)
					return
				}
			} else if method == "DELETE" {
				// Delete policy
				err = c.service.Delete(*pol)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error deleting %s", c.schemaName), http.StatusInternalServerError)
					return
				}
			}
			// Redirect or render a success message
			http.Redirect(w, r, fmt.Sprintf("%s/%s", c.adminHomeUrl, policySlug), http.StatusSeeOther)
			return
		}

	}

	// If not POST, ie. GET
	// Find all policies
	found, err := c.service.FindAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("%s not found", c.schemaName), http.StatusNotFound)
		return
	}

	// Filter by search query
	groupsSlice := searchPoliciesForExactResouceMatch(found, policyUnslug)
	// Prepare policies for rendering
	policies := convertMapToPolicyRule(groupsSlice)

	// Init new role selector values
	rolesCurrentlyInPolicy := c.formSelectors.RoleSelection()
	// Remove roles that are already in the policy
	rolesCurrentlyInPolicy = buildOnlyMissingRoleSelector(policies, rolesCurrentlyInPolicy)

	// Data to be injected into template
	data := PageRenderData{
		PageTitle:    fmt.Sprintf("Edit %s: %s", c.schemaName, policyUnslug),
		SectionTitle: fmt.Sprintf("Edit %s: %s", c.schemaName, policyUnslug),
		SidebarList:  sidebar,
		PageType: PageType{
			EditPage:   true,
			ReadPage:   false,
			CreatePage: false,
			DeletePage: false,
			PolicyMode: "policy",
		},
		PolicySection: PolicySection{
			FocusedPolicies: policies,
			PolicyResource:  policyUnslug,
			Selectors: PolicyEditSelectors{
				RoleSelection:   rolesCurrentlyInPolicy,
				ActionSelection: c.formSelectors.ActionSelection()},
		},
		FormData: FormData{
			FormDetails: FormDetails{},
			FormFields:  []FormField{},
		},
		HeaderSection: header,
	}

	// Execute the template with data and write to response
	err = app.AdminTemplates.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// Success handlers
func (c adminAuthPolicyController) CreateSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Create %s", c.schemaName), fmt.Sprintf("%s Created Successfully!", c.schemaName))
}
func (c adminAuthPolicyController) CreateRoleSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Create %s Role", c.schemaName), fmt.Sprintf("%s Role Created Successfully!", c.schemaName))
}
func (c adminAuthPolicyController) CreateInheritanceSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Create %s Role Inheritance", c.schemaName), fmt.Sprintf("%s Role Inheritance Created Successfully!", c.schemaName))
}
func (c adminAuthPolicyController) DeleteInheritanceSuccess(w http.ResponseWriter, r *http.Request) {
	// Serve admin success page
	serveAdminSuccess(w, fmt.Sprintf("Delete %s Inheritance Inheritance", c.schemaName), fmt.Sprintf("%s Inheritance Inheritance Created Successfully!", c.schemaName))
}

// Form generation
// Used to build Create form
func (c adminAuthPolicyController) generateCreateForm() []FormField {
	return []FormField{
		{DbLabel: "Resource", Label: "Resource", Name: "resource", Placeholder: "eg. '/api/posts'", Value: "", Type: "text", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "Role", Label: "First Role", Name: "role", Placeholder: "", Value: "", Type: "select", Required: true, Disabled: false, Errors: []ErrorMessage{}, Selectors: c.formSelectors.RoleSelection()},
		{DbLabel: "Action", Label: "Action", Name: "action", Placeholder: "", Value: "", Type: "select", Required: false, Disabled: false, Errors: []ErrorMessage{}, Selectors: c.formSelectors.ActionSelection()},
	}
}
func (c adminAuthPolicyController) generateCreateRoleForm() []FormField {
	return []FormField{
		{DbLabel: "Role", Label: "New Role Name", Name: "role", Placeholder: "eg. 'Moderator'", Value: "", Type: "text", Required: true, Disabled: false, Errors: []ErrorMessage{}},
		{DbLabel: "User", Label: "First Member", Name: "user", Placeholder: "", Value: "", Type: "select", Required: true, Disabled: false, Errors: []ErrorMessage{}, Selectors: c.formSelectors.UserSelection()},
	}
}
func (c adminAuthPolicyController) generateCreateInheritanceForm() []FormField {
	return []FormField{
		{DbLabel: "Role", Label: "Role", Name: "role", Placeholder: "", Value: "", Type: "select", Required: true, Disabled: false, Errors: []ErrorMessage{}, Selectors: c.formSelectors.RoleSelection()},
		{DbLabel: "InheritsFrom", Label: "Inherits from (role)", Name: "inherits_from", Placeholder: "", Value: "", Type: "select", Required: true, Disabled: false, Errors: []ErrorMessage{}, Selectors: c.formSelectors.RoleSelection()},
	}
}

// Extract forms
// Used to extract form submission from request and build into service-ready format
func (c adminAuthPolicyController) extractCreateFormSubmission(r *http.Request) (models.PolicyRule, error) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		return models.PolicyRule{}, errors.New("Error parsing form")
	}

	// Build struct for validation
	toValidate := models.PolicyRule{
		Role:     r.FormValue("role"),
		Resource: r.FormValue("resource"),
		Action:   r.FormValue("action"),
	}

	return toValidate, nil
}
func (c adminAuthPolicyController) extractCreateRoleFormSubmission(r *http.Request) (models.CasbinRoleAssignment, error) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		return models.CasbinRoleAssignment{}, errors.New("Error parsing form")
	}

	// Build struct for validation
	toValidate := models.CasbinRoleAssignment{
		Role:   r.FormValue("role"),
		UserId: r.FormValue("user"),
	}

	return toValidate, nil
}
func (c adminAuthPolicyController) extractCreateInheritanceFormSubmission(r *http.Request) (models.G2Record, error) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		return models.G2Record{}, errors.New("Error parsing form")
	}

	// Build struct for validation
	toValidate := models.G2Record{
		Role:         r.FormValue("role"),
		InheritsFrom: r.FormValue("inherits_from"),
	}

	return toValidate, nil
}

// Basic helper functions
// Used to extract form submission from request and build into map[string]string (Used in populateValuesWithForm)
func (c adminAuthPolicyController) extractFormFromRequest(r *http.Request) (map[string]string, error) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing form: %s", err.Error()))
	}
	// Map of user fields
	fieldMap := map[string]string{
		"Name":     r.FormValue("name"),
		"Username": r.FormValue("username"),
		"Email":    r.FormValue("email"),
		"Role":     r.FormValue("role"),
		"Verified": r.FormValue("verified"),
	}
	return fieldMap, nil
}

// Used to build standardize controller fields for admin panel sidebar generation
func (c adminAuthPolicyController) ObtainFields() BasicAdminController {
	return basicAdminController{
		AdminHomeUrl:     c.adminHomeUrl,
		SchemaName:       c.schemaName,
		PluralSchemaName: c.pluralSchemaName,
	}
}

// Search helpers
// Searches a list of policies for a given resource based on search term
func searchPoliciesByResource(maps []map[string]interface{}, searchTerm string) []map[string]interface{} {
	var result []map[string]interface{}

	// Iterate through map of policies
	for _, m := range maps {
		// Grab resource
		resource, ok := m["resource"].(string)
		// If success and resource contains search term
		if ok && containsString(resource, searchTerm) {
			result = append(result, m)
		}
	}

	return result
}

func searchMapUsingKeys(maps []map[string]string, mapKeysToSearch []string, searchTerm string) []map[string]string {
	var result []map[string]string
	// Init to record if already added to results
	addedToResult := false

	// Iterate through map of policies
	for _, m := range maps {
		// Reset added to result
		addedToResult = false
		// Iterate through list of keys to search for term
		for _, keyToSearch := range mapKeysToSearch {
			// Grab value
			value, ok := m[keyToSearch]
			// If success, and the record hasn't been added already and value contains search term
			if ok && containsString(value, searchTerm) && !addedToResult {
				// Append
				result = append(result, m)
				// Set added to true
				addedToResult = true
			}
		}
	}

	return result
}

// Searches a list of policies for a given resource based on search term
func searchPoliciesForExactResouceMatch(maps []map[string]interface{}, searchTerm string) []map[string]interface{} {
	var result []map[string]interface{}

	// Iterate through map of policies
	for _, m := range maps {
		// Grab resource
		resource, ok := m["resource"].(string)
		// If success and resource contains search term
		if ok && resource == searchTerm {
			result = append(result, m)
		}
	}

	return result
}

// Convert the map received from the service to a slice of models.PolicyRule
func convertMapToPolicyRule(m []map[string]interface{}) []PolicyEditDataRow {
	var policies []PolicyEditDataRow
	// Iterate through map of policies
	for _, policy := range m {
		var actions []PolicyActionCell

		// Iterate through actions
		for _, action := range possibleActions {
			// Create policy action cell
			actionToAdd := PolicyActionCell{
				Action: action,
				// Make false as default
				Granted: false,
			}

			// Check if array contains a string
			if arrayContainsString(policy["action"].([]string), action) {
				actionToAdd.Granted = true
			}

			// Append to actions
			actions = append(actions, actionToAdd)
		}

		// Build policy edit row
		policyToAdd := PolicyEditDataRow{
			Role:     policy["role"].(string),
			Resource: policy["resource"].(string),
			Actions:  actions,
		}

		// Append to policies
		policies = append(policies, policyToAdd)
	}
	return policies
}

// Takes a slice of PolicyEditDataRow and role selector and returns a slice of role selector with only missing roles
func buildOnlyMissingRoleSelector(policies []PolicyEditDataRow, roleSelector []FormFieldSelector) []FormFieldSelector {
	for _, p := range policies {
		// Iterate through roleSelector
		for i, role := range roleSelector {
			// If the role matches
			if role.Value == p.Role {
				// Remove from slice
				roleSelector = append(roleSelector[:i], roleSelector[i+1:]...)
				break
			}
		}
	}
	return roleSelector
}
