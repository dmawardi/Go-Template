package models

import "gorm.io/gorm"

// Used in API setup to standardize the array of setup configurations
type EntityConfig struct {
	Name          string
	NewRepo       func(*gorm.DB) interface{}
	NewService    func(interface{}) interface{}
	NewController func(interface{}) interface{}
}

// Wrappers used for API setup function to ensure different repositories are standardized
type RepoFactory func(*gorm.DB) interface{}
type ServiceFactory func(interface{}) interface{}
type ControllerFactory func(interface{}) interface{}

// Module set is used within SetupBasicModules to store the different modules
type ModuleSet = struct {
	Repo       interface{}
	Service    interface{}
	Controller interface{}
}
