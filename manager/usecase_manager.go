package manager

import "github.com/epulskyyy/majoo-test-2022/usecase"

type UseCaseManager interface {
	UserUseCase() usecase.IUserUseCase
	AuthUseCase() usecase.IAuthUseCase
	MerchantUseCase() usecase.IMerchantUseCase
	OutletUseCase() usecase.IOutletUseCase
	TransactionUseCase() usecase.ITransactionUseCase

}

type useCaseManager struct {
	repo   RepoManager
	client ClientManager
}

func (uc *useCaseManager) UserUseCase() usecase.IUserUseCase {
	return usecase.NewUserUseCase(uc.repo.UserRepo(), uc.client.Redis())
}
func (uc *useCaseManager) MerchantUseCase() usecase.IMerchantUseCase {
	return usecase.NewMerchantUseCase(uc.repo.MerchantRepo(), uc.client.Redis())
}
func (uc *useCaseManager) OutletUseCase() usecase.IOutletUseCase {
	return usecase.NewOutletUseCase(uc.repo.OutletRepo(), uc.client.Redis())
}
func (uc *useCaseManager) TransactionUseCase() usecase.ITransactionUseCase {
	return usecase.NewTransactionUseCase(uc.repo.TransactionRepo(),uc.repo.MerchantRepo(), uc.client.Redis())
}

func (uc *useCaseManager) AuthUseCase() usecase.IAuthUseCase {
	return usecase.NewAuthUseCase(uc.repo.UserRepo(), uc.client.Redis())
}

func NewUseCaseManger(repo RepoManager, client ClientManager) UseCaseManager {
	return &useCaseManager{repo: repo, client: client}
}
