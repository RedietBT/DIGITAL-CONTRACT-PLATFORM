# ✍️ Signature Service Documentation

This service handles the technical execution and legal integrity of digital signatures.

## 🛠️ Functional Overview
* **Multi-Type Support**: Handles Drawing (vector), Uploading (image), and Typing (text) signature styles.
* **Authorization**: Only users assigned as participants in the Contract Service can submit signatures.
* **Integrity**: Generates a unique SHA-256 hash for every signature mark to prevent tampering.
* **Positioning**: Stores `page_number`, `x/y` coordinates, and dimensions for PDF rendering.

## 🗄️ Database Architecture
* **Schema**: `signature_schema`
* **Primary Table**: `signatures` (ContractID, UserID, SignatureType, Hash, FileURL, VectorData)
* **Metadata**: Stores IP address and device info for the audit trail.

## 📡 Event Handling (RabbitMQ)
* **ContractPublishedEvent (Consume)**: When a contract is published, this service creates "placeholder" records for expected signers.
* **SignatureCreatedEvent (Publish)**: Emitted when a user successfully signs. Notifies the Analytics and Document (PDF) services.

## 🧪 Testing with Swagger
* **Local URL**: `http://localhost:8084/swagger/index.html`
* **Requirements**: Must provide a `Bearer <token>` in the Authorization header.