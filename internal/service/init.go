package service

import "github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/usecase"

var (
	APIKey     usecase.APIKeyUseCase
	User       usecase.UserUseCase
	Follow     usecase.FollowUseCase
	Permission usecase.PermissionUseCase
)

func Init() {
	APIKey = usecase.NewAPIKeyUseCase()
	User = usecase.NewUserUseCase()
	Follow = usecase.NewFollowUseCase()
	Permission = usecase.NewPermissionUseCase()
}
