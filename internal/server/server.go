package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	addrServer   = ":7878"
	readTimeout  = 10 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 120 * time.Second
)

// RunServerWithGracefulShutdown - запускает сервер с возможностью корректного завершения работы при получении сигналов
// остановки.
func RunServerWithGracefulShutdown(mux *http.ServeMux) {
	server := setupServer(mux)

	serverErr := make(chan error, 1)
	go startServer(server, serverErr)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		log.Println("Получен сигнал завершения!")
	case err := <-serverErr:
		log.Printf("Ошибка сервера: %v", err)
	}

	shutdownServer(server)
}

// setupServer - создает и настраивает HTTP-сервер с предопределенными параметрами.
func setupServer(mux *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:         addrServer,
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
}

// startServer - функция для асинхронного запуска сервера с обработкой ошибок.
func startServer(server *http.Server, serverErr chan<- error) {
	log.Printf("Сервер запущен на %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		serverErr <- err
	}
	close(serverErr)
}

// shutdownServer - реализует безопасное и корректное завершение работы сервера.
func shutdownServer(server *http.Server) {
	log.Println("Начинаем отключение...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Предупреждение: отключение не удалось: %v", err)
	}

	log.Println("Сервер остановлен")
}
