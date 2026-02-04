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
```
### 2. Login
- **Endpoint**: `POST /auth/login`
- **Payload**:
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```
Expected Response: 200 OK with a JSON token.

### 3. Forgot Password (Request Token)
- **Endpoint**: `POST /auth/forgot-password`
- **Payload**: 
```json
{
  "email": "user@example.com"
}
```
Action: Check Mailhog for the 8-character token.

### 4. Reset Password (Update Password)
- **Endpoint**: `POST /auth/reset-password`
- **Payload**: 
```json
{
  "email": "user@example.com",
  "token": "PASTE_THE_CODE_FROM_MAILHOG",
  "new_password": "newSecurePassword456"
}
```
### 🔐 Security Overview
- **Passwords**: Hashed using Bcrypt. We never store plain text passwords.
- **Authentication**: Stateless via JWT (JSON Web Tokens).
- **Recovery**: One-time-use tokens with a 15-minute expiry stored in auth_schema.password_resets.
- **Validation**: Strict input checks for email format and minimum password length (8 chars).

### 🛠️ Database Debugging
If you need to see what is happening inside the database:

View Reset Tokens: 
```
docker exec -it deployments-db-1 psql -U postgres -d digital_contract_db -c "SELECT * FROM auth_schema.password_resets;"
```

View Users: 
```
docker exec -it deployments-db-1 psql -U postgres -d digital_contract_db -c "SELECT email, password_hash FROM auth_schema.users;"
```
---
It is the perfect time to update your documentation. Your README is currently a "User-facing" guide, so we need to add the **Admin Features** and the **Account Management** (Me) section. This helps anyone else (or future you) understand how the RBAC (Role-Based Access Control) works.

Here is the markdown block you can append to your existing `README.md`:

---

### 🛡️ Role-Based Access Control (RBAC)

The service distinguishes between regular `users` and `admins`. Roles are enforced via the `RoleMiddleware`.

* **To Promote a User to Admin (SQL):**

```sql
UPDATE auth_schema.users SET role = 'admin' WHERE email = 'your@email.com';

```

### 👤 User Account Management ("Me" Routes)

These endpoints allow users to manage their own data. They require a valid JWT in the `Authorization: Bearer <token>` header.

| Endpoint | Method | Description |
| --- | --- | --- |
| `/auth/me` | `GET` | Returns the current user's profile info. |
| `/auth/me/email` | `PUT` | Updates the user's email address (Checks for uniqueness). |
| `/auth/me/password` | `PUT` | Changes password (Requires `old_password` verification). |
| `/auth/me` | `DELETE` | Allows a user to delete their own account. |

### 🔑 Admin Management

Endpoints restricted to users with the `admin` role.

* **List All Users**: `GET /auth/admin/users`
* Returns a full list of registered users (excluding password hashes).


* **Delete Any User**: `DELETE /auth/admin/users?id=<uuid>`
* Forcefully removes a user account from the system.



---

### 🚦 Middleware Architecture

The service uses a "Chained Middleware" pattern to secure routes:

1. **LoggerMiddleware**: Tracks every incoming request.
2. **AuthMiddleware**: Extracts and validates the JWT. Injects `UserID` and `Role` into the Request Context.
3. **RoleMiddleware**: Checks the context for specific permissions (e.g., `admin`) before allowing the request to hit the handler.

### 🏗️ Database Schema: Extended

We have two primary tables in the `auth_schema`:

* `users`: Stores core identity (ID, Email, Role, Password Hash).
* `password_resets`: Temporary storage for recovery tokens (Linked via Email with `ON DELETE CASCADE`).

---

### 🧼 Code Quality & Clean Architecture

The project follows the **Repository -> Service -> Handler** pattern:

* **Repository**: Pure SQL logic via `database/sql`.
* **Service**: Business logic (hashing, uniqueness checks, token validation).
* **Handler**: HTTP protocol handling (JSON parsing, status codes).

---
