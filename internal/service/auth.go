package service

import (
	"github.com/dmawardi/Go-Template/internal/models"
	"github.com/dmawardi/Go-Template/internal/repository"
)

type AuthPolicyService interface {
	FindAll() ([]map[string]interface{}, error)
	FindAllRoles() ([]string, error)
	AssignUserRole(userId, roleToApply string) (*bool, error)
	Create(policy models.CasbinRule) error
	Update(oldPolicy, newPolicy models.CasbinRule) error
	Delete(policy models.CasbinRule) error
}

type authPolicyService struct {
	repo repository.AuthPolicyRepository
}

func NewAuthPolicyService(repo repository.AuthPolicyRepository) AuthPolicyService {
	return &authPolicyService{repo}
}

// // FindAll returns all Casbin policies.
func (s *authPolicyService) FindAll() ([]map[string]interface{}, error) {
	data, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	// Transform data
	organizedData := transformDataToResponse(data)

	return organizedData["policies"], nil
}

func (s *authPolicyService) AssignUserRole(userId, roleToApply string) (*bool, error) {
	success, err := s.repo.AssignUserRole(userId, roleToApply)
	if err != nil {
		return nil, err
	}
	return success, nil
}

func (s *authPolicyService) FindAllRoles() ([]string, error) {
	return s.repo.FindAllRoles()
}
func (s *authPolicyService) Create(policy models.CasbinRule) error {
	return s.repo.Create(policy)
}
func (s *authPolicyService) Update(oldPolicy, newPolicy models.CasbinRule) error {
	return s.repo.Update(oldPolicy, newPolicy)
}
func (s *authPolicyService) Delete(policy models.CasbinRule) error {
	return s.repo.Delete(policy)
}

// Transform data from enforcer policies to User friendly response
func transformDataToResponse(data [][]string) map[string][]map[string]interface{} {
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

	return response
}