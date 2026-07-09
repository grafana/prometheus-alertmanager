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

// Package nflog implements a garbage-collected and snapshottable append-only log of
// active/resolved notifications. Each log entry stores the active/resolved state,
// the notified receiver, and a hash digest of the notification's identifying contents.
// The log can be queried along different parameters.
package nflog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/matttproud/golang_protobuf_extensions/pbutil"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/alertmanager/cluster/clusterutil"
	pb "github.com/prometheus/alertmanager/nflog/nflogpb"
)

// ErrNotFound is returned for empty query results.
var ErrNotFound = errors.New("not found")

// ErrInvalidState is returned if the state isn't valid.
var ErrInvalidState = errors.New("invalid state")

// query currently allows filtering by and/or receiver group key.
// It is configured via QueryParameter functions.
//
// TODO(fabxc): Future versions could allow querying a certain receiver,
// group or a given time interval.
type query struct {
	recv     *pb.Receiver
	groupKey string
}

// QueryParam is a function that modifies a query to incorporate
// a set of parameters. Returns an error for invalid or conflicting
// parameters.
type QueryParam func(*query) error

// QReceiver adds a receiver parameter to a query.
func QReceiver(r *pb.Receiver) QueryParam {
	return func(q *query) error {
		q.recv = r
		return nil
	}
}

// QGroupKey adds a group key as querying argument.
func QGroupKey(gk string) QueryParam {
	return func(q *query) error {
		q.groupKey = gk
		return nil
	}
}

// Log holds the notification log state for alerts that have been notified.
type Log struct {
	clock clock.Clock

	logger    log.Logger
	metrics   *metrics
	retention time.Duration
	limits    Limits

	// For now we only store the most recently added log entry.
	// The key is a serialized concatenation of group key and receiver.
	mtx sync.RWMutex
	st  state
	// size is the running sum of the marshaled sizes of all entries in st, used to
	// enforce Limits.MaxSizeBytes without recomputing the whole state on each write.
	size               int
	broadcast          func([]byte)
	isReliableDelivery func([]byte) bool
}

// Limits contains the limits for the notification log.
type Limits struct {
	// MaxSizeBytes limits the total marshaled size of the notification log in
	// bytes. If nil, or the returned value is negative or 0, there is no limit.
	//
	// When recording a notification would push the log over the limit, existing
	// entries are evicted until it fits again: resolved entries first, then oldest
	// by last-notification time. Because the notification has already been sent by
	// the time it is logged, an evicted group is at worst re-notified later (a
	// duplicate), never missed.
	MaxSizeBytes func() int
}

// MaintenanceFunc represents the function to run as part of the periodic maintenance for the nflog.
// It returns the size of the snapshot taken or an error if it failed.
type MaintenanceFunc func() (int64, error)

type metrics struct {
	gcDuration              prometheus.Summary
	snapshotDuration        prometheus.Summary
	snapshotSize            prometheus.Gauge
	queriesTotal            prometheus.Counter
	queryErrorsTotal        prometheus.Counter
	queryDuration           prometheus.Histogram
	propagatedMessagesTotal prometheus.Counter
	maintenanceTotal        prometheus.Counter
	maintenanceErrorsTotal  prometheus.Counter
	size                    prometheus.Gauge
	entriesEvictedTotal     prometheus.Counter
}

func newMetrics(r prometheus.Registerer) *metrics {
	m := &metrics{}

	m.gcDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "alertmanager_nflog_gc_duration_seconds",
		Help:       "Duration of the last notification log garbage collection cycle.",
		Objectives: map[float64]float64{},
	})
	m.snapshotDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "alertmanager_nflog_snapshot_duration_seconds",
		Help:       "Duration of the last notification log snapshot.",
		Objectives: map[float64]float64{},
	})
	m.snapshotSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "alertmanager_nflog_snapshot_size_bytes",
		Help: "Size of the last notification log snapshot in bytes.",
	})
	m.maintenanceTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "alertmanager_nflog_maintenance_total",
		Help: "How many maintenances were executed for the notification log.",
	})
	m.maintenanceErrorsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "alertmanager_nflog_maintenance_errors_total",
		Help: "How many maintenances were executed for the notification log that failed.",
	})
	m.queriesTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "alertmanager_nflog_queries_total",
		Help: "Number of notification log queries were received.",
	})
	m.queryErrorsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "alertmanager_nflog_query_errors_total",
		Help: "Number notification log received queries that failed.",
	})
	m.queryDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:                            "alertmanager_nflog_query_duration_seconds",
		Help:                            "Duration of notification log query evaluation.",
		Buckets:                         prometheus.DefBuckets,
		NativeHistogramBucketFactor:     1.1,
		NativeHistogramMaxBucketNumber:  100,
		NativeHistogramMinResetDuration: 1 * time.Hour,
	})
	m.propagatedMessagesTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "alertmanager_nflog_gossip_messages_propagated_total",
		Help: "Number of received gossip messages that have been further gossiped.",
	})
	m.size = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "alertmanager_nflog_size_bytes",
		Help: "Current total marshaled size of the notification log in bytes.",
	})
	m.entriesEvictedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "alertmanager_nflog_entries_evicted_total",
		Help: "Number of notification log entries evicted (oldest first) to keep the log within the configured maximum size.",
	})

	if r != nil {
		r.MustRegister(
			m.gcDuration,
			m.snapshotDuration,
			m.snapshotSize,
			m.queriesTotal,
			m.queryErrorsTotal,
			m.queryDuration,
			m.propagatedMessagesTotal,
			m.maintenanceTotal,
			m.maintenanceErrorsTotal,
			m.size,
			m.entriesEvictedTotal,
		)
	}
	return m
}

type state map[string]*pb.MeshEntry

func (s state) clone() state {
	c := make(state, len(s))
	for k, v := range s {
		c[k] = v
	}
	return c
}

// merge stores the entry if it is newer than the existing one for its key. It
// returns whether the entry was stored (used to decide whether to gossip the
// message further) and the resulting change (in bytes) to the total marshaled size
// of the state.
func (s state) merge(e *pb.MeshEntry, now time.Time) (bool, int) {
	if e.ExpiresAt.Before(now) {
		return false, 0
	}
	k := stateKey(string(e.Entry.GroupKey), e.Entry.Receiver)

	prev, ok := s[k]
	if !ok {
		s[k] = e
		return true, e.Size()
	}
	if prev.Entry.Timestamp.Before(e.Entry.Timestamp) {
		s[k] = e
		return true, e.Size() - prev.Size()
	}
	return false, 0
}

// sizeBytes returns the sum of the marshaled sizes of all entries. This matches
// the unit tracked incrementally by Log.size.
func (s state) sizeBytes() int {
	var n int
	for _, e := range s {
		n += e.Size()
	}
	return n
}

func (s state) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	for _, e := range s {
		if _, err := pbutil.WriteDelimited(&buf, e); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func decodeState(r io.Reader) (state, error) {
	st := state{}
	for {
		var e pb.MeshEntry
		_, err := pbutil.ReadDelimited(r, &e)
		if err == nil {
			if e.Entry == nil || e.Entry.Receiver == nil {
				return nil, ErrInvalidState
			}
			st[stateKey(string(e.Entry.GroupKey), e.Entry.Receiver)] = &e
			continue
		}
		if errors.Is(err, io.EOF) {
			break
		}
		return nil, err
	}
	return st, nil
}

func marshalMeshEntry(e *pb.MeshEntry) ([]byte, error) {
	var buf bytes.Buffer
	if _, err := pbutil.WriteDelimited(&buf, e); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Options configures a new Log implementation.
type Options struct {
	SnapshotReader io.Reader
	SnapshotFile   string

	Retention time.Duration
	Limits    Limits

	Logger  log.Logger
	Metrics prometheus.Registerer
}

func (o *Options) validate() error {
	if o.SnapshotFile != "" && o.SnapshotReader != nil {
		return errors.New("only one of SnapshotFile and SnapshotReader must be set")
	}

	return nil
}

// New creates a new notification log based on the provided options.
// The snapshot is loaded into the Log if it is set.
func New(o Options) (*Log, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	l := &Log{
		clock:              clock.New(),
		retention:          o.Retention,
		limits:             o.Limits,
		logger:             log.NewNopLogger(),
		st:                 state{},
		broadcast:          func([]byte) {},
		isReliableDelivery: clusterutil.OversizedMessage,
		metrics:            newMetrics(o.Metrics),
	}

	if o.Logger != nil {
		l.logger = o.Logger
	}

	if o.SnapshotFile != "" {
		if r, err := os.Open(o.SnapshotFile); err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
			level.Debug(l.logger).Log("msg", "notification log snapshot file doesn't exist", "err", err)
		} else {
			o.SnapshotReader = r
			defer r.Close()
		}
	}

	if o.SnapshotReader != nil {
		if err := l.loadSnapshot(o.SnapshotReader); err != nil {
			return l, err
		}
	}

	return l, nil
}

func (l *Log) now() time.Time {
	return l.clock.Now()
}

// Maintenance garbage collects the notification log state at the given interval. If the snapshot
// file is set, a snapshot is written to it afterwards.
// Terminates on receiving from stopc.
// If not nil, the last argument is an override for what to do as part of the maintenance - for advanced usage.
func (l *Log) Maintenance(interval time.Duration, snapf string, stopc <-chan struct{}, override MaintenanceFunc) {
	if interval == 0 || stopc == nil {
		level.Error(l.logger).Log("msg", "interval or stop signal are missing - not running maintenance")
		return
	}
	t := l.clock.Ticker(interval)
	defer t.Stop()

	var doMaintenance MaintenanceFunc
	doMaintenance = func() (int64, error) {
		var size int64
		if _, err := l.GC(); err != nil {
			return size, err
		}
		if snapf == "" {
			return size, nil
		}
		f, err := openReplace(snapf)
		if err != nil {
			return size, err
		}
		if size, err = l.Snapshot(f); err != nil {
			f.Close()
			return size, err
		}
		return size, f.Close()
	}

	if override != nil {
		doMaintenance = override
	}

	runMaintenance := func(do func() (int64, error)) error {
		l.metrics.maintenanceTotal.Inc()
		start := l.now().UTC()
		level.Debug(l.logger).Log("msg", "Running maintenance")
		size, err := do()
		l.metrics.snapshotSize.Set(float64(size))
		if err != nil {
			l.metrics.maintenanceErrorsTotal.Inc()
			return err
		}
		level.Debug(l.logger).Log("msg", "Maintenance done", "duration", l.now().Sub(start), "size", size)
		return nil
	}

Loop:
	for {
		select {
		case <-stopc:
			break Loop
		case <-t.C:
			if err := runMaintenance(doMaintenance); err != nil {
				level.Error(l.logger).Log("msg", "Running maintenance failed", "err", err)
			}
		}
	}

	// No need to run final maintenance if we don't want to snapshot.
	if snapf == "" {
		return
	}
	if err := runMaintenance(doMaintenance); err != nil {
		level.Error(l.logger).Log("msg", "Creating shutdown snapshot failed", "err", err)
	}
}

func receiverKey(r *pb.Receiver) string {
	return fmt.Sprintf("%s/%s/%d", r.GroupName, r.Integration, r.Idx)
}

// stateKey returns a string key for a log entry consisting of the group key
// and receiver.
func stateKey(k string, r *pb.Receiver) string {
	return fmt.Sprintf("%s:%s", k, receiverKey(r))
}

func (l *Log) Log(r *pb.Receiver, gkey string, firingAlerts, resolvedAlerts []uint64, expiry time.Duration) error {
	// Write all st with the same timestamp.
	now := l.now()
	key := stateKey(gkey, r)

	l.mtx.Lock()
	defer l.mtx.Unlock()

	_, isUpdate := l.st[key]
	if isUpdate {
		// Entry already exists, only overwrite if timestamp is newer.
		// This may happen with raciness or clock-drift across AM nodes.
		if l.st[key].Entry.Timestamp.After(now) {
			return nil
		}
	}

	expiresAt := now.Add(l.retention)
	if expiry > 0 && l.retention > expiry {
		expiresAt = now.Add(expiry)
	}

	e := &pb.MeshEntry{
		Entry: &pb.Entry{
			Receiver:       r,
			GroupKey:       []byte(gkey),
			Timestamp:      now,
			FiringAlerts:   firingAlerts,
			ResolvedAlerts: resolvedAlerts,
		},
		ExpiresAt: expiresAt,
	}

	b, err := marshalMeshEntry(e)
	if err != nil {
		return err
	}

	_, delta := l.st.merge(e, l.now())
	l.size += delta

	// Enforce the size limit by evicting existing entries (resolved first, then
	// oldest by last-notification time) until the log is within the limit. Evicting
	// rather than rejecting the new entry keeps the most recently active groups
	// recorded. The evicted entries are the least likely to be re-notified. The
	// notification has already been sent by this point, so at worst this re-notifies
	// an evicted group (a duplicate) later, never a missed notification.
	l.evictToLimit()

	l.metrics.size.Set(float64(l.size))
	l.broadcast(b)

	return nil
}

// evictToLimit removes entries until the total marshaled size is within
// Limits.MaxSizeBytes, preferring resolved entries and then the oldest (by last
// notification time). Caller must hold l.mtx.
func (l *Log) evictToLimit() {
	if l.limits.MaxSizeBytes == nil {
		return
	}
	max := l.limits.MaxSizeBytes()
	if max <= 0 || l.size <= max {
		return
	}

	// Order the evictable keys resolved-first, then oldest-first within each tier.
	// Evicting a resolved entry (no firing alerts) risks at most a duplicate
	// resolved notification, whereas evicting a firing entry risks a duplicate page,
	// so drop the harmless records before the harmful ones.
	keys := make([]string, 0, len(l.st))
	for k := range l.st {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		ei, ej := l.st[keys[i]].Entry, l.st[keys[j]].Entry
		ri, rj := len(ei.FiringAlerts) == 0, len(ej.FiringAlerts) == 0
		if ri != rj {
			return ri // resolved (no firing alerts) evicted first
		}
		return ei.Timestamp.Before(ej.Timestamp) // then oldest-first
	})

	for _, k := range keys {
		if l.size <= max {
			break
		}
		l.size -= l.st[k].Size()
		delete(l.st, k)
		l.metrics.entriesEvictedTotal.Inc()
	}
}

// GC implements the Log interface.
func (l *Log) GC() (int, error) {
	start := time.Now()
	defer func() { l.metrics.gcDuration.Observe(time.Since(start).Seconds()) }()

	now := l.now()
	var n int

	l.mtx.Lock()
	defer l.mtx.Unlock()

	for k, le := range l.st {
		if le.ExpiresAt.IsZero() {
			return n, errors.New("unexpected zero expiration timestamp")
		}
		if !le.ExpiresAt.After(now) {
			l.size -= le.Size()
			delete(l.st, k)
			n++
		}
	}
	l.metrics.size.Set(float64(l.size))

	return n, nil
}

// Query implements the Log interface.
func (l *Log) Query(params ...QueryParam) ([]*pb.Entry, error) {
	start := time.Now()
	l.metrics.queriesTotal.Inc()

	entries, err := func() ([]*pb.Entry, error) {
		q := &query{}
		for _, p := range params {
			if err := p(q); err != nil {
				return nil, err
			}
		}
		// TODO(fabxc): For now our only query mode is the most recent entry for a
		// receiver/group_key combination.
		if q.recv == nil || q.groupKey == "" {
			// TODO(fabxc): allow more complex queries in the future.
			// How to enable pagination?
			return nil, errors.New("no query parameters specified")
		}

		l.mtx.RLock()
		defer l.mtx.RUnlock()

		if le, ok := l.st[stateKey(q.groupKey, q.recv)]; ok {
			return []*pb.Entry{le.Entry}, nil
		}
		return nil, ErrNotFound
	}()
	if err != nil {
		l.metrics.queryErrorsTotal.Inc()
	}
	l.metrics.queryDuration.Observe(time.Since(start).Seconds())
	return entries, err
}

// loadSnapshot loads a snapshot generated by Snapshot() into the state.
func (l *Log) loadSnapshot(r io.Reader) error {
	st, err := decodeState(r)
	if err != nil {
		return err
	}

	l.mtx.Lock()
	l.st = st
	l.size = st.sizeBytes()
	l.mtx.Unlock()

	return nil
}

// Snapshot implements the Log interface.
func (l *Log) Snapshot(w io.Writer) (int64, error) {
	start := time.Now()
	defer func() { l.metrics.snapshotDuration.Observe(time.Since(start).Seconds()) }()

	l.mtx.RLock()
	defer l.mtx.RUnlock()

	b, err := l.st.MarshalBinary()
	if err != nil {
		return 0, err
	}

	return io.Copy(w, bytes.NewReader(b))
}

// MarshalBinary serializes all contents of the notification log.
func (l *Log) MarshalBinary() ([]byte, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	return l.st.MarshalBinary()
}

// Merge merges notification log state received from the cluster with the local state.
func (l *Log) Merge(b []byte) error {
	st, err := decodeState(bytes.NewReader(b))
	if err != nil {
		return err
	}
	l.mtx.Lock()
	defer l.mtx.Unlock()
	now := l.now()

	for _, e := range st {
		merged, delta := l.st.merge(e, now)
		l.size += delta
		if merged && !l.isReliableDelivery(b) {
			// If this is the first we've seen the message and it was
			// not sent reliably to all nodes, gossip it to other nodes.
			// We don't propagate reliable messages because they're
			// sent to all nodes already.
			l.broadcast(b)
			l.metrics.propagatedMessagesTotal.Inc()
			level.Debug(l.logger).Log("msg", "gossiping new entry", "entry", e)
		}
	}
	l.metrics.size.Set(float64(l.size))
	return nil
}

// SetBroadcast sets a broadcast callback that will be invoked with serialized state
// on updates.
func (l *Log) SetBroadcast(f func([]byte)) {
	l.mtx.Lock()
	l.broadcast = f
	l.mtx.Unlock()
}

// SetIsReliableDelivery sets a callback that returns true if the given message
// was delivered reliably to all peers and should not be re-broadcast.
func (l *Log) SetIsReliableDelivery(f func([]byte) bool) {
	l.mtx.Lock()
	l.isReliableDelivery = f
	l.mtx.Unlock()
}

// replaceFile wraps a file that is moved to another filename on closing.
type replaceFile struct {
	*os.File
	filename string
}

func (f *replaceFile) Close() error {
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.File.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), f.filename)
}

// openReplace opens a new temporary file that is moved to filename on closing.
func openReplace(filename string) (*replaceFile, error) {
	tmpFilename := fmt.Sprintf("%s.%x", filename, uint64(rand.Int63()))

	f, err := os.Create(tmpFilename)
	if err != nil {
		return nil, err
	}

	rf := &replaceFile{
		File:     f,
		filename: filename,
	}
	return rf, nil
}
