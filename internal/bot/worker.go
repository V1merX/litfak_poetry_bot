package bot

import "time"

const (
	sleepDuration time.Duration = 24 * time.Hour
)

func (b *bot) WorkerStart() {
	for {
		time.Sleep(sleepDuration)
		b.log.Debug("worker started")
	}
}