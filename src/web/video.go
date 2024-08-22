package web

import (
	"encoding/json"
	"log"
	"net/http"
	"videomanager/database"
	"videomanager/helper"
)

type Job struct {
	ConverterId int
	Video       *database.Video
	Media       *database.Media
	Format      *database.Format
	Error       string
}

type ProgressFormat struct {
	Pass   int
	Size   int
	Speed  string
	TimeMs int
}

// VideoGetJob creates the video and returns it to converter
func (s *Service) VideoGetJob(w http.ResponseWriter, r *http.Request) {
	nodeName := r.URL.Query().Get("nodeName")
	cpuQty := helper.StrToInt(r.URL.Query().Get("cpuQty"))
	version := r.URL.Query().Get("version")

	converter := &database.Converter{Name: nodeName, CpuQty: cpuQty, Version: version}

	converter.Register(s.dbService)

	converter.Load(s.dbService)

	converter.UpdateVersion(s.dbService)

	if !converter.Active {
		log.Printf("blocked converter %s", converter.Name)

		job := Job{
			Error: "you are blocked",
		}
		res, _ := json.Marshal(job)
		w.Write(res)
		return
	}

	media, video, format, err := database.VideoFindJobForConverter(s.dbService, converter.ID)
	if err != nil {
		log.Println(err)
		job := Job{
			Error: err.Error(),
		}
		res, _ := json.Marshal(job)
		w.Write(res)
		return
	}

	job := Job{
		ConverterId: converter.ID,
		Video:       video,
		Media:       media,
		Format:      format,
	}
	res, _ := json.Marshal(job)

	if media != nil && video != nil {
		database.VideoLogAdd(s.dbService, 0, converter.ID, media.PostID, media.ID, video.ID, "task taken")
		log.Printf("task found for converter (%d)%s: postID: %d, mediaID:%d, videoID: %d", converter.ID, converter.Name, media.PostID, media.ID, video.ID)
		w.Write(res)
		return
	}

	log.Printf("no task found for converter %s :(", converter.Name)

	w.Write(res)

}

// VideoProgress for converters to report tasks progress
func (s *Service) VideoProgress(_ http.ResponseWriter, r *http.Request) {
	// log.Println("handler progress")
	pass := helper.StrToInt(r.URL.Query().Get("pass"))
	videoId := helper.StrToInt(r.URL.Query().Get("videoId"))
	converterId := helper.StrToInt(r.URL.Query().Get("converterId"))
	mediaId := helper.StrToInt(r.URL.Query().Get("mediaId"))
	postId := helper.StrToInt(r.URL.Query().Get("postId"))
	size := helper.StrToInt(r.URL.Query().Get("size"))
	speed := r.URL.Query().Get("speed")
	timeMs := helper.StrToInt(r.URL.Query().Get("timeMs"))

	/*taskProgress := &TaskProgress{
		Pass:    pass,
		Speed:   r.URL.Query().Get("speed"),
		Frame:   r.URL.Query().Get("frame"),
		Bitrate: r.URL.Query().Get("bitrate"),
		Size:    r.URL.Query().Get("size"),
		TimeMs:  r.URL.Query().Get("timeMs"),
	}*/

	report := &ProgressFormat{
		Pass:   pass,
		Size:   size,
		Speed:  speed,
		TimeMs: timeMs,
	}
	msg, _ := json.Marshal(report)
	//log.Printf("[VideoProgress] %s", msg)
	database.VideoLogAdd(s.dbService, 1, converterId, postId, mediaId, videoId, string(msg))

	// save progress to DB

	/*reqDump, _ := httputil.DumpRequest(r, true)
	// body, _ := io.ReadAll(r.Body)
	log.Printf("XXXXXXXX: %s", reqDump)
	log.Println("\n-----XXXXXXX")*/
}

// VideoDone for converters to report tasks
func (s *Service) VideoDone(w http.ResponseWriter, r *http.Request) {
	log.Println("handler task report")
	videoId := helper.StrToInt(r.URL.Query().Get("videoId"))
	totalSize := helper.StrToInt64(r.URL.Query().Get("totalSize"))
	lengthSeconds := helper.StrToInt(r.URL.Query().Get("lengthSeconds"))

	log.Printf("task report: videoId:%d totalSize: %d lengthSeconds:%d \n", videoId, totalSize, lengthSeconds)

	// update video
	video, err := database.VideoDone(s.dbService, videoId, totalSize)
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// update media
	err = database.MediaReadyToPlay(s.dbService, video.MediaId, lengthSeconds)
	if err != nil {
		log.Println(err)
		return
	}

	// TODO: update lang in post
	media, err := database.MediaGet(s.dbService, video.MediaId)
	if err != nil {
		log.Println(err)
		return
	}
	err = database.PostAddLang(s.dbService, media.PostID, media.LangId)
	if err != nil {
		log.Println(err)
		return
	}
	/*if task.Status == "error" {
		s.telegramService.Send(telegram.ChanVideo, fmt.Sprintf("Error converting video %s", task.SourcePath))
	}

	*/

	database.VideoLogAdd(s.dbService, 2, video.ConverterId, media.PostID, media.ID, video.ID, "task done")

	_, _ = w.Write([]byte("OK"))
}
