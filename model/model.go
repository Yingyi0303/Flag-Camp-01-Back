package model

type User struct {
	Id int	`json:"id"`
	Username string	`json:"username"`
	Password string	`json:"password"`
	Role string	`json:"role"`
}

type Discussion struct {
	Id int	`json:"id"`
	Username string	`json:"username"`
	Topic string	`json:"topic"`
	Content string	`json:"content"`
	LastUpdateTime string	`json:"last_update_time"`
}

type Reply struct {
	Id int	`json:"id"`
	Username string	`json:"username"`
	DiscussionId int	`json:"discussion_id"`
	Content string	`json:"content"`
	ReplyTime string	`json:"reply_time"`
}

type DiscussionDto struct {
	Discussion Discussion	`json:"discussion"`
	Replies *[]Reply	`json:"replies,omitempty"`
}
