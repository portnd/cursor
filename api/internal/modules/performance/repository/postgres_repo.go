package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
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

const disciplineStartDateKey = "discipline_start_date"
const hiddenPulseUsersSettingKey = "pulse_hidden_user_ids"

func engineerJobDoneEventClause(alias string) string {
	return fmt.Sprintf(`(
		(%[1]s.action = 'READY_FOR_TEST'
			AND COALESCE(%[1]s.payload->>'from_status', '') = 'IN_PROGRESS'
			AND COALESCE(%[1]s.payload->>'to_status', '') = 'READY_FOR_TEST')
		OR
		(%[1]s.action = 'STATUS_CHANGED'
			AND COALESCE(%[1]s.payload->>'from_status', '') = 'IN_PROGRESS'
			AND COALESCE(%[1]s.payload->>'to_status', '') = 'READY_FOR_TEST')
	)`, alias)
}

// NewPostgresRepository returns a performance repository that queries existing tables
func NewPostgresRepository(db *gorm.DB) perfDomain.Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) GetUserTaskDeliveryStats(userID uint) (tasksWithDue int, completedOnTime int, err error) {
	type agg struct {
		TasksWithDue    int64 `gorm:"column:tasks_with_due"`
		CompletedOnTime int64 `gorm:"column:completed_on_time"`
	}
	var a agg
	query := fmt.Sprintf(`
		WITH first_job_done AS (
			SELECT tae.task_id, MIN(tae.created_at) AS job_done_at
			FROM task_activity_events tae
			JOIN tasks src ON src.id = tae.task_id
			WHERE src.assigned_to = ?
			  AND src.due_at IS NOT NULL
			  AND %s
			GROUP BY tae.task_id
		)
		SELECT
			COUNT(t.id) AS tasks_with_due,
			COUNT(t.id) FILTER (
				WHERE fjd.job_done_at IS NOT NULL
				  AND fjd.job_done_at <= t.due_at
			) AS completed_on_time
		FROM tasks t
		LEFT JOIN first_job_done fjd ON fjd.task_id = t.id
		WHERE t.assigned_to = ?
		  AND t.due_at IS NOT NULL
	`, engineerJobDoneEventClause("tae"))
	if err = r.db.Raw(query, userID, userID).Scan(&a).Error; err != nil {
		return 0, 0, err
	}
	return int(a.TasksWithDue), int(a.CompletedOnTime), nil
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

func (r *postgresRepo) GetUserReworkStats(userID uint) (jobDoneCount int, reworkCount int, err error) {
	var resetAt *time.Time
	type resetRow struct{ ReworkResetAt *time.Time }
	var rr resetRow
	r.db.Raw("SELECT rework_reset_at FROM users WHERE id = ?", userID).Scan(&rr)
	resetAt = rr.ReworkResetAt

	var jobDone int64
	jobDoneQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM task_activity_events tae
		JOIN tasks t ON t.id = tae.task_id
		WHERE t.assigned_to = ?
		  AND t.task_type IN ('TASK', 'BUG')
		  AND %s
	`, engineerJobDoneEventClause("tae"))
	args := []interface{}{userID}
	if resetAt != nil {
		jobDoneQuery += " AND tae.created_at > ?"
		args = append(args, *resetAt)
	}
	if err = r.db.Raw(jobDoneQuery, args...).Scan(&jobDone).Error; err != nil {
		return 0, 0, err
	}

	var reworks int64
	reworkQuery := `
		SELECT COUNT(*)
		FROM task_comments tc
		JOIN tasks t ON t.id = tc.task_id
		WHERE t.assigned_to = ?
		  AND tc.content LIKE '[REJECTED]%'
	`
	reworkArgs := []interface{}{userID}
	if resetAt != nil {
		reworkQuery += " AND tc.created_at > ?"
		reworkArgs = append(reworkArgs, *resetAt)
	}
	if err = r.db.Raw(reworkQuery, reworkArgs...).Scan(&reworks).Error; err != nil {
		return int(jobDone), 0, err
	}

	return int(jobDone), int(reworks), nil
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
	query := fmt.Sprintf(`
		WITH first_job_done AS (
			SELECT tae.task_id, MIN(tae.created_at) AS job_done_at
			FROM task_activity_events tae
			WHERE %s
			GROUP BY tae.task_id
		)
		SELECT ROW_NUMBER() OVER (ORDER BY s.end_date DESC NULLS LAST) AS sprint_order,
			COALESCE(SUM(CASE WHEN fjd.job_done_at IS NOT NULL THEN t.story_points ELSE 0 END), 0)::int AS sp
		FROM sprints s
		LEFT JOIN tasks t ON t.sprint_id = s.id AND t.assigned_to = ?
		LEFT JOIN first_job_done fjd ON fjd.task_id = t.id
		WHERE s.status = 'COMPLETED'
		GROUP BY s.id, s.end_date
		ORDER BY s.end_date DESC NULLS LAST
		LIMIT ?
	`, engineerJobDoneEventClause("tae"))
	err = r.db.Raw(query, userID, lastNSprints).Scan(&rows).Error
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
	type agg struct {
		TasksWithDue    int64 `gorm:"column:tasks_with_due"`
		CompletedOnTime int64 `gorm:"column:completed_on_time"`
	}
	var a agg
	query := fmt.Sprintf(`
		WITH first_job_done AS (
			SELECT tae.task_id, MIN(tae.created_at) AS job_done_at
			FROM task_activity_events tae
			JOIN tasks src ON src.id = tae.task_id
			WHERE src.assigned_to = ?
			  AND src.assigned_by_id = ?
			  AND src.due_at IS NOT NULL
			  AND %s
			GROUP BY tae.task_id
		)
		SELECT
			COUNT(t.id) AS tasks_with_due,
			COUNT(t.id) FILTER (
				WHERE fjd.job_done_at IS NOT NULL
				  AND fjd.job_done_at <= t.due_at
			) AS completed_on_time
		FROM tasks t
		LEFT JOIN first_job_done fjd ON fjd.task_id = t.id
		WHERE t.assigned_to = ?
		  AND t.assigned_by_id = ?
		  AND t.due_at IS NOT NULL
	`, engineerJobDoneEventClause("tae"))
	if err = r.db.Raw(query, devID, assignedByID, devID, assignedByID).Scan(&a).Error; err != nil {
		return 0, 0, err
	}
	return int(a.TasksWithDue), int(a.CompletedOnTime), nil
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

func (r *postgresRepo) GetUserReworkStatsForAssignedBy(devID, assignedByID uint) (jobDoneCount int, reworkCount int, err error) {
	var resetAt *time.Time
	type resetRow struct{ ReworkResetAt *time.Time }
	var rr resetRow
	r.db.Raw("SELECT rework_reset_at FROM users WHERE id = ?", devID).Scan(&rr)
	resetAt = rr.ReworkResetAt

	var jobDone int64
	jobDoneQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM task_activity_events tae
		JOIN tasks t ON t.id = tae.task_id
		WHERE t.assigned_to = ?
		  AND t.assigned_by_id = ?
		  AND t.task_type IN ('TASK', 'BUG')
		  AND %s
	`, engineerJobDoneEventClause("tae"))
	jobDoneArgs := []interface{}{devID, assignedByID}
	if resetAt != nil {
		jobDoneQuery += " AND tae.created_at > ?"
		jobDoneArgs = append(jobDoneArgs, *resetAt)
	}
	if err = r.db.Raw(jobDoneQuery, jobDoneArgs...).Scan(&jobDone).Error; err != nil {
		return 0, 0, err
	}

	var reworks int64
	reworkQuery := `
		SELECT COUNT(*)
		FROM task_comments tc
		JOIN tasks t ON t.id = tc.task_id
		WHERE t.assigned_to = ?
		  AND t.assigned_by_id = ?
		  AND tc.content LIKE '[REJECTED]%'
	`
	reworkArgs := []interface{}{devID, assignedByID}
	if resetAt != nil {
		reworkQuery += " AND tc.created_at > ?"
		reworkArgs = append(reworkArgs, *resetAt)
	}
	if err = r.db.Raw(reworkQuery, reworkArgs...).Scan(&reworks).Error; err != nil {
		return int(jobDone), 0, err
	}

	return int(jobDone), int(reworks), nil
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
	query := fmt.Sprintf(`
		WITH first_job_done AS (
			SELECT tae.task_id, MIN(tae.created_at) AS job_done_at
			FROM task_activity_events tae
			WHERE %s
			GROUP BY tae.task_id
		)
		SELECT ROW_NUMBER() OVER (ORDER BY s.end_date DESC NULLS LAST) AS sprint_order,
			COALESCE(SUM(CASE WHEN fjd.job_done_at IS NOT NULL THEN t.story_points ELSE 0 END), 0)::int AS sp
		FROM sprints s
		LEFT JOIN tasks t ON t.sprint_id = s.id AND t.assigned_to = ? AND t.assigned_by_id = ?
		LEFT JOIN first_job_done fjd ON fjd.task_id = t.id
		WHERE s.status = 'COMPLETED'
		GROUP BY s.id, s.end_date
		ORDER BY s.end_date DESC NULLS LAST
		LIMIT ?
	`, engineerJobDoneEventClause("tae"))
	err = r.db.Raw(query, devID, assignedByID, lastNSprints).Scan(&rows).Error
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
	// Sprints with >= 80% of planned story points reaching "job done"
	type agg struct {
		Total   int64
		Success int64
	}
	var a agg
	query := fmt.Sprintf(`
		WITH first_job_done AS (
			SELECT tae.task_id, MIN(tae.created_at) AS job_done_at
			FROM task_activity_events tae
			WHERE %s
			GROUP BY tae.task_id
		),
		sprint_sp AS (
			SELECT s.id,
				COALESCE(SUM(t.story_points), 0) AS planned,
				COALESCE(SUM(CASE WHEN fjd.job_done_at IS NOT NULL THEN t.story_points ELSE 0 END), 0) AS completed
			FROM sprints s
			LEFT JOIN tasks t ON t.sprint_id = s.id
			LEFT JOIN first_job_done fjd ON fjd.task_id = t.id
			WHERE s.status = 'COMPLETED'
			GROUP BY s.id
		)
		SELECT COUNT(*) AS total,
			SUM(CASE WHEN planned = 0 OR (completed::float / NULLIF(planned, 0)) >= 0.8 THEN 1 ELSE 0 END) AS success
		FROM sprint_sp
	`, engineerJobDoneEventClause("tae"))
	err = r.db.Raw(query).Scan(&a).Error
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
	// Active projects where >= 70% of tasks (with due_at) are on schedule based on first "job done" timestamp.
	type agg struct {
		Total   int64
		OnTrack int64
	}
	var a agg
	query := fmt.Sprintf(`
		WITH first_job_done AS (
			SELECT tae.task_id, MIN(tae.created_at) AS job_done_at
			FROM task_activity_events tae
			WHERE %s
			GROUP BY tae.task_id
		),
		project_stats AS (
			SELECT p.id,
				COUNT(t.id) FILTER (WHERE t.due_at IS NOT NULL) AS with_due,
				COUNT(t.id) FILTER (
					WHERE t.due_at IS NOT NULL
					  AND (
						(fjd.job_done_at IS NOT NULL AND fjd.job_done_at <= t.due_at)
						OR
						(fjd.job_done_at IS NULL AND t.due_at >= NOW())
					  )
				) AS on_schedule
			FROM projects p
			LEFT JOIN tasks t ON t.project_id = p.id AND t.parent_id IS NULL
			LEFT JOIN first_job_done fjd ON fjd.task_id = t.id
			WHERE p.status = 'ACTIVE'
			GROUP BY p.id
		)
		SELECT COUNT(*) AS total,
			SUM(CASE WHEN with_due = 0 OR (on_schedule::float / NULLIF(with_due, 0)) >= 0.7 THEN 1 ELSE 0 END) AS on_track
		FROM project_stats
	`, engineerJobDoneEventClause("tae"))
	err = r.db.Raw(query).Scan(&a).Error
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
	query := fmt.Sprintf(`
		WITH first_job_done AS (
			SELECT tae.task_id, MIN(tae.created_at) AS job_done_at
			FROM task_activity_events tae
			WHERE %s
			GROUP BY tae.task_id
		)
		SELECT COALESCE(SUM(CASE WHEN fjd.job_done_at IS NOT NULL THEN t.story_points ELSE 0 END), 0)::int AS completed_sp
		FROM sprints s
		LEFT JOIN tasks t ON t.sprint_id = s.id
		LEFT JOIN first_job_done fjd ON fjd.task_id = t.id
		WHERE s.status = 'COMPLETED'
		GROUP BY s.id, s.end_date
		ORDER BY s.end_date DESC NULLS LAST
		LIMIT 4
	`, engineerJobDoneEventClause("tae"))
	err = r.db.Raw(query).Scan(&rows).Error
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
	growthPct = (float64(recent) - float64(prev)) / float64(prev) * 100
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

// buildDisplayName returns the best human-readable name for a user.
// Priority: "FirstName LastName" → DisplayName → email prefix.
func buildDisplayName(firstName, lastName, displayName, email string) string {
	fn := strings.TrimSpace(firstName)
	ln := strings.TrimSpace(lastName)
	if fn != "" || ln != "" {
		return strings.TrimSpace(fn + " " + ln)
	}
	if dn := strings.TrimSpace(displayName); dn != "" {
		return dn
	}
	if idx := strings.Index(email, "@"); idx > 0 {
		return email[:idx]
	}
	return email
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

	// 2. Load globally hidden pulse users (CEO setting), then working-level users only.
	hiddenSet := map[uint]bool{}
	{
		var row struct {
			Value string `gorm:"column:value"`
		}
		if err := r.db.Raw("SELECT value FROM app_settings WHERE key = ?", hiddenPulseUsersSettingKey).Scan(&row).Error; err == nil && strings.TrimSpace(row.Value) != "" {
			var hiddenIDs []uint
			if err := json.Unmarshal([]byte(row.Value), &hiddenIDs); err == nil {
				for _, id := range hiddenIDs {
					hiddenSet[id] = true
				}
			}
		}
	}

	var users []authDomain.User
	if err := r.db.Select("id", "email", "first_name", "last_name", "display_name", "role", "avatar_url").
		Where("role NOT IN ?", []string{"CEO", "SUPPORT"}).
		Find(&users).Error; err != nil {
		return nil, err
	}
	if len(hiddenSet) > 0 {
		filtered := make([]authDomain.User, 0, len(users))
		for _, u := range users {
			if hiddenSet[u.ID] {
				continue
			}
			filtered = append(filtered, u)
		}
		users = filtered
	}

	// 3. Bulk-query stats with raw SQL for efficiency
	// Job Done indexes by role:
	// - Engineer/Chief Engineer: assigned TASK/BUG moved IN_PROGRESS -> READY_FOR_TEST. Counts both
	//   dedicated READY_FOR_TEST events (MarkReadyForTest) and STATUS_CHANGED from Kanban bulk update.
	// - Product Owner: approved task moved READY_FOR_TEST -> WAIT_FOR_DEPLOY.
	// - Chief Engineer (reviewer): deployment_requests marked DEPLOYED (merged into tasks_closed; see job_done_items).
	type tasksClosedRow struct {
		UserID uint
		Date   string
		Count  int
	}
	var engineerClosedRows []tasksClosedRow
	r.db.Raw(`
		SELECT t.assigned_to AS user_id, tae.created_at::date::text AS date, COUNT(*) AS count
		FROM task_activity_events tae
		JOIN tasks t ON t.id = tae.task_id
		WHERE t.assigned_to IS NOT NULL
		  AND t.task_type IN ('TASK', 'BUG')
		  AND tae.created_at::date BETWEEN ? AND ?
		  AND (
		    (tae.action = 'READY_FOR_TEST'
		      AND COALESCE(tae.payload->>'from_status', '') = 'IN_PROGRESS'
		      AND COALESCE(tae.payload->>'to_status', '') = 'READY_FOR_TEST')
		    OR
		    (tae.action = 'STATUS_CHANGED'
		      AND COALESCE(tae.payload->>'from_status', '') = 'IN_PROGRESS'
		      AND COALESCE(tae.payload->>'to_status', '') = 'READY_FOR_TEST')
		  )
		GROUP BY t.assigned_to, tae.created_at::date
	`, from, to).Scan(&engineerClosedRows)

	var productOwnerClosedRows []tasksClosedRow
	r.db.Raw(`
		SELECT tae.actor_id AS user_id, tae.created_at::date::text AS date, COUNT(*) AS count
		FROM task_activity_events tae
		WHERE tae.action = 'PM_APPROVED_TEST'
		  AND COALESCE(tae.payload->>'from_status', '') = 'READY_FOR_TEST'
		  AND COALESCE(tae.payload->>'to_status', '') = 'WAIT_FOR_DEPLOY'
		  AND tae.actor_id IS NOT NULL
		  AND tae.created_at::date BETWEEN ? AND ?
		GROUP BY tae.actor_id, tae.created_at::date
	`, from, to).Scan(&productOwnerClosedRows)

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

	type leaveDayRow struct {
		UserID      uint
		Date        string
		IsHalfDay   bool
		HalfSession string
	}
	var leaveDayRows []leaveDayRow
	r.db.Raw(`
		SELECT lr.user_id AS user_id,
		       gs::date::text AS date,
		       lr.is_half_day AS is_half_day,
		       COALESCE(lr.half_day_session, '') AS half_session
		FROM leave_requests lr
		CROSS JOIN LATERAL generate_series(lr.start_date::date, lr.end_date::date, interval '1 day') gs
		WHERE lr.status = 'APPROVED'
		  AND lr.start_date::date <= ?::date
		  AND lr.end_date::date >= ?::date
	`, to, from).Scan(&leaveDayRows)

	type holidayRow struct {
		Date string
	}
	var holidayRows []holidayRow
	r.db.Raw(`
		SELECT date::text AS date
		FROM holiday_calendars
		WHERE date BETWEEN ? AND ?
	`, from, to).Scan(&holidayRows)

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

	// WAIT_FOR_DEPLOY → READY_FOR_UAT when logged as DEPLOYED activity but no matching deployment_requests
	// on the same Bangkok calendar day for that task (avoids double-count with deployment-based Job Done).
	type uatAdvanceRow struct {
		UserID           uint
		TaskID           string
		TaskCode         string
		TaskTitle        string
		TaskType         string
		CreatedAt        time.Time
		ActorID          sql.NullInt64
		ActorEmail       string
		ActorDisplayName string
	}
	var uatAdvanceRows []uatAdvanceRow
	r.db.Raw(`
		SELECT tae.actor_id AS user_id,
		       t.id::text AS task_id,
		       COALESCE(t.code, '') AS task_code,
		       t.title AS task_title,
		       COALESCE(t.task_type, '') AS task_type,
		       tae.created_at AS created_at,
		       tae.actor_id AS actor_id,
		       COALESCE(actor.email, '') AS actor_email,
		       COALESCE(actor.display_name, '') AS actor_display_name
		FROM task_activity_events tae
		JOIN tasks t ON t.id = tae.task_id
		LEFT JOIN users actor ON actor.id = tae.actor_id
		WHERE tae.action = 'DEPLOYED'
		  AND COALESCE(tae.payload->>'from_status', '') = 'WAIT_FOR_DEPLOY'
		  AND COALESCE(tae.payload->>'to_status', '') = 'READY_FOR_UAT'
		  AND tae.actor_id IS NOT NULL
		  AND tae.created_at::date BETWEEN ? AND ?
		  AND NOT EXISTS (
		    SELECT 1 FROM deployment_requests dr
		    WHERE dr.task_id = tae.task_id
		      AND dr.status = 'DEPLOYED'
		      AND dr.deployed_at IS NOT NULL
		      AND (timezone('Asia/Bangkok', dr.deployed_at))::date = (timezone('Asia/Bangkok', tae.created_at))::date
		  )
		ORDER BY tae.created_at DESC
	`, from, to).Scan(&uatAdvanceRows)

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
	for _, row := range engineerClosedRows {
		if closedIndex[row.UserID] == nil {
			closedIndex[row.UserID] = map[string]int{}
		}
		closedIndex[row.UserID][row.Date] += row.Count
	}
	for _, row := range productOwnerClosedRows {
		if closedIndex[row.UserID] == nil {
			closedIndex[row.UserID] = map[string]int{}
		}
		closedIndex[row.UserID][row.Date] += row.Count
	}
	for _, row := range deployRows {
		if closedIndex[row.UserID] == nil {
			closedIndex[row.UserID] = map[string]int{}
		}
		closedIndex[row.UserID][row.Date] += row.Count
	}
	for _, row := range uatAdvanceRows {
		d := row.CreatedAt.In(ict).Format("2006-01-02")
		if closedIndex[row.UserID] == nil {
			closedIndex[row.UserID] = map[string]int{}
		}
		closedIndex[row.UserID][d]++
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
	type leaveDayData struct {
		OnLeave      bool
		LeaveSession string
	}
	leaveIndex := map[uint]map[string]leaveDayData{}
	for _, row := range leaveDayRows {
		if leaveIndex[row.UserID] == nil {
			leaveIndex[row.UserID] = map[string]leaveDayData{}
		}
		session := "FULL"
		if row.IsHalfDay {
			hs := strings.ToUpper(strings.TrimSpace(row.HalfSession))
			if hs == "AM" || hs == "PM" {
				session = hs
			}
		}
		leaveIndex[row.UserID][row.Date] = leaveDayData{OnLeave: true, LeaveSession: session}
	}
	holidayIndex := map[string]bool{}
	for _, row := range holidayRows {
		holidayIndex[row.Date] = true
	}
	// Deploy completions are included in tasks_closed / total_tasks_closed (Job Done) for the CE reviewer.
	deployIndex := map[uint]map[string]int{}

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

	// 4b. Job Done line items (matches engineer + PO aggregates; for UI drill-down)
	type jobDoneEventRow struct {
		UserID           uint
		TaskID           string
		TaskCode         string
		TaskTitle        string
		TaskType         string
		CreatedAt        time.Time
		EventKind        string
		ActorID          sql.NullInt64
		ActorEmail       string
		ActorDisplayName string
	}
	var engJobRows []jobDoneEventRow
	r.db.Raw(`
		SELECT t.assigned_to AS user_id,
		       t.id::text AS task_id,
		       COALESCE(t.code, '') AS task_code,
		       t.title AS task_title,
		       COALESCE(t.task_type, '') AS task_type,
		       tae.created_at AS created_at,
		       'READY_FOR_TEST' AS event_kind,
		       tae.actor_id AS actor_id,
		       COALESCE(actor.email, '') AS actor_email,
		       COALESCE(actor.display_name, '') AS actor_display_name
		FROM task_activity_events tae
		JOIN tasks t ON t.id = tae.task_id
		LEFT JOIN users actor ON actor.id = tae.actor_id
		WHERE t.assigned_to IS NOT NULL
		  AND t.task_type IN ('TASK', 'BUG')
		  AND tae.created_at::date BETWEEN ? AND ?
		  AND (
		    (tae.action = 'READY_FOR_TEST'
		      AND COALESCE(tae.payload->>'from_status', '') = 'IN_PROGRESS'
		      AND COALESCE(tae.payload->>'to_status', '') = 'READY_FOR_TEST')
		    OR
		    (tae.action = 'STATUS_CHANGED'
		      AND COALESCE(tae.payload->>'from_status', '') = 'IN_PROGRESS'
		      AND COALESCE(tae.payload->>'to_status', '') = 'READY_FOR_TEST')
		  )
		ORDER BY tae.created_at DESC
	`, from, to).Scan(&engJobRows)

	var poJobRows []jobDoneEventRow
	r.db.Raw(`
		SELECT tae.actor_id AS user_id,
		       t.id::text AS task_id,
		       COALESCE(t.code, '') AS task_code,
		       t.title AS task_title,
		       COALESCE(t.task_type, '') AS task_type,
		       tae.created_at AS created_at,
		       'PM_APPROVED_TEST' AS event_kind,
		       tae.actor_id AS actor_id,
		       COALESCE(actor.email, '') AS actor_email,
		       COALESCE(actor.display_name, '') AS actor_display_name
		FROM task_activity_events tae
		JOIN tasks t ON t.id = tae.task_id
		LEFT JOIN users actor ON actor.id = tae.actor_id
		WHERE tae.action = 'PM_APPROVED_TEST'
		  AND COALESCE(tae.payload->>'from_status', '') = 'READY_FOR_TEST'
		  AND COALESCE(tae.payload->>'to_status', '') = 'WAIT_FOR_DEPLOY'
		  AND tae.actor_id IS NOT NULL
		  AND tae.created_at::date BETWEEN ? AND ?
		ORDER BY tae.created_at DESC
	`, from, to).Scan(&poJobRows)

	// Chief Engineer: deployment marked DEPLOYED (linked task advances to READY_FOR_UAT when task_id set)
	type deployJobDoneRow struct {
		UserID           uint
		TaskID           string
		TaskCode         string
		TaskTitle        string
		TaskType         string
		DeployedAt       time.Time
		ActorID          sql.NullInt64
		ActorEmail       string
		ActorDisplayName string
		DeploymentID     uint
		DeploymentTitle  string
		Branch           string
		Environment      string
	}
	var deployJobDoneRows []deployJobDoneRow
	r.db.Raw(`
		SELECT dr.reviewer_id AS user_id,
		       COALESCE(dr.task_id::text, '') AS task_id,
		       COALESCE(t.code, '') AS task_code,
		       COALESCE(NULLIF(TRIM(t.title), ''), dr.title) AS task_title,
		       COALESCE(t.task_type, '') AS task_type,
		       dr.deployed_at AS deployed_at,
		       dr.reviewer_id AS actor_id,
		       COALESCE(actor.email, '') AS actor_email,
		       COALESCE(actor.display_name, '') AS actor_display_name,
		       dr.id AS deployment_id,
		       dr.title AS deployment_title,
		       dr.branch AS branch,
		       dr.environment AS environment
		FROM deployment_requests dr
		LEFT JOIN tasks t ON t.id = dr.task_id
		LEFT JOIN users actor ON actor.id = dr.reviewer_id
		WHERE dr.status = 'DEPLOYED'
		  AND dr.deployed_at IS NOT NULL
		  AND dr.reviewer_id IS NOT NULL
		  AND dr.deployed_at::date BETWEEN ? AND ?
		ORDER BY dr.deployed_at DESC
	`, from, to).Scan(&deployJobDoneRows)

	jobDoneItemFromRow := func(row jobDoneEventRow) perfDomain.DisciplineJobDoneItem {
		ts := row.CreatedAt.In(ict)
		item := perfDomain.DisciplineJobDoneItem{
			TaskID:           row.TaskID,
			TaskCode:         row.TaskCode,
			TaskTitle:        row.TaskTitle,
			TaskType:         row.TaskType,
			DoneDate:         ts.Format("2006-01-02"),
			DoneTime:         ts.Format("15:04"),
			EventKind:        row.EventKind,
			ActorEmail:       row.ActorEmail,
			ActorDisplayName: row.ActorDisplayName,
		}
		if row.ActorID.Valid && row.ActorID.Int64 > 0 {
			item.ActorID = uint(row.ActorID.Int64)
		}
		return item
	}

	jobDoneByUser := map[uint][]perfDomain.DisciplineJobDoneItem{}
	for _, row := range engJobRows {
		jobDoneByUser[row.UserID] = append(jobDoneByUser[row.UserID], jobDoneItemFromRow(row))
	}
	for _, row := range poJobRows {
		jobDoneByUser[row.UserID] = append(jobDoneByUser[row.UserID], jobDoneItemFromRow(row))
	}
	for _, row := range deployJobDoneRows {
		ts := row.DeployedAt.In(ict)
		item := perfDomain.DisciplineJobDoneItem{
			TaskID:           row.TaskID,
			TaskCode:         row.TaskCode,
			TaskTitle:        row.TaskTitle,
			TaskType:         row.TaskType,
			DoneDate:         ts.Format("2006-01-02"),
			DoneTime:         ts.Format("15:04"),
			EventKind:        "DEPLOYMENT_DEPLOYED",
			ActorEmail:       row.ActorEmail,
			ActorDisplayName: row.ActorDisplayName,
			DeploymentID:     row.DeploymentID,
			DeploymentTitle:  row.DeploymentTitle,
			Branch:           row.Branch,
			Environment:      row.Environment,
		}
		if row.ActorID.Valid && row.ActorID.Int64 > 0 {
			item.ActorID = uint(row.ActorID.Int64)
		}
		jobDoneByUser[row.UserID] = append(jobDoneByUser[row.UserID], item)
	}
	for _, row := range uatAdvanceRows {
		ts := row.CreatedAt.In(ict)
		item := perfDomain.DisciplineJobDoneItem{
			TaskID:           row.TaskID,
			TaskCode:         row.TaskCode,
			TaskTitle:        row.TaskTitle,
			TaskType:         row.TaskType,
			DoneDate:         ts.Format("2006-01-02"),
			DoneTime:         ts.Format("15:04"),
			EventKind:        "DEPLOYED_TO_UAT",
			ActorEmail:       row.ActorEmail,
			ActorDisplayName: row.ActorDisplayName,
		}
		if row.ActorID.Valid && row.ActorID.Int64 > 0 {
			item.ActorID = uint(row.ActorID.Int64)
		}
		jobDoneByUser[row.UserID] = append(jobDoneByUser[row.UserID], item)
	}
	for uid, items := range jobDoneByUser {
		sort.Slice(items, func(i, j int) bool {
			if items[i].DoneDate != items[j].DoneDate {
				return items[i].DoneDate > items[j].DoneDate
			}
			if items[i].DoneTime != items[j].DoneTime {
				return items[i].DoneTime > items[j].DoneTime
			}
			if items[i].EventKind != items[j].EventKind {
				return items[i].EventKind > items[j].EventKind
			}
			if items[i].DeploymentID != items[j].DeploymentID {
				return items[i].DeploymentID > items[j].DeploymentID
			}
			return items[i].TaskCode > items[j].TaskCode
		})
		jobDoneByUser[uid] = items
	}

	// 4c. Rework line items ([REJECTED] comments → assignee’s discipline rework count)
	type reworkDetailRow struct {
		UserID            uint
		TaskID            string
		TaskCode          string
		TaskTitle         string
		TaskType          string
		CreatedAt         time.Time
		Content           string
		AuthorID          sql.NullInt64
		AuthorEmail       string
		AuthorDisplayName string
	}
	var reworkDetailRows []reworkDetailRow
	r.db.Raw(`
		SELECT t.assigned_to AS user_id,
		       t.id::text AS task_id,
		       COALESCE(t.code, '') AS task_code,
		       t.title AS task_title,
		       COALESCE(t.task_type, '') AS task_type,
		       tc.created_at AS created_at,
		       tc.content AS content,
		       tc.user_id AS author_id,
		       COALESCE(au.email, '') AS author_email,
		       COALESCE(au.display_name, '') AS author_display_name
		FROM task_comments tc
		JOIN tasks t ON t.id = tc.task_id
		LEFT JOIN users au ON au.id = tc.user_id
		WHERE tc.content LIKE '[REJECTED]%'
		  AND tc.created_at::date BETWEEN ? AND ?
		  AND t.assigned_to IS NOT NULL
		ORDER BY tc.created_at DESC
	`, from, to).Scan(&reworkDetailRows)

	reworkByUser := map[uint][]perfDomain.DisciplineReworkItem{}
	for _, row := range reworkDetailRows {
		ts := row.CreatedAt.In(ict)
		snippet := row.Content
		runes := []rune(snippet)
		if len(runes) > 220 {
			snippet = string(runes[:220]) + "…"
		}
		rw := perfDomain.DisciplineReworkItem{
			TaskID:            row.TaskID,
			TaskCode:          row.TaskCode,
			TaskTitle:         row.TaskTitle,
			TaskType:          row.TaskType,
			EventDate:         ts.Format("2006-01-02"),
			EventTime:         ts.Format("15:04"),
			CommentSnippet:    snippet,
			AuthorEmail:       row.AuthorEmail,
			AuthorDisplayName: row.AuthorDisplayName,
		}
		if row.AuthorID.Valid && row.AuthorID.Int64 > 0 {
			rw.AuthorID = uint(row.AuthorID.Int64)
		}
		reworkByUser[row.UserID] = append(reworkByUser[row.UserID], rw)
	}
	for uid, items := range reworkByUser {
		sort.Slice(items, func(i, j int) bool {
			if items[i].EventDate != items[j].EventDate {
				return items[i].EventDate > items[j].EventDate
			}
			if items[i].EventTime != items[j].EventTime {
				return items[i].EventTime > items[j].EventTime
			}
			return items[i].TaskCode > items[j].TaskCode
		})
		reworkByUser[uid] = items
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
			UserDisplayName: buildDisplayName(u.FirstName, u.LastName, u.DisplayName, u.Email),
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
			leaveData := leaveIndex[u.ID][d]
			isOnApprovedLeave := leaveData.OnLeave
			isCompanyHoliday := holidayIndex[d]
			if !hasPulse && !isOnApprovedLeave && !isCompanyHoliday {
				missedPulse++
			}
			if att.IsLate {
				totalLate++
			}
			if att.EarlyCheckout {
				totalEarlyOut++
			}
			attendanceStatus := att.Status
			leaveSession := ""
			if isOnApprovedLeave {
				attendanceStatus = "on_leave"
				leaveSession = leaveData.LeaveSession
			} else if isCompanyHoliday && attendanceStatus == "" {
				attendanceStatus = "holiday"
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
				AttendanceStatus:     attendanceStatus,
				LeaveSession:         leaveSession,
			})
		}
		du.TotalTasksClosed = totalClosed
		du.TotalReworks = totalReworks
		du.TotalLoggedHours = float64(totalMins) / 60.0
		du.MissedPulseCount = missedPulse
		du.TotalDeployments = totalDeploys
		du.TotalLateDays = totalLate
		du.TotalEarlyCheckoutDays = totalEarlyOut
		jdi := jobDoneByUser[u.ID]
		if jdi == nil {
			jdi = []perfDomain.DisciplineJobDoneItem{}
		}
		du.JobDoneItems = jdi
		ri := reworkByUser[u.ID]
		if ri == nil {
			ri = []perfDomain.DisciplineReworkItem{}
		}
		du.ReworkItems = ri
		resp.Users = append(resp.Users, du)
	}
	return resp, nil
}

// GetDisciplineDayDetail returns drill-down activity for one user on one date.
func (r *postgresRepo) GetDisciplineDayDetail(userID uint, date string) (*perfDomain.DisciplineDayDetail, error) {
	// Resolve user info
	var u authDomain.User
	if err := r.db.Select("id", "email", "first_name", "last_name", "display_name").First(&u, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	detail := &perfDomain.DisciplineDayDetail{
		UserID:          u.ID,
		UserEmail:       u.Email,
		UserDisplayName: buildDisplayName(u.FirstName, u.LastName, u.DisplayName, u.Email),
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
		TaskID         string
		TaskCode       string
		TaskTitle      string
		Minutes        int
		Description    string
		WorkType       string
		IsTimerSession bool
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

	// Job Done tasks on this date (role-based)
	// - Product Owner: approved READY_FOR_TEST -> WAIT_FOR_DEPLOY
	// - Engineer/Chief Engineer and others: moved IN_PROGRESS -> READY_FOR_TEST on own assigned task
	type ctRow struct {
		TaskID      string
		TaskCode    string
		TaskTitle   string
		StoryPoints int
		TaskType    string
	}
	var ctRows []ctRow
	if u.Role == authDomain.RoleProductOwner {
		r.db.Raw(`
			SELECT t.id::text AS task_id,
			       COALESCE(t.code, '') AS task_code,
			       t.title AS task_title,
			       t.story_points,
			       t.task_type
			FROM task_activity_events tae
			JOIN tasks t ON t.id = tae.task_id
			WHERE tae.action = 'PM_APPROVED_TEST'
			  AND COALESCE(tae.payload->>'from_status', '') = 'READY_FOR_TEST'
			  AND COALESCE(tae.payload->>'to_status', '') = 'WAIT_FOR_DEPLOY'
			  AND tae.actor_id = ?
			  AND tae.created_at::date = ?
			ORDER BY tae.created_at
		`, userID, date).Scan(&ctRows)
	} else {
		r.db.Raw(`
			SELECT t.id::text AS task_id,
			       COALESCE(t.code, '') AS task_code,
			       t.title AS task_title,
			       t.story_points,
			       t.task_type
			FROM task_activity_events tae
			JOIN tasks t ON t.id = tae.task_id
			WHERE t.assigned_to = ?
			  AND t.task_type IN ('TASK', 'BUG')
			  AND tae.created_at::date = ?
			  AND (
			    (tae.action = 'READY_FOR_TEST'
			      AND COALESCE(tae.payload->>'from_status', '') = 'IN_PROGRESS'
			      AND COALESCE(tae.payload->>'to_status', '') = 'READY_FOR_TEST')
			    OR
			    (tae.action = 'STATUS_CHANGED'
			      AND COALESCE(tae.payload->>'from_status', '') = 'IN_PROGRESS'
			      AND COALESCE(tae.payload->>'to_status', '') = 'READY_FOR_TEST')
			  )
			ORDER BY tae.created_at
		`, userID, date).Scan(&ctRows)
	}

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

func (r *postgresRepo) GetDisciplineStartDate() (*perfDomain.DisciplineStartDateResponse, error) {
	var row struct {
		Value string
	}
	err := r.db.Raw("SELECT value FROM app_settings WHERE key = ?", disciplineStartDateKey).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	return &perfDomain.DisciplineStartDateResponse{StartDate: row.Value}, nil
}

func (r *postgresRepo) SetDisciplineStartDate(startDate string) (*perfDomain.DisciplineStartDateResponse, error) {
	err := r.db.Exec(`
		INSERT INTO app_settings (key, value)
		VALUES (?, ?)
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value
	`, disciplineStartDateKey, startDate).Error
	if err != nil {
		return nil, err
	}
	return &perfDomain.DisciplineStartDateResponse{StartDate: startDate}, nil
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
