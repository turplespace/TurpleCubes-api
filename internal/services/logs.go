package services

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	DEBUG LogLevel = "DEBUG"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     LogLevel  `json:"level"`
	Message   string    `json:"message"`
}

// logWriter implements io.Writer to capture standard log output
type logWriter struct {
	service *LogService
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	msg := string(p)
	// Only broadcast the message, don't write to stdout here
	w.service.broadcast(msg)
	return len(p), nil
}

type LogService struct {
	subscribers map[chan string]bool
	mu          sync.RWMutex
	logger      *log.Logger
}

var (
	logService *LogService
	once       sync.Once
)

func GetLogService() *LogService {
	once.Do(func() {
		ls := &LogService{
			subscribers: make(map[chan string]bool),
		}

		// Create a custom writer that only broadcasts
		writer := &logWriter{service: ls}

		// Create a multi-writer that writes to stdout and our broadcast writer
		multiWriter := io.MultiWriter(os.Stdout, writer)

		// Set the default Go logger to use our multi-writer
		log.SetOutput(multiWriter)

		// Initialize the service's logger with the same multi-writer
		ls.logger = log.New(multiWriter, "", log.LstdFlags)

		logService = ls
	})
	return logService
}

func (s *LogService) Subscribe() chan string {
	s.mu.Lock()
	defer s.mu.Unlock()
	ch := make(chan string, 100)
	s.subscribers[ch] = true
	return ch
}

func (s *LogService) Unsubscribe(ch chan string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.subscribers, ch)
	close(ch)
}

func (s *LogService) formatLogEntry(level LogLevel, message string) string {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}
	return fmt.Sprintf("[%s] %s: %s",
		entry.Timestamp.Format("2006/01/02 15:04:05"),
		entry.Level,
		entry.Message,
	)
}

func (s *LogService) broadcast(message string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Trim any trailing newlines for consistency
	message = strings.TrimSuffix(message, "\n")

	for ch := range s.subscribers {
		select {
		case ch <- message:
		default:
			// If channel is full, skip this message for this subscriber
		}
	}
}

func (s *LogService) Info(message string) {
	formattedMsg := s.formatLogEntry(INFO, message)
	s.logger.Println(formattedMsg)
}

func (s *LogService) Error(message string) {
	formattedMsg := s.formatLogEntry(ERROR, message)
	s.logger.Println(formattedMsg)
}

func (s *LogService) Debug(message string) {
	formattedMsg := s.formatLogEntry(DEBUG, message)
	s.logger.Println(formattedMsg)
}
