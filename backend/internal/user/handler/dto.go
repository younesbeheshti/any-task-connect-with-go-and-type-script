package handler

import (
	"time"

	userdomain "github.com/younesbeheshti/any-task-connect/backend/internal/user/domain"
)

// UserDTO matches front/docs/api-contracts.md User schema.
type UserDTO struct {
	ID             string                `json:"id"`
	FullName       string                `json:"fullName"`
	Phone          string                `json:"phone"`
	Email          *string               `json:"email"`
	City           string                `json:"city,omitempty"`
	Role           string                `json:"role"`
	AvatarURL      *string               `json:"avatarUrl"`
	Verification   userdomain.Verification `json:"verification"`
	Rating         float64               `json:"rating"`
	CompletedCount int                   `json:"completedCount"`
	CreatedAt      time.Time             `json:"createdAt"`
}

// PublicUserDTO is the public profile response.
type PublicUserDTO struct {
	ID             string                `json:"id"`
	FullName       string                `json:"fullName"`
	AvatarURL      *string               `json:"avatarUrl"`
	Rating         float64               `json:"rating"`
	CompletedCount int                   `json:"completedCount"`
	Verification   userdomain.Verification `json:"verification"`
	CreatedAt      time.Time             `json:"createdAt"`
}

// UpdateProfileRequest matches PATCH /me contract.
type UpdateProfileRequest struct {
	FullName  *string `json:"fullName"`
	City      *string `json:"city"`
	AvatarURL *string `json:"avatarUrl"`
	Bio       *string `json:"bio"`
}

// ToUserDTO maps domain user to API DTO.
func ToUserDTO(u *userdomain.User) UserDTO {
	avatar := u.Avatar
	return UserDTO{
		ID: u.ID.String(), FullName: u.FullName, Phone: u.Phone, Email: u.Email,
		City: u.CityTitle, Role: u.Role.APIRole(), AvatarURL: avatar,
		Verification: u.Verification(), Rating: u.Rating,
		CompletedCount: u.CompletedTasks, CreatedAt: u.CreatedAt,
	}
}

// ToPublicUserDTO maps public profile to DTO.
func ToPublicUserDTO(p *userdomain.PublicProfile) PublicUserDTO {
	return PublicUserDTO{
		ID: p.ID.String(), FullName: p.FullName, AvatarURL: p.Avatar,
		Rating: p.Rating, CompletedCount: p.CompletedCount,
		Verification: p.Verification, CreatedAt: p.CreatedAt,
	}
}
