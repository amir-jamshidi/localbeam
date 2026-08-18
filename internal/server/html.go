package server

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>LocalBeam</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Space+Mono:wght@400;700&family=DM+Sans:wght@300;400;500;600&display=swap" rel="stylesheet">
<style>
  :root {
    --bg: #0a0a0f;
    --surface: #13131a;
    --surface2: #1c1c27;
    --border: #2a2a3a;
    --accent: #6c63ff;
    --accent2: #ff6584;
    --accent3: #43e97b;
    --text: #e8e8f0;
    --text2: #8888a8;
    --text3: #5555778;
    --radius: 16px;
    --mono: 'Space Mono', monospace;
    --sans: 'DM Sans', sans-serif;
  }

  * { box-sizing: border-box; margin: 0; padding: 0; }

  body {
    background: var(--bg);
    color: var(--text);
    font-family: var(--sans);
    min-height: 100vh;
    overflow-x: hidden;
  }

  /* Animated background */
  body::before {
    content: '';
    position: fixed;
    top: -50%;
    left: -50%;
    width: 200%;
    height: 200%;
    background: radial-gradient(ellipse at 20% 20%, rgba(108,99,255,0.08) 0%, transparent 50%),
                radial-gradient(ellipse at 80% 80%, rgba(255,101,132,0.06) 0%, transparent 50%),
                radial-gradient(ellipse at 50% 50%, rgba(67,233,123,0.04) 0%, transparent 60%);
    pointer-events: none;
    z-index: 0;
    animation: bgShift 20s ease-in-out infinite alternate;
  }

  @keyframes bgShift {
    0% { transform: translate(0,0) rotate(0deg); }
    100% { transform: translate(2%, 2%) rotate(1deg); }
  }

  .app { position: relative; z-index: 1; min-height: 100vh; display: flex; flex-direction: column; }

  /* Header */
  header {
    padding: 20px 32px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid var(--border);
    backdrop-filter: blur(20px);
    background: rgba(10,10,15,0.8);
    position: sticky;
    top: 0;
    z-index: 100;
  }

  .logo {
    display: flex;
    align-items: center;
    gap: 12px;
    text-decoration: none;
    color: var(--text);
  }

  .logo-icon {
    width: 36px;
    height: 36px;
    background: linear-gradient(135deg, var(--accent), var(--accent2));
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
  }

  .logo-text {
    font-family: var(--mono);
    font-size: 18px;
    font-weight: 700;
    letter-spacing: -0.5px;
  }

  .logo-text span { color: var(--accent); }

  .header-badge {
    background: rgba(108,99,255,0.15);
    border: 1px solid rgba(108,99,255,0.3);
    color: var(--accent);
    padding: 4px 12px;
    border-radius: 20px;
    font-size: 12px;
    font-family: var(--mono);
  }

  /* Main */
  main {
    flex: 1;
    display: flex;
    align-items: flex-start;
    justify-content: center;
    padding: 48px 24px;
  }

  .container { width: 100%; max-width: 640px; }

  /* Hero */
  .hero {
    text-align: center;
    margin-bottom: 48px;
  }

  .hero h1 {
    font-size: clamp(32px, 6vw, 52px);
    font-weight: 600;
    letter-spacing: -1.5px;
    line-height: 1.1;
    margin-bottom: 16px;
  }

  .hero h1 .grad {
    background: linear-gradient(135deg, var(--accent), var(--accent2));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .hero p {
    color: var(--text2);
    font-size: 16px;
    line-height: 1.6;
    max-width: 440px;
    margin: 0 auto;
  }

  /* Tab navigation */
  .tabs {
    display: flex;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: calc(var(--radius) + 4px);
    padding: 4px;
    margin-bottom: 32px;
    gap: 4px;
  }

  .tab {
    flex: 1;
    padding: 12px;
    border: none;
    background: transparent;
    color: var(--text2);
    font-family: var(--sans);
    font-size: 14px;
    font-weight: 500;
    border-radius: var(--radius);
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
  }

  .tab:hover { color: var(--text); background: var(--surface2); }

  .tab.active {
    background: var(--accent);
    color: white;
    box-shadow: 0 4px 20px rgba(108,99,255,0.4);
  }

  /* Card */
  .card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: calc(var(--radius) + 4px);
    overflow: hidden;
    animation: fadeUp 0.3s ease;
  }

  @keyframes fadeUp {
    from { opacity: 0; transform: translateY(16px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .card-header {
    padding: 24px 28px 0;
  }

  .card-title {
    font-size: 20px;
    font-weight: 600;
    margin-bottom: 6px;
  }

  .card-sub {
    font-size: 13px;
    color: var(--text2);
  }

  .card-body { padding: 24px 28px; }

  /* Type selector */
  .type-selector {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
    margin-bottom: 24px;
  }

  .type-btn {
    padding: 16px;
    border: 1.5px solid var(--border);
    background: var(--surface2);
    border-radius: var(--radius);
    cursor: pointer;
    color: var(--text2);
    font-family: var(--sans);
    transition: all 0.2s;
    text-align: center;
  }

  .type-btn:hover { border-color: var(--accent); color: var(--text); }

  .type-btn.selected {
    border-color: var(--accent);
    background: rgba(108,99,255,0.12);
    color: var(--text);
  }

  .type-btn .type-icon { font-size: 28px; margin-bottom: 8px; display: block; }
  .type-btn .type-label { font-size: 14px; font-weight: 500; display: block; }
  .type-btn .type-desc { font-size: 11px; color: var(--text2); margin-top: 4px; display: block; }

  /* Drop zone */
  .dropzone {
    border: 2px dashed var(--border);
    border-radius: var(--radius);
    padding: 40px 24px;
    text-align: center;
    cursor: pointer;
    transition: all 0.2s;
    position: relative;
    background: var(--surface2);
    margin-bottom: 20px;
  }

  .dropzone:hover, .dropzone.dragover {
    border-color: var(--accent);
    background: rgba(108,99,255,0.08);
  }

  .dropzone input { position: absolute; inset: 0; opacity: 0; cursor: pointer; width: 100%; }

  .dropzone-icon { font-size: 40px; margin-bottom: 12px; }

  .dropzone-text { font-size: 15px; font-weight: 500; margin-bottom: 6px; }

  .dropzone-hint { font-size: 12px; color: var(--text2); }

  /* File list */
  .file-list { list-style: none; display: flex; flex-direction: column; gap: 8px; margin-bottom: 20px; }

  .file-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 14px;
    background: var(--surface2);
    border: 1px solid var(--border);
    border-radius: 12px;
    animation: fadeUp 0.2s ease;
  }

  .file-icon { font-size: 24px; flex-shrink: 0; }

  .file-info { flex: 1; min-width: 0; }

  .file-name {
    font-size: 14px;
    font-weight: 500;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    margin-bottom: 2px;
  }

  .file-size { font-size: 11px; color: var(--text2); font-family: var(--mono); }

  .file-remove {
    background: none;
    border: none;
    color: var(--text2);
    cursor: pointer;
    padding: 4px;
    border-radius: 6px;
    font-size: 16px;
    transition: all 0.15s;
    flex-shrink: 0;
  }

  .file-remove:hover { background: rgba(255,101,132,0.2); color: var(--accent2); }

  /* Textarea */
  textarea {
    width: 100%;
    background: var(--surface2);
    border: 1.5px solid var(--border);
    border-radius: var(--radius);
    padding: 16px;
    color: var(--text);
    font-family: var(--sans);
    font-size: 14px;
    line-height: 1.6;
    resize: vertical;
    min-height: 140px;
    outline: none;
    transition: border-color 0.2s;
    margin-bottom: 20px;
  }

  textarea:focus { border-color: var(--accent); }
  textarea::placeholder { color: var(--text2); }

  /* Buttons */
  .btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 14px 24px;
    border: none;
    border-radius: var(--radius);
    font-family: var(--sans);
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    width: 100%;
  }

  .btn-primary {
    background: var(--accent);
    color: white;
    box-shadow: 0 4px 24px rgba(108,99,255,0.35);
  }

  .btn-primary:hover {
    transform: translateY(-1px);
    box-shadow: 0 8px 32px rgba(108,99,255,0.5);
  }

  .btn-primary:active { transform: translateY(0); }

  .btn-primary:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none;
    box-shadow: none;
  }

  .btn-secondary {
    background: var(--surface2);
    color: var(--text);
    border: 1px solid var(--border);
  }

  .btn-secondary:hover { background: var(--border); }

  .btn-success {
    background: var(--accent3);
    color: #0a1a0f;
    box-shadow: 0 4px 24px rgba(67,233,123,0.3);
  }

  .btn-danger {
    background: rgba(255,101,132,0.15);
    color: var(--accent2);
    border: 1px solid rgba(255,101,132,0.3);
  }

  /* PIN input */
  .pin-wrapper {
    display: flex;
    gap: 10px;
    margin-bottom: 24px;
  }

  .pin-digit {
    flex: 1;
    width:64px;
    height: 64px;
    background: var(--surface2);
    border: 1.5px solid var(--border);
    border-radius: 12px;
    text-align: center;
    font-family: var(--mono);
    font-size: 28px;
    font-weight: 700;
    color: var(--text);
    outline: none;
    transition: all 0.2s;
  }

  .pin-digit:focus { border-color: var(--accent); background: rgba(108,99,255,0.08); }

  /* Session display */
  .session-display {
    animation: fadeUp 0.4s ease;
  }

  .session-info-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
    margin-bottom: 24px;
  }

  .info-box {
    background: var(--surface2);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 16px;
  }

  .info-label {
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--text2);
    margin-bottom: 8px;
  }

  .info-value {
    font-family: var(--mono);
    font-size: 22px;
    font-weight: 700;
    color: var(--accent);
    letter-spacing: 4px;
  }

  .info-value.small {
    font-size: 13px;
    letter-spacing: 0;
    color: var(--text);
    word-break: break-all;
  }

  /* QR container */
  .qr-container {
    background: white;
    border-radius: var(--radius);
    padding: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 20px;
    aspect-ratio: 1;
    max-width: 260px;
    margin-left: auto;
    margin-right: auto;
  }

  .qr-container img { width: 100%; height: 100%; image-rendering: pixelated; }

  /* Progress */
  .progress-bar {
    height: 4px;
    background: var(--border);
    border-radius: 4px;
    overflow: hidden;
    margin-bottom: 12px;
  }

  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, var(--accent), var(--accent2));
    border-radius: 4px;
    transition: width 0.3s ease;
    width: 0%;
  }

  /* Alert */
  .alert {
    padding: 14px 16px;
    border-radius: 12px;
    font-size: 14px;
    margin-bottom: 16px;
    display: flex;
    align-items: flex-start;
    gap: 10px;
    animation: fadeUp 0.2s ease;
  }

  .alert-success { background: rgba(67,233,123,0.12); border: 1px solid rgba(67,233,123,0.3); color: var(--accent3); }
  .alert-error { background: rgba(255,101,132,0.12); border: 1px solid rgba(255,101,132,0.3); color: var(--accent2); }
  .alert-info { background: rgba(108,99,255,0.12); border: 1px solid rgba(108,99,255,0.3); color: var(--accent); }

  /* Download items */
  .download-list { display: flex; flex-direction: column; gap: 10px; margin-bottom: 20px; }

  .download-item {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 14px 16px;
    background: var(--surface2);
    border: 1px solid var(--border);
    border-radius: var(--radius);
  }

  .download-item .file-info { flex: 1; min-width: 0; }

  .btn-download {
    background: rgba(108,99,255,0.15);
    border: 1px solid rgba(108,99,255,0.3);
    color: var(--accent);
    border-radius: 10px;
    padding: 8px 16px;
    cursor: pointer;
    font-family: var(--sans);
    font-size: 13px;
    font-weight: 600;
    transition: all 0.2s;
    flex-shrink: 0;
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    gap: 6px;
  }

  .btn-download:hover {
    background: var(--accent);
    color: white;
    border-color: var(--accent);
  }

  /* Copy button */
  .copy-btn {
    background: var(--surface2);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 10px 16px;
    cursor: pointer;
    color: var(--text2);
    font-size: 13px;
    font-family: var(--sans);
    transition: all 0.2s;
    display: flex;
    align-items: center;
    gap: 6px;
    margin-top: 12px;
    width: 100%;
    justify-content: center;
  }

  .copy-btn:hover { background: var(--border); color: var(--text); }

  .copy-btn.copied { background: rgba(67,233,123,0.15); color: var(--accent3); border-color: rgba(67,233,123,0.3); }

  /* Text content display */
  .text-content {
    background: var(--surface2);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 20px;
    font-size: 14px;
    line-height: 1.7;
    white-space: pre-wrap;
    word-break: break-word;
    max-height: 320px;
    overflow-y: auto;
    margin-bottom: 16px;
  }

  /* Timer */
  .session-timer {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 12px;
    color: var(--text2);
    font-family: var(--mono);
    justify-content: center;
    margin-top: 16px;
  }

  .timer-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--accent3); animation: pulse 2s infinite; }

  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.3; }
  }

  .timer-dot.warning { background: #ffd36b; }
  .timer-dot.danger { background: var(--accent2); }

  /* Divider */
  .divider {
    display: flex;
    align-items: center;
    gap: 16px;
    margin: 20px 0;
    color: var(--text2);
    font-size: 12px;
  }

  .divider::before, .divider::after {
    content: '';
    flex: 1;
    height: 1px;
    background: var(--border);
  }

  /* Spinner */
  .spinner {
    width: 20px;
    height: 20px;
    border: 2.5px solid rgba(255,255,255,0.2);
    border-top-color: white;
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
  }

  @keyframes spin { to { transform: rotate(360deg); } }

  /* Footer */
  footer {
    text-align: center;
    padding: 24px;
    color: var(--text2);
    font-size: 12px;
    border-top: 1px solid var(--border);
  }

  footer a { color: var(--accent); text-decoration: none; }

  /* Responsive */
  @media (max-width: 480px) {
    header { padding: 16px 20px; }
    main { padding: 32px 16px; }
    .card-body { padding: 20px; }
    .card-header { padding: 20px 20px 0; }
    .session-info-grid { grid-template-columns: 1fr; }
    .pin-digit { height: 56px; font-size: 24px; }
  }

  /* Hidden */
  .hidden { display: none !important; }
</style>
</head>
<body>
<div class="app">
  <header>
    <a class="logo" href="/">
      <div class="logo-icon">⚡</div>
      <span class="logo-text">Local<span>Beam</span></span>
    </a>
    <div class="header-badge">LOCAL NETWORK</div>
  </header>

  <main>
    <div class="container">
      <!-- Home view -->
      <div id="view-home">
        <div class="hero">
          <h1><span class="grad">Powered By Frontend Team</span><br></h1>
          <p>Share files and text between any devices on your local network — no internet required.</p>
        </div>

        <div class="tabs">
          <button class="tab active" onclick="showSendView()" id="tab-send">
            <span>📤</span> Send
          </button>
          <button class="tab" onclick="showReceiveView()" id="tab-receive">
            <span>📥</span> Receive
          </button>
        </div>

        <!-- Send Panel -->
        <div id="panel-send">
          <div class="card">
            <div class="card-header">
              <div class="card-title">Send something</div>
              <div class="card-sub">Choose what you want to transfer</div>
            </div>
            <div class="card-body">
              <div class="type-selector">
                <button class="type-btn selected" onclick="selectType('file')" id="typebtn-file">
                  <span class="type-icon">📁</span>
                  <span class="type-label">Files</span>
                  <span class="type-desc">Any file up to 500MB</span>
                </button>
                <button class="type-btn" onclick="selectType('text')" id="typebtn-text">
                  <span class="type-icon">📝</span>
                  <span class="type-label">Text</span>
                  <span class="type-desc">Clipboard, notes, links</span>
                </button>
              </div>

              <!-- File upload area -->
              <div id="file-section">
                <div class="dropzone" id="dropzone">
                  <input type="file" id="file-input" multiple onchange="handleFileSelect(this.files)">
                  <div class="dropzone-icon">🗂️</div>
                  <div class="dropzone-text">Drop files here</div>
                  <div class="dropzone-hint">or click to browse · Multiple files supported</div>
                </div>
                <ul class="file-list hidden" id="file-list"></ul>
              </div>

              <!-- Text area -->
              <div id="text-section" class="hidden">
                <textarea id="text-input" placeholder="Paste or type text, links, code..." rows="6"></textarea>
              </div>

              <div id="send-alert"></div>

              <button class="btn btn-primary" id="btn-send" onclick="createSendSession()">
                <span>⚡</span> Create Beam
              </button>
            </div>
          </div>
        </div>

        <!-- Receive Panel (join by PIN) -->
        <div id="panel-receive" class="hidden">
          <div class="card">
            <div class="card-header">
              <div class="card-title">Join a Beam</div>
              <div class="card-sub">Enter the 6-digit PIN shown on the sender's device</div>
            </div>
            <div class="card-body">
              <div class="pin-wrapper" id="pin-inputs">
                <input class="pin-digit" maxlength="1" type="tel" oninput="pinInput(this,0)" onkeydown="pinKey(event,0)" autofocus>
                <input class="pin-digit" maxlength="1" type="tel" oninput="pinInput(this,1)" onkeydown="pinKey(event,1)">
                <input class="pin-digit" maxlength="1" type="tel" oninput="pinInput(this,2)" onkeydown="pinKey(event,2)">
                <input class="pin-digit" maxlength="1" type="tel" oninput="pinInput(this,3)" onkeydown="pinKey(event,3)">
                <input class="pin-digit" maxlength="1" type="tel" oninput="pinInput(this,4)" onkeydown="pinKey(event,4)">
                <input class="pin-digit" maxlength="1" type="tel" oninput="pinInput(this,5)" onkeydown="pinKey(event,5)">
              </div>

              <div id="join-alert"></div>

              <button class="btn btn-primary" id="btn-join" onclick="joinByPIN()">
                <span>🔗</span> Connect
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Session (sender) view -->
      <div id="view-session" class="hidden session-display">
        <div class="card">
          <div class="card-header">
            <div style="display:flex;align-items:center;justify-content:space-between">
              <div>
                <div class="card-title">✅ Beam created!</div>
                <div class="card-sub">Waiting for receiver to connect…</div>
              </div>
              <button class="btn-danger" style="width:auto;padding:8px 16px;border-radius:10px;font-size:13px;border:none;cursor:pointer" onclick="cancelSession()">Cancel</button>
            </div>
          </div>
          <div class="card-body">
            <div class="session-info-grid">
              <div class="info-box">
                <div class="info-label">PIN Code</div>
                <div class="info-value" id="session-pin">------</div>
              </div>
              <div class="info-box">
                <div class="info-label">Session ID</div>
                <div class="info-value small" id="session-id-display">--</div>
              </div>
            </div>

            <div class="qr-container" id="qr-container">
              <div class="spinner"></div>
            </div>
            <div style="text-align:center;font-size:12px;color:var(--text2);margin-bottom:20px">
              Scan QR code or share the PIN manually
            </div>

            <div id="upload-progress" class="hidden">
              <div class="progress-bar"><div class="progress-fill" id="progress-fill"></div></div>
              <div style="font-size:12px;color:var(--text2);text-align:center" id="progress-text">Uploading…</div>
            </div>

            <div id="session-alert"></div>

            <div class="session-timer">
              <div class="timer-dot" id="timer-dot"></div>
              <span id="timer-text">Session expires in --:--</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Receive view (after QR scan or PIN join) -->
      <div id="view-receive" class="hidden session-display">
        <div class="card">
          <div class="card-header">
            <div class="card-title" id="receive-title">📥 Ready to receive</div>
            <div class="card-sub" id="receive-sub">Content is ready for you</div>
          </div>
          <div class="card-body">
            <!-- Files -->
            <div id="receive-files" class="hidden">
              <div class="download-list" id="download-list"></div>
              <button class="btn btn-success" id="btn-download-all" onclick="downloadAll()">
                ⬇️ Download All Files
              </button>
            </div>

            <!-- Text -->
            <div id="receive-text" class="hidden">
              <div class="text-content" id="text-display"></div>
              <button class="copy-btn" id="copy-btn" onclick="copyText()">
                <span>📋</span> Copy to Clipboard
              </button>
            </div>

            <div id="receive-alert"></div>

            <button class="btn btn-secondary" style="margin-top:16px" onclick="goHome()">← New Transfer</button>
          </div>
        </div>
      </div>

    </div>
  </main>

  <footer>
    LocalBeam &mdash; Secure local file transfer &nbsp;·&nbsp;
    <a href="/api/info">API</a>
  </footer>
</div>

<script>
const API = '';
let currentSession = null;
let currentType = 'file';
let selectedFiles = [];
let timerInterval = null;

// --- Navigation ---
function showSendView() {
  document.getElementById('tab-send').classList.add('active');
  document.getElementById('tab-receive').classList.remove('active');
  document.getElementById('panel-send').classList.remove('hidden');
  document.getElementById('panel-receive').classList.add('hidden');
}

function showReceiveView() {
  document.getElementById('tab-receive').classList.add('active');
  document.getElementById('tab-send').classList.remove('active');
  document.getElementById('panel-receive').classList.remove('hidden');
  document.getElementById('panel-send').classList.add('hidden');
}

function goHome() {
  document.getElementById('view-home').classList.remove('hidden');
  document.getElementById('view-session').classList.add('hidden');
  document.getElementById('view-receive').classList.add('hidden');
  if (timerInterval) clearInterval(timerInterval);
  currentSession = null;
  selectedFiles = [];
  document.getElementById('file-list').innerHTML = '';
  document.getElementById('file-list').classList.add('hidden');
  document.getElementById('text-input').value = '';
}

// --- Type selector ---
function selectType(type) {
  currentType = type;
  document.getElementById('typebtn-file').classList.toggle('selected', type === 'file');
  document.getElementById('typebtn-text').classList.toggle('selected', type === 'text');
  document.getElementById('file-section').classList.toggle('hidden', type !== 'file');
  document.getElementById('text-section').classList.toggle('hidden', type !== 'text');
}

// --- File handling ---
function handleFileSelect(files) {
  for (const f of files) selectedFiles.push(f);
  renderFileList();
}

function renderFileList() {
  const list = document.getElementById('file-list');
  if (selectedFiles.length === 0) { list.classList.add('hidden'); return; }
  list.classList.remove('hidden');
  list.innerHTML = selectedFiles.map((f, i) => {
    return '<li class="file-item"><span class="file-icon">' + fileIcon(f.name) + '</span><div class="file-info"><div class="file-name">' + escHtml(f.name) + '</div><div class="file-size">' + formatSize(f.size) + '</div></div><button class="file-remove" onclick="removeFile(' + i + ')">✕</button></li>';
  }).join('');
}

function removeFile(i) {
  selectedFiles.splice(i, 1);
  renderFileList();
}

function fileIcon(name) {
  const ext = name.split('.').pop().toLowerCase();
  const map = { pdf:'📄', jpg:'🖼️', jpeg:'🖼️', png:'🖼️', gif:'🖼️', webp:'🖼️', mp4:'🎬', mov:'🎬', avi:'🎬', mp3:'🎵', wav:'🎵', zip:'📦', rar:'📦', tar:'📦', doc:'📝', docx:'📝', txt:'📄', js:'💻', ts:'💻', py:'💻', go:'💻', html:'🌐', css:'🎨' };
  return map[ext] || '📎';
}

function formatSize(bytes) {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024*1024) return (bytes/1024).toFixed(1) + ' KB';
  if (bytes < 1024*1024*1024) return (bytes/(1024*1024)).toFixed(1) + ' MB';
  return (bytes/(1024*1024*1024)).toFixed(2) + ' GB';
}

function escHtml(s) {
  return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
}

// Drag and drop
const dz = document.getElementById('dropzone');
dz.addEventListener('dragover', e => { e.preventDefault(); dz.classList.add('dragover'); });
dz.addEventListener('dragleave', () => dz.classList.remove('dragover'));
dz.addEventListener('drop', e => {
  e.preventDefault();
  dz.classList.remove('dragover');
  handleFileSelect(e.dataTransfer.files);
});

// --- PIN input ---
function pinInput(el, idx) {
  const digits = document.querySelectorAll('.pin-digit');
  if (el.value && idx < 5) digits[idx+1].focus();
  checkPINComplete();
}

function pinKey(e, idx) {
  const digits = document.querySelectorAll('.pin-digit');
  if (e.key === 'Backspace' && !digits[idx].value && idx > 0) {
    digits[idx-1].focus();
  }
  if (e.key === 'Enter') joinByPIN();
}

function getPIN() {
  return Array.from(document.querySelectorAll('.pin-digit')).map(d => d.value).join('');
}

function checkPINComplete() {
  const pin = getPIN();
  if (pin.length === 6) joinByPIN();
}

// --- Create session ---
async function createSendSession() {
  const btn = document.getElementById('btn-send');
  btn.disabled = true;
  btn.innerHTML = '<div class="spinner"></div> Creating…';

  try {
    const res = await fetch(API + '/api/session/create', {
      method: 'POST',
      headers: {'Content-Type':'application/json'},
      body: JSON.stringify({ type: currentType, direction: 'push' })
    });
    const data = await res.json();
    if (data.error) throw new Error(data.error);

    currentSession = data;

    document.getElementById('view-home').classList.add('hidden');
    document.getElementById('view-session').classList.remove('hidden');
    document.getElementById('session-pin').textContent = data.pin;
    document.getElementById('session-id-display').textContent = data.session_id.substring(0,12) + '…';

    startTimer(new Date(data.expires_at));
    loadQR(data.session_id);

    if (currentType === 'file') {
      await uploadFiles(data.session_id);
    } else {
      const content = document.getElementById('text-input').value.trim();
      if (content) await sendText(data.session_id, content);
    }
  } catch(e) {
    showAlert('send-alert', 'error', e.message);
  } finally {
    btn.disabled = false;
    btn.innerHTML = '<span>⚡</span> Create Beam';
  }
}

async function uploadFiles(sessionID) {
  if (selectedFiles.length === 0) return;

  const progress = document.getElementById('upload-progress');
  const fill = document.getElementById('progress-fill');
  const text = document.getElementById('progress-text');
  progress.classList.remove('hidden');

  for (let i = 0; i < selectedFiles.length; i++) {
    const f = selectedFiles[i];
    text.textContent = 'Uploading ' + f.name + '…';
    fill.style.width = ((i / selectedFiles.length) * 100) + '%';

    const form = new FormData();
    form.append('file', f);

    try {
      const res = await fetch(API + '/api/upload/' + sessionID, { method: 'POST', body: form });
      if (!res.ok) throw new Error('Upload failed');
    } catch(e) {
      showAlert('session-alert', 'error', 'Failed to upload ' + f.name);
    }
  }

  fill.style.width = '100%';
  text.textContent = 'All files uploaded ✓';
  showAlert('session-alert', 'success', selectedFiles.length + ' file(s) ready to beam!');
}

async function sendText(sessionID, content) {
  await fetch(API + '/api/text/' + sessionID, {
    method: 'POST',
    headers: {'Content-Type':'application/json'},
    body: JSON.stringify({ content })
  });
  showAlert('session-alert', 'success', 'Text ready to beam!');
}

async function loadQR(sessionID) {
  try {
    const res = await fetch(API + '/api/qr/' + sessionID + '?format=base64');
    const data = await res.json();
    const container = document.getElementById('qr-container');
    container.innerHTML = '<img src="data:image/png;base64,' + data.data + '" alt="QR Code">';
  } catch(e) {
    document.getElementById('qr-container').innerHTML = '<div style="color:var(--text2);font-size:13px;text-align:center">QR unavailable<br>Use PIN instead</div>';
  }
}

function cancelSession() {
  goHome();
}

// --- Join ---
async function joinByPIN() {
  const pin = getPIN();
  if (pin.length !== 6) { showAlert('join-alert','error','Please enter all 6 digits'); return; }

  const btn = document.getElementById('btn-join');
  btn.disabled = true;
  btn.innerHTML = '<div class="spinner"></div> Connecting…';

  try {
    const res = await fetch(API + '/api/join', {
      method: 'POST',
      headers: {'Content-Type':'application/json'},
      body: JSON.stringify({ pin })
    });
    const data = await res.json();
    if (data.error) throw new Error(data.error);
    renderReceiveView(data);
  } catch(e) {
    showAlert('join-alert', 'error', e.message);
  } finally {
    btn.disabled = false;
    btn.innerHTML = '<span>🔗</span> Connect';
  }
}

function renderReceiveView(session) {
  document.getElementById('view-home').classList.add('hidden');
  document.getElementById('view-receive').classList.remove('hidden');

  if (session.type === 'file' && session.files && session.files.length > 0) {
    document.getElementById('receive-title').textContent = '📁 Files ready';
    document.getElementById('receive-sub').textContent = session.files.length + ' file(s) available for download';
    document.getElementById('receive-files').classList.remove('hidden');

    const list = document.getElementById('download-list');
    list.innerHTML = session.files.map(f => {
      return '<div class="download-item"><span class="file-icon">' + fileIcon(f.name) + '</span><div class="file-info"><div class="file-name">' + escHtml(f.name) + '</div><div class="file-size">' + formatSize(f.size) + '</div></div><a class="btn-download" href="/api/download/' + session.id + '/' + f.id + '" download="' + escHtml(f.name) + '">⬇ Save</a></div>';
    }).join('');

    window._receiveSession = session;
  } else if (session.type === 'text' && session.content) {
    document.getElementById('receive-title').textContent = '📝 Text received';
    document.getElementById('receive-sub').textContent = 'Content ready to copy';
    document.getElementById('receive-text').classList.remove('hidden');
    document.getElementById('text-display').textContent = session.content;
    window._receivedText = session.content;
  } else {
    // Poll for content
    document.getElementById('receive-sub').textContent = 'Waiting for content…';
    pollSession(session.id);
  }
}

async function pollSession(id) {
  for (let i = 0; i < 60; i++) {
    await new Promise(r => setTimeout(r, 2000));
    try {
      const res = await fetch(API + '/api/session/' + id);
      const data = await res.json();
      if ((data.type === 'file' && data.files && data.files.length > 0) ||
          (data.type === 'text' && data.content)) {
        renderReceiveView(data);
        return;
      }
    } catch(e) { break; }
  }
}

function downloadAll() {
  const session = window._receiveSession;
  if (!session) return;
  for (const f of session.files) {
    const a = document.createElement('a');
    a.href = '/api/download/' + session.id + '/' + f.id;
    a.download = f.name;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
  }
}

async function copyText() {
  const text = window._receivedText || document.getElementById('text-display').textContent;
  try {
    await navigator.clipboard.writeText(text);
    const btn = document.getElementById('copy-btn');
    btn.classList.add('copied');
    btn.innerHTML = '✅ Copied!';
    setTimeout(() => {
      btn.classList.remove('copied');
      btn.innerHTML = '<span>📋</span> Copy to Clipboard';
    }, 2000);
  } catch(e) {
    showAlert('receive-alert', 'error', 'Copy failed. Please select text manually.');
  }
}

// --- Timer ---
function startTimer(expiresAt) {
  if (timerInterval) clearInterval(timerInterval);
  const dot = document.getElementById('timer-dot');
  const text = document.getElementById('timer-text');

  timerInterval = setInterval(() => {
    const now = new Date();
    const diff = Math.max(0, expiresAt - now);
    const mins = Math.floor(diff / 60000);
    const secs = Math.floor((diff % 60000) / 1000);
    text.textContent = 'Session expires in ' + String(mins).padStart(2,'0') + ':' + String(secs).padStart(2,'0');

    dot.className = 'timer-dot';
    if (mins < 2) dot.classList.add('danger');
    else if (mins < 5) dot.classList.add('warning');

    if (diff === 0) {
      clearInterval(timerInterval);
      showAlert('session-alert', 'error', 'Session expired. Please create a new beam.');
      text.textContent = 'Session expired';
    }
  }, 1000);
}

// --- Alert ---
function showAlert(id, type, msg) {
  const el = document.getElementById(id);
  if (!el) return;
  const icons = { success:'✅', error:'❌', info:'ℹ️' };
  el.innerHTML = '<div class="alert alert-' + type + '"><span>' + icons[type] + '</span><span>' + escHtml(msg) + '</span></div>';
  setTimeout(() => { if (el) el.innerHTML = ''; }, 5000);
}

// --- Route: handle /receive/{id} ---
(function() {
  const path = window.location.pathname;
  const m = path.match(/^\/receive\/(.+)$/);
  if (m) {
    const sessionID = m[1];
    fetch(API + '/api/session/' + sessionID)
      .then(r => r.json())
      .then(data => {
        if (data.error) {
          document.getElementById('view-home').classList.add('hidden');
          // Show error
          const main = document.querySelector('main .container');
          main.innerHTML = '<div class="card"><div class="card-body" style="text-align:center;padding:48px"><div style="font-size:48px;margin-bottom:16px">❌</div><div class="card-title" style="margin-bottom:8px">Session not found</div><div class="card-sub">This beam may have expired or been cancelled.</div><button class="btn btn-secondary" style="margin-top:24px" onclick="location.href=\'/\'">Go Home</button></div></div>';
        } else {
          renderReceiveView(data);
        }
      })
      .catch(() => location.href = '/');
  }
})();
</script>
</body>
</html>`
