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
	Auth      *controller.AuthController
	Profile   *controller.ProfileController
	Swipe     *controller.SwipeController
	Subscribe *controller.SubscriptionController
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
	// Repository components
	usr := repository.NewUserRepository(c.DB)
	pfr := repository.NewProfileRepository(c.DB)
	pmr := repository.NewPremiumFeatureRepository(c.DB)
	sbr := repository.NewSubscriptionRepository(c.DB)
	swr := repository.NewSwipeRepository(c.DB)

	// Service components
	pfs := service.NewProfileService(pfr, sbr)
	aus := service.NewAuthService(usr, pfr)
	sbs := service.NewSubscriptionService(sbr, pmr)
	sws := service.NewSwipeService(swr, sbr, 10)

	// Controller components
	c.Controllers.Profile = controller.NewProfileController(pfs)
	c.Controllers.Auth = controller.NewAuthController(aus)
	c.Controllers.Subscribe = controller.NewSubscriptionController(sbs)
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
