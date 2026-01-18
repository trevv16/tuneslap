package app

import (
	"fmt"
	"os"
	"strings"
	"time"

	"tuneslap/config"
	"tuneslap/database"
	"tuneslap/router"
	"tuneslap/tasks"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// SetupApp configures and returns a Fiber app instance without starting the server
func SetupApp() (*fiber.App, error) {
	// load env
	err := config.LoadENV()
	if err != nil {
		return nil, fmt.Errorf("error loading env: %w", err)
	}

	// validate required config
	err = config.ValidateRequiredConfig()
	if err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// start database
	err = database.StartMongoDB()
	if err != nil {
		return nil, fmt.Errorf("error starting MongoDB: %w", err)
	}

	// start redis
	database.InitRedis()

	// start asynq client
	tasks.InitClient()

	// create app
	app := fiber.New()

	// configure timeouts
	app.Server().ReadTimeout = 10 * time.Second
	app.Server().WriteTimeout = 10 * time.Second

	// attach middleware
	attachMiddleware(app)

	// setup routes
	router.SetupRoutes(app)

	return app, nil
}

// SetupAndRunApp sets up the app and starts the server
func SetupAndRunApp() error {
	app, err := SetupApp()
	if err != nil {
		return err
	}

	// defer closing services
	defer func() {
		database.CloseMongoDB()
		database.CloseRedis()
		tasks.CloseClient()
	}()

	// get the port and start
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // default port
	}

	return app.Listen(":" + port)
}

// attachMiddleware attaches all middleware to the Fiber app
func attachMiddleware(app *fiber.App) {
	// recover middleware
	app.Use(recover.New())

	// logger middleware
	app.Use(logger.New(logger.Config{
		Format:     "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "America/New_York",
	}))

	// cors middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://localhost:8081, https://tuneslap.com, https://*.tuneslap.com",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// rate limit middleware
	app.Use(limiter.New(limiter.Config{
		Max:               30,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	// compress responses middleware
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// cache control middleware
	app.Use(cache.New(cache.Config{
		Expiration: 10 * time.Second,
	}))

	// helmet middleware
	app.Use(helmet.New())

	// healthcheck middleware
	app.Use(healthcheck.New())

	// accept json on api routes middleware
	app.Use(func(c *fiber.Ctx) error {
		// if path contains /api
		if strings.Contains(c.Path(), "/api") {
			c.Accepts("application/json")
			return c.Next()
		}
		return c.Next()
	})
}
