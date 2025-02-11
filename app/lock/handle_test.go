package lock

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createTestFile(t *testing.T) string {
	// create test file
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "testfile")
	require.NotEmpty(t, filePath)
	err := os.WriteFile(filePath, []byte("content"), 0644)
	require.NoError(t, err)
	return filePath
}

func createTestDir(t *testing.T) string {
	tmpDir := t.TempDir()
	dirPath := filepath.Join(tmpDir, "testdir")
	require.NotEmpty(t, dirPath)
	err := os.Mkdir(dirPath, 0755)
	require.NoError(t, err)
	return dirPath
}

func TestLock_Handle_FileBasic(t *testing.T) {
	filePath := createTestFile(t)

	// check executed state
	executed := false
	err := LockFile(filePath, func() error {
		executed = true
		return nil
	})
	require.NoError(t, err)
	require.True(t, executed)

	// ensure the lockfile does not exist anymore
	lockPath := getFileLockPath(filePath)
	_, err = os.Stat(lockPath)
	require.True(t, os.IsNotExist(err))
}

func TestLock_Handle_DirBasic(t *testing.T) {
	dirPath := createTestDir(t)

	// check executed state
	executed := false
	err := LockDir(dirPath, func() error {
		executed = true
		return nil
	})
	require.NoError(t, err)
	require.True(t, executed)

	// ensure the lockfile does not exist anymore
	lockPath := getDirLockPath(dirPath)
	_, err = os.Stat(lockPath)
	require.True(t, os.IsNotExist(err))
}

func TestLock_Handle_FileErrorPropagation(t *testing.T) {
	filePath := createTestFile(t)

	// ensure error propigation through the function works
	expectedErr := errors.New("test error")
	err := LockFile(filePath, func() error {
		return expectedErr
	})
	require.Equal(t, expectedErr, err)

	// ensure the lockfile does not exist anymore
	lockPath := getFileLockPath(filePath)
	_, err = os.Stat(lockPath)
	require.True(t, os.IsNotExist(err))
}

func TestLock_Handle_DirErrorPropagation(t *testing.T) {
	dirPath := createTestDir(t)

	// ensure error propigation through the function works
	expectedErr := errors.New("dir test error")
	err := LockDir(dirPath, func() error {
		return expectedErr
	})
	require.Equal(t, expectedErr, err)

	// ensure the lockfile does not exist anymore
	lockPath := getDirLockPath(dirPath)
	_, err = os.Stat(lockPath)
	require.True(t, os.IsNotExist(err))
}

func TestLock_Handle_ConcurrentLockFile(t *testing.T) {
	filePath := createTestFile(t)

	// setup wait group
	var concurrentCount, maxConcurrent int32
	var wg sync.WaitGroup
	n := 5

	// spawn go routines that attempt to all access the same file
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := LockFile(filePath, func() error {
				// when seen, add to the atomic counter
				count := atomic.AddInt32(&concurrentCount, 1)
				for {
					// if the current value is creater than the seen, then swap the
					// values to ensure we always track the highest seen across ALL
					// goroutines
					curr := atomic.LoadInt32(&maxConcurrent)
					if count > curr {
						if atomic.CompareAndSwapInt32(&maxConcurrent, curr, count) {
							break
						}
					} else {
						break
					}
				}
				time.Sleep(100 * time.Millisecond)
				// when finished, remove from the counter
				atomic.AddInt32(&concurrentCount, -1)
				return nil
			})
			require.NoError(t, err)
		}()
	}
	wg.Wait()

	// ensure no concurrent access
	require.Less(t, maxConcurrent, int32(2))
}

func TestLock_Handle_ConcurrentLockDir(t *testing.T) {
	dirPath := createTestDir(t)

	var concurrentCount, maxConcurrent int32
	var wg sync.WaitGroup
	n := 5

	// spawn go routines that attempt to all access the same dir
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := LockDir(dirPath, func() error {
				// when seen, add to the atomic counter
				count := atomic.AddInt32(&concurrentCount, 1)
				for {
					// if the current value is creater than the seen, then swap the
					// values to ensure we always track the highest seen across ALL
					// goroutines
					curr := atomic.LoadInt32(&maxConcurrent)
					if count > curr {
						if atomic.CompareAndSwapInt32(&maxConcurrent, curr, count) {
							break
						}
					} else {
						break
					}
				}
				time.Sleep(100 * time.Millisecond)
				// when finished, remove from the counter
				atomic.AddInt32(&concurrentCount, -1)
				return nil
			})
			require.NoError(t, err)
		}()
	}
	wg.Wait()

	// ensure no concurrent access
	require.Less(t, maxConcurrent, int32(2))
}
