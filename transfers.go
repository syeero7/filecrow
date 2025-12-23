package main

import (
	"io"
	"sync"
)

type Session struct {
	reader *io.PipeReader
	writer *io.PipeWriter
	done   chan struct{}
}

type FileTransfer struct {
	name    string
	size    int
	session *Session
}

type Transfers struct {
	transfers map[string]*FileTransfer
	mu        sync.Mutex
}

func (t *Transfers) add(id string, ft *FileTransfer) {
	t.mu.Lock()
	t.transfers[id] = ft
	t.mu.Unlock()
}

func (t *Transfers) get(id string) (*FileTransfer, bool) {
	t.mu.Lock()
	ft, ok := t.transfers[id]
	t.mu.Unlock()
	return ft, ok
}

func (t *Transfers) remove(id string) {
	t.mu.Lock()
	delete(t.transfers, id)
	t.mu.Unlock()
}

var transfers = &Transfers{transfers: make(map[string]*FileTransfer)}
