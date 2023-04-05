package checkpoint

import (
	"encoding/binary"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
)

var (
	DefaultBasePath = ".checkpoint"

	defaultValue = uint64(0)
	extension    = ".cp"
)

type CheckpointHandler interface {
	SetCheckpoint(uint64) error
	Increase() error
	Decrease() error

	CheckpointReader
}

type CheckpointReader interface {
	Checkpoint() uint64
}

type Checkpoint struct {
	path string
	kind string
	n    uint64
}

func New(basePath string, kind string) *Checkpoint {
	path, err := defaultPath(basePath, runtime.GOOS)
	if err != nil {
		panic(err)
	}

	if _, err := ioutil.ReadDir(path); err != nil {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			panic(err)
		}
		return &Checkpoint{path, kind, defaultValue}
	}

	n, err := ioutil.ReadFile(filepath.Join(path, filepath.Base(kind+extension)))
	if err != nil {
		return &Checkpoint{path, kind, defaultValue}
	}

	return &Checkpoint{path, kind, binary.BigEndian.Uint64(n)}
}

func (c *Checkpoint) Checkpoint() uint64 {
	return atomic.LoadUint64(&c.n)
}

func (c *Checkpoint) SetCheckpoint(n uint64) error {
	if err := c.write(uint64ToBytes(n)); err != nil {
		return err
	}

	atomic.StoreUint64(&c.n, n)
	return nil
}

func (c *Checkpoint) Increase() error {
	if err := c.write(uint64ToBytes(atomic.LoadUint64(&c.n) + 1)); err != nil {
		return err
	}

	atomic.AddUint64(&c.n, 1)
	return nil
}

func (c *Checkpoint) Decrease() error {
	n := atomic.LoadUint64(&c.n)

	if err := c.write(uint64ToBytes(n - 1)); err != nil {
		return err
	}

	atomic.StoreUint64(&c.n, n-1)
	return nil
}

func (c *Checkpoint) write(b []byte) error {
	return ioutil.WriteFile(filepath.Join(c.path, filepath.Base(c.kind+extension)), b, fs.FileMode(0644))
}

func uint64ToBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}

// There is little possibility of adding a separate logic for each
// OS other than the path.
func defaultPath(base string, os string) (string, error) {
	switch os {
	case "darwin":
		fallthrough
	case "freebsd":
		fallthrough
	case "linux":
		return "./" + base, nil
	case "windows":
		return ".\\" + base, nil
	default:
		return "", errors.New("not supported OS: " + os)
	}
}
