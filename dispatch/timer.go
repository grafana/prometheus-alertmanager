// Copyright 2018 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dispatch

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/alertmanager/flushlog"
	"github.com/prometheus/alertmanager/flushlog/flushlogpb"
)

// TimerFactory is a function that creates a timer.
type TimerFactory func(context.Context, *time.Timer, log.Logger, time.Duration, uint64) Timer

func standardTimerFactory(
	_ context.Context,
	t *time.Timer,
	_ log.Logger,
	_ time.Duration,
	_ uint64,
) Timer {
	return &standardTimer{t}
}

type Timer interface {
	C() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}

type standardTimer struct {
	t *time.Timer
}

func (sat *standardTimer) C() <-chan time.Time {
	return sat.t.C
}

func (sat *standardTimer) Reset(d time.Duration) bool {
	return sat.t.Reset(d)
}

func (sat *standardTimer) Stop() bool {
	return sat.t.Stop()
}

type syncTimer struct {
	c                chan time.Time
	t                *time.Timer
	flushLog         FlushLog
	logger           log.Logger
	groupFingerprint uint64
	groupInterval    time.Duration
}

type FlushLog interface {
	Log(groupFingerprint uint64, flushTime time.Time) error
	Query(groupFingerprint uint64) ([]*flushlogpb.FlushLog, error)
}

func NewSyncTimerFactory(
	flushLog FlushLog,
) TimerFactory {
	return func(
		ctx context.Context,
		t *time.Timer,
		l log.Logger,
		groupInterval time.Duration,
		groupFingerprint uint64,
	) Timer {
		st := &syncTimer{
			t:                t,
			c:                make(chan time.Time),
			flushLog:         flushLog,
			logger:           l,
			groupInterval:    groupInterval,
			groupFingerprint: groupFingerprint,
		}

		go st.start(ctx)

		return st
	}
}

func (st *syncTimer) start(ctx context.Context) {
	defer close(st.c)

	for {
		select {
		case <-ctx.Done():
			return
		case now, ok := <-st.t.C:
			if !ok { // capture t.Stop()
				return
			}

			wait, err := st.getWaitForNextTick(now)
			if err != nil {
				if !errors.Is(err, flushlog.ErrNotFound) {
					// log the error and continue
					level.Error(st.logger).Log("msg", "failed to calculate next tick", "err", err)
				}
			} else if wait > 0 {
				level.Debug(st.logger).Log("msg", "next tick in the future, waiting..", "wait", wait)
				st.t.Reset(wait)
				continue
			}

			st.logFlush(now)
			st.c <- now
		}
	}
}

func (st *syncTimer) getLastFlushTime() (*time.Time, error) {
	entries, err := st.flushLog.Query(st.groupFingerprint)
	if err != nil && !errors.Is(err, flushlog.ErrNotFound) {
		return nil, fmt.Errorf("error querying log entry: %w", err)
	} else if errors.Is(err, flushlog.ErrNotFound) || len(entries) == 0 {
		return nil, flushlog.ErrNotFound
	} else if len(entries) > 1 {
		return nil, fmt.Errorf("unexpected entry result size: %d", len(entries))
	}

	ft := entries[0].Timestamp
	if ft.Equal(time.Time{}) {
		return nil, fmt.Errorf("flush time nil or empty")
	}

	return &ft, nil
}

func (st *syncTimer) getWaitForNextTick(now time.Time) (time.Duration, error) {
	ft, err := st.getLastFlushTime()
	if err != nil {
		return 0, err
	}

	level.Debug(st.logger).Log("msg", "found flush log entry", "flush_time", ft)

	if next := ft.Add(st.groupInterval); next.After(now) {
		return next.Sub(now), nil
	}

	return 0, nil
}

func (st *syncTimer) Reset(d time.Duration) bool {
	return st.t.Reset(d)
}

func (st *syncTimer) Stop() bool {
	return st.t.Stop()
}

func (st *syncTimer) C() <-chan time.Time {
	return st.c
}

func (st *syncTimer) logFlush(now time.Time) {
	if err := st.flushLog.Log(
		st.groupFingerprint,
		now,
	); err != nil {
		// log the error and continue
		level.Error(st.logger).Log("msg", "failed to log tick time", "err", err)
	}
}
