package usecases

type UserUsecases interface {
	ChangeUsersBlockStatus(userID int) error
}