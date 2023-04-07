package memdb

import (
	"errors"
	"sync"

	"github.com/dbadoy/grinder/pkg/database"
)

var _ database.Database = (*MemoryDB)(nil)

type MemoryDB struct {
	mu sync.Mutex
	v  map[interface{}]database.Data
}

func New() *MemoryDB {
	return &MemoryDB{v: make(map[interface{}]database.Data)}
}

func (m *MemoryDB) Get(key []byte) database.Data {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	return m.v[string(key)]
}

func (m *MemoryDB) Size() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.v)
}

// database.Database
func (m *MemoryDB) HealthCheck() error {
	return nil
}

func (m *MemoryDB) Insert(key []byte, data database.Data) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	lk := string(key)

	if _, ok := m.v[lk]; ok {
		return database.ErrAlreadyExist
	}

	m.v[lk] = data
	return nil
}

func (m *MemoryDB) Put(key []byte, data database.Data) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.v[string(key)] = data
	return nil
}

func (m *MemoryDB) Delete(key []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	lk := string(key)

	if _, ok := m.v[lk]; !ok {
		return errors.New("not exist key")
	}

	delete(m.v, lk)
	return nil
}

func (m *MemoryDB) Exist(_ string, key []byte) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.v[string(key)]
	return ok, nil
}
