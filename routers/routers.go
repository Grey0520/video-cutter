package routers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	ffmpeg "github.com/u2takey/ffmpeg-go"

	"video_cutter/dao"
	"video_cutter/logger"
	"video_cutter/model"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	v1.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	v1.POST("/clip", ClipHandler)
	v1.GET("/clip", CheckHandler)

	v1.GET("/file", func(c *gin.Context) {
		c.File("./output/video-1672478378.mp4")
	})

	return r
}

// CheckHandler 用于检查后台处理的状态
func CheckHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is empty"})
		return
	}
	clip, err := dao.GetClipById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("当前状态为：%v", clip.IsDone),
	})
}

// ClipHandler 处理客户端的视频剪辑请求
func ClipHandler(c *gin.Context) {
	var ClipReq model.ClipRequest
	if err := c.ShouldBindJSON(&ClipReq); err != nil {
		fmt.Errorf("bind json failed, err:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileName, err := download(ClipReq.VideoUrl)
	if err != nil {
		fmt.Errorf("download failed, err:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ClipReq.FileName = fileName
	fmt.Println(ClipReq.FileName)
	id, err := dao.RecordHistory(&ClipReq)
	if err != nil {
		fmt.Errorf("record history failed, err:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	go clip(ClipReq.FileName, ClipReq.StartTime, ClipReq.StartTime)
	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("请求已成功提交，正在处理中，请稍后！id为%d", id),
	})
}

// clip 用于剪辑视频
func clip(url, start, end string) {
	err := ffmpeg.Input("./files/"+url, ffmpeg.KwArgs{"ss": start}).
		Output("./output/"+url, ffmpeg.KwArgs{"t": end}).OverWriteOutput().Run()
	if err != nil {
		panic(err)
	}
	dao.AmDone(url)
}

// download 用于下载视频
func download(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Errorf("req failed, err:%v", err)
		return "", err
	}
	req.Header.Set("User-Agent", "MyApp/1.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Errorf("resp failed, err:%v", err)
		return "", err
	}
	defer resp.Body.Close()

	t := time.Now().Unix()
	fileName := fmt.Sprintf("video-%d.mp4", t)
	file, err := os.Create("./files/" + fileName)
	if err != nil {
		fmt.Errorf("create file failed, err:%v", err)
		return "", err
	}
	defer file.Close()

	io.Copy(file, resp.Body)

	return fileName, nil
}
