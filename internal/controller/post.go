package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/go-chi/chi"
)

type PostController interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type postController struct {
	service service.PostService
}

func NewPostController(service service.PostService) PostController {
	return &postController{service}
}

// Used to init the query params for easy extraction in controller
// Returns: map[string]string{"age": "int", "name": "string", "active": "bool"}
func PostConditionQueryParams() map[string]string {
	return map[string]string{
		"title": "string",
		"body":  "string",
	}
}

// API/POSTS
// Find a list of posts
// @Summary      Find a list of posts
// @Description  Accepts limit, offset, and order params and returns list of posts
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        limit   query      int  true  "limit"
// @Param        offset   query      int  false  "offset"
// @Param        order   query      int  false  "order by"
// @Success      200 {object} []models.PaginatedPosts
// @Failure      400 {string} string "Can't find posts"
// @Failure      400 {string} string "Must include limit parameter with a max value of 50"
// @Router       /posts/{id} [get]
// @Security BearerToken
func (c postController) FindAll(w http.ResponseWriter, r *http.Request) {
	// Grab basic query params
	baseQueryParams, err := helpers.ExtractBasicFindAllQueryParams(r)
	if err != nil {
		http.Error(w, "Error extracting query params", http.StatusBadRequest)
		return
	}

	// Generate query params to extract
	queryParamsToExtract := PostConditionQueryParams()
	// Extract query params
	extractedConditionParams, err := helpers.ExtractSearchAndConditionParams(r, queryParamsToExtract)
	if err != nil {
		http.Error(w, "Error extracting query params", http.StatusBadRequest)
		return
	}

	// Check that limit is present as requirement
	if (baseQueryParams.Limit == 0) || (baseQueryParams.Limit > 50) {
		http.Error(w, "Must include limit parameter with a max value of 50", http.StatusBadRequest)
		return
	}

	// Query database for all users using query params
	posts, err := c.service.FindAll(baseQueryParams.Limit, baseQueryParams.Offset, baseQueryParams.Order, extractedConditionParams)
	if err != nil {
		http.Error(w, "Can't find posts", http.StatusBadRequest)
		return
	}
	err = helpers.WriteAsJSON(w, posts)
	if err != nil {
		http.Error(w, "Can't find posts", http.StatusBadRequest)
		fmt.Println("error writing users to response: ", err)
		return
	}
}

// Find a created post
// @Summary      Find post
// @Description  Find a post by ID
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Post ID"
// @Success      200 {object} db.Post
// @Failure      400 {string} string "Can't find post"
// @Router       /posts/{id} [get]
// @Security BearerToken
func (c postController) Find(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, err := strconv.Atoi(stringParameter)
	fmt.Println("id parameter from request: ", stringParameter)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	foundPost, err := c.service.FindById(idParameter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find post with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
	err = helpers.WriteAsJSON(w, foundPost)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't find post with ID: %v\n", idParameter), http.StatusBadRequest)
		return
	}
}

// Create a new post
// @Summary      Create Post
// @Description  Creates a new post
// @Tags         Post
// @Accept       json
// @Produce      plain
// @Param        post body models.CreatePost true "New Post Json"
// @Success      201 {string} string "Post creation successful!"
// @Failure      400 {string} string "Post creation failed."
// @Router       /posts [post]
func (c postController) Create(w http.ResponseWriter, r *http.Request) {
	// Init
	var post models.CreatePost
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&post)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Create post
	_, createErr := c.service.Create(&post)
	if createErr != nil {
		http.Error(w, "Post creation failed.", http.StatusBadRequest)
		return
	}

	// Set status to created
	w.WriteHeader(http.StatusCreated)
	// Send post success message in body
	w.Write([]byte("Post creation successful!"))
}

// Update a post (using URL parameter id)
// @Summary      Update Post
// @Description  Updates an existing post
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        post body models.UpdatePost true "Update Post Json"
// @Param        id   path      int  true  "User ID"
// @Success      200 {object} db.Post
// @Failure      400 {string} string "Failed post update"
// @Failure      403 {string} string "Authentication Token not detected"
// @Router       /posts/{id} [put]
// @Security BearerToken
func (c postController) Update(w http.ResponseWriter, r *http.Request) {
	// grab id parameter
	var post models.UpdatePost
	// Decode request body as JSON and store in login
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		fmt.Println("Decoding error: ", err)
	}

	// Validate the incoming DTO
	pass, valErrors := helpers.GoValidateStruct(&post)
	// If failure detected
	if !pass {
		// Write bad request header
		w.WriteHeader(http.StatusBadRequest)
		// Write validation errors to JSON
		helpers.WriteAsJSON(w, valErrors)
		return
	}
	// else, validation passes and allow through

	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Update post
	updatedPost, createErr := c.service.Update(idParameter, &post)
	if createErr != nil {
		http.Error(w, fmt.Sprintf("Failed post update: %s", createErr), http.StatusBadRequest)
		return
	}
	// Write post to output
	err = helpers.WriteAsJSON(w, updatedPost)
	if err != nil {
		fmt.Printf("Error encountered when writing to JSON. Err: %s", err)
	}
}

// Delete post (using URL parameter id)
// @Summary      Delete Post
// @Description  Deletes an existing post
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Post ID"
// @Success      200 {string} string "Deletion successful!"
// @Failure      400 {string} string "Failed post deletion"
// @Router       /posts/{id} [delete]
// @Security BearerToken
func (c postController) Delete(w http.ResponseWriter, r *http.Request) {
	// Grab URL parameter
	stringParameter := chi.URLParam(r, "id")
	// Convert to int
	idParameter, _ := strconv.Atoi(stringParameter)

	// Attampt to delete post using id
	err := c.service.Delete(idParameter)

	// If error detected
	if err != nil {
		http.Error(w, "Failed post deletion", http.StatusBadRequest)
		return
	}
	// Else write success
	w.Write([]byte("Deletion successful!"))
}
