package models

type CasbinRule struct {
	PType string `json:"ptype" gorm:"size:100;uniqueIndex:unique_index" valid:"required,in(p|g|g2)"`
	V0    string `json:"v0" gorm:"size:100;uniqueIndex:unique_index" valid:"required"`
	V1    string `json:"v1" gorm:"size:100;uniqueIndex:unique_index" valid:"required"`
	V2    string `json:"v2" gorm:"size:100;uniqueIndex:unique_index" valid:"in(read|create|update|delete)"`
}

type UpdateCasbinRule struct {
	OldPolicy CasbinRule `json:"old_policy" valid:"required"`
	NewPolicy CasbinRule `json:"new_policy" valid:"required"`
}

type CasbinRoleAssignment struct {
	UserId string `json:"user_id"  valid:"required"`
	Role   string `json:"role"  valid:"required"`
}
