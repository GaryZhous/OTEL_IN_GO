package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const name = "go.opentelemetry.io/otel/example/dice"

var (
	tracer        = otel.Tracer(name)
	meter         = otel.Meter(name)
	logger        = otelslog.NewLogger(name)
	rollCnt       metric.Int64Counter
	logs          []string
	logsMutex     sync.Mutex
	playerRoll    = make(map[string]int)
	leaderboard   = make(map[string]int)
	leaderboardMu sync.Mutex
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
	mode := r.URL.Query().Get("mode")

	if player == "" {
		player = "Anonymous"
	}

	var rollResult string
	switch mode {
	case "bestof3":
		rollResult = playBestOfThree(ctx, player)
	case "battle":
		opponent := r.URL.Query().Get("opponent")
		rollResult = playBattleMode(ctx, player, opponent)
	default:
		rollResult = playSingleRoll(ctx, player)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, rollResult)
}

func playSingleRoll(ctx context.Context, player string) string {
	roll := 1 + rand.Intn(6)
	event := ""
	if roll == 6 {
		event = fmt.Sprintf("%s hit the jackpot! 🎉", player)
	} else if roll == 1 {
		event = fmt.Sprintf("%s rolled a critical fail! 😢", player)
	}

	rollValueAttr := attribute.Int("roll.value", roll)
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("player", player), rollValueAttr)
	rollCnt.Add(ctx, 1, metric.WithAttributes(rollValueAttr))

	logger.InfoContext(ctx, fmt.Sprintf("%s rolled a %d 🎲 %s", player, roll, event))
	updateLeaderboard(player, roll)

	return fmt.Sprintf("Hi %s, your dice roll is: %d 🎲 %s\n", player, roll, event)
}

func playBestOfThree(ctx context.Context, player string) string {
	rolls := []int{1 + rand.Intn(6), 1 + rand.Intn(6), 1 + rand.Intn(6)}
	best := max(rolls)
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("player", player), attribute.IntSlice("rolls", rolls))
	rollCnt.Add(ctx, int64(len(rolls)))

	logger.InfoContext(ctx, fmt.Sprintf("%s played Best of Three: Rolls %v, Best %d 🎲", player, rolls, best))
	updateLeaderboard(player, best)

	return fmt.Sprintf("%s played Best of Three: Rolls %v, Best %d 🎲\n", player, rolls, best)
}

func playBattleMode(ctx context.Context, player, opponent string) string {
	if opponent == "" {
		opponent = "Anonymous Opponent"
	}
	playerRoll := 1 + rand.Intn(6)
	opponentRoll := 1 + rand.Intn(6)
	winner := determineWinner(player, opponent, playerRoll, opponentRoll)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("player", player),
		attribute.Int("player.roll", playerRoll),
		attribute.String("opponent", opponent),
		attribute.Int("opponent.roll", opponentRoll),
		attribute.String("winner", winner),
	)
	logger.InfoContext(ctx, fmt.Sprintf("%s vs %s: %s wins 🎲", player, opponent, winner))
	updateLeaderboard(player, playerRoll)
	updateLeaderboard(opponent, opponentRoll)

	return fmt.Sprintf("%s vs %s: %s wins 🎲\n", player, opponent, winner)
}

func determineWinner(player, opponent string, playerRoll, opponentRoll int) string {
	if playerRoll > opponentRoll {
		return player
	} else if opponentRoll > playerRoll {
		return opponent
	}
	return "It's a tie!"
}

func max(nums []int) int {
	best := nums[0]
	for _, num := range nums {
		if num > best {
			best = num
		}
	}
	return best
}

func updateLeaderboard(player string, roll int) {
	leaderboardMu.Lock()
	defer leaderboardMu.Unlock()
	if roll > leaderboard[player] {
		leaderboard[player] = roll
	}
}

func getLeaderboard() string {
	leaderboardMu.Lock()
	defer leaderboardMu.Unlock()

	result := "Leaderboard:\n"
	for player, score := range leaderboard {
		result += fmt.Sprintf("%s: %d\n", player, score)
	}
	return result
}

func viewLogs(w http.ResponseWriter, r *http.Request) {
	logsMutex.Lock()
	defer logsMutex.Unlock()
	for _, logEntry := range logs {
		io.WriteString(w, logEntry+"\n")
	}
}

func viewLeaderboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, getLeaderboard())
}

func resetMetrics(w http.ResponseWriter, r *http.Request) {
	playerRoll = make(map[string]int)
	leaderboardMu.Lock()
	leaderboard = make(map[string]int)
	leaderboardMu.Unlock()
	io.WriteString(w, "Metrics have been reset!\n")
}

func addLog(entry string) {
	logsMutex.Lock()
	defer logsMutex.Unlock()
	logs = append(logs, entry)
}