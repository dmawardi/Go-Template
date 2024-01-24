package service

import (
	"fmt"
	"sort"

	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
)

type AuthPolicyService interface {
	// Policies
	FindAll(searchQuery string) ([]map[string]interface{}, error)
	FindByResource(policyResource string) ([]map[string]interface{}, error)
	Create(policy models.PolicyRule) error
	Update(oldPolicy, newPolicy models.PolicyRule) error
	Delete(policy models.PolicyRule) error
	// Roles
	FindAllRoles() ([]string, error)
	AssignUserRole(userId, roleToApply string) (*bool, error)
	// Inheritance
	FindAllRoleInheritance() ([]map[string]string, error)
	CreateInheritance(inherit models.G2Record) error
	DeleteInheritance(inherit models.G2Record) error
	// Not for controller usage (used in auth)
	FindRoleByUserId(userId int) (string, error)
}

type authPolicyService struct {
	repo repository.AuthPolicyRepository
}

func NewAuthPolicyService(repo repository.AuthPolicyRepository) AuthPolicyService {
	return &authPolicyService{repo}
}

// Policies
//

func (s *authPolicyService) FindAll(searchQuery string) ([]map[string]interface{}, error) {
	data, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Transform data for easier consumption
	groupsSlice := transformDataToResponse(data)

	if searchQuery != "" {
		// Filter policies []map[string]interface{} by search query
		groupsSlice = helpers.SearchPoliciesByResource(groupsSlice, searchQuery)
	}

	// Sort by resource alphabetically
	sort.Slice(groupsSlice, func(i, j int) bool {
		// Give two items to compare to role resource alpha sorter
		return helpers.SortMapStringInterfaceAlphabetically(groupsSlice[i], groupsSlice[j], "resource")
	})

	return groupsSlice, nil
}
func (s *authPolicyService) FindByResource(policyResource string) ([]map[string]interface{}, error) {
	data, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Transform data for easier consumption
	organizedData := transformDataToResponse(data)
	// Filter by search query
	resourceMatchRecords := helpers.SearchPoliciesForExactResouceMatch(organizedData, policyResource)

	return resourceMatchRecords, nil
}
func (s *authPolicyService) Create(policy models.PolicyRule) error {
	casbinPolicy := models.CasbinRule{
		PType: "p",
		V0:    policy.Role,
		V1:    policy.Resource,
		V2:    policy.Action,
	}

	return s.repo.Create(casbinPolicy)
}
func (s *authPolicyService) Update(oldPolicy, newPolicy models.PolicyRule) error {
	oldCasbinPolicy := models.CasbinRule{
		PType: "p",
		V0:    oldPolicy.Role,
		V1:    oldPolicy.Resource,
		V2:    oldPolicy.Action,
	}
	newCasbinPolicy := models.CasbinRule{
		PType: "p",
		V0:    newPolicy.Role,
		V1:    newPolicy.Resource,
		V2:    newPolicy.Action,
	}
	return s.repo.Update(oldCasbinPolicy, newCasbinPolicy)
}
func (s *authPolicyService) Delete(policy models.PolicyRule) error {
	casbinPolicy := models.CasbinRule{
		PType: "p",
		V0:    policy.Role,
		V1:    policy.Resource,
		V2:    policy.Action,
	}
	return s.repo.Delete(casbinPolicy)
}

// Roles
//

func (s *authPolicyService) FindAllRoles() ([]string, error) {
	return s.repo.FindAllRoles()
}
func (s *authPolicyService) FindRoleByUserId(userId int) (string, error) {
	// Convert the userId to string then pass to repo
	return s.repo.FindRoleByUserId(fmt.Sprint(userId))
}

// Assigning a user a role will automatically create a new role with the user as the first member if not already found
func (s *authPolicyService) AssignUserRole(userId, roleToApply string) (*bool, error) {
	success, err := s.repo.AssignUserRole(userId, roleToApply)
	if err != nil {
		return nil, err
	}
	return success, nil
}

// Inheritance
//

func (s *authPolicyService) FindAllRoleInheritance() ([]map[string]string, error) {
	formattedG2Records := []map[string]string{}
	g2Records, err := s.repo.FindAllRoleInheritance()
	if err != nil {
		return nil, err
	}

	// Iterate through g2Records and format into array of G2Record struct
	for _, policy := range g2Records {
		record := map[string]string{
			"role":          policy[0],
			"inherits_from": policy[1],
		}
		formattedG2Records = append(formattedG2Records, record)
	}
	return formattedG2Records, nil
}
func (s *authPolicyService) CreateInheritance(inherit models.G2Record) error {
	return s.repo.CreateInheritance(inherit)
}
func (s *authPolicyService) DeleteInheritance(inherit models.G2Record) error {
	return s.repo.DeleteInheritance(inherit)
}

// Transform data from enforcer policies to User friendly response
func transformDataToResponse(data [][]string) []map[string]interface{} {
	// Response format
	response := make(map[string][]map[string]interface{})
	// Init policy dictionary for sorting
	policyDict := make(map[string]map[string]interface{})

	// Loop through data and build policy dictionary
	for _, item := range data {
		// Assign policy vars
		role, resource, action := item[0], item[1], item[2]
		key := role + resource

		// If key does not exist, create new entry
		if _, ok := policyDict[key]; !ok {
			policyDict[key] = map[string]interface{}{
				"role":     role,
				"resource": resource,
				"action":   []string{action},
			}

		} else {
			// Else, if record exists with resource, append action to action slice
			policyDict[key]["action"] = append(policyDict[key]["action"].([]string), action)
		}
	}

	// Loop through policyDict and append to response
	for _, policy := range policyDict {
		response["policies"] = append(response["policies"], policy)
	}

	return response["policies"]
}
