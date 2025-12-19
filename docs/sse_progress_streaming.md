# SSE Progress Streaming

This document describes how to use the Server-Sent Events (SSE) progress streaming feature for real-time import monitoring.

## Overview

The import service supports real-time progress updates and log streaming via **Server-Sent Events (SSE)**. This allows clients to monitor long-running import operations with live progress bars and log messages.

## Architecture

1. Client generates a unique session ID (UUID recommended)
2. Client opens SSE connection to `/stream/:sessionId`
3. Client submits import request to `/file` with the same `session_id`
4. Server streams progress and log events to the connected client
5. Connection automatically closes when import completes

## API Endpoints

### GET `/stream/:sessionId`

Opens an SSE connection for receiving real-time updates.

**Headers:**
- `Content-Type: text/event-stream`
- `Cache-Control: no-cache`
- `Connection: keep-alive`

**Response:** Stream of SSE events

### POST `/file`

Submits a file import request. Include `session_id` for progress streaming.

**Request Body:**
```json
{
  "url": "https://example.com/file.dsv",
  "file_extension": "DSV",
  "file_type": "DEFINITION",
  "meeting": "meeting-123",
  "exclude_events": [],
  "include_events": [],
  "session_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

## Event Types

### Progress Event
```json
{
  "type": "progress",
  "progress": 45.5,
  "message": "Processing records..."
}
```

- `progress`: Float between 0 and 100
- `message`: Optional descriptive message

### Log Event
```json
{
  "type": "log",
  "message": "Imported 150 records successfully",
  "level": "info"
}
```

- `level`: One of `"info"`, `"error"`, `"warning"`, `"success"`
- `message`: Log message text

## Client Implementation Examples

### JavaScript (Browser)

```javascript
// Generate unique session ID
const sessionId = crypto.randomUUID();

// Open SSE connection
const eventSource = new EventSource(`http://localhost:8080/stream/${sessionId}`);

eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data);
  
  if (data.type === 'progress') {
    // Update progress bar
    updateProgressBar(data.progress);
    console.log(`Progress: ${data.progress}% - ${data.message}`);
  } else if (data.type === 'log') {
    // Display log message
    addLogMessage(data.message, data.level);
  }
};

eventSource.onerror = (error) => {
  console.error('SSE Error:', error);
  eventSource.close();
};

// Submit import request with same session ID
fetch('http://localhost:8080/file', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    url: 'https://example.com/file.dsv',
    file_extension: 'DSV',
    file_type: 'DEFINITION',
    meeting: 'meeting-123',
    session_id: sessionId
  })
})
.then(response => {
  if (response.ok) {
    console.log('Import started successfully');
  }
});

// Helper functions
function updateProgressBar(progress) {
  document.getElementById('progress-bar').style.width = `${progress}%`;
  document.getElementById('progress-text').textContent = `${progress.toFixed(1)}%`;
}

function addLogMessage(message, level) {
  const logEntry = document.createElement('div');
  logEntry.className = `log-${level}`;
  logEntry.textContent = `[${new Date().toLocaleTimeString()}] ${message}`;
  document.getElementById('log-container').appendChild(logEntry);
}
```

### cURL (Testing)

```bash
# Terminal 1: Open SSE stream
SESSION_ID="test-session-123"
curl -N "http://localhost:8080/stream/${SESSION_ID}"

# Terminal 2: Submit import
curl -X POST "http://localhost:8080/file" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com/file.dsv",
    "file_extension": "DSV",
    "file_type": "DEFINITION",
    "meeting": "meeting-123",
    "session_id": "test-session-123"
  }'
```

### Go Client Example

```go
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ProgressEvent struct {
	Type     string  `json:"type"`
	Progress float64 `json:"progress"`
	Message  string  `json:"message"`
}

type LogEvent struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Level   string `json:"level"`
}

func main() {
	sessionID := "test-session-123"
	
	// Start SSE stream in goroutine
	go func() {
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/stream/%s", sessionID))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				
				var event map[string]interface{}
				json.Unmarshal([]byte(data), &event)
				
				if event["type"] == "progress" {
					fmt.Printf("Progress: %.1f%% - %s\n", 
						event["progress"], event["message"])
				} else if event["type"] == "log" {
					fmt.Printf("[%s] %s\n", 
						event["level"], event["message"])
				}
			}
		}
	}()
	
	// Submit import request
	reqBody, _ := json.Marshal(map[string]interface{}{
		"url":            "https://example.com/file.dsv",
		"file_extension": "DSV",
		"file_type":      "DEFINITION",
		"meeting":        "meeting-123",
		"session_id":     sessionID,
	})
	
	resp, err := http.Post("http://localhost:8080/file", 
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	
	fmt.Println("Import request submitted")
	
	// Keep main alive to receive events
	select {}
}
```

## Adding Progress Updates to Import Functions

### Example: Adding progress to DSV import

In your importer functions, you can add progress updates like this:

```go
func ImportDsvDefinitionFile(url string, meeting string, exclude []int, include []int) (*ImportFileStats, error) {
	// Get session ID from context or pass it as parameter
	sessionID := "your-session-id"
	
	service.SendLog(sessionID, "Starting DSV file download", "info")
	service.SendProgress(sessionID, 10, "Downloading file")
	
	// Download file
	file, err := downloadFile(url)
	if err != nil {
		return nil, err
	}
	
	service.SendProgress(sessionID, 30, "Parsing DSV data")
	
	// Parse file
	records, err := parseFile(file)
	if err != nil {
		return nil, err
	}
	
	service.SendLog(sessionID, fmt.Sprintf("Found %d records", len(records)), "info")
	service.SendProgress(sessionID, 50, "Processing records")
	
	// Process records
	for i, record := range records {
		// Process each record
		processRecord(record)
		
		// Update progress periodically
		if i%10 == 0 {
			progress := 50 + (float64(i)/float64(len(records)))*40
			service.SendProgress(sessionID, progress, fmt.Sprintf("Processed %d/%d records", i, len(records)))
		}
	}
	
	service.SendProgress(sessionID, 90, "Finalizing import")
	
	// Finalize
	stats := finalizeImport()
	
	service.SendProgress(sessionID, 100, "Import completed")
	service.SendLog(sessionID, "Import finished successfully", "success")
	
	return stats, nil
}
```

## Helper Methods Available

### `service.SendProgress(sessionID string, progress float64, message string)`
Send a progress update (0-100).

### `service.SendLog(sessionID string, message string, level string)`
Send a log message. Levels: `"info"`, `"error"`, `"warning"`, `"success"`.

### `service.SendError(sessionID string, err error)`
Send an error log message.

### `service.SendComplete(sessionID string)`
Send completion event and close the session.

## Best Practices

1. **Generate Unique Session IDs**: Use UUIDs to avoid collisions
2. **Open SSE Before Import**: Establish the SSE connection before submitting the import request
3. **Handle Disconnections**: Implement reconnection logic on the client side
4. **Progress Granularity**: Update progress every 5-10% to avoid overwhelming the client
5. **Meaningful Messages**: Include descriptive messages with progress updates
6. **Error Handling**: Always send error logs before closing the connection
7. **Cleanup**: Sessions are automatically cleaned up after completion

## Troubleshooting

### No events received
- Ensure SSE connection is established before submitting import request
- Verify session IDs match between `/stream/:sessionId` and request body
- Check for CORS issues in browser

### Connection closes immediately
- Session might have been used before (reuse protection)
- Backend might have crashed during import
- Network timeout

### Missing events
- Channel buffer is full (100 messages) - increase update frequency
- Connection interrupted - implement reconnection
