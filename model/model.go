package model

type Response struct {
	Username string	`json:"username"`
	Role string	`json:"role"`
	Token string	`json:"token"`
}

type User struct {
	Id int	`json:"id"`
	Username string	`json:"username"`
	Password string	`json:"password"`
	Role string	`json:"role"`
}

type Discussion struct {
	Id int	`json:"id"`
	Username string	`json:"username"`
	Subject string	`json:"subject"`
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
	Replies []Reply	`json:"replies"`
}

type Maintenance struct {
	Id int	`json:"id"`
	Username string	`json:"username"`
	Subject string	`json:"subject"`
	Content string	`json:"content"`
	Reply string	`json:"reply"`
	Completed bool	`json:"completed"`
	LastUpdateTime string	`json:"last_update_time"`
}

type Bill struct {
	Id int	`json:"id"`
	Username string	`json:"username"`
	MaintenanceId int	`json:"maintenance_id"`
	Item string	`json:"item"`
	Amount int	`json:"amount"`
	BillTime string	`json:"bill_time"`
}

type Payment struct {
	Id int	`json:"id"`
	Username string	`json:"username"`
	Item string	`json:"item"`
	Amount int	`json:"amount"`
	PaymentTime string	`json:"payment_time"`
}

type BalanceDto struct {
	Balance int	`json:"balance"`
	Bills []Bill	`json:"bills"`
	Payments []Payment	`json:"payments"`
}
