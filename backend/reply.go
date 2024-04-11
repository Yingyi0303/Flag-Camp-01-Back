package backend

import (
	"log"
	"propertyManagement/model"
	"time"
)

func (backend *PostgresBackend) InsertReply(discussionId int, username, content string) (*model.Reply, error) {
	tx, _ := backend.db.Begin()
	defer tx.Rollback()
	
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	_, err := tx.Exec("UPDATE discussions SET last_update_time = $1 WHERE id = $2", formattedTime, discussionId)
	if (err != nil) {
		log.Println(err)
		return nil, err
	}

	query := "INSERT INTO replies (username, discussion_id, content, reply_time) VALUES ($1, $2, $3, $4) RETURNING id"
	var id int
	err = tx.QueryRow(query, username, discussionId, content, formattedTime).Scan(&id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return nil, err
	}

	return backend.SelectReplyById(int(id))
}

func (backend *PostgresBackend) SelectRepliesByDiscussionId(id int) ([]model.Reply, error) {
	rows, err := backend.db.Query("SELECT id, username, discussion_id, content, TO_CHAR(reply_time, 'YYYY-MM-DD HH24:MI:SS') from replies WHERE discussion_id = $1 ORDER BY reply_time DESC", id)
	if err != nil {
		log.Println(err)
        return nil, err
    }
	defer rows.Close()

	replies := []model.Reply{}
	for rows.Next() {
		var reply model.Reply
		err := rows.Scan(&reply.Id,
						 &reply.Username,
						 &reply.DiscussionId,
						 &reply.Content,
						 &reply.ReplyTime)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		replies = append(replies, reply)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return replies, nil
}

func (backend *PostgresBackend) SelectAllRepliesByUsername(username string) ([]model.Reply, error) {
	query := "SELECT id, username, discussion_id, content, TO_CHAR(reply_time, 'YYYY-MM-DD HH24:MI:SS') FROM replies WHERE username = $1 ORDER BY reply_time DESC"
	rows, err := backend.db.Query(query, username)
    if err != nil {
		log.Println(err)
        return nil, err
    }
	defer rows.Close()
	
	replies := []model.Reply{}
	for rows.Next() {
		var reply model.Reply
		err := rows.Scan(&reply.Id,
						 &reply.Username,
						 &reply.DiscussionId,
						 &reply.Content,
						 &reply.ReplyTime)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		replies = append(replies, reply)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return replies, nil
}

func (backend *PostgresBackend) SelectReplyById(id int) (*model.Reply, error) {
	var reply model.Reply
	err := backend.db.QueryRow("SELECT id, username, discussion_id, content, TO_CHAR(reply_time, 'YYYY-MM-DD HH24:MI:SS') FROM replies WHERE id = $1", id).
		   Scan(&reply.Id,
				&reply.Username,
				&reply.DiscussionId,
				&reply.Content,
				&reply.ReplyTime)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &reply, nil
}

func (backend *PostgresBackend) DeleteReply(username string, id int) error {
	tx, _ := backend.db.Begin()
	defer tx.Rollback()

	query := "DELETE FROM replies WHERE username = $1 AND id = $2"

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
