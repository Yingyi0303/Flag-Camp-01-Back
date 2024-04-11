package service

import (
	"fmt"
	"log"
	"propertyManagement/backend"
	"propertyManagement/model"
)

func AddReply(reply *model.Reply) (*model.Reply, error) {
	success, err := backend.PGBackend.DiscussionExists(reply.DiscussionId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if !success {
		return nil, nil
	}
	result, err := backend.PGBackend.InsertReply(reply.DiscussionId, reply.Username, reply.Content)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Printf("Service added reply: %d\n", result.Id)
	return result, nil
}

func GetMyReplies(username string) ([]model.Reply, error) {
	replies, err := backend.PGBackend.SelectAllRepliesByUsername(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Printf("Service fetched all replies by user: %s\n", username)
	return replies, nil
}

func RemoveReply(username string, id int) error {
	err := backend.PGBackend.DeleteReply(username, id)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Printf("Service removed reply: %d\n", id)
	return nil
}