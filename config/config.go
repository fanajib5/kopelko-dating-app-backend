package config

import (
	"log"
	"os"
	"strconv"

	"kopelko-dating-app-backend/controllers"
	"kopelko-dating-app-backend/repositories"
	"kopelko-dating-app-backend/services"
	"kopelko-dating-app-backend/utils"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Config struct {
	DB         *gorm.DB
	JWTKey     []byte
	APIPort    string
	LimitSwipe int
	Controllers
}

const defaultLimitSwipe = 10

type Controllers struct {
	Auth      *controllers.AuthController
	Profile   *controllers.ProfileController
	Swipe     *controllers.SwipeController
	Subscribe *controllers.SubscriptionController
}

func New() *Config {
	loadEnv()

	var c = new(Config)
	c.initializeDB()
	c.initializeControllers()
	c.LoadAPIPort()
	c.LoadLimitSwipe()
	c.LoadJWTKey()

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

func (c *Config) LoadLimitSwipe() {
	limitSwipe := os.Getenv("LIMIT_SWIPE")
	limitSwipeInt, err := strconv.Atoi(limitSwipe)
	if err != nil {
		limitSwipeInt = defaultLimitSwipe
		log.Printf("could not convert LIMIT_SWIPE to int: %v, falling back to default value: %d", err, defaultLimitSwipe)
	}
	c.LimitSwipe = limitSwipeInt
}

// Initialize controllers
func (c *Config) initializeControllers() {
	// repositories. components
	usr := repositories.NewUserRepository(c.DB)
	pfr := repositories.NewProfileRepository(c.DB)
	pmr := repositories.NewPremiumFeatureRepository(c.DB)
	sbr := repositories.NewSubscriptionRepository(c.DB)
	swr := repositories.NewSwipeRepository(c.DB)
	pvr := repositories.NewProfileViewRepository(c.DB)

	// Service components
	pfs := services.NewProfileService(pfr, pvr, sbr, c.LimitSwipe)
	aus := services.NewAuthService(usr, pfr)
	sbs := services.NewSubscriptionService(sbr, pmr, pfr)
	sws := services.NewSwipeService(swr, sbr, c.LimitSwipe)

	// Controller components
	c.Controllers.Profile = controllers.NewProfileController(pfs)
	c.Controllers.Auth = controllers.NewAuthController(aus)
	c.Controllers.Subscribe = controllers.NewSubscriptionController(sbs)
	c.Controllers.Swipe = controllers.NewSwipeController(sws)
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

func NewValidator() *utils.CustomValidator {
	return utils.NewValidator()
}
