package handlers 

import (
	"Lighthouse/internal/database"
	"github.com/google/uuid"
)

type Connection struct {
	ID uuid.UUID `json:"id"`
	Service database.Service `json:"service"`
	ChannelID string `json:"channelId"`
	ChannelHandle string `json:"channelHandle"`
	Scopes string `json:"scopes"`
	Active bool `json:"active"`
}