package client

import "context"

// GraphQL Queries - exact copies from the API specification

// QueryGetUserNotifications fetches user notifications
const QueryGetUserNotifications = `query getUserNotifications($paging: PagingInput!) {
  s21Notification {
    getS21Notifications(paging: $paging) {
      notifications {
        id
        relatedObjectType
        relatedObjectId
        message
        time
        wasRead
        groupName
        __typename
      }
      totalCount
      groupNames
      __typename
    }
    __typename
  }
}`

// QueryGetUserNotificationsCount fetches notifications count
const QueryGetUserNotificationsCount = `query getUserNotificationsCount($wasReadIncluded: Boolean) {
  s21Notification {
    getS21NotificationsCount(wasReadIncluded: $wasReadIncluded)
    __typename
  }
}`

// QueryGetCurrentUser fetches current user info
const QueryGetCurrentUser = `query getCurrentUser {
  user {
    getCurrentUser {
      ...CurrentUser
      __typename
    }
    __typename
  }
}

fragment CurrentUser on User {
  id
  avatarUrl
  login
  firstName
  middleName
  lastName
  currentSchoolStudentId
  __typename
}`

// QueryProjectMapGetStudentStageGroups fetches student stage groups
const QueryProjectMapGetStudentStageGroups = `query ProjectMapGetStudentStageGroups($studentId: UUID!) {
  school21 {
    loadStudentStageGroups(studentId: $studentId) {
      studentId
      stageGroupStudentId
      stageGroupS21 {
        waveId
        waveName
        eduForm
        active
        __typename
      }
      __typename
    }
    __typename
  }
}`

// QueryProjectMapGetStudentGraphTemplate fetches student project graph
const QueryProjectMapGetStudentGraphTemplate = `query ProjectMapGetStudentGraphTemplate($studentId: UUID, $stageGroupId: Int) {
  holyGraph {
    getStudentGraphTemplate(studentId: $studentId, stageGroupId: $stageGroupId) {
      edges {
        id
        source
        target
        sourceHandle
        targetHandle
        data {
          sourceGap
          targetGap
          points {
            x
            y
            __typename
          }
          __typename
        }
        __typename
      }
      nodes {
        id
        label
        handles
        position {
          x
          y
          __typename
        }
        items {
          id
          code
          handles
          entityType
          entityId
          parentNodeCodes
          childrenNodeCodes
          skills {
            id
            name
            color
            textColor
            __typename
          }
          goal {
            projectId
            projectName
            projectDescription
            projectPoints
            goalExecutionType
            isMandatory
            __typename
          }
          course {
            projectId
            projectName
            projectDescription
            projectPoints
            courseType
            isMandatory
            __typename
          }
          __typename
        }
        __typename
      }
      __typename
    }
    __typename
  }
}`

// QueryCalendarGetMyReviews fetches upcoming reviews
const QueryCalendarGetMyReviews = `query calendarGetMyReviews($to: DateTime, $limit: Int) {
  student {
    getMyUpcomingBookings(to: $to, limit: $limit) {
      ...Review
      __typename
    }
    __typename
  }
}

fragment Review on CalendarBooking {
  id
  answerId
  eventSlot {
    id
    start
    end
    __typename
  }
  task {
    id
    title
    assignmentType
    goalId
    goalName
    studentTaskAdditionalAttributes {
      cookiesCount
      __typename
    }
    __typename
  }
  verifierUser {
    ...UserInBooking
    __typename
  }
  verifiableStudent {
    id
    user {
      ...UserInBooking
      __typename
    }
    __typename
  }
  team {
    ...ProjectTeamMembers
    __typename
  }
  bookingStatus
  isOnline
  vcLinkUrl
  __typename
}

fragment UserInBooking on User {
  id
  login
  avatarUrl
  userExperience {
    level {
      id
      range {
        levelCode
        __typename
      }
      __typename
    }
    __typename
  }
  __typename
}

fragment ProjectTeamMembers on ProjectTeamMembers {
  id
  teamLead {
    ...ProjectTeamMember
    __typename
  }
  members {
    ...ProjectTeamMember
    __typename
  }
  invitedUsers {
    ...ProjectTeamMember
    __typename
  }
  teamName
  teamStatus
  minTeamMemberCount
  maxTeamMemberCount
  __typename
}

fragment ProjectTeamMember on User {
  id
  avatarUrl
  login
  userExperience {
    level {
      id
      range {
        levelCode
        __typename
      }
      __typename
    }
    cookiesCount
    codeReviewPoints
    __typename
  }
  activeSchoolShortName
  __typename
}`

// QueryCalendarGetEvents fetches calendar events
const QueryCalendarGetEvents = `query calendarGetEvents($from: DateTime!, $to: DateTime!) {
  calendarEventS21 {
    getMyCalendarEvents(from: $from, to: $to) {
      ...CalendarEvent
      __typename
    }
    __typename
  }
}

fragment CalendarEvent on CalendarEvent {
  id
  start
  end
  description
  eventType
  eventCode
  eventSlots {
    id
    type
    start
    end
    event {
      eventUserRole
      __typename
    }
    school {
      shortName
      __typename
    }
    __typename
  }
  bookings {
    ...CalendarReviewBooking
    __typename
  }
  exam {
    ...CalendarEventExam
    __typename
  }
  studentCodeReview {
    studentGoalId
    __typename
  }
  activity {
    ...CalendarEventActivity
    studentFeedback {
      id
      rating
      comment
      isEmpty
      __typename
    }
    status
    activityType
    isMandatory
    isWaitListActive
    isVisible
    comments {
      type
      createTs
      comment
      __typename
    }
    organizers {
      id
      login
      __typename
    }
    __typename
  }
  goals {
    goalId
    goalName
    __typename
  }
  penalty {
    ...Penalty
    __typename
  }
  __typename
}

fragment CalendarReviewBooking on CalendarBooking {
  id
  answerId
  eventSlotId
  task {
    id
    goalId
    goalName
    studentTaskAdditionalAttributes {
      cookiesCount
      __typename
    }
    assignmentType
    __typename
  }
  eventSlot {
    id
    start
    end
    event {
      eventUserRole
      eventCode
      __typename
    }
    school {
      shortName
      __typename
    }
    __typename
  }
  verifierUser {
    ...CalendarReviewUser
    __typename
  }
  verifiableInfo {
    verifiableStudents {
      ...VerifiableStudentItem
      __typename
    }
    team {
      name
      __typename
    }
    __typename
  }
  bookingStatus
  isOnline
  vcLinkUrl
  additionalChecklist {
    filledChecklistId
    filledChecklistStatusRecordingEnum
    __typename
  }
  __typename
}

fragment CalendarReviewUser on User {
  id
  login
  __typename
}

fragment VerifiableStudentItem on VerifiableStudent {
  userId
  login
  avatarUrl
  levelCode
  isTeamLead
  cookiesCount
  codeReviewPoints
  school {
    shortName
    __typename
  }
  __typename
}

fragment CalendarEventExam on Exam {
  examId
  eventId
  beginDate
  endDate
  name
  location
  currentStudentsCount
  maxStudentCount
  updateDate
  goalId
  goalName
  isWaitListActive
  isInWaitList
  stopRegisterDate
  __typename
}

fragment CalendarEventActivity on ActivityEvent {
  activityEventId
  eventId
  name
  beginDate
  endDate
  isRegistered
  description
  currentStudentsCount
  maxStudentCount
  location
  updateDate
  isWaitListActive
  isInWaitList
  stopRegisterDate
  __typename
}

fragment Penalty on Penalty {
  comment
  id
  duration
  status
  startTime
  createTime
  penaltySlot {
    currentStudentsCount
    description
    duration
    startTime
    id
    endTime
    __typename
  }
  reasonId
  __typename
}`

// MutationCalendarDeleteEventSlot deletes an event slot
const MutationCalendarDeleteEventSlot = `mutation calendarDeleteEventSlot($eventSlotId: ID!) {
  student {
    deleteEventSlot(eventSlotId: $eventSlotId)
    __typename
  }
}`

// MutationCalendarChangeEventSlot changes an event slot
const MutationCalendarChangeEventSlot = `mutation calendarChangeEventSlot($id: ID!, $start: DateTime!, $end: DateTime!) {
  student {
    changeEventSlot(eventSlotId: $id, start: $start, end: $end) {
      ...CalendarEvent
      __typename
    }
    __typename
  }
}

fragment CalendarEvent on CalendarEvent {
  id
  start
  end
  description
  eventType
  eventCode
  eventSlots {
    id
    type
    start
    end
    event {
      eventUserRole
      __typename
    }
    school {
      shortName
      __typename
    }
    __typename
  }
  bookings {
    ...CalendarReviewBooking
    __typename
  }
  exam {
    ...CalendarEventExam
    __typename
  }
  studentCodeReview {
    studentGoalId
    __typename
  }
  activity {
    ...CalendarEventActivity
    studentFeedback {
      id
      rating
      comment
      isEmpty
      __typename
    }
    status
    activityType
    isMandatory
    isWaitListActive
    isVisible
    comments {
      type
      createTs
      comment
      __typename
    }
    organizers {
      id
      login
      __typename
    }
    __typename
  }
  goals {
    goalId
    goalName
    __typename
  }
  penalty {
    ...Penalty
    __typename
  }
  __typename
}

fragment CalendarReviewBooking on CalendarBooking {
  id
  answerId
  eventSlotId
  task {
    id
    goalId
    goalName
    studentTaskAdditionalAttributes {
      cookiesCount
      __typename
    }
    assignmentType
    __typename
  }
  eventSlot {
    id
    start
    end
    event {
      eventUserRole
      eventCode
      __typename
    }
    school {
      shortName
      __typename
    }
    __typename
  }
  verifierUser {
    ...CalendarReviewUser
    __typename
  }
  verifiableInfo {
    verifiableStudents {
      ...VerifiableStudentItem
      __typename
    }
    team {
      name
      __typename
    }
    __typename
  }
  bookingStatus
  isOnline
  vcLinkUrl
  additionalChecklist {
    filledChecklistId
    filledChecklistStatusRecordingEnum
    __typename
  }
  __typename
}

fragment CalendarReviewUser on User {
  id
  login
  __typename
}

fragment VerifiableStudentItem on VerifiableStudent {
  userId
  login
  avatarUrl
  levelCode
  isTeamLead
  cookiesCount
  codeReviewPoints
  school {
    shortName
    __typename
  }
  __typename
}

fragment CalendarEventExam on Exam {
  examId
  eventId
  beginDate
  endDate
  name
  location
  currentStudentsCount
  maxStudentCount
  updateDate
  goalId
  goalName
  isWaitListActive
  isInWaitList
  stopRegisterDate
  __typename
}

fragment CalendarEventActivity on ActivityEvent {
  activityEventId
  eventId
  name
  beginDate
  endDate
  isRegistered
  description
  currentStudentsCount
  maxStudentCount
  location
  updateDate
  isWaitListActive
  isInWaitList
  stopRegisterDate
  __typename
}

fragment Penalty on Penalty {
  comment
  id
  duration
  status
  startTime
  createTime
  penaltySlot {
    currentStudentsCount
    description
    duration
    startTime
    id
    endTime
    __typename
  }
  reasonId
  __typename
}`

// MutationCalendarAddEvent adds an event to timetable
const MutationCalendarAddEvent = `mutation calendarAddEvent($start: DateTime!, $end: DateTime!) {
  student {
    addEventToTimetable(start: $start, end: $end) {
      ...CalendarEvent
      __typename
    }
    __typename
  }
}

fragment CalendarEvent on CalendarEvent {
  id
  start
  end
  description
  eventType
  eventCode
  eventSlots {
    id
    type
    start
    end
    event {
      eventUserRole
      __typename
    }
    school {
      shortName
      __typename
    }
    __typename
  }
  bookings {
    ...CalendarReviewBooking
    __typename
  }
  exam {
    ...CalendarEventExam
    __typename
  }
  studentCodeReview {
    studentGoalId
    __typename
  }
  activity {
    ...CalendarEventActivity
    studentFeedback {
      id
      rating
      comment
      isEmpty
      __typename
    }
    status
    activityType
    isMandatory
    isWaitListActive
    isVisible
    comments {
      type
      createTs
      comment
      __typename
    }
    organizers {
      id
      login
      __typename
    }
    __typename
  }
  goals {
    goalId
    goalName
    __typename
  }
  penalty {
    ...Penalty
    __typename
  }
  __typename
}

fragment CalendarReviewBooking on CalendarBooking {
  id
  answerId
  eventSlotId
  task {
    id
    goalId
    goalName
    studentTaskAdditionalAttributes {
      cookiesCount
      __typename
    }
    assignmentType
    __typename
  }
  eventSlot {
    id
    start
    end
    event {
      eventUserRole
      eventCode
      __typename
    }
    school {
      shortName
      __typename
    }
    __typename
  }
  verifierUser {
    ...CalendarReviewUser
    __typename
  }
  verifiableInfo {
    verifiableStudents {
      ...VerifiableStudentItem
      __typename
    }
    team {
      name
      __typename
    }
    __typename
  }
  bookingStatus
  isOnline
  vcLinkUrl
  additionalChecklist {
    filledChecklistId
    filledChecklistStatusRecordingEnum
    __typename
  }
  __typename
}

fragment CalendarReviewUser on User {
  id
  login
  __typename
}

fragment VerifiableStudentItem on VerifiableStudent {
  userId
  login
  avatarUrl
  levelCode
  isTeamLead
  cookiesCount
  codeReviewPoints
  school {
    shortName
    __typename
  }
  __typename
}

fragment CalendarEventExam on Exam {
  examId
  eventId
  beginDate
  endDate
  name
  location
  currentStudentsCount
  maxStudentCount
  updateDate
  goalId
  goalName
  isWaitListActive
  isInWaitList
  stopRegisterDate
  __typename
}

fragment CalendarEventActivity on ActivityEvent {
  activityEventId
  eventId
  name
  beginDate
  endDate
  isRegistered
  description
  currentStudentsCount
  maxStudentCount
  location
  updateDate
  isWaitListActive
  isInWaitList
  stopRegisterDate
  __typename
}

fragment Penalty on Penalty {
  comment
  id
  duration
  status
  startTime
  createTime
  penaltySlot {
    currentStudentsCount
    description
    duration
    startTime
    id
    endTime
    __typename
  }
  reasonId
  __typename
}`

// PagingInput represents pagination parameters
type PagingInput struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// GetUserNotifications fetches user notifications
func (c *Client) GetUserNotifications(ctx context.Context, paging PagingInput) (*GetUserNotificationsData, error) {
	req := &GraphQLRequest{
		OperationName: "getUserNotifications",
		Query:         QueryGetUserNotifications,
		Variables: map[string]interface{}{
			"paging": paging,
		},
	}

	var resp GetUserNotificationsData
	if err := c.Do(ctx, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetUserNotificationsCount fetches notifications count
func (c *Client) GetUserNotificationsCount(ctx context.Context, wasReadIncluded bool) (*GetNotificationsCountData, error) {
	req := &GraphQLRequest{
		OperationName: "getUserNotificationsCount",
		Query:         QueryGetUserNotificationsCount,
		Variables: map[string]interface{}{
			"wasReadIncluded": wasReadIncluded,
		},
	}

	var resp GetNotificationsCountData
	if err := c.Do(ctx, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetCurrentUser fetches current user info
func (c *Client) GetCurrentUser(ctx context.Context) (*GetCurrentUserData, error) {
	req := &GraphQLRequest{
		OperationName: "getCurrentUser",
		Query:         QueryGetCurrentUser,
		Variables:     map[string]interface{}{},
	}

	var resp GetCurrentUserData
	if err := c.Do(ctx, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetStudentStageGroups fetches student stage groups
func (c *Client) GetStudentStageGroups(ctx context.Context, studentID string) (*LoadStudentStageGroupsData, error) {
	req := &GraphQLRequest{
		OperationName: "ProjectMapGetStudentStageGroups",
		Query:         QueryProjectMapGetStudentStageGroups,
		Variables: map[string]interface{}{
			"studentId": studentID,
		},
	}

	var resp LoadStudentStageGroupsData
	if err := c.Do(ctx, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetStudentGraphTemplate fetches student project graph template
func (c *Client) GetStudentGraphTemplate(ctx context.Context, studentID string, stageGroupID *int) (*GetStudentGraphTemplateData, error) {
	vars := map[string]interface{}{
		"studentId": studentID,
	}
	if stageGroupID != nil {
		vars["stageGroupId"] = *stageGroupID
	}

	req := &GraphQLRequest{
		OperationName: "ProjectMapGetStudentGraphTemplate",
		Query:         QueryProjectMapGetStudentGraphTemplate,
		Variables:     vars,
	}

	var resp GetStudentGraphTemplateData
	if err := c.Do(ctx, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetMyReviews fetches upcoming reviews
func (c *Client) GetMyReviews(ctx context.Context, to string, limit int) (*GetMyUpcomingBookingsData, error) {
	req := &GraphQLRequest{
		OperationName: "calendarGetMyReviews",
		Query:         QueryCalendarGetMyReviews,
		Variables: map[string]interface{}{
			"to":    to,
			"limit": limit,
		},
	}

	var resp GetMyUpcomingBookingsData
	if err := c.Do(ctx, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetCalendarEvents fetches calendar events
func (c *Client) GetCalendarEvents(ctx context.Context, from, to string) (*GetMyCalendarEventsData, error) {
	req := &GraphQLRequest{
		OperationName: "calendarGetEvents",
		Query:         QueryCalendarGetEvents,
		Variables: map[string]interface{}{
			"from": from,
			"to":   to,
		},
	}

	var resp GetMyCalendarEventsData
	if err := c.Do(ctx, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// DeleteEventSlot deletes an event slot
func (c *Client) DeleteEventSlot(ctx context.Context, eventSlotID string) (*DeleteEventSlotData, error) {
	req := &GraphQLRequest{
		OperationName: "calendarDeleteEventSlot",
		Query:         MutationCalendarDeleteEventSlot,
		Variables: map[string]interface{}{
			"eventSlotId": eventSlotID,
		},
	}

	var resp DeleteEventSlotData
	if err := c.Do(ctx, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ChangeEventSlot changes an event slot
func (c *Client) ChangeEventSlot(ctx context.Context, id, start, end string) (*ChangeEventSlotData, error) {
	req := &GraphQLRequest{
		OperationName: "calendarChangeEventSlot",
		Query:         MutationCalendarChangeEventSlot,
		Variables: map[string]interface{}{
			"id":    id,
			"start": start,
			"end":   end,
		},
	}

	var resp ChangeEventSlotData
	if err := c.Do(ctx, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// AddEventToTimetable adds an event to timetable
func (c *Client) AddEventToTimetable(ctx context.Context, start, end string) (*AddEventToTimetableData, error) {
	req := &GraphQLRequest{
		OperationName: "calendarAddEvent",
		Query:         MutationCalendarAddEvent,
		Variables: map[string]interface{}{
			"start": start,
			"end":   end,
		},
	}

	var resp AddEventToTimetableData
	if err := c.Do(ctx, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
