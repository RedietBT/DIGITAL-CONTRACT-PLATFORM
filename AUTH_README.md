# Auth Service Technical Guide

## 🛠️ Swagger Documentation
If you change any code in the handlers, you **must** update the documentation so the UI stays in sync.

**Command to update:**
Run this from the `backend` folder:
`swag init -g cmd/auth/main.go --output cmd/auth/docs`

## 🧪 Manual Test Scenarios

### 1. Register a User
- **Endpoint**: `POST /auth/register`
- **Payload**: 
  ```json
  {
    "email": "user@example.com",
    "password": "securePassword123"
  }
2. Login
Endpoint: POST /auth/login

Payload:

JSON
{
  "email": "user@example.com",
  "password": "securePassword123"
}
Expected Response: A 200 OK with a JSON object containing the token.

🔐 Security Overview
Passwords: Hashed using Bcrypt (never stored in plain text).

Authentication: Stateless via JSON Web Tokens (JWT).

Validation: Strict input validation for emails and password length.


---

### 💡 Why we included the extra parts:
You asked if you should only paste the first bit—you *could*, but including the **Login** test data and the **Security Overview** is better because:
1.  **Context**: It reminds you that `bcrypt` is working in the background.
2.  **Efficiency**: When you come back to this project in a month, you won't have to look through your Go code to remember what the JSON keys were named.

### 🚀 Final Step
Now that you have the files:
1.  Open your terminal in the `backend` folder.
2.  Run that `swag init` command one last time to make sure your Swagger is 100% updated.
3.  Go to `http://localhost:8080/swagger/index.html` and **Register** then **Login**.

**Would you like me to show you how to check if your JWT is valid using an online tool
