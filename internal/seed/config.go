package seed

import "gorm.io/gorm"

func factoryRegister(db *gorm.DB) []FactoryRegistration {
	// Define the items to seed
	return []FactoryRegistration{
		// Add the list of factories here, with the name of the factory and the factory itself
		// This will be accessible in the factoryMap eg. factoryMap["Name"]
		// UserFactory
		{
			Factory: NewUserFactory(db),
			Name:    "User",
		},
	}
}

type FactoryRegistration struct {
	Factory BasicFactory
	Name    string
}
