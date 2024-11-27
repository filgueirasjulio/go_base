package app

import (
	"log"
	"sync"
	"time"
)

// Queue representa uma fila de mensagens
type Queue struct {
	messages chan string
	wg       sync.WaitGroup
}

var QueueInstance *Queue

// InitQueue inicializa a fila única
func InitQueue(capacity int) {
	if QueueInstance == nil {
		QueueInstance = NewQueue(capacity)
	}
}

// GetQueue retorna a instância da fila
func GetQueue() *Queue {
	return QueueInstance
}

// NewQueue retorna uma nova instância da fila
func NewQueue(capacity int) *Queue {
	return &Queue{
		messages: make(chan string, capacity),
	}
}

// AddMessage adiciona uma mensagem à fila
func (q *Queue) AddMessage(message string) {
	select {
	case q.messages <- message:
	default:
		log.Println("Fila cheia")
	}
}

// ProcessMessages processa as mensagens da fila
func (q *Queue) ProcessMessages(process func(string), interval time.Duration) {
	go func() {
		for message := range q.messages {
			process(message)
			time.Sleep(interval)
		}
	}()
}

// Wait aguarda o processamento das mensagens
func (q *Queue) Wait() {
	q.wg.Wait()
}

// Close fecha a fila
func (q *Queue) Close() {
	close(q.messages)
}