package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dmawardi/Go-Template/ent"
	"github.com/dmawardi/Go-Template/ent/car"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/services"
)

// Init state variable
var app *config.AppConfig

// Function called in main.go to connect app state to current file
func SetStateInHandlers(a *config.AppConfig) {
	app = a
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
		return
	}

	var login models.Login
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}
	fmt.Printf("JSON Received: %+v\n", login)

	// Check if user exists in db
	foundUser, err := services.FindUserByEmail(app.Ctx, app.DbClient, login.Email)
	fmt.Println("founduser: ", foundUser)
	// If user found
	if err == nil {
		fmt.Println("User logging in: ", foundUser)

		// Compare stored (hashed) password with input password
		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(login.Password))

		// If match found (no errors)
		if err == nil {
			// Set login status to true
			tokenString, err := auth.GenerateJWT(foundUser.Username, foundUser.Email, foundUser.Role)
			// helpers.WriteAsJSON(w, )
			if err != nil {
				fmt.Println("Failed to create JWT")
			}
			// Build login response
			var loginResponse = LoginResponse{Token: tokenString}
			// Send to user in body
			helpers.WriteAsJSON(w, loginResponse)
			return
			// else if user password doesn't match
		} else {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
		// w.Write([]byte("Login successful!"))
	}
}

// Login URL check
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome!"))
}

// Detail to display a user's profile details
func UserDetails(w http.ResponseWriter, r *http.Request) {

}

// Users

// Handler to create a new user
func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	var user models.CreateUser
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}
	fmt.Printf("JSON Received: %+v\n", user)

	// Create user
	createdUser, createErr := services.CreateUser(app.Ctx, app.DbClient, &user)
	if createErr != nil {
		http.Error(w, "Failed user creation", http.StatusBadRequest)
		return
	}
	// Write user to output
	err = helpers.WriteAsJSON(w, createdUser)
	if err != nil {
		fmt.Println(err)
	}
}

// Handler to update a user (using URL parameter id)
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// grab id parameter
	var user models.UpdateUser
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)
	// Set within request body data object
	user.Id = idParameter

	fmt.Printf("JSON Received: %+v\n", user)

	// Update user
	updatedUser, createErr := services.UpdateUser(app.Ctx, app.DbClient, &user)
	if createErr != nil {
		http.Error(w, "Failed user creation", http.StatusBadRequest)
		return
	}
	// Write user to output
	err = helpers.WriteAsJSON(w, updatedUser)
	fmt.Println(err)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Attampt to delete user using id
	err := services.DeleteUser(app.Ctx, app.DbClient, idParameter)

	// If error detected
	if err != nil {
		http.Error(w, "Failed user deletion", http.StatusBadRequest)
		return
	}
	// Else write success
	w.Write([]byte("Deletion successful!"))
	return
}

// func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
// 	u, err := client.User.
// 		Query().
// 		Where(user.Name("a8m")).
// 		// `Only` fails if no user found,
// 		// or more than 1 user returned.
// 		Only(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed querying user: %w", err)
// 	}
// 	log.Println("user returned: ", u)
// 	return u, nil
// }

// Cars
func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
	// Create a new car with model "Tesla".

	tesla, err := client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", tesla)

	// Create a new car with model "Ford".
	ford, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", ford)

	// Create a new user, and add it the 2 cars.
	a8m, err := client.User.
		Create().
		SetName("a8m").
		AddCars(tesla, ford).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", a8m)
	return a8m, nil
}

func QueryCars(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}
	log.Println("returned cars:", cars)

	// What about filtering specific cars.
	ford, err := a8m.QueryCars().
		Where(car.Model("Ford")).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}
	log.Println(ford)
	return nil
}

func QueryCarUsers(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}
	// Query the inverse edge.
	for _, c := range cars {
		owner, err := c.QueryOwner().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying car %q owner: %w", c.Model, err)
		}
		log.Printf("car %q owner: %q\n", c.Model, owner.Name)
	}
	return nil
}

// Jobs
func GetJobs(w http.ResponseWriter, r *http.Request) {
	var jobs []models.Job

	jobs = append(jobs, models.Job{ID: 1, Name: "Accounting"})
	jobs = append(jobs, models.Job{ID: 2, Name: "Programming"})

	// Set header
	w.Header().Set("Content-Type", "application/json")

	// Build new JSON encoder to write to, then write jobs data
	json.NewEncoder(w).Encode(jobs)
}

// Create sample relationship graph
func CreateGraph(ctx context.Context, client *ent.Client) error {
	// First, create the users.
	a8m, err := client.User.
		Create().
		SetName("Ariel").
		Save(ctx)
	if err != nil {
		return err
	}
	neta, err := client.User.
		Create().
		SetName("Neta").
		Save(ctx)
	if err != nil {
		return err
	}
	// Then, create the cars, and attach them to the users created above.
	err = client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		// Attach this car to Ariel.
		SetOwner(a8m).
		Exec(ctx)
	if err != nil {
		return err
	}
	err = client.Car.
		Create().
		SetModel("Mazda").
		SetRegisteredAt(time.Now()).
		// Attach this car to Ariel.
		SetOwner(a8m).
		Exec(ctx)
	if err != nil {
		return err
	}
	err = client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		// Attach this graph to Neta.
		SetOwner(neta).
		Exec(ctx)
	if err != nil {
		return err
	}
	// Create the groups, and add their users in the creation.
	err = client.Group.
		Create().
		SetName("GitLab").
		AddUsers(neta, a8m).
		Exec(ctx)
	if err != nil {
		return err
	}
	err = client.Group.
		Create().
		SetName("GitHub").
		AddUsers(a8m).
		Exec(ctx)
	if err != nil {
		return err
	}
	log.Println("The graph was created successfully")
	return nil
}
