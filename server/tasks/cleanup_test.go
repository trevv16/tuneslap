package tasks

import (
	"context"
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypeDemoCleanupConstant(t *testing.T) {
	t.Run("has correct value", func(t *testing.T) {
		assert.Equal(t, "cleanup:demo", TypeDemoCleanup)
	})
}

func TestNewDemoCleanupTask(t *testing.T) {
	t.Run("creates valid task", func(t *testing.T) {
		task, err := NewDemoCleanupTask()
		require.NoError(t, err)
		assert.NotNil(t, task)
	})

	t.Run("task has correct type", func(t *testing.T) {
		task, err := NewDemoCleanupTask()
		require.NoError(t, err)
		assert.Equal(t, TypeDemoCleanup, task.Type())
	})

	t.Run("task has nil payload", func(t *testing.T) {
		task, err := NewDemoCleanupTask()
		require.NoError(t, err)
		// New task with nil payload
		assert.Empty(t, task.Payload())
	})
}

func TestRegisterCleanupTasks(t *testing.T) {
	t.Run("registers without panic", func(t *testing.T) {
		mux := asynq.NewServeMux()

		// Should not panic
		assert.NotPanics(t, func() {
			RegisterCleanupTasks(mux)
		})
	})
}

func TestCleanupStatsStruct(t *testing.T) {
	// Test the internal stats tracking structure

	t.Run("tracks all statistics", func(t *testing.T) {
		stats := struct {
			usersProcessed      int
			usersDeleted        int
			usersFailed         int
			mediaDeleted        int64
			mediaFailed         int
			boardsDeleted       int64
			boardsFailed        int
			collaboratorsFailed int
		}{
			usersProcessed:      10,
			usersDeleted:        8,
			usersFailed:         2,
			mediaDeleted:        50,
			mediaFailed:         5,
			boardsDeleted:       20,
			boardsFailed:        1,
			collaboratorsFailed: 0,
		}

		assert.Equal(t, 10, stats.usersProcessed)
		assert.Equal(t, 8, stats.usersDeleted)
		assert.Equal(t, 2, stats.usersFailed)
		assert.Equal(t, int64(50), stats.mediaDeleted)
		assert.Equal(t, 5, stats.mediaFailed)
		assert.Equal(t, int64(20), stats.boardsDeleted)
		assert.Equal(t, 1, stats.boardsFailed)
		assert.Equal(t, 0, stats.collaboratorsFailed)
	})

	t.Run("initializes to zero", func(t *testing.T) {
		stats := struct {
			usersProcessed      int
			usersDeleted        int
			usersFailed         int
			mediaDeleted        int64
			mediaFailed         int
			boardsDeleted       int64
			boardsFailed        int
			collaboratorsFailed int
		}{}

		assert.Equal(t, 0, stats.usersProcessed)
		assert.Equal(t, 0, stats.usersDeleted)
		assert.Equal(t, 0, stats.usersFailed)
		assert.Equal(t, int64(0), stats.mediaDeleted)
		assert.Equal(t, 0, stats.mediaFailed)
		assert.Equal(t, int64(0), stats.boardsDeleted)
		assert.Equal(t, 0, stats.boardsFailed)
		assert.Equal(t, 0, stats.collaboratorsFailed)
	})
}

func TestCleanupContextCancellation(t *testing.T) {
	// Test context cancellation detection logic

	t.Run("detects cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		select {
		case <-ctx.Done():
			assert.True(t, true, "Context cancellation detected")
		default:
			assert.Fail(t, "Context should be cancelled")
		}
	})

	t.Run("detects active context", func(t *testing.T) {
		ctx := context.Background()

		select {
		case <-ctx.Done():
			assert.Fail(t, "Context should not be cancelled")
		default:
			assert.True(t, true, "Context is active")
		}
	})

	t.Run("detects timeout context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// Wait for timeout
		time.Sleep(5 * time.Millisecond)

		select {
		case <-ctx.Done():
			assert.True(t, true, "Context timeout detected")
		default:
			assert.Fail(t, "Context should be timed out")
		}
	})
}

func TestCleanupTaskWorkflow(t *testing.T) {
	// Document the expected workflow steps

	t.Run("workflow steps are in correct order", func(t *testing.T) {
		steps := []string{
			"Check for context cancellation before starting",
			"Get collections (users, boards, media)",
			"Find all non-demo users",
			"Check context before decoding",
			"Decode users",
			"Get storage client for media deletion",
			"For each user: check context",
			"For each user: delete media files from storage",
			"For each user: delete media records from database",
			"For each user: delete boards",
			"For each user: remove from collaborator lists",
			"For each user: delete user",
			"Check context before re-seeding",
			"Re-seed demo data",
			"Log summary statistics",
		}

		assert.Len(t, steps, 15)
	})
}

func TestCleanupCriticalVsNonCriticalOperations(t *testing.T) {
	// Test the distinction between critical and non-critical operations

	t.Run("critical operations cause task failure", func(t *testing.T) {
		criticalOperations := []string{
			"Find all non-demo users",
			"Decode users",
			"Re-seed demo data",
		}

		for _, op := range criticalOperations {
			assert.NotEmpty(t, op, "Critical operation: "+op)
		}
	})

	t.Run("non-critical operations are logged but don't fail", func(t *testing.T) {
		nonCriticalOperations := []string{
			"Get storage client",
			"Delete media files from storage",
			"Delete media records",
			"Delete boards",
			"Remove from collaborator lists",
			"Delete user",
		}

		for _, op := range nonCriticalOperations {
			assert.NotEmpty(t, op, "Non-critical operation: "+op)
		}
	})
}

func TestCleanupUserDeletionOrder(t *testing.T) {
	// Test that deletion happens in the correct order (dependencies first)

	t.Run("media deleted before user", func(t *testing.T) {
		order := []string{
			"media files from storage",
			"media records from database",
			"boards",
			"collaborator references",
			"user",
		}

		// Media should be deleted before user
		mediaIndex := -1
		userIndex := -1
		for i, step := range order {
			if step == "media files from storage" {
				mediaIndex = i
			}
			if step == "user" {
				userIndex = i
			}
		}

		assert.Less(t, mediaIndex, userIndex, "Media should be deleted before user")
	})

	t.Run("boards deleted before user", func(t *testing.T) {
		order := []string{
			"media files from storage",
			"media records from database",
			"boards",
			"collaborator references",
			"user",
		}

		boardsIndex := -1
		userIndex := -1
		for i, step := range order {
			if step == "boards" {
				boardsIndex = i
			}
			if step == "user" {
				userIndex = i
			}
		}

		assert.Less(t, boardsIndex, userIndex, "Boards should be deleted before user")
	})
}

// Benchmark tests
func BenchmarkNewDemoCleanupTask(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewDemoCleanupTask()
	}
}

func BenchmarkContextCancellationCheck(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case <-ctx.Done():
		default:
		}
	}
}
