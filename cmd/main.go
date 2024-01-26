package main

import (
	"github.com/Renewdxin/selfMade/internal/adapters/app/auth"
	job2 "github.com/Renewdxin/selfMade/internal/adapters/app/job"
	"github.com/Renewdxin/selfMade/internal/adapters/app/middleware"
	"github.com/Renewdxin/selfMade/internal/adapters/app/user"
	authApp "github.com/Renewdxin/selfMade/internal/adapters/core/auth"
	"github.com/Renewdxin/selfMade/internal/adapters/core/job"
	userApp "github.com/Renewdxin/selfMade/internal/adapters/core/user"
	"github.com/Renewdxin/selfMade/internal/adapters/core/verify"
	"github.com/Renewdxin/selfMade/internal/adapters/framework/database"
	"github.com/Renewdxin/selfMade/internal/adapters/framework/logger"
	"github.com/Renewdxin/selfMade/internal/adapters/framework/mail"
	"github.com/Renewdxin/selfMade/internal/adapters/framework/vaidate"
	"github.com/Renewdxin/selfMade/internal/adapters/framework/web"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	logger.Logger = logger.NewLogger()
	err := godotenv.Load(".env")
	if err != nil {
		logger.Logger.Log(logger.FatalLevel, "无法加载 .env 文件: %v")
	}

	redisClient := database.NewRedisClient()

	mailSender := mail.NewMail()
	verification := verify.NewVerificationCodeService()
	validator := vaidate.NewValidator(redisClient)

	jwtAPI := middleware.NewJWTAdapters()
	jwtHandler := web.NewJWTHandlerAdapter(jwtAPI)

	userCore := userApp.NewUserService()
	userDao, err := database.NewUserDao(os.Getenv("DRIVER_NAME"), os.Getenv("DRIVER_SOURCE_NAME"), userCore)
	if err != nil {
		logger.Logger.Log(logger.FatalLevel, "falied to connect to the user database")
	}

	authCore := authApp.NewAdapter()
	authDao, err := database.NewauthDao(os.Getenv("DRIVER_NAME"), os.Getenv("DRIVER_SOURCE_NAME"), authCore)
	if err != nil {
		logger.Logger.Log(logger.FatalLevel, "falied to connect to the user database")
	}

	userAPI := user.NewUserAPI(userCore, userDao, mailSender, verification, redisClient, validator)
	userHandler := web.NewUserHandler(userAPI)

	authAPI := auth.NewAuthCaseAdapter(authCore, authDao, verification, redisClient, validator, mailSender)
	authHandler := web.NewAuthHandler(authAPI, jwtAPI)

	homeHandler := web.NewHomeHandlerAdapter()

	jobCore := job.NewJobsAdapter()
	jobDao := database.NewJobsDaoAdapter(os.Getenv("DRIVER_NAME"), os.Getenv("DRIVER_SOURCE_NAME"), jobCore)
	jobApp := job2.NewJobCaseAdapter(jobCore, jobDao)
	jobHandler := web.NewJobHandlerAdapter(jobApp)

	adminCore := userApp.NewAdminCoreAdapter()
	_ = database.NewAdminDaoAdapter(os.Getenv("DRIVER_NAME"), os.Getenv("DRIVER_SOURCE_NAME"), adminCore)
	adminApp := user.NewAdminAppAdapter(jobApp, jobDao, userDao)
	adminHandler := web.NewAdminHandlerAdapter(adminApp, jobApp, userAPI)

	r := gin.New()
	//home page
	r.POST("/home", homeHandler.HomePage)

	// auth setting
	apiAccount := r.Group("/auth")
	apiAccount.Use()
	{
		apiAccount.POST("/login", authHandler.Login)
		//apiAccount.POST("/logout")
		apiAccount.POST("/signup", authHandler.Register)
		apiAccount.POST("/password/forget", authHandler.ForgetPassword)
		apiAccount.POST("/password/change", authHandler.ChangePassword)
	}

	// personal info
	apiProfile := r.Group("/profile")
	apiProfile.Use(jwtHandler.JWTHandler())
	{
		// GetUserInfo 通过id得到用户信息，返回项为姓名、性别、出生日期、邮箱、手机号
		apiProfile.GET("/Info/:id", userHandler.GetUserInfo)
		// DeleteUser 用户在删除前需要进行手机验证码验证才能删除
		apiProfile.DELETE("/delete/:id", userHandler.DeleteUser)
		// 更新用户信息，仅限手机号、邮箱
		apiProfile.PUT("/update/:id", userHandler.UpdateUserInfo)
		// 查询是否被录取
		apiProfile.GET("/status/:id", userHandler.GetUserStatus)
	}

	apiJob := r.Group("/recruitment")
	apiJob.Use()
	{
		//查看岗位总览
		apiJob.GET("/jobs", jobHandler.GetJobs)
		//查看岗位详细信息
		apiJob.GET("/job/:id", jobHandler.GetJobInfo)
		//申请投递
		apiJob.POST("/job/:id/apply")
	}

	apiAdmin := r.Group("/admin")
	apiAdmin.Use() // 使用JWT中间件进行管理员身份验证
	{
		// 管理员仪表板或主页
		apiAdmin.GET("/dashboard", adminHandler.HomePage)

		// 查看所有职位发布（管理员）
		apiAdmin.GET("/jobs", adminHandler.ShowAllJobs)

		// 查看职位详情（管理员）
		apiAdmin.GET("/job/:jobID", adminHandler.ShowJobDetails)

		// 查看职位申请（管理员）
		apiAdmin.GET("/applications/:jobID", adminHandler.ShowJobApply)

		// 审批或拒绝职位申请（管理员）
		apiAdmin.PUT("/application/:appID", adminHandler.ShowJobApply)
	}

	err = r.Run(":8080")

	if err != nil {
		logger.Logger.Logf(logger.FatalLevel, "falied to start : %v", err)
	}
}
