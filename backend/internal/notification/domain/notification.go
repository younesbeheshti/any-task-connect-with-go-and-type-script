package domain

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType categorizes notifications.
type NotificationType string

const (
	NotificationTypeTaskAssigned  NotificationType = "TASK_ASSIGNED"
	NotificationTypeTaskCompleted NotificationType = "TASK_COMPLETED"
	NotificationTypePayment         NotificationType = "PAYMENT"
	NotificationTypeApplication     NotificationType = "APPLICATION"
	NotificationTypeMessage         NotificationType = "MESSAGE"
	NotificationTypeDispute         NotificationType = "DISPUTE"
	NotificationTypeSystem          NotificationType = "SYSTEM"
	NotificationTypeReview          NotificationType = "REVIEW"
)

// Notification represents an in-app notification.
type Notification struct {
	ID        uuid.UUID        `json:"id"`
	UserID    uuid.UUID        `json:"-"`
	Title     string           `json:"title"`
	Body      string           `json:"desc"`
	Type      NotificationType `json:"type"`
	IsRead    bool             `json:"isRead"`
	DeepLink  *string          `json:"deepLink"`
	CreatedAt time.Time        `json:"createdAt"`
}

// CreateNotificationInput holds notification creation data.
type CreateNotificationInput struct {
	UserID   uuid.UUID
	Title    string
	Body     string
	Type     NotificationType
	DeepLink *string
}
