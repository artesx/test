package banner

import "test-work/internal/storages/postgres"

type Repository struct {
	Storage *postgres.Storage
}

func NewRepository(storage *postgres.Storage) *Repository {
	return &Repository{
		Storage: storage,
	}
}
