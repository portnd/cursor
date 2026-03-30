package domain

import "errors"

// ErrStandupNotRequiredForRole is returned when CEO or SUPPORT submits a standup (not required).
var ErrStandupNotRequiredForRole = errors.New("daily standup is not required for your role")
