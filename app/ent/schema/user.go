package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("first_name").
			NotEmpty().
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(20)",
			}),
		field.String("last_name").
			NotEmpty().
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(20)",
			}),
		field.String("email").
			NotEmpty().
			Unique().
			SchemaType(map[string]string{
				dialect.MySQL: "varchar(50)",
			}),
		field.Int("age").
			Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("cars", Car.Type),
		// Create an inverse-edge called "groups" of type `Group`
		// and reference it to the "users" edge (in Group schema)
		// explicitly using the `Ref` method.
		edge.From("group", Group.Type).
			Ref("users"),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
