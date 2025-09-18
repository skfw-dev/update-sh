package runner

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Encoding int

const (
	UTF8    Encoding = iota
	UTF16LE          // UTF-16 Little Endian for Windows
	UTF16BE          // UTF-16 Big Endian (not commonly used)
)

func (e Encoding) String() string {
	switch e {
	case UTF8:
		return "UTF-8"
	case UTF16LE:
		return "UTF-16LE"
	case UTF16BE:
		return "UTF-16BE"
	default:
		return "Unknown Encoding"
	}
}

type CommandOptions struct {
	Description string
	DryRun      bool
	Name        string
	User        string
	Env         []string
	Args        []string
	Encoding    Encoding
}

func NewCommandOptions(description string, dryRun bool, name string, env []string, args ...string) *CommandOptions {
	return &CommandOptions{
		Description: description,
		DryRun:      dryRun,
		Name:        name,
		Env:         env,
		Args:        args,
		Encoding:    defaultEncoding(),
	}
}

func defaultEncoding() Encoding {
	// Default encoding can be set based on the platform or user preference.
	// For Linux/Unix-like systems, UTF-8 is standard.
	// For Windows, UTF-16 Little Endian is often used for command output.
	// if config.IsWindows() {
	// 	return UTF16LE
	// }
	return UTF8
}

func makeDecoder(encoding Encoding) (transform.Transformer, error) {
	// Use a transformer for encoding if specified
	switch encoding {
	case UTF8:
		return unicode.UTF8.NewDecoder(), nil
	case UTF16LE:
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), nil
	case UTF16BE:
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(), nil
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", encoding.String())
	}
}

func RunCommandWithOptions(opts *CommandOptions) error {
	if opts.DryRun {
		log.Info().Msgf("Dry Run: Would execute '%s': %s %v", opts.Description, opts.Name, opts.Args)
		return nil
	}

	log.Info().Msgf("%s...", opts.Description)

	cmd := exec.Command(opts.Name, opts.Args...)
	if len(opts.Env) > 0 {
		cmd.Env = append(cmd.Env, opts.Env...)
	}

	// Use a transformer for encoding if specified
	decoder, err := makeDecoder(opts.Encoding)
	if err != nil {
		return fmt.Errorf("failed to create decoder for encoding %s: %w", opts.Encoding.String(), err)
	}

	// Custom zerolog console writer
	// cmd.Stdout = zerolog.ConsoleWriter{Out: log.Logger.Output(os.Stdout), TimeFormat: zerolog.TimeFormatUnix}
	// cmd.Stderr = zerolog.ConsoleWriter{Out: log.Logger.Output(os.Stderr), TimeFormat: zerolog.TimeFormatUnix}
	return streamAndWait(cmd, decoder, opts.Description, opts.User)
}

// RunCommand executes a command and streams its output in real-time
func RunCommand(description string, dryRun bool, name string, env []string, arg ...string) error {
	opts := NewCommandOptions(description, dryRun, name, env, arg...)
	return RunCommandWithOptions(opts)
}

// streamAndWait runs the command, streams live output, and logs exit status
func streamAndWait(cmd *exec.Cmd, transformer transform.Transformer, description string, userTag string) error {
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe error: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe error: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Use tagged prefix for user-based logs if applicable
	tag := func() string {
		if userTag != "" {
			return fmt.Sprintf("User(tag=%s)", strconv.Quote(userTag))
		}
		return ""
	}

	go streamOutput(stdoutPipe, transformer, log.Info, tag)
	go streamOutput(stderrPipe, transformer, log.Warn, tag)

	if err := cmd.Wait(); err != nil {
		log.Error().Err(err).Msgf("Failed to %s", description)
		return err
	}

	log.Debug().Msgf("%s complete.", description)
	return nil
}

// streamOutput pipes output line-by-line to specified logger level
func streamOutput(r io.Reader, transformer transform.Transformer, level func() *zerolog.Event, tagFunc func() string) {
	// Use a transformer if specified
	if transformer != nil {
		r = transform.NewReader(r, transformer)
	}

	// Use a bufio.Reader for efficient, byte-level reading
	reader := bufio.NewReader(r)

	// Use a bytes.Buffer to collect lines
	var buf bytes.Buffer
	var line []byte
	var isPrefix bool
	var err error

	for {
		line, isPrefix, err = reader.ReadLine()
		buf.Write(line)

		// A line has been fully read
		if !isPrefix {
			line := strings.TrimSpace(buf.String())
			buf.Reset()

			if line == "" {
				continue // Skip empty content
			}

			// Skip specific warning message
			warnMsg := "WARNING: apt does not have a stable CLI interface. Use with caution in scripts."
			if strings.EqualFold(line, warnMsg) {
				continue
			}

			// Render output
			renderOutput(line, level, tagFunc)
		}

		// Handle end of stream or other errors
		if err != nil {
			if !errors.Is(err, io.EOF) && !errors.Is(err, os.ErrClosed) {
				log.Error().Err(err).Msg("error streaming command output")
			}
			break
		}
	}
}

func renderOutput(content string, level func() *zerolog.Event, tagFunc func() string) {
	var prefix string
	if tagFunc != nil {
		prefix = tagFunc()
	}

	if prefix != "" {
		if strings.Contains(content, "\r") {
			values := strings.SplitSeq(content, "\r")

			// Handle re-rendering
			for value := range values {

				// Skip empty string
				if value = strings.TrimSpace(value); value == "" {
					continue
				}

				// FIXME: "- \r 	 \\ \r 	 | \r 	 / \r 	 - \r 	 \\ \r 	 | \r 	 / \r 	 - \r"

				// Pass through with prefix
				level().Msgf("%s %s", prefix, value)
			}
			return
		}

		// Pass through with prefix
		level().Msgf("%s %s", prefix, content)
		return
	}

	// Pass through
	level().Msg(content)
}

// CommandExists checks if a command exists in the PATH.
// This function is common to both Linux and Windows.
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// RunCommandAndCaptureOutput executes a command and returns its standard output.
// This is useful for commands whose output needs to be parsed by the calling function.
func RunCommandAndCaptureOutput(description string, name string, env []string, arg ...string) (string, error) {
	log.Info().Msgf("Capturing output of '%s'...", description)
	cmd := exec.Command(name, arg...)
	if len(env) > 0 {
		cmd.Env = append(cmd.Env, env...)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to run command: %s", stderr.String())
		return "", err
	}
	return stdout.String(), nil
}
