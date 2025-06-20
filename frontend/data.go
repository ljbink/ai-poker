package frontend

import (
	"sync"
	"time"
)

// UserData represents player information
type UserData struct {
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	LastSeen    time.Time `json:"last_seen"`
	GamesPlayed int       `json:"games_played"`
	GamesWon    int       `json:"games_won"`
}

// SettingsData represents application settings
type SettingsData struct {
	Theme             string `json:"theme"` // "dark", "light", "auto"
	SoundEnabled      bool   `json:"sound_enabled"`
	AnimationsEnabled bool   `json:"animations_enabled"`
	AutoSave          bool   `json:"auto_save"`
	DefaultBuyIn      int    `json:"default_buy_in"`
	ShowProbabilities bool   `json:"show_probabilities"`

	// Game Setup Settings
	SmallBlind int `json:"small_blind"`
	BigBlind   int `json:"big_blind"`
	NumBots    int `json:"num_bots"`
}

// Data represents the central data store for the application
type Data struct {
	lock     sync.RWMutex
	user     *UserData
	settings *SettingsData
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
		ShowProbabilities: false,
		SmallBlind:        5,
		BigBlind:          10,
		NumBots:           3,
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
	case "show_probabilities":
		if v, ok := value.(bool); ok {
			d.settings.ShowProbabilities = v
		}
	case "small_blind":
		if v, ok := value.(int); ok {
			d.settings.SmallBlind = v
		}
	case "big_blind":
		if v, ok := value.(int); ok {
			d.settings.BigBlind = v
		}
	case "num_bots":
		if v, ok := value.(int); ok {
			d.settings.NumBots = v
		}
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
		ShowProbabilities: false,
		SmallBlind:        5,
		BigBlind:          10,
		NumBots:           3,
	}
}

// Game Setup Methods
func (d *Data) GetGameSetup() (smallBlind, bigBlind, numBots int) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	settings := d.GetSettings()
	return settings.SmallBlind, settings.BigBlind, settings.NumBots
}

func (d *Data) SetGameSetup(smallBlind, bigBlind, numBots int) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.settings == nil {
		d.settings = d.getDefaultSettings()
	}
	d.settings.SmallBlind = smallBlind
	d.settings.BigBlind = bigBlind
	d.settings.NumBots = numBots
}

// Utility Methods
func (d *Data) Reset() {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.user = nil
	d.settings = nil
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
