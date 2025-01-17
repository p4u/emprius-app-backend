package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/emprius/emprius-app-backend/db"
)

// convertBookingToResponse converts a db.Booking to a BookingResponse
func convertBookingToResponse(booking *db.Booking) BookingResponse {
	return BookingResponse{
		ID:            booking.ID.Hex(),
		ToolID:        booking.ToolID,
		FromUserID:    booking.FromUserID.Hex(),
		ToUserID:      booking.ToUserID.Hex(),
		StartDate:     booking.StartDate.Unix(),
		EndDate:       booking.EndDate.Unix(),
		Contact:       booking.Contact,
		Comments:      booking.Comments,
		BookingStatus: string(booking.BookingStatus),
		CreatedAt:     booking.CreatedAt,
		UpdatedAt:     booking.UpdatedAt,
	}
}

// HandleGetBookingRequests handles GET /bookings/requests
func (a *API) HandleGetBookingRequests(r *Request) (interface{}, error) {
	if r.UserID == "" {
		return nil, ErrUnauthorized
	}

	// Get user from database
	user, err := a.database.UserService.GetUserByEmail(r.Context.Request.Context(), r.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	bookings, err := a.database.BookingService.GetUserRequests(r.Context.Request.Context(), user.ID)
	if err != nil {
		return nil, ErrInternalServerError
	}

	response := make([]BookingResponse, len(bookings))
	for i, booking := range bookings {
		response[i] = convertBookingToResponse(booking)
	}

	return response, nil
}

// HandleGetBookingPetitions handles GET /bookings/petitions
func (a *API) HandleGetBookingPetitions(r *Request) (interface{}, error) {
	if r.UserID == "" {
		return nil, ErrUnauthorized
	}

	// Get user from database
	user, err := a.database.UserService.GetUserByEmail(r.Context.Request.Context(), r.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	bookings, err := a.database.BookingService.GetUserPetitions(r.Context.Request.Context(), user.ID)
	if err != nil {
		return nil, ErrInternalServerError
	}

	response := make([]BookingResponse, len(bookings))
	for i, booking := range bookings {
		response[i] = convertBookingToResponse(booking)
	}

	return response, nil
}

// HandleGetBooking handles GET /bookings/{bookingId}
func (a *API) HandleGetBooking(r *Request) (interface{}, error) {
	bookingID, err := primitive.ObjectIDFromHex(chi.URLParam(r.Context.Request, "bookingId"))
	if err != nil {
		return nil, ErrInvalidRequestBodyData
	}

	booking, err := a.database.BookingService.Get(r.Context.Request.Context(), bookingID)
	if err != nil {
		return nil, ErrInternalServerError
	}
	if booking == nil {
		return nil, ErrBookingNotFound
	}

	return convertBookingToResponse(booking), nil
}

// HandleAcceptPetition handles POST /bookings/petitions/{petitionId}/accept
func (a *API) HandleAcceptPetition(r *Request) (interface{}, error) {
	if r.UserID == "" {
		return nil, ErrUnauthorized
	}

	// Get user from database
	user, err := a.database.UserService.GetUserByEmail(r.Context.Request.Context(), r.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	petitionID, err := primitive.ObjectIDFromHex(chi.URLParam(r.Context.Request, "petitionId"))
	if err != nil {
		return nil, ErrInvalidRequestBodyData
	}

	booking, err := a.database.BookingService.Get(r.Context.Request.Context(), petitionID)
	if err != nil {
		return nil, ErrInternalServerError
	}
	if booking == nil {
		return nil, ErrBookingNotFound
	}

	// Verify user is the tool owner
	if booking.ToUserID != user.ID {
		return nil, ErrOnlyOwnerCanAccept
	}

	// Verify booking is in PENDING state
	if booking.BookingStatus != db.BookingStatusPending {
		return nil, ErrCanOnlyAcceptPending
	}

	err = a.database.BookingService.UpdateStatus(r.Context.Request.Context(), petitionID, db.BookingStatusAccepted)
	if err != nil {
		return nil, ErrInternalServerError
	}

	return nil, nil
}

// HandleDenyPetition handles POST /bookings/petitions/{petitionId}/deny
func (a *API) HandleDenyPetition(r *Request) (interface{}, error) {
	if r.UserID == "" {
		return nil, ErrUnauthorized
	}

	// Get user from database
	user, err := a.database.UserService.GetUserByEmail(r.Context.Request.Context(), r.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	petitionID, err := primitive.ObjectIDFromHex(chi.URLParam(r.Context.Request, "petitionId"))
	if err != nil {
		return nil, ErrInvalidRequestBodyData
	}

	booking, err := a.database.BookingService.Get(r.Context.Request.Context(), petitionID)
	if err != nil {
		return nil, ErrInternalServerError
	}
	if booking == nil {
		return nil, ErrBookingNotFound
	}

	// Verify user is the tool owner
	if booking.ToUserID != user.ID {
		return nil, ErrOnlyOwnerCanDeny
	}

	// Verify booking is in PENDING state
	if booking.BookingStatus != db.BookingStatusPending {
		return nil, ErrCanOnlyDenyPending
	}

	err = a.database.BookingService.UpdateStatus(r.Context.Request.Context(), petitionID, db.BookingStatusRejected)
	if err != nil {
		return nil, ErrInternalServerError
	}

	return nil, nil
}

// HandleCancelRequest handles POST /bookings/request/{petitionId}/cancel
func (a *API) HandleCancelRequest(r *Request) (interface{}, error) {
	if r.UserID == "" {
		return nil, ErrUnauthorized
	}

	// Get user from database
	user, err := a.database.UserService.GetUserByEmail(r.Context.Request.Context(), r.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	petitionID, err := primitive.ObjectIDFromHex(chi.URLParam(r.Context.Request, "petitionId"))
	if err != nil {
		return nil, ErrInvalidRequestBodyData
	}

	booking, err := a.database.BookingService.Get(r.Context.Request.Context(), petitionID)
	if err != nil {
		return nil, ErrInternalServerError
	}
	if booking == nil {
		return nil, ErrBookingNotFound
	}

	// Verify user is the requester
	if booking.FromUserID != user.ID {
		return nil, ErrOnlyRequesterCanCancel
	}

	// Verify booking is in PENDING state
	if booking.BookingStatus != db.BookingStatusPending {
		return nil, ErrCanOnlyCancelPending
	}

	err = a.database.BookingService.UpdateStatus(r.Context.Request.Context(), petitionID, db.BookingStatusCancelled)
	if err != nil {
		return nil, ErrInternalServerError
	}

	return nil, nil
}

// HandleReturnBooking handles POST /bookings/{bookingId}/return
func (a *API) HandleReturnBooking(r *Request) (interface{}, error) {
	if r.UserID == "" {
		return nil, ErrUnauthorized
	}

	// Get user from database
	user, err := a.database.UserService.GetUserByEmail(r.Context.Request.Context(), r.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	bookingID, err := primitive.ObjectIDFromHex(chi.URLParam(r.Context.Request, "bookingId"))
	if err != nil {
		return nil, ErrInvalidRequestBodyData
	}

	booking, err := a.database.BookingService.Get(r.Context.Request.Context(), bookingID)
	if err != nil {
		return nil, ErrInternalServerError
	}
	if booking == nil {
		return nil, ErrBookingNotFound
	}

	// Verify user is the tool owner
	if booking.ToUserID != user.ID {
		return nil, ErrOnlyOwnerCanReturn
	}

	err = a.database.BookingService.UpdateStatus(r.Context.Request.Context(), bookingID, db.BookingStatusReturned)
	if err != nil {
		return nil, ErrInternalServerError
	}

	return nil, nil
}

// HandleGetPendingRatings handles GET /bookings/rates
func (a *API) HandleGetPendingRatings(r *Request) (interface{}, error) {
	if r.UserID == "" {
		return nil, ErrUnauthorized
	}

	// Get user from database
	user, err := a.database.UserService.GetUserByEmail(r.Context.Request.Context(), r.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	bookings, err := a.database.BookingService.GetPendingRatings(r.Context.Request.Context(), user.ID)
	if err != nil {
		return nil, ErrInternalServerError
	}

	response := make([]BookingResponse, len(bookings))
	for i, booking := range bookings {
		response[i] = convertBookingToResponse(booking)
	}

	return response, nil
}

// RateRequest represents the request body for rating a booking
type RateRequest struct {
	Rating    int    `json:"rating"`
	BookingID string `json:"bookingId"`
}

// HandleCreateBooking handles POST /bookings
func (a *API) HandleCreateBooking(r *Request) (interface{}, error) {
	if r.UserID == "" {
		return nil, ErrUnauthorized
	}

	// Get user from database
	fromUser, err := a.database.UserService.GetUserByEmail(r.Context.Request.Context(), r.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	var req CreateBookingRequest
	if err := json.Unmarshal(r.Data, &req); err != nil {
		return nil, ErrInvalidRequestBodyData
	}

	toolID, err := strconv.ParseInt(req.ToolID, 10, 64)
	if err != nil {
		return nil, ErrInvalidRequestBodyData
	}

	// Get tool to verify it exists and get owner ID
	tool, err := a.database.ToolService.GetToolByID(r.Context.Request.Context(), toolID)
	if err != nil {
		return nil, ErrInternalServerError
	}
	if tool == nil {
		return nil, ErrToolNotFound
	}

	toUser, err := a.database.UserService.GetUserByID(r.Context.Request.Context(), tool.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Create booking request
	dbReq := &db.CreateBookingRequest{
		ToolID:    fmt.Sprintf("%d", toolID),
		StartDate: time.Unix(req.StartDate, 0),
		EndDate:   time.Unix(req.EndDate, 0),
		Contact:   req.Contact,
		Comments:  req.Comments,
	}

	booking, err := a.database.BookingService.Create(r.Context.Request.Context(), dbReq, fromUser.ID, toUser.ID)
	if err != nil {
		if err.Error() == "booking dates conflict with existing booking" {
			return nil, ErrBookingDatesConflict
		}
		return nil, ErrInternalServerError
	}

	return convertBookingToResponse(booking), nil
}

// HandleRateBooking handles POST /bookings/rates
func (a *API) HandleRateBooking(r *Request) (interface{}, error) {
	if r.UserID == "" {
		return nil, ErrUnauthorized
	}

	// Get user from database
	user, err := a.database.UserService.GetUserByEmail(r.Context.Request.Context(), r.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	var rateReq RateRequest
	if err := json.Unmarshal(r.Data, &rateReq); err != nil {
		return nil, ErrInvalidRequestBodyData
	}

	bookingID, err := primitive.ObjectIDFromHex(rateReq.BookingID)
	if err != nil {
		return nil, ErrInvalidRequestBodyData
	}

	booking, err := a.database.BookingService.Get(r.Context.Request.Context(), bookingID)
	if err != nil {
		return nil, ErrInternalServerError
	}
	if booking == nil {
		return nil, ErrBookingNotFound
	}

	// Verify user is involved in the booking
	if booking.FromUserID != user.ID && booking.ToUserID != user.ID {
		return nil, ErrUserNotInvolved
	}

	// Verify rating value
	if rateReq.Rating < 1 || rateReq.Rating > 5 {
		return nil, ErrInvalidRating
	}

	// TODO: Implement rating logic once rating schema is defined

	return nil, nil
}
