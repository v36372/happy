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

func (db *PGDB) CreateSong(songs []*Song) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if len(songs) == 0 {
		return errors.New("error when creating song, empty array")
	}

	for _, song := range songs {
		tx.Exec(`INSERT INTO song(name, link, provider, thumbnail) 
		VALUES ($1, $2, $3, $4)`,
			song.Name, song.Link, song.Provider, song.Thumbnail)
	}
	tx.Commit()

	return nil
}

func (db *PGDB) GetSongList() ([]*Song, error) {
	songList := []*Song{}

	err := db.Select(&songList, "SELECT * FROM song ORDER BY id")
	if err != nil {
		return nil, err
	}

	return songList, nil
}

func (db *PGDB) DeleteSong(songID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	tx.Exec(`DELETE FROM song WHERE id = $1`, songID)
	tx.Commit()

	return nil
}
