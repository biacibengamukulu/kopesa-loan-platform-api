package app

import (
	"fmt"

	auditapp "github.com/biangacila/kopesa-loan-platform-api/internal/audit/application"
	auditcql "github.com/biangacila/kopesa-loan-platform-api/internal/audit/infrastructure/cassandra"
	auditctrl "github.com/biangacila/kopesa-loan-platform-api/internal/audit/interfaces/controllers"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	arrearsapp "github.com/biangacila/kopesa-loan-platform-api/internal/arrears/application"
	arrearscql "github.com/biangacila/kopesa-loan-platform-api/internal/arrears/infrastructure/cassandra"
	arrearsctrl "github.com/biangacila/kopesa-loan-platform-api/internal/arrears/interfaces/controllers"
	attachmentsapp "github.com/biangacila/kopesa-loan-platform-api/internal/attachments/application"
	attachmentscql "github.com/biangacila/kopesa-loan-platform-api/internal/attachments/infrastructure/cassandra"
	attachmentsprovider "github.com/biangacila/kopesa-loan-platform-api/internal/attachments/infrastructure/provider"
	attachmentsctrl "github.com/biangacila/kopesa-loan-platform-api/internal/attachments/interfaces/controllers"
	campaignapp "github.com/biangacila/kopesa-loan-platform-api/internal/campaign/application"
	campaigncql "github.com/biangacila/kopesa-loan-platform-api/internal/campaign/infrastructure/cassandra"
	campaignctrl "github.com/biangacila/kopesa-loan-platform-api/internal/campaign/interfaces/controllers"
	customerapp "github.com/biangacila/kopesa-loan-platform-api/internal/customer/application"
	customercql "github.com/biangacila/kopesa-loan-platform-api/internal/customer/infrastructure/cassandra"
	customerctrl "github.com/biangacila/kopesa-loan-platform-api/internal/customer/interfaces/controllers"
	loanapp "github.com/biangacila/kopesa-loan-platform-api/internal/loan/application"
	loancql "github.com/biangacila/kopesa-loan-platform-api/internal/loan/infrastructure/cassandra"
	loanctrl "github.com/biangacila/kopesa-loan-platform-api/internal/loan/interfaces/controllers"
	messagingapp "github.com/biangacila/kopesa-loan-platform-api/internal/messaging/application"
	messagingcql "github.com/biangacila/kopesa-loan-platform-api/internal/messaging/infrastructure/cassandra"
	messagingprovider "github.com/biangacila/kopesa-loan-platform-api/internal/messaging/infrastructure/provider"
	messagingctrl "github.com/biangacila/kopesa-loan-platform-api/internal/messaging/interfaces/controllers"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/auth"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/cassandra"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/config"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/httpx"
	"github.com/biangacila/kopesa-loan-platform-api/internal/platform/kafka"
	reportingapp "github.com/biangacila/kopesa-loan-platform-api/internal/reporting/application"
	reportingcql "github.com/biangacila/kopesa-loan-platform-api/internal/reporting/infrastructure/cassandra"
	reportingctrl "github.com/biangacila/kopesa-loan-platform-api/internal/reporting/interfaces/controllers"
)

const apiBasePath = "/backend-kopesa/api/v1"

func Build(cfg config.Config) (*fiber.App, func(), error) {
	session, err := cassandra.NewSession(cfg, true)
	if err != nil {
		return nil, nil, fmt.Errorf("connect cassandra: %w", err)
	}

	api := fiber.New()
	api.Use(recover.New())
	api.Use(httpx.RequestIDMiddleware())
	api.Use(httpx.ErrorMiddleware())

	authManager := auth.NewManager(cfg.JWTSecret)

	registerInfraRoutes(api)
	buildRoutes(api, authManager, session, cfg)

	return api, func() { session.Close() }, nil
}

func buildRoutes(api *fiber.App, authManager *auth.Manager, session *gocql.Session, cfg config.Config) {
	_ = kafka.NewNoopPublisher()
	auditSvc := auditapp.NewService(auditcql.NewRepository(session))
	customerSvc := customerapp.NewService(
		customercql.NewUserRepository(session),
		customercql.NewRoleRepository(session),
		customercql.NewBranchRepository(session),
		customercql.NewAreaRepository(session),
		authManager,
	)
	loanSvc := loanapp.NewService(loancql.NewRepository(session))
	arrearsSvc := arrearsapp.NewService(arrearscql.NewRepository(session))
	campaignSvc := campaignapp.NewService(campaigncql.NewRepository(session))
	messagingSvc := messagingapp.NewService(messagingcql.NewRepository(session), messagingprovider.NewGateway(cfg))
	attachmentsSvc := attachmentsapp.NewService(attachmentscql.NewRepository(session), attachmentsprovider.NewDropbox(cfg))
	reportingSvc := reportingapp.NewService(reportingcql.NewRepository(session))

	v1 := api.Group(apiBasePath)
	customerctrl.RegisterRoutes(v1, customerSvc, authManager)

	protected := v1.Group("", auth.Middleware(authManager, false))
	protected.Use(httpx.AuditMiddleware(auditSvc))
	loanctrl.RegisterRoutes(protected, loanSvc)
	arrearsctrl.RegisterRoutes(protected, arrearsSvc)
	campaignctrl.RegisterRoutes(protected, campaignSvc)
	messagingctrl.RegisterRoutes(protected, messagingSvc)
	attachmentsctrl.RegisterRoutes(protected, attachmentsSvc)
	reportingctrl.RegisterRoutes(protected, reportingSvc)
	auditctrl.RegisterRoutes(protected, auditSvc)
}

func registerInfraRoutes(api *fiber.App) {
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	api.Get("/backend-kopesa/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	api.Get(apiBasePath+"/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	api.Get("/backend-kopesa/api/openapi.yaml", func(c *fiber.Ctx) error {
		return c.SendFile("/app/api/openapi.yaml")
	})
}
