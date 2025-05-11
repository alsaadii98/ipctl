package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alsaadii98/ipctl/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type APIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ApiIpRequest struct {
	TelegramUsername  string `json:"telegram_username"`
	TelegramChatId    string `json:"telegram_chat_id"`
	TelegramFirstName string `json:"telegram_first_name"`
	TelegramLastName  string `json:"telegram_last_name"`
	IpAddress         string `json:"ip_address"`
	Note              string `json:"note"`
}

type ValidateIpAndTelegramUserRequest struct {
	Username string `json:"username"`
}

func sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIResponse{
		Status:  statusCode,
		Message: message,
	})
}

func sendSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"status":  statusCode,
		"message": message,
	}
	if data != nil {
		response["data"] = data
	}
	json.NewEncoder(w).Encode(response)
}

func Validate(w http.ResponseWriter, r *http.Request) {
	var req ValidateIpAndTelegramUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	var userId int
	query := "SELECT id FROM users WHERE telegram_username = ?"
	err := config.DB.QueryRow(query, req.Username).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			sendError(w, http.StatusNotFound, "User not found")
		} else {
			sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	sendSuccess(w, http.StatusOK, "User found", map[string]interface{}{
		"user_id": userId,
	})
}

func AddIp(w http.ResponseWriter, r *http.Request) {
	var req ApiIpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	var (
		userID           int
		telegramUsername string
	)

	queryUser := "SELECT id, telegram_username FROM users WHERE telegram_username = ?"
	err := config.DB.QueryRow(queryUser, req.TelegramUsername).Scan(&userID, &telegramUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			sendError(w, http.StatusNotFound, "User not found")
		} else {
			sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	queryUpdateUser := `
		UPDATE users
		SET telegram_first_name = ?, telegram_last_name = ?, telegram_chat_id = ?
		WHERE id = ?
	`
	_, err = config.DB.Exec(queryUpdateUser, req.TelegramFirstName, req.TelegramLastName, req.TelegramChatId, userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to update user info")
		return
	}

	note := req.Note
	if note == "" {
		note = fmt.Sprintf("User %s added IP %s at %s", telegramUsername, req.IpAddress, time.Now().Format(time.RFC3339))
	}

	insertQuery := `
		INSERT INTO ip_addresses (user_id, ip_address, note, created_at)
		VALUES (?, ?, ?, ?)
	`
	_, err = config.DB.Exec(insertQuery, userID, req.IpAddress, note, time.Now())
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to add IP")
		return
	}

	sendSuccess(w, http.StatusCreated, "IP address added successfully", nil)
}

func DeleteExpiredIPs(w http.ResponseWriter, r *http.Request) {
	cutoff := time.Now().Add(-24 * time.Hour)

	deleteQuery := `
		DELETE FROM ip_addresses
		WHERE created_at < ?
		RETURNING user_id, ip_address
	`
	rows, err := config.DB.Query(deleteQuery, cutoff)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to delete expired IP addresses")
		return
	}
	defer rows.Close()

	var userId int
	var ipAddress string
	var telegramChatId string
	var message string

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_API_KEY"))
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to connect to Telegram")
		return
	}

	for rows.Next() {
		err := rows.Scan(&userId, &ipAddress)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "Error scanning rows")
			return
		}

		queryUser := "SELECT telegram_chat_id FROM users WHERE id = ?"
		err = config.DB.QueryRow(queryUser, userId).Scan(&telegramChatId)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			} else {
				sendError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		message = fmt.Sprintf("Your IP address %s has been removed because it was added over 24 hours ago.", ipAddress)
		msg := tgbotapi.NewMessageToChannel(telegramChatId, message)
		_, err = bot.Send(msg)
		if err != nil {
			fmt.Printf("Failed to send message to user %d: %v\n", userId, err)
		}
	}

	sendSuccess(w, http.StatusOK, "Expired IP addresses deleted and notifications sent", nil)
}
