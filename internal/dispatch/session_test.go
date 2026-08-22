package dispatch_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	agentdv1 "github.com/macrox-pro/agentd/gen/agentd/v1"
	"github.com/macrox-pro/agentd/internal/dispatch"
)

func TestSessionsSameIDSerializes(t *testing.T) {
	t.Parallel()
	s := &dispatch.Sessions{}
	var concurrent atomic.Int32
	var maxConcurrent atomic.Int32

	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			unlock := s.Lock("sess-a")
			defer unlock()
			n := concurrent.Add(1)
			for {
				cur := maxConcurrent.Load()
				if n <= cur || maxConcurrent.CompareAndSwap(cur, n) {
					break
				}
			}
			time.Sleep(5 * time.Millisecond)
			concurrent.Add(-1)
		})
	}
	wg.Wait()
	assert.Equal(t, int32(1), maxConcurrent.Load())
	assert.Equal(t, uint32(0), s.Active())
}

func TestSessionsDifferentIDsParallel(t *testing.T) {
	t.Parallel()
	s := &dispatch.Sessions{}
	started := make(chan struct{})
	release := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		unlock := s.Lock("a")
		defer unlock()
		started <- struct{}{}
		<-release
	}()
	go func() {
		defer wg.Done()
		unlock := s.Lock("b")
		defer unlock()
		started <- struct{}{}
		<-release
	}()

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("first lock not acquired")
	}
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("second lock not acquired in parallel")
	}
	assert.Equal(t, uint32(2), s.Active())
	close(release)
	wg.Wait()
}

func TestEngineSessionIDOf(t *testing.T) {
	t.Parallel()
	typed, err := dispatch.DecodeTyped(context.Background(), agentdv1.Provider_PROVIDER_CLAUDE_CODE, agentdv1.InvocationMode_INVOCATION_MODE_STDIN, claudeToolPre(t, "echo"))
	require.NoError(t, err)
	assert.Equal(t, "s", dispatch.SessionIDOf(typed))
	assert.Equal(t, "", dispatch.SessionIDOf(nil))
}
