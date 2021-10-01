package blockchain

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/MadBase/MadNet/blockchain/interfaces"
	"github.com/MadBase/MadNet/logging"
	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

//
var (
	ErrUnknownRequest = errors.New("unknown request type")
)

//
type Request struct {
	name   string
	txn    *types.Transaction
	group  int
	respch chan *Response
}

//
type Response struct {
	message string
	err     error
	rcpt    *types.Receipt
	rcpts   []*types.Receipt
}

// TransactionProfile
type TransactionProfile struct {
	AverageGas   uint64
	MinimumGas   uint64
	MaximumGas   uint64
	TotalCount   uint64
	TotalGas     uint64
	TotalSuccess uint64
}

// Behind is the struct used while monitoring Ethereum transactions
type Behind struct {
	sync.Mutex
	waitingTxns    []common.Hash                                  // Just a list of transactions whose receipts we're looking for
	readyTxns      map[common.Hash]*types.Receipt                 // All the transaction -> receipt pairs we know of
	selectors      map[common.Hash]interfaces.FuncSelector        // Maps a transaction to it's function selector
	groups         map[int][]common.Hash                          // A group is just an ID and a list of transactions
	aggregates     map[interfaces.FuncSelector]TransactionProfile //
	client         interfaces.GethClient                          // An interface with the Geth functionality we need
	knownSelectors interfaces.SelectorMap                         //
	logger         *logrus.Entry                                  //
	reqch          <-chan *Request                                //
	timeout        time.Duration                                  // How long will we wait for a receipt
}

func (b *Behind) Loop() {

	done := false
	for !done {
		select {
		case req, ok := <-b.reqch:
			// Some sort of request came in
			if !ok {
				b.logger.Debugf("command channel closed")
				b.status(nil)
				done = true
				break
			}

			b.logger.Debugf("received request: %v channel open: %v", req.name, ok)

			var handler func(*Request) *Response
			switch req.name {
			case "queue":
				handler = b.queue
			case "status":
				handler = b.status
			case "wait":
				handler = b.wait
			default:
				handler = b.unknown
			}

			go b.process(req, handler)

		case <-time.After(100 * time.Millisecond):
			// No request, so let's do some work
			b.collectReceipts()
		}
	}
	b.logger.Debug("finished")
}

func (b *Behind) collectReceipts() {
	ctx, cf := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cf()

	b.Lock()
	defer b.Unlock()

	n := len(b.waitingTxns)
	if n < 1 {
		return
	}

	// loop over transactions in need of receipts while building new list
	remainingTxns := make([]common.Hash, 0, n)
	for _, txn := range b.waitingTxns {
		rcpt, err := b.client.TransactionReceipt(ctx, txn)
		if err == geth.NotFound || (err == nil && rcpt == nil) {
			b.logger.Debugf("receipt not found: %v", txn.Hex())
		} else if err != nil {
			b.logger.Errorf("error getting receipt: %v", txn)
		} else if rcpt != nil {
			b.readyTxns[txn] = rcpt

			var profile TransactionProfile
			var selector [4]byte
			var sig string
			var present bool

			if selector, present = b.selectors[txn]; present {
				profile = b.aggregates[selector]
				sig = b.knownSelectors.Signature(selector)
			} else {
				profile = TransactionProfile{}
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

			b.aggregates[selector] = profile
			logEntry := b.logger.WithField("Transaction", rcpt.TxHash.Hex()).
				WithField("Function", sig).
				WithField("Selector", fmt.Sprintf("%x", selector)).
				WithField("Successful", rcpt.Status == 1)

			// This is hideous but useful when troubleshooting with simulator
			if b.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				fullTxn, _, err := b.client.TransactionByHash(ctx, txn)
				if err == nil {
					signer := types.NewEIP155Signer(big.NewInt(1337))
					msg, err := fullTxn.AsMessage(signer, nil)
					if err == nil {
						logEntry = logEntry.WithField("From", msg.From().Hash().Hex())
					}
				}
			}

			logEntry.Debugf("Receipt collected")
		}

		if _, present := b.readyTxns[txn]; !present {
			remainingTxns = append(remainingTxns, txn)
		}
	}

	b.logger.Debugf("started with %v txns, %v are remaining", len(b.waitingTxns), len(remainingTxns))
	b.waitingTxns = remainingTxns
}

func (b *Behind) process(req *Request, handler func(req *Request) *Response) {

	b.logger.Debug("processing request...")

	resp := handler(req)

	b.logger.Debugf("response channel present: %v", req.respch != nil)
	if req.respch != nil {
		req.respch <- resp
		close(req.respch)
	}
	b.logger.Debug("...done processing request")
}

func (b *Behind) queue(req *Request) *Response {

	b.Lock()
	defer b.Unlock()

	txnHash := req.txn.Hash()

	selector := ExtractSelector(req.txn.Data())

	sig := b.knownSelectors.Signature(selector)

	logEntry := b.logger.WithField("Transaction", txnHash).
		WithField("Function", sig).
		WithField("Selector", fmt.Sprintf("%x", selector))

	b.selectors[txnHash] = selector
	b.waitingTxns = append(b.waitingTxns, txnHash)

	if _, present := b.groups[req.group]; !present {
		b.groups[req.group] = make([]common.Hash, 0, 10)
	}
	b.groups[req.group] = append(b.groups[req.group], txnHash)

	// This is hideous but useful when troubleshooting with simulator
	if b.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		signer := types.NewEIP155Signer(big.NewInt(1337))
		msg, err := req.txn.AsMessage(signer, nil)
		if err == nil {
			logEntry = logEntry.WithField("From", msg.From().Hash().Hex())
		}
	}

	logEntry.Debug("Transaction queued")

	return &Response{message: "queued transaction"}
}

func (b *Behind) status(req *Request) *Response {
	b.Lock()
	defer b.Unlock()

	b.logger.WithField("Completed", len(b.readyTxns)).
		WithField("Pending", len(b.waitingTxns)).
		Info("Transaction counts")

	for selector, profile := range b.aggregates {
		sig := b.knownSelectors.Signature(selector)
		b.logger.WithField("Selector", fmt.Sprintf("%x", selector)).
			WithField("Function", sig).
			WithField("Profile", fmt.Sprintf("%+v", profile)).
			Info("Status")
	}

	return &Response{message: "status check"}
}

func (b *Behind) wait(req *Request) *Response {

	ctx, cf := context.WithTimeout(context.Background(), b.timeout)
	defer cf()

	resp := &Response{message: "status check"}
	done := false

	check := func() {
		b.Lock()
		defer b.Unlock()

		if req.txn != nil {
			// waiting for a specific transaction to complete
			txn := req.txn.Hash()
			if rcpt, present := b.readyTxns[txn]; present {
				resp.rcpt = rcpt
				// delete(b.readyTxns, txn) // TODO Add an explicit purge/cleanup
				done = true
			} else {
				b.logger.Debugf("rcpt not ready yet for %v", txn.Hex())
			}
		} else {
			// waiting for a group of transactions to complete
			allPresent := true
			for _, txn := range b.groups[req.group] {
				if _, present := b.readyTxns[txn]; !present {
					allPresent = false
				}
			}

			if allPresent {
				resp.rcpts = make([]*types.Receipt, 0, len(b.groups[req.group]))
				for _, txn := range b.groups[req.group] {
					resp.rcpts = append(resp.rcpts, b.readyTxns[txn])
					// delete(b.readyTxns, txn) // TODO Add an explicit purge/cleanup
				}
				// delete(b.groups, req.group) // TODO Add an explicit purge/cleanup
				done = true
			}
		}
	}

	check()
	for !done {
		select {
		case <-time.After(500 * time.Millisecond):
			check()
		case <-ctx.Done():
			done = true
			resp.err = ctx.Err()
		}
	}

	return resp
}

func (b *Behind) unknown(req *Request) *Response {
	return &Response{err: ErrUnknownRequest}
}

type TxnQueueDetail struct {
	backend *Behind
	logger  *logrus.Entry
	reqch   chan<- *Request
}

func NewTxnQueue(client interfaces.GethClient, sm interfaces.SelectorMap, to time.Duration) *TxnQueueDetail {
	reqch := make(chan *Request, 100)

	b := &Behind{
		reqch:          reqch,
		client:         client,
		logger:         logging.GetLogger("ethereum").WithField("Component", "behind"),
		waitingTxns:    make([]common.Hash, 0, 20),
		readyTxns:      make(map[common.Hash]*types.Receipt),
		selectors:      make(map[common.Hash]interfaces.FuncSelector),
		aggregates:     make(map[interfaces.FuncSelector]TransactionProfile),
		knownSelectors: sm,
		timeout:        to,
		groups:         make(map[int][]common.Hash)}

	q := &TxnQueueDetail{
		reqch:   reqch,
		backend: b,
		logger:  logging.GetLogger("ethereum").WithField("Component", "infront")}

	return q
}

func (f *TxnQueueDetail) StartLoop() {
	go f.backend.Loop()
}

func (f *TxnQueueDetail) QueueTransaction(ctx context.Context, txn *types.Transaction) {
	f.logger.WithField("Txn", string(txn.Hash().Bytes())).Debug("Queueing")
	req := &Request{name: "queue", txn: txn} // no response channel because I don't want to wait
	f.requestWait(ctx, req)
}

func (f *TxnQueueDetail) QueueGroupTransaction(ctx context.Context, grp int, txn *types.Transaction) {
	f.logger.WithFields(logrus.Fields{
		"Txn":   string(txn.Hash().Bytes()),
		"Group": grp}).Debug("Queueing for group")
	req := &Request{name: "queue", txn: txn, group: grp} // no response channel because I don't want to wait
	f.requestWait(ctx, req)
}

func (f *TxnQueueDetail) QueueAndWait(ctx context.Context, txn *types.Transaction) (*types.Receipt, error) {
	f.QueueTransaction(ctx, txn)
	return f.WaitTransaction(ctx, txn)
}

func (f *TxnQueueDetail) WaitTransaction(ctx context.Context, txn *types.Transaction) (*types.Receipt, error) {

	f.logger.WithField("Txn", string(txn.Hash().Bytes())).Debug("Waiting")
	req := &Request{name: "wait", txn: txn, respch: make(chan *Response)}
	resp := f.requestWait(ctx, req)

	if resp.err != nil {
		return nil, resp.err
	}

	return resp.rcpt, nil
}

func (f *TxnQueueDetail) WaitGroupTransactions(ctx context.Context, grp int) ([]*types.Receipt, error) {
	f.logger.WithField("Group", grp).Debug("Waiting")
	req := &Request{name: "wait", group: grp, respch: make(chan *Response)}
	resp := f.requestWait(ctx, req)

	if resp.err != nil {
		return nil, resp.err
	}

	return resp.rcpts, nil
}

func (f *TxnQueueDetail) Status(ctx context.Context) error {
	req := &Request{name: "status"}
	logger := f.logger.WithField("Command", req.name)
	logger.Debug("waiting...")
	f.requestWait(ctx, req)
	logger.Debug("...done waiting")
	return nil
}

func (f *TxnQueueDetail) Close() {
	f.logger.Debug("closing request channel...")
	close(f.reqch)
}

func (f *TxnQueueDetail) requestWait(ctx context.Context, req *Request) *Response {
	f.reqch <- req

	logReciept := func(message string, rcpt *types.Receipt) {
		f.logger.WithFields(logrus.Fields{
			"Message":     message,
			"Transaction": rcpt.TxHash.Hex(),
			"Block":       rcpt.BlockNumber.String(),
			"GasUsed":     rcpt.GasUsed,
			"Status":      rcpt.Status,
		}).Debugf("Received response")
	}

	if req.respch != nil {
		select {
		case resp, ok := <-req.respch:
			if !ok {
				f.logger.Error("response channel closed")
			} else if resp != nil {
				if resp.err != nil {
					f.logger.Infof("response error: %v", resp.err.Error())
				}

				if resp.rcpt != nil {
					logReciept(resp.message, resp.rcpt)
				} else if resp.rcpts != nil {
					for _, rcpt := range resp.rcpts {
						logReciept(resp.message, rcpt)
					}
				}

			} else {
				f.logger.Warn("no response or error")
			}
			return resp

		case <-ctx.Done():
			f.logger.Infof("context cancelled: %v", ctx.Err())
			return &Response{err: ctx.Err()}
		}
	}

	return nil // no response channel, so no meaningful response
}
