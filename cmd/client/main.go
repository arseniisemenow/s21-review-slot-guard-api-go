package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arseni/s21-api-client/pkg/client"
)

// formatDateTime formats time for the API (RFC3339 without nanoseconds)
func formatDateTime(t time.Time) string {
	return t.Format("2006-01-02T15:04:05Z")
}

// formatDateTimeMilli formats time with milliseconds for the API
func formatDateTimeMilli(t time.Time) string {
	return t.Format("2006-01-02T15:04:05.000Z")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	login := os.Getenv("S21_LOGIN")
	password := os.Getenv("S21_PASSWORD")

	if login == "" || password == "" {
		log.Fatal("S21_LOGIN and S21_PASSWORD environment variables must be set")
	}

	authConfig := &client.AuthConfig{
		Login:    login,
		Password: password,
	}

	// Create client with optional headers from environment
	opts := []client.ClientOption{}
	if schoolID := os.Getenv("S21_SCHOOL_ID"); schoolID != "" {
		opts = append(opts, client.WithSchoolID(schoolID))
	}
	if userRole := os.Getenv("S21_USER_ROLE"); userRole != "" {
		opts = append(opts, client.WithUserRole(userRole))
	}
	if productID := os.Getenv("S21_EDU_PRODUCT_ID"); productID != "" {
		opts = append(opts, client.WithEduProductID(productID))
	}
	if orgUnitID := os.Getenv("S21_EDU_ORG_UNIT_ID"); orgUnitID != "" {
		opts = append(opts, client.WithEduOrgUnitID(orgUnitID))
	}

	c := client.NewClient(authConfig, opts...)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := os.Args[1]

	switch cmd {
	case "user":
		getCurrentUser(ctx, c)
	case "notifications":
		getNotifications(ctx, c)
	case "reviews":
		getReviews(ctx, c)
	case "projects":
		getProjects(ctx, c)
	case "calendar":
		getCalendar(ctx, c)
	case "review-slots":
		handleReviewSlots(ctx, c)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

// parseDateTime parses a datetime string in multiple formats
func parseDateTime(s string) (time.Time, error) {
	// Try various formats
	formats := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, f := range formats {
		t, err := time.Parse(f, s)
		if err == nil {
			return t, nil
		}
	}

	// Try parsing as date only (set time to now)
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("unable to parse datetime: %s (try formats like: 2025-01-15, 2025-01-15 14:30, 2025-01-15T14:30:00Z)", s)
}

func handleReviewSlots(ctx context.Context, c *client.Client) {
	if len(os.Args) < 3 {
		printReviewSlotsUsage()
		os.Exit(1)
	}

	subCmd := os.Args[2]

	switch subCmd {
	case "get":
		getReviewSlotsCmd(ctx, c)
	case "add":
		addReviewSlotCmd(ctx, c)
	case "update":
		updateReviewSlotCmd(ctx, c)
	case "remove", "rm":
		removeReviewSlotCmd(ctx, c)
	default:
		fmt.Printf("Unknown review-slots command: %s\n", subCmd)
		printReviewSlotsUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: client <command>")
	fmt.Println("\nCommands:")
	fmt.Println("  user          - Get current user info")
	fmt.Println("  notifications - Get user notifications")
	fmt.Println("  reviews       - Get upcoming reviews")
	fmt.Println("  projects      - Get available projects")
	fmt.Println("  calendar      - Get calendar events")
	fmt.Println("  review-slots  - Manage review slots (get/add/update/remove)")
	fmt.Println("\nEnvironment variables:")
	fmt.Println("  S21_LOGIN           - Your 21-school login")
	fmt.Println("  S21_PASSWORD        - Your 21-school password")
	fmt.Println("  S21_SCHOOL_ID       - School ID (from browser, required for some operations)")
	fmt.Println("  S21_USER_ROLE       - User role (e.g., STUDENT)")
	fmt.Println("  S21_EDU_PRODUCT_ID  - Edu Product ID (from browser)")
	fmt.Println("  S21_EDU_ORG_UNIT_ID - Edu Org Unit ID (from browser)")
}

func getCurrentUser(ctx context.Context, c *client.Client) {
	user, err := c.GetCurrentUser(ctx)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	u := user.User.GetCurrentUser
	fmt.Printf("ID: %s\n", u.ID)
	fmt.Printf("Login: %s\n", u.Login)
	fmt.Printf("Name: %s %s\n", u.FirstName, u.LastName)
	fmt.Printf("Student ID: %s\n", u.CurrentSchoolStudentID)
}

func getNotifications(ctx context.Context, c *client.Client) {
	resp, err := c.GetUserNotifications(ctx, client.PagingInput{Offset: 0, Limit: 10})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	notif := resp.S21Notification.GetS21Notifications
	fmt.Printf("Total notifications: %d\n", notif.TotalCount)
	fmt.Printf("Showing %d:\n", len(notif.Notifications))

	for i, n := range notif.Notifications {
		fmt.Printf("\n%d. [%s] %s\n", i+1, n.RelatedObjectType, n.GroupName)
		fmt.Printf("   %s\n", n.Message)
		fmt.Printf("   %s (read: %v)\n", n.Time, n.WasRead)
	}
}

func getReviews(ctx context.Context, c *client.Client) {
	to := formatDateTime(time.Now().AddDate(0, 1, 0))
	resp, err := c.GetMyReviews(ctx, to, 10)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	reviews := resp.Student.GetMyUpcomingBookings
	fmt.Printf("Upcoming reviews: %d\n", len(reviews))

	for i, r := range reviews {
		start, _ := time.Parse(time.RFC3339, r.EventSlot.Start)
		fmt.Printf("\n%d. %s\n", i+1, r.Task.GoalName)
		fmt.Printf("   Time: %s\n", start.Format("2006-01-02 15:04 MST"))
		fmt.Printf("   Status: %s\n", r.BookingStatus)
	}
}

func getProjects(ctx context.Context, c *client.Client) {
	// First get current user to obtain student ID
	user, err := c.GetCurrentUser(ctx)
	if err != nil {
		log.Fatalf("Error getting user: %v", err)
	}

	studentID := user.User.GetCurrentUser.CurrentSchoolStudentID

	// Get graph template
	graph, err := c.GetStudentGraphTemplate(ctx, studentID, nil)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	g := graph.HolyGraph.GetStudentGraphTemplate
	fmt.Printf("Available projects (%d nodes):\n", len(g.Nodes))

	for _, node := range g.Nodes {
		if len(node.Items) > 0 {
			fmt.Printf("\n[%s] %s\n", node.Label, node.Label)
			for _, item := range node.Items {
				var name string
				if item.Goal != nil && item.Goal.ProjectName != "" {
					name = item.Goal.ProjectName
				} else if item.Course != nil && item.Course.ProjectName != "" {
					name = item.Course.ProjectName
				}
				if name != "" {
					fmt.Printf("  - %s\n", name)
				}
			}
		}
	}
}

func getCalendar(ctx context.Context, c *client.Client) {
	from := formatDateTimeMilli(time.Now())
	to := formatDateTimeMilli(time.Now().AddDate(0, 0, 7))

	resp, err := c.GetCalendarEvents(ctx, from, to)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	events := resp.CalendarEventS21.GetMyCalendarEvents
	fmt.Printf("Calendar events (next 7 days): %d\n", len(events))

	for i, e := range events {
		start, _ := time.Parse(time.RFC3339, e.Start)
		fmt.Printf("\n%d. %s (%s)\n", i+1, e.EventType, e.EventCode)
		fmt.Printf("   Time: %s - %s\n", start.Format("2006-01-02 15:04"), e.End)
		if len(e.EventSlots) > 0 {
			fmt.Printf("   Slots: %d\n", len(e.EventSlots))
		}
		if len(e.Bookings) > 0 {
			fmt.Printf("   Bookings: %d\n", len(e.Bookings))
		}
	}
}

func getReviewSlots(ctx context.Context, c *client.Client) {
	// Get slots for the next week
	from := time.Now()
	to := time.Now().AddDate(0, 0, 7)

	slots, bookings, err := c.GetReviewSlots(ctx, from, to)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Count available vs booked
	availableCount := 0
	bookedCount := 0
	for _, s := range slots {
		if s.Type == "FREE_TIME" {
			availableCount++
		} else {
			bookedCount++
		}
	}

	fmt.Printf("Review slots (next 7 days):\n")
	fmt.Printf("  Available: %d\n", availableCount)
	fmt.Printf("  Booked: %d\n\n", bookedCount)

	// Show available slots
	if availableCount > 0 {
		fmt.Println("Available slots:")
		for i, s := range slots {
			if s.Type == "FREE_TIME" {
				fmt.Printf("  %d. %s - %s [%s] (ID: %s)\n", i+1,
					s.Start.Format("2006-01-02 15:04"),
					s.End.Format("15:04"),
					s.School,
					s.ID)
			}
		}
	}

	// Show booked reviews
	if len(bookings) > 0 {
		fmt.Println("\nBooked reviews:")
		for i, b := range bookings {
			fmt.Printf("  %d. %s on %s\n", i+1, b.ProjectName,
				b.Start.Format("2006-01-02 15:04"))
			fmt.Printf("     Verifier: %s | Status: %s | SlotID: %s\n", b.VerifierLogin, b.Status, b.SlotID)
		}
	}
}

func printReviewSlotsUsage() {
	fmt.Println("Usage: client review-slots <command>")
	fmt.Println("\nCommands:")
	fmt.Println("  get [days]           - Show available and booked review slots (default: 7 days)")
	fmt.Println("  add <start> <end>    - Add a new review slot")
	fmt.Println("                        Format: YYYY-MM-DD HH:MM or YYYY-MM-DDTHH:MM:SSZ")
	fmt.Println("                        Example: client review-slots add '2025-01-15 14:00' '2025-01-15 14:30'")
	fmt.Println("  update <id> <start> <end> - Update an existing review slot")
	fmt.Println("                        Example: client review-slots update slot-123 '2025-01-15 15:00' '2025-01-15 15:30'")
	fmt.Println("  remove|rm <id>       - Remove a review slot")
	fmt.Println("                        Example: client review-slots remove slot-123")
	fmt.Println("\nExamples:")
	fmt.Println("  client review-slots get           # Show slots for next 7 days")
	fmt.Println("  client review-slots get 30        # Show slots for next 30 days")
	fmt.Println("  client review-slots add '2025-01-15 14:00' '2025-01-15 14:30'")
}

func getReviewSlotsCmd(ctx context.Context, c *client.Client) {
	days := 7 // default
	if len(os.Args) >= 4 {
		fmt.Sscanf(os.Args[3], "%d", &days)
	}

	from := time.Now()
	to := time.Now().AddDate(0, 0, days)

	slots, bookings, err := c.GetReviewSlots(ctx, from, to)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Count available vs booked
	availableCount := 0
	bookedCount := 0
	for _, s := range slots {
		if s.Type == "FREE_TIME" {
			availableCount++
		} else {
			bookedCount++
		}
	}

	fmt.Printf("Review slots (next %d days):\n", days)
	fmt.Printf("  Available: %d\n", availableCount)
	fmt.Printf("  Booked: %d\n\n", bookedCount)

	// Show available slots
	if availableCount > 0 {
		fmt.Println("Available slots:")
		idx := 1
		for _, s := range slots {
			if s.Type == "FREE_TIME" {
				fmt.Printf("  %d. %s - %s [%s]\n", idx,
					s.Start.Format("2006-01-02 15:04"),
					s.End.Format("15:04"),
					s.School)
				fmt.Printf("     ID: %s\n", s.ID)
				idx++
			}
		}
	}

	// Show booked reviews
	if len(bookings) > 0 {
		fmt.Println("\nBooked reviews:")
		for i, b := range bookings {
			fmt.Printf("  %d. %s on %s\n", i+1, b.ProjectName,
				b.Start.Format("2006-01-02 15:04"))
			fmt.Printf("     Verifier: %s | Status: %s\n", b.VerifierLogin, b.Status)
			fmt.Printf("     SlotID: %s\n", b.SlotID)
		}
	}
}

func addReviewSlotCmd(ctx context.Context, c *client.Client) {
	if len(os.Args) < 5 {
		fmt.Println("Usage: client review-slots add <start> <end>")
		fmt.Println("Example: client review-slots add '2025-01-15 14:00' '2025-01-15 14:30'")
		os.Exit(1)
	}

	startStr := os.Args[3]
	endStr := os.Args[4]

	start, err := parseDateTime(startStr)
	if err != nil {
		log.Fatalf("Error parsing start time: %v", err)
	}

	end, err := parseDateTime(endStr)
	if err != nil {
		log.Fatalf("Error parsing end time: %v", err)
	}

	if end.Before(start) {
		log.Fatal("Error: end time must be after start time")
	}

	fmt.Printf("Adding review slot: %s - %s\n", start.Format("2006-01-02 15:04"), end.Format("15:04"))

	slots, err := c.AddReviewSlot(ctx, start, end)
	if err != nil {
		log.Fatalf("Error adding review slot: %v", err)
	}

	fmt.Println("\nReview slot added successfully!")
	if len(slots) > 0 {
		fmt.Println("Created slots:")
		for i, s := range slots {
			fmt.Printf("  %d. ID: %s\n", i+1, s.ID)
			fmt.Printf("     Type: %s\n", s.Type)
			fmt.Printf("     Time: %s - %s\n", s.Start.Format("2006-01-02 15:04"), s.End.Format("15:04"))
		}
	}
}

func updateReviewSlotCmd(ctx context.Context, c *client.Client) {
	if len(os.Args) < 6 {
		fmt.Println("Usage: client review-slots update <slot-id> <start> <end>")
		fmt.Println("Example: client review-slots update slot-123 '2025-01-15 15:00' '2025-01-15 15:30'")
		os.Exit(1)
	}

	slotID := os.Args[3]
	startStr := os.Args[4]
	endStr := os.Args[5]

	start, err := parseDateTime(startStr)
	if err != nil {
		log.Fatalf("Error parsing start time: %v", err)
	}

	end, err := parseDateTime(endStr)
	if err != nil {
		log.Fatalf("Error parsing end time: %v", err)
	}

	if end.Before(start) {
		log.Fatal("Error: end time must be after start time")
	}

	fmt.Printf("Updating review slot %s: %s - %s\n", slotID, start.Format("2006-01-02 15:04"), end.Format("15:04"))

	slot, err := c.UpdateReviewSlot(ctx, slotID, start, end)
	if err != nil {
		log.Fatalf("Error updating review slot: %v", err)
	}

	fmt.Println("\nReview slot updated successfully!")
	fmt.Printf("  ID: %s\n", slot.ID)
	fmt.Printf("  Type: %s\n", slot.Type)
	fmt.Printf("  Time: %s - %s\n", slot.Start.Format("2006-01-02 15:04"), slot.End.Format("15:04"))
	if slot.School != "" {
		fmt.Printf("  School: %s\n", slot.School)
	}
}

func removeReviewSlotCmd(ctx context.Context, c *client.Client) {
	if len(os.Args) < 4 {
		fmt.Println("Usage: client review-slots remove <slot-id>")
		fmt.Println("Example: client review-slots remove slot-123")
		os.Exit(1)
	}

	slotID := os.Args[3]

	fmt.Printf("Removing review slot: %s\n", slotID)

	err := c.RemoveReviewSlot(ctx, slotID)
	if err != nil {
		log.Fatalf("Error removing review slot: %v", err)
	}

	fmt.Println("Review slot removed successfully!")
}
