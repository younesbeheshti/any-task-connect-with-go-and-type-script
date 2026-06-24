package domain

import (
	"time"

	"github.com/google/uuid"
)

type Account string

const (
	AccountUserWallet      Account = "USER_WALLET"
	AccountEscrow          Account = "ESCROW_ACCOUNT"
	AccountPlatformRevenue Account = "PLATFORM_REVENUE"
	AccountWithdraw        Account = "WITHDRAW_ACCOUNT"
	AccountSystem          Account = "SYSTEM_ACCOUNT"
)

type Entry struct {
	ID            uuid.UUID
	TransactionID uuid.UUID
	Account       Account
	Debit         int64
	Credit        int64
	BalanceAfter  int64
	CreatedAt     time.Time
}

type EntryPair struct {
	Debit  Entry
	Credit Entry
}

func ValidateDoubleEntry(entries []Entry) bool {
	var totalDebit, totalCredit int64
	for _, e := range entries {
		totalDebit += e.Debit
		totalCredit += e.Credit
	}
	return totalDebit == totalCredit
}
