package storage

import (
	"github.com/SaidovZohid/hotel-project/storage/postgres"
	"github.com/SaidovZohid/hotel-project/storage/repo"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	User() repo.UserStorageI
}

type StoragePg struct {
	userRepo repo.UserStorageI
}

func NewStorage(db *sqlx.DB) StorageI {
	return &StoragePg{
		userRepo: postgres.NewUserStorage(db),
	}
}

func (s *StoragePg) User() repo.UserStorageI {
	return s.userRepo
}