# Auth Service Technical Guide

This service handles user identity, secure authentication, and account recovery for the Digital Contract Platform.

## 🛠️ Swagger Documentation
If you change any code in the handlers (annotations), you **must** update the documentation so the UI stays in sync.

**Command to update:**
1. Open your terminal in the `backend` folder.
2. Run:
   `swag init -g cmd/auth/main.go --output cmd/auth/docs`
3. Rebuild your containers: `docker-compose up --build`

---

## 📧 Email Testing (Mailhog)
We use **Mailhog** as a local SMTP server to test emails without sending them to real addresses.

- **Web Interface**: [http://localhost:8025](http://localhost:8025)
- **SMTP Port**: `1025` (Internal to Docker)

---

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
Expected Response: 200 OK with a JSON token.

3. Forgot Password (Request Token)
Endpoint: POST /auth/forgot-password

Payload:

JSON
{
  "email": "user@example.com"
}
Action: Check Mailhog for the 8-character token.

4. Reset Password (Update Password)
Endpoint: POST /auth/reset-password

Payload:

JSON
{
  "email": "user@example.com",
  "token": "PASTE_THE_CODE_FROM_MAILHOG",
  "new_password": "newSecurePassword456"
}
🔐 Security Overview
Passwords: Hashed using Bcrypt. We never store plain text passwords.

Authentication: Stateless via JWT (JSON Web Tokens).

Recovery: One-time-use tokens with a 15-minute expiry stored in auth_schema.password_resets.

Validation: Strict input checks for email format and minimum password length (8 chars).

🛠️ Database Debugging
If you need to see what is happening inside the database:

View Reset Tokens: docker exec -it deployments-db-1 psql -U postgres -d digital_contract_db -c "SELECT * FROM auth_schema.password_resets;"

View Users: docker exec -it deployments-db-1 psql -U postgres -d digital_contract_db -c "SELECT email, password_hash FROM auth_schema.users;"


---

### 🏁 What this covers:
* **Registration**: The initial setup.
* **Login**: How to get your JWT.
* **Forgot/Reset**: The full loop using Mailhog.
* **Database**: The exact commands to verify your data.

**Once you've pasted this, your documentation is perfect. Should we move on to building the Auth Middleware so you can start protecting your future routes?**