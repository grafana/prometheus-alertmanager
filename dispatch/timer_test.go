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
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/flushlog"
	"github.com/prometheus/alertmanager/flushlog/flushlogpb"
	"github.com/prometheus/alertmanager/notify"
	"github.com/prometheus/alertmanager/provider/mem"
	"github.com/prometheus/alertmanager/types"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/require"
)

func TestSyncTimer(t *testing.T) {
	now := time.Now()

	buf := &logBuf{t: t, b: []string{}}
	logger := log.NewJSONLogger(buf)
	marker := types.NewMarker(prometheus.NewRegistry())
	alerts, err := mem.NewAlerts(context.Background(), marker, time.Hour, nil, logger, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer alerts.Close()
	flushlog := &mockLog{
		t: t,
		logCalls: []mockLogCall{
			{
				expGroupFingerprint: 13705263069144098434,
				expExpiry:           time.Millisecond * 20,
			},
		},
		queryCalls: []mockQueryCall{
			{ // first call to query doesn't find state
				expGroupFingerprint: 13705263069144098434,
				err:                 flushlog.ErrNotFound,
			},
			{
				expGroupFingerprint: 13705263069144098434,
				res: []*flushlogpb.FlushLog{
					{
						GroupFingerprint: 13705263069144098434,
						Timestamp:        now,
					},
				},
			},
			{
				expGroupFingerprint: 13705263069144098434,
				res: []*flushlogpb.FlushLog{
					{
						GroupFingerprint: 13705263069144098434,
						Timestamp:        now.Add(time.Millisecond * 80),
					},
				},
			},
		},
		deleteCalls: []mockDeleteCall{
			{
				expGroupFingerprint: 13705263069144098434,
				err:                 nil,
			},
		},
	}

	// wait for 3 notification cycles
	n := 3
	nfC := make(chan struct{}, n)
	stage := &pubStage{nfC}

	dispatcher, err := newTestDispatcher(time.Millisecond*10, alerts, stage, marker, logger, NewSyncTimerFactory(flushlog, func() int { return 0 }))
	if err != nil {
		t.Fatal(err)
	}

	go dispatcher.Run()

	alerts.Put(newAlert(model.LabelSet{"alertname": "TestingAlert"}))

	var i int
	for {
		select {
		case <-nfC:
			i++
			if i == n {
				dispatcher.Stop() // ensure we stop the dispatcher before making assertions to mitigate flakiness

				flushlog.requireCalls() // make sure the call stacks are empty
				buf.requireLogs(        // require the logs in order
					"Received alert",
					// flush ticks
					"calculated next tick",  // 1.1. doesn't find an entry (no logs)
					"flushing",              // 1.2. logs the flush
					"found flush log entry", // 2.1. finds an entry, so calculates next tick
					"calculated next tick",  // 2.2. finds an entry from the past, so flushes immediately
					"flushing",              // 2.3. logs the flush
					"found flush log entry", // 3.1. finds an entry, so calculates next tick
					"calculated next tick",  // 3.2. calculates next tick
					"flushing",              // 3.3. flushes the entry
				)
				return
			}
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for dispatcher to finish")
		}
	}
}

type mockLog struct {
	t     *testing.T
	bench bool

	mtx         sync.Mutex
	queryCalls  []mockQueryCall
	logCalls    []mockLogCall
	deleteCalls []mockDeleteCall
}

type mockQueryCall struct {
	expGroupFingerprint uint64

	res []*flushlogpb.FlushLog
	err error
}

func (m *mockLog) Query(groupFingerprint uint64) ([]*flushlogpb.FlushLog, error) {
	if m.bench {
		// we want to measure the impact of the extra go routine in the sync timer, so always return a time in the past
		return []*flushlogpb.FlushLog{
			{
				GroupFingerprint: groupFingerprint,
				Timestamp:        time.Now().Add(-time.Minute * 10),
			},
		}, nil
	}

	var c mockQueryCall
	func() {
		m.mtx.Lock()
		defer m.mtx.Unlock()
		if len(m.queryCalls) == 0 {
			require.FailNow(m.t, "no query calls")
		}
		c = m.queryCalls[0]
		m.queryCalls = m.queryCalls[1:]
	}()

	require.Equal(m.t, c.expGroupFingerprint, groupFingerprint)

	return c.res, c.err
}

type mockLogCall struct {
	expGroupFingerprint uint64
	expExpiry           time.Duration
	err                 error
}

func (m *mockLog) Log(groupFingerprint uint64, timestamp time.Time, expiry time.Duration) error {
	if m.bench {
		return nil
	}

	var c mockLogCall
	func() {
		m.mtx.Lock()
		defer m.mtx.Unlock()
		if len(m.logCalls) == 0 {
			require.FailNow(m.t, "no log calls")
		}
		c = m.logCalls[0]
		m.logCalls = m.logCalls[1:]
	}()

	require.Equal(m.t, c.expGroupFingerprint, groupFingerprint)
	require.Equal(m.t, c.expExpiry, expiry)

	return c.err
}

type mockDeleteCall struct {
	expGroupFingerprint uint64
	err                 error
}

func (m *mockLog) Delete(groupFingerprint uint64) error {
	if m.bench {
		return nil
	}

	var c mockDeleteCall
	func() {
		m.mtx.Lock()
		defer m.mtx.Unlock()
		if len(m.deleteCalls) == 0 {
			require.FailNow(m.t, "no delete calls")
		}
		c = m.deleteCalls[0]
		m.deleteCalls = m.deleteCalls[1:]
	}()

	require.Equal(m.t, c.expGroupFingerprint, groupFingerprint)

	return c.err
}

func (m *mockLog) requireCalls() {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	require.Empty(m.t, m.queryCalls)
	require.Empty(m.t, m.logCalls)
	require.Empty(m.t, m.deleteCalls)
}

type logBuf struct {
	t *testing.T
	m sync.Mutex
	b []string
}

func (pb *logBuf) Write(b []byte) (int, error) {
	pb.m.Lock()
	defer pb.m.Unlock()

	var l struct {
		Msg string `json:"msg"`
		Err string `json:"err,omitempty"`
	}
	if err := json.Unmarshal(b, &l); err != nil {
		return 0, err
	}
	if l.Err != "" {
		fmt.Println("error in log:", l.Err)
	}

	pb.b = append(pb.b, l.Msg)
	return len(b), nil
}

func (pb *logBuf) requireLogs(expLogs ...string) {
	pb.m.Lock()
	defer pb.m.Unlock()
	require.Equal(pb.t, expLogs, pb.b)
}

func BenchmarkSyncTimer(b *testing.B) {
	for b.Loop() {
		benchTimer(func(m *mockLog) TimerFactory { return NewSyncTimerFactory(m, func() int { return 0 }) }, b)
	}
}

func BenchmarkStdTimer(b *testing.B) {
	for b.Loop() {
		benchTimer(func(*mockLog) TimerFactory { return standardTimerFactory }, b)
	}
}

func benchTimer(timerFactoryBuilder func(*mockLog) TimerFactory, b *testing.B) {
	logger := log.NewNopLogger()
	marker := types.NewMarker(prometheus.NewRegistry())
	alerts, err := mem.NewAlerts(context.Background(), marker, time.Hour, nil, logger, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer alerts.Close()

	n := rand.Intn(10) + 1
	nfC := make(chan struct{}, n)
	stage := &pubStage{nfC}

	flushlog := &mockLog{
		bench:      true,
		logCalls:   []mockLogCall{},
		queryCalls: []mockQueryCall{},
	}

	dispatcher, err := newTestDispatcher(time.Minute*1, alerts, stage, marker, logger, timerFactoryBuilder(flushlog))
	if err != nil {
		b.Fatal(err)
	}

	as := make([]*types.Alert, 0, n)
	for i := 0; i < n; i++ {
		as = append(as, newAlert(model.LabelSet{"alertname": model.LabelValue(fmt.Sprintf("TestingAlert_%d", i))}))
	}

	go dispatcher.Run()

	for i := 0; i < n; i++ {
		alerts.Put(as...)
	}

	var i int
	for {
		select {
		case <-nfC:
			i++
			if i == n {
				dispatcher.Stop() // ensure we stop the dispatcher before making assertions to mitigate flakiness
				return
			}
		case <-time.After(20 * time.Second):
			b.Fatal("timed out waiting for dispatcher to finish")
		}
	}
}

func newTestDispatcher(
	groupInterval time.Duration,
	alerts *mem.Alerts,
	stage notify.Stage,
	marker types.Marker,
	logger log.Logger,
	timerFactory TimerFactory,
) (*Dispatcher, error) {
	confData := fmt.Sprintf(`receivers:
- name: 'testing'

route:
  group_by: ['alertname']
  group_wait: 10ms
  group_interval: %s
  receiver: 'testing'
  routes: []`, groupInterval)
	conf, err := config.Load(confData)
	if err != nil {
		return nil, err
	}

	route := NewRoute(conf.Route, nil)
	timeout := func(d time.Duration) time.Duration { return time.Duration(0) }

	return NewDispatcher(alerts, route, stage, marker, timeout, nil, logger, NewDispatcherMetrics(false, prometheus.NewRegistry()), timerFactory), nil
}

type pubStage struct {
	c chan struct{}
}

func (b *pubStage) Exec(ctx context.Context, l log.Logger, alerts ...*types.Alert) (context.Context, []*types.Alert, error) {
	b.c <- struct{}{}
	return ctx, nil, nil
}
