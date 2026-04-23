CREATE TABLE IF NOT EXISTS iam_areas (
    id text PRIMARY KEY,
    name text
);

CREATE TABLE IF NOT EXISTS iam_branches (
    id text PRIMARY KEY,
    code text,
    name text,
    area_id text
);

CREATE TABLE IF NOT EXISTS iam_roles (
    id text PRIMARY KEY,
    label text,
    scope text,
    permissions list<text>
);

CREATE TABLE IF NOT EXISTS iam_users (
    id text PRIMARY KEY,
    full_name text,
    email text,
    password_hash text,
    role text,
    allowed_roles list<text>,
    branch_id text,
    area_id text,
    avatar_color text,
    active boolean
);

CREATE TABLE IF NOT EXISTS iam_users_by_email (
    email text PRIMARY KEY,
    id text,
    full_name text,
    password_hash text,
    role text,
    allowed_roles list<text>,
    branch_id text,
    area_id text,
    avatar_color text,
    active boolean
);

CREATE TABLE IF NOT EXISTS loans_applications (
    id text PRIMARY KEY,
    loan_id text,
    client_name text,
    client_id text,
    client_phone text,
    monthly_income bigint,
    monthly_expenses bigint,
    requested_amount bigint,
    term_months int,
    branch_id text,
    created_by text,
    status text,
    affordability bigint,
    assessment_note text,
    approval_note text
);

CREATE TABLE IF NOT EXISTS loans_loans (
    id text PRIMARY KEY,
    reference text,
    client_name text,
    client_phone text,
    branch_id text,
    principal bigint,
    term_months int,
    rate_bp int,
    status text,
    outstanding_balance bigint,
    next_due_date text,
    created_by text,
    assessed_by text,
    approved_by text,
    disbursed_by text,
    created_at text,
    updated_at text
);

CREATE TABLE IF NOT EXISTS arrears_cases (
    id text PRIMARY KEY,
    loan_id text,
    client_name text,
    client_phone text,
    branch_id text,
    days_past_due int,
    arrears_amount bigint,
    outstanding_balance bigint,
    status text,
    assigned_to text,
    last_action_at text,
    ptps_json text
);

CREATE TABLE IF NOT EXISTS arrears_payments (
    id text PRIMARY KEY,
    case_id text,
    amount bigint,
    method text,
    reference text,
    captured_by text,
    captured_at text,
    attachment_id text
);

CREATE TABLE IF NOT EXISTS campaigns_campaigns (
    id text PRIMARY KEY,
    name text,
    branch_id text,
    status text,
    start_date text,
    end_date text,
    target_leads int,
    captured_leads int
);

CREATE TABLE IF NOT EXISTS campaigns_routes (
    campaign_id text,
    id text,
    date text,
    assigned_to text,
    status text,
    stops_json text,
    PRIMARY KEY ((campaign_id), id)
);

CREATE TABLE IF NOT EXISTS campaigns_leads (
    id text PRIMARY KEY,
    campaign_id text,
    full_name text,
    phone text,
    suburb text,
    captured_by text,
    captured_at text,
    qualified text,
    estimated_amount bigint
);

CREATE TABLE IF NOT EXISTS messaging_templates (
    id text PRIMARY KEY,
    name text,
    channel text,
    context text,
    body text,
    description text
);

CREATE TABLE IF NOT EXISTS messaging_logs (
    id text PRIMARY KEY,
    context text,
    entity_id text,
    channel text,
    template_id text,
    recipient_to text,
    body text,
    status text,
    sent_at text,
    sent_by text,
    next_touch_at text,
    provider_ref text
);

CREATE TABLE IF NOT EXISTS attachments_attachments (
    id text PRIMARY KEY,
    context text,
    entity_id text,
    file_name text,
    mime_type text,
    size_bytes bigint,
    url text,
    captured_by text,
    captured_at text,
    sync text,
    note text,
    provider text,
    path text,
    revision text
);

CREATE TABLE IF NOT EXISTS audit_events (
    id text PRIMARY KEY,
    occurred_at text,
    actor text,
    actor_role text,
    action text,
    entity_type text,
    entity_id text,
    ip text,
    user_agent text,
    request_id text,
    metadata map<text, text>
);
