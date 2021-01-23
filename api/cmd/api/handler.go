package main

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/syncedvideo/syncedvideo"
)

// WsActionHandler manages WebSocket action handler methods
type WsActionHandler struct {
	WsAction *syncedvideo.WsAction
	Room     *syncedvideo.Room
	User     *syncedvideo.User
}

// NewWsActionHandler returns a new WsActionHandler
func NewWsActionHandler(a *syncedvideo.WsAction, r *syncedvideo.Room, u *syncedvideo.User) *WsActionHandler {
	return &WsActionHandler{
		WsAction: a,
		Room:     r,
		User:     u,
	}
}

// Handle incoming WebSocket actions
func (handler *WsActionHandler) Handle() {
	switch handler.WsAction.Name {

	// User actions
	case syncedvideo.WsActionUserSetBuffering:
		handler.handleUserSetBuffering()
	case syncedvideo.WsActionUserSetUsername:
		handler.handleUserSetUsername()
	case syncedvideo.WsActionUserSetColor:
		handler.handleUserSetColor()

	// Player actions
	case syncedvideo.WsActionPlayerPlay:
		handler.handlePlayerPlay()
	case syncedvideo.WsActionPlayerPause:
		handler.handlePlayerPause()
	case syncedvideo.WsActionPlayerSkip:
		handler.handlePlayerSkip()
	case syncedvideo.WsActionPlayerSeek:
		handler.handlePlayerSeek()

	// Queue actions
	case syncedvideo.WsActionQueueAdd:
		handler.handleQueueAdd()
	case syncedvideo.WsActionQueueRemove:
		handler.handleQueueRemove()
	case syncedvideo.WsActionQueueVote:
		handler.handleQueueVote()

	// Chat actions
	case syncedvideo.WsActionChatMessage:
		handler.handleChatMessage()
	}

	// Sync room state after handling the action
	handler.Room.BroadcastSync()
}

func (handler *WsActionHandler) handleUserSetBuffering() {
	var buffering bool
	err := json.Unmarshal(handler.WsAction.Data, &buffering)
	if err != nil {
		log.Println("handleUserSetBuffering error:", err)
		return
	}
	handler.User.SetBuffering(buffering)
}

func (handler *WsActionHandler) handleUserSetUsername() {
	var username string
	err := json.Unmarshal(handler.WsAction.Data, &username)
	if err != nil {
		log.Println("handleUserSetUsername error:", err)
		return
	}
	handler.User.SetUsername(username)
}

func (handler *WsActionHandler) handleUserSetColor() {
	var color string
	err := json.Unmarshal(handler.WsAction.Data, &color)
	if err != nil {
		log.Println("handleUserSetColor error:", err)
		return
	}
	handler.User.SetChatColor(color)
}

func (handler *WsActionHandler) handlePlayerPlay() {
	if handler.Room.Player.Video == nil {
		log.Println("handlePlayerPlay: Video is nil")
		return
	}
	handler.Room.Player.Play(handler.Room.Player.Video)
}

func (handler *WsActionHandler) handlePlayerPause() {
	if handler.Room.Player.Video == nil {
		log.Println("handlePlayerPause: Video is nil")
		return
	}
	handler.Room.Player.Playing = false
	log.Println("Player paused")
}

func (handler *WsActionHandler) handlePlayerSkip() {
	if len(handler.Room.Player.Queue.Videos) >= 1 {
		handler.Room.Player.Play(handler.Room.Player.Queue.Videos[0])
		handler.Room.Player.Queue.Remove(handler.Room.Player.Video.ID)
		log.Println("handlePlayerSkip: Video skipped by user:", handler.User)
		return
	}
	log.Println("handlePlayerSkip: Queue is empty")
}

func (handler *WsActionHandler) handlePlayerSeek() {
	var t int64
	err := json.Unmarshal(handler.WsAction.Data, &t)
	if err != nil {
		log.Println("e error:", err)
		return
	}
	handler.Room.Player.Time = t
	handler.Room.BroadcastRoomSeeked(t)
}

func (handler *WsActionHandler) handleQueueAdd() {
	var video *syncedvideo.Video
	err := json.Unmarshal(handler.WsAction.Data, &video)
	if err != nil {
		log.Println("handleQueueAdd error:", err)
		return
	}
	log.Println("handleQueueAdd:", video)
	if handler.Room.Player.Video == nil {
		handler.Room.Player.Play(video)
		return
	}
	video.AddVote(handler.User)
	handler.Room.Player.Queue.Add(handler.User, video)
}

func (handler *WsActionHandler) handleQueueRemove() {
	var idString string
	err := json.Unmarshal(handler.WsAction.Data, &idString)
	if err != nil {
		log.Println("handleQueueRemove error:", err)
		return
	}
	id, err := uuid.Parse(idString)
	if err != nil {
		log.Println("handleQueueRemove error:", err)
		return
	}
	handler.Room.Player.Queue.Remove(id)
}

func (handler *WsActionHandler) handleQueueVote() {
	var idString string
	err := json.Unmarshal(handler.WsAction.Data, &idString)
	if err != nil {
		log.Println("handleQueueVote error:", err)
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		log.Println("handleQueueVote error:", err)
		return
	}

	video := handler.Room.Player.Queue.Find(id)
	if video == nil {
		log.Println("handleQueueVote: Video %w not found", id)
		return
	}

	handler.Room.Player.Queue.ToggleVote(handler.User, video)
}

func (handler *WsActionHandler) handleChatMessage() {
	var text string
	err := json.Unmarshal(handler.WsAction.Data, &text)
	if err != nil {
		log.Println("ChatMessage error:", err)
		return
	}
	handler.Room.Chat.NewMessage(handler.User, text)
}
