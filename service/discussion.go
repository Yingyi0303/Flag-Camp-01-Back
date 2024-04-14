package service

import (
	"fmt"
	"log"
	"propertyManagement/backend"
	"propertyManagement/model"
)

func AddDiscussion(discussion *model.Discussion) (*model.Discussion, error) {
	result, err := backend.PGBackend.InsertDiscussion(discussion.Username, discussion.Subject, discussion.Content)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Printf("Service added discussion: %d\n", result.Id)
	return result, nil
}

func GetAllDiscussions() ([]model.Discussion, error) {
	discussions, err := backend.PGBackend.SelectAllDiscussions()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println("Service fetched all discussions")
	return discussions, nil
}

func GetMyDiscussions(username string) ([]model.Discussion, error) {
	discussions, err := backend.PGBackend.SelectAllDiscussionsByUsername(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Printf("Service fetched all discussions by user: %s\n", username)
	return discussions, nil
}

func GetDiscussionDetails(id int) (*model.DiscussionDto, error) {
	success, err := backend.PGBackend.DiscussionExists(id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if !success {
		return nil, nil
	}
	discussion, err := backend.PGBackend.SelectDiscussionById(id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	replies, err := backend.PGBackend.SelectRepliesByDiscussionId(id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	discussionDto := model.DiscussionDto {
		Discussion: *discussion,
		Replies: replies,
	}
	fmt.Printf("Service fetched details for discussion: %d\n", id)
	return &discussionDto, nil
}

func RemoveDiscussion(username string, id int) error {
	err := backend.PGBackend.DeleteDiscussion(username, id)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Printf("Service removed discussion: %d\n", id)
	return nil
}
