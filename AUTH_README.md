---

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

* **Web Interface**: [http://localhost:8025](https://www.google.com/search?q=http://localhost:8025)
* **SMTP Port**: `1025`

---

## 🧪 Manual Test Scenarios

### 1. Register & Login

* **Register**: `POST /auth/register` (Email/Password)
* **Login**: `POST /auth/login`
* Returns `access_token` (short-lived) and `refresh_token` (7 days).


* **Refresh**: `POST /auth/refresh`
* **Payload**: `{"refresh_token": "..."}`
* Returns a new access token without requiring a password.



### 2. Password Recovery

* **Forgot Password**: `POST /auth/forgot-password` (Sends token to Mailhog).
* **Reset Password**: `POST /auth/reset-password` (Requires email, token, and new_password).

---

### 🛡️ Role-Based Access Control (RBAC)

The service distinguishes between regular `users` and `admins`. Roles are enforced via the `RoleMiddleware`.

* **To Promote a User to Admin (SQL):**

```sql
UPDATE auth_schema.users SET role = 'admin' WHERE email = 'your@email.com';

```

### 👤 User Account Management ("Me" Routes)

These require a valid JWT in the `Authorization: Bearer <token>` header.

| Endpoint | Method | Description |
| --- | --- | --- |
| `/auth/me/email` | `PUT` | Updates the user's email address. |
| `/auth/me/password` | `PUT` | Changes password (Checks `old_password`). |
| `/auth/me` | `DELETE` | Self-deletion of account. |
| `/auth/logout` | `POST` | Informs client to clear session. |

### 🔑 Admin Management

Endpoints restricted to users with the `admin` role.

* **List All Users**: `GET /auth/admin/users`
* **Update User Status**: `PUT /auth/admin/user-status`
* **Payload**: `{"user_id": "...", "status": "suspended"}`
* Options: `active`, `suspended`, `deactivated`.


* **Hard Delete User**: `DELETE /auth/admin/users?id=<uuid>`

---

### 🔐 Security Overview

* **Passwords**: Hashed using **Bcrypt**.
* **Tokens**: Dual-token system (JWT Access + Persistent Refresh Token).
* **Status Check**: Accounts marked `suspended` or `deactivated` are blocked from logging in or refreshing tokens.
* **Middleware**: Chained architecture (Logger -> Auth -> Role).

### 🛠️ Database Debugging

```bash
# View all users and their status
docker exec -it deployments-db-1 psql -U postgres -d digital_contract_db -c "SELECT email, role, status, last_login_at FROM auth_schema.users;"

# View active refresh tokens
docker exec -it deployments-db-1 psql -U postgres -d digital_contract_db -c "SELECT * FROM auth_schema.refresh_tokens;"

```

---
