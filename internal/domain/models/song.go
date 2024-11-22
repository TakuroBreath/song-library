package models

type Song struct {
	ID          int    `json:"id" `
	Group       string `json:"group" binding:"required,min=1,max=255"`
	Song        string `json:"song"  binding:"required,min=1,max=255"`
	ReleaseDate string `json:"release_date" binding:"required,datetime=2006-01-02"`
	Text        string `json:"text" binding:"required"`
	Link        string `json:"link" binding:"required,url"`
}
