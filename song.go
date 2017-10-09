package happy

import (
	"fmt"
	"time"
)

type Song struct {
	ID          uint64    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Artists     []string  `json:"artists"`
	ArtistLabel string    `json:"artist_label" db:"artist_label"`
	Link        string    `json:"link" db:"link"`
	DateCreated time.Time `json:"date_created" db:"date_created"`
}

func (db *PGDB) CreateSong(link string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	tx.QueryRow(`INSERT INTO song(name, link, artist_label, date_created) 
		VALUES ($1, $2, $3, $4)`,
		link, link, link, TimeNow())
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (db *PGDB) GetSongList(perBatch int) ([]*Song, error) {
	var songList []*Song

	queryLimit := fmt.Sprintf(` LIMIT %d`, perBatch)
	err := db.Select(&songList, fmt.Sprintf(`SELECT * FROM song ORDER BY id %s`, queryLimit))
	if err != nil {
		return nil, err
	}

	return songList, nil
}
