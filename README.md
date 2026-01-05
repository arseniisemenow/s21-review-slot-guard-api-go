# THIS README.md IS VIBECODED!
# S21 API Client

Go client for the 21-school GraphQL API with full review slot management support.

## Features

- JWT authentication with auto-refresh
- Abstract GraphQL request builder
- Type-safe responses for all API operations
- Review slot management (get, add, update, remove)
- Docker support for consistent builds
- Integration tests with real credentials
- Unit tests with mock API server

## Installation

```bash
go get github.com/arseni/s21-api-client
```

## Usage

```go
import "github.com/arseni/s21-api-client/pkg/client"
```

### Authentication

```go
authConfig := &client.AuthConfig{
    Login:    os.Getenv("S21_LOGIN"),
    Password: os.Getenv("S21_PASSWORD"),
}

c := client.NewClient(authConfig)
ctx := context.Background()

// Authenticate is called automatically when needed
token, err := c.Authenticate(ctx)
```

### Review Slot Management

```go
// Get all review slots (available and booked)
from := time.Now()
to := time.Now().AddDate(0, 0, 7)
slots, bookings, err := c.GetReviewSlots(ctx, from, to)

// Get only available slots
available, err := c.GetAvailableReviewSlots(ctx, from, to)

// Add a new review slot
start := time.Date(2025, 1, 15, 14, 0, 0, 0, time.UTC)
end := time.Date(2025, 1, 15, 14, 30, 0, 0, time.UTC)
slots, err := c.AddReviewSlot(ctx, start, end)

// Update an existing slot
slotID := "slot-123"
newStart := time.Date(2025, 1, 15, 15, 0, 0, 0, time.UTC)
newEnd := time.Date(2025, 1, 15, 15, 30, 0, 0, time.UTC)
slot, err := c.UpdateReviewSlot(ctx, slotID, newStart, newEnd)

// Remove/delete a slot
err = c.RemoveReviewSlot(ctx, slotID)
// or
err = c.CancelReview(ctx, slotID)
```

### Get Current User

```go
user, err := c.GetCurrentUser(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Hello, %s %s!\n",
    user.User.GetCurrentUser.FirstName,
    user.User.GetCurrentUser.LastName)
```

### Get Notifications

```go
resp, err := c.GetUserNotifications(ctx, client.PagingInput{
    Offset: 0,
    Limit:  30,
})

notifications := resp.S21Notification.GetS21Notifications.Notifications
```

### Delete Event Slot

```go
resp, err := c.DeleteEventSlot(ctx, "event-slot-id")
if resp.Student.DeleteEventSlot {
    fmt.Println("Event slot deleted successfully")
}
```

### Get Student Projects

```go
user, _ := c.GetCurrentUser(ctx)
studentID := user.User.GetCurrentUser.CurrentSchoolStudentID

graph, err := c.GetStudentGraphTemplate(ctx, studentID, nil)
nodes := graph.HolyGraph.GetStudentGraphTemplate.Nodes
```

## CLI

Build and run the CLI tool:

```bash
make build
export S21_LOGIN="your@login.com"
export S21_PASSWORD="yourpassword"
./build/client user
```

### Available Commands

| Command | Description |
|---------|-------------|
| `user` | Get current user info |
| `notifications` | Get user notifications |
| `reviews` | Get upcoming reviews |
| `projects` | Get available projects |
| `calendar` | Get calendar events |
| `review-slots` | Manage review slots |

### Review Slots CLI

```bash
# Show available and booked slots (next 7 days)
./build/client review-slots get

# Show slots for next 30 days
./build/client review-slots get 30

# Add a new review slot
./build/client review-slots add '2025-01-15 14:00' '2025-01-15 14:30'

# Update an existing slot
./build/client review-slots update slot-123 '2025-01-15 15:00' '2025-01-15 15:30'

# Remove a slot
./build/client review-slots remove slot-123
# or using shorthand
./build/client review-slots rm slot-123
```

#### Accepted Datetime Formats

- `2025-01-15`
- `2025-01-15 14:00`
- `2025-01-15 14:00:00`
- `2025-01-15T14:00:00Z`
- `2025-01-15T14:00:00.000Z`

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `S21_LOGIN` | Yes | Your 21-school login |
| `S21_PASSWORD` | Yes | Your 21-school password |
| `S21_SCHOOL_ID` | No* | School ID (from browser, required for some operations) |
| `S21_USER_ROLE` | No* | User role (e.g., STUDENT) |
| `S21_EDU_PRODUCT_ID` | No* | Edu Product ID (from browser) |
| `S21_EDU_ORG_UNIT_ID` | No* | Edu Org Unit ID (from browser) |

*May be required depending on the API operation.

## Docker

```bash
# Build and start container
make docker-up

# Run tests in Docker
make docker-test

# Build project in Docker
make docker-build-project

# Open shell in container
make docker-shell

# Stop container
make docker-down
```

### Available Docker Commands

| Command | Description |
|---------|-------------|
| `make docker-build` | Build Docker image |
| `make docker-up` | Start Docker container |
| `make docker-down` | Stop Docker container |
| `make docker-shell` | Open shell in container |
| `make docker-build-project` | Build project in Docker |
| `make docker-test` | Run all tests in Docker |
| `make docker-test-mock` | Run mock tests in Docker |
| `make docker-test-real` | Run integration tests in Docker |
| `make docker-all` | Build and test in Docker |

## Testing

### Unit Tests (Mock API)

```bash
make test-mock
# or in Docker
make docker-test-mock
```

### Integration Tests (Real Credentials)

```bash
export S21_LOGIN="your@login.com"
export S21_PASSWORD="yourpassword"
make test-real
# or in Docker
docker compose exec s21-api-client make test-real
```

## Available API Operations

| Operation | Description |
|-----------|-------------|
| `Authenticate` | Get JWT access token |
| `GetCurrentUser` | Get current user info |
| `GetUserNotifications` | Get user notifications |
| `GetUserNotificationsCount` | Get notifications count |
| `GetStudentStageGroups` | Get student stage groups |
| `GetStudentGraphTemplate` | Get student project graph |
| `GetMyReviews` | Get upcoming reviews |
| `GetCalendarEvents` | Get calendar events |
| `GetReviewSlots` | Get review slots (available + booked) |
| `GetAvailableReviewSlots` | Get only available review slots |
| `GetBookedReviews` | Get only booked reviews |
| `AddReviewSlot` | Add a new review slot |
| `UpdateReviewSlot` | Update an existing review slot |
| `RemoveReviewSlot` | Delete a review slot |
| `CancelReview` | Cancel a review (alias for RemoveReviewSlot) |
| `DeleteEventSlot` | Delete an event slot |
| `ChangeEventSlot` | Change an event slot |
| `AddEventToTimetable` | Add event to timetable |

## Project Structure

```
.
├── cmd/
│   └── client/           # CLI application
├── pkg/
│   └── client/           # API client library
│       ├── client.go     # Core client with auth
│       ├── operations.go # GraphQL queries/mutations
│       ├── types.go      # Response types
│       └── review_slots.go # Review slot operations
├── tests/
│   ├── integration/      # Real API tests
│   └── unit/             # Mock API tests
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

## License

MIT
