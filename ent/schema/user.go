package schema

import (
	"errors"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/asaskevich/govalidator"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Default("unknown"),
		field.String("email").Unique().Validate(func(s string) error {
			if s == "" {
				return errors.New("can't create with empty email")
			}
			if !govalidator.IsEmail(s) {
				return errors.New("email not valid")
			}

			return nil
		}),
		field.String("username"),
		field.String("role").Default("user"),
		field.String("password").Sensitive(),
		field.Time("created_at").
			Default(time.Now).Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	// Build relationship to cars
	return []ent.Edge{
		edge.To("cars", Car.Type),
		// Create an inverse-edge called "groups" of type `Group`
		// and reference it to the "users" edge (in Group schema)
		// explicitly using the `Ref` method.
		edge.From("groups", Group.Type).
			Ref("users"),
	}
}
