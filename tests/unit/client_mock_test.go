//go:build mock
// +build mock

package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/arseni/s21-api-client/pkg/client"
)

// mockAuthServer creates a mock auth server
func mockAuthServer(tokenResp *client.TokenResponse) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/realms/EduPowerKeycloak/protocol/openid-connect/token" {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			// Check content type
			if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tokenResp)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
}

// mockGraphQLServer creates a mock GraphQL server
func mockGraphQLServer(handler func(http.ResponseWriter, *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/services/graphql" {
			// Check auth header
			auth := r.Header.Get("Authorization")
			if auth == "" || auth != "Bearer mock-token" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"errors": []map[string]interface{}{
						{"message": "Unauthorized"},
					},
				})
				return
			}

			// Check content type
			if r.Header.Get("Content-Type") != "application/json" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			handler(w, r)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestMockClient_AuthenticateSuccess(t *testing.T) {
	tokenResp := &client.TokenResponse{
		AccessToken: "mock-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	authServer := mockAuthServer(tokenResp)
	defer authServer.Close()

	authConfig := &client.AuthConfig{
		Login:    "test@example.com",
		Password: "testpass",
	}

	c := client.NewClient(authConfig, client.WithAuthURL(authServer.URL))

	ctx := context.Background()
	resp, err := c.Authenticate(ctx)

	if err != nil {
		t.Fatalf("Authenticate() failed = %v", err)
	}

	if resp.AccessToken != "mock-token" {
		t.Errorf("AccessToken = %s, want mock-token", resp.AccessToken)
	}

	if resp.TokenType != "Bearer" {
		t.Errorf("TokenType = %s, want Bearer", resp.TokenType)
	}

	if resp.ExpiresIn != 3600 {
		t.Errorf("ExpiresIn = %d, want 3600", resp.ExpiresIn)
	}
}

func TestMockClient_AuthenticateFailure(t *testing.T) {
	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "invalid_grant"}`))
	}))
	defer authServer.Close()

	authConfig := &client.AuthConfig{
		Login:    "test@example.com",
		Password: "wrongpass",
	}

	c := client.NewClient(authConfig, client.WithAuthURL(authServer.URL))

	ctx := context.Background()
	_, err := c.Authenticate(ctx)

	if err == nil {
		t.Fatal("Expected error for invalid credentials, got nil")
	}
}

func TestMockClient_GetCurrentUser(t *testing.T) {
	tokenResp := &client.TokenResponse{
		AccessToken: "mock-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	authServer := mockAuthServer(tokenResp)
	defer authServer.Close()

	graphqlServer := mockGraphQLServer(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"user": map[string]interface{}{
					"getCurrentUser": map[string]interface{}{
						"id":                     "test-user-id",
						"avatarUrl":              "/avatar.png",
						"login":                  "testuser",
						"firstName":              "Test",
						"middleName":             "",
						"lastName":               "User",
						"currentSchoolStudentId": "test-student-id",
						"__typename":             "User",
					},
					"__typename": "UserQueries",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer graphqlServer.Close()

	authConfig := &client.AuthConfig{
		Login:    "test@example.com",
		Password: "testpass",
	}

	c := client.NewClient(
		authConfig,
		client.WithAuthURL(authServer.URL),
		client.WithBaseURL(graphqlServer.URL),
	)

	// Pre-set the token to avoid re-authentication
	c.SetToken(tokenResp, time.Now().Add(time.Hour))

	ctx := context.Background()
	resp, err := c.GetCurrentUser(ctx)

	if err != nil {
		t.Fatalf("GetCurrentUser() failed = %v", err)
	}

	user := resp.User.GetCurrentUser
	if user.ID != "test-user-id" {
		t.Errorf("ID = %s, want test-user-id", user.ID)
	}

	if user.Login != "testuser" {
		t.Errorf("Login = %s, want testuser", user.Login)
	}

	if user.FirstName != "Test" {
		t.Errorf("FirstName = %s, want Test", user.FirstName)
	}

	if user.LastName != "User" {
		t.Errorf("LastName = %s, want User", user.LastName)
	}
}

func TestMockClient_GetUserNotifications(t *testing.T) {
	tokenResp := &client.TokenResponse{
		AccessToken: "mock-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	authServer := mockAuthServer(tokenResp)
	defer authServer.Close()

	graphqlServer := mockGraphQLServer(func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)

		response := map[string]interface{}{
			"data": map[string]interface{}{
				"s21Notification": map[string]interface{}{
					"getS21Notifications": map[string]interface{}{
						"notifications": []map[string]interface{}{
							{
								"id":               "notif-1",
								"relatedObjectType": "CALENDAR",
								"relatedObjectId":   "slot-123",
								"message":          "Someone registered for a review of the project <b>C3_s21_stringplus</b> by you on <b>2025.12.20, 18:00</b>",
								"time":             "2025-12-20T09:20:57.373Z",
								"wasRead":          false,
								"groupName":        "PROJECTS",
								"__typename":       "S21Notification",
							},
							{
								"id":               "notif-2",
								"relatedObjectType": "PROFILE",
								"relatedObjectId":   "0",
								"message":          "You achieved the new level <b>Level 12</b>!",
								"time":             "2025-12-20T18:00:48.583Z",
								"wasRead":          true,
								"groupName":        "PROFILE",
								"__typename":       "S21Notification",
							},
						},
						"totalCount": 430,
						"groupNames": []string{"SALE", "EVENTS", "PROJECTS", "AWARDS", "PROFILE", "PENALTIES"},
						"__typename": "S21NotificationReport",
					},
					"__typename": "S21NotificationQueries",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer graphqlServer.Close()

	authConfig := &client.AuthConfig{
		Login:    "test@example.com",
		Password: "testpass",
	}

	c := client.NewClient(
		authConfig,
		client.WithAuthURL(authServer.URL),
		client.WithBaseURL(graphqlServer.URL),
	)
	c.SetToken(tokenResp, time.Now().Add(time.Hour))

	ctx := context.Background()
	resp, err := c.GetUserNotifications(ctx, client.PagingInput{Offset: 0, Limit: 30})

	if err != nil {
		t.Fatalf("GetUserNotifications() failed = %v", err)
	}

	notifications := resp.S21Notification.GetS21Notifications.Notifications
	if len(notifications) != 2 {
		t.Fatalf("Got %d notifications, want 2", len(notifications))
	}

	if notifications[0].ID != "notif-1" {
		t.Errorf("First notification ID = %s, want notif-1", notifications[0].ID)
	}

	if notifications[0].RelatedObjectType != "CALENDAR" {
		t.Errorf("First notification type = %s, want CALENDAR", notifications[0].RelatedObjectType)
	}

	if resp.S21Notification.GetS21Notifications.TotalCount != 430 {
		t.Errorf("TotalCount = %d, want 430", resp.S21Notification.GetS21Notifications.TotalCount)
	}
}

func TestMockClient_DeleteEventSlot(t *testing.T) {
	tokenResp := &client.TokenResponse{
		AccessToken: "mock-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	authServer := mockAuthServer(tokenResp)
	defer authServer.Close()

	graphqlServer := mockGraphQLServer(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"student": map[string]interface{}{
					"deleteEventSlot": true,
					"__typename":     "StudentMutations",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer graphqlServer.Close()

	authConfig := &client.AuthConfig{
		Login:    "test@example.com",
		Password: "testpass",
	}

	c := client.NewClient(
		authConfig,
		client.WithAuthURL(authServer.URL),
		client.WithBaseURL(graphqlServer.URL),
	)
	c.SetToken(tokenResp, time.Now().Add(time.Hour))

	ctx := context.Background()
	resp, err := c.DeleteEventSlot(ctx, "slot-123")

	if err != nil {
		t.Fatalf("DeleteEventSlot() failed = %v", err)
	}

	if !resp.Student.DeleteEventSlot {
		t.Error("DeleteEventSlot returned false, want true")
	}
}

func TestMockClient_GraphQLError(t *testing.T) {
	tokenResp := &client.TokenResponse{
		AccessToken: "mock-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	authServer := mockAuthServer(tokenResp)
	defer authServer.Close()

	graphqlServer := mockGraphQLServer(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{},
			"errors": []map[string]interface{}{
				{
					"message": "Validation error",
					"path":    []interface{}{"user", "getCurrentUser"},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer graphqlServer.Close()

	authConfig := &client.AuthConfig{
		Login:    "test@example.com",
		Password: "testpass",
	}

	c := client.NewClient(
		authConfig,
		client.WithAuthURL(authServer.URL),
		client.WithBaseURL(graphqlServer.URL),
	)
	c.SetToken(tokenResp, time.Now().Add(time.Hour))

	ctx := context.Background()
	_, err := c.GetCurrentUser(ctx)

	if err == nil {
		t.Fatal("Expected error for GraphQL error, got nil")
	}

	if err.Error() == "" {
		t.Error("Error message is empty")
	}
}

func TestMockClient_GetStudentGraphTemplate(t *testing.T) {
	tokenResp := &client.TokenResponse{
		AccessToken: "mock-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	authServer := mockAuthServer(tokenResp)
	defer authServer.Close()

	graphqlServer := mockGraphQLServer(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"holyGraph": map[string]interface{}{
					"getStudentGraphTemplate": map[string]interface{}{
						"edges": []map[string]interface{}{
							{
								"id":     "edge-1",
								"source": "node-1",
								"target": "node-2",
								"data": map[string]interface{}{
									"sourceGap": 0.0,
									"targetGap": 0.0,
									"points":    []map[string]float64{},
								},
							},
						},
						"nodes": []map[string]interface{}{
							{
								"id":     "node-1",
								"label":  "Common Core",
								"items": []map[string]interface{}{
									{
										"id":     "item-1",
										"goal": map[string]interface{}{
											"projectId":          "project-1",
											"projectName":        "libft",
											"projectDescription": "Libft project",
											"projectPoints":      42,
											"goalExecutionType":  "DEFAULT",
											"isMandatory":        true,
										},
										"course": map[string]interface{}{
											"projectName": "Common Core",
										},
									},
								},
								"position": map[string]float64{"x": 100, "y": 200},
							},
							{
								"id":     "node-2",
								"label":  "Algorithmics",
								"items": []map[string]interface{}{
									{
										"id":     "item-2",
										"course": map[string]interface{}{
											"projectId":          "project-2",
											"projectName":        "CPP0_module",
											"projectDescription": "CPP project",
											"projectPoints":      100,
											"courseType":         "CPP",
											"isMandatory":        true,
										},
									},
								},
								"position": map[string]float64{"x": 300, "y": 200},
							},
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer graphqlServer.Close()

	authConfig := &client.AuthConfig{
		Login:    "test@example.com",
		Password: "testpass",
	}

	c := client.NewClient(
		authConfig,
		client.WithAuthURL(authServer.URL),
		client.WithBaseURL(graphqlServer.URL),
	)
	c.SetToken(tokenResp, time.Now().Add(time.Hour))

	ctx := context.Background()
	resp, err := c.GetStudentGraphTemplate(ctx, "student-123", nil)

	if err != nil {
		t.Fatalf("GetStudentGraphTemplate() failed = %v", err)
	}

	graph := resp.HolyGraph.GetStudentGraphTemplate
	if len(graph.Nodes) != 2 {
		t.Fatalf("Got %d nodes, want 2", len(graph.Nodes))
	}

	if len(graph.Edges) != 1 {
		t.Fatalf("Got %d edges, want 1", len(graph.Edges))
	}

	if graph.Nodes[0].Label != "Common Core" {
		t.Errorf("First node label = %s, want Common Core", graph.Nodes[0].Label)
	}

	if graph.Nodes[0].Items[0].Goal.ProjectName != "libft" {
		t.Errorf("First item ProjectName = %s, want libft", graph.Nodes[0].Items[0].Goal.ProjectName)
	}
}

func TestMockClient_GetCalendarEvents(t *testing.T) {
	tokenResp := &client.TokenResponse{
		AccessToken: "mock-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	authServer := mockAuthServer(tokenResp)
	defer authServer.Close()

	graphqlServer := mockGraphQLServer(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"calendarEventS21": map[string]interface{}{
					"getMyCalendarEvents": []map[string]interface{}{
						{
							"id":          "event-1",
							"start":       "2025-12-20T14:00:00Z",
							"end":         "2025-12-20T14:30:00Z",
							"description": "",
							"eventType":   "Я проверяю",
							"eventCode":   "student_check",
							"eventSlots": []map[string]interface{}{
								{
									"id":    "slot-1",
									"type":  "FREE_TIME",
									"start": "2025-12-20T14:00:00Z",
									"end":   "2025-12-20T14:30:00Z",
									"event": map[string]interface{}{
										"eventUserRole": "STUDENT",
									},
									"school": map[string]interface{}{
										"shortName": "21 Kazan",
									},
								},
							},
							"bookings":          []map[string]interface{}{},
							"exam":              nil,
							"studentCodeReview": nil,
							"activity":          nil,
							"goals":             []map[string]interface{}{},
							"penalty":           nil,
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer graphqlServer.Close()

	authConfig := &client.AuthConfig{
		Login:    "test@example.com",
		Password: "testpass",
	}

	c := client.NewClient(
		authConfig,
		client.WithAuthURL(authServer.URL),
		client.WithBaseURL(graphqlServer.URL),
	)
	c.SetToken(tokenResp, time.Now().Add(time.Hour))

	ctx := context.Background()
	resp, err := c.GetCalendarEvents(ctx, "2025-01-01T00:00:00Z", "2025-12-31T23:59:59Z")

	if err != nil {
		t.Fatalf("GetCalendarEvents() failed = %v", err)
	}

	events := resp.CalendarEventS21.GetMyCalendarEvents
	if len(events) != 1 {
		t.Fatalf("Got %d events, want 1", len(events))
	}

	if events[0].ID != "event-1" {
		t.Errorf("Event ID = %s, want event-1", events[0].ID)
	}

	if events[0].EventType != "Я проверяю" {
		t.Errorf("EventType = %s, want 'Я проверяю'", events[0].EventType)
	}
}

func TestMockClient_TokenExpiry(t *testing.T) {
	tokenResp := &client.TokenResponse{
		AccessToken: "mock-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	authServer := mockAuthServer(tokenResp)
	defer authServer.Close()

	authConfig := &client.AuthConfig{
		Login:    "test@example.com",
		Password: "testpass",
	}

	c := client.NewClient(
		authConfig,
		client.WithAuthURL(authServer.URL),
	)

	// Set token that's already expired
	c.SetToken(tokenResp, time.Now().Add(-time.Hour))

	graphqlServer := mockGraphQLServer(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"user": map[string]interface{}{
					"getCurrentUser": map[string]interface{}{
						"id":                     "test-user-id",
						"login":                  "testuser",
						"firstName":              "Test",
						"lastName":               "User",
						"currentSchoolStudentId": "test-student-id",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
	defer graphqlServer.Close()
	c.SetBaseURL(graphqlServer.URL)

	ctx := context.Background()
	_, err := c.GetCurrentUser(ctx)

	// Should auto-reauthenticate and succeed
	if err != nil {
		t.Fatalf("GetCurrentUser() failed = %v", err)
	}
}

func TestMockClient_Unauthorized(t *testing.T) {
	tokenResp := &client.TokenResponse{
		AccessToken: "bad-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	graphqlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer graphqlServer.Close()

	authConfig := &client.AuthConfig{
		Login:    "test@example.com",
		Password: "testpass",
	}

	c := client.NewClient(
		authConfig,
		client.WithBaseURL(graphqlServer.URL),
	)
	c.SetToken(tokenResp, time.Now().Add(time.Hour))

	ctx := context.Background()
	_, err := c.GetCurrentUser(ctx)

	if err == nil {
		t.Fatal("Expected error for unauthorized request, got nil")
	}
}
