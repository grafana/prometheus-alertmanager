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
	"math"
	"time"

	"github.com/prometheus/alertmanager/flushlog"
	"github.com/prometheus/alertmanager/flushlog/flushlogpb"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// TimerFactory is a function that creates a timer.
type TimerFactory func(
	context.Context,
	*RouteOpts,
	log.Logger,
	uint64,
) Timer

func standardTimerFactory(
	_ context.Context,
	o *RouteOpts,
	_ log.Logger,
	_ uint64,
) Timer {
	return &standardTimer{
		time.NewTimer(o.GroupWait),
		o.GroupInterval,
	}
}

// Timer is the interface for the dispatcher timer.
type Timer interface {
	C() <-chan time.Time
	Reset(time.Time) bool
	Stop(bool) bool
	Flush() bool
}

type standardTimer struct {
	t             *time.Timer
	groupInterval time.Duration
}

func (sat *standardTimer) C() <-chan time.Time {
	return sat.t.C
}

func (sat *standardTimer) Reset(_ time.Time) bool {
	return sat.t.Reset(sat.groupInterval)
}

func (sat *standardTimer) Flush() bool {
	return sat.t.Reset(0)
}

func (sat *standardTimer) Stop(_ bool) bool {
	return sat.t.Stop()
}

type syncTimer struct {
	t                *time.Timer
	flushLog         FlushLog
	position         func() int
	logger           log.Logger
	groupFingerprint uint64
	groupInterval    time.Duration
}

type FlushLog interface {
	Log(groupFingerprint uint64, flushTime time.Time, expiry time.Duration) error
	Query(groupFingerprint uint64) ([]*flushlogpb.FlushLog, error)
	Delete(groupFingerprint uint64) error
}

func NewSyncTimerFactory(
	flushLog FlushLog,
	position func() int,
) TimerFactory {
	return func(
		ctx context.Context,
		o *RouteOpts,
		l log.Logger,
		groupFingerprint uint64,
	) Timer {
		st := &syncTimer{
			t:                time.NewTimer(o.GroupWait),
			flushLog:         flushLog,
			position:         position,
			logger:           l,
			groupInterval:    o.GroupInterval,
			groupFingerprint: groupFingerprint,
		}

		return st
	}
}

func (st *syncTimer) getFirstFlushTime() (*time.Time, error) {
	entries, err := st.flushLog.Query(st.groupFingerprint)
	if err != nil && !errors.Is(err, flushlog.ErrNotFound) {
		return nil, fmt.Errorf("error querying log entry: %w", err)
	} else if errors.Is(err, flushlog.ErrNotFound) || len(entries) == 0 {
		return nil, flushlog.ErrNotFound
	} else if len(entries) > 1 {
		return nil, fmt.Errorf("unexpected entry result size: %d", len(entries))
	}

	ft := entries[0].Timestamp
	if ft.IsZero() {
		return nil, flushlog.ErrNotFound
	}

	return &ft, nil
}

func (st *syncTimer) getNextTick(now time.Time) (time.Duration, bool, error) {
	isPosZero := st.position() == 0

	ft, err := st.getFirstFlushTime()
	if err != nil {
		isNotFound := errors.Is(err, flushlog.ErrNotFound)
		return st.groupInterval, isPosZero && isNotFound, err
	}

	level.Debug(st.logger).Log("msg", "found flush log entry", "flush_time", ft)

	next := ft.Add(time.Duration(st.nextFlushIteration(*ft, now)) * st.groupInterval)
	nextTick := next.Sub(now)
	if !next.After(now) {
		// edge case, now is exactly on the boundary (shouldn't happen)
		// subtract overshoot to maintain interval alignment
		delta := now.Sub(next)
		nextTick = st.groupInterval - delta
	}

	// minimize gossip by logging once per expiry period
	// - expiry is based on the time given by the flush log clock and is set on the mesh struct.
	// - we don't have any of that here, so based on flush time (which is before the flushlog clock time)
	// the idea is to keep the logging frequency low but also ensure that entries that shouldn't expire don't
	// (on expire, the flushlog gets recreated but that introduces drift / de-syncs the flushes)
	closeToExpiry := next.Add(st.groupInterval * 2).After(ft.Add(st.flushLogExpiry()))

	return nextTick, isPosZero && closeToExpiry, nil
}

func (st *syncTimer) Reset(pipelineTime time.Time) bool {
	nextTick, shouldLog, err := st.getNextTick(time.Now())
	if err != nil && !errors.Is(err, flushlog.ErrNotFound) {
		level.Error(st.logger).Log("msg", "failed to calculate next tick", "err", err)
	}

	level.Debug(st.logger).Log("msg", "calculated next tick", "next_tick", nextTick, "should_log", shouldLog)

	if shouldLog {
		st.logFlush(pipelineTime)
	}

	return st.t.Reset(nextTick)
}

func (st *syncTimer) Flush() bool {
	return st.t.Reset(0)
}

func (st *syncTimer) Stop(cleanState bool) bool {
	if st.position() == 0 && cleanState {
		if err := st.flushLog.Delete(st.groupFingerprint); err != nil && !errors.Is(err, flushlog.ErrNotFound) {
			level.Warn(st.logger).Log("msg", "failed to delete flush log entry", "err", err)
		}
	}
	return st.t.Stop()
}

func (st *syncTimer) C() <-chan time.Time {
	return st.t.C
}

func (st *syncTimer) flushLogExpiry() time.Duration {
	// minimum expiry of 24 hours to avoid excessive log churn
	return max(st.groupInterval*2, time.Hour*24)
}

func (st *syncTimer) logFlush(now time.Time) {
	if err := st.flushLog.Log(
		st.groupFingerprint,
		now,
		st.flushLogExpiry(),
	); err != nil {
		// log the error and continue
		level.Error(st.logger).Log("msg", "failed to log tick time", "err", err)
	}
}

func (st *syncTimer) nextFlushIteration(firstFlush, now time.Time) int64 {
	if now.Before(firstFlush) {
		level.Warn(st.logger).Log("msg", "now is before first flush", "first flush", firstFlush, "now", now)
		return 0
	}

	// convert it all to milliseconds
	ns := now.UnixMilli()
	fs := firstFlush.UnixMilli()
	gs := st.groupInterval.Milliseconds()

	return int64(math.Ceil(float64(ns-fs) / float64(gs)))
}
