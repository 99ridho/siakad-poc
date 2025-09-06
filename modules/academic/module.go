package academic

import (
	"siakad-poc/common"
	"siakad-poc/constants"
	"siakad-poc/db/repositories"
	"siakad-poc/middlewares"
	"siakad-poc/modules"
	"siakad-poc/modules/academic/handlers"
	"siakad-poc/modules/academic/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AcademicModule struct {
	academicRepository      repositories.AcademicRepository
	courseOfferingUseCase   *usecases.CourseOfferingUseCase
	courseEnrollmentUseCase *usecases.CourseEnrollmentUseCase
	courseOfferingHandler   *handlers.CourseOfferingHandler
	courseEnrollmentHandler *handlers.CourseEnrollmentHandler
}

// Compile time interface conformance check
var _ modules.RoutableModule = (*AcademicModule)(nil)

func NewModule(pool *pgxpool.Pool) *AcademicModule {
	txExecutor := common.NewPgxTransactionExecutor(pool)
	academicRepository := repositories.NewDefaultAcademicRepository(pool)

	courseOfferingUseCase := usecases.NewCourseOfferingUseCase(academicRepository)
	courseEnrollmentUseCase := usecases.NewCourseEnrollmentUseCase(academicRepository, txExecutor)

	courseOfferingHandler := handlers.NewCourseOfferingHandler(courseOfferingUseCase)
	courseEnrollmentHandler := handlers.NewEnrollmentHandler(courseEnrollmentUseCase)

	return &AcademicModule{
		academicRepository:      academicRepository,
		courseOfferingUseCase:   courseOfferingUseCase,
		courseEnrollmentUseCase: courseEnrollmentUseCase,
		courseOfferingHandler:   courseOfferingHandler,
		courseEnrollmentHandler: courseEnrollmentHandler,
	}
}

func (m *AcademicModule) SetupRoutes(fiberApp *fiber.App, prefix string) {
	academicGroup := fiberApp.Group(prefix)
	academicGroup.Use(middlewares.JWT())
	academicGroup.Post(
		"/course-offering/:id/enroll",
		middlewares.ShouldBeAccessedByRoles([]constants.RoleType{constants.RoleStudent}),
		m.courseEnrollmentHandler.HandleCourseEnrollment,
	)

	// Course offering CRUD routes (Admin and Koorprodi only)
	academicGroup.Get(
		"/course-offerings",
		middlewares.ShouldBeAccessedByRoles([]constants.RoleType{constants.RoleAdmin, constants.RoleKoorprodi}),
		m.courseOfferingHandler.HandleListCourseOfferings,
	)
	academicGroup.Post(
		"/course-offering",
		middlewares.ShouldBeAccessedByRoles([]constants.RoleType{constants.RoleAdmin, constants.RoleKoorprodi}),
		m.courseOfferingHandler.HandleCreateCourseOffering,
	)
	academicGroup.Put(
		"/course-offering/:id",
		middlewares.ShouldBeAccessedByRoles([]constants.RoleType{constants.RoleAdmin, constants.RoleKoorprodi}),
		m.courseOfferingHandler.HandleUpdateCourseOffering,
	)
	academicGroup.Delete(
		"/course-offering/:id",
		middlewares.ShouldBeAccessedByRoles([]constants.RoleType{constants.RoleAdmin, constants.RoleKoorprodi}),
		m.courseOfferingHandler.HandleDeleteCourseOffering,
	)
}
