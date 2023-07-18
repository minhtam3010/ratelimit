package api

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Limiter struct {
	RPS     float64 // Request per second
	Burst   int     // How many requests can be handled at once
	Enabled bool    // Enable rate limiter
}

func ExceedRequestsLimit(handler gin.HandlerFunc, limit Limiter) gin.HandlerFunc {
	type Client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*Client)
	)

	// Go routine to remove old entries from the map
	go func() {
		for {
			time.Sleep(time.Minute) // Go routine will sleep within 1 minute before removing old entries
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return gin.HandlerFunc(func(c *gin.Context) {
		if limit.Enabled {
			// Get the IP Address of the request client
			ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			log.Println("IP Address: ", ip)

			// Lock()
			mu.Lock()

			// Perform our processing logic. If the IP address is in the map otherwise we add it into our map
			if _, exists := clients[ip]; !exists {
				clients[ip] = &Client{
					limiter: rate.NewLimiter(rate.Limit(limit.RPS), limit.Burst),
				}
			}

			clients[ip].lastSeen = time.Now()
			// Check if there is still room for the request allowed
			if !clients[ip].limiter.Allow() {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
				mu.Unlock()
				return
			}

			// Unlock()
			mu.Unlock()
		}
		handler(c)
	})
}
