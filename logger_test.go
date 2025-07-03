package sendgrid

import (
	"bytes"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	buf := bytes.NewBufferString("")
	logger := internalLog{logger: log.New(buf, "", 0|log.Lshortfile)}
	logger.Println("test line 123")
	assert.Equal(t, buf.String(), "logger_test.go:15: test line 123\n")
	buf.Truncate(0)
	logger.Print("test line 123")
	assert.Equal(t, buf.String(), "logger_test.go:18: test line 123\n")
	buf.Truncate(0)
	logger.Printf("test line 123\n")
	assert.Equal(t, buf.String(), "logger_test.go:21: test line 123\n")
	buf.Truncate(0)
	if err := logger.Output(1, "test line 123\n"); err != nil {
		log.Println(err)
	}
	assert.Equal(t, buf.String(), "logger_test.go:24: test line 123\n")
	buf.Truncate(0)
}

// errorLogger implements logger interface but returns error on Output
type errorLogger struct {
	shouldError bool
}

func (m *errorLogger) Output(calldepth int, s string) error {
	if m.shouldError {
		return errors.New("mock output error")
	}
	return nil
}

// fatalCalled tracks if logFatal was called
var fatalCalled bool

// mockLogFatal replaces logFatal for testing
func mockLogFatal(v ...interface{}) {
	fatalCalled = true
	// Don't actually call log.Fatal as it would exit the test
}

func TestLoggingErrors(t *testing.T) {
	// Save original logFatal and restore it after test
	originalLogFatal := logFatal
	defer func() {
		logFatal = originalLogFatal
	}()

	// Replace logFatal with our mock
	logFatal = mockLogFatal

	tests := []struct {
		name   string
		method string
	}{
		{"Println", "Println"},
		{"Printf", "Printf"},
		{"Print", "Print"},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_Error", func(t *testing.T) {
			fatalCalled = false
			logger := internalLog{logger: &errorLogger{shouldError: true}}

			switch tt.method {
			case "Println":
				logger.Println("test")
			case "Printf":
				logger.Printf("test %s", "message")
			case "Print":
				logger.Print("test")
			}

			assert.True(t, fatalCalled, "logFatal should have been called when Output returns error")
		})
	}
}

func TestInternalLogInterface(t *testing.T) {
	// Test that internalLog implements ilogger interface
	var _ ilogger = internalLog{}

	// Test with real logger for successful Output calls
	buf := bytes.NewBufferString("")
	logger := internalLog{logger: log.New(buf, "test: ", 0)}

	// Test all methods exist and work
	logger.Print("print test")
	assert.Contains(t, buf.String(), "print test")

	buf.Reset()
	logger.Printf("printf %s", "test")
	assert.Contains(t, buf.String(), "printf test")

	buf.Reset()
	logger.Println("println test")
	assert.Contains(t, buf.String(), "println test")

	buf.Reset()
	err := logger.Output(2, "output test")
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "output test")
}
