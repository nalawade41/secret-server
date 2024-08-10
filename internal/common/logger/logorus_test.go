package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// TestHook is a logrus hook to capture logs for testing
type TestHook struct {
	Entries []*logrus.Entry
}

// NewTestHook creates a new TestHook
func NewTestHook() *TestHook {
	return &TestHook{
		Entries: []*logrus.Entry{},
	}
}

// Levels returns the log levels that this hook should capture
func (hook *TestHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire is called when a log event occurs
func (hook *TestHook) Fire(entry *logrus.Entry) error {
	hook.Entries = append(hook.Entries, entry)
	return nil
}

func TestDebug(t *testing.T) {
	hook := NewTestHook()
	logrus.StandardLogger().Hooks.Add(hook)

	// Set the log level to Debug to ensure Debug logs are captured
	logrus.SetLevel(logrus.DebugLevel)

	Debug("Debug message")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.DebugLevel, hook.Entries[0].Level)
	assert.Equal(t, "Debug message", hook.Entries[0].Message)

	logrus.StandardLogger().Hooks = make(logrus.LevelHooks)
}

func TestDebugf(t *testing.T) {
	hook := NewTestHook()
	logrus.StandardLogger().Hooks.Add(hook)

	// Set the log level to Debug to ensure Debug logs are captured
	logrus.SetLevel(logrus.DebugLevel)

	Debugf("Debug %s", "message")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.DebugLevel, hook.Entries[0].Level)
	assert.Equal(t, "Debug message", hook.Entries[0].Message)

	logrus.StandardLogger().Hooks = make(logrus.LevelHooks)
}

func TestInfo(t *testing.T) {
	hook := NewTestHook()
	logrus.StandardLogger().Hooks.Add(hook)

	Info("Info message")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.InfoLevel, hook.Entries[0].Level)
	assert.Equal(t, "Info message", hook.Entries[0].Message)

	logrus.StandardLogger().Hooks = make(logrus.LevelHooks)
}

func TestInfof(t *testing.T) {
	hook := NewTestHook()
	logrus.StandardLogger().Hooks.Add(hook)

	Infof("Info %s", "message")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.InfoLevel, hook.Entries[0].Level)
	assert.Equal(t, "Info message", hook.Entries[0].Message)

	logrus.StandardLogger().Hooks = make(logrus.LevelHooks)
}

func TestWarn(t *testing.T) {
	hook := NewTestHook()
	logrus.StandardLogger().Hooks.Add(hook)

	Warn("Warn message")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.WarnLevel, hook.Entries[0].Level)
	assert.Equal(t, "Warn message", hook.Entries[0].Message)

	logrus.StandardLogger().Hooks = make(logrus.LevelHooks)
}

func TestWarnf(t *testing.T) {
	hook := NewTestHook()
	logrus.StandardLogger().Hooks.Add(hook)

	Warnf("Warn %s", "message")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.WarnLevel, hook.Entries[0].Level)
	assert.Equal(t, "Warn message", hook.Entries[0].Message)

	logrus.StandardLogger().Hooks = make(logrus.LevelHooks)
}

func TestError(t *testing.T) {
	hook := NewTestHook()
	logrus.StandardLogger().Hooks.Add(hook)

	Error("Error message")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.Entries[0].Level)
	assert.Equal(t, "Error message", hook.Entries[0].Message)

	logrus.StandardLogger().Hooks = make(logrus.LevelHooks)
}

func TestErrorf(t *testing.T) {
	hook := NewTestHook()
	logrus.StandardLogger().Hooks.Add(hook)

	Errorf("Error %s", "message")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.Entries[0].Level)
	assert.Equal(t, "Error message", hook.Entries[0].Message)

	logrus.StandardLogger().Hooks = make(logrus.LevelHooks)
}
