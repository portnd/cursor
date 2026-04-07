package domain

import "errors"

var (
	ErrOutsideOffice     = errors.New("check-in denied: outside office GPS radius")
	ErrOutsideOfficeGPS  = errors.New("check-in denied: outside office GPS radius")
	ErrOutsideOfficeIP   = errors.New("check-in denied: client IP is not in allowed office network")
	ErrOutsideOfficeBoth = errors.New("check-in denied: outside office GPS radius and client IP is not in allowed office network")
	ErrNoOfficeConfig    = errors.New("office attendance is not configured")
	ErrAlreadyCheckedIn  = errors.New("already checked in today")
	ErrNotCheckedIn      = errors.New("no check-in today to check out")
	ErrAlreadyCheckedOut = errors.New("already checked out")
	ErrNotWorkDay        = errors.New("today is not a configured office or WFH day")
	ErrForbiddenAdmin    = errors.New("forbidden: CEO or MANAGER only")
	ErrInvalidCursor     = errors.New("invalid or tampered pagination cursor")
	ErrInvalidSchedule   = errors.New("at least one office day or WFH day is required")
	ErrInvalidDateRange  = errors.New("invalid leave date range")
	ErrLeaveNotFound     = errors.New("leave request not found")
	ErrLeaveNotPending   = errors.New("leave request is already reviewed")
	ErrUserNotFound      = errors.New("user not found")
	ErrAttendanceRecordNotFound = errors.New("attendance record not found")
	ErrHalfDayAMLateCheckIn = errors.New("half-day morning leave: check-in must be before 13:00")
	ErrHalfDayPMEarlyCheckOut = errors.New("half-day afternoon leave: check-out allowed from 12:00")
)
