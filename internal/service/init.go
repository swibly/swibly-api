package service

import "github.com/swibly/swibly-api/internal/service/usecase"

var (
	APIKey        usecase.APIKeyUseCase
	User          usecase.UserUseCase
	Follow        usecase.FollowUseCase
	Permission    usecase.PermissionUseCase
	Project       usecase.ProjectUseCase
	Component     usecase.ComponentUseCase
	PasswordReset usecase.PasswordResetUseCase
	Notification  usecase.NotificationUseCase
)

func Init() {
	APIKey = usecase.NewAPIKeyUseCase()
	User = usecase.NewUserUseCase()
	Follow = usecase.NewFollowUseCase()
	Permission = usecase.NewPermissionUseCase()
	Project = usecase.NewProjectUseCase()
	Component = usecase.NewComponentUseCase()
	PasswordReset = usecase.NewPasswordResetUseCase()
	Notification = usecase.NewNotificationUseCase()
}
