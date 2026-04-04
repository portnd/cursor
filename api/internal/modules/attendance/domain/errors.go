package domain

import "errors"

var (
	ErrOutsideOffice     = errors.New("check-in denied: not within office GPS radius or allowed office network")
	ErrNoOfficeConfig    = errors.New("office attendance is not configured")
	ErrAlreadyCheckedIn  = errors.New("already checked in today")
	ErrNotCheckedIn      = errors.New("no check-in today to check out")
	ErrAlreadyCheckedOut = errors.New("already checked out")
	ErrNotWorkDay        = errors.New("today is not a configured office or WFH day")
	ErrForbiddenAdmin    = errors.New("forbidden: CEO or MANAGER only")
	ErrInvalidCursor     = errors.New("invalid or tampered pagination cursor")
	ErrInvalidSchedule   = errors.New("at least one office day or WFH day is required")
)
