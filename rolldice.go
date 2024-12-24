package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const name = "go.opentelemetry.io/otel/example/dice"

var (
	tracer     = otel.Tracer(name)
	meter      = otel.Meter(name)
	logger     = otelslog.NewLogger(name)
	rollCnt    metric.Int64Counter
	logs       []string
	logsMutex  sync.Mutex
	playerRoll = make(map[string]int)
)

func init() {
	var err error
	rollCnt, err = meter.Int64Counter("dice.rolls", metric.WithDescription("The number of dice rolls by value"))
	if err != nil {
		panic(err)
	}
}

func rolldice(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "roll")
	defer span.End()

	player := r.URL.Query().Get("player")
	if player == "" {
		player = "Anonymous"
	}

	roll := 1 + rand.Intn(6)
	rollValueAttr := attribute.Int("roll.value", roll)
	span.SetAttributes(rollValueAttr)
	rollCnt.Add(ctx, 1, metric.WithAttributes(rollValueAttr))

	msg := fmt.Sprintf("%s rolled a %d 🎲", player, roll)
	logger.InfoContext(ctx, msg)
	addLog(msg)

	resp := fmt.Sprintf("Hi %s, your dice roll is: %d\n", player, roll)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, resp)
}

func viewLogs(w http.ResponseWriter, r *http.Request) {
	logsMutex.Lock()
	defer logsMutex.Unlock()
	for _, logEntry := range logs {
		io.WriteString(w, logEntry+"\n")
	}
}

func resetMetrics(w http.ResponseWriter, r *http.Request) {
	playerRoll = make(map[string]int)
	io.WriteString(w, "Metrics have been reset!\n")
}

func addLog(entry string) {
	logsMutex.Lock()
	defer logsMutex.Unlock()
	logs = append(logs, entry)
}