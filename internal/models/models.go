package models

type Job struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Login
type Login struct {
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required"`
}

type ValidationError struct {
	Validation_errors map[string][]string `json:"validation_errors"`
}
type SchemaMetaData struct {
	Total_Records    int64 `json:"total_records"`    // Total number of records in the entire dataset
	Records_Per_Page int   `json:"records_per_page"` // Number of records displayed per page
	Current_Page     int   `json:"current_page"`     // Current page number
	Next_Page        *int  `json:"next_page"`        // Next page number (null if there is no next page)
	Prev_Page        *int  `json:"prev_page"`        // Previous page number (null if there is no previous page)
}
