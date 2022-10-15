package scanner

import "time"

type Config struct {
	DeepInspection bool
	Timeout        time.Duration
	Threads        int
	Interval       time.Duration

	HostDiscoverClosed bool
}

func DefaultConfig() *Config {
	return &Config{
		DeepInspection:     false,
		Timeout:            time.Second * 2,
		Threads:            800,
		Interval:           time.Millisecond * 300,
		HostDiscoverClosed: false,
	}
}
