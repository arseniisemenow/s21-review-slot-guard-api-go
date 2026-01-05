package client

// NotificationTypes
type S21Notification struct {
	ID               string `json:"id"`
	RelatedObjectType string `json:"relatedObjectType"`
	RelatedObjectID  string `json:"relatedObjectId"`
	Message          string `json:"message"`
	Time             string `json:"time"`
	WasRead          bool   `json:"wasRead"`
	GroupName        string `json:"groupName"`
	Typename         string `json:"__typename"`
}

type NotificationResponse struct {
	Notifications []S21Notification `json:"notifications"`
	TotalCount    int              `json:"totalCount"`
	GroupNames    []string         `json:"groupNames"`
	Typename      string           `json:"__typename"`
}

type S21NotificationQueries struct {
	GetS21Notifications NotificationResponse `json:"getS21Notifications"`
	Typename            string               `json:"__typename"`
}

type GetUserNotificationsData struct {
	S21Notification S21NotificationQueries `json:"s21Notification"`
}

type GetNotificationsCountData struct {
	S21Notification struct {
		GetS21NotificationsCount int `json:"getS21NotificationsCount"`
		Typename              string `json:"__typename"`
	} `json:"s21Notification"`
}

// User Types
type CurrentUser struct {
	ID                     string `json:"id"`
	AvatarURL              string `json:"avatarUrl"`
	Login                  string `json:"login"`
	FirstName              string `json:"firstName"`
	MiddleName             string `json:"middleName"`
	LastName               string `json:"lastName"`
	CurrentSchoolStudentID string `json:"currentSchoolStudentId"`
	Typename               string `json:"__typename"`
}

type UserQueries struct {
	GetCurrentUser CurrentUser `json:"getCurrentUser"`
	Typename      string      `json:"__typename"`
}

type GetCurrentUserData struct {
	User UserQueries `json:"user"`
}

// User Experience for bookings
type LevelRange struct {
	LevelCode int `json:"levelCode"`
	Typename  string `json:"__typename"`
}

type Level struct {
	ID     string      `json:"id"`
	Range  LevelRange  `json:"range"`
	Typename string    `json:"__typename"`
}

type UserExperience struct {
	Level Level `json:"level"`
	Typename string `json:"__typename"`
}

// User in booking (simplified version)
type UserInBooking struct {
	ID     string `json:"id"`
	Login  string `json:"login"`
	Typename string `json:"__typename"`
}

// User with experience
type UserWithExperience struct {
	ID               string         `json:"id"`
	AvatarURL        string         `json:"avatarUrl"`
	Login            string         `json:"login"`
	UserExperience   UserExperience `json:"userExperience"`
	Typename         string         `json:"__typename"`
}

// Project Team Members
type ProjectTeamMember struct {
	ID                  string         `json:"id"`
	AvatarURL           string         `json:"avatarUrl"`
	Login               string         `json:"login"`
	UserExperience      UserExperience `json:"userExperience"`
	CookiesCount        int            `json:"cookiesCount"`
	CodeReviewPoints    int            `json:"codeReviewPoints"`
	ActiveSchoolShortName string        `json:"activeSchoolShortName"`
	Typename            string         `json:"__typename"`
}

type ProjectTeamMembers struct {
	ID                 string              `json:"id"`
	TeamLead            ProjectTeamMember   `json:"teamLead"`
	Members             []ProjectTeamMember `json:"members"`
	InvitedUsers        []ProjectTeamMember `json:"invitedUsers"`
	TeamName            string              `json:"teamName"`
	TeamStatus          string              `json:"teamStatus"`
	MinTeamMemberCount  int                 `json:"minTeamMemberCount"`
	MaxTeamMemberCount  int                 `json:"maxTeamMemberCount"`
	Typename            string              `json:"__typename"`
}

// Verifiable Student in CalendarReviewBooking
type VerifiableStudentItem struct {
	UserID        string `json:"userId"`
	Login         string `json:"login"`
	AvatarURL     string `json:"avatarUrl"`
	LevelCode     int    `json:"levelCode"`
	IsTeamLead    *bool  `json:"isTeamLead"`
	CookiesCount  int    `json:"cookiesCount"`
	CodeReviewPoints int  `json:"codeReviewPoints"`
	School        struct {
		ShortName string `json:"shortName"`
		Typename string `json:"__typename"`
	} `json:"school"`
	Typename     string `json:"__typename"`
}

// VerifiableInfo
type VerifiableInfo struct {
	VerifiableStudents []VerifiableStudentItem `json:"verifiableStudents"`
	Team               *struct {
		Name      string `json:"name"`
		Typename string `json:"__typename"`
	} `json:"team"`
	Typename string `json:"__typename"`
}

// VerifiableStudent (for getMyReviews response)
type VerifiableStudent struct {
	ID   string `json:"id"`
	User UserWithExperience `json:"user"`
	Typename string `json:"__typename"`
}

// Calendar Types
type CalendarTimeSlot struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Start  string `json:"start"`
	End    string `json:"end"`
	Event  struct {
		EventUserRole string `json:"eventUserRole"`
		Typename      string `json:"__typename"`
	} `json:"event"`
	School struct {
		ShortName string `json:"shortName"`
		Typename string `json:"__typename"`
	} `json:"school"`
	Typename string `json:"__typename"`
}

type TaskInBooking struct {
	ID                            string                  `json:"id"`
	Title                         string                  `json:"title"`
	AssignmentType                string                  `json:"assignmentType"`
	GoalID                        string                  `json:"goalId"`
	GoalName                      string                  `json:"goalName"`
	StudentTaskAdditionalAttributes *StudentTaskAttributes `json:"studentTaskAdditionalAttributes"`
	Typename                      string                  `json:"__typename"`
}

type EventSlotInBooking struct {
	ID     string `json:"id"`
	Start  string `json:"start"`
	End    string `json:"end"`
	Event  struct {
		EventUserRole string `json:"eventUserRole"`
		EventCode     string `json:"eventCode"`
		Typename      string `json:"__typename"`
	} `json:"event"`
	School struct {
		ShortName string `json:"shortName"`
		Typename string `json:"__typename"`
	} `json:"school"`
	Typename string `json:"__typename"`
}

type CalendarBooking struct {
	ID              string          `json:"id"`
	AnswerID        string          `json:"answerId"`
	EventSlotID     string          `json:"eventSlotId"`
	Task            TaskInBooking   `json:"task"`
	EventSlot       EventSlotInBooking `json:"eventSlot"`
	VerifierUser    UserInBooking   `json:"verifierUser"`
	VerifiableStudent *VerifiableStudent `json:"verifiableStudent,omitempty"`
	Team            *ProjectTeamMembers `json:"team,omitempty"`
	VerifiableInfo  *VerifiableInfo  `json:"verifiableInfo,omitempty"`
	BookingStatus   string          `json:"bookingStatus"`
	IsOnline        bool            `json:"isOnline"`
	VCLinkURL       *string         `json:"vcLinkUrl"`
	Typename        string          `json:"__typename"`
}

type Exam struct {
	ExamID              string `json:"examId"`
	EventID             string `json:"eventId"`
	BeginDate           string `json:"beginDate"`
	EndDate             string `json:"endDate"`
	Name                string `json:"name"`
	Location            string `json:"location"`
	CurrentStudentsCount int    `json:"currentStudentsCount"`
	MaxStudentCount     int    `json:"maxStudentCount"`
	UpdateDate          string `json:"updateDate"`
	GoalID              string `json:"goalId"`
	GoalName            string `json:"goalName"`
	IsWaitListActive    bool   `json:"isWaitListActive"`
	IsInWaitList        bool   `json:"isInWaitList"`
	StopRegisterDate    string `json:"stopRegisterDate"`
	Typename            string `json:"__typename"`
}

type StudentFeedback struct {
	ID      string `json:"id"`
	Rating  *int   `json:"rating"`
	Comment string `json:"comment"`
	IsEmpty bool   `json:"isEmpty"`
	Typename string `json:"__typename"`
}

type Organizer struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Typename string `json:"__typename"`
}

type ActivityComment struct {
	Type     string `json:"type"`
	CreateTs string `json:"createTs"`
	Comment  string `json:"comment"`
	Typename string `json:"__typename"`
}

type ActivityEvent struct {
	ActivityEventID      string           `json:"activityEventId"`
	EventID              string           `json:"eventId"`
	Name                 string           `json:"name"`
	BeginDate            string           `json:"beginDate"`
	EndDate              string           `json:"endDate"`
	IsRegistered         bool             `json:"isRegistered"`
	Description          string           `json:"description"`
	CurrentStudentsCount int              `json:"currentStudentsCount"`
	MaxStudentCount      int              `json:"maxStudentCount"`
	Location             string           `json:"location"`
	UpdateDate           string           `json:"updateDate"`
	IsWaitListActive     bool             `json:"isWaitListActive"`
	IsInWaitList         bool             `json:"isInWaitList"`
	StopRegisterDate     string           `json:"stopRegisterDate"`
	StudentFeedback      *StudentFeedback `json:"studentFeedback,omitempty"`
	Status               *string          `json:"status,omitempty"`
	ActivityType         string           `json:"activityType,omitempty"`
	IsMandatory          bool             `json:"isMandatory,omitempty"`
	IsVisible            bool             `json:"isVisible,omitempty"`
	Comments             []ActivityComment `json:"comments,omitempty"`
	Organizers           []Organizer      `json:"organizers,omitempty"`
	Typename             string           `json:"__typename"`
}

type Penalty struct {
	Comment   string `json:"comment"`
	ID        string `json:"id"`
	Duration  int    `json:"duration"`
	Status    string `json:"status"`
	StartTime string `json:"startTime"`
	CreateTime string `json:"createTime"`
	PenaltySlot *struct {
		CurrentStudentsCount int    `json:"currentStudentsCount"`
		Description          string `json:"description"`
		Duration             int    `json:"duration"`
		StartTime            string `json:"startTime"`
		ID                   string `json:"id"`
		EndTime              string `json:"endTime"`
		Typename             string `json:"__typename"`
	} `json:"penaltySlot"`
	ReasonID string `json:"reasonId"`
	Typename string `json:"__typename"`
}

type CalendarGoal struct {
	GoalID   string `json:"goalId"`
	GoalName string `json:"goalName"`
	Typename string `json:"__typename"`
}

type CalendarEvent struct {
	ID                string             `json:"id"`
	Start             string             `json:"start"`
	End               string             `json:"end"`
	Description       string             `json:"description"`
	EventType         string             `json:"eventType"`
	EventCode         string             `json:"eventCode"`
	EventSlots        []CalendarTimeSlot `json:"eventSlots"`
	Bookings          []CalendarBooking  `json:"bookings"`
	Exam              *Exam              `json:"exam"`
	StudentCodeReview *struct {
		StudentGoalID string `json:"studentGoalId"`
		Typename      string `json:"__typename"`
	} `json:"studentCodeReview"`
	Activity   *ActivityEvent `json:"activity"`
	Goals      []CalendarGoal  `json:"goals"`
	Penalty    *Penalty        `json:"penalty"`
	Typename   string          `json:"__typename"`
}

type GetMyUpcomingBookingsData struct {
	Student struct {
		GetMyUpcomingBookings []CalendarBooking `json:"getMyUpcomingBookings"`
		Typename              string            `json:"__typename"`
	} `json:"student"`
}

type GetMyCalendarEventsData struct {
	CalendarEventS21 struct {
		GetMyCalendarEvents []CalendarEvent `json:"getMyCalendarEvents"`
		Typename            string          `json:"__typename"`
	} `json:"calendarEventS21"`
}

type DeleteEventSlotData struct {
	Student struct {
		DeleteEventSlot bool   `json:"deleteEventSlot"`
		Typename       string `json:"__typename"`
	} `json:"student"`
}

type ChangeEventSlotData struct {
	Student struct {
		ChangeEventSlot CalendarEvent `json:"changeEventSlot"`
		Typename        string        `json:"__typename"`
	} `json:"student"`
}

type AddEventToTimetableData struct {
	Student struct {
		AddEventToTimetable []CalendarEvent `json:"addEventToTimetable"`
		Typename            string          `json:"__typename"`
	} `json:"student"`
}

// Project Map Types
type StageGroupS21 struct {
	WaveID   int    `json:"waveId"`
	WaveName string `json:"waveName"`
	EduForm  string `json:"eduForm"`
	Active   bool   `json:"active"`
	Typename string `json:"__typename"`
}

type StageGroupS21Student struct {
	StudentID          string        `json:"studentId"`
	StageGroupStudentID string       `json:"stageGroupStudentId"`
	StageGroupS21      StageGroupS21 `json:"stageGroupS21"`
	Typename           string        `json:"__typename"`
}

type LoadStudentStageGroupsData struct {
	School21 struct {
		LoadStudentStageGroups []StageGroupS21Student `json:"loadStudentStageGroups"`
		Typename               string                 `json:"__typename"`
	} `json:"school21"`
}

type GraphPoint struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Typename string  `json:"__typename"`
}

type EdgeData struct {
	SourceGap float64      `json:"sourceGap"`
	TargetGap float64      `json:"targetGap"`
	Points    []GraphPoint `json:"points"`
	Typename  string       `json:"__typename"`
}

type GraphEdge struct {
	ID         string   `json:"id"`
	Source     string   `json:"source"`
	Target     string   `json:"target"`
	SourceHandle string `json:"sourceHandle"`
	TargetHandle string `json:"targetHandle"`
	Data       EdgeData `json:"data"`
	Typename   string `json:"__typename"`
}

type Skill struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	TextColor string `json:"textColor"`
	Typename   string `json:"__typename"`
}

type GoalItem struct {
	ProjectID          interface{} `json:"projectId"`
	ProjectName        string `json:"projectName"`
	ProjectDescription string `json:"projectDescription"`
	ProjectPoints      int    `json:"projectPoints"`
	GoalExecutionType  string `json:"goalExecutionType"`
	IsMandatory        bool   `json:"isMandatory"`
	Typename           string `json:"__typename"`
}

type CourseItem struct {
	ProjectID          interface{} `json:"projectId"`
	ProjectName        string `json:"projectName"`
	ProjectDescription string `json:"projectDescription"`
	ProjectPoints      int    `json:"projectPoints"`
	CourseType         string `json:"courseType"`
	IsMandatory        bool   `json:"isMandatory"`
	Typename           string `json:"__typename"`
}

type GraphNodeItem struct {
	ID               string      `json:"id"`
	Code             string      `json:"code"`
	Handles          interface{} `json:"handles"`
	EntityType       string      `json:"entityType"`
	EntityID         interface{} `json:"entityId"`
	ParentNodeCodes  []string    `json:"parentNodeCodes"`
	ChildrenNodeCodes []string   `json:"childrenNodeCodes"`
	Skills           []Skill     `json:"skills"`
	Goal             *GoalItem   `json:"goal"`
	Course           *CourseItem `json:"course"`
	Typename         string      `json:"__typename"`
}

type NodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Typename string `json:"__typename"`
}

type GraphNode struct {
	ID       string         `json:"id"`
	Label    string         `json:"label"`
	Handles  interface{}    `json:"handles"`
	Position NodePosition   `json:"position"`
	Items    []GraphNodeItem `json:"items"`
	Typename string         `json:"__typename"`
}

type StudentGraphTemplate struct {
	Edges    []GraphEdge `json:"edges"`
	Nodes    []GraphNode `json:"nodes"`
	Typename string      `json:"__typename"`
}

type GetStudentGraphTemplateData struct {
	HolyGraph struct {
		GetStudentGraphTemplate StudentGraphTemplate `json:"getStudentGraphTemplate"`
		Typename                 string              `json:"__typename"`
	} `json:"holyGraph"`
}

type StudentTaskAttributes struct {
	CookiesCount int `json:"cookiesCount"`
	Typename    string `json:"__typename"`
}
