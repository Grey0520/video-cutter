package dao

import (
	"fmt"

	"video_cutter/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"video_cutter/configs"
)

var db *sqlx.DB

// Init 初始化MySQL连接
func Init(cfg *configs.MySQLConfig) (err error) {
	// "user:password@tcp(host:port)/dbname"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return
}

// Close 关闭MySQL连接
func Close() {
	_ = db.Close()
}

// RecordHistory inserts a new row into the history table
func RecordHistory(ClipReq *model.ClipRequest) (id int, err error) {
	sqlStr := "insert into history (video_url, start_time, end_time, file_name) values (?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, ClipReq.VideoUrl, ClipReq.StartTime, ClipReq.EndTime, ClipReq.FileName)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	// Select the last inserted id

	err = db.Get(&id, "SELECT LAST_INSERT_ID()")
	return
}

// GetClipById returns a ClipRequest struct from the history table
func GetClipById(id string) (ClipReq model.ClipRequest, err error) {
	sqlStr := "select * from history where id = ?"
	err = db.Get(&ClipReq, sqlStr, id)
	if err != nil {
		fmt.Printf("get ClipReq by id failed, err:%v\n", err)
		return
	}
	return
}

// AmDone updates the is_done field in the history table
func AmDone(url string) (err error) {
	sqlStr := "update history set is_done = 1 where file_name = ?"
	_, err = db.Exec(sqlStr, url)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	return
}
