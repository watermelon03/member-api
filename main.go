package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/watermelon03/member-api/controllers"
	"github.com/watermelon03/member-api/middlewere"

	_ "github.com/go-sql-driver/mysql"
)

const basePath = "/member-api"

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	controllers.SetupDB()

	r := gin.Default()
	// r.Use(controllers.Cors())
	// r.Use(cors.Default())
	config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    config.AllowHeaders = []string{"Authorization", "Content-Type"} // Use "Authorization" instead of "authorization"
    r.Use(cors.New(config))

	r.GET("/ping", func(c *gin.Context) {
		fmt.Println(controllers.DB)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	baseGroup := r.Group(basePath)
	{
		baseGroup.POST("/admin/register", controllers.RegisterAdmin())
		baseGroup.POST("/admin/login", controllers.LoginAdmin())

		baseGroup.POST("/user/register", controllers.RegisterUser())
		baseGroup.POST("/user/login", controllers.LoginUser())

		baseGroup.GET("/testHash", controllers.TestHash())
		baseGroup.POST("/uploadForm", controllers.UploadFormData())
		baseGroup.POST("/getdataB", controllers.GetDataB())
	}

	adminGroup := r.Group(basePath+"/admin", middlewere.Authen())
	{
		adminGroup.GET("/all", controllers.GetAdminAll())
		adminGroup.GET("/profile", controllers.GetAdminProfile())
		adminGroup.PUT("/profile/update/password", controllers.UpdateAdminPassword())
		adminGroup.PUT("/profile/update/info", controllers.UpdateAdminInfo())
		adminGroup.GET("/getUser", controllers.GetUserAll())
	}

	userGroup := r.Group(basePath+"/user", middlewere.Authen())
	{
		userGroup.GET("/profile", controllers.GetUserProfile())
		userGroup.PUT("/profile/update/image", controllers.UpdateUserImage())
		userGroup.PUT("/profile/update/password", controllers.UpdateUserPassword())
		userGroup.PUT("/profile/update/info", controllers.UpdateUserInfo())
	}

	r.Run(":5050")
}
