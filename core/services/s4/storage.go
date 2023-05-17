package s4

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/ethereum/go-ethereum/common"
)

const (
	serviceName = "S4-Storage"
)

// Storage represents S4 interface
type Storage interface {
	job.ServiceCtx

	// Constraints returns a copy of Constraints struct specified during service creation.
	// The implementation is thread-safe.
	Constraints() Constraints

	// Get returns a copy of record (with metadata) associated with the specified address and slotId.
	// The implementation is thread-safe. The returned Record & Metadata are always a copy.
	// If no record is found for the given parameters, ErrRecordNotFound is returned.
	Get(ctx context.Context, address common.Address, slotId int) (*Record, *Metadata, error)

	// Put creates (or updates) a record identified by the specified address and slotId.
	// The implementation is thread-safe.
	// For signature calculation see envelope.go
	Put(ctx context.Context, address common.Address, slotId int, record *Record, signature []byte) error
}

type inMemoryStorage struct {
	utils.StartStopOnce

	lggr           logger.Logger
	contraints     Constraints
	expiryInterval time.Duration

	records  map[string]Record
	metadata map[string]Metadata
	doneCh   chan struct{}
	mu       sync.RWMutex
}

func NewInMemoryStorage(lggr logger.Logger, contraints Constraints, expiryInterval time.Duration) Storage {
	return &inMemoryStorage{
		lggr:           lggr.Named(serviceName),
		contraints:     contraints,
		expiryInterval: expiryInterval,
		records:        map[string]Record{},
		metadata:       map[string]Metadata{},
		doneCh:         make(chan struct{}),
	}
}

func (s *inMemoryStorage) Start(ctx context.Context) error {
	return s.StartOnce(serviceName, func() error {
		go func() {
			ticker := time.NewTicker(s.expiryInterval)
			defer ticker.Stop()

			for {
				select {
				case <-s.doneCh:
					return
				case <-ticker.C:
					s.expirationLoop()
				}
			}
		}()
		return nil
	})
}

func (s *inMemoryStorage) Close() error {
	return s.StopOnce(serviceName, func() (err error) {
		// Acquring mu allows any pending operations to complete
		s.mu.Lock()
		defer s.mu.Unlock()
		close(s.doneCh)
		return
	})
}

func (s *inMemoryStorage) Constraints() Constraints {
	return s.contraints
}

func (s *inMemoryStorage) Get(ctx context.Context, address common.Address, slotId int) (*Record, *Metadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.StartStopOnce.State() != utils.StartStopOnce_Started {
		return nil, nil, ErrServiceNotStarted
	}

	key := fmt.Sprintf("%s_%d", address, slotId)
	record, ok := s.records[key]
	if !ok {
		return nil, nil, ErrRecordNotFound
	}
	metadata, ok := s.metadata[key]
	if !ok {
		return nil, nil, ErrRecordNotFound
	}
	if metadata.State == ExpiredRecordState {
		return nil, nil, ErrRecordExpired
	}

	recordClone := record.Clone()
	metadataClone := metadata.Clone()
	return &recordClone, &metadataClone, nil
}

func (s *inMemoryStorage) Put(ctx context.Context, address common.Address, slotId int, record *Record, signature []byte) error {
	if slotId >= s.contraints.MaxSlotsPerUser {
		return ErrTooBigSlotId
	}
	if len(record.Payload) > s.contraints.MaxPayloadSizeBytes {
		return ErrTooBigPayload
	}
	if time.Now().UnixMilli() > record.Expiration {
		return ErrPastExpiration
	}

	envelope := NewEnvelopeFromRecord(address, slotId, record)
	signer, err := envelope.GetSignerAddress(signature)
	if err != nil || signer != address {
		return ErrWrongSignature
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.StartStopOnce.State() != utils.StartStopOnce_Started {
		return ErrServiceNotStarted
	}

	key := fmt.Sprintf("%s_%d", address, slotId)

	existing, ok := s.records[key]
	if ok && existing.Version <= record.Version {
		return ErrOlderVersion
	}

	metadata := s.metadata[key]
	metadata.State = NewRecordState
	metadata.Signature = make([]byte, len(signature))
	copy(metadata.Signature, signature)
	if metadata.HighestExpiration < record.Expiration {
		metadata.HighestExpiration = record.Expiration
	}
	s.metadata[key] = metadata
	s.records[key] = record.Clone()

	return nil
}

func (s *inMemoryStorage) expirationLoop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	for k, v := range s.records {
		if v.Expiration < now {
			m := s.metadata[k]
			m.State = ExpiredRecordState
			s.metadata[k] = m
		}
	}
}
