package db

import (
	"errors"
	"fmt"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Player struct {
	ID    string  `gorm:"primarykey"` // Firebase UID
	Games []*Game `gorm:"many2many:player_games;"`

	// One-to-many relationships
	Submissions []Submission `gorm:"constraint:OnDelete:CASCADE;"`
	Rankings    []Ranking    `gorm:"constraint:OnDelete:CASCADE;"`
}

type Game struct {
	ID         string    `gorm:"primarykey"`
	Name       string    `gorm:"not null"`
	Players    []*Player `gorm:"many2many:player_games;"`
	Deadline   uint      `gorm:"not null"`
	NSongs     uint      `gorm:"not null"`
	Playlist   string    `gorm:"not null"`
	AddedSongs bool      `gorm:"not null"` // Were songs added to playlist yet or not

	// One-to-many relationships - each game has exactly two tierlists
	Tierlists []Tierlist `gorm:"constraint:OnDelete:CASCADE;"`
	// GuessList   Tierlist   `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE;"`
	// RankingList Tierlist   `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE;"`
	Rankings []Ranking `gorm:"constraint:OnDelete:CASCADE;"`

	// One-to-many relationship with submissions
	Submissions []Submission `gorm:"constraint:OnDelete:CASCADE;"`
}

type Tierlist struct {
	gorm.Model
	GameID string `gorm:"not null"`
	Type   string `gorm:"not null"` // "guess" or "ranking"

	// One-to-many relationships
	Tiers    []Tier    `gorm:"constraint:OnDelete:CASCADE;"`
	Rankings []Ranking `gorm:"constraint:OnDelete:CASCADE;"`
}

type Tier struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Rank         int    `gorm:"not null"` // Lower number = higher rank
	TierlistID   uint   `gorm:"not null"`
	SubmissionID *uint  // For guess list, associate player tier w/ their submission
}

type Submission struct {
	gorm.Model
	PlayerID string `gorm:"not null"`
	GameID   string `gorm:"not null"`
	Nickname string `gorm:"not null"`
	Songs    []Song `gorm:"constraint:OnDelete:CASCADE;"`
	Drawing  string `gorm:"not null"`
	Tier     Tier   `gorm:"constraint:OnDelete:CASCADE;"`

	// Unique constraint to ensure one submission per player per game
	UniqueSubmission string `gorm:"uniqueIndex:idx_player_game"`
	UniqueNickname   string `gorm:"uniqueIndex:idx_nickname_game"`
}

type Song struct {
	gorm.Model
	SubmissionID uint   `gorm:"not null"`
	Spotify      string `gorm:"not null"`
	AlbumArt     string `gorm:"not null"` // Album cover art URL
	Name         string `gorm:"not null"`
	GameID       string `gorm:"not null"` // Added to enforce unique songs per game

	// Unique constraint to prevent duplicate songs in a game
	UniqueSong string `gorm:"uniqueIndex:idx_game_song"`
}

type Ranking struct {
	gorm.Model
	PlayerID   string `gorm:"not null"`
	TierlistID uint   `gorm:"not null"`
	GameID     string `gorm:"not null"`
	Ranking    string `gorm:"not null"` // JSON tier: []songs

	// Unique constraint to ensure one ranking per player per tierlist
	UniqueRanking string `gorm:"uniqueIndex:idx_player_tierlist"`

	// Foreign key constraints
	Player   Player   `gorm:"constraint:OnDelete:CASCADE;"`
	Tierlist Tierlist `gorm:"constraint:OnDelete:CASCADE;"`
}

// Hooks to enforce business rules

func (g *Game) BeforeCreate(tx *gorm.DB) error {
	// If Deadline is before the current time, return an error
	if g.Deadline < uint(time.Now().Unix()) {
		return errors.New("deadline must be in the future")
	}

	// Name must be between 1 and 50 characters
	if len(g.Name) < 1 || len(g.Name) > 50 {
		return errors.New("name must be between 1 and 50 characters")
	}

	// Num songs must be between 1 and 5
	if g.NSongs < 1 || g.NSongs > 5 {
		return errors.New("number of songs must be between 1 and 5")
	}

	id, err := gonanoid.New()
	if err != nil {
		return err
	}

	g.ID = id
	g.AddedSongs = false
	g.Tierlists = append(g.Tierlists, Tierlist{Type: "guess"})
	ranking := Tierlist{Type: "ranking"}
	// Populate RankingList with default tiers (S, A, B, C, D)
	for i, tier := range []string{"S", "A", "B", "C", "D"} {
		ranking.Tiers = append(ranking.Tiers, Tier{Name: tier, Rank: i})
	}
	g.Tierlists = append(g.Tierlists, ranking)

	return nil
}

func (s *Submission) BeforeCreate(tx *gorm.DB) error {
	// Set the unique constraint value
	s.UniqueSubmission = fmt.Sprintf("%s-%s", s.PlayerID, s.GameID)
	s.UniqueNickname = fmt.Sprintf("%s-%s", s.Nickname, s.GameID)

	// Create tier in guesslist associated with submission
	var tierlist Tierlist
	err := tx.Where("type = ? and game_id = ?", "guess", s.GameID).First(&tierlist).Error
	if err != nil {
		return err
	}
	tier := Tier{
		Name:       s.Nickname,
		Rank:       0,
		TierlistID: tierlist.ID,
	}
	s.Tier = tier

	return nil
}

// func (s *Submission) AfterCreate(tx *gorm.DB) error {
// 	// Create tier in guesslist associated with submission
// 	var tierlist Tierlist
// 	err := tx.Where("type = ? and game_id = ?", "guess", s.GameID).First(&tierlist).Error
// 	if err != nil {
// 		return err
// 	}
// 	tier := Tier{
// 		Name:         s.Nickname,
// 		Rank:         0,
// 		SubmissionID: &s.ID,
// 		TierlistID:   tierlist.ID,
// 	}
// 	if err = tx.Create(tier).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

func (s *Song) BeforeCreate(tx *gorm.DB) error {
	// Get the GameID from the associated Submission
	var submission Submission
	if err := tx.First(&submission, s.SubmissionID).Error; err != nil {
		return err
	}

	// Set the GameID and unique constraint value
	s.GameID = submission.GameID
	s.UniqueSong = fmt.Sprintf("%s-%s", s.GameID, s.Spotify)

	// Check if the song already exists in this game
	var count int64
	if err := tx.Model(&Song{}).
		Where("game_id = ? AND spotify = ? AND id != ?", s.GameID, s.Spotify, s.ID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("this song has already been submitted to this game")
	}

	return nil
}

func (r *Ranking) BeforeCreate(tx *gorm.DB) error {
	// Set the unique constraint value
	r.UniqueRanking = fmt.Sprintf("%s-%d", r.PlayerID, r.TierlistID)

	// Check Tierlist is part of game
	var tierlist Tierlist
	if err := tx.First(&tierlist, r.TierlistID).Error; err != nil {
		return err
	}
	if tierlist.GameID != r.GameID {
		return errors.New("tierlist does not belong to game")
	}

	return nil
}
