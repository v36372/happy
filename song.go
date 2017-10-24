package happy

import (
	"errors"
)

type Song struct {
	ID        uint64 `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Link      string `json:"link" db:"link"`
	Provider  string `json:"provider" db:"provider"`
	Thumbnail string `json:"thumbnail" db:"thumbnail"`
}

func (db *PGDB) CreateSong(song *Song) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if song == nil {
		return errors.New("error when creating song, song is nil")
	}

	tx.QueryRow(`INSERT INTO song(name, link, provider, thumbnail) 
		VALUES ($1, $2, $3, $4)`,
		song.Name, song.Link, song.Provider, song.Thumbnail)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (db *PGDB) GetSongList() ([]*Song, error) {
	var songList []*Song

	err := db.Select(&songList, "SELECT * FROM song ORDER BY id")
	if err != nil {
		return nil, err
	}

	return songList, nil
}
