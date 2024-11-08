package config

import (
	"log"

	"kopelko-dating-app-backend/controllers"
	"kopelko-dating-app-backend/repositories"
	"kopelko-dating-app-backend/services"
	"kopelko-dating-app-backend/utils"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Config struct {
	DB     *gorm.DB
	JWTKey []byte
	Controllers
}

type Controllers struct {
	Auth    *controllers.AuthController
	Profile *controllers.ProfileController
	Swipe   *controllers.SwipeController
}

func New() *Config {
	loadEnv()

	var c = new(Config)
	c.initializeDB()
	c.initializeControllers()

	return c
}

// Initialize database connection
func (c *Config) initializeDB() {
	db, err := utils.InitDB()
	if err != nil {
		log.Fatalf("could not set up database: %v", err)
	}
	c.DB = db
}

func (c *Config) LoadJWTKey() {
	c.JWTKey = utils.LoadJWTKey()
}

// Initialize controllers
func (c *Config) initializeControllers() {
	// Profile component
	pfr := repositories.NewProfileRepository(c.DB)
	pfs := services.NewProfileService(pfr)
	c.Controllers.Profile = controllers.NewProfileController(pfs)

	// User component
	usr := repositories.NewUserRepository(c.DB)

	// Auth component
	aus := services.NewAuthService(usr, pfr)
	c.Controllers.Auth = controllers.NewAuthController(aus)

	// Swipe component
	swr := repositories.NewSwipeRepository(c.DB)
	sws := services.NewSwipeService(swr, 10)
	c.Controllers.Swipe = controllers.NewSwipeController(sws)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func NewValidator() *utils.CustomValidator {
	return utils.NewValidator()
}
