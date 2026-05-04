package routes

import (
	"schoolbackend/controllers"
	"schoolbackend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/login", controllers.Login)
	r.POST("/OTP", controllers.SendOTP)
	r.POST("/register", controllers.Register)
	r.GET("/getrole", controllers.GetRole)
	r.Static("/images", "./public/images")
	r.GET("/api/studentenrollreport", controllers.GetStudentEnrollmentReport)
	r.GET("/api/studentgenderreport", controllers.GetGenderStatsReport)
	r.GET("/api/studentbyaddress", controllers.GetStudentByAddress)
	r.GET("/api/studentpoor", controllers.Detailedlistofpoorstudentsbyclassandacademicyear)
	r.GET("/api/Numberofpoorstudentsbyclassandgender", controllers.Numberofpoorstudentsbyclassandgender)
	r.GET("/api/NumberOfPoorStudentsByProvinceDistrict", controllers.NumberOfPoorStudentsByProvinceDistrict)
	r.GET("/api/PoorVsNonPoorStudentsComparisonByYear", controllers.PoorVsNonPoorStudentsComparisonByYear)
	r.GET("/api/TotalSummaryStudentsTeachersClassesAcademicYears", controllers.TotalSummaryStudentsTeachersClassesAcademicYears)
	r.GET("/api/StudentSubjectExamScoresByAcademicYear", controllers.StudentSubjectExamScoresByAcademicYear)
	r.GET("/api/DetailedDisabledStudentsByDisabilityType", controllers.DetailedDisabledStudentsByDisabilityType)
	// Protected routes
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		// Role
		auth.POST("/role", middleware.PermissionMiddleware("add-role"), controllers.CreateRole)
		auth.GET("/role", middleware.PermissionMiddleware("view-role"), controllers.GetRole)
		auth.PUT("/role/:id", middleware.PermissionMiddleware("edit-role"), controllers.UpdateRole)
		auth.PUT("/changestatusrole/:id", middleware.PermissionMiddleware("changestatusrole"), controllers.ChangeStatusRole)

		// Permission
		auth.POST("/rolepermission", middleware.PermissionMiddleware("assign-permission"), controllers.CreateRolePermissions)
		auth.DELETE("/rolepermission", middleware.PermissionMiddleware("remove-permission"), controllers.DeleteRolePermissions)
		auth.GET("/roleassignpermission/:id", middleware.PermissionMiddleware("view-roleassignpermission"), controllers.GetRolePermissions)

		// Location
		auth.GET("/province", middleware.PermissionMiddleware("view-province"), controllers.GetProvince)
		auth.GET("/district/:id", middleware.PermissionMiddleware("view-district"), controllers.GetDistrict)
		auth.GET("/district-by-id/:id", middleware.PermissionMiddleware("view-district"), controllers.GetDistrictByID)
		auth.GET("/communce/:id", middleware.PermissionMiddleware("view-communce"), controllers.GetCommune)
		auth.GET("/commune-by-id/:id", middleware.PermissionMiddleware("view-communce"), controllers.GetCommuneByID)
		auth.GET("/village/:id", middleware.PermissionMiddleware("view-village"), controllers.GetVillage)
		auth.GET("/village-by-id/:id", middleware.PermissionMiddleware("view-village"), controllers.GetVillageByID)

		// User
		auth.GET("/viewuser", middleware.PermissionMiddleware("view-user"), controllers.GetUsers)
		auth.GET("/user/:id", middleware.PermissionMiddleware("view-user"), controllers.GetUser)
		auth.PUT("/changestatususer/:id", middleware.PermissionMiddleware("change-status-user"), controllers.ChangeStatusUser)
		auth.PUT("/user/:id", middleware.PermissionMiddleware("edit-user"), controllers.UpdateUser)
		auth.PUT("/changepassword/:id", middleware.PermissionMiddleware("change-password"), controllers.ChangePassword)

		// Class
		auth.POST("/class", middleware.PermissionMiddleware("add-class"), controllers.SaveClass)
		auth.PUT("/class", middleware.PermissionMiddleware("edit-class"), controllers.SaveClass)
		auth.GET("/class", middleware.PermissionMiddleware("view-class"), controllers.Handleclass)
		auth.PUT("/changestatusclass/:id", middleware.PermissionMiddleware("change-status-class"), controllers.Handleclass)

		// Subject
		auth.POST("/subject", middleware.PermissionMiddleware("add-subject"), controllers.SaveSubject)
		auth.GET("/subject", middleware.PermissionMiddleware("view-subject"), controllers.Getsubject)
		auth.PUT("/subject", middleware.PermissionMiddleware("edit-subject"), controllers.SaveSubject)
		auth.PUT("/changestatussubject/:id", middleware.PermissionMiddleware("change-status-subject"), controllers.ChangestatusSubject)
		auth.POST("/aissignsubjecttoclass", middleware.PermissionMiddleware("assign-subject-to-class"), controllers.CreateClassSubject)
		auth.GET("/viewclassasignsubject", middleware.PermissionMiddleware("view-class-assign-subject"), controllers.GetClassSubjects)
		auth.GET("/viewclassasignsubjectnotinexamcomponent", middleware.PermissionMiddleware("view-class-assign-subject"), controllers.GetClassSubjects)
		auth.GET("/viewclassnotassignsubject", middleware.PermissionMiddleware("view-class-assign-subject"), controllers.GetSubjectNotAssignedToClass)

		// Acedemicyear
		auth.POST("/academicyear", middleware.PermissionMiddleware("add-academicyear"), controllers.SaveAcademicyear)
		auth.PUT("/academicyear", middleware.PermissionMiddleware("edit-academicyear"), controllers.SaveAcademicyear)
		auth.GET("viewacademicyear", middleware.PermissionMiddleware("view-academicyear"), controllers.Handleacademicyear)
		auth.PUT("/changestatusacademic/:id", middleware.PermissionMiddleware("change-status-academicyear"), controllers.Handleacademicyear)

		// Student
		auth.GET("/getstudent", middleware.PermissionMiddleware("view-student"), controllers.Getstudent)
		auth.POST("/student", middleware.PermissionMiddleware("add-student"), controllers.SaveStudent)
		auth.PUT("/student/:id", middleware.PermissionMiddleware("edit-student"), controllers.SaveStudent)
		auth.GET("/viewstudent", middleware.PermissionMiddleware("view-student"), controllers.HandlStudent)
		auth.PUT("/changestatusstudent/:id", middleware.PermissionMiddleware("change-status-student"), controllers.HandlStudent)
		auth.PUT("/Suspendstudies/:id", middleware.PermissionMiddleware("change-status-student"), controllers.SuspendStudies)
		auth.PUT("/changeschool/:id", middleware.PermissionMiddleware("change-status-student"), controllers.ChangeSchool)
		auth.PUT("/comeback/:id", middleware.PermissionMiddleware("change-status-student"), controllers.Comeback)

		// StudentClass
		auth.POST("/assignstudenttoclass", middleware.PermissionMiddleware("assign-student-to-class"), controllers.CreateStudentClass)
		auth.PUT("/updatestudentclass/:id", middleware.PermissionMiddleware("edit-student-to-class"), controllers.UpdateStudentClass)
		auth.GET("/viewstudentclass/:id", middleware.PermissionMiddleware("view-student-class"), controllers.GetStudentClassByStudentID)
		auth.GET("/viewstudentclasstoaddscore", middleware.PermissionMiddleware("view-student-class"), controllers.GetStudentClassbyClassIDandAcademicyearID)

		// SubjectClass
		auth.GET("/getsubjectclassnotassigntoteacher", middleware.PermissionMiddleware("view-subject"), controllers.GetClassSubjectsNotAssigntoTeacher)
		auth.DELETE("/deleteclasssubject/:id", middleware.PermissionMiddleware("delete-class-subject"), controllers.DeleteClassSubjectByID)
		// TeacherSubject
		auth.POST("/teachersubject", middleware.PermissionMiddleware("assign-teacher-to-class"), controllers.CreateTeacherSubject)
		auth.GET("/teachersubject", middleware.PermissionMiddleware("view-teacher-subject"), controllers.GetClassandSubjectteachbyteacher)
		auth.PUT("/teachersubject/:id", middleware.PermissionMiddleware("change-status-teacher-subject"), controllers.ChangestatusTeachersubject)
		auth.POST("/addteachertoclass", middleware.PermissionMiddleware("assign-teacher-to-class"), controllers.CreateClassTeacher)
		auth.DELETE("/deleteteachersubject/:id", middleware.PermissionMiddleware("delete-teacher-subject"), controllers.DeleteTeachersubject)
		auth.GET("/getteachersubjectbyteacherid", middleware.PermissionMiddleware("view-teacher-subject"), controllers.GetTeachersubjectBYTeacherID)
		auth.PUT("/chnagestatusofclassteacher/:id", middleware.PermissionMiddleware("update-status-class-teacher"), controllers.UpdatestatusClassTeacher)

		// Exam Component
		auth.GET("/examcomponent", middleware.PermissionMiddleware("add-exam-component"), controllers.GetExamComponent)
		auth.POST("/examcomponent", middleware.PermissionMiddleware("add-exam-component"), controllers.CreateExamComponent)
		auth.PUT("/examcomponent/:id", middleware.PermissionMiddleware("edit-exam-component"), controllers.UpdateExamComponent)
		auth.PUT("/changestausexamcomponent/:id", middleware.PermissionMiddleware("change-status-examcomponent"), controllers.ChangeStatusExamcomponent)

		// Score
		auth.POST("/score", middleware.PermissionMiddleware("add-score"), controllers.CreateScore)
		auth.GET("/score", middleware.PermissionMiddleware("view-score"), controllers.GetScore)
		auth.GET("/scoreavg", middleware.PermissionMiddleware("view-score"), controllers.GetAverageScore)
		auth.PUT("/score/:id", middleware.PermissionMiddleware("edit-score"), controllers.UpdateScore)
		auth.DELETE("/score/:id", middleware.PermissionMiddleware("delete-score"), controllers.DeleteScore)
		auth.GET("/scoreavg>5", middleware.PermissionMiddleware("view-score"), controllers.GetAverageScoreBigthan5)
		auth.GET("/scoreAnnual", middleware.PermissionMiddleware("view-score"), controllers.GetAnnualAverageScore)

		// TypeExam
		auth.GET("/typeexam", middleware.PermissionMiddleware("view-type-exam"), controllers.GetTypeExam)

		// Promote
		auth.POST("/promote", middleware.PermissionMiddleware("promote-student"), controllers.PromoteStudent)
		auth.DELETE("/promote/:id", middleware.PermissionMiddleware("delete-promote"), controllers.DeletePromotion)
		auth.GET("/promote", middleware.PermissionMiddleware("view-promote"), controllers.GetPromote)

		// Disabilities
		auth.GET("/disabilities", middleware.PermissionMiddleware("view-disabilities"), controllers.GetDisabilities)
	}
}
