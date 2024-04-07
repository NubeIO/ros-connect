package backend

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func (inst *App) SetupRoutes(router *gin.Engine) {
	router.GET("/api/hosts", inst.getAllInstancesHandler)
	router.GET("/api/hosts/:uuid", inst.getInstanceHandler)
	router.POST("/api/hosts", inst.createInstanceHandler)
	router.DELETE("/api/hosts/:uuid", inst.deleteInstanceHandler)
	router.PATCH("/api/hosts/:uuid", inst.updateInstanceHandler)

	router.GET("/ws", inst.webSocketHandler)

	router.Any("/proxy/*path", inst.ProxyHandler)
	router.Any("/ws/proxy/*path", inst.ProxyHandlerWS)

	router.GET("/api/hosts/polling", inst.getAllPollingHandler)
	router.POST("/api/hosts/polling", inst.createPollingHandler)
	router.DELETE("/api/hosts/polling/:uuid", inst.deletePollingHandler)

	router.GET("/api/config", inst.getConfigHandler)
	router.POST("/api/config", inst.updateConfigHandler)
	router.GET("/api/config/git/token", inst.gitTokenHandler)
	router.POST("/api/config/git/token", inst.updateGitTokenHandler)

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

}

func (inst *App) createInstanceHandler(c *gin.Context) {
	var instance *Instance
	if err := c.BindJSON(&instance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if instance, err := inst.CreateInstance(instance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": instance})
	}

}

func (inst *App) getAllInstancesHandler(c *gin.Context) {
	instances := inst.GetAllInstances()
	c.JSON(http.StatusOK, instances)
}

func (inst *App) deleteInstanceHandler(c *gin.Context) {
	id := c.Param("uuid")
	if err := inst.DeleteInstance(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "instance deleted successfully"})
}

func (inst *App) getInstanceHandler(c *gin.Context) {
	id := c.Param("uuid")
	instance, err := inst.GetInstance(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": instance})
}

func (inst *App) updateInstanceHandler(c *gin.Context) {
	var instance *Instance
	id := c.Param("uuid")
	if err := c.BindJSON(&instance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	instance.UUID = id
	if err := inst.UpdateInstance(instance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "instance updated successfully"})
}

func (inst *App) ProxyHandler(c *gin.Context) {
	name := c.GetHeader("Instance-Name")
	instance, exists := inst.Instances[name]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instance not found"})
		return
	}

	target := fmt.Sprintf("%s:%d", instance.IP, instance.Port)
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   target,
	})

	c.Request.URL.Path = c.Param("path") // Include the endpoint path
	proxy.ServeHTTP(c.Writer, c.Request)
}

func newUUID() string {
	id := uuid.New()
	return strings.ReplaceAll(id.String(), "-", "")
}
