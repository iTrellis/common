/*
Copyright © 2021 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package backoff

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Config configures a Backoff
type Config struct {
	MinBackoff time.Duration `yaml:"min_period"`  // start backoff at this level
	MaxBackoff time.Duration `yaml:"max_period"`  // increase exponentially to this level
	MaxRetries int           `yaml:"max_retries"` // give up after this many; zero means infinite retries
}

// Backoff implements exponential backoff with randomized wait times
type Backoff struct {
	cfg          Config
	ctx          context.Context
	numRetries   int
	nextDelayMin time.Duration
	nextDelayMax time.Duration
}

// New creates a Backoff object. Pass a Context that can also terminate the operation.
func New(ctx context.Context, cfg Config) *Backoff {
	return &Backoff{
		cfg:          cfg,
		ctx:          ctx,
		nextDelayMin: cfg.MinBackoff,
		nextDelayMax: doubleDuration(cfg.MinBackoff, cfg.MaxBackoff),
	}
}

// Reset the Backoff back to its initial condition
func (b *Backoff) Reset() {
	b.numRetries = 0
	b.nextDelayMin = b.cfg.MinBackoff
	b.nextDelayMax = doubleDuration(b.cfg.MinBackoff, b.cfg.MaxBackoff)
}

// Ongoing returns true if caller should keep going
func (b *Backoff) Ongoing() bool {
	// Stop if Context has errored or max retry count is exceeded
	return b.ctx.Err() == nil && (b.cfg.MaxRetries == 0 || b.numRetries < b.cfg.MaxRetries)
}

// Err returns the reason for terminating the backoff, or nil if it didn't terminate
func (b *Backoff) Err() error {
	if b.ctx.Err() != nil {
		return b.ctx.Err()
	}
	if b.cfg.MaxRetries != 0 && b.numRetries >= b.cfg.MaxRetries {
		return fmt.Errorf("terminated after %d retries", b.numRetries)
	}
	return nil
}

// NumRetries returns the number of retries so far
func (b *Backoff) NumRetries() int {
	return b.numRetries
}

// Wait sleeps for the backoff time then increases the retry count and backoff time
// Returns immediately if Context is terminated
func (b *Backoff) Wait() {
	// Increase the number of retries and get the next delay
	sleepTime := b.NextDelay()

	if b.Ongoing() {
		select {
		case <-b.ctx.Done():
		case <-time.After(sleepTime):
		}
	}
}

func (b *Backoff) NextDelay() time.Duration {
	b.numRetries++

	// Handle the edge case the min and max have the same value
	// (or due to some misconfig max is < min)
	if b.nextDelayMin >= b.nextDelayMax {
		return b.nextDelayMin
	}

	// Add a jitter within the next exponential backoff range
	sleepTime := b.nextDelayMin + time.Duration(rand.Int63n(int64(b.nextDelayMax-b.nextDelayMin)))

	// Apply the exponential backoff to calculate the next jitter
	// range, unless we've already reached the max
	if b.nextDelayMax < b.cfg.MaxBackoff {
		b.nextDelayMin = doubleDuration(b.nextDelayMin, b.cfg.MaxBackoff)
		b.nextDelayMax = doubleDuration(b.nextDelayMax, b.cfg.MaxBackoff)
	}

	return sleepTime
}

func doubleDuration(value time.Duration, max time.Duration) time.Duration {
	value = value * 2

	if value <= max {
		return value
	}

	return max
}
