package main

import (
	"dezzles-apps/build-engine/model"
	"dezzles-apps/build-engine/repository"
	"dezzles-apps/build-engine/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var eventService *service.EventService = &service.EventService{}
var buildService *service.BuildService = &service.BuildService{}
var portfolioRepository *repository.PortfolioRepository = &repository.PortfolioRepository{}
var database *repository.Database = &repository.Database{}
var eventRepository *repository.EventsRepository = &repository.EventsRepository{}

func main() {

	config, err := model.LoadConfig()
	if err != nil {
		panic(err)
	}
	database.Connect(&config.Database)
	eventRepository.Initialise(database)
	eventService.Initialise(&config.EventService)
	buildService.Initialise(eventRepository)
	portfolioRepository.Initialise(database)
	router := gin.Default()
	router.GET("/api/v1/builds", getAllBuilds)
	router.GET("/api/v1/builds/:org/:repository", getBuildsByRepository)
	router.GET("/api/v1/builds/:org/:repository/:build_number", getBuild)
	router.GET("/api/v1/repositories/:org/:repository", getRepositoryConfiguration)
	router.POST("/api/v1/events", notifyEvent)
	router.POST("/api/v1/publish", publish)
	router.Run(":8080")
}

func getAllBuilds(c *gin.Context) {
	builds, err := eventRepository.GetBuilds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, builds)
}

func getBuildsByRepository(c *gin.Context) {
	org := c.Param("org")
	repositoryName := c.Param("repository")
	builds, err := eventRepository.GetBuildsByRepository(org, repositoryName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, builds)
}

func getRepositoryConfiguration(c *gin.Context) {
	org := c.Param("org")
	repositoryName := c.Param("repository")
	config, err := eventRepository.GetRepositoryConfiguration(org, repositoryName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, config)
}

func getBuild(c *gin.Context) {
	org := c.Param("org")
	repositoryName := c.Param("repository")
	buildNumberStr := c.Param("build_number")
	buildNumber, err := strconv.Atoi(buildNumberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid build number"})
		return
	}
	build, err := eventRepository.GetBuild(org, repositoryName, buildNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, build)
}

func notifyEvent(c *gin.Context) {
	var event model.BuildEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := event.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := buildService.ValidateBuild(event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	eventService.SendEvent(event)
	c.JSON(http.StatusOK, gin.H{"success": true, "event": event})
}

func publish(c *gin.Context) {
	var message model.PublishEvent
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := message.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := portfolioRepository.UpdatePortfolio(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": message})
}
