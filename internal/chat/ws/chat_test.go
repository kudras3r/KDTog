package ws

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kudras3r/KDTog/pkg/logger"
)

func TestWebSocketBroadcast(t *testing.T) {
	log := logger.New()
	SetLogger(log)

	hub := NewHub(log)
	go hub.Run()

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	}))
	defer server.Close()

	// Подключаем клиентов
	url := "ws" + server.URL[4:] + "/ws"
	clientCount := 3
	clients := make([]*websocket.Conn, clientCount)
	var err error

	for i := 0; i < clientCount; i++ {
		clients[i], _, err = websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("failed to connect client %d: %v", i+1, err)
		}
		defer clients[i].Close()
	}

	// Отправляем сообщение от первого клиента
	message := "hello from client 1"
	if err := clients[0].WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("failed to send message from client 1: %v", err)
	}

	// Проверяем, что остальные клиенты получили сообщение
	var wg sync.WaitGroup
	wg.Add(clientCount - 1)

	for i := 1; i < clientCount; i++ {
		go func(clientIndex int) {
			defer wg.Done()
			_, received, err := clients[clientIndex].ReadMessage()
			if err != nil {
				t.Errorf("failed to read message for client %d: %v", clientIndex+1, err)
				return
			}
			if string(received) != message {
				t.Errorf("client %d received incorrect message: expected %q, got %q", clientIndex+1, message, string(received))
			}
		}(i)
	}

	// Ждем завершения всех проверок
	wg.Wait()

	// Проверяем, что первый клиент не получил свое сообщение обратно
	clients[0].SetReadDeadline(time.Now().Add(1 * time.Second))
	_, _, err = clients[0].ReadMessage()
	if err == nil {
		t.Errorf("client 1 should not receive its own message")
	}
}
