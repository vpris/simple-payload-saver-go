package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

var payloadList []json.RawMessage

func main() {
	// Checking if there is a file with saved payloads
	if _, err := os.Stat("payloads.json"); err == nil {
		loadPayloadsFromFile()
	}

	router := gin.Default()

	router.POST("/webhook", handleWebhook)
	router.GET("/payloads", viewPayloads)

	router.Run(":8080")
}

func handleWebhook(c *gin.Context) {
	var payload json.RawMessage
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}
	payloadList = append(payloadList, payload)

	// We save payloadList to a file every time we update
	savePayloadsToFile()

	c.JSON(200, gin.H{"success": true})
}

func viewPayloads(c *gin.Context) {
	// Reading the list of saved payloads from a file
	payloadListFromFile := loadPayloadsFromFile()

	c.JSON(200, payloadListFromFile)
}

func savePayloadsToFile() {
	// Encoding payloadList in JSON
	data, err := json.Marshal(payloadList)
	if err != nil {
		fmt.Println("Failed to encode payloadList to JSON:", err)
		return
	}

	// Save to file payloads.json
	err = os.WriteFile("payloads.json", data, 0644)
	if err != nil {
		fmt.Println("Failed to write payloads to file:", err)
	}
}

func loadPayloadsFromFile() []json.RawMessage {
	// Reading data from the payloads.json file
	data, err := os.ReadFile("payloads.json")
	if err != nil {
		fmt.Println("Failed to read payloads file:", err)
		return nil
	}

	// Decoding data in payloadList
	err = json.Unmarshal(data, &payloadList)
	if err != nil {
		fmt.Println("Failed to decode payloads data:", err)
		return nil
	}

	return payloadList
}
