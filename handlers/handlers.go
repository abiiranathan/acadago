package handlers

import (
	"github.com/abiiranathan/acada/middleware"
	"github.com/abiiranathan/acada/services"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(svcs *services.Services) *gin.Engine {
	r := gin.Default()
	r.RedirectTrailingSlash = false
	r.SetTrustedProxies(nil)
	gin.SetMode(gin.ReleaseMode)

	api := r.Group("/api")

	LoginRequired := middleware.LoginRequired(svcs.UserService)
	AdminRequired := middleware.AdminRequired()

	user_handle := NewUserHandler(svcs.UserService)
	api.POST("/auth/login", user_handle.Login())
	api.GET("/auth/verify", LoginRequired, user_handle.VerifyToken())

	user_api := api.Group("users", LoginRequired, AdminRequired)
	{
		user_api.GET("", user_handle.GetAllUsers())
		user_api.POST("", user_handle.CreateUser())
		user_api.GET("/:id", user_handle.GetUser())
		user_api.PUT("/:id", user_handle.UpdateUser())
		user_api.DELETE("/:id", user_handle.DeleteUser())
	}

	// Programs
	program_handle := NewProgramHandler(svcs.ProgramService)
	program_api := api.Group("programs", LoginRequired)
	{
		program_api.GET("", program_handle.GetAllPrograms())
		program_api.POST("", program_handle.CreateProgram())
		program_api.GET("/:id", program_handle.GetProgram())
		program_api.PUT("/:id", program_handle.UpdateProgram())
		program_api.DELETE("/:id", program_handle.DeleteProgram())
	}

	// Course Units
	course_unit_handle := NewCourseUnitHandler(svcs.CourseUnitService)
	course_unit_api := api.Group("courseunits", LoginRequired)
	{
		course_unit_api.GET("", course_unit_handle.GetAllCourseUnits())
		course_unit_api.POST("", course_unit_handle.CreateCourseUnit())
		course_unit_api.GET("/:id", course_unit_handle.GetCourseUnit())
		course_unit_api.PUT("/:id", course_unit_handle.UpdateCourseUnit())
		course_unit_api.DELETE("/:id", course_unit_handle.DeleteCourseUnit())
	}

	// Resources
	resource_handle := NewResourceHandler(svcs.ResourceService)
	resource_api := api.Group("resources", LoginRequired)
	{
		resource_api.GET("", resource_handle.GetAllResources())
		resource_api.POST("", resource_handle.CreateResource())
		resource_api.GET("/:id", resource_handle.GetResource())
		resource_api.DELETE("/:id", resource_handle.DeleteResource())
	}

	return r
}
