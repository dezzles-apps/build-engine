package main
import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"dezzles-apps/build-engine/model"
	"dezzles-apps/build-engine/repository"
	"dezzles-apps/build-engine/service"
)

var eventService *service.EventService = &service.EventService{}

func main() {
	
	config, err := model.LoadConfig()
	if err != nil {
		panic(err)
	}
	repository.Connect(&config.Database)
	eventService.Initialise(&config.EventService)
	router := gin.Default()
	router.GET("/api/v1/builds", getAllBuilds)
	router.GET("/api/v1/builds/:org/:repository", getBuildsByRepository)
	router.GET("/api/v1/builds/:org/:repository/:build_number", getBuild)
	router.GET("/api/v1/repositories/:org/:repository", getRepositoryConfiguration)
	router.POST("/api/v1/events", notifyEvent)
	router.Run(":8080")
}

func getAllBuilds(c *gin.Context) {
	builds, err := repository.GetBuilds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, builds)
}

func getBuildsByRepository(c *gin.Context) {
	org := c.Param("org")
	repositoryName := c.Param("repository")
	builds, err := repository.GetBuildsByRepository(org, repositoryName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, builds)
}

func getRepositoryConfiguration(c *gin.Context) {
	org := c.Param("org")
	repositoryName := c.Param("repository")
	config, err := repository.GetRepositoryConfiguration(org, repositoryName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, config)
}

func getBuild(c * gin.Context) {
	org := c.Param("org")
	repositoryName := c.Param("repository")
	buildNumberStr := c.Param("build_number")
	buildNumber, err := strconv.Atoi(buildNumberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid build number"})
		return
	}
	build, err := repository.GetBuild(org, repositoryName, buildNumber)
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

	_, err := service.ValidateBuild(event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	eventService.SendEvent(event)
	c.JSON(http.StatusOK, gin.H{"success": true, "event": event})
}