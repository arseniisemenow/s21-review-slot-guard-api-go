//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/arseni/s21-api-client/pkg/client"
)

// getCredentials reads credentials from environment variables
func getCredentials(t *testing.T) (login, password string) {
	login = os.Getenv("S21_LOGIN")
	password = os.Getenv("S21_PASSWORD")

	if login == "" || password == "" {
		t.Skip("S21_LOGIN and S21_PASSWORD environment variables must be set for integration tests")
	}

	return login, password
}

func TestIntegration_Authenticate(t *testing.T) {
	login, password := getCredentials(t)

	authConfig := &client.AuthConfig{
		Login:    login,
		Password: password,
	}

	c := client.NewClient(authConfig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tokenResp, err := c.Authenticate(ctx)
	if err != nil {
		t.Fatalf("Authenticate() failed = %v", err)
	}

	if tokenResp == nil {
		t.Fatal("tokenResp is nil")
	}

	if tokenResp.AccessToken == "" {
		t.Error("AccessToken is empty")
	}

	if tokenResp.TokenType != "Bearer" {
		t.Errorf("TokenType = %s, want Bearer", tokenResp.TokenType)
	}

	if tokenResp.ExpiresIn <= 0 {
		t.Errorf("ExpiresIn = %d, want > 0", tokenResp.ExpiresIn)
	}

	t.Logf("Successfully authenticated, token expires in %d seconds", tokenResp.ExpiresIn)
}

func TestIntegration_GetCurrentUser(t *testing.T) {
	login, password := getCredentials(t)

	authConfig := &client.AuthConfig{
		Login:    login,
		Password: password,
	}

	c := client.NewClient(authConfig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := c.GetCurrentUser(ctx)
	if err != nil {
		t.Fatalf("GetCurrentUser() failed = %v", err)
	}

	if user == nil {
		t.Fatal("user is nil")
	}

	if user.User.GetCurrentUser.Login != login {
		t.Errorf("Login = %s, want %s", user.User.GetCurrentUser.Login, login)
	}

	if user.User.GetCurrentUser.ID == "" {
		t.Error("ID is empty")
	}

	if user.User.GetCurrentUser.CurrentSchoolStudentID == "" {
		t.Error("CurrentSchoolStudentID is empty")
	}

	t.Logf("Successfully got current user: %s %s (login: %s, id: %s)",
		user.User.GetCurrentUser.FirstName,
		user.User.GetCurrentUser.LastName,
		user.User.GetCurrentUser.Login,
		user.User.GetCurrentUser.ID)
}

func TestIntegration_GetUserNotifications(t *testing.T) {
	login, password := getCredentials(t)

	authConfig := &client.AuthConfig{
		Login:    login,
		Password: password,
	}

	c := client.NewClient(authConfig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	paging := client.PagingInput{
		Offset: 0,
		Limit:  10,
	}

	resp, err := c.GetUserNotifications(ctx, paging)
	if err != nil {
		t.Fatalf("GetUserNotifications() failed = %v", err)
	}

	if resp == nil {
		t.Fatal("resp is nil")
	}

	notifications := resp.S21Notification.GetS21Notifications.Notifications
	if len(notifications) == 0 {
		t.Log("No notifications found")
		return
	}

	t.Logf("Successfully got %d notifications (total: %d)",
		len(notifications), resp.S21Notification.GetS21Notifications.TotalCount)

	// Check first notification
	first := notifications[0]
	if first.ID == "" {
		t.Error("First notification ID is empty")
	}
	if first.Message == "" {
		t.Error("First notification message is empty")
	}
}

func TestIntegration_GetUserNotificationsCount(t *testing.T) {
	login, password := getCredentials(t)

	authConfig := &client.AuthConfig{
		Login:    login,
		Password: password,
	}

	c := client.NewClient(authConfig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	countUnread, err := c.GetUserNotificationsCount(ctx, false)
	if err != nil {
		t.Fatalf("GetUserNotificationsCount(false) failed = %v", err)
	}

	countTotal, err := c.GetUserNotificationsCount(ctx, true)
	if err != nil {
		t.Fatalf("GetUserNotificationsCount(true) failed = %v", err)
	}

	t.Logf("Notification counts - unread: %d, total: %d",
		countUnread.S21Notification.GetS21NotificationsCount,
		countTotal.S21Notification.GetS21NotificationsCount)
}

func TestIntegration_GetStudentStageGroups(t *testing.T) {
	login, password := getCredentials(t)

	authConfig := &client.AuthConfig{
		Login:    login,
		Password: password,
	}

	c := client.NewClient(authConfig)

	// First get current user to obtain student ID
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := c.GetCurrentUser(ctx)
	if err != nil {
		t.Fatalf("GetCurrentUser() failed = %v", err)
	}

	studentID := user.User.GetCurrentUser.CurrentSchoolStudentID
	if studentID == "" {
		t.Fatal("CurrentSchoolStudentID is empty")
	}

	resp, err := c.GetStudentStageGroups(ctx, studentID)
	if err != nil {
		t.Fatalf("GetStudentStageGroups() failed = %v", err)
	}

	groups := resp.School21.LoadStudentStageGroups
	if len(groups) == 0 {
		t.Log("No stage groups found")
		return
	}

	t.Logf("Successfully got %d stage groups", len(groups))

	for _, group := range groups {
		t.Logf("  - Wave: %s, EduForm: %s, Active: %v",
			group.StageGroupS21.WaveName,
			group.StageGroupS21.EduForm,
			group.StageGroupS21.Active)
	}
}

func TestIntegration_GetStudentGraphTemplate(t *testing.T) {
	login, password := getCredentials(t)

	authConfig := &client.AuthConfig{
		Login:    login,
		Password: password,
	}

	c := client.NewClient(authConfig)

	// First get current user to obtain student ID
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := c.GetCurrentUser(ctx)
	if err != nil {
		t.Fatalf("GetCurrentUser() failed = %v", err)
	}

	studentID := user.User.GetCurrentUser.CurrentSchoolStudentID

	resp, err := c.GetStudentGraphTemplate(ctx, studentID, nil)
	if err != nil {
		t.Fatalf("GetStudentGraphTemplate() failed = %v", err)
	}

	graph := resp.HolyGraph.GetStudentGraphTemplate
	t.Logf("Successfully got graph template: %d nodes, %d edges",
		len(graph.Nodes), len(graph.Edges))

	// Print first few nodes
	for i, node := range graph.Nodes {
		if i >= 3 {
			break
		}
		t.Logf("  Node %d: %s (type: %s, label: %s)",
			i, node.ID, node.Type, node.Label)
		if len(node.Items) > 0 {
			for j, item := range node.Items {
				if j >= 2 {
					break
				}
				name := item.ProjectName
				if name == "" {
					name = item.Course.ProjectName
				}
				t.Logf("    - Item: %s", name)
			}
		}
	}
}

func TestIntegration_GetMyReviews(t *testing.T) {
	login, password := getCredentials(t)

	authConfig := &client.AuthConfig{
		Login:    login,
		Password: password,
	}

	c := client.NewClient(authConfig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get reviews for the next month
	to := time.Now().AddDate(0, 1, 0).Format(time.RFC3339)
	limit := 10

	resp, err := c.GetMyReviews(ctx, to, limit)
	if err != nil {
		t.Fatalf("GetMyReviews() failed = %v", err)
	}

	reviews := resp.Student.GetMyUpcomingBookings
	t.Logf("Successfully got %d upcoming reviews", len(reviews))

	for i, review := range reviews {
		if i >= 3 {
			break
		}
		t.Logf("  Review %d: %s on %s (status: %s)",
			i, review.Task.GoalName, review.EventSlot.Start, review.BookingStatus)
	}
}

func TestIntegration_GetCalendarEvents(t *testing.T) {
	login, password := getCredentials(t)

	authConfig := &client.AuthConfig{
		Login:    login,
		Password: password,
	}

	c := client.NewClient(authConfig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get events for the next week
	from := time.Now().Format(time.RFC3339)
	to := time.Now().AddDate(0, 0, 7).Format(time.RFC3339)

	resp, err := c.GetCalendarEvents(ctx, from, to)
	if err != nil {
		t.Fatalf("GetCalendarEvents() failed = %v", err)
	}

	events := resp.CalendarEventS21.GetMyCalendarEvents
	t.Logf("Successfully got %d calendar events", len(events))

	for i, event := range events {
		if i >= 5 {
			break
		}
		t.Logf("  Event %d: %s (%s) from %s to %s",
			i, event.EventType, event.EventCode, event.Start, event.End)
	}
}

func TestIntegration_AddDeleteEventSlot(t *testing.T) {
	login, password := getCredentials(t)

	authConfig := &client.AuthConfig{
		Login:    login,
		Password: password,
	}

	c := client.NewClient(authConfig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Add an event slot for tomorrow
	tomorrow := time.Now().AddDate(0, 0, 1)
	start := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 0, 0, 0, time.UTC)
	end := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 30, 0, 0, time.UTC)

	addResp, err := c.AddEventToTimetable(ctx, start.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("AddEventToTimetable() failed = %v", err)
	}

	events := addResp.Student.AddEventToTimetable
	if len(events) == 0 {
		t.Fatal("No events returned from AddEventToTimetable")
	}

	addedEvent := events[0]
	if len(addedEvent.EventSlots) == 0 {
		t.Fatal("No event slots in added event")
	}

	eventSlotID := addedEvent.EventSlots[0].ID
	t.Logf("Successfully added event slot: %s", eventSlotID)

	// Clean up - delete the event slot
	delResp, err := c.DeleteEventSlot(ctx, eventSlotID)
	if err != nil {
		t.Fatalf("DeleteEventSlot() failed = %v", err)
	}

	if !delResp.Student.DeleteEventSlot {
		t.Error("DeleteEventSlot returned false")
	}

	t.Logf("Successfully deleted event slot: %s", eventSlotID)
}

func TestIntegration_FullWorkflow(t *testing.T) {
	login, password := getCredentials(t)

	authConfig := &client.AuthConfig{
		Login:    login,
		Password: password,
	}

	c := client.NewClient(authConfig)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Step 1: Get current user
	t.Log("Step 1: Getting current user...")
	user, err := c.GetCurrentUser(ctx)
	if err != nil {
		t.Fatalf("GetCurrentUser() failed = %v", err)
	}
	t.Logf("  Logged in as: %s %s (%s)",
		user.User.GetCurrentUser.FirstName,
		user.User.GetCurrentUser.LastName,
		user.User.GetCurrentUser.Login)

	// Step 2: Get notifications
	t.Log("Step 2: Getting notifications...")
	notifResp, err := c.GetUserNotifications(ctx, client.PagingInput{Offset: 0, Limit: 5})
	if err != nil {
		t.Fatalf("GetUserNotifications() failed = %v", err)
	}
	t.Logf("  Got %d notifications", len(notifResp.S21Notification.GetS21Notifications.Notifications))

	// Step 3: Get stage groups
	t.Log("Step 3: Getting stage groups...")
	studentID := user.User.GetCurrentUser.CurrentSchoolStudentID
	groupsResp, err := c.GetStudentStageGroups(ctx, studentID)
	if err != nil {
		t.Fatalf("GetStudentStageGroups() failed = %v", err)
	}
	t.Logf("  Got %d stage groups", len(groupsResp.School21.LoadStudentStageGroups))

	// Step 4: Get graph template
	t.Log("Step 4: Getting graph template...")
	graphResp, err := c.GetStudentGraphTemplate(ctx, studentID, nil)
	if err != nil {
		t.Fatalf("GetStudentGraphTemplate() failed = %v", err)
	}
	graph := graphResp.HolyGraph.GetStudentGraphTemplate
	t.Logf("  Got graph with %d nodes and %d edges", len(graph.Nodes), len(graph.Edges))

	// Find allowed projects (nodes with items)
	var allowedProjects []string
	for _, node := range graph.Nodes {
		if len(node.Items) > 0 {
			for _, item := range node.Items {
				name := item.ProjectName
				if name == "" {
					name = item.Course.ProjectName
				}
				if name != "" {
					allowedProjects = append(allowedProjects, name)
				}
			}
		}
	}
	t.Logf("  Found %d projects", len(allowedProjects))

	// Step 5: Get calendar events
	t.Log("Step 5: Getting calendar events...")
	from := time.Now().Format(time.RFC3339)
	to := time.Now().AddDate(0, 0, 7).Format(time.RFC3339)
	eventsResp, err := c.GetCalendarEvents(ctx, from, to)
	if err != nil {
		t.Fatalf("GetCalendarEvents() failed = %v", err)
	}
	t.Logf("  Got %d calendar events for next week", len(eventsResp.CalendarEventS21.GetMyCalendarEvents))

	// Step 6: Get upcoming reviews
	t.Log("Step 6: Getting upcoming reviews...")
	toMonth := time.Now().AddDate(0, 1, 0).Format(time.RFC3339)
	reviewsResp, err := c.GetMyReviews(ctx, toMonth, 10)
	if err != nil {
		t.Fatalf("GetMyReviews() failed = %v", err)
	}
	t.Logf("  Got %d upcoming reviews", len(reviewsResp.Student.GetMyUpcomingBookings))

	t.Log("Full workflow completed successfully!")
}

// BenchmarkGetCurrentUser benchmarks the GetCurrentUser request
func BenchmarkGetCurrentUser(b *testing.B) {
	login, password := getCredentials(&testing.T{})

	authConfig := &client.AuthConfig{
		Login:    login,
		Password: password,
	}

	c := client.NewClient(authConfig)
	ctx := context.Background()

	// Authenticate once
	if _, err := c.Authenticate(ctx); err != nil {
		b.Fatalf("Authenticate() failed = %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.GetCurrentUser(ctx)
		if err != nil {
			b.Fatalf("GetCurrentUser() failed = %v", err)
		}
	}
}

// ExampleGetCurrentUser demonstrates how to use the client
func ExampleGetCurrentUser() {
	// This example shows how to get the current user
	// In real usage, you would get credentials from environment or config
	authConfig := &client.AuthConfig{
		Login:    os.Getenv("S21_LOGIN"),
		Password: os.Getenv("S21_PASSWORD"),
	}

	c := client.NewClient(authConfig)
	ctx := context.Background()

	user, err := c.GetCurrentUser(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Hello, %s %s!\n",
		user.User.GetCurrentUser.FirstName,
		user.User.GetCurrentUser.LastName)
}
