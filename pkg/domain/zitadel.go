package domain

type ZitadelEventUserCreatedRequest struct {
	Email struct {
		Email      string `json:"email" binding:"required,email"`
		IsVerified bool   `json:"isVerified" binding:"required"`
	} `json:"email" binding:"required"`
	Organization struct {
		OrganizationID string `json:"orgId" binding:"required"`
	} `json:"organization" binding:"required"`
	Profile struct {
		FirstName string `binding:"required" json:"givenName"`
		LastName  string `binding:"required" json:"familyName"`
	} `json:"profile" binding:"required"`
	UserName string `json:"userName" binding:"required"`
}
type ZitadelEventUserCreated struct {
	InstanceID     string                         `json:"instanceId" binding:"required"`
	OrganizationID string                         `json:"orgID" binding:"required"`
	Request        ZitadelEventUserCreatedRequest `json:"request" binding:"required"`
	UserID         string                         `json:"userID" binding:"required"`
}
