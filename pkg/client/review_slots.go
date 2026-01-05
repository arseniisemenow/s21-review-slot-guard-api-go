package client

import (
	"context"
	"fmt"
	"time"
)

// ReviewSlot represents a review slot from the calendar
type ReviewSlot struct {
	ID        string    `json:"id"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Type      string    `json:"type"` // FREE_TIME or BOOKED_TIME
	IsOnline  bool      `json:"isOnline"`
	VCLink    *string   `json:"vcLink"`
	School    string    `json:"school"`
}

// ReviewBooking represents a booked review
type ReviewBooking struct {
	ID            string        `json:"id"`
	SlotID        string        `json:"slotId"`
	Start         time.Time     `json:"start"`
	End           time.Time     `json:"end"`
	ProjectName   string        `json:"projectName"`
	VerifierLogin string        `json:"verifierLogin"`
	IsOnline      bool          `json:"isOnline"`
	Status        string        `json:"status"`
}

// GetReviewSlots fetches available and booked review slots within a date range
func (c *Client) GetReviewSlots(ctx context.Context, from, to time.Time) ([]ReviewSlot, []ReviewBooking, error) {
	fromStr := from.Format("2006-01-02T15:04:05.000Z")
	toStr := to.Format("2006-01-02T15:04:05.000Z")

	resp, err := c.GetCalendarEvents(ctx, fromStr, toStr)
	if err != nil {
		return nil, nil, err
	}

	var slots []ReviewSlot
	var bookings []ReviewBooking

	for _, event := range resp.CalendarEventS21.GetMyCalendarEvents {
		// Filter for student_check events (review slots)
		if event.EventCode != "student_check" {
			continue
		}

		for _, slot := range event.EventSlots {
			start, _ := time.Parse(time.RFC3339, slot.Start)
			end, _ := time.Parse(time.RFC3339, slot.End)

			schoolShortName := ""
			if slot.School.ShortName != "" {
				schoolShortName = slot.School.ShortName
			}

			reviewSlot := ReviewSlot{
				ID:       slot.ID,
				Start:    start,
				End:      end,
				Type:     slot.Type,
				IsOnline: false,
				School:   schoolShortName,
			}

			slots = append(slots, reviewSlot)
		}

		// Process bookings (already booked reviews)
		for _, booking := range event.Bookings {
			start, _ := time.Parse(time.RFC3339, booking.EventSlot.Start)
			end, _ := time.Parse(time.RFC3339, booking.EventSlot.End)

			reviewBooking := ReviewBooking{
				ID:            booking.ID,
				SlotID:        booking.EventSlot.ID,
				Start:         start,
				End:           end,
				ProjectName:   booking.Task.GoalName,
				VerifierLogin:  booking.VerifierUser.Login,
				IsOnline:      booking.IsOnline,
				Status:        booking.BookingStatus,
			}

			bookings = append(bookings, reviewBooking)
		}
	}

	return slots, bookings, nil
}

// GetAvailableReviewSlots fetches only available (free) review slots
func (c *Client) GetAvailableReviewSlots(ctx context.Context, from, to time.Time) ([]ReviewSlot, error) {
	slots, _, err := c.GetReviewSlots(ctx, from, to)
	if err != nil {
		return nil, err
	}

	var available []ReviewSlot
	for _, slot := range slots {
		if slot.Type == "FREE_TIME" {
			available = append(available, slot)
		}
	}

	return available, nil
}

// GetBookedReviews fetches only booked review slots
func (c *Client) GetBookedReviews(ctx context.Context, from, to time.Time) ([]ReviewBooking, error) {
	_, bookings, err := c.GetReviewSlots(ctx, from, to)
	if err != nil {
		return nil, err
	}

	// Deduplicate bookings by ID
	seen := make(map[string]bool)
	var result []ReviewBooking
	for _, b := range bookings {
		if !seen[b.ID] {
			seen[b.ID] = true
			result = append(result, b)
		}
	}

	return result, nil
}

// AddReviewSlot adds a new review slot to the timetable
func (c *Client) AddReviewSlot(ctx context.Context, start, end time.Time) ([]ReviewSlot, error) {
	startStr := start.Format("2006-01-02T15:04:05.000Z")
	endStr := end.Format("2006-01-02T15:04:05.000Z")

	resp, err := c.AddEventToTimetable(ctx, startStr, endStr)
	if err != nil {
		return nil, err
	}

	var slots []ReviewSlot
	for _, event := range resp.Student.AddEventToTimetable {
		if event.EventCode != "student_check" {
			continue
		}

		for _, slot := range event.EventSlots {
			s, _ := time.Parse(time.RFC3339, slot.Start)
			e, _ := time.Parse(time.RFC3339, slot.End)

			schoolShortName := ""
			if slot.School.ShortName != "" {
				schoolShortName = slot.School.ShortName
			}

			reviewSlot := ReviewSlot{
				ID:       slot.ID,
				Start:    s,
				End:      e,
				Type:     slot.Type,
				IsOnline: false,
				School:   schoolShortName,
			}

			slots = append(slots, reviewSlot)
		}
	}

	return slots, nil
}

// UpdateReviewSlot changes the time of an existing review slot
func (c *Client) UpdateReviewSlot(ctx context.Context, slotID string, newStart, newEnd time.Time) (*ReviewSlot, error) {
	startStr := newStart.Format("2006-01-02T15:04:05.000Z")
	endStr := newEnd.Format("2006-01-02T15:04:05.000Z")

	resp, err := c.ChangeEventSlot(ctx, slotID, startStr, endStr)
	if err != nil {
		return nil, err
	}

	// Find the updated slot
	for _, event := range resp.Student.ChangeEventSlot.EventSlots {
		if event.ID == slotID {
			start, _ := time.Parse(time.RFC3339, event.Start)
			end, _ := time.Parse(time.RFC3339, event.End)

			schoolShortName := ""
			if event.School.ShortName != "" {
				schoolShortName = event.School.ShortName
			}

			return &ReviewSlot{
				ID:       event.ID,
				Start:    start,
				End:      end,
				Type:     event.Type,
				IsOnline: false,
				School:   schoolShortName,
			}, nil
		}
	}

	return nil, fmt.Errorf("updated slot not found in response")
}

// RemoveReviewSlot cancels/deletes a review slot
func (c *Client) RemoveReviewSlot(ctx context.Context, slotID string) error {
	_, err := c.DeleteEventSlot(ctx, slotID)
	return err
}

// CancelReview cancels a review by deleting its slot
// This is an alias for RemoveReviewSlot for clarity
func (c *Client) CancelReview(ctx context.Context, slotID string) error {
	return c.RemoveReviewSlot(ctx, slotID)
}
