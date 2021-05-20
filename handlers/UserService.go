package handlers

type UserService struct {
	UserServ interface{}
}

type UserOps interface {
	Subscribe()
	Subscribed()
	CheckAndPersist()
}
