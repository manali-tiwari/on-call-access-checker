# On-Call Access Checker (Backend)

You can run the backend server directly in your browser using CodeSandbox.io - no local setup required!

## 🚀 Quick Start Guide

1. **Open in CodeSandbox**  
   [![Open in CodeSandbox](https://codesandbox.io/static/img/play-codesandbox.svg)](https://codesandbox.io/p/github/manali-tiwari/on-call-access-checker/master?file=/backend/cmd/server/main.go)

2. **Configure Environment**  
   - Open the `.env` file in the `backend` folder
   - Add your credentials (or use mock mode by leaving blank):
     ```
     OKTA_HOST=your-org.okta.com
     OKTA_TOKEN=your-api-token
     AWS_REGION=us-east-1
     AWS_PROFILE=default
     ```

3. **Start the Server**  
   In the terminal:
   ```bash
   cd backend
   go run ./cmd/server/main.go

4. 🔧 **Mock Mode Features**

    When no credentials are provided:
    - ✅ Auto-enables mock Okta responses
    - ✅ Simulates AWS profile data
    - 🌐 All API endpoints remain functional

5. 🌐 **API Endpoints**
   
   *POST* /api/check-access

    Request:

    ```json
    {
    "email": "user@example.com",
    "environment": "Production"
    }
    ```

    Response:


    ```json
    {
    "vpn": true,
    "production": true,
    "configTool": true,
    "currentProfile": "mock-profile",
    "validUntil": "2025-04-22T14:30:00Z"
    }
    ```

6. 💻 **Tips**

    Terminal Shortcuts:
    - Ctrl+`` to toggle terminal
    - Right-click ports to "Open in New Tab"

    Port-Forwarding:

    The server auto-opens on port 8080. To change:

    ```go
    // In cmd/server/main.go
    port := "3000" // CodeSandbox prefers port 3000
    ```

    Debugging Common Issues:

    ```bash
    # Port in use? Find and kill process:
    lsof -i :8080
    kill -9 PID

    # Reset environment:
    go clean -modcache
    go mod tidy
    ```

7. 🛠️ **Dependencies**
    
    Automatically installed in CodeSandbox:
    - Go 1.22+
    - Gin Web Framework
    - Okta SDK
    - AWS SDK   

8. 📦 **Alternative Local Setup**
    
    ```bash
    git clone https://github.com/yourusername/on-call-access-checker.git
    cd on-call-access-checker/backend
    go mod download
    go run ./cmd/server/main.go
    ```
