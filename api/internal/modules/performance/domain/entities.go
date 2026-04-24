package domain

// PersonalKPIs is the response for GET /performance/me (all roles)
type PersonalKPIs struct {
	UserID      uint    `json:"user_id"`
	Email       string  `json:"email"`
	Role        string  `json:"role"`
	HealthScore float64 `json:"health_score"`

	// DEV-focused metrics (populated for DEV; zeros for others)
	DeliveryRatePct  float64 `json:"delivery_rate_pct"`  // first job-done timestamp on time / tasks with due_date
	CodeQualityIndex float64 `json:"code_quality_index"` // AVG(ai_score) from submissions
	ReworkRatePct    float64 `json:"rework_rate_pct"`    // rejected-comment events / (job-done events + rejected-comment events)
	TimeAccuracyPct  float64 `json:"time_accuracy_pct"`  // 1 - |logged - estimated|/estimated
	SprintVelocitySP float64 `json:"sprint_velocity_sp"` // avg story points last 3 sprints
	VelocityTrend    string  `json:"velocity_trend"`     // "up" | "down" | "stable"
}

// TeamMemberKPI is one row in the team leaderboard (for /performance/team)
type TeamMemberKPI struct {
	UserID           uint    `json:"user_id"`
	Email            string  `json:"email"`
	Role             string  `json:"role"`
	HealthScore      float64 `json:"health_score"`
	DeliveryRatePct  float64 `json:"delivery_rate_pct"`
	CodeQualityIndex float64 `json:"code_quality_index"`
	ReworkRatePct    float64 `json:"rework_rate_pct"`
	TimeAccuracyPct  float64 `json:"time_accuracy_pct"`
	SprintVelocitySP float64 `json:"sprint_velocity_sp"`
	CompositeScore   float64 `json:"composite_score"` // 0-100 for ranking
}

// TeamKPIsResponse is the response for GET /performance/team (CEO + Product Owner)
type TeamKPIsResponse struct {
	Members []TeamMemberKPI `json:"members"`
}

// OverviewKPIs is the response for GET /performance/overview (CEO only)
type OverviewKPIs struct {
	EngineeringHealthIndex float64 `json:"engineering_health_index"` // weighted composite 0-100
	SprintSuccessRatePct   float64 `json:"sprint_success_rate_pct"`
	ProjectOnTrackRatePct  float64 `json:"project_on_track_rate_pct"`
	MilestoneHitRatePct    float64 `json:"milestone_hit_rate_pct"`
	CursorAdoptionScore    int     `json:"cursor_adoption_score"`   // from system_config
	TeamVelocityTrendPct   float64 `json:"team_velocity_trend_pct"` // sprint-over-sprint growth %
}

// ─── Discipline Dashboard ──────────────────────────────────────────────────────

// DisciplineJobDoneItem is one credited "Job Done" event in the queried range.
type DisciplineJobDoneItem struct {
	TaskID    string `json:"task_id"`
	TaskCode  string `json:"task_code,omitempty"`
	TaskTitle string `json:"task_title"`
	TaskType  string `json:"task_type"`
	DoneDate  string `json:"done_date"`  // YYYY-MM-DD (Bangkok)
	DoneTime  string `json:"done_time"`  // HH:MM (Bangkok)
	EventKind string `json:"event_kind"` // READY_FOR_TEST | PM_APPROVED_TEST | DEPLOYMENT_DEPLOYED
	// Actor = user who changed status (task_activity_events.actor_id or deployment reviewer)
	ActorID          uint   `json:"actor_id,omitempty"`
	ActorEmail       string `json:"actor_email,omitempty"`
	ActorDisplayName string `json:"actor_display_name,omitempty"`
	// When EventKind is DEPLOYMENT_DEPLOYED (Chief Engineer marked deployment deployed → task may advance to READY_FOR_UAT)
	DeploymentID    uint   `json:"deployment_id,omitempty"`
	DeploymentTitle string `json:"deployment_title,omitempty"`
	Branch          string `json:"branch,omitempty"`
	Environment     string `json:"environment,omitempty"`
}

// DisciplineReworkItem is one [REJECTED] comment credited to the task assignee in the range.
type DisciplineReworkItem struct {
	TaskID         string `json:"task_id"`
	TaskCode       string `json:"task_code,omitempty"`
	TaskTitle      string `json:"task_title"`
	TaskType       string `json:"task_type,omitempty"`
	EventDate      string `json:"event_date"` // YYYY-MM-DD (Bangkok)
	EventTime      string `json:"event_time"` // HH:MM (Bangkok)
	CommentSnippet string `json:"comment_snippet"`
	// Author = user who posted the [REJECTED] comment (task_comments.user_id)
	AuthorID          uint   `json:"author_id,omitempty"`
	AuthorEmail       string `json:"author_email,omitempty"`
	AuthorDisplayName string `json:"author_display_name,omitempty"`
}

// DisciplineUserDayStat holds one employee's activity metrics for a single date.
type DisciplineUserDayStat struct {
	Date                 string `json:"date"`                  // YYYY-MM-DD
	TasksClosed          int    `json:"tasks_closed"`          // tasks completed on this date
	Reworks              int    `json:"reworks"`               // [REJECTED] comments on tasks owned by user on this date
	LoggedMinutes        int    `json:"logged_minutes"`        // total minutes logged via time_logs
	HasDailyPulse        bool   `json:"has_daily_pulse"`       // whether a daily standup was submitted
	DeploymentsCompleted int    `json:"deployments_completed"` // deployment_requests marked DEPLOYED by this reviewer on this date
	// Attendance fields
	IsLate           bool   `json:"is_late"`
	EarlyCheckout    bool   `json:"early_checkout"`
	CheckInAt        string `json:"check_in_at,omitempty"`       // HH:MM ICT
	CheckOutAt       string `json:"check_out_at,omitempty"`      // HH:MM ICT
	AttendanceStatus string `json:"attendance_status,omitempty"` // present|late|absent|wfh|on_leave|holiday
	LeaveSession     string `json:"leave_session,omitempty"`     // AM|PM|FULL
}

// DisciplineUser aggregates one employee's discipline stats across the queried range.
type DisciplineUser struct {
	UserID                 uint                    `json:"user_id"`
	UserEmail              string                  `json:"user_email"`
	UserDisplayName        string                  `json:"user_display_name,omitempty"`
	UserAvatarURL          string                  `json:"user_avatar_url,omitempty"`
	Role                   string                  `json:"role"`
	MissedPulseCount       int                     `json:"missed_pulse_count"` // working days without a standup in range
	TotalTasksClosed       int                     `json:"total_tasks_closed"`
	TotalReworks           int                     `json:"total_reworks"`
	TotalLoggedHours       float64                 `json:"total_logged_hours"`
	TotalDeployments       int                     `json:"total_deployments"`         // total deployments completed (Chief Engineer)
	TotalLateDays          int                     `json:"total_late_days"`           // times checked in late
	TotalEarlyCheckoutDays int                     `json:"total_early_checkout_days"` // times left early
	Days                   []DisciplineUserDayStat `json:"days"`
	JobDoneItems           []DisciplineJobDoneItem `json:"job_done_items"`
	ReworkItems            []DisciplineReworkItem  `json:"rework_items"`
}

// DisciplineAttendanceRecord holds attendance check-in/out data for one day.
type DisciplineAttendanceRecord struct {
	CheckInAt     string `json:"check_in_at,omitempty"`  // HH:MM ICT
	CheckOutAt    string `json:"check_out_at,omitempty"` // HH:MM ICT
	IsLate        bool   `json:"is_late"`
	EarlyCheckout bool   `json:"early_checkout"`
	Status        string `json:"status"` // present | late | absent | wfh
	CheckInMethod string `json:"check_in_method,omitempty"`
}

// DisciplineResponse is the full payload for GET /performance/discipline.
type DisciplineResponse struct {
	FromDate string           `json:"from_date"` // YYYY-MM-DD
	ToDate   string           `json:"to_date"`   // YYYY-MM-DD
	Dates    []string         `json:"dates"`     // all calendar dates in range
	Users    []DisciplineUser `json:"users"`
}

// DisciplineStartDateResponse defines global discipline start date.
type DisciplineStartDateResponse struct {
	StartDate string `json:"start_date"` // YYYY-MM-DD
}

// ─── Discipline Day Detail ────────────────────────────────────────────────────

// DisciplineTimeLogEntry is one time-log entry for a specific day.
type DisciplineTimeLogEntry struct {
	TaskID      string  `json:"task_id"`
	TaskCode    string  `json:"task_code,omitempty"`
	TaskTitle   string  `json:"task_title"`
	Minutes     int     `json:"minutes"`
	Hours       float64 `json:"hours"`
	Description string  `json:"description,omitempty"`
	WorkType    string  `json:"work_type,omitempty"`
	IsTimer     bool    `json:"is_timer"`
}

// DisciplineCompletedTask is a task closed on a specific day.
type DisciplineCompletedTask struct {
	TaskID      string `json:"task_id"`
	TaskCode    string `json:"task_code,omitempty"`
	TaskTitle   string `json:"task_title"`
	StoryPoints int    `json:"story_points"`
	TaskType    string `json:"task_type"`
}

// DisciplineDeployedRequest is a deployment request marked DEPLOYED on a specific day.
type DisciplineDeployedRequest struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Branch      string `json:"branch"`
	Environment string `json:"environment"`
}

// DisciplineReworkEntry is a task that received a [REJECTED] comment on a specific day.
type DisciplineReworkEntry struct {
	TaskID          string `json:"task_id"`
	TaskCode        string `json:"task_code,omitempty"`
	TaskTitle       string `json:"task_title"`
	RejectedComment string `json:"rejected_comment"`
}

// DisciplineDayDetail is the drill-down payload for one user on one day.
type DisciplineDayDetail struct {
	UserID           uint                        `json:"user_id"`
	UserEmail        string                      `json:"user_email"`
	UserDisplayName  string                      `json:"user_display_name,omitempty"`
	Date             string                      `json:"date"` // YYYY-MM-DD
	HasDailyPulse    bool                        `json:"has_daily_pulse"`
	TotalLoggedMin   int                         `json:"total_logged_minutes"`
	Attendance       *DisciplineAttendanceRecord `json:"attendance,omitempty"` // nil if no record
	TimeLogs         []DisciplineTimeLogEntry    `json:"time_logs"`
	CompletedTasks   []DisciplineCompletedTask   `json:"completed_tasks"`
	Reworks          []DisciplineReworkEntry     `json:"reworks"`
	DeployedRequests []DisciplineDeployedRequest `json:"deployed_requests"` // deployments completed as reviewer on this day
}

// Usecase defines the performance business logic interface
type Usecase interface {
	GetPersonalKPIs(userID uint, role string) (*PersonalKPIs, error)
	GetTeamKPIs(requestingUserID uint, requestingRole string) (*TeamKPIsResponse, error)
	GetOverviewKPIs(requestingUserID uint, requestingRole string) (*OverviewKPIs, error)
	ResetReworkRate(devUserID uint, requesterRole string) error // CEO only: reset rework history for a dev
	GetDiscipline(from, to string) (*DisciplineResponse, error)
	GetDisciplineDayDetail(userID uint, date string) (*DisciplineDayDetail, error)
	GetDisciplineStartDate() (*DisciplineStartDateResponse, error)
	SetDisciplineStartDate(startDate string) (*DisciplineStartDateResponse, error)
}

// Repository defines the data access interface for performance aggregations
type Repository interface {
	// Raw aggregates for a single user (DEV metrics) — all tasks
	GetUserTaskDeliveryStats(userID uint) (tasksWithDue int, completedOnTime int, err error)
	GetUserSubmissionStats(userID uint) (avgScore float64, totalSubs int, failCount int, err error)
	GetUserReworkStats(userID uint) (jobDoneCount int, reworkCount int, err error)
	GetUserTimeAccuracy(userID uint) (avgAccuracyPct float64, sampleCount int, err error)
	GetUserSprintVelocity(userID uint, lastNSprints int) (avgStoryPoints float64, trend string, err error)

	// Product Owner–scoped: only tasks assigned by this Product Owner (assigned_by_id = pmID)
	GetDevUserIDsAssignedByPM(pmID uint) ([]uint, error)
	GetUserTaskDeliveryStatsForAssignedBy(devID, assignedByID uint) (tasksWithDue int, completedOnTime int, err error)
	GetUserSubmissionStatsForAssignedBy(devID, assignedByID uint) (avgScore float64, totalSubs int, failCount int, err error)
	GetUserReworkStatsForAssignedBy(devID, assignedByID uint) (jobDoneCount int, reworkCount int, err error)
	GetUserTimeAccuracyForAssignedBy(devID, assignedByID uint) (avgAccuracyPct float64, sampleCount int, err error)
	GetUserSprintVelocityForAssignedBy(devID, assignedByID uint, lastNSprints int) (avgStoryPoints float64, trend string, err error)

	// For team: same metrics for every user (DEV role)
	GetAllDevUserIDs() ([]uint, error)
	GetUserEmailAndRole(userID uint) (email string, role string, healthScore float64, err error)

	// Discipline dashboard — cross-module daily stats per user
	GetDisciplineStats(from, to string) (*DisciplineResponse, error)
	GetDisciplineDayDetail(userID uint, date string) (*DisciplineDayDetail, error)
	GetDisciplineStartDate() (*DisciplineStartDateResponse, error)
	SetDisciplineStartDate(startDate string) (*DisciplineStartDateResponse, error)

	// Company-wide (for overview)
	GetSprintSuccessRate() (ratePct float64, err error)
	GetMilestoneHitRate() (reached, missed int, err error)
	GetProjectOnTrackRate() (onTrackPct float64, err error)
	GetCursorAdoptionScore() (score int, err error)
	GetTeamVelocityTrend() (growthPct float64, err error)
	GetCompanyWideDeliveryAndQuality() (avgDeliveryPct, avgCodeQuality float64, err error)
	GetCompanyWideReworkAndTimeAccuracy() (avgReworkPct, avgTimeAccuracyPct float64, err error)
}
