---
name: kopesa-backend-handoff
description: Use when continuing work on this repository so the agent follows the established DDD structure, knows the migration/bootstrap workflow, and leaves updated handoff notes for the next agent.
---

# Kopesa Backend Handoff

Use this skill when working inside `kopesa-loan-platform-api`.

## Project structure

- Keep business modules under `internal/<domain>/domain`, `application`, `interfaces/controllers`, and `infrastructure/cassandra`.
- Core business modules are `customer`, `loan`, `arrears`, and `campaign`.
- Supporting modules follow the same layout: `messaging`, `attachments`, `reporting`, `audit`.
- Keep HTTP controllers thin. Business rules belong in `domain` and orchestration/use-cases in `application`.

## Runtime entrypoints

- API entrypoint: `cmd/api/main.go`
- Cassandra migration and seeding: `cmd/migrate/main.go`
- Schema: `migrations/001_init.cql`

## Cassandra notes

- Keyspace in use: `kopesa_loan_platform`
- Remote host used previously: `safer.easipath.com`
- The migration command creates the keyspace and seeds the frontend-aligned data.
- Existing deployments require schema evolution through `cmd/migrate/main.go`, not only `migrations/001_init.cql`.
- Current compatibility alters add:
  - `attachments_attachments.provider`
  - `attachments_attachments.path`
  - `attachments_attachments.revision`
  - `messaging_logs.provider_ref`

## Workflow

1. Read `documents/DDD_Guide.pdf` and `documents/kopesa-backend-api-spec.md` before structural changes.
2. Preserve the current DDD folder layout.
3. After changes, run `gofmt -w cmd internal`.
4. Verify with:
```bash
env GOCACHE=/home/biangacila/GolandProjects/kopesa-loan-platform-api/.cache/go-build \
GOMODCACHE=/home/biangacila/GolandProjects/kopesa-loan-platform-api/.cache/gomod \
go test ./...
```
5. Update the handoff memory note at `/home/biangacila/.codex/memories/kopesa-loan-platform-api.md`.

## External integrations

- Messaging provider adapter is implemented in `internal/messaging/infrastructure/provider`.
  - SMS + email use `cloudcalls.easipath.com/backend-email-service/api/v1`
  - WhatsApp uses Evolution API on `safer.easipath.com:8080`
- Dropbox adapter is implemented in `internal/attachments/infrastructure/provider`.
- Kafka is still abstracted behind a no-op publisher stub in `internal/platform/kafka`.
- Public API base path is `/backend-kopesa/api/v1`.
- Public OpenAPI path is `/backend-kopesa/api/openapi.yaml`.
- Frontend consumer document lives at `documents/frontend-consumer-api.md`.
- If changing provider payloads, read:
  - `documents/external-provider/API_EMAIL_SMS_GATEWAY.md`
  - `documents/external-provider/API_INTEGRATION_GUIDE-DROPBOX.md`
  - `documents/external-provider/API_INTEGRATION_WHATSAPP_SENDING_MESSANGE.md`
