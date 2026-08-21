package application

import "time"

func Backoff(attempt int, base, max time.Duration) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	if base <= 0 {
		base = time.Second
	}
	if max <= 0 {
		max = 5 * time.Minute
	}
	d := base
	for i := 0; i < attempt && d < max; i++ {
		d *= 2
		if d > max {
			d = max
		}
	}
	if d > max {
		d = max
	}
	return d
}
func Jitter(d time.Duration, seed int64) time.Duration {
	if d <= 0 {
		return 0
	}
	if seed < 0 {
		seed = -seed
	}
	return d + time.Duration(seed%int64(d/4+1))
}
