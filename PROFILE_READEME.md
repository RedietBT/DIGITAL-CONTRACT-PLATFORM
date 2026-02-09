# Profile Service 👤

This microservice manages user profiles within the Digital Contract Platform. it is built with **Go (Gin)** and maintains data integrity with the Auth Service via **RabbitMQ**.

## 🚀 Features

* **Profile Management**: Retrieve and update user-specific data.
* **Event-Driven Sync**: Automatically creates a profile when a user registers and deletes the profile when a user is removed from the Auth service.
* **Input Validation**: Strict validation on all profile update requests.
* **Swagger Documentation**: Interactive API testing available at `/swagger/index.html`.

## 🏗 Architecture

* **Database**: PostgreSQL (`profile_schema.profile` table).
* **Messaging**: RabbitMQ (Consumer for `user.created` and `user.deleted` events).
* **Security**: JWT Authentication shared with Auth Service.

## 🛠 API Endpoints

| Method | Endpoint | Description | Auth Required |
| --- | --- | --- | --- |
| **GET** | `/profile/me` | Retrieve current user's profile | **Yes (Bearer Token)** |
| **PUT** | `/profile/me` | Update profile (Validated) | **Yes (Bearer Token)** |
| **GET** | `/health` | Service health check | No |

## 📡 Event Handling (RabbitMQ)

The service listens for the following events from the `user_events` exchange:

1. **User Created**: Initializes a new entry in `profile_schema.profile`.
2. **User Deleted**: Performs a permanent cleanup of the profile data when the auth account is closed.

---