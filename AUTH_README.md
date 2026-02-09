# Auth Service Technical Guide (Corrected)

This service handles user identity, secure authentication, and account recovery for the Digital Contract Platform.

## 🛠️ Swagger Documentation

If you change any code in the handlers (annotations), you **must** update the documentation so the UI stays in sync.

**Command to update:**
To ensure Swagger picks up models from your internal packages, use the `--parseDependency` flag.

1. Open your terminal in the `backend` folder.
2. Run:
`swag init -g cmd/auth/main.go --output cmd/auth/docs --parseDependency --parseInternal`
3. Rebuild your containers: `docker-compose up --build auth`

---

## 📧 Email Testing (Mailhog)

We use **Mailhog** as a local SMTP server to test emails (Welcome emails, Password Resets) without sending them to real addresses.

* **Web Interface**: [http://localhost:8025](https://www.google.com/search?q=http://localhost:8025)
* **SMTP Port**: `1025`

---

## 🧪 Manual Test Scenarios

### 1. Register & Login

* **Register**: `POST /auth/register`. **Note**: This automatically triggers a RabbitMQ event to the Profile Service.
* **Login**: `POST /auth/login`.
* **Fix Applied**: Requires a `UNIQUE` constraint on `user_id` in the `refresh_tokens` table to handle the `ON CONFLICT` logic.

### 2. Token Refreshing

* **Refresh**: `POST /auth/refresh`
* **Payload**: `{"refresh_token": "..."}`
* Returns a new access token. If a user is `suspended`, this will fail.

---

### 🛡️ Role-Based Access Control (RBAC)

* **To Promote a User to Admin (SQL):**

```sql
UPDATE auth_schema.users SET role = 'admin' WHERE email = 'your@email.com';

```

### 👤 User Account Management ("Me" Routes)

These require a valid JWT in the `Authorization: Bearer <token>` header.

| Endpoint | Method | Description |
| --- | --- | --- |
| `/auth/me` | `GET` | Returns current user's auth data. |
| `/auth/me/email` | `PUT` | Updates the user's email address. |
| `/auth/me/password` | `PUT` | Changes password (requires validation). |
| `/auth/me` | `DELETE` | **Sync Delete**: Triggers RabbitMQ to delete the user's profile. |

---

### 🔐 Security & Integration Overview

* **Passwords**: Hashed using **Bcrypt**.
* **Event-Driven**: All deletions (Self-delete or Admin-delete) broadcast a `user.deleted` message to RabbitMQ to maintain cross-service data integrity.
* **Database**: Uses `auth_schema` to isolate identity data from business logic.

### 🛠️ Database Debugging

```bash
# View all users and their status
docker exec -it db-1 psql -U guest -d digital_contract_db -c "SELECT email, role, status FROM auth_schema.users;"

# View active refresh tokens (Verify UNIQUE constraint)
docker exec -it db-1 psql -U guest -d digital_contract_db -c "\d auth_schema.refresh_tokens"

```

---
