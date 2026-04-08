package repository

import (
	"time"

	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	perfDomain "github.com/portnd/the-sentinel-core/internal/modules/performance/domain"
	sentinelDomain "github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
	"gorm.io/gorm"
)

// ict is used to format timestamps to Bangkok time.
var ict *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		loc = time.FixedZone("ICT", 7*3600)
	}
	ict = loc
}

type postgresRepo struct {
	db *gorm.DB
}

// NewPostgresRepository returns a performance repository that queries existing tables
func NewPostgresRepository(db *gorm.DB) perfDomain.Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) GetUserTaskDeliveryStats(userID uint) (tasksWithDue int, completedOnTime int, err error) {
	var withDue int64
	err = r.db.Model(&sentinelDomain.Task{}).
		Where("assigned_to = ? AND due_at IS NOT NULL", userID).
		Count(&withDue).Error
	if err != nil {
		return 0, 0, err
	}
	tasksWithDue = int(withDue)

	var onTime int64
	err = r.db.Model(&sentinelDomain.Task{}).
		Where("assigned_to = ? AND due_at IS NOT NULL AND status = ? AND completed_at IS NOT NULL AND completed_at <= due_at", userID, "COMPLETED").
		Count(&onTime).Error
	if err != nil {
		return tasksWithDue, 0, err
	}
	completedOnTime = int(onTime)
	return tasksWithDue, completedOnTime, nil
}

func (r *postgresRepo) GetUserSubmissionStats(userID uint) (avgScore float64, totalSubs int, failCount int, err error) {
	// Fetch rework_reset_at cutoff for this user (NULL means no reset ever done)
	var resetAt *time.Time
	type resetRow struct{ ReworkResetAt *time.Time }
	var rr resetRow
	r.db.Raw("SELECT rework_reset_at FROM users WHERE id = ?", userID).Scan(&rr)
	resetAt = rr.ReworkResetAt

	// totalSubs = tasks assigned to this dev that have at least one submission
	// failCount = tasks assigned to this dev that have at least one [REJECTED] comment (rework events)
	var total int64
	err = r.db.Model(&sentinelDomain.Task{}).
		Where("assigned_to = ? AND id IN (SELECT DISTINCT task_id FROM submissions WHERE dev_id = ?)", userID, userID).
		Count(&total).Error
	if err != nil {
		return 0, 0, 0, err
	}
	if total == 0 {
		return 0, 0, 0, nil
	}

	var fails int64
	if resetAt != nil {
		err = r.db.Raw(`
			SELECT COUNT(DISTINCT t.id)
			FROM tasks t
			WHERE t.assigned_to = ?
			  AND EXISTS (
			      SELECT 1 FROM task_comments tc
			      WHERE tc.task_id = t.id
			        AND tc.content LIKE '[REJECTED]%'
			        AND tc.created_at > ?
			  )
		`, userID, resetAt).Scan(&fails).Error
	} else {
		err = r.db.Raw(`
			SELECT COUNT(DISTINCT t.id)
			FROM tasks t
			WHERE t.assigned_to = ?
			  AND EXISTS (
			      SELECT 1 FROM task_comments tc
			      WHERE tc.task_id = t.id
			        AND tc.content LIKE '[REJECTED]%'
			  )
		`, userID).Scan(&fails).Error
	}
	if err != nil {
		return 0, int(total), 0, err
	}

	// avgScore repurposed as approval rate (0–100) for the composite score
	avgScore = (1.0 - float64(fails)/float64(total)) * 100
	return avgScore, int(total), int(fails), nil
}

func (r *postgresRepo) GetUserTimeAccuracy(userID uint) (avgAccuracyPct float64, sampleCount int, err error) {
	// For each task assigned to user that has estimated_minutes > 0 and has time_logs, compute
	// accuracy = 1 - |sum(logged_mins) - estimated| / estimated, capped 0..1. Then average.
	type row struct {
		Accuracy float64
	}
	var rows []row
	// MEETING time is excluded — it doesn't represent dev productivity
	err = r.db.Raw(`
		WITH task_totals AS (
			SELECT t.id, t.estimated_minutes AS est,
				COALESCE(SUM(tl.minutes), 0)::int AS logged
			FROM tasks t
			LEFT JOIN time_logs tl ON tl.task_id = t.id
				AND tl.user_id = ?
				AND COALESCE(tl.work_type, 'DEV') != 'MEETING'
			WHERE t.assigned_to = ? AND t.estimated_minutes > 0
			GROUP BY t.id, t.estimated_minutes
		)
		SELECT LEAST(1.0, GREATEST(0.0, 1.0 - ABS(logged - est)::float / NULLIF(est, 0))) AS accuracy
		FROM task_totals
		WHERE est > 0
	`, userID, userID).Scan(&rows).Error
	if err != nil {
		return 0, 0, err
	}
	if len(rows) == 0 {
		return 0, 0, nil
	}
	var sum float64
	for _, r := range rows {
		sum += r.Accuracy
	}
	avgAccuracyPct = sum / float64(len(rows)) * 100
	sampleCount = len(rows)
	return avgAccuracyPct, sampleCount, nil
}

func (r *postgresRepo) GetUserSprintVelocity(userID uint, lastNSprints int) (avgStoryPoints float64, trend string, err error) {
	// Last N completed sprints (by end_date), story points completed by this user in those sprints
	type row struct {
		SprintOrder int
		SP          int
	}
	var rows []row
	err = r.db.Raw(`
		SELECT ROW_NUMBER() OVER (ORDER BY s.end_date DESC NULLS LAST) AS sprint_order,
			COALESCE(SUM(CASE WHEN t.status = 'COMPLETED' THEN t.story_points ELSE 0 END), 0)::int AS sp
		FROM sprints s
		LEFT JOIN tasks t ON t.sprint_id = s.id AND t.assigned_to = ?
		WHERE s.status = 'COMPLETED'
		GROUP BY s.id, s.end_date
		ORDER BY s.end_date DESC NULLS LAST
		LIMIT ?
	`, userID, lastNSprints).Scan(&rows).Error
	if err != nil {
		return 0, "stable", err
	}
	if len(rows) == 0 {
		return 0, "stable", nil
	}
	var total int
	for _, r := range rows {
		total += r.SP
	}
	avgStoryPoints = float64(total) / float64(len(rows))

	if len(rows) >= 2 {
		first, last := rows[len(rows)-1].SP, rows[0].SP
		if last > first {
			trend = "up"
		} else if last < first {
			trend = "down"
		} else {
			trend = "stable"
		}
	} else {
		trend = "stable"
	}
	return avgStoryPoints, trend, nil
}

// GetDevUserIDsAssignedByPM returns distinct engineer user IDs who have at least one task assigned by this Product Owner (assigned_by_id).
func (r *postgresRepo) GetDevUserIDsAssignedByPM(pmID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&sentinelDomain.Task{}).
		Where("assigned_by_id = ? AND assigned_to IS NOT NULL", pmID).
		Distinct("assigned_to").
		Pluck("assigned_to", &ids).Error
	return ids, err
}

func (r *postgresRepo) GetUserTaskDeliveryStatsForAssignedBy(devID, assignedByID uint) (tasksWithDue int, completedOnTime int, err error) {
	var withDue int64
	err = r.db.Model(&sentinelDomain.Task{}).
		Where("assigned_to = ? AND assigned_by_id = ? AND due_at IS NOT NULL", devID, assignedByID).
		Count(&withDue).Error
	if err != nil {
		return 0, 0, err
	}
	tasksWithDue = int(withDue)
	var onTime int64
	err = r.db.Model(&sentinelDomain.Task{}).
		Where("assigned_to = ? AND assigned_by_id = ? AND due_at IS NOT NULL AND status = ? AND completed_at IS NOT NULL AND completed_at <= due_at", devID, assignedByID, "COMPLETED").
		Count(&onTime).Error
	if err != nil {
		return tasksWithDue, 0, err
	}
	completedOnTime = int(onTime)
	return tasksWithDue, completedOnTime, nil
}

func (r *postgresRepo) GetUserSubmissionStatsForAssignedBy(devID, assignedByID uint) (avgScore float64, totalSubs int, failCount int, err error) {
	// Fetch rework_reset_at cutoff for this dev
	var resetAt *time.Time
	type resetRow struct{ ReworkResetAt *time.Time }
	var rr resetRow
	r.db.Raw("SELECT rework_reset_at FROM users WHERE id = ?", devID).Scan(&rr)
	resetAt = rr.ReworkResetAt

	// totalSubs = tasks assigned to devID by assignedByID that have at least one submission
	// failCount = those tasks that have at least one [REJECTED] comment (rework events)
	var total int64
	err = r.db.Model(&sentinelDomain.Task{}).
		Where("assigned_to = ? AND assigned_by_id = ? AND id IN (SELECT DISTINCT task_id FROM submissions WHERE dev_id = ?)", devID, assignedByID, devID).
		Count(&total).Error
	if err != nil {
		return 0, 0, 0, err
	}
	if total == 0 {
		return 0, 0, 0, nil
	}

	var fails int64
	if resetAt != nil {
		err = r.db.Raw(`
			SELECT COUNT(DISTINCT t.id)
			FROM tasks t
			WHERE t.assigned_to = ?
			  AND t.assigned_by_id = ?
			  AND EXISTS (
			      SELECT 1 FROM task_comments tc
			      WHERE tc.task_id = t.id
			        AND tc.content LIKE '[REJECTED]%'
			        AND tc.created_at > ?
			  )
		`, devID, assignedByID, resetAt).Scan(&fails).Error
	} else {
		err = r.db.Raw(`
			SELECT COUNT(DISTINCT t.id)
			FROM tasks t
			WHERE t.assigned_to = ?
			  AND t.assigned_by_id = ?
			  AND EXISTS (
			      SELECT 1 FROM task_comments tc
			      WHERE tc.task_id = t.id
			        AND tc.content LIKE '[REJECTED]%'
			  )
		`, devID, assignedByID).Scan(&fails).Error
	}
	if err != nil {
		return 0, int(total), 0, err
	}

	avgScore = (1.0 - float64(fails)/float64(total)) * 100
	return avgScore, int(total), int(fails), nil
}

func (r *postgresRepo) GetUserTimeAccuracyForAssignedBy(devID, assignedByID uint) (avgAccuracyPct float64, sampleCount int, err error) {
	type row struct {
		Accuracy float64
	}
	var rows []row
	// MEETING time excluded from dev productivity metrics
	err = r.db.Raw(`
		WITH task_totals AS (
			SELECT t.id, t.estimated_minutes AS est,
				COALESCE(SUM(tl.minutes), 0)::int AS logged
			FROM tasks t
			LEFT JOIN time_logs tl ON tl.task_id = t.id
				AND tl.user_id = ?
				AND COALESCE(tl.work_type, 'DEV') != 'MEETING'
			WHERE t.assigned_to = ? AND t.assigned_by_id = ? AND t.estimated_minutes > 0
			GROUP BY t.id, t.estimated_minutes
		)
		SELECT LEAST(1.0, GREATEST(0.0, 1.0 - ABS(logged - est)::float / NULLIF(est, 0))) AS accuracy
		FROM task_totals
		WHERE est > 0
	`, devID, devID, assignedByID).Scan(&rows).Error
	if err != nil {
		return 0, 0, err
	}
	if len(rows) == 0 {
		return 0, 0, nil
	}
	var sum float64
	for _, row := range rows {
		sum += row.Accuracy
	}
	avgAccuracyPct = sum / float64(len(rows)) * 100
	sampleCount = len(rows)
	return avgAccuracyPct, sampleCount, nil
}

func (r *postgresRepo) GetUserSprintVelocityForAssignedBy(devID, assignedByID uint, lastNSprints int) (avgStoryPoints float64, trend string, err error) {
	type row struct {
		SprintOrder int
		SP          int
	}
	var rows []row
	err = r.db.Raw(`
		SELECT ROW_NUMBER() OVER (ORDER BY s.end_date DESC NULLS LAST) AS sprint_order,
			COALESCE(SUM(CASE WHEN t.status = 'COMPLETED' THEN t.story_points ELSE 0 END), 0)::int AS sp
		FROM sprints s
		LEFT JOIN tasks t ON t.sprint_id = s.id AND t.assigned_to = ? AND t.assigned_by_id = ?
		WHERE s.status = 'COMPLETED'
		GROUP BY s.id, s.end_date
		ORDER BY s.end_date DESC NULLS LAST
		LIMIT ?
	`, devID, assignedByID, lastNSprints).Scan(&rows).Error
	if err != nil {
		return 0, "stable", err
	}
	if len(rows) == 0 {
		return 0, "stable", nil
	}
	var total int
	for _, row := range rows {
		total += row.SP
	}
	avgStoryPoints = float64(total) / float64(len(rows))
	if len(rows) >= 2 {
		first, last := rows[len(rows)-1].SP, rows[0].SP
		if last > first {
			trend = "up"
		} else if last < first {
			trend = "down"
		} else {
			trend = "stable"
		}
	} else {
		trend = "stable"
	}
	return avgStoryPoints, trend, nil
}

func (r *postgresRepo) GetAllDevUserIDs() ([]uint, error) {
	var ids []uint
	err := r.db.Model(&authDomain.User{}).Where("role IN ?", []string{authDomain.RoleEngineer, authDomain.RoleChiefEngineer}).Pluck("id", &ids).Error
	return ids, err
}

func (r *postgresRepo) GetUserEmailAndRole(userID uint) (email string, role string, healthScore float64, err error) {
	var u authDomain.User
	err = r.db.Select("email", "role", "health_score").First(&u, "id = ?", userID).Error
	if err != nil {
		return "", "", 0, err
	}
	return u.Email, u.Role, u.HealthScore, nil
}

func (r *postgresRepo) GetSprintSuccessRate() (ratePct float64, err error) {
	// Sprints with >= 80% of planned story points completed
	type agg struct {
		Total    int64
		Success  int64
	}
	var a agg
	err = r.db.Raw(`
		WITH sprint_sp AS (
			SELECT s.id,
				COALESCE(SUM(t.story_points), 0) AS planned,
				COALESCE(SUM(CASE WHEN t.status = 'COMPLETED' THEN t.story_points ELSE 0 END), 0) AS completed
			FROM sprints s
			LEFT JOIN tasks t ON t.sprint_id = s.id
			WHERE s.status = 'COMPLETED'
			GROUP BY s.id
		)
		SELECT COUNT(*) AS total,
			SUM(CASE WHEN planned = 0 OR (completed::float / NULLIF(planned, 0)) >= 0.8 THEN 1 ELSE 0 END) AS success
		FROM sprint_sp
	`).Scan(&a).Error
	if err != nil {
		return 0, err
	}
	if a.Total == 0 {
		return 0, nil
	}
	ratePct = float64(a.Success) / float64(a.Total) * 100
	return ratePct, nil
}

func (r *postgresRepo) GetMilestoneHitRate() (reached, missed int, err error) {
	var rCount, mCount int64
	err = r.db.Model(&sentinelDomain.Milestone{}).
		Where("status = ?", "REACHED").Count(&rCount).Error
	if err != nil {
		return 0, 0, err
	}
	err = r.db.Model(&sentinelDomain.Milestone{}).
		Where("status = ?", "MISSED").Count(&mCount).Error
	if err != nil {
		return 0, 0, err
	}
	return int(rCount), int(mCount), nil
}

func (r *postgresRepo) GetProjectOnTrackRate() (onTrackPct float64, err error) {
	// Active projects where >= 70% of tasks (with due_at) are on schedule (completed on time or not yet due)
	type agg struct {
		Total   int64
		OnTrack int64
	}
	var a agg
	err = r.db.Raw(`
		WITH project_stats AS (
			SELECT p.id,
				COUNT(t.id) FILTER (WHERE t.due_at IS NOT NULL) AS with_due,
				COUNT(t.id) FILTER (WHERE t.due_at IS NOT NULL AND (t.status != 'COMPLETED' OR t.completed_at <= t.due_at)) AS on_schedule
			FROM projects p
			LEFT JOIN tasks t ON t.project_id = p.id AND t.parent_id IS NULL
			WHERE p.status = 'ACTIVE'
			GROUP BY p.id
		)
		SELECT COUNT(*) AS total,
			SUM(CASE WHEN with_due = 0 OR (on_schedule::float / NULLIF(with_due, 0)) >= 0.7 THEN 1 ELSE 0 END) AS on_track
		FROM project_stats
	`).Scan(&a).Error
	if err != nil {
		return 0, err
	}
	if a.Total == 0 {
		return 0, nil
	}
	onTrackPct = float64(a.OnTrack) / float64(a.Total) * 100
	return onTrackPct, nil
}

func (r *postgresRepo) GetCursorAdoptionScore() (score int, err error) {
	var s struct{ CursorAssistance int }
	err = r.db.Table("system_configs").Select("cursor_assistance").First(&s).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 80, nil
		}
		return 0, err
	}
	return s.CursorAssistance, nil
}

func (r *postgresRepo) GetTeamVelocityTrend() (growthPct float64, err error) {
	type row struct {
		CompletedSP int
	}
	var rows []row
	err = r.db.Raw(`
		SELECT COALESCE(SUM(CASE WHEN t.status = 'COMPLETED' THEN t.story_points ELSE 0 END), 0)::int AS completed_sp
		FROM sprints s
		LEFT JOIN tasks t ON t.sprint_id = s.id
		WHERE s.status = 'COMPLETED'
		GROUP BY s.id, s.end_date
		ORDER BY s.end_date DESC NULLS LAST
		LIMIT 4
	`).Scan(&rows).Error
	if err != nil || len(rows) < 2 {
		return 0, err
	}
	recent, prev := rows[0].CompletedSP, rows[1].CompletedSP
	if prev == 0 {
		if recent > 0 {
			return 100, nil
		}
		return 0, nil
	}
	growthPct = (float64(recent)-float64(prev))/float64(prev) * 100
	return growthPct, nil
}

func (r *postgresRepo) GetCompanyWideDeliveryAndQuality() (avgDeliveryPct, avgCodeQuality float64, err error) {
	var ids []uint
	if err = r.db.Model(&authDomain.User{}).Where("role IN ?", []string{authDomain.RoleEngineer, authDomain.RoleChiefEngineer}).Pluck("id", &ids).Error; err != nil {
		return 0, 0, err
	}
	if len(ids) == 0 {
		return 0, 0, nil
	}
	var sumDelivery float64
	for _, uid := range ids {
		withDue, onTime, _ := r.GetUserTaskDeliveryStats(uid)
		if withDue > 0 {
			sumDelivery += float64(onTime) / float64(withDue) * 100
		}
	}
	avgDeliveryPct = sumDelivery / float64(len(ids))

	var sumQuality float64
	var countQuality int
	for _, uid := range ids {
		avgScore, total, _, _ := r.GetUserSubmissionStats(uid)
		if total > 0 {
			sumQuality += avgScore
			countQuality++
		}
	}
	if countQuality > 0 {
		avgCodeQuality = sumQuality / float64(countQuality)
	}
	return avgDeliveryPct, avgCodeQuality, nil
}

// GetDisciplineStats returns per-user per-day activity for the given date range.
// from/to are inclusive, format YYYY-MM-DD.
func (r *postgresRepo) GetDisciplineStats(from, to string) (*perfDomain.DisciplineResponse, error) {
	// 1. Build ordered date list
	fromTime, err := time.Parse("2006-01-02", from)
	if err != nil {
		return nil, err
	}
	toTime, err := time.Parse("2006-01-02", to)
	if err != nil {
		return nil, err
	}
	var dates []string
	for d := fromTime; !d.After(toTime); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}

	// 2. Load working-level users only (exclude CEO and SUPPORT)
	var users []authDomain.User
	if err := r.db.Select("id", "email", "display_name", "role", "avatar_url").
		Where("role NOT IN ?", []string{"CEO", "SUPPORT"}).
		Find(&users).Error; err != nil {
		return nil, err
	}

	// 3. Bulk-query stats with raw SQL for efficiency
	type tasksClosedRow struct {
		UserID uint
		Date   string
		Count  int
	}
	var closedRows []tasksClosedRow
	r.db.Raw(`
		SELECT assigned_to AS user_id, completed_at::date::text AS date, COUNT(*) AS count
		FROM tasks
		WHERE status = 'COMPLETED'
		  AND completed_at IS NOT NULL
		  AND completed_at::date BETWEEN ? AND ?
		GROUP BY assigned_to, completed_at::date
	`, from, to).Scan(&closedRows)

	type reworkRow struct {
		UserID uint
		Date   string
		Count  int
	}
	var reworkRows []reworkRow
	r.db.Raw(`
		SELECT t.assigned_to AS user_id, tc.created_at::date::text AS date, COUNT(*) AS count
		FROM task_comments tc
		JOIN tasks t ON t.id = tc.task_id
		WHERE tc.content LIKE '[REJECTED]%'
		  AND tc.created_at::date BETWEEN ? AND ?
		  AND t.assigned_to IS NOT NULL
		GROUP BY t.assigned_to, tc.created_at::date
	`, from, to).Scan(&reworkRows)

	type logRow struct {
		UserID  uint
		Date    string
		Minutes int
	}
	var logRows []logRow
	r.db.Raw(`
		SELECT user_id, logged_date::text AS date, COALESCE(SUM(minutes), 0)::int AS minutes
		FROM time_logs
		WHERE logged_date BETWEEN ? AND ?
		GROUP BY user_id, logged_date
	`, from, to).Scan(&logRows)

	type pulseRow struct {
		UserID uint
		Date   string
	}
	var pulseRows []pulseRow
	r.db.Raw(`
		SELECT user_id, date::text AS date
		FROM daily_standups
		WHERE date BETWEEN ? AND ?
	`, from, to).Scan(&pulseRows)

	// Deployment requests marked DEPLOYED (reviewer = Chief Engineer who deployed)
	type deployRow struct {
		UserID uint
		Date   string
		Count  int
	}
	var deployRows []deployRow
	r.db.Raw(`
		SELECT reviewer_id AS user_id, deployed_at::date::text AS date, COUNT(*) AS count
		FROM deployment_requests
		WHERE status = 'DEPLOYED'
		  AND deployed_at IS NOT NULL
		  AND reviewer_id IS NOT NULL
		  AND deployed_at::date BETWEEN ? AND ?
		GROUP BY reviewer_id, deployed_at::date
	`, from, to).Scan(&deployRows)

	// Attendance records (is_late, early_checkout, check-in/out times)
	type attDayRow struct {
		UserID        uint
		Date          string
		IsLate        bool
		EarlyCheckout bool
		Status        string
		CheckInTime   string
		CheckOutTime  string
	}
	var attDayRows []attDayRow
	r.db.Raw(`
		SELECT user_id,
		       attendance_date::text AS date,
		       is_late,
		       early_checkout,
		       COALESCE(status, '') AS status,
		       COALESCE(TO_CHAR(check_in_at AT TIME ZONE 'Asia/Bangkok', 'HH24:MI'), '') AS check_in_time,
		       COALESCE(TO_CHAR(check_out_at AT TIME ZONE 'Asia/Bangkok', 'HH24:MI'), '') AS check_out_time
		FROM attendance_records
		WHERE attendance_date BETWEEN ? AND ?
	`, from, to).Scan(&attDayRows)

	// 4. Index rows by userID+date for O(1) lookup
	closedIndex := map[uint]map[string]int{}
	for _, row := range closedRows {
		if closedIndex[row.UserID] == nil {
			closedIndex[row.UserID] = map[string]int{}
		}
		closedIndex[row.UserID][row.Date] = row.Count
	}
	reworkIndex := map[uint]map[string]int{}
	for _, row := range reworkRows {
		if reworkIndex[row.UserID] == nil {
			reworkIndex[row.UserID] = map[string]int{}
		}
		reworkIndex[row.UserID][row.Date] = row.Count
	}
	logIndex := map[uint]map[string]int{}
	for _, row := range logRows {
		if logIndex[row.UserID] == nil {
			logIndex[row.UserID] = map[string]int{}
		}
		logIndex[row.UserID][row.Date] = row.Minutes
	}
	pulseIndex := map[uint]map[string]bool{}
	for _, row := range pulseRows {
		if pulseIndex[row.UserID] == nil {
			pulseIndex[row.UserID] = map[string]bool{}
		}
		pulseIndex[row.UserID][row.Date] = true
	}
	deployIndex := map[uint]map[string]int{}
	for _, row := range deployRows {
		if deployIndex[row.UserID] == nil {
			deployIndex[row.UserID] = map[string]int{}
		}
		deployIndex[row.UserID][row.Date] = row.Count
	}

	type attDayData struct {
		IsLate        bool
		EarlyCheckout bool
		Status        string
		CheckInTime   string
		CheckOutTime  string
	}
	attIndex := map[uint]map[string]attDayData{}
	for _, row := range attDayRows {
		if attIndex[row.UserID] == nil {
			attIndex[row.UserID] = map[string]attDayData{}
		}
		attIndex[row.UserID][row.Date] = attDayData{
			IsLate:        row.IsLate,
			EarlyCheckout: row.EarlyCheckout,
			Status:        row.Status,
			CheckInTime:   row.CheckInTime,
			CheckOutTime:  row.CheckOutTime,
		}
	}

	// 5. Assemble per-user stats
	resp := &perfDomain.DisciplineResponse{
		FromDate: from,
		ToDate:   to,
		Dates:    dates,
		Users:    make([]perfDomain.DisciplineUser, 0, len(users)),
	}
	for _, u := range users {
		du := perfDomain.DisciplineUser{
			UserID:          u.ID,
			UserEmail:       u.Email,
			UserDisplayName: u.DisplayName,
			UserAvatarURL:   u.AvatarURL,
			Role:            u.Role,
		}
		var totalClosed, totalReworks, totalMins, missedPulse, totalDeploys, totalLate, totalEarlyOut int
		for _, d := range dates {
			closed := closedIndex[u.ID][d]
			rework := reworkIndex[u.ID][d]
			mins := logIndex[u.ID][d]
			hasPulse := pulseIndex[u.ID][d]
			deploys := deployIndex[u.ID][d]
			att := attIndex[u.ID][d]
			if !hasPulse {
				missedPulse++
			}
			if att.IsLate {
				totalLate++
			}
			if att.EarlyCheckout {
				totalEarlyOut++
			}
			totalClosed += closed
			totalReworks += rework
			totalMins += mins
			totalDeploys += deploys
			du.Days = append(du.Days, perfDomain.DisciplineUserDayStat{
				Date:                 d,
				TasksClosed:          closed,
				Reworks:              rework,
				LoggedMinutes:        mins,
				HasDailyPulse:        hasPulse,
				DeploymentsCompleted: deploys,
				IsLate:               att.IsLate,
				EarlyCheckout:        att.EarlyCheckout,
				CheckInAt:            att.CheckInTime,
				CheckOutAt:           att.CheckOutTime,
				AttendanceStatus:     att.Status,
			})
		}
		du.TotalTasksClosed = totalClosed
		du.TotalReworks = totalReworks
		du.TotalLoggedHours = float64(totalMins) / 60.0
		du.MissedPulseCount = missedPulse
		du.TotalDeployments = totalDeploys
		du.TotalLateDays = totalLate
		du.TotalEarlyCheckoutDays = totalEarlyOut
		resp.Users = append(resp.Users, du)
	}
	return resp, nil
}

// GetDisciplineDayDetail returns drill-down activity for one user on one date.
func (r *postgresRepo) GetDisciplineDayDetail(userID uint, date string) (*perfDomain.DisciplineDayDetail, error) {
	// Resolve user info
	var u authDomain.User
	if err := r.db.Select("id", "email", "display_name").First(&u, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	detail := &perfDomain.DisciplineDayDetail{
		UserID:          u.ID,
		UserEmail:       u.Email,
		UserDisplayName: u.DisplayName,
		Date:            date,
	}

	// Daily Pulse check
	var pulseCount int64
	r.db.Raw("SELECT COUNT(*) FROM daily_standups WHERE user_id = ? AND date = ?", userID, date).Scan(&pulseCount)
	detail.HasDailyPulse = pulseCount > 0

	// Attendance record for this user/date
	type attDetailRow struct {
		CheckInAt     *time.Time
		CheckOutAt    *time.Time
		IsLate        bool
		EarlyCheckout bool
		Status        string
		CheckInMethod string
	}
	var attRec attDetailRow
	r.db.Raw(`
		SELECT check_in_at, check_out_at, is_late, early_checkout,
		       COALESCE(status, '') AS status,
		       COALESCE(check_in_method, '') AS check_in_method
		FROM attendance_records
		WHERE user_id = ? AND attendance_date = ?
	`, userID, date).Scan(&attRec)
	if attRec.Status != "" {
		rec := &perfDomain.DisciplineAttendanceRecord{
			IsLate:        attRec.IsLate,
			EarlyCheckout: attRec.EarlyCheckout,
			Status:        attRec.Status,
			CheckInMethod: attRec.CheckInMethod,
		}
		if attRec.CheckInAt != nil {
			rec.CheckInAt = attRec.CheckInAt.In(ict).Format("15:04")
		}
		if attRec.CheckOutAt != nil {
			rec.CheckOutAt = attRec.CheckOutAt.In(ict).Format("15:04")
		}
		detail.Attendance = rec
	}

	// Time logs with task info
	type tlRow struct {
		TaskID          string
		TaskCode        string
		TaskTitle       string
		Minutes         int
		Description     string
		WorkType        string
		IsTimerSession  bool
	}
	var tlRows []tlRow
	r.db.Raw(`
		SELECT tl.task_id::text AS task_id,
		       COALESCE(t.code, '') AS task_code,
		       COALESCE(t.title, 'Unknown Task') AS task_title,
		       tl.minutes,
		       COALESCE(tl.description, '') AS description,
		       COALESCE(tl.work_type, 'DEV') AS work_type,
		       COALESCE(tl.is_timer_session, false) AS is_timer_session
		FROM time_logs tl
		LEFT JOIN tasks t ON t.id = tl.task_id
		WHERE tl.user_id = ?
		  AND tl.logged_date = ?
		ORDER BY tl.logged_at
	`, userID, date).Scan(&tlRows)

	totalMins := 0
	for _, row := range tlRows {
		totalMins += row.Minutes
		detail.TimeLogs = append(detail.TimeLogs, perfDomain.DisciplineTimeLogEntry{
			TaskID:      row.TaskID,
			TaskCode:    row.TaskCode,
			TaskTitle:   row.TaskTitle,
			Minutes:     row.Minutes,
			Hours:       float64(row.Minutes) / 60.0,
			Description: row.Description,
			WorkType:    row.WorkType,
			IsTimer:     row.IsTimerSession,
		})
	}
	detail.TotalLoggedMin = totalMins

	// Completed tasks on this date
	type ctRow struct {
		TaskID      string
		TaskCode    string
		TaskTitle   string
		StoryPoints int
		TaskType    string
	}
	var ctRows []ctRow
	r.db.Raw(`
		SELECT id::text AS task_id,
		       COALESCE(code, '') AS task_code,
		       title AS task_title,
		       story_points,
		       task_type
		FROM tasks
		WHERE assigned_to = ?
		  AND status = 'COMPLETED'
		  AND completed_at IS NOT NULL
		  AND completed_at::date = ?
		ORDER BY completed_at
	`, userID, date).Scan(&ctRows)

	for _, row := range ctRows {
		detail.CompletedTasks = append(detail.CompletedTasks, perfDomain.DisciplineCompletedTask{
			TaskID:      row.TaskID,
			TaskCode:    row.TaskCode,
			TaskTitle:   row.TaskTitle,
			StoryPoints: row.StoryPoints,
			TaskType:    row.TaskType,
		})
	}

	// Rework events (REJECTED comments) on this date
	type rwRow struct {
		TaskID          string
		TaskCode        string
		TaskTitle       string
		RejectedComment string
	}
	var rwRows []rwRow
	r.db.Raw(`
		SELECT t.id::text AS task_id,
		       COALESCE(t.code, '') AS task_code,
		       COALESCE(t.title, 'Unknown Task') AS task_title,
		       tc.content AS rejected_comment
		FROM task_comments tc
		JOIN tasks t ON t.id = tc.task_id
		WHERE t.assigned_to = ?
		  AND tc.content LIKE '[REJECTED]%'
		  AND tc.created_at::date = ?
		ORDER BY tc.created_at
	`, userID, date).Scan(&rwRows)

	for _, row := range rwRows {
		detail.Reworks = append(detail.Reworks, perfDomain.DisciplineReworkEntry{
			TaskID:          row.TaskID,
			TaskCode:        row.TaskCode,
			TaskTitle:       row.TaskTitle,
			RejectedComment: row.RejectedComment,
		})
	}

	// Deployment requests marked DEPLOYED by this reviewer on this date
	type drRow struct {
		ID          uint
		Title       string
		Branch      string
		Environment string
	}
	var drRows []drRow
	r.db.Raw(`
		SELECT id, title, branch, environment
		FROM deployment_requests
		WHERE reviewer_id = ?
		  AND status = 'DEPLOYED'
		  AND deployed_at IS NOT NULL
		  AND deployed_at::date = ?
		ORDER BY deployed_at
	`, userID, date).Scan(&drRows)

	for _, row := range drRows {
		detail.DeployedRequests = append(detail.DeployedRequests, perfDomain.DisciplineDeployedRequest{
			ID:          row.ID,
			Title:       row.Title,
			Branch:      row.Branch,
			Environment: row.Environment,
		})
	}

	if detail.TimeLogs == nil {
		detail.TimeLogs = []perfDomain.DisciplineTimeLogEntry{}
	}
	if detail.CompletedTasks == nil {
		detail.CompletedTasks = []perfDomain.DisciplineCompletedTask{}
	}
	if detail.Reworks == nil {
		detail.Reworks = []perfDomain.DisciplineReworkEntry{}
	}
	if detail.DeployedRequests == nil {
		detail.DeployedRequests = []perfDomain.DisciplineDeployedRequest{}
	}

	return detail, nil
}

func (r *postgresRepo) GetCompanyWideReworkAndTimeAccuracy() (avgReworkPct, avgTimeAccuracyPct float64, err error) {
	var ids []uint
	if err = r.db.Model(&authDomain.User{}).Where("role IN ?", []string{authDomain.RoleEngineer, authDomain.RoleChiefEngineer}).Pluck("id", &ids).Error; err != nil {
		return 0, 0, err
	}
	if len(ids) == 0 {
		return 0, 0, nil
	}
	var sumRework, sumTimeAcc float64
	var countRework, countTimeAcc int
	for _, uid := range ids {
		_, total, failCount, _ := r.GetUserSubmissionStats(uid)
		if total > 0 {
			sumRework += float64(failCount) / float64(total) * 100
			countRework++
		}
		acc, n, _ := r.GetUserTimeAccuracy(uid)
		if n > 0 {
			sumTimeAcc += acc
			countTimeAcc++
		}
	}
	if countRework > 0 {
		avgReworkPct = sumRework / float64(countRework)
	}
	if countTimeAcc > 0 {
		avgTimeAccuracyPct = sumTimeAcc / float64(countTimeAcc)
	}
	return avgReworkPct, avgTimeAccuracyPct, nil
}
 
