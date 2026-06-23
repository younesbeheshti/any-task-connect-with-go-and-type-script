package dto

// UserStats holds profile statistics.
type UserStats struct {
	Rating         float64 `json:"rating"`
	CompletedCount int     `json:"completedCount"`
	RatingCount    int     `json:"ratingCount"`
	WalletBalance  int64   `json:"walletBalance"`
	LockedBalance  int64   `json:"lockedBalance"`
}
