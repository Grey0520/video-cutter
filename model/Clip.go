package model

type ClipRequest struct {
	ID        int    `db:"id" json:"id"`
	VideoUrl  string `db:"video_url" json:"video_url"`
	StartTime string `db:"start_time" json:"start_time"`
	EndTime   string `db:"end_time" json:"end_time"`
	FileName  string `db:"file_name" json:"file_name"`
	IsDone    bool   `db:"is_done" json:"is_done"`
}
