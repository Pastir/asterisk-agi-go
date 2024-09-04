package asterisk_agi

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Response represents a response to an AGI
// request.
type Response struct {
	Error        error  // Error received, if any
	Status       int    // HTTP-style status code received
	Result       int    // Result is the numerical return (if parseable)
	ResultString string // Result value as a string
	Value        string // Value is the (optional) string value returned
}

// Session represents AGI session
type Session struct {
	// Variable AGI
	Variables map[string]string
	reader    *bufio.Reader
	writer    *bufio.Writer
}

// New creates a new AGI session
func New() *Session {
	return &Session{
		Variables: make(map[string]string),
		reader:    bufio.NewReader(os.Stdin),
		writer:    bufio.NewWriter(os.Stdout),
	}
}

// SendCommand sends an AGI command
func (s *Session) SendCommand(command string) (string, error) {
	// Writing a command in Asterisk
	_, err := s.writer.WriteString(command + "\n")
	if err != nil {
		return "", err
	}

	err = s.writer.Flush()
	if err != nil {
		return "", err
	}

	// Reading the response from Asterisk
	response, err := s.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(response), nil
}

// Close session close
func (s *Session) Close() error {
	return nil
}

// Answer отвечает на звонок
func (s *Session) Answer() error {
	_, err := s.SendCommand("ANSWER")
	return err
}

// StreamFile воспроизводит файл
func (s *Session) StreamFile(fileName, escapeDigits string) error {
	command := fmt.Sprintf("STREAM FILE %s %s", fileName, escapeDigits)
	_, err := s.SendCommand(command)
	return err
}

// RecordFile records the caller's voice to a file
func (s *Session) RecordFile(fileName string, format string, timeout int, silenceTimeout int) error {
	// AGI command to write file with silence timeout
	command := fmt.Sprintf("RECORD FILE %s %s %d %d", fileName, format, timeout, silenceTimeout)
	response, err := s.SendCommand(command)
	if err != nil {
		return err
	}

	// Проверяем, успешно ли выполнена команда
	if strings.Contains(response, "200 result=0") {
		return nil
	} else {
		return fmt.Errorf("Ошибка записи файла: %s", response)
	}
}

// Hangup завершает звонок
func (s *Session) Hangup() error {
	_, err := s.SendCommand("HANGUP")
	return err
}

// GetVariables читает переменные AGI
func (s *Session) GetVariables() (map[string]string, error) {
	variables := make(map[string]string)
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil || line == "\n" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			variables[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return variables, nil
}
