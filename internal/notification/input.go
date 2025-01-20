package notification

type GetByIdNotificationInput struct {
	Id int `uri:"id" binding:"required"`
}

type UpdateNotificationInput struct {
	Status string `json:"status" binding:"required"`
}
