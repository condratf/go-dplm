package userservice

type UserService interface {
	RegisterUser(login, password string) error
	LoginUser(login, password string) error
}
type userRepository interface {
	CreateUser(login, password string) error
	GetUserPassword(login string) (string, error)
}
type userService struct {
	userRepo userRepository
}
