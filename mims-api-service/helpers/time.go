package helpers

import (
	"time"
)

func GetBangkokTimeNow() (*time.Time, error) {
	bangkok, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, err
	}
	now := time.Now().In(bangkok)
	return &now, nil
}
