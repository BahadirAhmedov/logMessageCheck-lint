package linttest

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
)

// Эта функция может нигде не вызываться — это ок.
// Главное: файл должен компилироваться, чтобы линтер прошёл по AST.
func SampleLogs() {
	userID := 42
	email := "user@example.com"
	err := fmt.Errorf("db timeout: %w", ErrNotFound)

	// --- stdlib log: разные варианты сообщений ---

	// Потенциально "плохие" (часто линтеры ругаются на: заглавные буквы, точку в конце, лишние пробелы и т.п.)
	log.Println("User created.")                  // точка на конце
	log.Printf("Failed to create user: %v", err)  // заглавная буква
	log.Printf("  double  spaces  in message  ")  // двойные пробелы/пробелы по краям
	log.Printf("created user %d.", userID)        // точка
	log.Printf("created user: %d", userID)        // двоеточие (если запрещено)
	log.Printf("created user %d, email=%s", userID, email) // смешанный стиль

	// Потенциально "плохие", если твой линтер требует, чтобы message была литералом:
	msg := fmt.Sprintf("created user %d", userID)
	log.Println(msg) // сообщение не литерал

	// Более "хорошие" (обычно): маленькая буква, без точки, коротко, без мусора
	log.Printf("created user id=%d email=%s", userID, email)
	log.Printf("cannot create user id=%d err=%v", userID, err)
	log.Printf("request finished duration=%s", (1500 * time.Millisecond).String())

	// --- slog: структурированные логи (если твой линтер умеет slog) ---

	logger := slog.Default()

	// Потенциально "плохие"
	logger.Info("User created")                       
	logger.Info("created user.", "user_id", userID) 
	logger.Info(fmt.Sprintf("created user %d", userID))  // message не литерал

	// Более "хорошие"
	logger.Info("created user", "user_id", userID, "email", email)
	logger.Error("cannot create user", "user_id", userID, "err", err)
	logger.Warn("user not found", "user_id", userID, "err", ErrNotFound)
}