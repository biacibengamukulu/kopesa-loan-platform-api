# Kopesa - System Scope & Requirements Document

## Executive Summary

**Kopesa** is a comprehensive field lending and collections management platform designed for microfinance institutions operating in South Africa. The system enables end-to-end loan lifecycle management from lead capture through collections, with a strong focus on field operations, mobile-first workflows, and offline-capable functionality.

---

## Business Context

### Target Market
- Microfinance institutions in South Africa
- Branch-based lending operations
- Field agents collecting loan applications and payments
- Multi-branch operations with area management oversight

### Key Business Requirements
1. **Field-First Design**: System must work in areas with poor connectivity
2. **Offline Capability**: Core functions must work offline and sync later
3. **Multi-Role Access**: 10+ distinct user roles with granular permissions
4. **Audit Trail**: Complete tracking of all financial actions
5. **Maker-Checker**: Dual approval for critical financial operations
6. **Regulatory Compliance**: KYC, affordability assessments, NCA compliance

---

## System Architecture Overview

### Frontend Stack
- **Framework**: React 18 + TypeScript
- **Build Tool**: Vite 5
- **Styling**: Tailwind CSS v3 + shadcn/ui components
- **State Management**: Redux Toolkit
- **Routing**: React Router v6
- **Charts**: Recharts for analytics

### Backend Architecture (Planned)
- **API Framework**: Go (Fiber)
- **Database**: Cassandra (distributed, write-heavy workload)
- **Event Streaming**: Apache Kafka (event sourcing, audit logs)
- **Authentication**: JWT with refresh tokens
- **File Storage**: S3-compatible object storage

### Deployment
- **Frontend**: Lovable Cloud (current)
- **Backend**: Kubernetes cluster (planned)
- **CDN**: CloudFlare (planned)

---

## Domain Model & Core Entities

### 1. Identity & Access Management (IAM)

#### Roles (10 distinct roles)
| Role ID | Label | Scope | Primary Function |
|---------|-------|-------|------------------|
| `field_officer` | Field Officer | Field | Campaign planning, lead capture, loan applications |
| `collector` | Collector | Field | Arrears follow-up, PTP capture, payment collection |
| `branch_agent` | Branch Agent | Branch | Walk-in originations, basic collections |
| `consultant` | Loan Consultant | Branch | Loan assessment, affordability analysis |
| `branch_manager` | Branch Manager | Branch | Branch oversight, allocation, approvals |
| `area_manager` | Area Manager | Area | Multi-branch oversight, escalations, write-offs |
| `finance_officer` | Finance Officer | Global | Reconciliation, disbursements, write-offs |
| `compliance_officer` | Compliance Officer | Global | KYC, audit, regulatory compliance |
| `admin` | System Admin | Global | Users, branches, configuration |
| `executive` | Executive | Global | Read-only KPIs, dashboards |

#### Permission System
- **Granular Permissions**: 30+ specific permissions
- **Role-Permission Mapping**: Each role has a defined set of permissions
- **Scope Enforcement**: Field, Branch, Area, Global data visibility
- **Maker-Checker**: Critical operations require dual approval

### 2. Loans Domain

#### Loan Lifecycle States
```
application → assessment → approved → disbursed → active → closed
     ↓           ↓           ↓          ↓
  declined    declined    declined   written_off
```

#### Key Entities
- **Loan**: Core loan aggregate with principal, term, rate, status
- **LoanApplication**: Application data with affordability assessment
- **Disbursement**: Payment record with reconciliation tracking

### 3. Arrears Domain

#### Arrears Case Lifecycle
```
new → allocated → in_progress → ptp → paid
                    ↓              ↓
              escalated      broken_ptp
                    ↓              ↓
              written_off    escalated
```

#### Key Entities
- **ArrearsCase**: Aggregate tracking DPD, amounts, assignments
- **PTP (Promise to Pay)**: Scheduled payment commitment
- **PaymentEvidence**: Proof of payment with attachment

### 4. Campaigns Domain

#### Campaign Lifecycle
```
planning → active → completed
    ↓         ↓
cancelled  cancelled
```

#### Key Entities
- **Campaign**: Marketing initiative with targets and dates
- **CampaignRoute**: Daily route with GPS stops
- **RouteStop**: Individual location visit record
- **Lead**: Captured prospect with qualification status

### 5. Messaging Domain

#### Channels
- **SMS**: Primary channel for arrears and loan notifications
- **WhatsApp**: Rich messaging for campaigns and documents

#### Key Entities
- **MessageTemplate**: Reusable message with variable substitution
- **MessageLogEntry**: Audit trail of all sent messages

### 6. Attachments Domain

#### Contexts
- **arrears_pop**: Proof of payment
- **loan_doc**: Loan documentation
- **kyc**: Know Your Customer documents
- **lead_doc**: Lead capture documents

#### Key Features
- **Offline-First**: Base64 storage with sync queue
- **Sync Status**: pending → synced → failed

---

## API Specification Summary

### Authentication Endpoints
```
POST /api/v1/auth/login          → JWT tokens
POST /api/v1/auth/refresh        → Token refresh
POST /api/v1/auth/logout         → Token revocation
GET  /api/v1/auth/me             → Current user
```

### Loans Endpoints
```
GET    /api/v1/loans                    → List loans
POST   /api/v1/loans                    → Create loan
GET    /api/v1/loans/:id                → Get loan
POST   /api/v1/loans/:id/assess         → Submit assessment
POST   /api/v1/loans/:id/approve        → Approve loan
POST   /api/v1/loans/:id/disburse       → Disburse loan
GET    /api/v1/loans/applications       → List applications
POST   /api/v1/loans/applications       → Submit application
```

### Arrears Endpoints
```
GET    /api/v1/arrears                  → List cases
POST   /api/v1/arrears/import           → Bulk import
GET    /api/v1/arrears/:id              → Get case
POST   /api/v1/arrears/:id/allocate     → Assign collector
POST   /api/v1/arrears/:id/ptp          → Create PTP
POST   /api/v1/arrears/:id/payment      → Record payment
POST   /api/v1/arrears/:id/escalate     → Escalate case
POST   /api/v1/arrears/:id/writeoff     → Write off case
```

### Campaigns Endpoints
```
GET    /api/v1/campaigns                → List campaigns
POST   /api/v1/campaigns                → Create campaign
GET    /api/v1/campaigns/:id            → Get campaign
POST   /api/v1/campaigns/:id/activate   → Activate campaign
POST   /api/v1/campaigns/:id/complete   → Complete campaign
GET    /api/v1/campaigns/:id/routes     → List routes
POST   /api/v1/campaigns/:id/routes     → Create route
GET    /api/v1/routes/:id               → Get route
POST   /api/v1/routes/:id/stops         → Record stop visit
GET    /api/v1/leads                    → List leads
POST   /api/v1/leads                    → Capture lead
GET    /api/v1/leads/:id                → Get lead
POST   /api/v1/leads/:id/qualify        → Qualify lead
```

### Messaging Endpoints
```
GET    /api/v1/messaging/templates        → List templates
POST   /api/v1/messaging/templates        → Create template
GET    /api/v1/messaging/templates/:id    → Get template
POST   /api/v1/messaging/send             → Send message
GET    /api/v1/messaging/logs             → Message history
```

### Attachments Endpoints
```
POST   /api/v1/attachments                → Upload attachment
GET    /api/v1/attachments                → List attachments
GET    /api/v1/attachments/:id              → Get attachment
GET    /api/v1/attachments/:id/download     → Download attachment
POST   /api/v1/attachments/:id/sync         → Trigger sync
```

### Admin Endpoints
```
GET    /api/v1/admin/users                  → List users
POST   /api/v1/admin/users                  → Create user
GET    /api/v1/admin/users/:id              → Get user
PUT    /api/v1/admin/users/:id              → Update user
DELETE /api/v1/admin/users/:id              → Delete user
GET    /api/v1/admin/branches               → List branches
POST   /api/v1/admin/branches               → Create branch
GET    /api/v1/admin/branches/:id           → Get branch
PUT    /api/v1/admin/branches/:id           → Update branch
GET    /api/v1/admin/areas                  → List areas
POST   /api/v1/admin/areas                  → Create area
GET    /api/v1/reports/arrears              → Arrears analytics
GET    /api/v1/reports/loans                → Loan portfolio reports
GET    /api/v1/reports/campaigns            → Campaign effectiveness
GET    /api/v1/reports/compliance           → Regulatory reports
GET    /api/v1/audit/logs                   → Audit trail
```

---

## Data Models (Cassandra Schema)

### Key Design Principles
- **Denormalized tables** for read-heavy workloads
- **TimeUUID** for temporal ordering
- **Partition keys** aligned with query patterns
- **Materialized views** for secondary lookups

### Core Tables

```sql
-- Users & Authentication
CREATE TABLE users (
    id uuid PRIMARY KEY,
    email text,
    full_name text,
    role_id text,
    branch_id uuid,
    area_id uuid,
    password_hash text,
    avatar_color text,
    is_active boolean,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE user_sessions (
    user_id uuid,
    session_token text,
    expires_at timestamp,
    created_at timestamp,
    PRIMARY KEY (user_id, session_token)
);

-- Loans Domain
CREATE TABLE loans (
    id uuid PRIMARY KEY,
    client_name text,
    client_phone text,
    client_id text, -- SA ID
    branch_id uuid,
    principal_cents bigint,
    term_months int,
    rate_basis_points int, -- annual % * 100
    status text,
    created_by uuid,
    assessed_by uuid,
    approved_by uuid,
    disbursed_by uuid,
    outstanding_cents bigint,
    next_due_date date,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE loan_applications (
    id uuid PRIMARY KEY,
    loan_id uuid,
    client_name text,
    client_id text,
    monthly_income_cents bigint,
    monthly_expenses_cents bigint,
    requested_amount_cents bigint,
    term_months int,
    branch_id uuid,
    created_by uuid,
    status text,
    affordability_basis_points int, -- ratio * 10000
    assessment_note text,
    approval_note text,
    created_at timestamp,
    updated_at timestamp
);

-- Arrears Domain
CREATE TABLE arrears_cases (
    id uuid PRIMARY KEY,
    loan_id uuid,
    client_name text,
    client_phone text,
    branch_id uuid,
    days_past_due int,
    arrears_cents bigint,
    outstanding_cents bigint,
    status text,
    assigned_to uuid,
    last_action_at timestamp,
    next_touch_at timestamp,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE arrears_ptps (
    id uuid PRIMARY KEY,
    case_id uuid,
    amount_cents bigint,
    promised_at date,
    captured_by uuid,
    note text,
    status text,
    created_at timestamp
);

CREATE TABLE payment_evidences (
    id uuid PRIMARY KEY,
    case_id uuid,
    amount_cents bigint,
    payment_method text,
    reference text,
    captured_by uuid,
    attachment_id uuid,
    captured_at timestamp,
    synced_at timestamp
);

-- Campaigns Domain
CREATE TABLE campaigns (
    id uuid PRIMARY KEY,
    name text,
    branch_id uuid,
    status text,
    start_date date,
    end_date date,
    target_leads int,
    captured_leads int,
    created_by uuid,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE campaign_routes (
    id uuid PRIMARY KEY,
    campaign_id uuid,
    route_date date,
    assigned_to uuid,
    status text,
    planned_by uuid,
    created_at timestamp
);

CREATE TABLE route_stops (
    id uuid PRIMARY KEY,
    route_id uuid,
    address text,
    suburb text,
    status text,
    visited_at timestamp,
    lat double,
    lng double,
    note text,
    lead_id uuid
);

CREATE TABLE leads (
    id uuid PRIMARY KEY,
    campaign_id uuid,
    full_name text,
    phone text,
    suburb text,
    captured_by uuid,
    captured_at timestamp,
    qualified text,
    estimated_amount_cents bigint
);

-- Messaging Domain
CREATE TABLE message_templates (
    id uuid PRIMARY KEY,
    name text,
    channel text,
    context text,
    body text,
    description text,
    created_by uuid,
    created_at timestamp
);

CREATE TABLE message_logs (
    id uuid PRIMARY KEY,
    context text,
    entity_id uuid,
    channel text,
    template_id uuid,
    to_number text,
    body text,
    status text,
    sent_at timestamp,
    delivered_at timestamp,
    failed_reason text,
    sent_by uuid,
    next_touch_at timestamp
);

-- Attachments Domain
CREATE TABLE attachments (
    id uuid PRIMARY KEY,
    context text,
    entity_id uuid,
    file_name text,
    mime_type text,
    size_bytes bigint,
    data_url text,
    captured_by uuid,
    captured_at timestamp,
    sync_status text,
    remote_url text,
    note text
);

-- Audit & Reporting
CREATE TABLE audit_logs (
    id timeuuid PRIMARY KEY,
    user_id uuid,
    action text,
    entity_type text,
    entity_id uuid,
    old_values text,
    new_values text,
    ip_address text,
    user_agent text,
    created_at timestamp
);

-- Materialized Views for Common Queries
CREATE MATERIALIZED VIEW loans_by_branch AS
    SELECT * FROM loans
    WHERE branch_id IS NOT NULL AND id IS NOT NULL
    PRIMARY KEY (branch_id, id);

CREATE MATERIALIZED VIEW arrears_by_collector AS
    SELECT * FROM arrears_cases
    WHERE assigned_to IS NOT NULL AND id IS NOT NULL
    PRIMARY KEY (assigned_to, id);

CREATE MATERIALIZED VIEW leads_by_campaign AS
    SELECT * FROM leads
    WHERE campaign_id IS NOT NULL AND id IS NOT NULL
    PRIMARY KEY (campaign_id, id);
```

---

## Event Schema (Kafka Topics)

### Event Structure (CloudEvents-inspired)
```json
{
  "specversion": "1.0",
  "type": "loans.loan.disbursed",
  "source": "kopesa.loan-service",
  "id": "uuid-v7",
  "time": "2024-01-15T10:30:00Z",
  "datacontenttype": "application/json",
  "data": { ...event payload... },
  "kopesa": {
    "tenant_id": "tenant-uuid",
    "user_id": "user-uuid",
    "correlation_id": "request-uuid",
    "causation_id": "previous-event-uuid"
  }
}
```

### Topic Categories

#### loans.*
- `loans.loan.created`
- `loans.loan.assessed`
- `loans.loan.approved`
- `loans.loan.disbursed`
- `loans.loan.repayment.received`
- `loans.loan.closed`
- `loans.loan.written_off`

#### arrears.*
- `arrears.case.created`
- `arrears.case.allocated`
- `arrears.ptp.created`
- `arrears.ptp.kept`
- `arrears.ptp.broken`
- `arrears.payment.captured`
- `arrears.case.escalated`
- `arrears.case.written_off`
- `arrears.case.resolved`

#### campaigns.*
- `campaigns.campaign.created`
- `campaigns.campaign.activated`
- `campaigns.route.planned`
- `campaigns.route.assigned`
- `campaigns.stop.visited`
- `campaigns.lead.captured`
- `campaigns.lead.qualified`
- `campaigns.lead.converted`
- `campaigns.campaign.completed`

#### iam.*
- `iam.user.created`
- `iam.user.role_changed`
- `iam.user.deactivated`
- `iam.user.login`
- `iam.user.logout`
- `iam.permission.granted`
- `iam.permission.revoked`

#### audit.*
- `audit.action.recorded` (all actions)

---

## Frontend Implementation Details

### Route Structure
```
/                          → Landing page (marketing)
/login                     → Login with role selector
/select-role               → Role selection (subscription-gated)

/app                       → Main dashboard
/arrears                   → Arrears case list
/arrears/:id               → Arrears case detail
/arrears/allocation        → Case allocation to collectors
/arrears/payments          → Payment evidence capture
/arrears/import            → CSV bulk import

/campaigns                 → Campaign list
/campaigns/planner         → Route planning
/campaigns/routes          → Route execution
/campaigns/leads           → Lead management

/loans                    → Loan portfolio
/loans/applications       → Application processing
/loans/disbursements      → Disbursement management

/reports                  → Standard reports
/compliance               → Compliance dashboard
/admin                    → System administration

/dashboards/executive     → Executive KPI dashboard
/dashboards/operations    → Operations dashboard
/help                     → User guide & documentation
```

### State Management (Redux)

#### Slices
- `sessionSlice`: Authentication, current user, role switching
- `loansSlice`: Loan portfolio, applications, disbursements
- `arrearsSlice`: Arrears cases, PTPs, payments, allocations
- `campaignsSlice`: Campaigns, routes, leads
- `messagingSlice`: Templates, message logs
- `attachmentsSlice`: File uploads, sync status
- `uiSlice`: Sidebar state, toasts, loading states

#### Key Patterns
- Normalized state structure
- Async thunks for API calls
- Optimistic updates for UI responsiveness
- Offline queue for pending actions

### Component Architecture

#### Layout Components
- `AppLayout`: Main authenticated layout with sidebar
- `AppSidebar`: Collapsible navigation with role-based items
- `MobileBottomNav`: Mobile-optimized navigation
- `RoleSwitcher`: Quick role switching for demo/testing

#### Common Components
- `PermissionGate`: Route-level permission checking
- `StatusBadge`: Consistent status indicators
- `PageHeader`: Standardized page headers
- `EmptyState`: No-data states

#### Domain Components
- `MessageComposer`: SMS/WhatsApp composition with templates
- `MessageLog`: Communication history
- `AttachmentDropzone`: File upload with offline support

---

## Backend API Specification (Go + Cassandra + Kafka)

### Architecture Pattern: Domain-Driven Design (DDD) + CQRS

#### Project Structure
```
/kopesa-backend
├── /cmd
│   └── /api                 → Main application entry
├── /internal
│   ├── /domain              → Domain models & business logic
│   │   ├── /loans
│   │   ├── /arrears
│   │   ├── /campaigns
│   │   ├── /iam
│   │   ├── /messaging
│   │   └── /common
│   ├── /application         → Use cases / services
│   │   ├── /commands
│   │   └── /queries
│   ├── /infrastructure
│   │   ├── /cassandra       → Cassandra repositories
│   │   ├── /kafka           → Event publishing
│   │   ├── /http            → Fiber handlers
│   │   └── /auth            → JWT middleware
│   └── /interfaces
│       └── /http            → Route definitions
├── /pkg
│   └── /shared              → Shared utilities
└── /migrations              → Cassandra schema migrations
```

### API Standards

#### Request/Response Format
```json
// Standard Response Envelope
{
  "success": true,
  "data": { ... },
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  },
  "error": null
}

// Error Response
{
  "success": false,
  "data": null,
  "meta": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      { "field": "principal", "message": "Must be positive" }
    ]
  }
}
```

#### Pagination
- **Cursor-based**: For high-volume streams (events, audit logs)
- **Offset-based**: For standard lists with page numbers

#### Idempotency
- **Header**: `Idempotency-Key: uuid`
- **TTL**: 24 hours for idempotency keys
- **Scope**: Per-user idempotency

### Authentication & Authorization

#### JWT Token Structure
```json
{
  "sub": "user-uuid",
  "email": "user@kopesa.co.za",
  "role": "branch_manager",
  "scope": "branch",
  "branch_id": "br-jhb-cbd",
  "area_id": "area-gp",
  "permissions": ["arrears.view", "arrears.allocate", ...],
  "iat": 1705315200,
  "exp": 1705318800
}
```

#### Permission Middleware
```go
// Example: Require permission
router.Get("/arrears", requirePermission("arrears.view"), handler)

// Example: Require any of multiple permissions
router.Post("/arrears/:id/payment", 
    requireAnyPermission("arrears.payment.capture", "finance.reconcile"), 
    handler)
```

### Key API Endpoints by Domain

#### IAM Domain
```
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
GET    /api/v1/auth/me
GET    /api/v1/auth/roles               → Available roles for user
PUT    /api/v1/auth/role                 → Switch active role

GET    /api/v1/admin/users
POST   /api/v1/admin/users
GET    /api/v1/admin/users/:id
PUT    /api/v1/admin/users/:id
DELETE /api/v1/admin/users/:id

GET    /api/v1/admin/branches
POST   /api/v1/admin/branches
GET    /api/v1/admin/branches/:id
PUT    /api/v1/admin/branches/:id

GET    /api/v1/admin/areas
POST   /api/v1/admin/areas
```

#### Loans Domain
```
GET    /api/v1/loans
POST   /api/v1/loans
GET    /api/v1/loans/:id
GET    /api/v1/loans/:id/history
POST   /api/v1/loans/:id/assess
POST   /api/v1/loans/:id/approve
POST   /api/v1/loans/:id/decline
POST   /api/v1/loans/:id/disburse

GET    /api/v1/loans/applications
POST   /api/v1/loans/applications
GET    /api/v1/loans/applications/:id
PUT    /api/v1/loans/applications/:id

GET    /api/v1/loans/disbursements
POST   /api/v1/loans/disbursements
GET    /api/v1/loans/disbursements/:id
```

#### Arrears Domain
```
GET    /api/v1/arrears
POST   /api/v1/arrears
POST   /api/v1/arrears/import            → CSV bulk import
GET    /api/v1/arrears/:id
GET    /api/v1/arrears/:id/history
GET    /api/v1/arrears/:id/ptps
POST   /api/v1/arrears/:id/allocate       → Assign to collector
POST   /api/v1/arrears/:id/action        → Record action
POST   /api/v1/arrears/:id/ptp           → Create PTP
POST   /api/v1/arrears/:id/payment       → Record payment
POST   /api/v1/arrears/:id/escalate      → Escalate case
POST   /api/v1/arrears/:id/writeoff      → Write off case
POST   /api/v1/arrears/:id/cancel        → Cancel case

GET    /api/v1/arrears/allocations        → Allocation history
GET    /api/v1/arrears/payments           → Payment evidence list
```

#### Campaigns Domain
```
GET    /api/v1/campaigns
POST   /api/v1/campaigns
GET    /api/v1/campaigns/:id
PUT    /api/v1/campaigns/:id
POST   /api/v1/campaigns/:id/activate
POST   /api/v1/campaigns/:id/complete
POST   /api/v1/campaigns/:id/cancel

GET    /api/v1/campaigns/:id/routes
POST   /api/v1/campaigns/:id/routes
GET    /api/v1/routes/:id
PUT    /api/v1/routes/:id
POST   /api/v1/routes/:id/assign

GET    /api/v1/routes/:id/stops
POST   /api/v1/routes/:id/stops
PUT    /api/v1/stops/:id
POST   /api/v1/stops/:id/visit
POST   /api/v1/stops/:id/skip

GET    /api/v1/leads
POST   /api/v1/leads
GET    /api/v1/leads/:id
PUT    /api/v1/leads/:id
POST   /api/v1/leads/:id/qualify
POST   /api/v1/leads/:id/reject
POST   /api/v1/leads/:id/convert
```

#### Messaging Domain
```
GET    /api/v1/messaging/templates
POST   /api/v1/messaging/templates
GET    /api/v1/messaging/templates/:id
PUT    /api/v1/messaging/templates/:id
DELETE /api/v1/messaging/templates/:id

POST   /api/v1/messaging/send
POST   /api/v1/messaging/send-bulk
POST   /api/v1/messaging/schedule

GET    /api/v1/messaging/logs
GET    /api/v1/messaging/logs/:id
GET    /api/v1/messaging/stats

POST   /api/v1/messaging/webhooks/sms
POST   /api/v1/messaging/webhooks/whatsapp
```

#### Attachments Domain
```
POST   /api/v1/attachments                    → Upload (multipart)
GET    /api/v1/attachments
GET    /api/v1/attachments/:id
GET    /api/v1/attachments/:id/download
PUT    /api/v1/attachments/:id
DELETE /api/v1/attachments/:id

POST   /api/v1/attachments/:id/sync           → Trigger sync
GET    /api/v1/attachments/pending            → Pending sync list
POST   /api/v1/attachments/bulk-sync          → Bulk sync

GET    /api/v1/attachments/context/:context/:entity_id
```

#### Reporting Domain
```
GET    /api/v1/reports/arrears/summary
GET    /api/v1/reports/arrears/aging
GET    /api/v1/reports/arrears/collector-performance
GET    /api/v1/reports/arrears/ptp-analysis

GET    /api/v1/reports/loans/portfolio
GET    /api/v1/reports/loans/disbursements
GET    /api/v1/reports/loans/applications
GET    /api/v1/reports/loans/performance

GET    /api/v1/reports/campaigns/summary
GET    /api/v1/reports/campaigns/conversion
GET    /api/v1/reports/campaigns/routes

GET    /api/v1/reports/compliance/kyc
GET    /api/v1/reports/compliance/nca
GET    /api/v1/reports/compliance/audit

GET    /api/v1/reports/executive/kpis
GET    /api/v1/reports/executive/trends
```

#### Audit Domain
```
GET    /api/v1/audit/logs
GET    /api/v1/audit/logs/:id
GET    /api/v1/audit/logs/user/:user_id
GET    /api/v1/audit/logs/entity/:entity_type/:entity_id
GET    /api/v1/audit/stats
GET    /api/v1/audit/export
```

---

## Error Handling

### HTTP Status Codes
| Code | Usage |
|------|-------|
| 200 | Success |
| 201 | Created |
| 204 | No Content (deletions) |
| 400 | Bad Request (validation) |
| 401 | Unauthorized |
| 403 | Forbidden (permission) |
| 404 | Not Found |
| 409 | Conflict (duplicate, state) |
| 422 | Unprocessable (business rule) |
| 429 | Rate Limited |
| 500 | Internal Server Error |

### Error Response Format
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Request validation failed",
    "details": [
      { "field": "principal", "message": "Must be greater than 0" },
      { "field": "term_months", "message": "Must be between 1 and 60" }
    ]
  }
}
```

### Error Codes by Domain
```go
// Common errors
ErrInvalidInput       = "INVALID_INPUT"
ErrUnauthorized       = "UNAUTHORIZED"
ErrForbidden          = "FORBIDDEN"
ErrNotFound           = "NOT_FOUND"
ErrConflict           = "CONFLICT"
ErrValidation         = "VALIDATION_ERROR"
ErrInternal           = "INTERNAL_ERROR"

// Domain-specific errors
ErrLoanNotInState        = "LOAN_NOT_IN_STATE"
ErrInsufficientFunds     = "INSUFFICIENT_FUNDS"
ErrArrearsAlreadyAllocated = "ARREARS_ALREADY_ALLOCATED"
ErrPTPBroken            = "PTP_BROKEN"
ErrCampaignNotActive    = "CAMPAIGN_NOT_ACTIVE"
ErrRouteNotAssigned     = "ROUTE_NOT_ASSIGNED"
ErrMakerCheckerConflict = "MAKER_CHECKER_CONFLICT"
```

---

## Security Requirements

### Authentication
- JWT access tokens (15 min expiry)
- Refresh tokens (7 days, rotation)
- Secure httpOnly cookies
- CSRF protection for state-changing ops

### Authorization
- RBAC with 30+ granular permissions
- Data scope enforcement (field/branch/area/global)
- Maker-checker for financial mutations
- Subscription-gated role selection

### Data Protection
- Encryption at rest (AES-256)
- TLS 1.3 in transit
- PII masking in logs
- GDPR-compliant data retention

### Audit
- Immutable audit log (Cassandra time-series)
- Event sourcing for financial transactions
- Tamper-evident log chaining

---

## Infrastructure & Deployment

### Container Strategy
```dockerfile
# Multi-stage build for Go services
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]
```

### Kubernetes Resources
```yaml
# Deployment with HPA
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kopesa-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: kopesa-api
  template:
    metadata:
      labels:
        app: kopesa-api
    spec:
      containers:
      - name: api
        image: kopesa/api:latest
        ports:
        - containerPort: 8080
        env:
        - name: CASSANDRA_HOSTS
          valueFrom:
            configMapKeyRef:
              name: kopesa-config
              key: cassandra-hosts
        - name: KAFKA_BROKERS
          valueFrom:
            configMapKeyRef:
              name: kopesa-config
              key: kafka-brokers
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: kopesa-api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: kopesa-api
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### Environment Configuration
```yaml
# config.yaml
app:
  name: kopesa-api
  version: 1.0.0
  environment: production
  port: 8080

cassandra:
  hosts:
    - cassandra-1:9042
    - cassandra-2:9042
    - cassandra-3:9042
  keyspace: kopesa
  consistency: LOCAL_QUORUM
  timeout: 5s

kafka:
  brokers:
    - kafka-1:9092
    - kafka-2:9092
    - kafka-3:9092
  topics_prefix: kopesa
  consumer_group: kopesa-api

auth:
  jwt_secret: ${JWT_SECRET}
  access_token_ttl: 15m
  refresh_token_ttl: 7d
  issuer: kopesa

security:
  rate_limit_requests: 100
  rate_limit_window: 1m
  cors_origins:
    - https://loan.biacibenga.com
    - https://field-lend-flow.lovable.app
```

---

## Testing Strategy

### Unit Tests
```go
// Example: Loan domain test
func TestLoan_CanDisburse(t *testing.T) {
    tests := []struct {
        name    string
        status  LoanStatus
        want    bool
        wantErr error
    }{
        {"approved can disburse", StatusApproved, true, nil},
        {"active cannot disburse", StatusActive, false, ErrLoanNotInState},
        {"application cannot disburse", StatusApplication, false, ErrLoanNotInState},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            loan := &Loan{Status: tt.status}
            got, err := loan.CanDisburse()
            assert.Equal(t, tt.want, got)
            assert.ErrorIs(t, err, tt.wantErr)
        })
    }
}
```

### Integration Tests
```go
// API integration test
func TestArrearsAPI_CreateCase(t *testing.T) {
    // Setup test container for Cassandra
    cassandra := testcontainers.NewCassandraContainer()
    defer cassandra.Terminate()
    
    // Setup test Kafka
    kafka := testcontainers.NewKafkaContainer()
    defer kafka.Terminate()
    
    // Run API tests
    app := setupTestApp(cassandra, kafka)
    
    t.Run("create case success", func(t *testing.T) {
        req := httptest.NewRequest("POST", "/api/v1/arrears", 
            strings.NewReader(`{"loan_id": "...", "days_past_due": 30}`))
        resp, err := app.Test(req)
        assert.NoError(t, err)
        assert.Equal(t, 201, resp.StatusCode)
    })
}
```

### Contract Tests (Pact)
```json
{
  "consumer": {
    "name": "kopesa-web"
  },
  "provider": {
    "name": "kopesa-api"
  },
  "interactions": [
    {
      "description": "get arrears list",
      "request": {
        "method": "GET",
        "path": "/api/v1/arrears",
        "headers": {
          "Authorization": "Bearer token"
        }
      },
      "response": {
        "status": 200,
        "body": {
          "success": true,
          "data": [{
            "id": "uuid",
            "client_name": "string",
            "days_past_due": 30
          }]
        }
      }
    }
  ]
}
```

---

## Deployment & Operations

### CI/CD Pipeline
```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Build Docker image
        run: docker build -t kopesa/api:${{ github.sha }} .
      
      - name: Push to registry
        run: |
          echo ${{ secrets.DOCKER_PASSWORD }} | docker login -u ${{ secrets.DOCKER_USERNAME }} --password-stdin
          docker push kopesa/api:${{ github.sha }}

  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - name: Deploy to Kubernetes
        run: |
          kubectl set image deployment/kopesa-api api=kopesa/api:${{ github.sha }}
          kubectl rollout status deployment/kopesa-api
```

### Monitoring & Observability

#### Metrics (Prometheus)
```go
// Custom metrics
var (
    loanDisbursements = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "kopesa_loan_disbursements_total",
            Help: "Total loan disbursements",
        },
        []string{"branch_id", "currency"},
    )
    
    arrearsResolutionTime = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "kopesa_arrears_resolution_duration_days",
            Help:    "Time to resolve arrears case",
            Buckets: []float64{1, 7, 14, 30, 60, 90},
        },
        []string{"dpd_bucket"},
    )
)
```

#### Distributed Tracing (OpenTelemetry)
```go
// Trace context propagation
ctx, span := tracer.Start(ctx, "ProcessLoanApplication")
defer span.End()

span.SetAttributes(
    attribute.String("loan.id", loanID),
    attribute.String("applicant.id", applicantID),
    attribute.Int64("requested.amount", amountCents),
)

// Child spans for database operations
ctx, dbSpan := tracer.Start(ctx, "Cassandra.Insert")
err := session.Query("INSERT INTO loans ...").WithContext(ctx).Exec()
dbSpan.End()
```

#### Logging (Structured)
```go
// Zap logger configuration
logger, _ := zap.NewProduction()
defer logger.Sync()

// Contextual logging
logger.Info("loan_disbursed",
    zap.String("loan_id", loan.ID.String()),
    zap.String("branch_id", loan.BranchID.String()),
    zap.Int64("amount_cents", loan.PrincipalCents),
    zap.String("disbursed_by", userID.String()),
    zap.String("trace_id", traceID),
)
```

### Alerting Rules
```yaml
# prometheus-alerts.yml
groups:
  - name: kopesa-critical
    rules:
      - alert: HighErrorRate
        expr: rate(kopesa_http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          
      - alert: DatabaseLatency
        expr: histogram_quantile(0.95, rate(kopesa_cassandra_query_duration_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Cassandra query latency high"
          
      - alert: LoanDisbursementStuck
        expr: time() - kopesa_loan_approval_timestamp > 3600
        for: 10m
        labels:
          severity: critical
        annotations:
          summary: "Approved loan not disbursed for over 1 hour"
```

---

## Development Guidelines

### Code Organization
```
/internal
  /domain          → Pure business logic, no external deps
    /loans
      aggregate.go   → Loan entity, invariants
      events.go      → Domain events
      repository.go  → Repository interface
    /arrears
    /campaigns
    ...
  
  /application     → Use cases, orchestration
    /commands
      create_loan.go
      disburse_loan.go
    /queries
      get_loan.go
      list_loans.go
    /handlers
      event_handlers.go
  
  /infrastructure  → External implementations
    /cassandra
      repositories.go
      connection.go
    /kafka
      producer.go
      consumer.go
    /http
      handlers.go
      middleware.go
      routes.go
    /auth
      jwt.go
      middleware.go
  
  /interfaces      → Input adapters
    /http
      dto.go
      mappers.go
      validators.go
```

### Testing Strategy
1. **Unit Tests**: Domain logic, pure functions
2. **Integration Tests**: Repository implementations with test containers
3. **Contract Tests**: Pact for frontend/backend compatibility
4. **E2E Tests**: Cypress for critical user journeys

### Git Workflow
```
main
  ↓
feature/loans-disbursement
  ↓
PR → Code Review → CI Pass → Merge
```

### Commit Convention
```
feat(loans): add disbursement endpoint
fix(arrears): handle null assigned_to in allocation
docs(api): update arrears import specification
refactor(domain): extract common aggregate logic
test(campaigns): add route planning integration tests
chore(deps): update fiber to v2.50
```

---

## Appendix A: Glossary

| Term | Definition |
|------|------------|
| **Arrears** | Loan payments that are overdue |
| **DPD** | Days Past Due - measure of delinquency |
| **PTP** | Promise to Pay - commitment from borrower |
| **POP** | Proof of Payment - evidence of payment |
| **KYC** | Know Your Customer - identity verification |
| **NCA** | National Credit Act (South Africa) |
| **Maker-Checker** | Dual approval process for critical operations |
| **Denormalized** | Data duplication for query performance |
| **CQRS** | Command Query Responsibility Segregation |

---

## Appendix B: South African Regulatory Considerations

### National Credit Act (NCA) Compliance
1. **Affordability Assessment**: Must verify income vs expenses
2. **Interest Rate Caps**: Regulated maximum rates
3. **Cooling-off Period**: 5 days for certain products
4. **Dispute Resolution**: Formal process required

### POPIA (Data Protection)
1. **Consent Management**: Explicit consent for marketing
2. **Data Minimization**: Only collect necessary data
3. **Retention Limits**: Delete after statutory period
4. **Breach Notification**: 72-hour reporting requirement

### Financial Intelligence Centre Act (FICA)
1. **Customer Identification**: Verify identity documents
2. **Record Keeping**: 5-year retention
3. **Suspicious Transaction Reporting**: STR filing
4. **Ongoing Due Diligence**: Periodic verification

---

## Document Version History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2024-01-15 | System | Initial specification |
| 1.1 | 2024-01-20 | System | Added Cassandra schema details |
| 1.2 | 2024-01-25 | System | Added Kafka event specifications |

---

**Document Owner**: Kopesa Development Team
**Review Cycle**: Monthly
**Next Review**: 2024-02-15
