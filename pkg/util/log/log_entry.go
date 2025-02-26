// Copyright 2019 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package log

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/cockroachdb/cockroach/pkg/build"
	"github.com/cockroachdb/cockroach/pkg/util/caller"
	"github.com/cockroachdb/cockroach/pkg/util/log/channel"
	"github.com/cockroachdb/cockroach/pkg/util/log/eventpb"
	"github.com/cockroachdb/cockroach/pkg/util/log/logpb"
	"github.com/cockroachdb/cockroach/pkg/util/log/severity"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/cockroachdb/logtags"
	"github.com/cockroachdb/redact"
	"github.com/petermattis/goid"
)

// logEntry represents a logging event flowing through this package.
//
// It is different from logpb.Entry in that it is able to preserve
// more information about the structure of the source event, so that
// more details about this structure can be preserved by output
// formatters. logpb.Entry, in comparison, was tailored specifically
// to the legacy crdb-v1 formatter, and is a lossy representation.
type logEntry struct {
	// The entry timestamp.
	ts int64
	// The severity of the event.
	sev Severity
	// The channel on which the entry was sent.
	ch Channel
	// The goroutine where the event was generated.
	gid int64
	// The file/line where the event was generated.
	file string
	line int

	// The entry counter. Populated by outputLogEntry().
	counter uint64

	// The logging tags.
	tags *logtags.Buffer

	// The stack trace(s), when processing e.g. a fatal event.
	stacks []byte

	// Whether the entry is structured or not.
	structured bool

	// The entry payload.
	payload entryPayload
}

type entryPayload struct {
	// Whether the payload is redactable or not.
	redactable bool

	// The actual payload string.
	// For structured entries, this is the JSON
	// representation of the payload fields, without the
	// outer '{}'.
	// For unstructured entries, this is the (flat) message.
	//
	// If redactable is true, message is a RedactableString
	// in disguise. If it is false, message is a flat string with
	// no guarantees about content.
	message string
}

func makeRedactablePayload(m redact.RedactableString) entryPayload {
	return entryPayload{redactable: true, message: string(m)}
}

func makeUnsafePayload(m string) entryPayload {
	return entryPayload{redactable: false, message: m}
}

// makeEntry creates a logEntry.
func makeEntry(ctx context.Context, s Severity, c Channel, depth int) (res logEntry) {
	res = logEntry{
		ts:   timeutil.Now().UnixNano(),
		sev:  s,
		ch:   c,
		gid:  goid.Get(),
		tags: logtags.FromContext(ctx),
	}

	// Populate file/lineno.
	res.file, res.line, _ = caller.Lookup(depth + 1)

	return res
}

// makeStructuredEntry creates a logEntry using a structured payload.
func makeStructuredEntry(
	ctx context.Context, s Severity, c Channel, depth int, payload eventpb.EventPayload,
) (res logEntry) {
	res = makeEntry(ctx, s, c, depth+1)

	res.structured = true
	_, b := payload.AppendJSONFields(false, nil)
	res.payload = makeRedactablePayload(b.ToString())
	return res
}

// makeUnstructuredEntry creates a logEntry using an unstructured message.
func makeUnstructuredEntry(
	ctx context.Context,
	s Severity,
	c Channel,
	depth int,
	redactable bool,
	format string,
	args ...interface{},
) (res logEntry) {
	res = makeEntry(ctx, s, c, depth+1)

	res.structured = false

	if redactable {
		var buf redact.StringBuilder
		if len(args) == 0 {
			// TODO(knz): Remove this legacy case.
			buf.Print(redact.Safe(format))
		} else if len(format) == 0 {
			buf.Print(args...)
		} else {
			buf.Printf(format, args...)
		}
		res.payload = makeRedactablePayload(buf.RedactableString())
	} else {
		var buf strings.Builder
		formatArgs(&buf, format, args...)
		res.payload = makeUnsafePayload(buf.String())
	}

	return res
}

var configTagsBuffer = logtags.SingleTagBuffer("config", nil)

// makeStartLine creates a formatted log entry suitable for the start
// of a logging output using the canonical logging format.
func makeStartLine(formatter logFormatter, format string, args ...interface{}) *buffer {
	entry := makeUnstructuredEntry(
		context.Background(),
		severity.INFO,
		channel.DEV, /* DEV ensures the channel number is omitted in headers. */
		2,           /* depth */
		true,        /* redactable */
		format,
		args...)
	entry.tags = configTagsBuffer
	return formatter.formatEntry(entry)
}

// getStartLines retrieves the log entries for the start
// of a new log file output.
func (l *sinkInfo) getStartLines(now time.Time) []*buffer {
	f := l.formatter
	messages := make([]*buffer, 0, 6)
	messages = append(messages,
		makeStartLine(f, "file created at: %s", Safe(now.Format("2006/01/02 15:04:05"))),
		makeStartLine(f, "running on machine: %s", host),
		makeStartLine(f, "binary: %s", Safe(build.GetInfo().Short())),
		makeStartLine(f, "arguments: %s", os.Args),
	)

	logging.mu.Lock()
	if logging.mu.clusterID != "" {
		messages = append(messages, makeStartLine(f, "clusterID: %s", logging.mu.clusterID))
	}
	if logging.mu.nodeID != 0 {
		messages = append(messages, makeStartLine(f, "nodeID: n%d", logging.mu.nodeID))
	}
	if logging.mu.tenantID != "" {
		messages = append(messages, makeStartLine(f, "tenantID: %s", logging.mu.tenantID))
	}
	if logging.mu.sqlInstanceID != 0 {
		messages = append(messages, makeStartLine(f, "instanceID: %d", logging.mu.sqlInstanceID))
	}
	logging.mu.Unlock()

	// Including a non-ascii character in the first 1024 bytes of the log helps
	// viewers that attempt to guess the character encoding.
	messages = append(messages,
		makeStartLine(f, "line format: [IWEF]yymmdd hh:mm:ss.uuuuuu goid file:line msg utf8=\u2713"))
	return messages
}

// convertToLegacy turns the entry into a logpb.Entry.
func (e logEntry) convertToLegacy() (res logpb.Entry) {
	res = logpb.Entry{
		Severity:   e.sev,
		Channel:    e.ch,
		Time:       e.ts,
		File:       e.file,
		Line:       int64(e.line),
		Goroutine:  e.gid,
		Counter:    e.counter,
		Redactable: e.payload.redactable,
		Message:    e.payload.message,
	}

	if e.tags != nil {
		if e.payload.redactable {
			res.Tags = string(renderTagsAsRedactable(e.tags))
		} else {
			var buf strings.Builder
			e.tags.FormatToString(&buf)
			res.Tags = buf.String()
		}
	}

	if e.structured {
		// At this point, the message only contains the JSON fields of the
		// payload. Add the decoration suitable for our legacy file
		// format.
		res.Message = "Structured entry: {" + res.Message + "}"
	}

	if e.stacks != nil {
		res.Message += "\n" + string(e.stacks)
	}

	return res
}

// MakeLegacyEntry creates an logpb.Entry.
func MakeLegacyEntry(
	ctx context.Context,
	s Severity,
	c Channel,
	depth int,
	redactable bool,
	format string,
	args ...interface{},
) (res logpb.Entry) {
	return makeUnstructuredEntry(ctx, s, c, depth+1, redactable, format, args...).convertToLegacy()
}
