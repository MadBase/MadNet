package blockchain

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/MadBase/MadNet/blockchain/interfaces"
	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/logging"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

//
var (
	ErrUnknownRequest = errors.New("unknown request type")
)

type ErrRecoverable struct {
	message string
}

func (e *ErrRecoverable) Error() string {
	return e.message
}

type ErrTransactionStale struct {
	message string
}

func (e *ErrTransactionStale) Error() string {
	return e.message
}

type ErrNonRecoverable struct {
	message string
}

func (e *ErrNonRecoverable) Error() string {
	return e.message
}

type ErrInvalidMonitorRequest struct {
	message string
}

func (e *ErrInvalidMonitorRequest) Error() string {
	return e.message
}

type ErrTransactionAlreadyQueued struct {
	message string
}

func (e *ErrTransactionAlreadyQueued) Error() string {
	return e.message
}

type ErrInvalidTransactionRequest struct {
	message string
}

func (e *ErrInvalidTransactionRequest) Error() string {
	return e.message
}

// Internal struct to keep track of transactions that are being monitoring
type TransactionInfo struct {
	ctx                     context.Context                   // ctx used for calling the monitoring a certain tx
	txn                     *types.Transaction                // Transaction object
	selector                interfaces.FuncSelector           // 4 bytes that identify the function being called by the tx
	functionSignature       string                            // function signature as we see on the smart contracts
	startedMonitoringHeight uint64                            // ethereum height where we first added the tx to be watched. Mainly used to see if a tx was dropped
	receiptResponseChannel  *ResponseChannel[ReceiptResponse] // channel where the response will be sent
}

// Struct to keep track of the receipts
type ReceiptInfo struct {
	receipt           *types.Receipt // receipt object
	retrievedAtHeight uint64         // block where receipt was added to the cache
}

// Internal struct to keep track of what blocks we already checked during monitoring
type BlockInfo struct {
	height uint64      // block height
	hash   common.Hash // block header hash
}

// Compare if 2 blockInfo structs are equal by comparing the height and block
// hash. Return true in case they are equal, false otherwise.
func (a *BlockInfo) Equal(b *BlockInfo) bool {
	return bytes.Equal(a.hash[:], b.hash[:]) && a.height == b.height
}

// Type to do request against the tx receipt monitoring system. Ctx and response channel should be set
type MonitorRequest struct {
	ctx             context.Context                   // tx ctx used for tx monitoring cancellation
	txn             *types.Transaction                // the transaction that should watched
	responseChannel *ResponseChannel[MonitorResponse] // channel where we going to send the request response
}

// Constrain interface used by the Response channel generics
type transferable interface {
	MonitorResponse | ReceiptResponse
}

// Type that it's going to be used to reply a request
type MonitorResponse struct {
	txnHash                common.Hash                       // Hash of the txs which this response belongs
	err                    error                             // errors that happened when processing the monitor request
	receiptResponseChannel *ResponseChannel[ReceiptResponse] // non blocking channel where the result from the tx/receipt monitoring will be send
}

// Response of the monitoring system
type ReceiptResponse struct {
	txnHash common.Hash    // Hash of the txs which this response belongs
	err     error          // response error that happened during processing
	receipt *types.Receipt // tx receipt after txConfirmationBlocks of a tx that was not queued in txGroup
}

// A response channel is basically a non-blocking channel that can only be
// written and closed once. The internal channel is closed after the first
// message is sent. Additional tries to send a message result in no-op. The
// writes to the internal channel are non-blocking calls. If for some reason the
// internal channel is full, the message is dropped and log is recorded. Only
// first attempt to close the Response channel will result in the closing.
// Additional calls are no-op.
type ResponseChannel[T transferable] struct {
	writeOnce sync.Once
	channel   chan *T       // internal channel
	isClosed  bool          // flag to check if a channel is closed or not
	logger    *logrus.Entry // logger using for logging error when trying to write a response more than once
}

// Create a new response channel.
func NewResponseChannel[T transferable](logger *logrus.Entry) *ResponseChannel[T] {
	return &ResponseChannel[T]{channel: make(chan *T, 1), logger: logger}
}

// send a unique response and close the internal channel. Additional calls to
// this function will be no-op
func (rc *ResponseChannel[T]) SendResponse(response *T) {
	if !rc.isClosed {
		select {
		case rc.channel <- response:
		default:
			rc.logger.Debugf("Failed to write request to channel")
		}
		rc.CloseChannel()
	}
}

// Close the internal channel. Additional calls will be no-op
func (rc *ResponseChannel[T]) CloseChannel() {
	rc.writeOnce.Do(func() {
		rc.isClosed = true
		close(rc.channel)
	})
}

// TransactionProfile to keep track of gas metrics in the overall system
type TransactionProfile struct {
	AverageGas   uint64
	MinimumGas   uint64
	MaximumGas   uint64
	TotalCount   uint64
	TotalGas     uint64
	TotalSuccess uint64
}

// Internal struct used to send work requests to the workers that will retrieve
// the receipts
type MonitorWorkRequest struct {
	txn    TransactionInfo // TransactionInfo object that contains the data that will be used to retrieve the receipt from the blockchain
	height uint64          // Current height of the blockchain head
}

// Internal struct used by the workers to communicate the result from the receipt retrieval work
type MonitorWorkResponse struct {
	txnHash common.Hash    // hash of transaction
	err     error          // any error found during the receipt retrieve (can be NonRecoverable, Recoverable or TransactionState errors)
	receipt *types.Receipt // receipt retrieved (can be nil) if a receipt was not found or it's not ready yet
}

// Backend struct used to monitor Ethereum transactions and retrieve their receipts
type WatcherBackend struct {
	mainCtx              context.Context                                // main context for the background services
	lastProcessedBlock   *BlockInfo                                     // Last ethereum block that we checked for receipts
	monitoredTxns        map[common.Hash]TransactionInfo                // Map of transactions whose receipts we're looking for
	receiptCache         map[common.Hash]ReceiptInfo                    // Receipts retrieved from transactions
	txConfirmationBlocks uint64                                         // number of ethereum blocks that we should wait to consider a receipt valid
	aggregates           map[interfaces.FuncSelector]TransactionProfile // Struct to keep track of the gas metrics used by the system
	client               interfaces.GethClient                          // An interface with the Geth functionality we need
	knownSelectors       interfaces.SelectorMap                         // Map with signature -> name
	logger               *logrus.Entry                                  // Logger to log messages
	requestChannel       <-chan MonitorRequest                          // Channel used to send request to this backend service
}

func (b *WatcherBackend) Loop() {

	poolingTime := time.After(constants.TxPollingTime)
	for {
		select {
		case req, ok := <-b.requestChannel:
			if !ok {
				b.logger.Debugf("request channel closed, exiting")
				return
			}
			b.logger.Debugf("received request: %v channel open: %v", req, ok)
			if req.responseChannel == nil {
				b.logger.Debug("Invalid request for txn without a response channel, ignoring")
				continue
			}
			b.queue(req)

		case <-poolingTime:
			b.collectReceipts()
			poolingTime = time.After(constants.TxPollingTime)
		}
	}
	b.logger.Debug("finished")
}

func (b *WatcherBackend) queue(req MonitorRequest) {

	if req.txn == nil {
		req.responseChannel.SendResponse(&MonitorResponse{err: &ErrInvalidMonitorRequest{"invalid request, missing txn object"}})
		return
	}
	if req.ctx == nil {
		req.responseChannel.SendResponse(&MonitorResponse{err: &ErrInvalidMonitorRequest{"invalid request, missing ctx"}})
		return
	}

	txnHash := req.txn.Hash()
	receiptResponseChannel := NewResponseChannel[ReceiptResponse](b.logger)

	// if we already have the records of the receipt for this tx we try to send the
	// receipt back
	if receipt, ok := b.receiptCache[txnHash]; ok {
		receiptResponseChannel.SendResponse(&ReceiptResponse{receipt: receipt.receipt, txnHash: txnHash})
	} else {
		if _, ok := b.monitoredTxns[txnHash]; ok {
			req.responseChannel.SendResponse(&MonitorResponse{err: &ErrTransactionAlreadyQueued{"invalid request, tx already queued, try to get receipt later!"}})
			return
		}

		selector := ExtractSelector(req.txn.Data())
		//todo: replace this with a generated mapping from the bindings
		sig := b.knownSelectors.Signature(selector)

		logEntry := b.logger.WithField("Transaction", txnHash).
			WithField("Function", sig).
			WithField("Selector", fmt.Sprintf("%x", selector))

		b.monitoredTxns[txnHash] = TransactionInfo{
			ctx:                    req.ctx,
			txn:                    req.txn,
			selector:               selector,
			functionSignature:      sig,
			receiptResponseChannel: receiptResponseChannel,
		}
		logEntry.Debug("Transaction queued")
	}
	req.responseChannel.SendResponse(&MonitorResponse{receiptResponseChannel: receiptResponseChannel})
}

func (b *WatcherBackend) collectReceipts() {

	lenMonitoredTxns := len(b.monitoredTxns)

	// If there's no tx to be monitored just return
	if lenMonitoredTxns == 0 {
		b.logger.Tracef("TxMonitor: no transaction to watch")
		return
	}

	networkCtx, cf := context.WithTimeout(b.mainCtx, constants.TxNetworkTimeout)
	defer cf()

	blockHeader, err := b.client.HeaderByNumber(networkCtx, nil)
	if err != nil {
		b.logger.Debugf("TxMonitor: error getting latest block number from ethereum node: %v", err)
		return
	}
	blockInfo := &BlockInfo{
		blockHeader.Number.Uint64(),
		blockHeader.Hash(),
	}

	if b.lastProcessedBlock.Equal(blockInfo) {
		b.logger.Tracef("TxMonitor: block: %v with hash: %v already processed", blockInfo.height, blockInfo.hash.Hex())
		return
	}

	var expiredTxs []common.Hash
	var finishedTxs map[common.Hash]MonitorWorkResponse

	numWorkers := min(max(uint64(lenMonitoredTxns)/4, 128), 1)
	requestWorkChannel := make(chan MonitorWorkRequest, lenMonitoredTxns+3)
	responseWorkChannel := make(chan MonitorWorkResponse, lenMonitoredTxns+3)

	for txn, txnInfo := range b.monitoredTxns {
		select {
		case <-txnInfo.ctx.Done():
			// the go-routine who wanted this information has stopped caring. This most
			// likely indicates a failure, and cancellation of polling prevents a memory
			// leak
			b.logger.Debugf("context for tx %v is finished, marking it to be excluded!", txn.Hex())
			expiredTxs = append(expiredTxs, txnInfo.txn.Hash())
		default:
			// if this is the first time seeing a tx
			if txnInfo.startedMonitoringHeight == 0 {
				txnInfo.startedMonitoringHeight = blockInfo.height
			}
			requestWorkChannel <- MonitorWorkRequest{txnInfo, blockInfo.height}
		}
	}

	// close the request channel, so the workers know when to finish
	close(requestWorkChannel)

	workerPool := NewWorkerPool(b.mainCtx, b.client, b.logger, b.txConfirmationBlocks, requestWorkChannel, responseWorkChannel)

	// spawn the workers and wait for all to complete
	go workerPool.ExecuteWork(numWorkers)

	for workResponse := range responseWorkChannel {
		select {
		case <-b.mainCtx.Done():
			// main thread was killed
			return
		default:
			if workResponse.err != nil {
				if _, ok := err.(*ErrRecoverable); !ok {
					b.logger.Debugf("Couldn't get tx receipt for tx:%v cause: %v", workResponse.txnHash, workResponse.err)
					finishedTxs[workResponse.txnHash] = workResponse
				}
			} else if workResponse.receipt != nil {
				b.logger.Debugf("Successfully got receipt for tx:%v", workResponse.txnHash)
				b.receiptCache[workResponse.txnHash] = ReceiptInfo{receipt: workResponse.receipt, retrievedAtHeight: blockInfo.height}
				finishedTxs[workResponse.txnHash] = workResponse
			}
		}
	}

	// Cleaning finished and failed transactions
	for txnHash, workResponse := range finishedTxs {
		if txnInfo, ok := b.monitoredTxns[txnHash]; ok {
			txnInfo.receiptResponseChannel.SendResponse(&ReceiptResponse{txnHash: workResponse.txnHash, receipt: workResponse.receipt, err: workResponse.err})
			delete(b.monitoredTxns, txnHash)
		}
	}

	// Cleaning expired transactions
	for _, txnHash := range expiredTxs {
		if txnInfo, ok := b.monitoredTxns[txnHash]; ok {
			txnInfo.receiptResponseChannel.CloseChannel()
			delete(b.monitoredTxns, txnHash)
		}
	}

	var expiredReceipts []common.Hash
	// Marking expired receipts
	for receiptTxnHash, receiptInfo := range b.receiptCache {
		if blockInfo.height >= receiptInfo.retrievedAtHeight+constants.TxReceiptCacheMaxBlocks {
			expiredReceipts = append(expiredReceipts, receiptTxnHash)
		}
	}

	// being paranoic and excluding the expired receipts in another loop
	for _, receiptTxHash := range expiredTxs {
		delete(b.receiptCache, receiptTxHash)
	}

	b.lastProcessedBlock = blockInfo
}

// Structs that keep track of the data needed by the worker pool service. The
// workerPool spawn multiple go routines (workers) to check and retrieve the
// receipts.
type WorkerPool struct {
	wg                   *sync.WaitGroup
	ctx                  context.Context            // Main context passed by the parent, used to cancel worker and the pool
	ethClient            interfaces.GethClient      // An interface with the Geth functionality we need
	logger               *logrus.Entry              // Logger to log messages
	txConfirmationBlocks uint64                     // Number of blocks that we should wait in order to consider a receipt valid
	requestWorkChannel   <-chan MonitorWorkRequest  // Channel where will be send the work requests
	responseWorkChannel  chan<- MonitorWorkResponse // Channel where the work response will be send
}

// Creates a new WorkerPool service
func NewWorkerPool(ctx context.Context, ethClient interfaces.GethClient, logger *logrus.Entry, txConfirmationBlocks uint64, requestWorkChannel <-chan MonitorWorkRequest, responseWorkChannel chan<- MonitorWorkResponse) *WorkerPool {
	return &WorkerPool{new(sync.WaitGroup), ctx, ethClient, logger, txConfirmationBlocks, requestWorkChannel, responseWorkChannel}
}

// Function to spawn the workers and wait for the job to be done.
func (w *WorkerPool) ExecuteWork(numWorkers uint64) {
	for i := uint64(0); i < numWorkers; i++ {
		w.wg.Add(1)
		go w.worker()
	}
	w.wg.Wait()
	close(w.responseWorkChannel)
}

// Unit of work. A worker is spawned as go routine. A worker check and retrieve
// receipts for multiple transactions. The worker will be executing while
// there's transactions to be checked or there's a timeout (set by
// constants.TxWorkerTimeout)
func (w *WorkerPool) worker() {
	ctx, cf := context.WithTimeout(w.ctx, constants.TxWorkerTimeout)
	defer cf()
	defer w.wg.Done()
	// iterating over a closed channel
	for work := range w.requestWorkChannel {
		select {
		case <-ctx.Done():
			// worker context timed out or parent was cancelled, should return
			return
		default:
			monitoredTx := work.txn
			currentHeight := work.height
			txnHash := monitoredTx.txn.Hash()
			for i := uint64(1); i <= constants.TxWorkerMaxWorkRetries; i++ {
				select {
				case <-monitoredTx.ctx.Done():
					// the go-routine who wanted this information has stopped caring. This most
					// likely indicates a failure, and cancellation of polling prevents a memory
					// leak
					w.responseWorkChannel <- MonitorWorkResponse{txnHash: txnHash, err: &ErrNonRecoverable{fmt.Sprintf("context for tx %v is finished!", txnHash.Hex())}}
					//should continue getting other tx work
					break
				case <-ctx.Done():
					// worker context timed out or parent was cancelled, should return
					return
				default:
					rcpt, err := w.getReceipt(ctx, monitoredTx, currentHeight, txnHash)
					if err != nil {
						// retry on recoverable error `constants.TxWorkerMaxWorkRetries` times
						if _, ok := err.(*ErrRecoverable); ok && i < constants.TxWorkerMaxWorkRetries {
							continue
						}
						// send nonRecoverable errors back to main or recoverable errors after constants.TxWorkerMaxWorkRetries
						w.responseWorkChannel <- MonitorWorkResponse{txnHash: txnHash, err: err}
					} else {
						// send receipt (even if it nil) back to main thread
						w.responseWorkChannel <- MonitorWorkResponse{txnHash: txnHash, receipt: rcpt}
					}
					//should continue getting other tx work
					break
				}
			}
		}
	}
}

// Internal function used by the workers to check/retrieve the receipts for a given transaction
func (w *WorkerPool) getReceipt(ctx context.Context, monitoredTx TransactionInfo, currentHeight uint64, txnHash common.Hash) (*types.Receipt, error) {
	txnHex := txnHash.Hex()
	blockTimeSpan := currentHeight - monitoredTx.startedMonitoringHeight
	_, isPending, err := w.ethClient.TransactionByHash(ctx, txnHash)
	if err != nil {
		// if we couldn't locate a tx after NotFoundMaxBlocks blocks and we are still
		// failing in getting the tx data, probably means that it was dropped
		if blockTimeSpan >= constants.TxNotFoundMaxBlocks {
			return nil, &ErrNonRecoverable{fmt.Sprintf("could not find tx %v and %v blocks have passed!", txnHex, constants.TxNotFoundMaxBlocks)}
		}
		// probably a network error, should retry
		return nil, &ErrRecoverable{fmt.Sprintf("error getting tx: %v: %v", txnHex, err)}
	}
	if isPending {
		if blockTimeSpan >= constants.TxMaxStaleBlocks {
			return nil, &ErrTransactionStale{fmt.Sprintf("error tx: %v is stale on the memory pool for more than %v blocks, please retry!", txnHex, constants.TxMaxStaleBlocks)}
		}
	} else {
		// tx is not pending, so check for receipt
		rcpt, err := w.ethClient.TransactionReceipt(ctx, txnHash)
		if err != nil {
			// if we couldn't locate a tx receipt after NotFoundMaxBlocks blocks and we are still
			// failing in getting the tx data, probably means that it was dropped
			if blockTimeSpan >= constants.TxNotFoundMaxBlocks {
				return nil, &ErrNonRecoverable{fmt.Sprintf("could not find receipt for tx %v and %v blocks have passed!", txnHex, constants.TxNotFoundMaxBlocks)}
			}
			// 1. probably a network error, should retry
			// 2. in case receipt not found after tx not Pending check, we had an edge case,
			// probably tx was mined (isPending == false), then we had a chain re-org, now
			// the tx is back to the memPool or was dropped, we should retry, and the logic
			// above should see if the tx was dropped or not in the next iteration
			return nil, &ErrRecoverable{fmt.Sprintf("error getting receipt: %v: %v", txnHex, err)}
		}

		if currentHeight >= rcpt.BlockNumber.Uint64()+w.txConfirmationBlocks {
			return rcpt, nil
		}
	}
	return nil, nil
}

// Struct that has the data necessary by the Transaction Watcher service. The
// transaction watcher service is responsible for check, retrieve and cache
// transaction receipts.
type TransactionWatcher struct {
	backend          *WatcherBackend       // backend service responsible for check, retrieving and caching the receipts
	logger           *logrus.Entry         // logger used to log the message for the transaction watcher
	closeMainContext context.CancelFunc    // function used to cancel the main context in the backend service
	requestChannel   chan<- MonitorRequest // channel used to send request to the backend service to retrieve transactions
}

// Creates a new transaction watcher struct
func NewTransactionWatcher(client interfaces.GethClient, selectMap interfaces.SelectorMap, txConfirmationBlocks uint64) *TransactionWatcher {
	requestChannel := make(chan MonitorRequest, 100)
	// main context that will cancel all workers and go routine
	mainCtx, cf := context.WithCancel(context.Background())

	backend := &WatcherBackend{
		mainCtx:              mainCtx,
		requestChannel:       requestChannel,
		client:               client,
		logger:               logging.GetLogger("ethereum").WithField("Component", "TransactionWatcherBackend"),
		monitoredTxns:        make(map[common.Hash]TransactionInfo),
		receiptCache:         make(map[common.Hash]ReceiptInfo),
		aggregates:           make(map[interfaces.FuncSelector]TransactionProfile),
		knownSelectors:       selectMap,
		lastProcessedBlock:   &BlockInfo{0, common.HexToHash("")},
		txConfirmationBlocks: txConfirmationBlocks,
	}

	transactionWatcher := &TransactionWatcher{
		requestChannel:   requestChannel,
		closeMainContext: cf,
		backend:          backend,
		logger:           logging.GetLogger("ethereum").WithField("Component", "TransactionWatcher"),
	}
	return transactionWatcher
}

// Start the transaction watcher service
func (f *TransactionWatcher) StartLoop() {
	go f.backend.Loop()
}

// Close the transaction watcher service
func (f *TransactionWatcher) Close() {
	f.logger.Debug("closing request channel...")
	close(f.requestChannel)
	f.closeMainContext()
}

// Queue a transaction to be watched by the transaction watcher service. If a
// transaction was accepted to be watched, a response channel is returned. The
// response channel is where the receipt going to be sent by the tx watcher
// backend.
func (tw *TransactionWatcher) QueueTransaction(ctx context.Context, txn *types.Transaction) (<-chan *ReceiptResponse, error) {
	tw.logger.WithField("Txn", txn.Hash().Hex()).Debug("Queueing a transaction watcher")
	respChannel := NewResponseChannel[MonitorResponse](tw.logger)
	defer respChannel.CloseChannel()
	req := MonitorRequest{ctx: ctx, txn: txn, responseChannel: respChannel}

	select {
	case tw.requestChannel <- req:
	case <-ctx.Done():
		return nil, &ErrInvalidTransactionRequest{fmt.Sprintf("context cancelled reqChannel: %v", ctx.Err())}
	}

	select {
	case requestResponse := <-req.responseChannel.channel:
		return requestResponse.receiptResponseChannel.channel, requestResponse.err
	case <-ctx.Done():
		return nil, &ErrInvalidTransactionRequest{fmt.Sprintf("context cancelled: %v", ctx.Err())}
	}
}

// function that wait for a transaction receipt. This is blocking function that will wait for a response in the input IResponse channel
func (f *TransactionWatcher) WaitTransaction(ctx context.Context, receiptResponseChannel <-chan *ReceiptResponse) (*types.Receipt, error) {
	select {
	case receiptResponse := <-receiptResponseChannel:
		return receiptResponse.receipt, receiptResponse.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Queue a transaction and wait for its receipt
func (f *TransactionWatcher) QueueAndWait(ctx context.Context, txn *types.Transaction) (*types.Receipt, error) {
	receiptResponseChannel, err := f.QueueTransaction(ctx, txn)
	if err != nil {
		return nil, err
	}
	return f.WaitTransaction(ctx, receiptResponseChannel)
}

// function to compute the max between 2 uint64
func max(a uint64, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

// function to compute the min between 2 uint64
func min(a uint64, b uint64) uint64 {
	if a > b {
		return b
	}
	return a
}

func isTestRun() bool {
	return flag.Lookup("test.v") != nil
}
