package services

import (
	"crypto/rand"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/takumi/fastchem/internal/models"
)

// Room represents a custom (friendly) match room.
type Room struct {
	Code      string
	HostID    int64
	HostName  string
	GuestID   int64
	GuestName string
	MatchID   int64 // set once match starts
	CreatedAt time.Time
	Status    string // "waiting", "active", "finished"
}

// RoomService manages custom match rooms with 6-digit codes.
type RoomService struct {
	mu           sync.Mutex
	rooms        map[string]*Room    // code → room
	matchService *RankedMatchService // reuse the same match engine (without DB persist)
}

// NewRoomService creates a new room service.
func NewRoomService(matchService *RankedMatchService) *RoomService {
	return &RoomService{
		rooms:        make(map[string]*Room),
		matchService: matchService,
	}
}

// generateCode generates a random 6-digit numeric code.
func generateCode() string {
	b := make([]byte, 3)
	rand.Read(b)
	code := fmt.Sprintf("%06d", (int(b[0])<<16|int(b[1])<<8|int(b[2]))%1000000)
	return code
}

// CreateRoom creates a new room and returns its code.
func (rs *RoomService) CreateRoom(hostID int64, hostName string) string {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	// Generate unique code
	var code string
	for {
		code = generateCode()
		if _, exists := rs.rooms[code]; !exists {
			break
		}
	}

	room := &Room{
		Code:      code,
		HostID:    hostID,
		HostName:  hostName,
		CreatedAt: time.Now(),
		Status:    "waiting",
	}
	rs.rooms[code] = room
	slog.Info("room: created room", "code", code, "host", hostName, "hostID", hostID)
	return code
}

// JoinRoom lets a guest join a room. Returns the room if successful.
func (rs *RoomService) JoinRoom(code string, guestID int64, guestName string) (*Room, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	room, exists := rs.rooms[code]
	if !exists {
		return nil, fmt.Errorf("ไม่พบห้อง %s", code)
	}
	if room.Status != "waiting" {
		return nil, fmt.Errorf("ห้องนี้เริ่มเกมแล้ว")
	}
	if room.HostID == guestID {
		return nil, fmt.Errorf("ไม่สามารถเข้าห้องตัวเองได้")
	}
	if room.GuestID != 0 {
		return nil, fmt.Errorf("ห้องเต็มแล้ว")
	}

	room.GuestID = guestID
	room.GuestName = guestName
	return room, nil
}

// GetRoom returns a room by code.
func (rs *RoomService) GetRoom(code string) (*Room, bool) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	r, ok := rs.rooms[code]
	return r, ok
}

// StartMatch creates an in-memory match for the room (no DB persist).
func (rs *RoomService) StartMatch(code string) (*ActiveRankedMatch, error) {
	rs.mu.Lock()
	room, exists := rs.rooms[code]
	if !exists {
		rs.mu.Unlock()
		return nil, fmt.Errorf("room not found")
	}
	if room.Status != "waiting" || room.GuestID == 0 {
		rs.mu.Unlock()
		return nil, fmt.Errorf("room not ready")
	}
	room.Status = "active"
	rs.mu.Unlock()

	seed := time.Now().UnixNano()
	questions := GenerateRankedQuestions(seed)

	// Use a negative matchID (no DB entry) to distinguish from ranked matches
	matchID := -time.Now().UnixNano()

	match := &ActiveRankedMatch{
		MatchID: matchID,
		Seed:    seed,
		Player1: &RankedPlayerState{
			UserID:    room.HostID,
			Username:  room.HostName,
			Rating:    models.DefaultRating,
			Connected: true,
			Results:   make([]models.RankedQuestionResult, 0, models.RankedQuestionsPerMatch),
		},
		Player2: &RankedPlayerState{
			UserID:    room.GuestID,
			Username:  room.GuestName,
			Rating:    models.DefaultRating,
			Connected: true,
			Results:   make([]models.RankedQuestionResult, 0, models.RankedQuestionsPerMatch),
		},
		Questions:       questions,
		Status:          "active",
		CurrentP1Index:  0,
		CurrentP2Index:  0,
		LastSyncedIndex: 0,
		CreatedAt:       time.Now(),
		P1Send:          make(chan models.WSMessage, 32),
		P2Send:          make(chan models.WSMessage, 32),
		Done:            make(chan struct{}),
	}

	rs.matchService.mu.Lock()
	rs.matchService.matches[matchID] = match
	rs.matchService.byUser[room.HostID] = matchID
	rs.matchService.byUser[room.GuestID] = matchID
	rs.matchService.mu.Unlock()

	room.MatchID = matchID
	return match, nil
}

// RemoveRoom removes a room from the service.
func (rs *RoomService) RemoveRoom(code string) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	delete(rs.rooms, code)
}

// GetRoomByUserID finds a room where the user is host or guest.
func (rs *RoomService) GetRoomByUserID(userID int64) (*Room, bool) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	for _, r := range rs.rooms {
		if r.HostID == userID || r.GuestID == userID {
			return r, true
		}
	}
	return nil, false
}

// CleanupStaleRooms removes rooms older than maxAge that are still waiting.
func (rs *RoomService) CleanupStaleRooms(maxAge time.Duration) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	cutoff := time.Now().Add(-maxAge)
	for code, r := range rs.rooms {
		if r.CreatedAt.Before(cutoff) && r.Status == "waiting" {
			slog.Info("room: cleaning up stale room", "code", code)
			delete(rs.rooms, code)
		}
	}
}
