package scanner

import "github.com/lcvvvv/pool"

type client struct {
	config *Config
	pool   *pool.Pool

	deferFunc func()
}

func (c *client) Stop() {
	c.pool.Stop()
	c.deferFunc()
}

func (c *client) Start() {
	c.pool.Run()
}

func (c *client) Run() {
	c.pool.Run()
}

func (c *client) Defer(f func()) {
	c.deferFunc = f
}

func (c *client) IsDone() bool {
	return c.pool.Done
}

func (c *client) RunningThreads() int {
	return c.pool.RunningThreads()
}

func newConfig(config *Config, threads int) *client {
	return &client{config, pool.New(threads), func() {}}
}
