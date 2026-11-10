// Copyright 2016 Prometheus Team
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

package nflog

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	pb "github.com/prometheus/alertmanager/nflog/nflogpb"

	"github.com/benbjohnson/clock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
)

func TestLogGC(t *testing.T) {
	mockClock := clock.NewMock()
	now := mockClock.Now()
	// We only care about key names and expiration timestamps.
	newEntry := func(ts time.Time) *pb.MeshEntry {
		return &pb.MeshEntry{
			ExpiresAt: ts,
		}
	}

	l := &Log{
		st: state{
			"a1": newEntry(now),
			"a2": newEntry(now.Add(time.Second)),
			"a3": newEntry(now.Add(-time.Second)),
		},
		clock:   mockClock,
		metrics: newMetrics(nil),
	}
	n, err := l.GC()
	require.NoError(t, err, "unexpected error in garbage collection")
	require.Equal(t, 2, n, "unexpected number of removed entries")

	expected := state{
		"a2": newEntry(now.Add(time.Second)),
	}
	require.Equal(t, expected, l.st, "unexpected state after garbage collection")
}

func TestLogMaxSnapshotSizeBytes(t *testing.T) {
	mockClock := clock.NewMock()
	rcv := func(i int) *pb.Receiver {
		return &pb.Receiver{GroupName: "g", Integration: "test", Idx: uint32(i)}
	}
	newLog := func(limit int) *Log {
		return &Log{
			clock:              mockClock,
			retention:          time.Hour,
			limits:             Limits{MaxSnapshotSizeBytes: func() int { return limit }},
			st:                 state{},
			broadcast:          func([]byte) {},
			isReliableDelivery: func([]byte) bool { return true },
			metrics:            newMetrics(nil),
		}
	}

	// Measure one entry's serialized size with no limit.
	measure := newLog(0)
	require.NoError(t, measure.Log(rcv(1), "gk1", []uint64{1}, nil, 0))
	one, err := measure.MarshalBinary()
	require.NoError(t, err)
	entrySize := len(one)
	require.Positive(t, entrySize)

	// Room for ~2 entries once serialized.
	l := newLog(2*entrySize + entrySize/2)
	mockClock.Add(time.Second)
	require.NoError(t, l.Log(rcv(1), "gk1", []uint64{1}, nil, 0)) // firing, oldest
	mockClock.Add(time.Second)
	require.NoError(t, l.Log(rcv(2), "gk2", nil, []uint64{2}, 0)) // resolved
	mockClock.Add(time.Second)
	require.NoError(t, l.Log(rcv(3), "gk3", []uint64{3}, nil, 0)) // firing, newest

	// The in-memory log is never trimmed, regardless of the limit.
	require.Len(t, l.st, 3)

	// Serializing trims to the limit, keeping firing/newest and dropping the
	// resolved entry first.
	b, err := l.MarshalBinary()
	require.NoError(t, err)
	require.LessOrEqual(t, len(b), l.limits.MaxSnapshotSizeBytes())
	loaded, err := decodeState(bytes.NewReader(b))
	require.NoError(t, err)
	require.Len(t, loaded, 2)
	require.Contains(t, loaded, stateKey("gk3", rcv(3)))    // newest firing kept
	require.Contains(t, loaded, stateKey("gk1", rcv(1)))    // firing kept over resolved
	require.NotContains(t, loaded, stateKey("gk2", rcv(2))) // resolved dropped first
	require.Equal(t, 1.0, testutil.ToFloat64(l.metrics.snapshotEntriesDropped))

	// Serializing does not modify the in-memory log.
	require.Len(t, l.st, 3)

	// Expired entries are excluded from the serialized snapshot (they would be
	// discarded on merge anyway) and do not consume budget.
	le := newLog(1 << 20) // generous limit, so exclusion is due to expiry, not size
	mockClock.Add(time.Second)
	require.NoError(t, le.Log(rcv(1), "live", []uint64{1}, nil, 0))
	le.st[stateKey("expired", rcv(2))] = &pb.MeshEntry{
		Entry:     &pb.Entry{Receiver: rcv(2), GroupKey: []byte("expired"), Timestamp: mockClock.Now(), FiringAlerts: []uint64{2}},
		ExpiresAt: mockClock.Now().Add(-time.Hour),
	}
	require.Len(t, le.st, 2)
	be, err := le.MarshalBinary()
	require.NoError(t, err)
	loadedE, err := decodeState(bytes.NewReader(be))
	require.NoError(t, err)
	require.Len(t, loadedE, 1)
	require.Contains(t, loadedE, stateKey("live", rcv(1)))
	require.NotContains(t, loadedE, stateKey("expired", rcv(2)))

	// With no limit, everything is serialized.
	l0 := newLog(0)
	require.NoError(t, l0.Log(rcv(1), "gk1", []uint64{1}, nil, 0))
	require.NoError(t, l0.Log(rcv(2), "gk2", []uint64{2}, nil, 0))
	require.NoError(t, l0.Log(rcv(3), "gk3", []uint64{3}, nil, 0))
	b0, err := l0.MarshalBinary()
	require.NoError(t, err)
	loaded0, err := decodeState(bytes.NewReader(b0))
	require.NoError(t, err)
	require.Len(t, loaded0, 3)
}

func TestLogSnapshot(t *testing.T) {
	// Check whether storing and loading the snapshot is symmetric.
	mockClock := clock.NewMock()
	now := mockClock.Now().UTC()

	cases := []struct {
		entries []*pb.MeshEntry
	}{
		{
			entries: []*pb.MeshEntry{
				{
					Entry: &pb.Entry{
						GroupKey:  []byte("d8e8fca2dc0f896fd7cb4cb0031ba249"),
						Receiver:  &pb.Receiver{GroupName: "abc", Integration: "test1", Idx: 1},
						GroupHash: []byte("126a8a51b9d1bbd07fddc65819a542c3"),
						Resolved:  false,
						Timestamp: now,
					},
					ExpiresAt: now,
				}, {
					Entry: &pb.Entry{
						GroupKey:  []byte("d8e8fca2dc0f8abce7cb4cb0031ba249"),
						Receiver:  &pb.Receiver{GroupName: "def", Integration: "test2", Idx: 29},
						GroupHash: []byte("122c2331b9d1bbd07fddc65819a542c3"),
						Resolved:  true,
						Timestamp: now,
					},
					ExpiresAt: now,
				}, {
					Entry: &pb.Entry{
						GroupKey:  []byte("aaaaaca2dc0f896fd7cb4cb0031ba249"),
						Receiver:  &pb.Receiver{GroupName: "ghi", Integration: "test3", Idx: 0},
						GroupHash: []byte("126a8a51b9d1bbd07fddc6e3e3e542c3"),
						Resolved:  false,
						Timestamp: now,
					},
					ExpiresAt: now,
				},
			},
		},
	}

	for _, c := range cases {
		f, err := os.CreateTemp("", "snapshot")
		require.NoError(t, err, "creating temp file failed")

		l1 := &Log{
			st:      state{},
			metrics: newMetrics(nil),
		}
		// Setup internal state manually.
		for _, e := range c.entries {
			l1.st[stateKey(string(e.Entry.GroupKey), e.Entry.Receiver)] = e
		}
		_, err = l1.Snapshot(f)
		require.NoError(t, err, "creating snapshot failed")
		require.NoError(t, f.Close(), "closing snapshot file failed")

		f, err = os.Open(f.Name())
		require.NoError(t, err, "opening snapshot file failed")

		// Check again against new nlog instance.
		l2 := &Log{}
		err = l2.loadSnapshot(f)
		require.NoError(t, err, "error loading snapshot")
		require.Equal(t, l1.st, l2.st, "state after loading snapshot did not match snapshotted state")

		require.NoError(t, f.Close(), "closing snapshot file failed")
	}
}

func TestWithMaintenance_SupportsCustomCallback(t *testing.T) {
	f, err := os.CreateTemp("", "snapshot")
	require.NoError(t, err, "creating temp file failed")
	stopc := make(chan struct{})
	reg := prometheus.NewPedanticRegistry()
	opts := Options{
		Metrics:      reg,
		SnapshotFile: f.Name(),
	}

	l, err := New(opts)
	clock := clock.NewMock()
	l.clock = clock
	require.NoError(t, err)

	var calls atomic.Int32
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		l.Maintenance(100*time.Millisecond, f.Name(), stopc, func() (int64, error) {
			calls.Add(1)
			return 0, nil
		})
	}()
	gosched()

	// Before the first tick, no maintenance executed.
	clock.Add(99 * time.Millisecond)
	require.EqualValues(t, 0, calls.Load())

	// Tick once.
	clock.Add(1 * time.Millisecond)
	require.EqualValues(t, 1, calls.Load())

	// Stop the maintenance loop. We should get exactly one more execution of the maintenance func.
	close(stopc)
	wg.Wait()

	require.EqualValues(t, 2, calls.Load())
	// Check the maintenance metrics.
	require.NoError(t, testutil.GatherAndCompare(reg, bytes.NewBufferString(`
# HELP alertmanager_nflog_maintenance_errors_total How many maintenances were executed for the notification log that failed.
# TYPE alertmanager_nflog_maintenance_errors_total counter
alertmanager_nflog_maintenance_errors_total 0
# HELP alertmanager_nflog_maintenance_total How many maintenances were executed for the notification log.
# TYPE alertmanager_nflog_maintenance_total counter
alertmanager_nflog_maintenance_total 2
`), "alertmanager_nflog_maintenance_total", "alertmanager_nflog_maintenance_errors_total"))
}

func TestReplaceFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "replace_file")
	require.NoError(t, err, "creating temp dir failed")

	origFilename := filepath.Join(dir, "testfile")

	of, err := os.Create(origFilename)
	require.NoError(t, err, "creating file failed")

	nf, err := openReplace(origFilename)
	require.NoError(t, err, "opening replacement file failed")

	_, err = nf.Write([]byte("test"))
	require.NoError(t, err, "writing replace file failed")

	require.NotEqual(t, nf.Name(), of.Name(), "replacement file must have different name while editing")
	require.NoError(t, nf.Close(), "closing replacement file failed")
	require.NoError(t, of.Close(), "closing original file failed")

	ofr, err := os.Open(origFilename)
	require.NoError(t, err, "opening original file failed")
	defer ofr.Close()

	res, err := io.ReadAll(ofr)
	require.NoError(t, err, "reading original file failed")
	require.Equal(t, "test", string(res), "unexpected file contents")
}

func TestStateMerge(t *testing.T) {
	mockClock := clock.NewMock()
	now := mockClock.Now()

	// We only care about key names and timestamps for the
	// merging logic.
	newEntry := func(name string, ts, exp time.Time) *pb.MeshEntry {
		return &pb.MeshEntry{
			Entry: &pb.Entry{
				Timestamp: ts,
				GroupKey:  []byte("key"),
				Receiver: &pb.Receiver{
					GroupName:   name,
					Idx:         1,
					Integration: "integr",
				},
			},
			ExpiresAt: exp,
		}
	}

	exp := now.Add(time.Minute)

	cases := []struct {
		a, b  state
		final state
	}{
		{
			a: state{
				"key:a1/integr/1": newEntry("a1", now, exp),
				"key:a2/integr/1": newEntry("a2", now, exp),
				"key:a3/integr/1": newEntry("a3", now, exp),
			},
			b: state{
				"key:b1/integr/1": newEntry("b1", now, exp),                                          // new key, should be added
				"key:b2/integr/1": newEntry("b2", now.Add(-time.Minute), now.Add(-time.Millisecond)), // new key, expired, should not be added
				"key:a2/integr/1": newEntry("a2", now.Add(-time.Minute), exp),                        // older timestamp, should be dropped
				"key:a3/integr/1": newEntry("a3", now.Add(time.Minute), exp),                         // newer timestamp, should overwrite
			},
			final: state{
				"key:a1/integr/1": newEntry("a1", now, exp),
				"key:a2/integr/1": newEntry("a2", now, exp),
				"key:a3/integr/1": newEntry("a3", now.Add(time.Minute), exp),
				"key:b1/integr/1": newEntry("b1", now, exp),
			},
		},
	}

	for _, c := range cases {
		ca, cb := c.a.clone(), c.b.clone()

		res := c.a.clone()
		for _, e := range cb {
			res.merge(e, now)
		}
		require.Equal(t, c.final, res, "Merge result should match expectation")
		require.Equal(t, c.b, cb, "Merged state should remain unmodified")
		require.NotEqual(t, c.final, ca, "Merge should not change original state")
	}
}

func TestStateDataCoding(t *testing.T) {
	// Check whether encoding and decoding the data is symmetric.
	mockClock := clock.NewMock()
	now := mockClock.Now().UTC()

	cases := []struct {
		entries []*pb.MeshEntry
	}{
		{
			entries: []*pb.MeshEntry{
				{
					Entry: &pb.Entry{
						GroupKey:  []byte("d8e8fca2dc0f896fd7cb4cb0031ba249"),
						Receiver:  &pb.Receiver{GroupName: "abc", Integration: "test1", Idx: 1},
						GroupHash: []byte("126a8a51b9d1bbd07fddc65819a542c3"),
						Resolved:  false,
						Timestamp: now,
					},
					ExpiresAt: now,
				}, {
					Entry: &pb.Entry{
						GroupKey:  []byte("d8e8fca2dc0f8abce7cb4cb0031ba249"),
						Receiver:  &pb.Receiver{GroupName: "def", Integration: "test2", Idx: 29},
						GroupHash: []byte("122c2331b9d1bbd07fddc65819a542c3"),
						Resolved:  true,
						Timestamp: now,
					},
					ExpiresAt: now,
				}, {
					Entry: &pb.Entry{
						GroupKey:  []byte("aaaaaca2dc0f896fd7cb4cb0031ba249"),
						Receiver:  &pb.Receiver{GroupName: "ghi", Integration: "test3", Idx: 0},
						GroupHash: []byte("126a8a51b9d1bbd07fddc6e3e3e542c3"),
						Resolved:  false,
						Timestamp: now,
					},
					ExpiresAt: now,
				},
			},
		},
	}

	for _, c := range cases {
		// Create gossip data from input.
		in := state{}
		for _, e := range c.entries {
			in[stateKey(string(e.Entry.GroupKey), e.Entry.Receiver)] = e
		}
		msg, err := in.MarshalBinary()
		require.NoError(t, err)

		out, err := decodeState(bytes.NewReader(msg))
		require.NoError(t, err, "decoding message failed")

		require.Equal(t, in, out, "decoded data doesn't match encoded data")
	}
}

func TestQuery(t *testing.T) {
	opts := Options{Retention: time.Second}
	nl, err := New(opts)
	if err != nil {
		require.NoError(t, err, "constructing nflog failed")
	}

	recv := new(pb.Receiver)

	// no key param
	_, err = nl.Query(QGroupKey("key"))
	require.EqualError(t, err, "no query parameters specified")

	// no recv param
	_, err = nl.Query(QReceiver(recv))
	require.EqualError(t, err, "no query parameters specified")

	// no entry
	_, err = nl.Query(QGroupKey("nonexistentkey"), QReceiver(recv))
	require.EqualError(t, err, "not found")

	// existing entry
	firingAlerts := []uint64{1, 2, 3}
	resolvedAlerts := []uint64{4, 5}

	err = nl.Log(recv, "key", firingAlerts, resolvedAlerts, 0)
	require.NoError(t, err, "logging notification failed")

	entries, err := nl.Query(QGroupKey("key"), QReceiver(recv))
	require.NoError(t, err, "querying nflog failed")
	entry := entries[0]
	require.Equal(t, firingAlerts, entry.FiringAlerts)
	require.Equal(t, resolvedAlerts, entry.ResolvedAlerts)
}

func TestStateDecodingError(t *testing.T) {
	// Check whether decoding copes with erroneous data.
	s := state{"": &pb.MeshEntry{}}

	msg, err := s.MarshalBinary()
	require.NoError(t, err)

	_, err = decodeState(bytes.NewReader(msg))
	require.Equal(t, ErrInvalidState, err)
}

// runtime.Gosched() does not "suspend" the current goroutine so there's no guarantee that the main goroutine won't
// be able to continue. For more see https://pkg.go.dev/runtime#Gosched.
func gosched() {
	time.Sleep(1 * time.Millisecond)
}

func TestLogRefreshExpiry(t *testing.T) {
	rcv := &pb.Receiver{GroupName: "g", Integration: "test", Idx: 1}
	const (
		gkey      = "group"
		retention = time.Hour
	)

	newLog := func(clk clock.Clock) (*Log, *atomic.Int32) {
		var broadcasts atomic.Int32
		return &Log{
			clock:              clk,
			retention:          retention,
			st:                 state{},
			metrics:            newMetrics(nil),
			broadcast:          func([]byte) { broadcasts.Inc() },
			isReliableDelivery: func([]byte) bool { return true },
		}, &broadcasts
	}
	expiresAt := func(l *Log) time.Time {
		l.mtx.RLock()
		defer l.mtx.RUnlock()
		return l.st[stateKey(gkey, rcv)].ExpiresAt
	}

	t.Run("missing entry is not created", func(t *testing.T) {
		l, broadcasts := newLog(clock.NewMock())
		l.RefreshExpiry(rcv, gkey, l.now(), time.Hour)
		require.Empty(t, l.st)
		require.Zero(t, broadcasts.Load())
	})

	t.Run("entry is extended and otherwise unchanged", func(t *testing.T) {
		mockClock := clock.NewMock()
		l, broadcasts := newLog(mockClock)
		require.NoError(t, l.Log(rcv, gkey, []uint64{1, 2}, []uint64{3}, 10*time.Minute))
		entries, err := l.Query(QReceiver(rcv), QGroupKey(gkey))
		require.NoError(t, err)
		require.Len(t, entries, 1)
		original := entries[0]
		broadcastsBefore := broadcasts.Load()

		mockClock.Add(8 * time.Minute)
		now := mockClock.Now()
		l.RefreshExpiry(rcv, gkey, now, 10*time.Minute)
		require.Equal(t, now.Add(10*time.Minute), expiresAt(l))

		entries, err = l.Query(QReceiver(rcv), QGroupKey(gkey))
		require.NoError(t, err)
		require.Len(t, entries, 1)
		require.Equal(t, original, entries[0])
		require.Equal(t, []uint64{1, 2}, entries[0].FiringAlerts)
		require.Equal(t, []uint64{3}, entries[0].ResolvedAlerts)
		require.Equal(t, broadcastsBefore, broadcasts.Load(), "refresh must not be broadcast")
	})

	t.Run("expiry is capped by retention", func(t *testing.T) {
		mockClock := clock.NewMock()
		l, _ := newLog(mockClock)
		require.NoError(t, l.Log(rcv, gkey, []uint64{1}, nil, time.Minute))

		mockClock.Add(30 * time.Second)
		now := mockClock.Now()
		l.RefreshExpiry(rcv, gkey, now, 100*time.Hour)
		require.Equal(t, now.Add(retention), expiresAt(l))
	})

	t.Run("expiry is never shortened", func(t *testing.T) {
		mockClock := clock.NewMock()
		l, _ := newLog(mockClock)
		require.NoError(t, l.Log(rcv, gkey, []uint64{1}, nil, 30*time.Minute))
		before := expiresAt(l)

		l.RefreshExpiry(rcv, gkey, mockClock.Now(), time.Minute)
		require.Equal(t, before, expiresAt(l))
	})

	t.Run("result depends on the passed time and not on the local clock", func(t *testing.T) {
		now := clock.NewMock().Now()
		var results []time.Time
		for _, skew := range []time.Duration{0, 3 * time.Second, -2 * time.Second} {
			mockClock := clock.NewMock()
			l, _ := newLog(mockClock)
			require.NoError(t, l.Log(rcv, gkey, []uint64{1}, nil, 10*time.Minute))
			mockClock.Add(5*time.Minute + skew)

			l.RefreshExpiry(rcv, gkey, now.Add(5*time.Minute), 10*time.Minute)
			results = append(results, expiresAt(l))
		}
		require.Equal(t, results[0], results[1])
		require.Equal(t, results[0], results[2])
	})

	t.Run("entry survives GC when refreshed on every flush", func(t *testing.T) {
		const (
			groupInterval  = time.Minute
			repeatInterval = 3 * time.Hour // longer than retention
		)
		expiry := 2 * repeatInterval

		run := func(refresh bool) int {
			mockClock := clock.NewMock()
			l, _ := newLog(mockClock)
			require.NoError(t, l.Log(rcv, gkey, []uint64{1}, nil, expiry))
			for i := 0; i < 180; i++ {
				mockClock.Add(groupInterval)
				if refresh {
					l.RefreshExpiry(rcv, gkey, mockClock.Now(), expiry)
				}
				_, err := l.GC()
				require.NoError(t, err)
			}
			return len(l.st)
		}
		require.Equal(t, 0, run(false), "without refresh the entry must expire")
		require.Equal(t, 1, run(true), "with refresh the entry must stay")
	})

	t.Run("merge of a newer peer entry replaces the refreshed one", func(t *testing.T) {
		mockClock := clock.NewMock()
		l, _ := newLog(mockClock)
		require.NoError(t, l.Log(rcv, gkey, []uint64{1}, nil, 10*time.Minute))
		mockClock.Add(9 * time.Minute)
		l.RefreshExpiry(rcv, gkey, mockClock.Now(), 10*time.Minute)

		peerEntry := &pb.MeshEntry{
			Entry: &pb.Entry{
				Receiver:     rcv,
				GroupKey:     []byte(gkey),
				Timestamp:    mockClock.Now().Add(time.Second),
				FiringAlerts: []uint64{7},
			},
			ExpiresAt: mockClock.Now().Add(time.Hour),
		}
		b, err := marshalMeshEntry(peerEntry)
		require.NoError(t, err)
		require.NoError(t, l.Merge(b))

		entries, err := l.Query(QReceiver(rcv), QGroupKey(gkey))
		require.NoError(t, err)
		require.Equal(t, []uint64{7}, entries[0].FiringAlerts)
	})
}
