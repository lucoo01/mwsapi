package gogroup

import (
	"sync"
)

// syncWorker is an object that allows a set of asyncronous jobs
// to be run in a grouped fashion
type syncWorker struct {
	closed bool      // boolean indiciating if we are done
	done   chan bool // channel that gets triggered when done

	workGroup *sync.WaitGroup // group for adding

	hasActiveChannel bool  // do we have an active channel
	activeChannel    int64 // number of active channel

	messages chan *syncMessage // for receiving new messages
	backlog  []*syncMessage    // backlog of messages
}

// newSyncWorker makes a new SyncWorker with the given capacity
func newSyncWorker() (worker *syncWorker) {
	worker = &syncWorker{
		closed:   false,
		done:     make(chan bool),
		messages: make(chan *syncMessage, 1),
	}
	go worker.workerThread()
	return
}

// sync message represents a message
// passed from the main to the worker thread
type syncMessage struct {
	channel int64   // channel this message should be sent on
	close   bool    // if true, close the channel afterwards
	code    *func() // code to run (if any)
}

// SendWork instructs the SyncWorker to run the given code on the given channel
func (worker *syncWorker) SendWork(channel int64, code func()) {
	worker.messages <- &syncMessage{
		channel: channel,
		code:    &code,
	}
}

// SendClose closes a channel and informs the Worker
// that it may not continue sending messages on different channels
func (worker *syncWorker) SendClose(channel int64) {
	worker.messages <- &syncMessage{
		channel: channel,
		close:   true,
	}
}

// Close closes this worker
func (worker *syncWorker) Close() error {
	close(worker.messages)
	return nil
}

// Wait waits for processing to complete
func (worker *syncWorker) Wait() {
	worker.Close()
	<-worker.done
	return
}

// workerThread performs the work
// should only be called once
func (worker *syncWorker) workerThread() {
	// grab the current backlog
	backlog := worker.backlog
	worker.backlog = nil

	for i, m := range backlog {
		closed := worker.processMessage(m)
		if closed {
			worker.backlog = append(worker.backlog, backlog[(i+1):]...)
			go worker.workerThread()
			return
		}
	}

	for m := range worker.messages {
		closed := worker.processMessage(m)
		if closed {
			go worker.workerThread()
			return
		}
	}

	worker.done <- true
	close(worker.done)
}

// processMessage processes a single message
func (worker *syncWorker) processMessage(message *syncMessage) bool {

	// if we don't have an active channel set it to the message
	if !(worker.hasActiveChannel) {
		worker.hasActiveChannel = true
		worker.activeChannel = message.channel
	}

	// if it is a different channel, append it to the backlog
	if message.channel != worker.activeChannel {
		worker.backlog = append(worker.backlog, message)
		return false
	}

	// if we have some code, run it
	if message.code != nil {
		(*message.code)()
	}

	// if we are asked to close the channel, close it
	if message.close {
		worker.hasActiveChannel = false
		return true
	}

	return false
}
