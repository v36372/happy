package happy

import "time"

// TimeNow return the time now
func TimeNow() time.Time {
	return time.Now().UTC()
}
