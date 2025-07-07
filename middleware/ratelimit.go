package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)


type RateLimitStore struct {
	mu       sync.RWMutex
	clients  map[string]*ClientInfo
	cleanup  time.Duration
	lastClean time.Time
}


type ClientInfo struct {
	Requests  int
	LastReset time.Time
	Blocked   bool
	BlockedUntil time.Time
}


type RateLimitConfig struct {
	Max        int           
	Window     time.Duration 
	BlockTime  time.Duration 
	Message    string        
	SkipPaths  []string     
}


var globalStore = &RateLimitStore{
	clients:   make(map[string]*ClientInfo),
	cleanup:   5 * time.Minute,
	lastClean: time.Now(),
}


func RateLimit(config RateLimitConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		
		path := c.Path()
		for _, skipPath := range config.SkipPaths {
			if path == skipPath {
				return c.Next()
			}
		}

		
		clientIP := getClientIP(c)
		
		
		globalStore.cleanOldEntries()

		
		if !globalStore.checkRateLimit(clientIP, config) {
			message := config.Message
			if message == " " {
				message = "Demasiadas solicitudes. Intenta de nuevo más tarde."
			}
			
			return c.Status(429).JSON(fiber.Map{
				"error":     message,
				"retry_after": int(config.BlockTime.Seconds()),
				"limit":     config.Max,
				"window":    config.Window.String(),
			})
		}

		return c.Next()
	}
}


func (store *RateLimitStore) checkRateLimit(clientIP string, config RateLimitConfig) bool {
	store.mu.Lock()
	defer store.mu.Unlock()

	now := time.Now()
	client, exists := store.clients[clientIP]

	if !exists {
		
		store.clients[clientIP] = &ClientInfo{
			Requests:  1,
			LastReset: now,
			Blocked:   false,
		}
		return true
	}

	
	if client.Blocked && now.Before(client.BlockedUntil) {
		return false
	}

	
	if now.Sub(client.LastReset) >= config.Window {
		client.Requests = 1
		client.LastReset = now
		client.Blocked = false
		return true
	}

	
	client.Requests++

	
	if client.Requests > config.Max {
		client.Blocked = true
		client.BlockedUntil = now.Add(config.BlockTime)
		return false
	}

	return true
}


func (store *RateLimitStore) cleanOldEntries() {
	store.mu.Lock()
	defer store.mu.Unlock()

	now := time.Now()
	if now.Sub(store.lastClean) < store.cleanup {
		return
	}

	for ip, client := range store.clients {
		
		if now.Sub(client.LastReset) > time.Hour {
			delete(store.clients, ip)
		}
	}

	store.lastClean = now
}


func getClientIP(c *fiber.Ctx) string {
	
	if ip := c.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := c.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := c.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}
	return c.IP()
}


func AuthRateLimit() fiber.Handler {
	return RateLimit(RateLimitConfig{
		Max:       5,                    
		Window:    15 * time.Minute,     
		BlockTime: 2 * time.Minute,      
		Message:   "Demasiados intentos de autenticación. Intenta en 2 minutos.",
	})
}


func RegisterRateLimit() fiber.Handler {
	return RateLimit(RateLimitConfig{
		Max:       3,                    
		Window:    time.Hour,            
		BlockTime: 2 * time.Minute,      
		Message:   "Demasiados registros. Intenta en 2 minutos.",
	})
}


func GeneralAPIRateLimit() fiber.Handler {
	return RateLimit(RateLimitConfig{
		Max:       100,                  
		Window:    time.Minute,          
		BlockTime: 2 * time.Minute,      
		Message:   "Límite de API excedido. Intenta en 2 minutos.",
		SkipPaths: []string{"/health"},  
	})
}


func AdminRateLimit() fiber.Handler {
	return RateLimit(RateLimitConfig{
		Max:       50,                   
		Window:    time.Minute,          
		BlockTime: 2 * time.Minute,     
		Message:   "Límite de operaciones administrativas excedido. Intenta en 2 minutos.",
	})
}


func MedicalRateLimit() fiber.Handler {
	return RateLimit(RateLimitConfig{
		Max:       200,                 
		Window:    time.Minute,         
		BlockTime: 2 * time.Minute,      
		Message:   "Límite de operaciones médicas excedido. Intenta en 2 minutos.",
	})
}

func GetRateLimitStatus(clientIP string) map[string]interface{} {
	globalStore.mu.RLock()
	defer globalStore.mu.RUnlock()

	client, exists := globalStore.clients[clientIP]
	if !exists {
		return map[string]interface{}{
			"requests": 0,
			"blocked":  false,
		}
	}

	return map[string]interface{}{
		"requests":     client.Requests,
		"last_reset":   client.LastReset,
		"blocked":      client.Blocked,
		"blocked_until": client.BlockedUntil,
	}
}