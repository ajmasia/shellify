package handlers

import "net/http"

// APIDocs returns an HTML page with API documentation.
func APIDocs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(apiDocsHTML))
	}
}

const apiDocsHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Shellify API Documentation</title>
  <style>
    :root {
      --bg: #1e1e2e;
      --surface: #313244;
      --text: #cdd6f4;
      --subtext: #a6adc8;
      --blue: #89b4fa;
      --green: #a6e3a1;
      --yellow: #f9e2af;
      --red: #f38ba8;
      --mauve: #cba6f7;
      --peach: #fab387;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
      background: var(--bg);
      color: var(--text);
      line-height: 1.6;
      padding: 2rem;
      max-width: 1200px;
      margin: 0 auto;
    }
    h1 { color: var(--mauve); margin-bottom: 0.5rem; }
    h2 { color: var(--blue); margin: 2rem 0 1rem; border-bottom: 1px solid var(--surface); padding-bottom: 0.5rem; }
    h3 { color: var(--text); margin: 1.5rem 0 0.5rem; }
    .version { color: var(--subtext); margin-bottom: 2rem; }
    .endpoint {
      background: var(--surface);
      border-radius: 8px;
      padding: 1rem;
      margin: 1rem 0;
    }
    .method {
      display: inline-block;
      padding: 0.25rem 0.5rem;
      border-radius: 4px;
      font-weight: bold;
      font-size: 0.85rem;
      margin-right: 0.5rem;
    }
    .get { background: var(--green); color: var(--bg); }
    .post { background: var(--blue); color: var(--bg); }
    .put { background: var(--yellow); color: var(--bg); }
    .delete { background: var(--red); color: var(--bg); }
    .path {
      font-family: 'Monaco', 'Menlo', monospace;
      color: var(--peach);
    }
    .desc { color: var(--subtext); margin-top: 0.5rem; font-size: 0.9rem; }
    code {
      background: var(--bg);
      padding: 0.2rem 0.4rem;
      border-radius: 4px;
      font-family: 'Monaco', 'Menlo', monospace;
      font-size: 0.85rem;
    }
    pre {
      background: var(--bg);
      padding: 1rem;
      border-radius: 4px;
      overflow-x: auto;
      margin: 0.5rem 0;
    }
    .response-format {
      background: var(--surface);
      border-radius: 8px;
      padding: 1rem;
      margin: 1rem 0;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      margin: 1rem 0;
    }
    th, td {
      text-align: left;
      padding: 0.75rem;
      border-bottom: 1px solid var(--surface);
    }
    th { color: var(--blue); }
    footer {
      margin-top: 3rem;
      padding-top: 1rem;
      border-top: 1px solid var(--surface);
      color: var(--subtext);
      font-size: 0.85rem;
    }
  </style>
</head>
<body>
  <h1>Shellify API</h1>
  <p class="version">v0.6.0 - REST API for terminal multiplexer session management</p>

  <h2>Response Format</h2>
  <div class="response-format">
    <p>All responses follow this format:</p>
    <pre><code>{
  "success": true,
  "data": { ... },
  "error": ""  // Only on failure
}</code></pre>
  </div>

  <h2>Health</h2>
  <div class="endpoint">
    <span class="method get">GET</span>
    <span class="path">/api/health</span>
    <p class="desc">Check server health status</p>
  </div>

  <h2>Projects</h2>

  <div class="endpoint">
    <span class="method get">GET</span>
    <span class="path">/api/projects</span>
    <p class="desc">List all projects with session counts</p>
  </div>

  <div class="endpoint">
    <span class="method post">POST</span>
    <span class="path">/api/projects</span>
    <p class="desc">Create a new project</p>
    <pre><code>{ "name": "string", "description": "string" }</code></pre>
  </div>

  <div class="endpoint">
    <span class="method get">GET</span>
    <span class="path">/api/projects/{projectId}</span>
    <p class="desc">Get a single project by ID</p>
  </div>

  <div class="endpoint">
    <span class="method put">PUT</span>
    <span class="path">/api/projects/{projectId}</span>
    <p class="desc">Update a project</p>
    <pre><code>{ "name": "string?", "description": "string?" }</code></pre>
  </div>

  <div class="endpoint">
    <span class="method delete">DELETE</span>
    <span class="path">/api/projects/{projectId}</span>
    <p class="desc">Delete a project. Use <code>?force=true</code> to delete with sessions</p>
  </div>

  <div class="endpoint">
    <span class="method get">GET</span>
    <span class="path">/api/projects/{projectId}/sessions</span>
    <p class="desc">List all sessions in a project</p>
  </div>

  <div class="endpoint">
    <span class="method post">POST</span>
    <span class="path">/api/projects/{projectId}/backup</span>
    <p class="desc">Create a backup of project with all sessions</p>
  </div>

  <div class="endpoint">
    <span class="method post">POST</span>
    <span class="path">/api/projects/{projectId}/restore</span>
    <p class="desc">Restore sessions into a project</p>
    <pre><code>{ "sessions": [ { session data... } ] }</code></pre>
  </div>

  <h2>Sessions</h2>

  <div class="endpoint">
    <span class="method get">GET</span>
    <span class="path">/api/sessions</span>
    <p class="desc">List all sessions across all projects</p>
  </div>

  <div class="endpoint">
    <span class="method post">POST</span>
    <span class="path">/api/sessions</span>
    <p class="desc">Create a new session</p>
    <pre><code>{
  "projectId": "string",
  "session": {
    "name": "string",
    "description": "string?",
    "workingDirectory": "string?",
    "targetMultiplexer": "tmux" | "zellij",
    "environment": { "KEY": "value" },
    "preCommands": ["string"],
    "postCommands": ["string"],
    "windows": [{
      "name": "string",
      "workingDirectory": "string?",
      "command": "string?",
      "panes": [{
        "name": "string?",
        "command": "string?",
        "workingDirectory": "string?",
        "size": number
      }]
    }]
  }
}</code></pre>
  </div>

  <div class="endpoint">
    <span class="method get">GET</span>
    <span class="path">/api/sessions/running</span>
    <p class="desc">List all currently running multiplexer sessions</p>
  </div>

  <div class="endpoint">
    <span class="method get">GET</span>
    <span class="path">/api/sessions/{sessionId}</span>
    <p class="desc">Get a single session by ID</p>
  </div>

  <div class="endpoint">
    <span class="method put">PUT</span>
    <span class="path">/api/sessions/{sessionId}</span>
    <p class="desc">Update a session</p>
    <pre><code>{
  "name": "string?",
  "description": "string?",
  "workingDirectory": "string?",
  "targetMultiplexer": "tmux" | "zellij"
}</code></pre>
  </div>

  <div class="endpoint">
    <span class="method delete">DELETE</span>
    <span class="path">/api/sessions/{sessionId}</span>
    <p class="desc">Delete a session</p>
  </div>

  <div class="endpoint">
    <span class="method post">POST</span>
    <span class="path">/api/sessions/{sessionId}/clone</span>
    <p class="desc">Clone a session</p>
    <pre><code>{ "name": "string?" }  // New name, defaults to "Copy"</code></pre>
  </div>

  <div class="endpoint">
    <span class="method post">POST</span>
    <span class="path">/api/sessions/{sessionId}/launch</span>
    <p class="desc">Launch a session in detached mode</p>
  </div>

  <div class="endpoint">
    <span class="method post">POST</span>
    <span class="path">/api/sessions/{sessionId}/attach</span>
    <p class="desc">Attach to a running session</p>
  </div>

  <div class="endpoint">
    <span class="method post">POST</span>
    <span class="path">/api/sessions/{sessionId}/stop</span>
    <p class="desc">Stop a running session</p>
  </div>

  <div class="endpoint">
    <span class="method get">GET</span>
    <span class="path">/api/sessions/{sessionId}/status</span>
    <p class="desc">Check if session is running</p>
  </div>

  <h2>Settings</h2>

  <div class="endpoint">
    <span class="method get">GET</span>
    <span class="path">/api/settings</span>
    <p class="desc">Get application settings</p>
  </div>

  <div class="endpoint">
    <span class="method put">PUT</span>
    <span class="path">/api/settings</span>
    <p class="desc">Update application settings</p>
  </div>

  <h2>Error Codes</h2>
  <table>
    <tr><th>Status</th><th>Description</th></tr>
    <tr><td><code>200</code></td><td>Success</td></tr>
    <tr><td><code>201</code></td><td>Created</td></tr>
    <tr><td><code>204</code></td><td>No Content</td></tr>
    <tr><td><code>400</code></td><td>Bad Request - Invalid input</td></tr>
    <tr><td><code>404</code></td><td>Not Found - Resource doesn't exist</td></tr>
    <tr><td><code>409</code></td><td>Conflict - Resource already exists</td></tr>
    <tr><td><code>500</code></td><td>Internal Server Error</td></tr>
  </table>

  <footer>
    <p>Shellify - Terminal multiplexer session manager for tmux and zellij</p>
  </footer>
</body>
</html>`
