package manager

import (
	"github.com/epulskyyy/majoo-test-2022/repository"
)

type RepoManager interface {
	UserRepo() repository.IUserRepository
	MerchantRepo() repository.IMerchantRepository
	OutletRepo() repository.IOutletRepository
	TransactionRepo() repository.ITransactionRepository
}

type repoManager struct {
	infra Infra
}

func (rm *repoManager) UserRepo() repository.IUserRepository {
	return repository.NewUserRepository(rm.infra.GormDB())
}

func (rm *repoManager) MerchantRepo() repository.IMerchantRepository {
	return repository.NewMerchantRepository(rm.infra.GormDB())
}

func (rm *repoManager) OutletRepo() repository.IOutletRepository {
	return repository.NewOutletRepository(rm.infra.GormDB())
}
func (rm *repoManager) TransactionRepo() repository.ITransactionRepository {
	return repository.NewTransactionRepository(rm.infra.GormDB())
}


func NewRepoManager(infra Infra) RepoManager {
	return &repoManager{infra}
}
