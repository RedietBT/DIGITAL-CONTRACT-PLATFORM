Location: C:\Users\HP\Documents\projects\digital-contract-platform\README.md Purpose: Overall project summary.

Markdown
# Digital Contract Platform

A microservices-based platform for managing digital contracts.

## 🏗️ Architecture
- **Auth Service**: Handles User Registration, Login, and JWT generation.
- **Contract Service**: Manages the lifecycle of digital contracts.
- **Database**: PostgreSQL with schema-based isolation.

## 🚀 How to Run
1. Navigate to the deployments folder:
   `cd deployments`
2. Start the services:
   `docker-compose up --build`

## 🔗 Access Points
- **Auth Swagger UI**: http://localhost:8080/swagger/index.html
- **Database**: localhost:5432