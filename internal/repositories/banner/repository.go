package banner

import (
	"fmt"
	"test-work/internal/storages/postgres"
)

type Repository struct {
	Storage *postgres.Storage
}

func NewRepository(storage *postgres.Storage) *Repository {
	return &Repository{
		Storage: storage,
	}
}

func (r *Repository) AddStats(bannerID uint64, timestamp string, count uint64) error {
	err := r.Storage.Db.Exec(`
			INSERT INTO stats (count, banner_id, timestamp)
			VALUES (?, ?, ?)
			ON CONFLICT (banner_id, timestamp)
			DO UPDATE SET count = EXCLUDED.count;`, count, bannerID, timestamp).Error

	if err != nil {
		fmt.Println("Error inserting stats: ", err)
		return err
	}
	return nil
}

func (r *Repository) GetStatsSum(bannerID uint64, tsFrom string, tsTo string) (uint64, error) {
	var sum uint64
	err := r.Storage.Db.Raw(`
			select sum(count) as ct 
			from stats 
			where banner_id = ? and timestamp between ? and ?`,
		bannerID, tsFrom, tsTo).Scan(&sum).Error

	if err != nil {
		return 0, err
	}
	return sum, nil
}
