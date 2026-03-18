# 📄 Contract Service Documentation

This service manages the lifecycle, versioning, and participant permissions of digital agreements.

## 🛠️ Functional Overview
* **Drafting**: Create new contracts with an initial Version 1.
* **Participant Assignment**: Assign specific User IDs as signers via `POST /contracts/{id}/participants`.
* **Lifecycle Tracking**: Moves contracts from `draft` -> `pending` (once published) -> `completed`.
* **Ownership Control**: JWT-based middleware ensures only the owner can Edit, Delete, or Invite signers.

## 🗄️ Database Architecture
* **Schema**: `contract_schema`
* **Primary Table**: `contracts` (OwnerID, Title, Status)
* **Participants Table**: `contract_participants` (UserID, Role, IsRequired)
* **History Table**: `contract_versions` (Stores `content_json` snapshots)

## 📡 Event Handling (RabbitMQ)
* **UserDeletedEvent (Consume)**: Purges all contracts owned by the deleted user.
* **ContractPublishedEvent (Publish)**: Emitted when participants are assigned. Notifies the Signature Service to prepare the signing workflow.

## 🧪 Testing with Swagger
* **Local URL**: `http://localhost:8083/swagger/index.html`
* **Requirements**: Must provide a `Bearer <token>` in the Authorization header.