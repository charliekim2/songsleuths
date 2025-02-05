package db

import (
	"errors"
	"fmt"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Player struct {
	gorm.Model
	Name  string  `gorm:"not null"`
	Games []*Game `gorm:"many2many:player_games;"`

	// One-to-many relationships
	Submissions []Submission `gorm:"constraint:OnDelete:CASCADE;"`
	Rankings    []Ranking    `gorm:"constraint:OnDelete:CASCADE;"`
}

type Game struct {
	ID       string    `gorm:"primarykey"`
	Name     string    `gorm:"not null"`
	Players  []*Player `gorm:"many2many:player_games;"`
	Deadline uint      `gorm:"not null"`
	NSongs   uint      `gorm:"not null"`

	// One-to-many relationships - each game has exactly two tierlists
	Guesslist   Tierlist `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE;"`
	RankingList Tierlist `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE;"`

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
	Name       string `gorm:"not null"`
	Rank       int    `gorm:"not null"` // Lower number = higher rank
	TierlistID uint   `gorm:"not null"`

	// One-to-many relationship
	Rankings []Ranking `gorm:"constraint:OnDelete:CASCADE;"`
}

type Submission struct {
	gorm.Model
	PlayerID uint   `gorm:"not null"`
	GameID   string `gorm:"not null"`
	Nickname string `gorm:"not null"`
	Songs    []Song `gorm:"constraint:OnDelete:CASCADE;"`

	// Unique constraint to ensure one submission per player per game
	UniqueSubmission string `gorm:"uniqueIndex:idx_player_game"`
}

type Song struct {
	gorm.Model
	SubmissionID uint   `gorm:"not null"`
	YoutubeLink  string `gorm:"not null"`
	GameID       string `gorm:"not null"` // Added to enforce unique songs per game

	// Unique constraint to prevent duplicate songs in a game
	UniqueSong string `gorm:"uniqueIndex:idx_game_song"`

	// One-to-many relationship with rankings
	Rankings []Ranking `gorm:"constraint:OnDelete:CASCADE;"`
}

type Ranking struct {
	gorm.Model
	PlayerID   uint `gorm:"not null"`
	TierlistID uint `gorm:"not null"`
	TierID     uint `gorm:"not null"`
	SongID     uint `gorm:"not null"`

	// Unique constraint to ensure one ranking per player per song per tierlist
	UniqueRanking string `gorm:"uniqueIndex:idx_player_tierlist_song"`

	// Foreign key constraints
	Player   Player   `gorm:"constraint:OnDelete:CASCADE;"`
	Tierlist Tierlist `gorm:"constraint:OnDelete:CASCADE;"`
	Tier     Tier     `gorm:"constraint:OnDelete:CASCADE;"`
	Song     Song     `gorm:"constraint:OnDelete:CASCADE;"`
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
	g.Guesslist = Tierlist{Type: "guess"}
	g.RankingList = Tierlist{Type: "ranking"}
	// Populate RankingList with default tiers (S, A, B, C, D)
	for i, tier := range []string{"S", "A", "B", "C", "D"} {
		g.RankingList.Tiers = append(g.RankingList.Tiers, Tier{Name: tier, Rank: i})
	}

	return nil
}

func (s *Submission) BeforeCreate(tx *gorm.DB) error {
	// Set the unique constraint value
	s.UniqueSubmission = fmt.Sprintf("%d-%s", s.PlayerID, s.GameID)
	return nil
}

func (s *Song) BeforeCreate(tx *gorm.DB) error {
	// Get the GameID from the associated Submission
	var submission Submission
	if err := tx.First(&submission, s.SubmissionID).Error; err != nil {
		return err
	}

	// Set the GameID and unique constraint value
	s.GameID = submission.GameID
	s.UniqueSong = fmt.Sprintf("%s-%s", s.GameID, s.YoutubeLink)

	// Check if the song already exists in this game
	var count int64
	if err := tx.Model(&Song{}).
		Where("game_id = ? AND youtube_link = ? AND id != ?", s.GameID, s.YoutubeLink, s.ID).
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
	r.UniqueRanking = fmt.Sprintf("%d-%d-%d", r.PlayerID, r.TierlistID, r.SongID)

	// Verify that the Tier belongs to the Tierlist
	var count int64
	if err := tx.Model(&Tier{}).
		Where("id = ? AND tierlist_id = ?", r.TierID, r.TierlistID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("tier must belong to the specified tierlist")
	}

	// Verify that the Player is part of the Game associated with the Tierlist
	var tierlist Tierlist
	if err := tx.First(&tierlist, r.TierlistID).Error; err != nil {
		return err
	}

	var playerCount int64
	if err := tx.Table("player_games").
		Where("player_id = ? AND game_id = ?", r.PlayerID, tierlist.GameID).
		Count(&playerCount).Error; err != nil {
		return err
	}
	if playerCount == 0 {
		return errors.New("player must be part of the game to create rankings")
	}

	return nil
}
