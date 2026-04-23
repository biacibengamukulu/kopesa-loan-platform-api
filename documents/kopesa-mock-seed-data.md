# Kopesa — Mock Seed Data for Backend Implementation

> Purpose: a single, copy-pasteable seed dataset that mirrors the React app's mock state.
> Use this to (a) seed Cassandra for dev/QA, (b) drive Postman/Newman tests against the Go/Fiber API,
> and (c) validate event payloads on Kafka topics. All money is in **ZAR minor units (cents)**.
> All IDs shown as human-friendly strings — replace with TimeUUID v1 on insert and keep the mapping table at the end.

---

## 0. Conventions

- **Currency**: ZAR. `principal_minor = rand * 100`. Example: R 25 000,00 → `2500000`.
- **Phones**: E.164, `+27...`.
- **Timestamps**: ISO-8601 UTC. Where you see `T-1d`, replace with `now - 1 day` at seed time.
- **Branch code**: 5-char external reference (e.g. `JHB01`) — must be unique, used by external core-banking system.
- **Tenant**: single tenant `kopesa-za` for all rows.

---

## 1. Areas

| id        | name           |
|-----------|----------------|
| area-gp   | Gauteng        |
| area-wc   | Western Cape   |
| area-kzn  | KwaZulu-Natal  |

```json
[
  {"id":"area-gp","name":"Gauteng"},
  {"id":"area-wc","name":"Western Cape"},
  {"id":"area-kzn","name":"KwaZulu-Natal"}
]
```

---

## 2. Branches

`branch_code` is the **external reference** used by the legacy core-banking / GL system.

| id           | name              | branch_code | area_id   |
|--------------|-------------------|-------------|-----------|
| br-jhb-cbd   | Johannesburg CBD  | JHB01       | area-gp   |
| br-soweto    | Soweto            | JHB02       | area-gp   |
| br-cpt-cbd   | Cape Town CBD     | CPT01       | area-wc   |
| br-dbn-cbd   | Durban CBD        | DBN01       | area-kzn  |

```json
[
  {"id":"br-jhb-cbd","name":"Johannesburg CBD","branch_code":"JHB01","area_id":"area-gp"},
  {"id":"br-soweto","name":"Soweto","branch_code":"JHB02","area_id":"area-gp"},
  {"id":"br-cpt-cbd","name":"Cape Town CBD","branch_code":"CPT01","area_id":"area-wc"},
  {"id":"br-dbn-cbd","name":"Durban CBD","branch_code":"DBN01","area_id":"area-kzn"}
]
```

---

## 3. Roles & Permissions

10 roles. Permissions match `src/data/roles.ts` exactly.

```json
[
  {"id":"field_officer","label":"Field Officer","scope":"field",
   "permissions":["campaign.view","campaign.plan","campaign.lead.capture","campaign.lead.qualify","loan.view","messaging.send","messaging.view","attachments.upload","attachments.view"]},
  {"id":"collector","label":"Collector","scope":"field",
   "permissions":["arrears.view","arrears.action","arrears.ptp.create","arrears.payment.capture","loan.view","messaging.send","messaging.view","attachments.upload","attachments.view"]},
  {"id":"branch_agent","label":"Branch Agent","scope":"branch",
   "permissions":["loan.view","loan.application.create","arrears.view","arrears.action","arrears.ptp.create","campaign.view","messaging.send","messaging.view","attachments.upload","attachments.view"]},
  {"id":"consultant","label":"Loan Consultant","scope":"branch",
   "permissions":["loan.view","loan.application.create","loan.assess","maker","messaging.send","messaging.view","attachments.upload","attachments.view"]},
  {"id":"branch_manager","label":"Branch Manager","scope":"branch",
   "permissions":["arrears.view","arrears.allocate","arrears.escalate","campaign.view","campaign.plan","campaign.route.assign","loan.view","loan.assess","loan.approve","checker","reports.exec","messaging.send","messaging.view","attachments.view"]},
  {"id":"area_manager","label":"Area Manager","scope":"area",
   "permissions":["arrears.view","arrears.allocate","arrears.escalate","arrears.writeoff","campaign.view","campaign.plan","loan.view","loan.approve","checker","reports.exec","messaging.view","attachments.view"]},
  {"id":"finance_officer","label":"Finance Officer","scope":"global",
   "permissions":["arrears.view","arrears.payment.capture","arrears.writeoff","loan.view","loan.disburse","loan.servicing","finance.reconcile","checker","attachments.view","attachments.upload","messaging.view"]},
  {"id":"compliance_officer","label":"Compliance Officer","scope":"global",
   "permissions":["loan.view","loan.assess","arrears.view","reports.exec","checker","attachments.view","attachments.upload","messaging.view"]},
  {"id":"admin","label":"System Admin","scope":"global",
   "permissions":["admin.users","admin.config","reports.exec","arrears.view","campaign.view","loan.view","messaging.view","attachments.view"]},
  {"id":"executive","label":"Executive","scope":"global",
   "permissions":["arrears.view","campaign.view","loan.view","reports.exec","messaging.view","attachments.view"]}
]
```

---

## 4. Users (login credentials for demo)

Default password for **all** seed users: `Kopesa@2026`. Hash with bcrypt cost 12 before insert.

| id   | full_name            | email                       | role                | branch_id   | area_id    |
|------|----------------------|-----------------------------|---------------------|-------------|------------|
| u-1  | Thabo Mokoena        | thabo@kopesa.co.za          | field_officer       | br-soweto   | area-gp    |
| u-2  | Naledi Dlamini       | naledi@kopesa.co.za         | collector           | br-jhb-cbd  | area-gp    |
| u-3  | Sipho Khumalo        | sipho@kopesa.co.za          | branch_agent        | br-jhb-cbd  | area-gp    |
| u-4  | Lerato Nkosi         | lerato@kopesa.co.za         | consultant          | br-cpt-cbd  | area-wc    |
| u-5  | Pieter van Wyk       | pieter@kopesa.co.za         | branch_manager      | br-cpt-cbd  | area-wc    |
| u-6  | Zanele Mthembu       | zanele@kopesa.co.za         | area_manager        | —           | area-kzn   |
| u-7  | Hendrik Botha        | hendrik@kopesa.co.za        | finance_officer     | —           | —          |
| u-8  | Aisha Patel          | aisha@kopesa.co.za          | compliance_officer  | —           | —          |
| u-9  | System Admin         | admin@kopesa.co.za          | admin               | —           | —          |
| u-10 | Mervin Biangacila    | exec@kopesa.co.za           | executive           | —           | —          |

```json
[
  {"id":"u-1","full_name":"Thabo Mokoena","email":"thabo@kopesa.co.za","role":"field_officer","branch_id":"br-soweto","area_id":"area-gp","password":"Kopesa@2026"},
  {"id":"u-2","full_name":"Naledi Dlamini","email":"naledi@kopesa.co.za","role":"collector","branch_id":"br-jhb-cbd","area_id":"area-gp","password":"Kopesa@2026"},
  {"id":"u-3","full_name":"Sipho Khumalo","email":"sipho@kopesa.co.za","role":"branch_agent","branch_id":"br-jhb-cbd","area_id":"area-gp","password":"Kopesa@2026"},
  {"id":"u-4","full_name":"Lerato Nkosi","email":"lerato@kopesa.co.za","role":"consultant","branch_id":"br-cpt-cbd","area_id":"area-wc","password":"Kopesa@2026"},
  {"id":"u-5","full_name":"Pieter van Wyk","email":"pieter@kopesa.co.za","role":"branch_manager","branch_id":"br-cpt-cbd","area_id":"area-wc","password":"Kopesa@2026"},
  {"id":"u-6","full_name":"Zanele Mthembu","email":"zanele@kopesa.co.za","role":"area_manager","area_id":"area-kzn","password":"Kopesa@2026"},
  {"id":"u-7","full_name":"Hendrik Botha","email":"hendrik@kopesa.co.za","role":"finance_officer","password":"Kopesa@2026"},
  {"id":"u-8","full_name":"Aisha Patel","email":"aisha@kopesa.co.za","role":"compliance_officer","password":"Kopesa@2026"},
  {"id":"u-9","full_name":"System Admin","email":"admin@kopesa.co.za","role":"admin","password":"Kopesa@2026"},
  {"id":"u-10","full_name":"Mervin Biangacila","email":"exec@kopesa.co.za","role":"executive","password":"Kopesa@2026"}
]
```

---

## 5. Loans

`rate` = annual interest %. Money in cents.

```json
[
  {"id":"LN-44021","client_name":"Bongani Sithole","client_phone":"+27821234567","branch_id":"br-jhb-cbd","branch_code":"JHB01",
   "principal_minor":2500000,"term_months":12,"rate":28.0,"status":"active","created_by":"u-3",
   "outstanding_balance_minor":1890000,"next_due_date":"2025-11-30"},
  {"id":"LN-44078","client_name":"Nomvula Khoza","client_phone":"+27735550102","branch_id":"br-jhb-cbd","branch_code":"JHB01",
   "principal_minor":5000000,"term_months":18,"rate":26.0,"status":"active","created_by":"u-3",
   "outstanding_balance_minor":4120000,"next_due_date":"2025-11-15"},
  {"id":"LN-44102","client_name":"Tshepo Mabaso","client_phone":"+27842220190","branch_id":"br-soweto","branch_code":"JHB02",
   "principal_minor":3500000,"term_months":12,"rate":28.0,"status":"active","created_by":"u-3",
   "outstanding_balance_minor":2870000,"next_due_date":"2025-11-05"},
  {"id":"LN-44900","client_name":"Mpho Radebe","client_phone":"+27821001100","branch_id":"br-soweto","branch_code":"JHB02",
   "principal_minor":800000,"term_months":6,"rate":30.0,"status":"approved","created_by":"u-4","assessed_by":"u-4","approved_by":"u-5",
   "outstanding_balance_minor":800000},
  {"id":"LN-44901","client_name":"Lerato Maseko","client_phone":"+27732209981","branch_id":"br-soweto","branch_code":"JHB02",
   "principal_minor":500000,"term_months":6,"rate":30.0,"status":"assessment","created_by":"u-4",
   "outstanding_balance_minor":500000}
]
```

## 6. Loan Applications

```json
[
  {"id":"app-1","loan_id":"LN-44900","client_name":"Mpho Radebe","client_id_no":"8806125555084",
   "monthly_income_minor":1200000,"monthly_expenses_minor":750000,"requested_amount_minor":800000,"term_months":6,
   "branch_id":"br-soweto","branch_code":"JHB02","created_by":"u-4","status":"approved","affordability_minor":450000,
   "assessment_note":"Stable employment 3yrs","approval_note":"Approved within branch limit"},
  {"id":"app-2","loan_id":"LN-44901","client_name":"Lerato Maseko","client_id_no":"9203054444081",
   "monthly_income_minor":850000,"monthly_expenses_minor":520000,"requested_amount_minor":500000,"term_months":6,
   "branch_id":"br-soweto","branch_code":"JHB02","created_by":"u-4","status":"submitted"}
]
```

---

## 7. Arrears Cases

```json
[
  {"id":"ar-1001","loan_id":"LN-44021","client_name":"Bongani Sithole","client_phone":"+27821234567",
   "branch_id":"br-jhb-cbd","branch_code":"JHB01","days_past_due":12,
   "arrears_amount_minor":245000,"outstanding_balance_minor":1890000,
   "status":"allocated","assigned_to":"u-2","last_action_at":"T-1d","ptps":[]},
  {"id":"ar-1002","loan_id":"LN-44078","client_name":"Nomvula Khoza","client_phone":"+27735550102",
   "branch_id":"br-jhb-cbd","branch_code":"JHB01","days_past_due":34,
   "arrears_amount_minor":580000,"outstanding_balance_minor":4120000,
   "status":"ptp","assigned_to":"u-2","last_action_at":"T-2d",
   "ptps":[{"id":"ptp-1","amount_minor":250000,"promised_at":"T+3d","captured_by":"u-2","status":"pending","note":"Will pay after payday"}]},
  {"id":"ar-1003","loan_id":"LN-44102","client_name":"Tshepo Mabaso","client_phone":"+27842220190",
   "branch_id":"br-soweto","branch_code":"JHB02","days_past_due":67,
   "arrears_amount_minor":990000,"outstanding_balance_minor":2870000,
   "status":"broken_ptp","assigned_to":"u-2","ptps":[]},
  {"id":"ar-1004","loan_id":"LN-44211","client_name":"Refilwe Modise","client_phone":"+27719041188",
   "branch_id":"br-cpt-cbd","branch_code":"CPT01","days_past_due":5,
   "arrears_amount_minor":120000,"outstanding_balance_minor":980000,
   "status":"new","ptps":[]},
  {"id":"ar-1005","loan_id":"LN-44318","client_name":"Sandile Ngubane","client_phone":"+27827002244",
   "branch_id":"br-dbn-cbd","branch_code":"DBN01","days_past_due":92,
   "arrears_amount_minor":1420000,"outstanding_balance_minor":3650000,
   "status":"escalated","assigned_to":"u-3","ptps":[]},
  {"id":"ar-1006","loan_id":"LN-44402","client_name":"Karabo Mahlangu","client_phone":"+27796128801",
   "branch_id":"br-jhb-cbd","branch_code":"JHB01","days_past_due":18,
   "arrears_amount_minor":340000,"outstanding_balance_minor":2210000,
   "status":"in_progress","assigned_to":"u-2","ptps":[]}
]
```

### Arrears CSV import template (header row)

```
loan_id,client_name,client_phone,branch_code,days_past_due,arrears_amount,outstanding_balance
LN-99001,Example Client,+27820000000,JHB01,7,1500.00,12000.00
```

> Backend validation: `branch_code` must resolve to an existing branch; reject row otherwise with code `BRANCH_CODE_NOT_FOUND`.

---

## 8. Campaigns

```json
[
  {"id":"cmp-1","name":"Soweto Spring Drive","branch_id":"br-soweto","branch_code":"JHB02","status":"active",
   "start_date":"2025-09-01","end_date":"2025-12-15","target_leads":200,"captured_leads":84},
  {"id":"cmp-2","name":"JHB CBD Walk-in Push","branch_id":"br-jhb-cbd","branch_code":"JHB01","status":"active",
   "start_date":"2025-10-01","end_date":"2025-11-30","target_leads":150,"captured_leads":41},
  {"id":"cmp-3","name":"Cape Town Suburbs Q4","branch_id":"br-cpt-cbd","branch_code":"CPT01","status":"planning",
   "start_date":"2025-11-15","end_date":"2026-02-28","target_leads":300,"captured_leads":0}
]
```

## 9. Campaign Routes

```json
[
  {"id":"rt-1","campaign_id":"cmp-1","date":"TODAY","assigned_to":"u-1","status":"in_progress",
   "stops":[
     {"id":"s-1","address":"12 Vilakazi St","suburb":"Orlando West","status":"visited","visited_at":"NOW"},
     {"id":"s-2","address":"88 Klipspruit Valley Rd","suburb":"Klipspruit","status":"lead","visited_at":"NOW","lead_id":"ld-1"},
     {"id":"s-3","address":"4 Mofolo North Ext","suburb":"Mofolo","status":"no_contact"},
     {"id":"s-4","address":"203 Chris Hani Rd","suburb":"Diepkloof","status":"pending"},
     {"id":"s-5","address":"17 Pimville Zone 4","suburb":"Pimville","status":"pending"}
   ]},
  {"id":"rt-2","campaign_id":"cmp-2","date":"TODAY","assigned_to":"u-3","status":"planned",
   "stops":[
     {"id":"s-6","address":"Carlton Centre, 150 Commissioner","suburb":"Marshalltown","status":"pending"},
     {"id":"s-7","address":"90 Pritchard St","suburb":"JHB CBD","status":"pending"}
   ]}
]
```

## 10. Leads

```json
[
  {"id":"ld-1","campaign_id":"cmp-1","full_name":"Mpho Radebe","phone":"+27821001100",
   "suburb":"Orlando West","captured_by":"u-1","captured_at":"T-1h","qualified":"qualified","estimated_amount_minor":800000},
  {"id":"ld-2","campaign_id":"cmp-1","full_name":"Lerato Maseko","phone":"+27732209981",
   "suburb":"Klipspruit","captured_by":"u-1","captured_at":"T-2h","qualified":"pending","estimated_amount_minor":500000},
  {"id":"ld-3","campaign_id":"cmp-2","full_name":"Sibusiso Dube","phone":"+27845557711",
   "suburb":"Marshalltown","captured_by":"u-3","captured_at":"T-1d","qualified":"rejected"}
]
```

---

## 11. Message Templates

```json
[
  {"id":"tpl-ar-reminder","name":"Friendly reminder","channel":"both","context":"arrears",
   "body":"Hi {{client}}, your Kopesa account is {{daysPastDue}} days past due. Please pay {{amount}} to settle. Reply STOP to opt out."},
  {"id":"tpl-ar-ptp","name":"PTP confirmation","channel":"whatsapp","context":"arrears",
   "body":"Thanks {{client}}. We have recorded your promise to pay {{amount}} by {{dueDate}}. — Kopesa {{branch}}"},
  {"id":"tpl-ar-broken","name":"Broken PTP follow-up","channel":"sms","context":"arrears",
   "body":"Hi {{client}}, we did not receive the {{amount}} promised on {{dueDate}}. Please contact your branch today."},
  {"id":"tpl-ar-final","name":"Final demand","channel":"sms","context":"arrears",
   "body":"FINAL NOTICE: {{client}}, account {{loanId}} is {{daysPastDue}} days in arrears. Pay {{amount}} now to avoid escalation."},
  {"id":"tpl-ld-welcome","name":"Lead welcome","channel":"whatsapp","context":"lead",
   "body":"Hi {{client}}, thanks for chatting with our team in {{suburb}}. A consultant will call within 24h about your loan."},
  {"id":"tpl-rt-visit","name":"Visit confirmation","channel":"sms","context":"campaign_route",
   "body":"Hi {{client}}, our agent will visit you today at {{address}}. Have your ID and 3 payslips ready."},
  {"id":"tpl-app-received","name":"Application received","channel":"whatsapp","context":"loan_application",
   "body":"Hi {{client}}, we received your application for {{amount}}. We will respond within 1 business day."}
]
```

## 12. Message Log

```json
[
  {"id":"msg-1","context":"arrears","entity_id":"ar-1002","channel":"whatsapp",
   "template_id":"tpl-ar-ptp","to":"+27735550102",
   "body":"Thanks Nomvula, we have recorded your promise to pay R 2,500 by 25 Apr.",
   "status":"delivered","sent_at":"T-1h","delivered_at":"T-59m","sent_by":"u-2","next_touch_at":"T+3d"},
  {"id":"msg-2","context":"arrears","entity_id":"ar-1003","channel":"sms",
   "template_id":"tpl-ar-broken","to":"+27842220190",
   "body":"Hi Tshepo, we did not receive the R 4,000 promised on 15 Apr.",
   "status":"failed","sent_at":"T-2h","failed_reason":"Unreachable handset","sent_by":"u-2"}
]
```

---

## 13. Subscription Plans (role gating, scaffold only)

```json
[
  {"id":"plan-starter","name":"Starter","monthly_zar_minor":99900,
   "allowed_roles":["branch_agent","collector","field_officer","branch_manager","admin"]},
  {"id":"plan-growth","name":"Growth","monthly_zar_minor":249900,
   "allowed_roles":["branch_agent","collector","field_officer","consultant","branch_manager","area_manager","finance_officer","admin"]},
  {"id":"plan-enterprise","name":"Enterprise","monthly_zar_minor":599900,
   "allowed_roles":["field_officer","collector","branch_agent","consultant","branch_manager","area_manager","finance_officer","compliance_officer","admin","executive"]}
]
```

Demo tenant `kopesa-za` is on `plan-enterprise` so all roles are unlocked.

---

## 14. Sample Kafka events (post-seed)

After loading the data above, the backend should be able to **replay** these events without errors:

```json
{"specversion":"1.0","type":"za.kopesa.loans.loan.disbursed.v1","source":"loans-service",
 "id":"evt-001","time":"2025-04-22T10:15:00Z","subject":"loan/LN-44021",
 "data":{"loan_id":"LN-44021","branch_code":"JHB01","disbursed_by":"u-7","amount_minor":2500000}}

{"specversion":"1.0","type":"za.kopesa.arrears.ptp.created.v1","source":"arrears-service",
 "id":"evt-002","time":"2025-04-22T11:00:00Z","subject":"arrears/ar-1002",
 "data":{"case_id":"ar-1002","ptp_id":"ptp-1","amount_minor":250000,"promised_at":"2025-04-25","captured_by":"u-2"}}

{"specversion":"1.0","type":"za.kopesa.campaigns.lead.captured.v1","source":"campaigns-service",
 "id":"evt-003","time":"2025-04-22T09:30:00Z","subject":"lead/ld-1",
 "data":{"lead_id":"ld-1","campaign_id":"cmp-1","branch_code":"JHB02","captured_by":"u-1"}}
```

---

## 15. ID mapping table (string → TimeUUID)

When seeding production-like environments, generate a TimeUUID v1 for each string ID and store the mapping in
`seed_id_map (legacy_id text PRIMARY KEY, uuid timeuuid)` so foreign-key references in the JSON above can be
rewritten in a single pass before insert.

---

**End of seed file. Total volume: 4 areas, 4 branches, 10 roles, 10 users, 5 loans, 2 applications, 6 arrears cases, 3 campaigns, 2 routes (7 stops), 3 leads, 7 templates, 2 messages, 3 subscription plans.**
