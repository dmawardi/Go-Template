package modules

import (
	modulecontrollers "github.com/dmawardi/Go-Template/internal/controller/moduleControllers"
	webapi "github.com/dmawardi/Go-Template/internal/helpers/webApi"
	"github.com/dmawardi/Go-Template/internal/models"
	modulerepositories "github.com/dmawardi/Go-Template/internal/repository/module"
	moduleservices "github.com/dmawardi/Go-Template/internal/service/module"
)

// Define setup configurations (to use in setupModules within API setup function)
var ModulesToSetup = []models.EntityConfig{
	{
		// Used for module name in module map
		Name:          "Post",
		NewRepo:       webapi.NewRepository(modulerepositories.NewPostRepository),
		NewService:    webapi.NewService(moduleservices.NewPostService),
		NewController: webapi.NewController(modulecontrollers.NewPostController),
	},
	// ADD ADDITIONAL BASIC MODULES HERE
}
