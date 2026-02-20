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

package cluster

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/hashicorp/memberlist"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
)

func TestNormalMessagesGossiped(t *testing.T) {
	var sent bool
	c := newChannel(
		func(_ []byte) { sent = true },
		func() []*memberlist.Node { return nil },
		func(_ *memberlist.Node, _ []byte) error { return nil },
	)

	c.Broadcast([]byte{})

	if sent != true {
		t.Fatalf("small message not sent")
	}
}

func TestOversizedMessagesGossiped(t *testing.T) {
	var sent bool
	ctx, cancel := context.WithCancel(context.Background())
	c := newChannel(
		func(_ []byte) {},
		func() []*memberlist.Node { return []*memberlist.Node{{}} },
		func(_ *memberlist.Node, _ []byte) error { sent = true; cancel(); return nil },
	)

	f, err := os.Open("/dev/zero")
	if err != nil {
		t.Fatalf("failed to open /dev/zero: %v", err)
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	toCopy := int64(800)
	if n, err := io.CopyN(buf, f, toCopy); err != nil {
		t.Fatalf("failed to copy bytes: %v", err)
	} else if n != toCopy {
		t.Fatalf("wanted to copy %d bytes, only copied %d", toCopy, n)
	}

	c.Broadcast(buf.Bytes())

	<-ctx.Done()

	if sent != true {
		t.Fatalf("oversized message not sent")
	}
}

func TestWithQueueSizePanicsOnNonPositive(t *testing.T) {
	require.Panics(t, func() { WithQueueSize(0) })
	require.Panics(t, func() { WithQueueSize(-1) })
}

func TestDefaultQueueSize(t *testing.T) {
	c := newChannel(
		func(_ []byte) {},
		func() []*memberlist.Node { return nil },
		func(_ *memberlist.Node, _ []byte) error { return nil },
	)

	require.Equal(t, defaultQueueSize, cap(c.msgc))
}

func TestWithQueueSize(t *testing.T) {
	reg := prometheus.NewRegistry()
	stopc := make(chan struct{})
	defer close(stopc)

	sendEntered := make(chan struct{}, 1)
	blockSend := make(chan struct{})

	c := NewChannel(
		"test",
		func(_ []byte) {},
		func() []*memberlist.Node { return []*memberlist.Node{{}} },
		func(_ *memberlist.Node, _ []byte) error {
			sendEntered <- struct{}{}
			<-blockSend
			return nil
		},
		log.NewNopLogger(),
		stopc,
		reg,
		WithQueueSize(2),
	)

	require.Equal(t, 2, cap(c.msgc))

	oversizedMsg := make([]byte, MaxGossipPacketSize)

	// First message is consumed by handleOverSizedMessages, which blocks in sendOversize.
	c.Broadcast(oversizedMsg)
	<-sendEntered

	// Fill the remaining queue capacity.
	c.Broadcast(oversizedMsg)
	c.Broadcast(oversizedMsg)

	// The queue is now full; next message should be dropped.
	c.Broadcast(oversizedMsg)

	var m dto.Metric
	require.NoError(t, c.oversizeGossipMessageDroppedTotal.Write(&m))
	require.Equal(t, 1.0, m.GetCounter().GetValue())

	close(blockSend)
}

func newChannel(
	send func([]byte),
	peers func() []*memberlist.Node,
	sendOversize func(*memberlist.Node, []byte) error,
	opts ...ChannelOption,
) *Channel {
	return NewChannel(
		"test",
		send,
		peers,
		sendOversize,
		log.NewNopLogger(),
		make(chan struct{}),
		prometheus.NewRegistry(),
		opts...,
	)
}
