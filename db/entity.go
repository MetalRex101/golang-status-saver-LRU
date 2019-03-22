package db

import (
	"fmt"
	"gitlab.com/artilligence/http-db-saver/domain"
	"math/rand"
	"time"
)

type entityRepo struct {}

func NewEntityRepo () domain.EntityRepository {
	return &entityRepo{}
}

func (r *entityRepo) Insert(entity *domain.Entity) error {
	fmt.Println(fmt.Printf("insert %+v", entity))

	return nil
}

func (r *entityRepo) Get(id int64) *domain.Entity {
	rand.Seed(time.Now().UnixNano())

	return &domain.Entity{
		Status: true,
		ID: id,
	}
}