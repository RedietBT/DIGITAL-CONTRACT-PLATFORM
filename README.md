# Digital Contract Platform

A microservices-based platform for managing digital contracts, utilizing an event-driven architecture for high data consistency.

## 🏗️ Architecture

* **Auth Service (`:8080`)**: Handles user registration, login, and JWT generation. It broadcasts user events (created/deleted) to the message broker.
* **Profile Service (`:8082`)**: Manages detailed user information in `profile_schema.profile`. It stays in sync with the Auth Service via RabbitMQ.
* **Message Broker**: **RabbitMQ** handles asynchronous communication between services to ensure data integrity.
* **Database**: **PostgreSQL** with schema-based isolation (e.g., `auth_schema`, `profile_schema`).

## 🚀 How to Run

1. **Generate Documentation**:
Ensure Swagger docs are up to date for both services:
```bash
swag init -g cmd/auth/main.go --output cmd/auth/docs
swag init -g cmd/profile/main.go --output cmd/profile/docs

```


2. **Deploy with Docker**:
Navigate to the deployments folder and start the stack:
```bash
cd deployments
docker-compose up --build

```



## 🔗 Access Points

* **Auth Swagger UI**: [http://localhost:8080/swagger/index.html](https://www.google.com/search?q=http://localhost:8080/swagger/index.html)
* **Profile Swagger UI**: [http://localhost:8082/swagger/index.html](https://www.google.com/search?q=http://localhost:8082/swagger/index.html)
* **RabbitMQ Management**: [http://localhost:15672](https://www.google.com/search?q=http://localhost:15672) (User/Pass: guest/guest)
* **MailHog (Email Testing)**: [http://localhost:8025](https://www.google.com/search?q=http://localhost:8025)

## 📡 Event-Driven Logic

* **User Created**: When a user registers, the Profile Service automatically initializes their profile.
* **User Deleted**: When a user is deleted from the Auth service, the corresponding profile is purged from the database via a RabbitMQ event.

---