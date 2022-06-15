package transaction

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/MadBase/MadNet/blockchain/ethereum"
	"github.com/MadBase/MadNet/bridge/bindings"
	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/constants"
	"github.com/MadBase/MadNet/constants/dbprefix"
	"github.com/MadBase/MadNet/utils"
	"github.com/dgraph-io/badger/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

type FuncSelector [4]byte

// Internal struct to keep track of transactions that are being monitoring
type info struct {
	Txn               *types.Transaction `json:"txn"`               // Transaction object
	FromAddress       common.Address     `json:"fromAddress"`       // address of the transaction signer
	Selector          FuncSelector       `json:"selector"`          // 4 bytes that identify the function being called by the tx
	FunctionSignature string             `json:"functionSignature"` // function signature as we see on the smart contracts
	RetryGroup        common.Hash        `json:"retryGroup"`        // internal group Id to keep track of all tx that were created during the retry of a tx
	DisableAutoRetry  bool               `json:"disableAutoRetry"`  // whether we should disable the auto retry of a transaction

	MonitoringHeight uint64        `json:"monitoringHeight"` // ethereum height where we first added the tx to be watched or did a tx retry.
	RetryAmount      uint64        `json:"retryAmount"`      // counter to indicate how many times we tried to retry a transaction
	NotFoundBlocks   uint64        `json:"notFoundBlocks"`   // counter to indicate approximate number of blocks that we could not find a tx
	logger           *logrus.Entry `json:"-"`                // logger to log transaction info
}

// Internal struct to keep track of transactions retries groups
type group struct {
	InternalGroup   []common.Hash    `json:"internalGroup"` // slice where we keep track of all tx in a group
	receiptResponse *ReceiptResponse `json:"-"`             // struct used to send/share the receipt
}

// creates a new group
func newGroup() group {
	return group{receiptResponse: newReceiptResponse()}
}

// add a new hash to the group
func (g *group) add(txHash common.Hash) {
	g.InternalGroup = append(g.InternalGroup, txHash)
}

// remove a hash from the group
func (g *group) remove(txHash common.Hash) error {
	index := -1
	lastIndex := len(g.InternalGroup) - 1
	if lastIndex == -1 {
		return fmt.Errorf("invalid removal, empty group %v", txHash.Hex())
	}
	for i, internalInfo := range g.InternalGroup {
		if bytes.Equal(internalInfo.Bytes(), txHash.Bytes()) {
			index = i
		}
	}
	if index == -1 {
		return fmt.Errorf("txInfo %v not found", txHash.Hex())
	}
	if index != lastIndex {
		// copy the last element in the index that we want to delete
		g.InternalGroup[index] = g.InternalGroup[lastIndex]
	}
	// drop the last index
	g.InternalGroup = g.InternalGroup[:lastIndex]
	return nil
}

// check if a group is empty
func (g *group) isEmpty() bool {
	return len(g.InternalGroup) == 0
}

// send a receipt inc ase this group has an unique tx or we have the receipt
func (g *group) sendReceipt(logger *logrus.Entry, receipt *types.Receipt, err error) {
	if g.isEmpty() {
		logger.Trace("empty group, cannot send receipt")
		return
	}
	if err != nil {
		logger.Tracef("sending group err %v", err)
	}
	if receipt != nil {
		logger.Tracef(
			"sending response group with receipt status %v mined at block %v",
			receipt.Status,
			receipt.BlockHash,
		)
	}
	// if this is the unique tx in the retry group or we have the receipt, we are good to send the response
	if len(g.InternalGroup) == 1 || receipt != nil {
		logger.Trace("sending tx")
		// in case we are recovering the group from a serialization during a crash, receiptResponse will be nil
		if g.receiptResponse == nil {
			g.receiptResponse = newReceiptResponse()
		}
		g.receiptResponse.writeReceipt(receipt, err)
	} else {
		logger.Tracef("not sending tx since group has more than one txn, group.len: %v", len(g.InternalGroup))
	}
}

// making sure that struct conforms the interface
var _ IReceiptResponse = &ReceiptResponse{}

// Struct to send and share a receipt retrieved by the watcher
type ReceiptResponse struct {
	doneChan chan struct{}
	err      error          // response error that happened during processing
	receipt  *types.Receipt // tx receipt after txConfirmationBlocks of a tx that was not queued in txGroup
}

func newReceiptResponse() *ReceiptResponse {
	return &ReceiptResponse{doneChan: make(chan struct{}, 1)}
}

// Function to check if a receipt is ready
func (r *ReceiptResponse) IsReady() bool {
	select {
	case <-r.doneChan:
		return true
	default:
		return false
	}
}

// blocking function to get the receipt from a transaction. This function will
// block until the receipt is available and sent by the transaction watcher
// service.
func (r *ReceiptResponse) GetReceiptBlocking(ctx context.Context) (*types.Receipt, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("error waiting for receipt: %v", ctx.Err())
		case <-r.doneChan:
			return r.receipt, r.err
		}
	}
}

// function to write the receipt or error from a transaction being watched.
func (r *ReceiptResponse) writeReceipt(receipt *types.Receipt, err error) {
	if receipt == nil && err == nil {
		return
	}
	if !r.IsReady() {
		r.receipt = receipt
		r.err = err
		close(r.doneChan)
	}

}

// Internal struct to keep track of the receipts
type receipt struct {
	Receipt           *types.Receipt `json:"receipt"` // receipt object
	RetrievedAtHeight uint64         `json:"-"`       // block height where receipt was added to the cache
}

// Internal struct to keep track of what blocks we already checked during monitoring
type block struct {
	Height uint64      `json:"height"` // block height
	Hash   common.Hash `json:"hash"`   // block header hash
}

// Compare if 2 blockInfo structs are equal by comparing the height and block
// hash. Return true in case they are equal, false otherwise.
func (a *block) Equal(b *block) bool {
	return bytes.Equal(a.Hash[:], b.Hash[:]) && a.Height == b.Height
}

// Type to do subscription request against the tx watcher system. SubscribeResponseChannel should be set
type SubscribeRequest struct {
	txn              *types.Transaction        // the transaction that should watched
	disableAutoRetry bool                      // whether we should disable the auto retry of a transaction
	responseChannel  *SubscribeResponseChannel // channel where we going to send the request response
}

// creates a new subscribe request
func NewSubscribeRequest(txn *types.Transaction, disableAutoRetry bool) SubscribeRequest {
	return SubscribeRequest{txn: txn, responseChannel: NewResponseChannel(), disableAutoRetry: disableAutoRetry}
}

// blocking function to listen for the response of a subscribe request
func (a SubscribeRequest) Listen(ctx context.Context) (*ReceiptResponse, error) {
	select {
	case subscribeResponse := <-a.responseChannel.channel:
		return subscribeResponse.Response, subscribeResponse.Err
	case <-ctx.Done():
		return nil, &ErrInvalidTransactionRequest{fmt.Sprintf("context cancelled: %v", ctx.Err())}
	}
}

// Type that it's going to be used to reply a subscription request
type SubscribeResponse struct {
	Err      error            // errors that happened when processing the subscription request
	Response *ReceiptResponse // struct where the receipt from the tx monitoring will be send
}

// A response channel is basically a non-blocking channel that can only be
// written and closed once.
type SubscribeResponseChannel struct {
	writeOnce sync.Once
	channel   chan *SubscribeResponse // internal channel
}

// Create a new response channel.
func NewResponseChannel() *SubscribeResponseChannel {
	return &SubscribeResponseChannel{channel: make(chan *SubscribeResponse, 1)}
}

// send a unique response and close the internal channel. Additional calls to
// this function will be no-op
func (rc *SubscribeResponseChannel) sendResponse(response *SubscribeResponse) {
	rc.writeOnce.Do(func() {
		rc.channel <- response
		close(rc.channel)
	})
}

// Profile to keep track of gas metrics in the overall system
type Profile struct {
	AverageGas   uint64 `json:"averageGas"`
	MinimumGas   uint64 `json:"minimumGas"`
	MaximumGas   uint64 `json:"maximumGas"`
	TotalGas     uint64 `json:"totalGas"`
	TotalCount   uint64 `json:"totalCount"`
	TotalSuccess uint64 `json:"totalSuccess"`
}

// Backend struct used to monitor Ethereum transactions and retrieve their receipts
type WatcherBackend struct {
	mainCtx            context.Context          `json:"-"`             // main context for the background services
	lastProcessedBlock *block                   `json:"-"`             // Last ethereum block that we checked for receipts
	MonitoredTxns      map[common.Hash]info     `json:"monitoredTxns"` // Map of transactions whose receipts we're looking for
	ReceiptCache       map[common.Hash]receipt  `json:"receiptCache"`  // Receipts retrieved from transactions. The keys are are txGroup hashes
	Aggregates         map[FuncSelector]Profile `json:"aggregates"`    // Struct to keep track of the gas metrics used by the system
	RetryGroups        map[common.Hash]group    `json:"retryGroups"`   // Map of groups of transactions that were retried
	client             ethereum.Network         `json:"-"`             // An interface with the ethereum functionality we need
	logger             *logrus.Entry            `json:"-"`             // Logger to log messages
	requestChannel     <-chan SubscribeRequest  `json:"-"`             // Channel used to send request to this backend service
	database           *db.Database             `json:"-"`             // database where we are going to persist and load state
	metricsDisplay     bool                     `json:"-"`             // flag to display the metrics in the logs. The metrics are still collect even if this flag is false.
}

// Creates a new watcher backend
func newWatcherBackend(mainCtx context.Context, requestChannel <-chan SubscribeRequest, client ethereum.Network, logger *logrus.Logger, database *db.Database, metricsDisplay bool) *WatcherBackend {
	return &WatcherBackend{
		mainCtx:            mainCtx,
		requestChannel:     requestChannel,
		client:             client,
		logger:             logger.WithField("Component", "TransactionWatcherBackend"),
		database:           database,
		MonitoredTxns:      make(map[common.Hash]info),
		ReceiptCache:       make(map[common.Hash]receipt),
		Aggregates:         make(map[FuncSelector]Profile),
		RetryGroups:        make(map[common.Hash]group),
		lastProcessedBlock: &block{0, common.HexToHash("")},
		metricsDisplay:     metricsDisplay,
	}
}

func (wb *WatcherBackend) LoadState() error {
	if err := wb.database.View(func(txn *badger.Txn) error {
		key := dbprefix.PrefixTransactionWatcherState()
		wb.logger.WithField("Key", string(key)).Tracef("Looking up state")
		rawData, err := utils.GetValue(txn, key)
		if err != nil {
			return fmt.Errorf("failed to get value %v", err)
		}
		err = json.Unmarshal(rawData, wb)
		if err != nil {
			return fmt.Errorf("failed to unmarshal %v", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (wb *WatcherBackend) PersistState() error {
	rawData, err := json.Marshal(wb)
	if err != nil {
		return fmt.Errorf("failed to marshal %v", err)
	}
	err = wb.database.Update(func(txn *badger.Txn) error {
		key := dbprefix.PrefixTransactionWatcherState()
		wb.logger.WithField("Key", string(key)).Tracef("Saving state")
		if err := utils.SetValue(txn, key, rawData); err != nil {
			return fmt.Errorf("failed to set Value %v", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// synchronizing db state to disk
	if err := wb.database.Sync(); err != nil {
		return fmt.Errorf("Failed to set sync %v", err)
	}
	return nil
}

func (wb *WatcherBackend) Loop() {
	poolingTime := time.After(constants.TxPollingTime)
	statusTime := time.After(constants.TxStatusTime)
	for {
		select {
		case req, ok := <-wb.requestChannel:
			if !ok {
				wb.logger.Debugf("request channel closed, exiting")
				return
			}
			if req.responseChannel == nil {
				wb.logger.Debug("Invalid request for txn without a response channel, ignoring")
				continue
			}
			resp, err := wb.queue(req)
			req.responseChannel.sendResponse(&SubscribeResponse{Err: err, Response: resp})

		case <-poolingTime:
			wb.collectReceipts()
			poolingTime = time.After(constants.TxPollingTime)
			err := wb.PersistState()
			if err != nil {
				wb.logger.Errorf("Failed to persist state on the database %v", err)
			}

		case <-statusTime:
			if wb.metricsDisplay {
				for selector, profile := range wb.Aggregates {
					sig := bindings.FunctionMapping[selector]
					wb.logger.WithField("Selector", fmt.Sprintf("%x", selector)).
						WithField("Function", sig).
						WithField("Profile", fmt.Sprintf("%+v", profile)).
						Info("Status")
				}
			}
			statusTime = time.After(constants.TxStatusTime)
		}
	}
}

func (wb *WatcherBackend) queue(req SubscribeRequest) (*ReceiptResponse, error) {
	if req.txn == nil {
		return nil, &ErrInvalidMonitorRequest{"invalid request, missing txn object"}
	}
	txnHash := req.txn.Hash()
	fromAddr, err := wb.client.ExtractTransactionSender(req.txn)
	if err != nil {
		// faulty transaction
		return nil, &ErrInvalidMonitorRequest{fmt.Sprintf("cannot extract fromAddr from transaction: %v", txnHash)}
	}

	var receiptResponse *ReceiptResponse
	// if we already have the records of the receipt for this tx, we try to send the
	// receipt back
	if receipt, ok := wb.ReceiptCache[txnHash]; ok {
		receiptResponse = newReceiptResponse()
		receiptResponse.writeReceipt(receipt.Receipt, nil)
	} else {
		var txGroupHash common.Hash
		if _, ok = wb.MonitoredTxns[txnHash]; ok {
			txGroupHash = wb.MonitoredTxns[txnHash].RetryGroup
		} else if _, ok = wb.RetryGroups[txnHash]; ok {
			txGroupHash = txnHash
		} else {
			selector, err := ExtractSelector(req.txn.Data())
			if err != nil {
				return nil, &ErrInvalidTransactionRequest{
					fmt.Sprintf("invalid request, transaction data is not present %v, err %v!", txnHash.Hex(), err),
				}
			}
			sig := bindings.FunctionMapping[selector]

			logEntry := wb.logger.WithField("Transaction", txnHash).
				WithField("Function", sig).
				WithField("Selector", fmt.Sprintf("%x", selector))

			wb.MonitoredTxns[txnHash] = info{
				Txn:               req.txn,
				FromAddress:       fromAddr,
				Selector:          selector,
				FunctionSignature: sig,
				RetryGroup:        txnHash,
				DisableAutoRetry:  req.disableAutoRetry,
				logger:            logEntry,
			}
			txGroup := newGroup()
			txGroup.add(txnHash)
			wb.RetryGroups[txnHash] = txGroup
			logEntry.Debug("Transaction queued")
			txGroupHash = txnHash
		}
		receiptResponse = wb.RetryGroups[txGroupHash].receiptResponse
	}
	return receiptResponse, nil
}

func (wb *WatcherBackend) collectReceipts() {

	lenMonitoredTxns := len(wb.MonitoredTxns)

	// If there's no tx to be monitored just return
	if lenMonitoredTxns == 0 {
		wb.logger.Tracef("no transaction to watch")
		return
	}

	networkCtx, cf := context.WithTimeout(wb.mainCtx, constants.TxNetworkTimeout)
	defer cf()

	blockHeader, err := wb.client.GetHeaderByNumber(networkCtx, nil)
	if err != nil {
		wb.logger.Debugf("error getting latest block number from ethereum node: %v", err)
		return
	}
	blockInfo := &block{
		blockHeader.Number.Uint64(),
		blockHeader.Hash(),
	}

	if wb.lastProcessedBlock.Equal(blockInfo) {
		wb.logger.Tracef("already processed block: %v with hash: %v", blockInfo.Height, blockInfo.Hash.Hex())
		return
	}
	wb.logger.Tracef("processing block: %v with hash: %v", blockInfo.Height, blockInfo.Hash.Hex())

	baseFee, tipCap, err := wb.client.GetBlockBaseFeeAndSuggestedGasTip(networkCtx)
	if err != nil {
		wb.logger.Debugf("error getting baseFee and suggested gas tip from ethereum node: %v", err)
		return
	}

	finishedTxs := make(map[common.Hash]MonitorWorkResponse)

	numWorkers := utils.Min(utils.Max(uint64(lenMonitoredTxns)/4, 128), 1)
	requestWorkChannel := make(chan MonitorWorkRequest, lenMonitoredTxns+3)
	responseWorkChannel := make(chan MonitorWorkResponse, lenMonitoredTxns+3)

	for txn, txnInfo := range wb.MonitoredTxns {
		// if this is the first time seeing a tx or we have a reorg and
		// startedMonitoring is now greater than the current ethereum block height
		if txnInfo.MonitoringHeight == 0 || txnInfo.MonitoringHeight > blockInfo.Height {
			txnInfo.MonitoringHeight = blockInfo.Height
			wb.MonitoredTxns[txn] = txnInfo
		}
		requestWorkChannel <- MonitorWorkRequest{txnInfo, blockInfo.Height}
	}

	// close the request channel, so the workers know when to finish
	close(requestWorkChannel)

	workerPool := NewWorkerPool(wb.mainCtx, wb.client, baseFee, tipCap, wb.logger, requestWorkChannel, responseWorkChannel)

	// spawn the workers and wait for all to complete
	go workerPool.ExecuteWork(numWorkers)

	for workResponse := range responseWorkChannel {
		select {
		case <-wb.mainCtx.Done():
			// main thread was killed
			return
		default:
		}
		txInfo, ok := wb.MonitoredTxns[workResponse.txnHash]
		logEntry := txInfo.logger
		if !ok {
			// invalid tx, should not happen, but well if it happens we continue
			logEntry.Trace("got a invalid tx with hash from workers")
			continue
		}
		if workResponse.err != nil {
			err := workResponse.err
			switch err.(type) {
			case *ErrRecoverable:
				logEntry.Tracef("Retrying! Got a recoverable error when trying to get receipt, err: %v", workResponse.err)
			case *ErrTxNotFound:
				// since we only analyze a tx once per new block, the notFoundBlocks counter
				// should have approx the amount of blocks that we failed on finding the tx
				txInfo.NotFoundBlocks++
				if txInfo.NotFoundBlocks >= wb.client.GetTxNotFoundMaxBlocks() {
					logEntry.Debugf("Couldn't get tx receipt, err: %v", workResponse.err)
					finishedTxs[workResponse.txnHash] = workResponse
				}
				logEntry.Tracef("Retrying, couldn't get info, num attempts: %v, err: %v", txInfo.NotFoundBlocks, workResponse.err)
			case *ErrTransactionStale:
				// If we get this error it means that we should not retry or we cannot retry
				// automatically, should forward the error to the subscribers
				logEntry.Debugf("Stale transaction, err: %v", workResponse.err)
				finishedTxs[workResponse.txnHash] = workResponse
			}
		} else {
			if workResponse.retriedTxn != nil {
				// restart the monitoringHeight, so we don't retry the tx in the next block
				txInfo.MonitoringHeight = 0
				if workResponse.retriedTxn.err == nil && workResponse.retriedTxn.txn != nil {
					newTxnHash := workResponse.retriedTxn.txn.Hash()
					wb.MonitoredTxns[newTxnHash] = info{
						Txn:               workResponse.retriedTxn.txn,
						FromAddress:       txInfo.FromAddress,
						Selector:          txInfo.Selector,
						FunctionSignature: txInfo.FunctionSignature,
						RetryGroup:        txInfo.RetryGroup,
						DisableAutoRetry:  txInfo.DisableAutoRetry,
						logger:            txInfo.logger,
					}
					// update retry group
					txGroup := wb.RetryGroups[txInfo.RetryGroup]
					txGroup.add(newTxnHash)
					wb.RetryGroups[txInfo.RetryGroup] = txGroup
					logEntry.Tracef("successfully replaced a tx with %v", newTxnHash)
					txInfo.RetryAmount++
				} else {
					logEntry.Debugf("could not replace tx error %v", workResponse.retriedTxn.err)
				}
			}
			if workResponse.receipt != nil {
				logEntry.WithFields(
					logrus.Fields{
						"mined":          workResponse.receipt.BlockNumber,
						"current height": blockInfo.Height,
					},
				).Debug("Successfully got receipt")
				wb.ReceiptCache[txInfo.RetryGroup] = receipt{Receipt: workResponse.receipt, RetrievedAtHeight: blockInfo.Height}
				finishedTxs[workResponse.txnHash] = workResponse
			}
		}
		wb.MonitoredTxns[workResponse.txnHash] = txInfo
	}

	// Cleaning finished and failed transactions
	for txnHash, workResponse := range finishedTxs {
		if txnInfo, ok := wb.MonitoredTxns[txnHash]; ok {
			if txGroup, ok := wb.RetryGroups[txnInfo.RetryGroup]; ok {
				logger := txnInfo.logger.WithFields(logrus.Fields{
					"group": txnInfo.RetryGroup,
				})

				if workResponse.receipt != nil {
					rcpt := workResponse.receipt
					var profile Profile
					if _, present := wb.Aggregates[txnInfo.Selector]; present {
						profile = wb.Aggregates[txnInfo.Selector]
					} else {
						profile = Profile{}
					}
					// Update transaction profile
					profile.AverageGas = (profile.AverageGas*profile.TotalCount + rcpt.GasUsed) / (profile.TotalCount + 1)
					if profile.MaximumGas < rcpt.GasUsed {
						profile.MaximumGas = rcpt.GasUsed
					}
					if profile.MinimumGas == 0 || profile.MinimumGas > rcpt.GasUsed {
						profile.MinimumGas = rcpt.GasUsed
					}
					profile.TotalCount++
					profile.TotalGas += rcpt.GasUsed
					if rcpt.Status == uint64(1) {
						profile.TotalSuccess++
					}
					wb.Aggregates[txnInfo.Selector] = profile
				}
				txGroup.sendReceipt(logger, workResponse.receipt, workResponse.err)
				err = txGroup.remove(txnHash)
				if err != nil {
					logger.Debugf("Failed to remove txn from group: %v", err)
				} else {
					if txGroup.isEmpty() {
						logger.Tracef("empty group removing")
						delete(wb.RetryGroups, txnInfo.RetryGroup)
					}
				}
			} else {
				txnInfo.logger.Debugf("Failed to find a group for txn")
			}
			delete(wb.MonitoredTxns, txnHash)
		} else {
			wb.logger.Debugf("Failed to find txn to remove: %v", txnHash.Hex())
		}
	}

	var expiredReceipts []common.Hash
	// Marking expired receipts and restarting the height of state recovered receipts
	for receiptTxnHash, receiptInfo := range wb.ReceiptCache {
		if receiptInfo.RetrievedAtHeight == 0 || receiptInfo.RetrievedAtHeight > blockInfo.Height {
			receiptInfo.RetrievedAtHeight = blockInfo.Height
		}
		if blockInfo.Height >= receiptInfo.RetrievedAtHeight+constants.TxReceiptCacheMaxBlocks {
			expiredReceipts = append(expiredReceipts, receiptTxnHash)
		}
	}
	for _, receiptTxHash := range expiredReceipts {
		wb.logger.Tracef("cleaning %v from receipt cache", receiptTxHash.Hex())
		delete(wb.ReceiptCache, receiptTxHash)
	}

	wb.lastProcessedBlock = blockInfo
}

func ExtractSelector(data []byte) (FuncSelector, error) {
	var selector [4]byte
	if len(data) < 4 {
		return selector, fmt.Errorf("couldn't extract selector for data: %v", data)
	}
	for idx := 0; idx < 4; idx++ {
		selector[idx] = data[idx]
	}
	return selector, nil
}