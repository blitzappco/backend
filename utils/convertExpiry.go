package utils

import "time"

func ConvertExpiry(expiry string, timedate time.Time) time.Time {
	switch expiry {
	case "1h":
		return timedate.Add(time.Hour)

	case "1.5h":
		return timedate.Add(time.Hour).Add(30 * time.Minute)

	case "2h":
		return timedate.Add(2 * time.Hour)

	case "1d":
		return timedate.Add(24 * time.Hour)

		// expansion late

	default:
		return timedate
	}
}
