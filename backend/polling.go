package backend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (inst *App) getAllPollingHandler(c *gin.Context) {
	pollingInstances := inst.AllPolling()
	c.JSON(http.StatusOK, pollingInstances)
}

func (inst *App) createPollingHandler(c *gin.Context) {
	var req struct {
		InstanceUUID string `json:"instanceUUID"`
		Timeout      int    `json:"timeout"`
		RefreshRate  int    `json:"refreshRate"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.RefreshRate <= 0 {
		req.RefreshRate = 5
	}

	uuid := inst.CreatePolling(req.InstanceUUID, req.Timeout, time.Duration(req.RefreshRate)*time.Second, nil)
	if uuid == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "instance not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"uuid": uuid})
}

func (inst *App) deletePollingHandler(c *gin.Context) {
	uuid := c.Param("uuid")
	deletedUUID := inst.DeletePolling(uuid)
	if deletedUUID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "polling instance not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted_uuid": deletedUUID})
}

func (inst *App) AllPolling() []*PollingInstance {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()

	pollingInstances := make([]*PollingInstance, 0, len(inst.polling))
	for _, pollingInstance := range inst.polling {
		pollingInstances = append(pollingInstances, pollingInstance)
	}

	return pollingInstances
}

func (inst *App) CreatePolling(instanceUUID string, timeout int, refreshRate time.Duration, opts *PollingOptions) string {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()

	_, ok := inst.Instances[instanceUUID]
	if !ok {
		return ""
	}
	id := newUUID()
	pollingInstance := &PollingInstance{
		UUID:         id,
		InstanceUUID: instanceUUID,
		Timeout:      timeout,
		RefreshRate:  refreshRate,
		Options:      opts,
		stopChan:     make(chan struct{}),
		counter:      0,
	}

	go func() {
		ticker := time.NewTicker(refreshRate)
		for {
			select {
			case <-ticker.C:
				// Perform polling logic here
				pollingInstance.counter++
				fmt.Printf("Poll uuid: %s instance: %s count: %d \n", id, instanceUUID, pollingInstance.counter)
				if pollingInstance.counter >= timeout/int(refreshRate.Seconds()) {
					inst.DeletePolling(id)
					return
				}
			case <-pollingInstance.stopChan:
				fmt.Printf("Stop uuid: %s instance: %s count: %d \n", id, instanceUUID, pollingInstance.counter)
				ticker.Stop()
				return
			}
		}
	}()
	inst.polling[id] = pollingInstance
	return id
}

func (inst *App) DeletePolling(uuid string) string {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	pollingInstance, ok := inst.polling[uuid]
	if !ok {
		return ""
	}
	close(pollingInstance.stopChan)
	delete(inst.polling, uuid)

	return uuid
}
