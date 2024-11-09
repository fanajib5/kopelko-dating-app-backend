package config

import (
	"log"
	"os"

	controller "kopelko-dating-app-backend/controllers"
	repository "kopelko-dating-app-backend/repositories"
	service "kopelko-dating-app-backend/services"
	util "kopelko-dating-app-backend/utils"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Config struct {
	DB      *gorm.DB
	JWTKey  []byte
	APIPort string
	Controllers
}

type Controllers struct {
	Auth    *controller.AuthController
	Profile *controller.ProfileController
	Swipe   *controller.SwipeController
}

func New() *Config {
	loadEnv()

	var c = new(Config)
	c.initializeDB()
	c.initializeControllers()
	c.LoadAPIPort()

	return c
}

// Initialize database connection
func (c *Config) initializeDB() {
	db, err := util.InitDB()
	if err != nil {
		log.Fatalf("could not set up database: %v", err)
	}
	c.DB = db
}

func (c *Config) LoadJWTKey() {
	c.JWTKey = util.LoadJWTKey()
}

// Initialize controllers
func (c *Config) initializeControllers() {
	// Profile component
	pfr := repository.NewProfileRepository(c.DB)
	pfs := service.NewProfileService(pfr)
	c.Controllers.Profile = controller.NewProfileController(pfs)

	// User component
	usr := repository.NewUserRepository(c.DB)

	// Auth component
	aus := service.NewAuthService(usr, pfr)
	c.Controllers.Auth = controller.NewAuthController(aus)

	// Swipe component
	swr := repository.NewSwipeRepository(c.DB)
	sws := service.NewSwipeService(swr, 10)
	c.Controllers.Swipe = controller.NewSwipeController(sws)
}

func (c *Config) LoadAPIPort() {
	c.APIPort = os.Getenv("API_PORT")
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func NewValidator() *util.CustomValidator {
	return util.NewValidator()
}
