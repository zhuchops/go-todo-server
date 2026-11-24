package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Task holds the schema definition for the Task entity.
type Task struct {
	ent.Schema
}

// Fields of the Task.
func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.Int("owner_id"),
		field.String("content"),
		field.Time("created_at").Default(time.Now()).Immutable(),
		field.Enum("status").
			Values("todo", "done").
			Default("todo"),
	}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("tasks").
			Unique().
			Required().
			Field("owner_id"),
	}
}
