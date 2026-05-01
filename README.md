# msg_forwarder

A lightweight message forwarding service that receives HTTP POST requests and forwards the content to email.

## Quick Start

```bash
cp .env.example .env
go run .
```

## Configuration

All config via `.env` or environment variables:

| Variable | Required | Default | Description |
|---|---|---|---|
| `SMTP_HOST` | No | `smtp.163.com` | SMTP server host |
| `SMTP_PORT` | No | `465` | SMTP server port |
| `SMTP_USER` | Yes | - | SMTP login email |
| `SMTP_PASSWORD` | Yes | - | SMTP authorization code |
| `AUTH_TOKEN` | Yes | - | API bearer token |
| `SERVER_PORT` | No | `8080` | HTTP listen port |
| `SENDER_NAME` | No | `Msg Forwarder` | Display name of sender |

## API

### POST /api/send

**Authentication:** `Authorization: Bearer <token>` or `?token=<token>`

**Request body:**

```json
{
  "to": "recipient@example.com",
  "subject": "Alert",
  "content": "Something happened",
  "name": "Optional sender name"
}
```

| Field | Required | Description |
|---|---|---|
| `to` | Yes | Recipient email |
| `content` | Yes | Message body |
| `subject` | No | Email subject (defaults to "New Message" or "Message from {name}") |
| `name` | No | Sender name, appended to email body |

**Response:**

```json
{ "success": true, "message": "email sent to recipient@example.com" }
```

**Example:**

```bash
curl -X POST http://localhost:8080/api/send \
  -H "Authorization: Bearer your-secret-token" \
  -H "Content-Type: application/json" \
  -d '{"to":"user@example.com","content":"Hello!","name":"Alice"}'
```