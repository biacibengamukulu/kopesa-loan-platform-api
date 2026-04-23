package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"

	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/cassandra"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/config"
)

func main() {
	cfg := config.Load()
	if len(cfg.CassandraHosts) == 1 && cfg.CassandraHosts[0] == "127.0.0.1" {
		cfg.CassandraHosts = []string{"safer.easipath.com"}
	}
	if cfg.CassandraKeyspace == "" {
		cfg.CassandraKeyspace = "kopesa_loan_platform"
	}

	adminSession, err := cassandra.NewSession(cfg, false)
	if err != nil {
		log.Fatalf("connect cassandra admin session: %v", err)
	}
	defer adminSession.Close()

	createKS := fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`, cfg.CassandraKeyspace)
	if err := adminSession.Query(createKS).Exec(); err != nil {
		log.Fatalf("create keyspace: %v", err)
	}

	session, err := cassandra.NewSession(cfg, true)
	if err != nil {
		log.Fatalf("connect cassandra keyspace session: %v", err)
	}
	defer session.Close()

	schema, err := os.ReadFile("migrations/001_init.cql")
	if err != nil {
		log.Fatalf("read schema: %v", err)
	}

	if err := applySchema(session, string(schema)); err != nil {
		log.Fatalf("apply schema: %v", err)
	}
	if err := applyCompatibilityAlterations(session); err != nil {
		log.Fatalf("apply compatibility alterations: %v", err)
	}
	if err := seed(session); err != nil {
		log.Fatalf("seed data: %v", err)
	}

	log.Printf("migration complete on %s keyspace=%s", strings.Join(cfg.CassandraHosts, ","), cfg.CassandraKeyspace)
}

func applySchema(session *gocql.Session, schema string) error {
	stmts := strings.Split(schema, ";")
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if err := session.Query(stmt).Exec(); err != nil {
			return err
		}
	}
	return nil
}

func applyCompatibilityAlterations(session *gocql.Session) error {
	stmts := []string{
		`ALTER TABLE attachments_attachments ADD provider text`,
		`ALTER TABLE attachments_attachments ADD path text`,
		`ALTER TABLE attachments_attachments ADD revision text`,
		`ALTER TABLE messaging_logs ADD provider_ref text`,
	}
	for _, stmt := range stmts {
		if err := session.Query(stmt).Exec(); err != nil && !isIgnorableAlterError(err) {
			return fmt.Errorf("%s: %w", stmt, err)
		}
	}
	return nil
}

func isIgnorableAlterError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "already exists") || strings.Contains(msg, "conflicts with an existing column")
}

func seed(session *gocql.Session) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("Kopesa@2026"), 12)
	if err != nil {
		return err
	}
	mustExec(session, `INSERT INTO iam_areas (id, name) VALUES (?, ?)`, "area-gp", "Gauteng")
	mustExec(session, `INSERT INTO iam_areas (id, name) VALUES (?, ?)`, "area-wc", "Western Cape")
	mustExec(session, `INSERT INTO iam_areas (id, name) VALUES (?, ?)`, "area-kzn", "KwaZulu-Natal")

	branches := [][]any{
		{"br-jhb-cbd", "JHB01", "Johannesburg CBD", "area-gp"},
		{"br-soweto", "JHB02", "Soweto", "area-gp"},
		{"br-cpt-cbd", "CPT01", "Cape Town CBD", "area-wc"},
		{"br-dbn-cbd", "DBN01", "Durban CBD", "area-kzn"},
	}
	for _, row := range branches {
		mustExec(session, `INSERT INTO iam_branches (id, code, name, area_id) VALUES (?, ?, ?, ?)`, row...)
	}

	type roleSeed struct {
		ID          string
		Label       string
		Scope       string
		Permissions []string
	}
	roles := []roleSeed{
		{"field_officer", "Field Officer", "field", []string{"campaign.view", "campaign.plan", "campaign.lead.capture", "campaign.lead.qualify", "loan.view", "messaging.send", "messaging.view", "attachments.upload", "attachments.view"}},
		{"collector", "Collector", "field", []string{"arrears.view", "arrears.action", "arrears.ptp.create", "arrears.payment.capture", "loan.view", "messaging.send", "messaging.view", "attachments.upload", "attachments.view"}},
		{"branch_agent", "Branch Agent", "branch", []string{"loan.view", "loan.application.create", "arrears.view", "arrears.action", "arrears.ptp.create", "campaign.view", "messaging.send", "messaging.view", "attachments.upload", "attachments.view"}},
		{"consultant", "Loan Consultant", "branch", []string{"loan.view", "loan.application.create", "loan.assess", "maker", "messaging.send", "messaging.view", "attachments.upload", "attachments.view"}},
		{"branch_manager", "Branch Manager", "branch", []string{"arrears.view", "arrears.allocate", "arrears.escalate", "campaign.view", "campaign.plan", "campaign.route.assign", "loan.view", "loan.assess", "loan.approve", "checker", "reports.exec", "messaging.send", "messaging.view", "attachments.view"}},
		{"area_manager", "Area Manager", "area", []string{"arrears.view", "arrears.allocate", "arrears.escalate", "arrears.writeoff", "campaign.view", "campaign.plan", "loan.view", "loan.approve", "checker", "reports.exec", "messaging.view", "attachments.view"}},
		{"finance_officer", "Finance Officer", "global", []string{"arrears.view", "arrears.payment.capture", "arrears.writeoff", "loan.view", "loan.disburse", "loan.servicing", "finance.reconcile", "checker", "attachments.view", "attachments.upload", "messaging.view"}},
		{"compliance_officer", "Compliance Officer", "global", []string{"loan.view", "loan.assess", "arrears.view", "reports.exec", "checker", "attachments.view", "attachments.upload", "messaging.view"}},
		{"admin", "System Admin", "global", []string{"admin.users", "admin.config", "reports.exec", "arrears.view", "campaign.view", "loan.view", "messaging.view", "attachments.view"}},
		{"executive", "Executive", "global", []string{"arrears.view", "campaign.view", "loan.view", "reports.exec", "messaging.view", "attachments.view"}},
	}
	for _, role := range roles {
		mustExec(session, `INSERT INTO iam_roles (id, label, scope, permissions) VALUES (?, ?, ?, ?)`, role.ID, role.Label, role.Scope, role.Permissions)
	}

	type userSeed struct {
		ID, FullName, Email, Role, Avatar string
		Allowed                           []string
		BranchID, AreaID                  *string
	}
	brSoweto := "br-soweto"
	brJHB := "br-jhb-cbd"
	brCPT := "br-cpt-cbd"
	areaGP := "area-gp"
	areaWC := "area-wc"
	areaKZN := "area-kzn"
	users := []userSeed{
		{"u-1", "Thabo Mokoena", "thabo@kopesa.co.za", "field_officer", "hsl(38 80% 43%)", []string{"field_officer"}, &brSoweto, &areaGP},
		{"u-2", "Naledi Dlamini", "naledi@kopesa.co.za", "collector", "hsl(215 70% 44%)", []string{"collector"}, &brJHB, &areaGP},
		{"u-3", "Sipho Khumalo", "sipho@kopesa.co.za", "branch_agent", "hsl(115 48% 38%)", []string{"branch_agent"}, &brJHB, &areaGP},
		{"u-4", "Lerato Nkosi", "lerato@kopesa.co.za", "consultant", "hsl(347 63% 50%)", []string{"consultant"}, &brCPT, &areaWC},
		{"u-5", "Pieter van Wyk", "pieter@kopesa.co.za", "branch_manager", "hsl(28 62% 48%)", []string{"branch_manager"}, &brCPT, &areaWC},
		{"u-6", "Zanele Mthembu", "zanele@kopesa.co.za", "area_manager", "hsl(262 45% 48%)", []string{"area_manager"}, nil, &areaKZN},
		{"u-7", "Hendrik Botha", "hendrik@kopesa.co.za", "finance_officer", "hsl(202 65% 48%)", []string{"finance_officer"}, nil, nil},
		{"u-8", "Aisha Patel", "aisha@kopesa.co.za", "compliance_officer", "hsl(12 74% 52%)", []string{"compliance_officer"}, nil, nil},
		{"u-9", "System Admin", "admin@kopesa.co.za", "admin", "hsl(220 13% 25%)", []string{"admin"}, nil, nil},
		{"u-10", "Mervin Biangacila", "exec@kopesa.co.za", "executive", "hsl(40 82% 46%)", []string{"executive"}, nil, nil},
	}
	for _, user := range users {
		mustExec(session, `INSERT INTO iam_users (id, full_name, email, password_hash, role, allowed_roles, branch_id, area_id, avatar_color, active) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			user.ID, user.FullName, user.Email, string(passwordHash), user.Role, user.Allowed, user.BranchID, user.AreaID, user.Avatar, true)
		mustExec(session, `INSERT INTO iam_users_by_email (email, id, full_name, password_hash, role, allowed_roles, branch_id, area_id, avatar_color, active) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			user.Email, user.ID, user.FullName, string(passwordHash), user.Role, user.Allowed, user.BranchID, user.AreaID, user.Avatar, true)
	}

	now := time.Now().UTC()
	tomorrow := now.Add(24 * time.Hour).Format(time.RFC3339)
	yesterday := now.Add(-24 * time.Hour).Format(time.RFC3339)
	threeDays := now.Add(72 * time.Hour).Format(time.RFC3339)
	oneHour := now.Add(-1 * time.Hour).Format(time.RFC3339)
	twoHour := now.Add(-2 * time.Hour).Format(time.RFC3339)

	mustExec(session, `INSERT INTO loans_loans (id, reference, client_name, client_phone, branch_id, principal, term_months, rate_bp, status, outstanding_balance, next_due_date, created_by, assessed_by, approved_by, disbursed_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"loan-1", "LN-44021", "Bongani Sithole", "+27821234567", "br-jhb-cbd", int64(2500000), 12, 2800, "active", int64(1890000), "2025-11-30", "u-3", nil, nil, nil, now.Format(time.RFC3339), now.Format(time.RFC3339))
	mustExec(session, `INSERT INTO loans_loans (id, reference, client_name, client_phone, branch_id, principal, term_months, rate_bp, status, outstanding_balance, next_due_date, created_by, assessed_by, approved_by, disbursed_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"loan-2", "LN-44078", "Nomvula Khoza", "+27735550102", "br-jhb-cbd", int64(5000000), 18, 2600, "active", int64(4120000), "2025-11-15", "u-3", nil, nil, nil, now.Format(time.RFC3339), now.Format(time.RFC3339))
	mustExec(session, `INSERT INTO loans_loans (id, reference, client_name, client_phone, branch_id, principal, term_months, rate_bp, status, outstanding_balance, next_due_date, created_by, assessed_by, approved_by, disbursed_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"loan-3", "LN-44102", "Tshepo Mabaso", "+27842220190", "br-soweto", int64(3500000), 12, 2800, "active", int64(2870000), "2025-11-05", "u-3", nil, nil, nil, now.Format(time.RFC3339), now.Format(time.RFC3339))
	mustExec(session, `INSERT INTO loans_loans (id, reference, client_name, client_phone, branch_id, principal, term_months, rate_bp, status, outstanding_balance, next_due_date, created_by, assessed_by, approved_by, disbursed_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"loan-4", "LN-44900", "Mpho Radebe", "+27821001100", "br-soweto", int64(800000), 6, 3000, "approved", int64(800000), nil, "u-4", strp("u-4"), strp("u-5"), nil, now.Format(time.RFC3339), now.Format(time.RFC3339))

	mustExec(session, `INSERT INTO loans_applications (id, loan_id, client_name, client_id, client_phone, monthly_income, monthly_expenses, requested_amount, term_months, branch_id, created_by, status, affordability, assessment_note, approval_note) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"app-1", "LN-44900", "Mpho Radebe", "8806125555084", "+27821001100", int64(1200000), int64(750000), int64(800000), 6, "br-soweto", "u-4", "approved", int64(450000), "Stable employment 3yrs", "Approved within branch limit")
	mustExec(session, `INSERT INTO loans_applications (id, loan_id, client_name, client_id, client_phone, monthly_income, monthly_expenses, requested_amount, term_months, branch_id, created_by, status, affordability, assessment_note, approval_note) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"app-2", "LN-44901", "Lerato Maseko", "9203054444081", "+27732209981", int64(850000), int64(520000), int64(500000), 6, "br-soweto", "u-4", "submitted", nil, nil, nil)

	ptps, _ := json.Marshal([]map[string]any{{"id": "ptp-1", "amount": 250000, "promisedAt": threeDays, "capturedBy": "u-2", "status": "pending", "note": "Will pay after payday"}})
	mustExec(session, `INSERT INTO arrears_cases (id, loan_id, client_name, client_phone, branch_id, days_past_due, arrears_amount, outstanding_balance, status, assigned_to, last_action_at, ptps_json) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"ar-1001", "LN-44021", "Bongani Sithole", "+27821234567", "br-jhb-cbd", 12, int64(245000), int64(1890000), "allocated", "u-2", yesterday, "[]")
	mustExec(session, `INSERT INTO arrears_cases (id, loan_id, client_name, client_phone, branch_id, days_past_due, arrears_amount, outstanding_balance, status, assigned_to, last_action_at, ptps_json) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"ar-1002", "LN-44078", "Nomvula Khoza", "+27735550102", "br-jhb-cbd", 34, int64(580000), int64(4120000), "ptp", "u-2", yesterday, string(ptps))
	mustExec(session, `INSERT INTO arrears_cases (id, loan_id, client_name, client_phone, branch_id, days_past_due, arrears_amount, outstanding_balance, status, assigned_to, last_action_at, ptps_json) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"ar-1003", "LN-44102", "Tshepo Mabaso", "+27842220190", "br-soweto", 67, int64(990000), int64(2870000), "broken_ptp", "u-2", yesterday, "[]")

	mustExec(session, `INSERT INTO campaigns_campaigns (id, name, branch_id, status, start_date, end_date, target_leads, captured_leads) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		"cmp-1", "Soweto Spring Drive", "br-soweto", "active", "2025-09-01", "2025-12-15", 200, 84)
	mustExec(session, `INSERT INTO campaigns_campaigns (id, name, branch_id, status, start_date, end_date, target_leads, captured_leads) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		"cmp-2", "JHB CBD Walk-in Push", "br-jhb-cbd", "active", "2025-10-01", "2025-11-30", 150, 41)

	stops, _ := json.Marshal([]map[string]any{
		{"id": "s-1", "address": "12 Vilakazi St", "suburb": "Orlando West", "status": "visited", "visitedAt": now.Format(time.RFC3339)},
		{"id": "s-2", "address": "88 Klipspruit Valley Rd", "suburb": "Klipspruit", "status": "lead", "visitedAt": now.Format(time.RFC3339), "leadId": "ld-1"},
		{"id": "s-3", "address": "4 Mofolo North Ext", "suburb": "Mofolo", "status": "no_contact"},
	})
	mustExec(session, `INSERT INTO campaigns_routes (campaign_id, id, date, assigned_to, status, stops_json) VALUES (?, ?, ?, ?, ?, ?)`,
		"cmp-1", "rt-1", now.Format("2006-01-02"), "u-1", "in_progress", string(stops))

	mustExec(session, `INSERT INTO campaigns_leads (id, campaign_id, full_name, phone, suburb, captured_by, captured_at, qualified, estimated_amount) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"ld-1", "cmp-1", "Mpho Radebe", "+27821001100", "Orlando West", "u-1", oneHour, "qualified", int64(800000))
	mustExec(session, `INSERT INTO campaigns_leads (id, campaign_id, full_name, phone, suburb, captured_by, captured_at, qualified, estimated_amount) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"ld-2", "cmp-1", "Lerato Maseko", "+27732209981", "Klipspruit", "u-1", twoHour, "pending", int64(500000))

	mustExec(session, `INSERT INTO messaging_templates (id, name, channel, context, body, description) VALUES (?, ?, ?, ?, ?, ?)`,
		"tpl-ar-reminder", "Friendly reminder", "both", "arrears", "Hello {{client}}, your account is overdue by {{daysPastDue}} days.", "Arrears reminder")
	mustExec(session, `INSERT INTO messaging_templates (id, name, channel, context, body, description) VALUES (?, ?, ?, ?, ?, ?)`,
		"tpl-ar-ptp", "PTP follow-up", "whatsapp", "arrears", "Your promised payment of {{amount}} is due on {{dueDate}}.", "Promise to pay reminder")

	mustExec(session, `INSERT INTO attachments_attachments (id, context, entity_id, file_name, mime_type, size_bytes, url, captured_by, captured_at, sync, note) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"att-1", "arrears_pop", "ar-1002", "pop.jpg", "image/jpeg", int64(182311), "https://object-store.local/att-1", "u-2", tomorrow, "synced", "Imported demo attachment")

	return nil
}

func mustExec(session *gocql.Session, query string, values ...any) {
	if err := session.Query(query, values...).Exec(); err != nil {
		log.Fatalf("seed query failed: %v query=%s", err, query)
	}
}

func strp(value string) *string { return &value }
