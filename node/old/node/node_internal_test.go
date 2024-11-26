package node

import (
	"bufio"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/stretchr/testify/require"

	"github.com/blocklessnetwork/b7s/host"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/testing/mocks"
)

const (
	loopback = "127.0.0.1"

	// How long can the client wait for a published message before giving up.
	publishTimeout = 10 * time.Second

	// It seems like a delay is needed so that the hosts exchange information about the fact
	// that they are subscribed to the same topic. If that does not happen, node might publish
	// a message too soon and the client might miss it. It will then wait for a published message in vain.
	// This is the pause we make after subscribing to the topic and before publishing a message.
	// In reality as little as 250ms is enough, but lets allow a longer time for when
	// tests are executed in parallel or on weaker machines.
	subscriptionDiseminationPause = 2 * time.Second
)

func TestNode_New(t *testing.T) {

	var (
		logger          = mocks.NoopLogger
		store           = mocks.BaselineStore(t)
		functionHandler = mocks.BaselineFStore(t)
		executor        = mocks.BaselineExecutor(t)
	)

	host, err := host.New(logger, loopback, 0)
	require.NoError(t, err)

	t.Run("create a head node", func(t *testing.T) {
		t.Parallel()

		node, err := New(logger, host, store, functionHandler, WithRole(blockless.HeadNode))
		require.NoError(t, err)
		require.NotNil(t, node)

		// Creating a head node with executor fails.
		_, err = New(logger, host, store, functionHandler, WithRole(blockless.HeadNode), WithExecutor(executor))
		require.Error(t, err)
	})
	t.Run("create a worker node", func(t *testing.T) {
		t.Parallel()

		node, err := New(logger, host, store, functionHandler, WithRole(blockless.WorkerNode), WithExecutor(executor), WithWorkspace(t.TempDir()))
		require.NoError(t, err)
		require.NotNil(t, node)

		// Creating a worker node without executor fails.
		_, err = New(logger, host, store, functionHandler, WithRole(blockless.WorkerNode))
		require.Error(t, err)
	})
}

func createNode(t *testing.T, role blockless.NodeRole) *Node {
	t.Helper()

	var (
		logger          = mocks.NoopLogger
		store           = mocks.BaselineStore(t)
		functionHandler = mocks.BaselineFStore(t)
	)

	host, err := host.New(logger, loopback, 0)
	require.NoError(t, err)

	opts := []Option{
		WithRole(role),
	}

	if role == blockless.WorkerNode {
		executor := mocks.BaselineExecutor(t)
		opts = append(opts, WithExecutor(executor))
		opts = append(opts, WithWorkspace(t.TempDir()))
	}

	node, err := New(logger, host, store, functionHandler, opts...)
	require.NoError(t, err)

	return node
}

func getStreamPayload(t *testing.T, stream network.Stream, output any) {
	t.Helper()

	buf := bufio.NewReader(stream)
	payload, err := buf.ReadBytes('\n')
	require.ErrorIs(t, err, io.EOF)

	err = json.Unmarshal(payload, output)
	require.NoError(t, err)
}

func serialize(t *testing.T, message any) []byte {
	payload, err := json.Marshal(message)
	require.NoError(t, err)

	return payload
}
