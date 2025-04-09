package main

// Series define la estructura de una serie de TV/Streaming.
type Series struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Status         string `json:"status"`
	CurrentEpisode int    `json:"current_episode"`
	Score          int    `json:"score"`
}