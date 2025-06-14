package userservice

type UserService interface {
	RegisterUser(login, password, email string) error
	LoginUser(login, password string) error
}
type userRepository interface {
	CreateUser(login, password, email string) error
	GetUserPassword(login string) (string, error)
}
type userService struct {
	userRepo userRepository
}
