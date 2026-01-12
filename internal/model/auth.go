package model

import "time"

type UserLoginGoogle struct {
	TokenID string `json:"token_id" binding:"required"`
}

type UserRegisterDevice struct {
	DeviceID    string `json:"device_id" binding:"required"`
	DeviceType  string `json:"device_type"`
	DeviceToken string `json:"device_token" binding:"required"`
}

type UserApplication struct {
	CompanyName       string `json:"company_name" binding:"required"`
	LicenceIssueDate  string `json:"licence_issue_date" binding:"required"`
	LicenceExpiryDate string `json:"licence_expiry_date" binding:"required"`
	FullName          string `json:"full_name" binding:"required"`
	Email             string `json:"email" binding:"required"`
	Phone             string `json:"phone" binding:"required"`
	Address           string `json:"address" binding:"required"`
	Password          string `json:"password" binding:"omitempty"`
	VATNumber         string `json:"vat_number" binding:"required"`
	CompanyTypeID     int    `json:"company_type_id" binding:"required"`
	ActivityFieldID   int    `json:"activity_field_id" binding:"required"`
	RoleID            int    `json:"role_id" binding:"required"` // 1 user, 2 dealer, 3 logistic, 4 broker, 5 car service
}

type TempUser struct {
	LicenceIssueDate  *time.Time `json:"licence_issue_date" binding:"required"`
	LicenceExpiryDate *time.Time `json:"licence_expiry_date" binding:"required"`
	CompanyName       *string    `json:"company_name" binding:"required"`
	FullName          *string    `json:"full_name" binding:"required"`
	Email             *string    `json:"email" binding:"required"`
	Phone             *string    `json:"phone" binding:"required"`
	Address           *string    `json:"address" binding:"required"`
	VATNumber         *string    `json:"vat_number" binding:"required"`
	CompanyTypeID     *int       `json:"company_type_id" binding:"required"`
	ActivityFieldID   *int       `json:"activity_field_id" binding:"required"`
	RoleID            *int       `json:"role_id" binding:"required"` // 1 user, 2 dealer, 3 logistic, 4 broker, 5 car service
	DocumentsID       *int       `json:"documents_id"`
}

type UserApplicationDocuments struct {
	Licence    string `json:"licence" binding:"required"`
	Memorandum string `json:"memorandum" binding:"required"`
	CopyOfID   string `json:"copy_of_id" binding:"required"`
}

type UserEmailConfirmationRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

type UserPhoneConfirmationRequest struct {
	Phone string `json:"phone" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}

type UserLoginEmail struct {
	Email string `json:"email" binding:"required,email"`
}

type UserLoginPhone struct {
	Phone string `json:"phone" binding:"required"`
}

type ThirdPartyLogin struct {
	Password       string `json:"password" binding:"required"`
	ID             int    `json:"id"`
	RoleID         int    `json:"role_id" binding:"required"`
	FirstTimeLogin bool   `json:"first_time_login"`
}

type UserByEmail struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	OTP      string `json:"otp"`
	ID       int    `json:"id"`
	RoleID   int    `json:"role_id"`
}

// type ThirdPartyLogin struct {
// 	Email          string `json:"email" binding:"required,email"`
// 	Password       string `json:"password" binding:"required"`
// 	RoleID         int    `json:"role_id" binding:"required"`
// 	FirstTimeLogin bool   `json:"first_time_login"`
// }

type UserByPhone struct {
	Phone    string `json:"phone"`
	Username string `json:"username"`
	OTP      string `json:"otp"`
	ID       int    `json:"id"`
}

type LoginFiberResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type ThirdPartyLoginFiberResponse struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	FirstTimeLogin bool   `json:"first_time_login"`
}

type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone" binding:"required,min=6,max=15"`
	Username string `json:"username" binding:"required,min=3,max=20"`
}
