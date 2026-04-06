package usecase

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"sort"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	chromepdf "github.com/portnd/the-sentinel-core/internal/core/pdf"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

func (u *sentinelUsecase) GetEpicTimelineData(projectID uuid.UUID) (*domain.EpicTimelineData, error) {
	return u.repo.GetEpicTimelineData(projectID)
}

func (u *sentinelUsecase) GetSprintTimelineData(projectID uuid.UUID) (*domain.SprintTimelineData, error) {
	return u.repo.GetSprintTimelineData(projectID)
}

// ─── 2-Page Client Report:
//   Page 1 (Portrait): งวด/Milestone → Epic → Task list with dates
//   Page 2 (Landscape): Gantt chart (Sprint mode, exact system colours)
// ──────────────────────────────────────────────────────────────────────────────

// p1Task is one deliverable row under an epic.
type p1Task struct {
	Title     string
	StartDate string
	EndDate   string
}

// p1Epic is an epic block (header + tasks) inside a milestone group.
type p1Epic struct {
	Title string
	Tasks []p1Task
}

// p1Milestone is one งวด/delivery group.
type p1Milestone struct {
	Number  int
	Title   string
	DueDate string
	Epics   []p1Epic
	Count   int // total tasks in this milestone
}

// p2GanttCol is a month header column with percentage-based positioning.
type p2GanttCol struct {
	Label    string
	LeftPct  float64
	WidthPct float64
}

// p2GanttRow is one row (sprint header or task) in the Gantt chart.
type p2GanttRow struct {
	IsSprint  bool
	Label     string
	BarLeft   float64 // % from left edge of chart area
	BarWidth  float64 // % width
	BarLabel  string
	HasBar    bool
}

// ganttMonthRow is one task row for the month-timeframe Gantt chart.
// Bar position uses actual dates so bar width reflects real duration (e.g. 2 weeks ≠ full 2 months).
type ganttMonthRow struct {
	Label       string
	EpicTitle   string
	StartMonth  int // 0-based index (kept for sort/display)
	EndMonth    int
	StartDate   string
	EndDate     string
	BarLeftPct  float64 // 0–100: position of bar start as % of chart timeline
	BarWidthPct float64 // 0–100: bar width as % of chart timeline (real duration)
}

// clientReportData is the full payload injected into the HTML template.
type clientReportData struct {
	ProjectName string
	GeneratedAt string
	// Page 1
	MilestoneGroups []p1Milestone
	HasUnassigned   bool
	UnassignedEpics []p1Epic
	// Page 2: day-scale Gantt (sprint)
	GanttCols    []p2GanttCol
	GanttRows    []p2GanttRow
	HasGanttData bool
	// Gantt by month (timeframe เดือน)
	GanttMonthLabels []string
	GanttMonthRows   []ganttMonthRow
	HasGanttMonth    bool
}

func fmtDate(t *time.Time) string {
	if t == nil {
		return "—"
	}
	return t.Format("02 Jan 2006")
}

func fmtDateStr(s string) string {
	if s == "" {
		return "—"
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return s
	}
	return t.Format("02 Jan 2006")
}

// daysBetween returns the number of days between two times (can be negative).
func daysBetween(a, b time.Time) float64 {
	return b.Sub(a).Hours() / 24
}

// truncToMonth returns the first day of the month containing t.
func truncToMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// ExportTimelinePDF generates the 2-page client report via chromedp (same pattern as mims).
func (u *sentinelUsecase) ExportTimelinePDF(projectID uuid.UUID, _ string, templateDir string) ([]byte, string, error) {
	project, err := u.repo.GetProjectByID(projectID, domain.CallerContext{Role: domain.RoleCEO})
	if err != nil {
		return nil, "", fmt.Errorf("project not found: %w", err)
	}

	// ── Fetch data ──────────────────────────────────────────────────────────────
	epicData, err := u.repo.GetEpicTimelineData(projectID)
	if err != nil {
		return nil, "", fmt.Errorf("get epic timeline: %w", err)
	}
	sprintData, err := u.repo.GetSprintTimelineData(projectID)
	if err != nil {
		return nil, "", fmt.Errorf("get sprint timeline: %w", err)
	}
	milestones, _ := u.repo.GetMilestonesByProjectID(projectID) // optional

	// ── Page 1: group tasks by Milestone → Epic ─────────────────────────────────
	// Sort milestones by DueDate ascending
	sort.Slice(milestones, func(i, j int) bool {
		if milestones[i].DueDate == nil {
			return false
		}
		if milestones[j].DueDate == nil {
			return true
		}
		return milestones[i].DueDate.Before(*milestones[j].DueDate)
	})

	// milestoneMap: milestoneID → index in milestones slice (-1 = unassigned)
	milestoneIndexByID := map[string]int{}
	for i, m := range milestones {
		milestoneIndexByID[m.ID.String()] = i
	}

	// Collect all tasks from epics, annotated with their epic title.
	type taskWithEpic struct {
		epicTitle string
		task      domain.Task
	}
	var allTasksWithEpic []taskWithEpic
	for _, e := range epicData.Epics {
		for _, t := range e.Tasks {
			allTasksWithEpic = append(allTasksWithEpic, taskWithEpic{epicTitle: e.Title, task: t})
		}
	}

	// For each milestone, maintain ordered epic groups.
	type epicGroupBuilder struct {
		order []string // epic title order
		tasks map[string][]p1Task
	}
	milestoneBuilders := make([]epicGroupBuilder, len(milestones))
	for i := range milestoneBuilders {
		milestoneBuilders[i].tasks = map[string][]p1Task{}
	}
	var unassignedBuilder epicGroupBuilder
	unassignedBuilder.tasks = map[string][]p1Task{}

	assignTask := func(idx int, epicTitle string, t domain.Task) {
		taskEnd := coalesce(t.EndDate, t.DueAt)
		row := p1Task{Title: t.Title, StartDate: fmtDate(t.StartDate), EndDate: fmtDate(taskEnd)}
		if idx < 0 {
			if _, seen := unassignedBuilder.tasks[epicTitle]; !seen {
				unassignedBuilder.order = append(unassignedBuilder.order, epicTitle)
			}
			unassignedBuilder.tasks[epicTitle] = append(unassignedBuilder.tasks[epicTitle], row)
			return
		}
		b := &milestoneBuilders[idx]
		if _, seen := b.tasks[epicTitle]; !seen {
			b.order = append(b.order, epicTitle)
		}
		b.tasks[epicTitle] = append(b.tasks[epicTitle], row)
	}

	for _, tw := range allTasksWithEpic {
		t := tw.task
		// 1) Direct milestone_id link
		if t.MilestoneID != nil {
			if idx, ok := milestoneIndexByID[t.MilestoneID.String()]; ok {
				assignTask(idx, tw.epicTitle, t)
				continue
			}
		}
		// 2) Date-based: earliest milestone whose DueDate >= task end date
		taskEnd := coalesce(t.EndDate, t.DueAt)
		assigned := false
		if taskEnd != nil {
			for idx, m := range milestones {
				if m.DueDate != nil && !m.DueDate.Before(*taskEnd) {
					assignTask(idx, tw.epicTitle, t)
					assigned = true
					break
				}
			}
		}
		if !assigned {
			assignTask(-1, tw.epicTitle, t)
		}
	}

	// Convert builders → []p1Milestone
	var mGroups []p1Milestone
	for i, m := range milestones {
		b := milestoneBuilders[i]
		var epics []p1Epic
		count := 0
		for _, et := range b.order {
			tasks := b.tasks[et]
			epics = append(epics, p1Epic{Title: et, Tasks: tasks})
			count += len(tasks)
		}
		mGroups = append(mGroups, p1Milestone{
			Number:  i + 1,
			Title:   m.Title,
			DueDate: fmtDate(m.DueDate),
			Epics:   epics,
			Count:   count,
		})
	}
	var unassignedEpics []p1Epic
	for _, et := range unassignedBuilder.order {
		unassignedEpics = append(unassignedEpics, p1Epic{Title: et, Tasks: unassignedBuilder.tasks[et]})
	}

	// ── Gantt by month (timeframe เดือน) ───────────────────────────────────────────
	var ganttMonthLabels []string
	var ganttMonthRows []ganttMonthRow
	hasGanttMonth := false
	var minMonth, maxMonth time.Time
	for _, e := range epicData.Epics {
		for _, t := range e.Tasks {
			start := t.StartDate
			end := coalesce(t.EndDate, t.DueAt)
			if start == nil || end == nil {
				continue
			}
			sM, eM := truncToMonth(*start), truncToMonth(*end)
			if !hasGanttMonth {
				minMonth, maxMonth = sM, eM
				hasGanttMonth = true
			} else {
				if sM.Before(minMonth) {
					minMonth = sM
				}
				if eM.After(maxMonth) {
					maxMonth = eM
				}
			}
		}
	}
	if hasGanttMonth {
		var monthList []time.Time
		for m := minMonth; !m.After(maxMonth); m = m.AddDate(0, 1, 0) {
			monthList = append(monthList, m)
		}
		monthIndex := make(map[string]int)
		for i, m := range monthList {
			monthIndex[m.Format("2006-01")] = i
		}
		for _, lab := range monthList {
			ganttMonthLabels = append(ganttMonthLabels, lab.Format("Jan 06"))
		}
		chartStart := minMonth
		chartEnd := maxMonth.AddDate(0, 1, 0)
		chartDurationDays := daysBetween(chartStart, chartEnd)
		if chartDurationDays <= 0 {
			chartDurationDays = 1
		}
		for _, e := range epicData.Epics {
			for _, t := range e.Tasks {
				start := t.StartDate
				end := coalesce(t.EndDate, t.DueAt)
				if start == nil || end == nil {
					continue
				}
				sM, eM := truncToMonth(*start), truncToMonth(*end)
				si, okS := monthIndex[sM.Format("2006-01")]
				ei, okE := monthIndex[eM.Format("2006-01")]
				if !okS || !okE {
					continue
				}
				if ei < si {
					ei = si
				}
				leftPct := daysBetween(chartStart, *start) / chartDurationDays * 100
				widthPct := daysBetween(*start, *end) / chartDurationDays * 100
				if widthPct <= 0 {
					widthPct = 2
				}
				if leftPct < 0 {
					leftPct = 0
				}
				if leftPct+widthPct > 100 {
					widthPct = 100 - leftPct
				}
				ganttMonthRows = append(ganttMonthRows, ganttMonthRow{
					Label:       t.Title,
					EpicTitle:   e.Title,
					StartMonth:  si,
					EndMonth:    ei,
					StartDate:   fmtDate(start),
					EndDate:     fmtDate(end),
					BarLeftPct:  leftPct,
					BarWidthPct: widthPct,
				})
			}
		}
		sort.Slice(ganttMonthRows, func(i, j int) bool {
			if ganttMonthRows[i].StartMonth != ganttMonthRows[j].StartMonth {
				return ganttMonthRows[i].StartMonth < ganttMonthRows[j].StartMonth
			}
			if ganttMonthRows[i].EpicTitle != ganttMonthRows[j].EpicTitle {
				return ganttMonthRows[i].EpicTitle < ganttMonthRows[j].EpicTitle
			}
			return ganttMonthRows[i].Label < ganttMonthRows[j].Label
		})
	}

	// ── Page 2: build sprint Gantt rows with exact status colours ───────────────
	type rowBuilder struct {
		isSprint bool
		label    string
		start    *time.Time
		end      *time.Time
		status   string
	}
	var rows []rowBuilder
	var chartStart, chartEnd *time.Time

	for _, s := range sprintData.Sprints {
		sStart, sEnd := s.StartDate, s.EndDate
		rows = append(rows, rowBuilder{isSprint: true, label: s.Name, start: sStart, end: sEnd, status: s.Status})
		if sStart != nil && (chartStart == nil || sStart.Before(*chartStart)) {
			t := *sStart
			chartStart = &t
		}
		if sEnd != nil && (chartEnd == nil || sEnd.After(*chartEnd)) {
			t := *sEnd
			chartEnd = &t
		}
		for _, t := range s.Tasks {
			tStart := t.StartDate
			tEnd := coalesce(t.EndDate, t.DueAt)
			rows = append(rows, rowBuilder{isSprint: false, label: t.Title, start: tStart, end: tEnd, status: t.Status})
			if tStart != nil && (chartStart == nil || tStart.Before(*chartStart)) {
				tt := *tStart
				chartStart = &tt
			}
			if tEnd != nil && (chartEnd == nil || tEnd.After(*chartEnd)) {
				tt := *tEnd
				chartEnd = &tt
			}
		}
	}

	hasGantt := chartStart != nil && chartEnd != nil && !chartEnd.Before(*chartStart)
	var cols []p2GanttCol
	var outRows []p2GanttRow

	if hasGantt {
		totalDays := daysBetween(*chartStart, *chartEnd)
		if totalDays < 1 {
			totalDays = 1
		}
		pct := func(d time.Time) float64 {
			return daysBetween(*chartStart, d) / totalDays * 100
		}

		for m := truncToMonth(*chartStart); !m.After(*chartEnd); m = m.AddDate(0, 1, 0) {
			next := m.AddDate(0, 1, 0)
			left := pct(m)
			w := pct(next) - left
			if w < 0 {
				w = 0
			}
			cols = append(cols, p2GanttCol{Label: m.Format("Jan 2006"), LeftPct: left, WidthPct: w})
		}

		for _, r := range rows {
			gr := p2GanttRow{IsSprint: r.isSprint, Label: r.label}
			if r.start != nil && r.end != nil && !r.end.Before(*r.start) {
				left := pct(*r.start)
				width := pct(*r.end) - left
				if width < 1 {
					width = 1
				}
				gr.HasBar = true
				gr.BarLeft = left
				gr.BarWidth = width
				gr.BarLabel = fmt.Sprintf("%s → %s", fmtDateStr(r.start.Format(time.RFC3339)), fmtDateStr(r.end.Format(time.RFC3339)))
			}
			outRows = append(outRows, gr)
		}
	}

	data := clientReportData{
		ProjectName:      project.Name,
		GeneratedAt:      time.Now().Format("02 Jan 2006 15:04"),
		MilestoneGroups:  mGroups,
		HasUnassigned:    len(unassignedEpics) > 0,
		UnassignedEpics:  unassignedEpics,
		GanttCols:        cols,
		GanttRows:        outRows,
		HasGanttData:     hasGantt,
		GanttMonthLabels: ganttMonthLabels,
		GanttMonthRows:   ganttMonthRows,
		HasGanttMonth:    hasGanttMonth,
	}

	tmplPath := templateDir + "timeline_report.html"
	funcMap := template.FuncMap{
		"add":     func(a, b int) int { return a + b },
		"mod":     func(a, b int) int { return a % b },
		"ge":      func(a, b int) bool { return a >= b },
		"le":      func(a, b int) bool { return a <= b },
		"printf2": func(f float64) string { return fmt.Sprintf("%.4f", f) },
	}
	tmpl, err := template.New("timeline_report.html").Funcs(funcMap).ParseFiles(tmplPath)
	if err != nil {
		return nil, "", fmt.Errorf("parse template: %w", err)
	}
	var htmlBuf bytes.Buffer
	if err := tmpl.Execute(&htmlBuf, data); err != nil {
		return nil, "", fmt.Errorf("execute template: %w", err)
	}

	ctx, cancel := chromepdf.NewChromedpContext(context.Background())
	defer cancel()

	var pdfBuf []byte
	if err := chromedp.Run(ctx, chromepdf.PrintToPDF(htmlBuf.String(), &pdfBuf, true)); err != nil {
		return nil, "", fmt.Errorf("chromedp print to pdf: %w", err)
	}

	filename := fmt.Sprintf("project-plan-%s-%s.pdf",
		slugify(project.Name),
		time.Now().Format("20060102"),
	)
	return pdfBuf, filename, nil
}
