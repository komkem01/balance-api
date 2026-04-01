package entitiesinf

import (
	"context"
	"time"

	"balance/app/modules/entities/ent"

	"github.com/google/uuid"
)

// ObjectEntity defines the interface for object entity operations such as create, retrieve, update, and soft delete.
type ExampleEntity interface {
	CreateExample(ctx context.Context, userID uuid.UUID) (*ent.Example, error)
	GetExampleByID(ctx context.Context, id uuid.UUID) (*ent.Example, error)
	UpdateExampleByID(ctx context.Context, id uuid.UUID, status ent.ExampleStatus) (*ent.Example, error)
	SoftDeleteExampleByID(ctx context.Context, id uuid.UUID) error
	ListExamplesByStatus(ctx context.Context, status ent.ExampleStatus) ([]*ent.Example, error)
}
type ExampleTwoEntity interface {
	CreateExampleTwo(ctx context.Context, userID uuid.UUID) (*ent.Example, error)
}

type GenderEntity interface {
	CreateGender(ctx context.Context, name string, isActive bool) (*ent.GenderEntity, error)
	GetGenderByID(ctx context.Context, id string) (*ent.GenderEntity, error)
	UpdateGender(ctx context.Context, id string, name *string, isActive *bool) (*ent.GenderEntity, error)
	DeleteGender(ctx context.Context, id string) error
	ListGenders(ctx context.Context, isActive *bool) ([]*ent.GenderEntity, error)
}

type PrefixEntity interface {
	CreatePrefix(ctx context.Context, genderID string, name string, isActive bool) (*ent.PrefixEntity, error)
	GetPrefixByID(ctx context.Context, id string) (*ent.PrefixEntity, error)
	UpdatePrefix(ctx context.Context, id string, name *string, isActive *bool) (*ent.PrefixEntity, error)
	DeletePrefix(ctx context.Context, id string) error
	ListPrefixes(ctx context.Context, isActive *bool) ([]*ent.PrefixEntity, error)
}

type MemberEntity interface {
	CreateMember(ctx context.Context, genderID *string, prefixID *string, firstName string, lastName string, displayName string, phone string) (*ent.MemberEntity, error)
	CreateMemberWithAccount(ctx context.Context, genderID *string, prefixID *string, firstName string, lastName string, displayName string, phone string, username string, password string) (*ent.MemberEntity, error)
	GetMemberByID(ctx context.Context, id string) (*ent.MemberEntity, error)
	UpdateMember(ctx context.Context, id string, genderID *string, prefixID *string, firstName *string, lastName *string, displayName *string, phone *string, lastLogin *time.Time) (*ent.MemberEntity, error)
	UpdateMemberSettings(ctx context.Context, id string, preferredCurrency *string, preferredLanguage *string, notifyBudget *bool, notifySecurity *bool, notifyWeekly *bool) (*ent.MemberEntity, error)
	DeleteMember(ctx context.Context, id string) error
	DeleteMemberWithAccounts(ctx context.Context, id string) error
	ListMembers(ctx context.Context) ([]*ent.MemberEntity, error)
}

type MemberAccountEntity interface {
	CreateMemberAccount(ctx context.Context, memberID *string, username string, password string) (*ent.MemberAccountEntity, error)
	GetMemberAccountByID(ctx context.Context, id string) (*ent.MemberAccountEntity, error)
	UpdateMemberAccount(ctx context.Context, id string, memberID *string, username *string, password *string) (*ent.MemberAccountEntity, error)
	DeleteMemberAccount(ctx context.Context, id string) error
	DeleteMemberAccountByMemberID(ctx context.Context, memberID string) error
	ListMemberAccounts(ctx context.Context) ([]*ent.MemberAccountEntity, error)
}

type WalletEntity interface {
	CreateWallet(ctx context.Context, memberID *string, name string, balance float64, currency string, colorCode string, isActive bool) (*ent.WalletEntity, error)
	GetWalletByID(ctx context.Context, id string) (*ent.WalletEntity, error)
	UpdateWallet(ctx context.Context, id string, memberID *string, name *string, balance *float64, currency *string, colorCode *string, isActive *bool) (*ent.WalletEntity, error)
	DeleteWallet(ctx context.Context, id string) error
	ListWallets(ctx context.Context, isActive *bool) ([]*ent.WalletEntity, error)
}

type CategoryEntity interface {
	CreateCategory(ctx context.Context, memberID *string, name string, categoryType ent.CategoryType, iconName string, colorCode string) (*ent.CategoryEntity, error)
	GetCategoryByID(ctx context.Context, id string) (*ent.CategoryEntity, error)
	UpdateCategory(ctx context.Context, id string, memberID *string, name *string, categoryType *ent.CategoryType, iconName *string, colorCode *string) (*ent.CategoryEntity, error)
	DeleteCategory(ctx context.Context, id string) error
	ListCategories(ctx context.Context, memberID *string, categoryType *ent.CategoryType) ([]*ent.CategoryEntity, error)
}

type TransactionEntity interface {
	CreateTransaction(ctx context.Context, walletID *string, categoryID *string, amount float64, transactionType ent.TransactionType, transactionDate *time.Time, note string, imageURL string) (*ent.TransactionEntity, error)
	CreateTransactionWithWalletAdjust(ctx context.Context, walletID *string, categoryID *string, amount float64, transactionType ent.TransactionType, transactionDate *time.Time, note string, imageURL string) (*ent.TransactionEntity, error)
	GetTransactionByID(ctx context.Context, id string) (*ent.TransactionEntity, error)
	UpdateTransaction(ctx context.Context, id string, walletID *string, categoryID *string, amount *float64, transactionType *ent.TransactionType, transactionDate *time.Time, note *string, imageURL *string) (*ent.TransactionEntity, error)
	UpdateTransactionWithWalletAdjust(ctx context.Context, id string, walletID *string, categoryID *string, amount *float64, transactionType *ent.TransactionType, transactionDate *time.Time, note *string, imageURL *string) (*ent.TransactionEntity, error)
	DeleteTransaction(ctx context.Context, id string) error
	DeleteTransactionWithWalletAdjust(ctx context.Context, id string) error
	ListTransactions(ctx context.Context, walletID *string, categoryID *string, transactionType *ent.TransactionType) ([]*ent.TransactionEntity, error)
}

type BudgetEntity interface {
	CreateBudget(ctx context.Context, memberID *string, categoryID *string, amount float64, period ent.BudgetPeriod, startDate *time.Time, endDate *time.Time) (*ent.BudgetEntity, error)
	GetBudgetByID(ctx context.Context, id string) (*ent.BudgetEntity, error)
	UpdateBudget(ctx context.Context, id string, memberID *string, categoryID *string, amount *float64, period *ent.BudgetPeriod, startDate *time.Time, endDate *time.Time) (*ent.BudgetEntity, error)
	DeleteBudget(ctx context.Context, id string) error
	ListBudgets(ctx context.Context, memberID *string, categoryID *string, period *ent.BudgetPeriod) ([]*ent.BudgetEntity, error)
}
