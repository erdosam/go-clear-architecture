package v1

import (
	"fmt"
	"github.com/erdosam/go-clear-architecture/internal/usecase"
	"github.com/gin-gonic/gin"
)

// pingRoute is for testing only
type pingRoute struct {
	usecase.Ping
}

func (rh *routerHandler) initPingRoutes(parent *gin.RouterGroup) {
	r := &pingRoute{rh.feature.Ping}

	parent.GET("/pub-sub", r.testSendingPubsub)
}

func (r *pingRoute) testSendingPubsub(c *gin.Context) {
	msg := c.Query("message")
	err := r.TestPubsub(msg)
	if err != nil {
		errorJSON(c, 411, err, 0)
		return
	}
	detailJSON(c, gin.H{"message": fmt.Sprintf("Message '%s' has been sent", msg)})
}
