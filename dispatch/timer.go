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

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/alertmanager/nflog"
	"github.com/prometheus/alertmanager/nflog/nflogpb"
	"github.com/prometheus/alertmanager/notify"
)

type (
	TimerFactory func(
		context.Context,
		*time.Timer,
		log.Logger,
		string, // routeKey
		string, // receiver
		time.Duration, // groupInterval
	) Timer
)

func standardTimerFactory(
	_ context.Context,
	t *time.Timer,
	_ log.Logger,
	_ string,
	_ string,
	_ time.Duration,
) Timer {
	return &standardTimer{t}
}

type Timer interface {
	GetC() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}

type standardTimer struct {
	t *time.Timer
}

func (sat *standardTimer) GetC() <-chan time.Time {
	return sat.t.C
}

func (sat *standardTimer) Reset(d time.Duration) bool {
	return sat.t.Reset(d)
}

func (sat *standardTimer) Stop() bool {
	return sat.t.Stop()
}

type syncTimer struct {
	c             chan time.Time
	t             *time.Timer
	nflog         notify.NotificationLog
	routeKey      string
	receiver      string
	logger        log.Logger
	groupInterval time.Duration
}

func NewSyncTimerFactory(
	nflog notify.NotificationLog,
) TimerFactory {
	return func(
		ctx context.Context,
		t *time.Timer,
		l log.Logger,
		routeKey string,
		receiver string,
		groupInterval time.Duration,
	) Timer {
		st := &syncTimer{
			t:             t,
			c:             make(chan time.Time),
			nflog:         nflog,
			routeKey:      routeKey,
			receiver:      receiver,
			logger:        l,
			groupInterval: groupInterval,
		}

		go st.start(ctx)

		return st
	}
}

func (st *syncTimer) start(ctx context.Context) {
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
				if !errors.Is(err, nflog.ErrNotFound) {
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
	entries, err := st.nflog.Query(
		nflog.QGroupKey(st.routeKey),
		nflog.QReceiver(&nflogpb.Receiver{
			GroupName:   st.routeKey,
			Integration: st.receiver,
			Idx:         math.MaxUint32,
		}),
	)
	if err != nil && !errors.Is(err, nflog.ErrNotFound) {
		return nil, fmt.Errorf("error querying log entry: %w", err)
	} else if errors.Is(err, nflog.ErrNotFound) || len(entries) == 0 {
		return nil, nflog.ErrNotFound
	} else if len(entries) > 1 {
		return nil, fmt.Errorf("unexpected entry result size: %d", len(entries))
	}

	ft := entries[0].FlushTime
	if ft == nil || ft.Equal(time.Time{}) {
		return nil, fmt.Errorf("flush time nil or empty")
	}

	return entries[0].FlushTime, nil
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

func (st *syncTimer) GetC() <-chan time.Time {
	return st.c
}

func (st *syncTimer) logFlush(now time.Time) {
	if err := st.nflog.Log(
		&nflogpb.Receiver{
			GroupName:   st.routeKey,
			Integration: st.receiver,
			Idx:         math.MaxUint32,
		},
		st.routeKey,
		nil,
		nil,
		st.groupInterval*2,
		&now,
	); err != nil {
		// log the error and continue
		level.Error(st.logger).Log("msg", "failed to log tick time", "err", err)
	}
}
