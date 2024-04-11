package backend

import (
	"log"
	"propertyManagement/model"
	"time"
)

func (backend *PostgresBackend) InsertDiscussion(username, subject, content string) (*model.Discussion, error) {
	tx, _ := backend.db.Begin()
	defer tx.Rollback()
	
	query := "INSERT INTO discussions (username, subject, content, last_update_time) VALUES ($1, $2, $3, $4) RETURNING id"

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	var id int
	err := tx.QueryRow(query, username, subject, content, formattedTime).Scan(&id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return nil, err
	}

	return backend.SelectDiscussionById(int(id))
}

func (backend *PostgresBackend) SelectAllDiscussions() ([]model.Discussion, error) {
	rows, err := backend.db.Query("SELECT id, username, subject, content, TO_CHAR(last_update_time, 'YYYY-MM-DD HH24:MI:SS') FROM discussions ORDER BY last_update_time DESC")
    if err != nil {
		log.Println(err)
        return nil, err
    }
	defer rows.Close()
	
	discussions := []model.Discussion{}
	for rows.Next() {
		var discussion model.Discussion
		err := rows.Scan(&discussion.Id,
						 &discussion.Username,
						 &discussion.Subject,
						 &discussion.Content,
						 &discussion.LastUpdateTime)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		discussions = append(discussions, discussion)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return discussions, nil
}

func (backend *PostgresBackend) SelectAllDiscussionsByUsername(username string) ([]model.Discussion, error) {
	query := "SELECT id, username, subject, content, TO_CHAR(last_update_time, 'YYYY-MM-DD HH24:MI:SS') FROM discussions WHERE username = $1 ORDER BY last_update_time DESC"
	rows, err := backend.db.Query(query, username)
    if err != nil {
		log.Println(err)
        return nil, err
    }
	defer rows.Close()
	
	discussions := []model.Discussion{}
	for rows.Next() {
		var discussion model.Discussion
		err := rows.Scan(&discussion.Id,
						 &discussion.Username,
						 &discussion.Subject,
						 &discussion.Content,
						 &discussion.LastUpdateTime)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		discussions = append(discussions, discussion)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return discussions, nil
}

func (backend *PostgresBackend) SelectDiscussionById(id int) (*model.Discussion, error) {
	var discussion model.Discussion
	err := backend.db.QueryRow("SELECT id, username, subject, content, TO_CHAR(last_update_time, 'YYYY-MM-DD HH24:MI:SS') FROM discussions WHERE id = $1", id).
		   Scan(&discussion.Id,
				&discussion.Username,
				&discussion.Subject,
				&discussion.Content,
				&discussion.LastUpdateTime)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &discussion, nil
}

func (backend *PostgresBackend) DeleteDiscussion(username string, id int) error {
	tx, _ := backend.db.Begin()
	defer tx.Rollback()
	
	query := "DELETE FROM discussions WHERE username = $1 AND id = $2"

	_, err := tx.Exec(query, username, id)
	if err != nil {
		log.Println(err)
		return err
	}
	
	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (backend *PostgresBackend) DiscussionExists(id int) (bool, error) {
	var count int
	err := backend.db.QueryRow("SELECT COUNT(*) FROM discussions WHERE id = $1", id).Scan(&count)
	if err != nil {
		log.Println(err)
        return false, err
    }
	return count > 0, nil
}
