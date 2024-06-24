package ingester_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/duneanalytics/blockchain-ingester/ingester"
	duneapi_mock "github.com/duneanalytics/blockchain-ingester/mocks/duneapi"
	jsonrpc_mock "github.com/duneanalytics/blockchain-ingester/mocks/jsonrpc"
	"github.com/duneanalytics/blockchain-ingester/models"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestRunLoopUntilCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	maxBlockNumber := int64(10)
	sentBlockNumber := int64(0)
	producedBlockNumber := int64(0)
	duneapi := &duneapi_mock.BlockchainIngesterMock{
		SendBlockFunc: func(_ context.Context, block models.RPCBlock) error {
			atomic.StoreInt64(&sentBlockNumber, block.BlockNumber)
			if block.BlockNumber == maxBlockNumber {
				// cancel execution when we send the last block
				cancel()
				return context.Canceled
			}
			return nil
		},
		SendBlocksFunc: func(_ context.Context, blocks []models.RPCBlock) error {
			if len(blocks) == 0 {
				return nil
			}

			next := sentBlockNumber + 1
			for _, block := range blocks {
				// We cannot send blocks out of order to DuneAPI
				if block.BlockNumber != next {
					return fmt.Errorf("expected block %d, got %d", next, block.BlockNumber)
				}
				next++
			}

			lastBlockNumber := blocks[len(blocks)-1].BlockNumber
			atomic.StoreInt64(&sentBlockNumber, lastBlockNumber)
			if lastBlockNumber >= maxBlockNumber {
				// cancel execution when we have sent the last block
				cancel()
				return context.Canceled
			}

			return nil
		},
		PostProgressReportFunc: func(_ context.Context, _ models.BlockchainIndexProgress) error {
			return nil
		},
	}
	rpcClient := &jsonrpc_mock.BlockchainClientMock{
		LatestBlockNumberFunc: func() (int64, error) {
			return maxBlockNumber + 1, nil
		},
		BlockByNumberFunc: func(_ context.Context, blockNumber int64) (models.RPCBlock, error) {
			atomic.StoreInt64(&producedBlockNumber, blockNumber)
			return models.RPCBlock{
				BlockNumber: blockNumber,
				Payload:     []byte(`block`),
			}, nil
		},
		CloseFunc: func() error {
			return nil
		},
	}
	// Swap these to see logs
	// logOutput := os.Stderr
	logOutput := io.Discard
	ing := ingester.New(slog.New(slog.NewTextHandler(logOutput, nil)), rpcClient, duneapi, ingester.Config{
		MaxBatchSize:         1,
		PollInterval:         1000 * time.Millisecond,
		BatchRequestInterval: time.Nanosecond,
	})

	err := ing.Run(ctx, 1, -1)                // run until canceled
	require.ErrorIs(t, err, context.Canceled) // this is expected
	require.GreaterOrEqual(t, sentBlockNumber, maxBlockNumber)
}

func TestProduceBlockNumbers(t *testing.T) {
	duneapi := &duneapi_mock.BlockchainIngesterMock{
		SendBlockFunc: func(_ context.Context, _ models.RPCBlock) error {
			return nil
		},
		PostProgressReportFunc: func(_ context.Context, _ models.BlockchainIndexProgress) error {
			return nil
		},
	}
	rpcClient := &jsonrpc_mock.BlockchainClientMock{
		LatestBlockNumberFunc: func() (int64, error) {
			return 100_000, nil
		},
		BlockByNumberFunc: func(_ context.Context, blockNumber int64) (models.RPCBlock, error) {
			return models.RPCBlock{BlockNumber: blockNumber}, nil
		},
		CloseFunc: func() error {
			return nil
		},
	}
	logOutput := io.Discard
	ing := ingester.New(slog.New(slog.NewTextHandler(logOutput, nil)), rpcClient, duneapi, ingester.Config{
		MaxBatchSize:         1,
		PollInterval:         1000 * time.Millisecond,
		BatchRequestInterval: time.Nanosecond,
	})
	blockNumbers := make(chan int64)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ing.ProduceBlockNumbers(context.Background(), blockNumbers, 1, 100_000)
	}()
	for i := 1; i <= 100_000; i++ {
		require.Equal(t, int64(i), <-blockNumbers)
	}
	wg.Wait()
}

func TestSendBlocks(t *testing.T) {
	sentBlockNumber := int64(0)
	duneapi := &duneapi_mock.BlockchainIngesterMock{
		SendBlockFunc: func(_ context.Context, block models.RPCBlock) error {
			// We cannot send blocks out of order to DuneAPI
			if block.BlockNumber != sentBlockNumber+1 {
				return fmt.Errorf("expected block %d, got %d", sentBlockNumber+1, block.BlockNumber)
			}

			atomic.StoreInt64(&sentBlockNumber, block.BlockNumber)
			return nil
		},
		SendBlocksFunc: func(_ context.Context, blocks []models.RPCBlock) error {
			if len(blocks) == 0 {
				return nil
			}

			next := sentBlockNumber + 1
			for _, block := range blocks {
				// We cannot send blocks out of order to DuneAPI
				if block.BlockNumber != next {
					return fmt.Errorf("expected block %d, got %d", next, block.BlockNumber)
				}
				next++
			}

			lastBlockNumber := blocks[len(blocks)-1].BlockNumber
			atomic.StoreInt64(&sentBlockNumber, lastBlockNumber)
			return nil
		},
		PostProgressReportFunc: func(_ context.Context, _ models.BlockchainIndexProgress) error {
			return nil
		},
	}
	// Swap these to see logs
	// logOutput := os.Stderr
	logOutput := io.Discard
	ing := ingester.New(
		slog.New(slog.NewTextHandler(logOutput, nil)),
		nil, // node client isn't used in this unit test
		duneapi,
		ingester.Config{
			MaxBatchSize:         10,                      // this won't matter as we only run SendBlocks
			PollInterval:         1000 * time.Millisecond, // ditto
			BatchRequestInterval: time.Nanosecond,         // this matters
		},
	)

	blocks := make(chan models.RPCBlock)

	startFromBlock := 1

	group, _ := errgroup.WithContext(context.Background())
	group.Go(func() error {
		return ing.SendBlocks(context.Background(), blocks, int64(startFromBlock))
	})

	// Send blocks except the next block, ensure none are sent to the API
	for _, n := range []int64{2, 3, 4, 5, 10} {
		blocks <- models.RPCBlock{BlockNumber: n}
		require.Equal(t, int64(0), sentBlockNumber)

		time.Sleep(time.Nanosecond) // need to sleep to allow a few ticks of BatchRequestInterval
	}
	// Now send the first block
	blocks <- models.RPCBlock{BlockNumber: 1}
	close(blocks)
	require.NoError(t, group.Wait())

	// Ensure the last correct block was sent
	require.Equal(t, int64(5), sentBlockNumber)
}

// TestRunLoopBlocksOutOfOrder asserts that we can fetch blocks concurrently and that we ingest them in order
// even if they are produced out of order. We ensure they are produced out of order by sleeping a random amount of time.
func TestRunLoopBlocksOutOfOrder(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	maxBlockNumber := int64(1000)
	sentBlockNumber := int64(0)
	producedBlockNumber := int64(0)
	duneapi := &duneapi_mock.BlockchainIngesterMock{
		SendBlockFunc: func(_ context.Context, block models.RPCBlock) error {
			// We cannot send blocks out of order to DuneAPI
			if block.BlockNumber != sentBlockNumber+1 {
				return fmt.Errorf("expected block %d, got %d", sentBlockNumber+1, block.BlockNumber)
			}

			atomic.StoreInt64(&sentBlockNumber, block.BlockNumber)
			return nil
		},
		SendBlocksFunc: func(_ context.Context, blocks []models.RPCBlock) error {
			if len(blocks) == 0 {
				return nil
			}

			next := sentBlockNumber + 1
			for _, block := range blocks {
				// We cannot send blocks out of order to DuneAPI
				if block.BlockNumber != next {
					return fmt.Errorf("expected block %d, got %d", next, block.BlockNumber)
				}
				next++
			}

			lastBlockNumber := blocks[len(blocks)-1].BlockNumber
			atomic.StoreInt64(&sentBlockNumber, lastBlockNumber)
			if lastBlockNumber >= maxBlockNumber {
				// cancel execution when we have sent the last block
				cancel()
				return context.Canceled
			}

			return nil
		},
		PostProgressReportFunc: func(_ context.Context, _ models.BlockchainIndexProgress) error {
			return nil
		},
	}
	rpcClient := &jsonrpc_mock.BlockchainClientMock{
		LatestBlockNumberFunc: func() (int64, error) {
			return maxBlockNumber + 1, nil
		},
		BlockByNumberFunc: func(_ context.Context, blockNumber int64) (models.RPCBlock, error) {
			// Get blocks out of order by sleeping for a random amount of ms
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			atomic.StoreInt64(&producedBlockNumber, blockNumber)
			return models.RPCBlock{BlockNumber: blockNumber}, nil
		},
		CloseFunc: func() error {
			return nil
		},
	}
	// Swap these to see logs
	// logOutput := os.Stderr
	logOutput := io.Discard
	ing := ingester.New(
		slog.New(slog.NewTextHandler(logOutput, nil)),
		rpcClient,
		duneapi,
		ingester.Config{
			MaxBatchSize:         10, // fetch blocks in multiple goroutines
			PollInterval:         1000 * time.Millisecond,
			BatchRequestInterval: time.Nanosecond,
		},
	)

	err := ing.Run(ctx, 1, -1)                // run until canceled
	require.ErrorIs(t, err, context.Canceled) // this is expected
	require.GreaterOrEqual(t, sentBlockNumber, maxBlockNumber)
}
