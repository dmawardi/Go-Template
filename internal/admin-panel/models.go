package adminpanel

// Used for bulk delete on find all pages
type BulkDeleteRequest struct {
	SelectedItems []string `json:"selected_items"`
}
