package backend

import (
	"github.com/NubeIO/git/pkg/git"
	"github.com/NubeIO/ros-connect/backend/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (inst *App) getConfigHandler(c *gin.Context) {
	c.JSON(http.StatusOK, inst.config.GetConfig())
}

func (inst *App) gitTokenHandler(c *gin.Context) {
	decoded := inst.config.GetGitToken()
	c.JSON(http.StatusOK, gin.H{"token": decoded})
}

func (inst *App) updateGitTokenHandler(c *gin.Context) {
	var req *config.Config

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	encoded := git.EncodeToken(req.GitToken)
	con := inst.config.GetConfig()
	con.GitToken = encoded
	if err := inst.config.SaveConfig(con); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Config updated successfully"})
}

func (inst *App) updateConfigHandler(c *gin.Context) {
	var req *config.Config

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := inst.config.SaveConfig(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Config updated successfully"})
}
