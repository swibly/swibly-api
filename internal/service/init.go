package service

import "github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/usecase"

func Init() {
	usecase.UserInstance = usecase.NewUserUseCase()
}
