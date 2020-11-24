package ecron

// Cron could scheduled by ego or customized scheduler
type Cron interface {
	Run() error
	Stop() error
}
