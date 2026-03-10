package repository

import (
	"time"

	authDomain "github.com/portnd/the-sentinel-core/internal/modules/auth/domain"
	perfDomain "github.com/portnd/the-sentinel-core/internal/modules/performance/domain"
	sentinelDomain "github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
	"gorm.io/gorm"
)

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
	err = r.db.Raw(`
		WITH task_totals AS (
			SELECT t.id, t.estimated_minutes AS est,
				COALESCE(SUM(tl.minutes), 0)::int AS logged
			FROM tasks t
			LEFT JOIN time_logs tl ON tl.task_id = t.id AND tl.user_id = ?
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

// GetDevUserIDsAssignedByPM returns distinct dev user IDs who have at least one task assigned by this PM.
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
	err = r.db.Raw(`
		WITH task_totals AS (
			SELECT t.id, t.estimated_minutes AS est,
				COALESCE(SUM(tl.minutes), 0)::int AS logged
			FROM tasks t
			LEFT JOIN time_logs tl ON tl.task_id = t.id AND tl.user_id = ?
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
	err := r.db.Model(&authDomain.User{}).Where("role = ?", "DEV").Pluck("id", &ids).Error
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
	if err = r.db.Model(&authDomain.User{}).Where("role = ?", "DEV").Pluck("id", &ids).Error; err != nil {
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

func (r *postgresRepo) GetCompanyWideReworkAndTimeAccuracy() (avgReworkPct, avgTimeAccuracyPct float64, err error) {
	var ids []uint
	if err = r.db.Model(&authDomain.User{}).Where("role = ?", "DEV").Pluck("id", &ids).Error; err != nil {
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
 
