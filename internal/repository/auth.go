package repository

// import (
//     "github.com/casbin/casbin/v2"
//     "gorm.io/gorm"
// )

// // CasbinPolicyRepository represents a repository for Casbin policies.
// type CasbinPolicyRepository interface {
//     FindAll() ([]casbin.Policy, error)
//     FindByUserID(userID string) ([]casbin.Policy, error)
//     Create(policy casbin.Policy) error
//     Update(oldPolicy, newPolicy casbin.Policy) error
//     Delete(policy casbin.Policy) error
// }

// // GormCasbinPolicyRepository is a GORM implementation of CasbinPolicyRepository.
// type GormCasbinPolicyRepository struct {
//     db *gorm.DB
// }

// // NewGormCasbinPolicyRepository creates a new instance of GormCasbinPolicyRepository.
// func NewGormCasbinPolicyRepository(db *gorm.DB) *GormCasbinPolicyRepository {
//     return &GormCasbinPolicyRepository{
//         db: db,
//     }
// }

// // FindAll returns all Casbin policies.
// func (r *GormCasbinPolicyRepository) FindAll() ([]casbin.Policy, error) {
//     var policies []casbin.Policy
//     result := r.db.Model(&casbin.CasbinRule{}).Find(&policies)
//     if result.Error != nil {
//         return nil, result.Error
//     }
//     return policies, nil
// }

// // FindByUserID returns policies associated with a specific user ID.
// func (r *GormCasbinPolicyRepository) FindByUserID(userID string) ([]casbin.Policy, error) {
//     var policies []casbin.Policy
//     result := r.db.Model(&casbin.CasbinRule{}).Where("v0 = ?", userID).Find(&policies)
//     if result.Error != nil {
//         return nil, result.Error
//     }
//     return policies, nil
// }

// // Create adds a new policy.
// func (r *GormCasbinPolicyRepository) Create(policy casbin.Policy) error {
//     casbinRule := &casbin.CasbinRule{
//         Ptype: policy[0],
//         V0:    policy[1],
//         V1:    policy[2],
//         V2:    policy[3],
//         V3:    policy[4],
//         V4:    policy[5],
//         V5:    policy[6],
//     }
//     result := r.db.Create(casbinRule)
//     return result.Error
// }

// // Update updates an existing policy.
// func (r *GormCasbinPolicyRepository) Update(oldPolicy, newPolicy casbin.Policy) error {
//     casbinRule := &casbin.CasbinRule{
//         Ptype: oldPolicy[0],
//         V0:    oldPolicy[1],
//         V1:    oldPolicy[2],
//         V2:    oldPolicy[3],
//         V3:    oldPolicy[4],
//         V4:    oldPolicy[5],
//         V5:    oldPolicy[6],
//     }
//     result := r.db.Model(casbinRule).
//         Updates(map[string]interface{}{
//             "ptype": newPolicy[0],
//             "v0":    newPolicy[1],
//             "v1":    newPolicy[2],
//             "v2":    newPolicy[3],
//             "v3":    newPolicy[4],
//             "v4":    newPolicy[5],
//             "v5":    newPolicy[6],
//         })
//     return result.Error
// }

// // Delete removes an existing policy.
// func (r *GormCasbinPolicyRepository) Delete(policy casbin.Policy) error {
//     casbinRule := &casbin.CasbinRule{
//         Ptype: policy[0],
//         V0:    policy[1],
//         V1:    policy[2],
//         V2:    policy[3],
//         V3:    policy[4],
//         V4:    policy[5],
//         V5:    policy[6],
//     }
//     result := r.db.Delete(casbinRule)
//     return result.Error
// }
