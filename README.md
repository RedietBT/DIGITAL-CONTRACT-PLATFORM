# 🚀 Digital Contract Platform

A high-performance, microservices-based platform for managing digital agreements, utilizing an event-driven architecture and PostgreSQL schema isolation.

## 🏗️ System Architecture
* **Auth Service (`:8080`)**: Identity management and JWT issuance.
* **Contract Service (`:8083`)**: Document lifecycle and versioning.
* **Profile Service (`:8082`)**: User metadata and personal details.
* **Global Swagger (`:8085`)**: Centralized API documentation for all services.

## 🚀 Getting Started

### 1. Generate All Documentation
Run the following commands to update the blueprints for the Global Swagger UI:
```bash
swag init -g cmd/auth/main.go -o cmd/auth/docs --instanceName auth
swag init -g cmd/profile/main.go -o cmd/profile/docs --instanceName profile
swag init -g cmd/contract/main.go -o cmd/contract/docs --instanceName contract

```

### 2. Launch Platform

```bash
cd deployments
docker-compose up --build

```

## 📡 Core Event Flows

| Event | Source | Action |
| --- | --- | --- |
| **UserCreated** | Auth Service | Profile Service initializes user data. |
| **UserDeleted** | Auth Service | Profile Service deletes profile; Contract Service purges owned contracts. |

## 🔗 Infrastructure Access

* **Global Docs**: [http://localhost:8085](https://www.google.com/search?q=http://localhost:8085)
* **MailHog**: [http://localhost:8025](https://www.google.com/search?q=http://localhost:8025) (SMTP testing)
* **RabbitMQ**: [http://localhost:15672](https://www.google.com/search?q=http://localhost:15672) (guest/guest)

---
