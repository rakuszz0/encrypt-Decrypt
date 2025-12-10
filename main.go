package main

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

type EncryptRequest struct {
	PlainText string `json:"plainText" binding:"required"`
	Key       string `json:"key" binding:"required"`
}

type DecryptRequest struct {
	EncryptedText string `json:"encryptedText" binding:"required"`
	Key           string `json:"key" binding:"required"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type KeyResponse struct {
	Key string `json:"key"`
}

type EncryptResponse struct {
	EncryptedText string            `json:"encryptedText"`
	Mapping       map[string]string `json:"mapping,omitempty"`
}

type DecryptResponse struct {
	DecryptedText string `json:"decryptedText"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	router.Use(cors.New(config))

	router.GET("/api/generate-key", generateKeyHandler)
	router.POST("/api/encrypt", encryptHandler)
	router.POST("/api/decrypt", decryptHandler)

	router.Run(":8080")
}

func generateKey() string {
	chars := strings.Split(alphabet, "")

	for i := len(chars) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		chars[i], chars[j] = chars[j], chars[i]
	}

	return strings.Join(chars, "")
}

func isValidKey(key string) bool {
	if len(key) != 26 {
		return false
	}

	key = strings.ToLower(key)
	seen := make(map[rune]bool)

	for _, char := range key {
		if char < 'a' || char > 'z' {
			return false
		}
		if seen[char] {
			return false
		}
		seen[char] = true
	}

	return len(seen) == 26
}

func generateMapping(key string) map[string]string {
	mapping := make(map[string]string)

	for i, char := range alphabet {
		mapping[string(char)] = string(key[i])
		mapping[strings.ToUpper(string(char))] = strings.ToUpper(string(key[i]))
	}

	return mapping
}

func generateReverseMapping(key string) map[string]string {
	reverseMapping := make(map[string]string)

	for i, char := range key {
		reverseMapping[string(char)] = string(alphabet[i])
		reverseMapping[strings.ToUpper(string(char))] = strings.ToUpper(string(alphabet[i]))
	}

	return reverseMapping
}

func encryptText(plainText, key string) (string, map[string]string) {
	mapping := generateMapping(key)
	var result strings.Builder

	for _, char := range plainText {
		strChar := string(char)

		if strChar == " " {
			result.WriteString(" ")
			continue
		}

		if encryptedChar, exists := mapping[strChar]; exists {
			result.WriteString(encryptedChar)
		} else {
			result.WriteString(strChar)
		}
	}

	return result.String(), mapping
}

func decryptText(encryptedText, key string) string {
	reverseMapping := generateReverseMapping(key)
	var result strings.Builder

	for _, char := range encryptedText {
		strChar := string(char)

		if strChar == " " {
			result.WriteString(" ")
			continue
		}

		if decryptedChar, exists := reverseMapping[strChar]; exists {
			result.WriteString(decryptedChar)
		} else {
			result.WriteString(strChar)
		}
	}

	return result.String()
}

func generateKeyHandler(c *gin.Context) {
	key := generateKey()

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: KeyResponse{
			Key: key,
		},
	})
}

func encryptHandler(c *gin.Context) {
	var req EncryptRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	if !isValidKey(req.Key) {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Key must be 26 unique lowercase letters",
		})
		return
	}

	encryptedText, mapping := encryptText(req.PlainText, req.Key)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: EncryptResponse{
			EncryptedText: encryptedText,
			Mapping:       mapping,
		},
	})
}

func decryptHandler(c *gin.Context) {
	var req DecryptRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	if !isValidKey(req.Key) {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Key must be 26 unique lowercase letters",
		})
		return
	}

	decryptedText := decryptText(req.EncryptedText, req.Key)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: DecryptResponse{
			DecryptedText: decryptedText,
		},
	})
}
