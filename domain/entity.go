package domain

const TypeEntity EntityType = "entity"

type Entity struct {
	ID     int64 `json:"id"`
	Status bool  `json:"status"`
}

type EntityRepository interface {
	Insert(*Entity) error
	Get(id int64) *Entity
}
