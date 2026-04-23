# Kopesa Frontend Consumer API

Tested against production on `2026-04-23`.

## Base URLs

- API base: `https://cloudcalls.easipath.com/backend-kopesa/api/v1`
- OpenAPI spec: `https://cloudcalls.easipath.com/backend-kopesa/api/openapi.yaml`
- Public health: `https://cloudcalls.easipath.com/backend-kopesa/api/v1/health`

## Response Envelope

All API responses use this envelope:

```json
{
  "data": {},
  "meta": {
    "requestId": "req_xxx"
  },
  "error": null
}
```

On failure:

```json
{
  "data": null,
  "meta": {
    "requestId": "req_xxx"
  },
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  }
}
```

## Auth

### Login

- `POST /auth/login`
- Public endpoint

Request:

```json
{
  "email": "admin@kopesa.co.za",
  "password": "Kopesa@2026"
}
```

Response data:

```json
{
  "accessToken": "jwt",
  "refreshToken": "jwt",
  "expiresIn": 900,
  "user": {
    "id": "u-9",
    "fullName": "System Admin",
    "email": "admin@kopesa.co.za",
    "role": "admin",
    "allowedRoles": ["admin"],
    "avatarColor": "hsl(220 13% 25%)",
    "active": true
  }
}
```

### Register

- `POST /auth/register`
- Public endpoint

Request:

```json
{
  "fullName": "QA Smoke User",
  "email": "qa.user@example.com",
  "password": "Kopesa@2026",
  "role": "branch_agent",
  "allowedRoles": ["branch_agent"],
  "branchId": "br-jhb-cbd",
  "areaId": "area-gp"
}
```

### Current User

- `GET /auth/me`
- Bearer token required

## Reference Data

- `GET /users`
- `GET /roles`
- `GET /branches`
- `GET /areas`

All require Bearer token.

## Loans

- `GET /loans`
- `GET /loans/:id`
- `GET /loans/applications`
- `GET /loans/applications/:id`
- `POST /loans/applications`
- `POST /loans/applications/:id/assess`
- `POST /loans/applications/:id/approve`
- `POST /loans/:id/disburse`

All require Bearer token.

## Arrears

- `GET /arrears/cases`
- `GET /arrears/cases/:id`
- `POST /arrears/cases/:id/allocate`
- `POST /arrears/cases/:id/ptps`
- `POST /arrears/cases/:id/payments`

All require Bearer token.

## Campaigns

- `GET /campaigns`
- `POST /campaigns`
- `GET /campaigns/:id/routes`
- `POST /campaigns/:id/routes`
- `POST /campaigns/:id/leads`
- `GET /leads`
- `POST /leads/:id/qualify`

All require Bearer token.

## Reporting

### Executive Overview

- `GET /reports/exec/overview`
- Bearer token required

Example response data:

```json
{
  "period": "mtd",
  "loans": {
    "activeCount": 1840,
    "disbursedCount": 312,
    "disbursedAmount": 450000000
  },
  "arrears": {
    "casesOpen": 412,
    "totalArrears": 87500000,
    "ptpRate": 0.41
  },
  "campaigns": {
    "active": 3,
    "leadsCaptured": 125,
    "leadsQualified": 71,
    "conversionRate": 0.18
  },
  "trend": [
    {
      "date": "2026-04-01",
      "disbursed": 15000000,
      "collected": 12000000
    }
  ]
}
```

## Attachments and Dropbox

### Direct Upload and Get Stream URL

- `POST /attachments/upload?context={context}&entityId={entityId}&fileName={fileName}`
- Bearer token required
- `multipart/form-data`
- File field name: `file`

Example:

- `context=loan_doc`
- `entityId=loan-123` or `ar-1002`
- `fileName=contract.pdf`

Success response data:

```json
{
  "id": "e4fbcc7e-3f04-11f1-b326-0242ac190002",
  "context": "loan_doc",
  "entityId": "qa-smoke",
  "fileName": "smoke.txt",
  "mimeType": "text/plain",
  "sizeBytes": 39,
  "url": "https://cloudcalls.easipath.com/backend-biatechdropbox/api/stream/6501e6d48209b5259b47c",
  "capturedBy": "u-9",
  "capturedAt": "2026-04-23T11:09:31Z",
  "sync": "synced",
  "provider": "dropbox",
  "path": "loan_doc/qa-smoke/smoke.txt",
  "revision": "6501e6d48209b5259b47c"
}
```

Frontend should persist the returned `url` for preview or streaming.

### List Attachments

- `GET /attachments?entityId={entityId}&context={context}`
- Bearer token required

### Get Single Attachment

- `GET /attachments/:id`
- Bearer token required

### Presign Flow

- `POST /attachments/presign`
- `POST /attachments`

These endpoints still exist, but for the current frontend the direct upload endpoint above is the simpler integration path.

## Messaging

### List Templates

- `GET /messaging/templates`
- Bearer token required

### List Logs

- `GET /messaging/log?context={context}&entityId={entityId}`
- Bearer token required

### Send Message

- `POST /messaging/send`
- Bearer token required

Request:

```json
{
  "context": "arrears",
  "entityId": "ar-1002",
  "channel": "sms",
  "to": "27821234567",
  "subject": "Optional for email",
  "templateId": "tpl-ar-reminder",
  "body": "Hello customer",
  "variables": {
    "client": "Bongani",
    "daysPastDue": "12"
  },
  "nextTouchAt": "2026-04-25T10:00:00Z"
}
```

Supported `channel` values:

- `sms`
- `email`
- `whatsapp`
- `both`

Channel guidance:

- For SMS: set `channel` to `sms`, `to` must be a phone number.
- For Email: set `channel` to `email`, `to` can be one email or multiple emails separated by `,` or `;`.
- For WhatsApp: set `channel` to `whatsapp`, `to` must be a phone number.
- For Both: set `channel` to `both`; current behavior sends SMS and WhatsApp using the same `to` value.

Provider-backed frontend use cases:

- Send SMS: `POST /messaging/send` with `channel: "sms"`
- Send Email: `POST /messaging/send` with `channel: "email"`
- Send WhatsApp: `POST /messaging/send` with `channel: "whatsapp"`

Response:

- HTTP `202 Accepted`
- Message log entry is returned in `data`

## Tested Production Paths

These were verified on production on `2026-04-23`:

- `GET /backend-kopesa/api/v1/health`
- `POST /auth/login`
- `GET /auth/me`
- `POST /auth/register`
- `GET /loans`
- `GET /reports/exec/overview`
- `GET /messaging/templates`
- `POST /messaging/send`
- `GET /messaging/log`
- `POST /attachments/upload`
- `GET /attachments`
- Dropbox stream URL returned by upload

## Seed Admin Account

For QA only:

- Email: `admin@kopesa.co.za`
- Password: `Kopesa@2026`
