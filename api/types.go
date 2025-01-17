package api

import (
	"time"

	"github.com/emprius/emprius-app-backend/db"
	"github.com/emprius/emprius-app-backend/types"
)

// Response is the default response of the API
type Response struct {
	Header ResponseHeader `json:"header"`
	Data   any            `json:"data,omitempty"`
}

// ResponseHeader is the header of the response
type ResponseHeader struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	ErrorCode int    `json:"errorCode,omitempty"`
}
type Register struct {
	UserEmail         string `json:"email"`
	RegisterAuthToken string `json:"invitationToken"`
	UserProfile
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Token    string    `json:"token"`
	Expirity time.Time `json:"expirity"`
}

type UserProfile struct {
	Name      string       `json:"name"`
	Community string       `json:"community"`
	Location  *db.Location `json:"location,omitempty"`
	Active    *bool        `json:"active,omitempty"`
	Avatar    []byte       `json:"avatar,omitempty"`
	Password  string       `json:"password,omitempty"`
}

type UsersWrapper struct {
	Users []db.User `json:"users"`
}

// Tool is the type of the tool
type Tool struct {
	ID               int64            `json:"id"`
	Title            string           `json:"title"`
	Description      string           `json:"description"`
	IsAvailable      *bool            `json:"isAvailable"`
	MayBeFree        *bool            `json:"mayBeFree"`
	AskWithFee       *bool            `json:"askWithFee"`
	Cost             *uint64          `json:"cost"`
	Images           []types.HexBytes `json:"images"`
	TransportOptions []int            `json:"transportOptions"`
	Category         int              `json:"category"`
	Location         db.Location      `json:"location"`
	EstimatedValue   uint64           `json:"estimatedValue"`
	Height           uint32           `json:"height"`
	Weight           uint32           `json:"weight"`
}

type ToolID struct {
	ID int64 `json:"id"`
}

type ToolsWrapper struct {
	Tools []db.Tool `json:"tools"`
}

// ToolSearch is the type of the tool search
type ToolSearch struct {
	Term             string  `json:"term"`
	Categories       []int   `json:"categories"`
	Distance         int     `json:"distance"`
	MaxCost          *uint64 `json:"maxCost"`
	MayBeFree        *bool   `json:"mayBeFree"`
	AvailableFrom    int     `json:"availableFrom"`
	TransportOptions []int   `json:"transportOptions"`
}

type Info struct {
	Users      int               `json:"users"`
	Tools      int               `json:"tools"`
	Categories []db.ToolCategory `json:"categories"`
	Transports []db.Transport    `json:"transports"`
}

// CreateBookingRequest represents the request to create a new booking
type CreateBookingRequest struct {
	ToolID    string `json:"toolId"`
	StartDate int64  `json:"startDate"`
	EndDate   int64  `json:"endDate"`
	Contact   string `json:"contact"`
	Comments  string `json:"comments"`
}

// BookingResponse represents the API response for a booking
type BookingResponse struct {
	ID            string    `json:"id"`
	ToolID        string    `json:"toolId"`
	FromUserID    string    `json:"fromUserId"`
	ToUserID      string    `json:"toUserId"`
	StartDate     int64     `json:"startDate"`
	EndDate       int64     `json:"endDate"`
	Contact       string    `json:"contact"`
	Comments      string    `json:"comments"`
	BookingStatus string    `json:"bookingStatus"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
