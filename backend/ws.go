package backend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func (inst *App) webSocketHandler(c *gin.Context) {
	conn, err := inst.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Add the client to the clients map
	inst.clients[conn] = true
	defer delete(inst.clients, conn)

	// Keep the connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (inst *App) sendWSMessage(data []byte) {
	for client := range inst.clients {
		if err := client.WriteMessage(websocket.BinaryMessage, data); err != nil {
			fmt.Println("Error sending message:", err)
			client.Close()
			delete(inst.clients, client)
		}
	}
}

func (inst *App) ProxyHandlerWS(c *gin.Context) {
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

	// Create a custom ResponseWriter to capture the response data
	responseCapture := &responseWriterCapture{ResponseWriter: c.Writer}
	proxy.ServeHTTP(responseCapture, c.Request)
	inst.sendWSMessage(responseCapture.body)
}

// responseWriterCapture is a custom http.ResponseWriter that captures the response status code and body
type responseWriterCapture struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rwc *responseWriterCapture) WriteHeader(statusCode int) {
	rwc.statusCode = statusCode
	rwc.ResponseWriter.WriteHeader(statusCode)
}

func (rwc *responseWriterCapture) Write(body []byte) (int, error) {
	rwc.body = append(rwc.body, body...)
	return rwc.ResponseWriter.Write(body)
}
