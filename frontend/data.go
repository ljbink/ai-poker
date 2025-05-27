package frontend

import (
	"sync"
	"time"
)

// UserData represents player information
type UserData struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
	// Game stats could be added here later
	GamesPlayed int `json:"games_played"`
	GamesWon    int `json:"games_won"`
}

// SettingsData represents application settings
type SettingsData struct {
	Theme             string `json:"theme"` // "dark", "light", "auto"
	SoundEnabled      bool   `json:"sound_enabled"`
	AnimationsEnabled bool   `json:"animations_enabled"`
	AutoSave          bool   `json:"auto_save"`
	// Game-specific settings
	DefaultBuyIn      int  `json:"default_buy_in"`
	PreferredSeats    int  `json:"preferred_seats"`
	ShowProbabilities bool `json:"show_probabilities"`
}

// GameSessionData represents current game session state
type GameSessionData struct {
	SessionID    string    `json:"session_id"`
	StartTime    time.Time `json:"start_time"`
	IsActive     bool      `json:"is_active"`
	CurrentChips int       `json:"current_chips"`
	// Additional game state could be stored here
}

// Data represents the central data store for the application
type Data struct {
	lock        sync.RWMutex
	user        *UserData
	settings    *SettingsData
	gameSession *GameSessionData
}

// User Data Methods
func (d *Data) SetUser(user *UserData) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if user != nil {
		user.LastSeen = time.Now()
	}
	d.user = user
}

func (d *Data) GetUser() *UserData {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if d.user != nil {
		// Return a copy to prevent external modification
		userCopy := *d.user
		return &userCopy
	}
	return nil
}

func (d *Data) SetPlayerName(name string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.user == nil {
		d.user = &UserData{
			Name:      name,
			CreatedAt: time.Now(),
			LastSeen:  time.Now(),
		}
	} else {
		d.user.Name = name
		d.user.LastSeen = time.Now()
	}
}

func (d *Data) GetPlayerName() string {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if d.user != nil {
		return d.user.Name
	}
	return ""
}

func (d *Data) UpdateGameStats(won bool) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.user != nil {
		d.user.GamesPlayed++
		if won {
			d.user.GamesWon++
		}
		d.user.LastSeen = time.Now()
	}
}

// Settings Methods
func (d *Data) SetSettings(settings *SettingsData) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.settings = settings
}

func (d *Data) GetSettings() *SettingsData {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if d.settings != nil {
		// Return a copy to prevent external modification
		settingsCopy := *d.settings
		return &settingsCopy
	}
	// Return default settings if none set
	return &SettingsData{
		Theme:             "dark",
		SoundEnabled:      true,
		AnimationsEnabled: true,
		AutoSave:          true,
		DefaultBuyIn:      1000,
		PreferredSeats:    6,
		ShowProbabilities: false,
	}
}

func (d *Data) UpdateSetting(key string, value interface{}) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.settings == nil {
		d.settings = d.getDefaultSettings()
	}

	switch key {
	case "theme":
		if v, ok := value.(string); ok {
			d.settings.Theme = v
		}
	case "sound_enabled":
		if v, ok := value.(bool); ok {
			d.settings.SoundEnabled = v
		}
	case "animations_enabled":
		if v, ok := value.(bool); ok {
			d.settings.AnimationsEnabled = v
		}
	case "auto_save":
		if v, ok := value.(bool); ok {
			d.settings.AutoSave = v
		}
	case "default_buy_in":
		if v, ok := value.(int); ok {
			d.settings.DefaultBuyIn = v
		}
	case "preferred_seats":
		if v, ok := value.(int); ok {
			d.settings.PreferredSeats = v
		}
	case "show_probabilities":
		if v, ok := value.(bool); ok {
			d.settings.ShowProbabilities = v
		}
	}
}

// Game Session Methods
func (d *Data) StartGameSession(sessionID string, initialChips int) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.gameSession = &GameSessionData{
		SessionID:    sessionID,
		StartTime:    time.Now(),
		IsActive:     true,
		CurrentChips: initialChips,
	}
}

func (d *Data) EndGameSession() {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.gameSession != nil {
		d.gameSession.IsActive = false
	}
}

func (d *Data) GetGameSession() *GameSessionData {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if d.gameSession != nil {
		sessionCopy := *d.gameSession
		return &sessionCopy
	}
	return nil
}

func (d *Data) UpdateChips(chips int) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.gameSession != nil && d.gameSession.IsActive {
		d.gameSession.CurrentChips = chips
	}
}

// Helper method to get default settings
func (d *Data) getDefaultSettings() *SettingsData {
	return &SettingsData{
		Theme:             "dark",
		SoundEnabled:      true,
		AnimationsEnabled: true,
		AutoSave:          true,
		DefaultBuyIn:      1000,
		PreferredSeats:    6,
		ShowProbabilities: false,
	}
}

// Utility Methods
func (d *Data) Reset() {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.user = nil
	d.settings = nil
	d.gameSession = nil
}

func (d *Data) IsGameActive() bool {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.gameSession != nil && d.gameSession.IsActive
}

// Singleton pattern
var (
	dataInstance *Data
	once         sync.Once
)

func GetData() *Data {
	once.Do(func() {
		dataInstance = &Data{
			lock: sync.RWMutex{},
		}
	})
	return dataInstance
}
