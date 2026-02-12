# 📄 Contract Service Documentation

This service manages the lifecycle, versioning, and security of digital agreements within the platform.

## 🛠️ Functional Overview
* **Drafting**: Create new contracts with an initial Version 1.
* **Metadata Management**: Update titles and descriptions for existing drafts.
* **Versioning**: Every content change is tracked in `contract_schema.contract_versions` using JSON storage.
* **Ownership Control**: Uses JWT-based middleware to ensure only the owner can Edit or Delete a contract.

## 🗄️ Database Architecture
* **Schema**: `contract_schema`
* **Primary Table**: `contracts` (Stores OwnerID, Title, Status)
* **History Table**: `contract_versions` (Stores `content_json` and timestamp)

## 📡 Event Handling (RabbitMQ)
As per the platform's data consistency policy:
* **UserDeletedEvent**: When a user is removed via the Auth service, this service consumes the event and purges all contracts owned by that `userID` to prevent orphaned data.

## 🧪 Testing with Swagger
* **Local URL**: `http://localhost:8081/swagger/index.html`
* **Requirements**: Must provide a `Bearer <token>` in the Authorization header.