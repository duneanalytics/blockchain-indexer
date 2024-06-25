// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package duneapi_mock

import (
	"context"
	"github.com/duneanalytics/blockchain-ingester/client/duneapi"
	"github.com/duneanalytics/blockchain-ingester/models"
	"sync"
)

// Ensure, that BlockchainIngesterMock does implement duneapi.BlockchainIngester.
// If this is not the case, regenerate this file with moq.
var _ duneapi.BlockchainIngester = &BlockchainIngesterMock{}

// BlockchainIngesterMock is a mock implementation of duneapi.BlockchainIngester.
//
//	func TestSomethingThatUsesBlockchainIngester(t *testing.T) {
//
//		// make and configure a mocked duneapi.BlockchainIngester
//		mockedBlockchainIngester := &BlockchainIngesterMock{
//			GetProgressReportFunc: func(ctx context.Context) (*models.BlockchainIndexProgress, error) {
//				panic("mock out the GetProgressReport method")
//			},
//			PostProgressReportFunc: func(ctx context.Context, progress models.BlockchainIndexProgress) error {
//				panic("mock out the PostProgressReport method")
//			},
//			SendBlocksFunc: func(ctx context.Context, payloads []models.RPCBlock) error {
//				panic("mock out the SendBlocks method")
//			},
//		}
//
//		// use mockedBlockchainIngester in code that requires duneapi.BlockchainIngester
//		// and then make assertions.
//
//	}
type BlockchainIngesterMock struct {
	// GetProgressReportFunc mocks the GetProgressReport method.
	GetProgressReportFunc func(ctx context.Context) (*models.BlockchainIndexProgress, error)

	// PostProgressReportFunc mocks the PostProgressReport method.
	PostProgressReportFunc func(ctx context.Context, progress models.BlockchainIndexProgress) error

	// SendBlocksFunc mocks the SendBlocks method.
	SendBlocksFunc func(ctx context.Context, payloads []models.RPCBlock) error

	// calls tracks calls to the methods.
	calls struct {
		// GetProgressReport holds details about calls to the GetProgressReport method.
		GetProgressReport []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// PostProgressReport holds details about calls to the PostProgressReport method.
		PostProgressReport []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Progress is the progress argument value.
			Progress models.BlockchainIndexProgress
		}
		// SendBlocks holds details about calls to the SendBlocks method.
		SendBlocks []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Payloads is the payloads argument value.
			Payloads []models.RPCBlock
		}
	}
	lockGetProgressReport  sync.RWMutex
	lockPostProgressReport sync.RWMutex
	lockSendBlocks         sync.RWMutex
}

// GetProgressReport calls GetProgressReportFunc.
func (mock *BlockchainIngesterMock) GetProgressReport(ctx context.Context) (*models.BlockchainIndexProgress, error) {
	if mock.GetProgressReportFunc == nil {
		panic("BlockchainIngesterMock.GetProgressReportFunc: method is nil but BlockchainIngester.GetProgressReport was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockGetProgressReport.Lock()
	mock.calls.GetProgressReport = append(mock.calls.GetProgressReport, callInfo)
	mock.lockGetProgressReport.Unlock()
	return mock.GetProgressReportFunc(ctx)
}

// GetProgressReportCalls gets all the calls that were made to GetProgressReport.
// Check the length with:
//
//	len(mockedBlockchainIngester.GetProgressReportCalls())
func (mock *BlockchainIngesterMock) GetProgressReportCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockGetProgressReport.RLock()
	calls = mock.calls.GetProgressReport
	mock.lockGetProgressReport.RUnlock()
	return calls
}

// PostProgressReport calls PostProgressReportFunc.
func (mock *BlockchainIngesterMock) PostProgressReport(ctx context.Context, progress models.BlockchainIndexProgress) error {
	if mock.PostProgressReportFunc == nil {
		panic("BlockchainIngesterMock.PostProgressReportFunc: method is nil but BlockchainIngester.PostProgressReport was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		Progress models.BlockchainIndexProgress
	}{
		Ctx:      ctx,
		Progress: progress,
	}
	mock.lockPostProgressReport.Lock()
	mock.calls.PostProgressReport = append(mock.calls.PostProgressReport, callInfo)
	mock.lockPostProgressReport.Unlock()
	return mock.PostProgressReportFunc(ctx, progress)
}

// PostProgressReportCalls gets all the calls that were made to PostProgressReport.
// Check the length with:
//
//	len(mockedBlockchainIngester.PostProgressReportCalls())
func (mock *BlockchainIngesterMock) PostProgressReportCalls() []struct {
	Ctx      context.Context
	Progress models.BlockchainIndexProgress
} {
	var calls []struct {
		Ctx      context.Context
		Progress models.BlockchainIndexProgress
	}
	mock.lockPostProgressReport.RLock()
	calls = mock.calls.PostProgressReport
	mock.lockPostProgressReport.RUnlock()
	return calls
}

// SendBlocks calls SendBlocksFunc.
func (mock *BlockchainIngesterMock) SendBlocks(ctx context.Context, payloads []models.RPCBlock) error {
	if mock.SendBlocksFunc == nil {
		panic("BlockchainIngesterMock.SendBlocksFunc: method is nil but BlockchainIngester.SendBlocks was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		Payloads []models.RPCBlock
	}{
		Ctx:      ctx,
		Payloads: payloads,
	}
	mock.lockSendBlocks.Lock()
	mock.calls.SendBlocks = append(mock.calls.SendBlocks, callInfo)
	mock.lockSendBlocks.Unlock()
	return mock.SendBlocksFunc(ctx, payloads)
}

// SendBlocksCalls gets all the calls that were made to SendBlocks.
// Check the length with:
//
//	len(mockedBlockchainIngester.SendBlocksCalls())
func (mock *BlockchainIngesterMock) SendBlocksCalls() []struct {
	Ctx      context.Context
	Payloads []models.RPCBlock
} {
	var calls []struct {
		Ctx      context.Context
		Payloads []models.RPCBlock
	}
	mock.lockSendBlocks.RLock()
	calls = mock.calls.SendBlocks
	mock.lockSendBlocks.RUnlock()
	return calls
}
