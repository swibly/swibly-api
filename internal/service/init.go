package service

import "github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/usecase"

var (
	User       usecase.UserUseCase
	Follow     usecase.FollowUseCase
	Permission usecase.PermissionUseCase
)

func Init() {
	User = usecase.NewUserUseCase()
	Follow = usecase.NewFollowUseCase()
	Permission = usecase.NewPermissionUseCase()
}
