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
	var err error

	converter := &database.Converter{}

	err = converter.Register(s.dbService, nodeName)
	if err != nil {
		log.Println(err)
	}

	converter.CpuQty = cpuQty
	converter.Version = version

	err = converter.Save(s.dbService)
	if err != nil {
		log.Println(err)
	}

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
		converter.PostID = media.PostID
		converter.MediaID = media.ID
		converter.VideoID = video.ID
		converter.Save(s.dbService)

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
	//pass := helper.StrToInt(r.URL.Query().Get("pass"))
	videoId := helper.StrToInt(r.URL.Query().Get("videoId"))
	//converterId := helper.StrToInt(r.URL.Query().Get("converterId"))
	//mediaId := helper.StrToInt(r.URL.Query().Get("mediaId"))
	//postId := helper.StrToInt(r.URL.Query().Get("postId"))
	//size := helper.StrToInt(r.URL.Query().Get("size"))
	//speed := r.URL.Query().Get("speed")
	seconds := helper.StrToInt(r.URL.Query().Get("seconds"))

	/*taskProgress := &TaskProgress{
		Pass:    pass,
		Speed:   r.URL.Query().Get("speed"),
		Frame:   r.URL.Query().Get("frame"),
		Bitrate: r.URL.Query().Get("bitrate"),
		Size:    r.URL.Query().Get("size"),
		TimeMs:  r.URL.Query().Get("timeMs"),
	}*/

	// report := &ProgressFormat{
	// 	Pass:   pass,
	// 	Size:   size,
	// 	Speed:  speed,
	// 	TimeMs: timeMs,
	// }
	//msg, _ := json.Marshal(report)

	// update video
	err := database.VideoUpdateProgress(s.dbService, videoId, seconds)
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	/*reqDump, _ := httputil.DumpRequest(r, true)
	// body, _ := io.ReadAll(r.Body)
	log.Printf("XXXXXXXX: %s", reqDump)
	log.Println("\n-----XXXXXXX")*/
}

// VideoUpdate for converters to report tasks
func (s *Service) VideoUpdate(w http.ResponseWriter, r *http.Request) {
	//log.Println("VideoUpdate")
	videoId := helper.StrToInt(r.URL.Query().Get("videoId"))
	progress := helper.StrToInt64(r.URL.Query().Get("progress"))
	status := helper.StrToInt(r.URL.Query().Get("status"))
	msg := r.URL.Query().Get("msg")

	log.Printf("[VideoUpdate] videoId:%d status:%d \n", status, videoId)

	video := &database.Video{}
	video.Load(s.dbService, videoId)

	media := &database.Media{}
	media.Load(s.dbService, video.MediaId)

	//if status == -1 { // error
	//TODO: send message to telegram
	//}

	if status == 2 { // probe
		// dirty hack - we use progress to send duration
		media.Duration = int(progress)
		video.Duration = int(progress)
		progress = 0
		media.Save(s.dbService)
	}

	video.Status = status
	video.Msg = msg
	video.Progress = int(progress)
	video.Save(s.dbService)

	if status == 5 {
		//totalSize := helper.StrToInt64(r.URL.Query().Get("totalSize"))
		//lengthSeconds := helper.StrToInt(r.URL.Query().Get("lengthSeconds"))

		video.FileSize = progress
		video.Progress = 0
		video.Save(s.dbService)

		media.Status = 2
		media.Save(s.dbService)

		// TODO: update lang in post

		err := database.PostAddLang(s.dbService, media.PostID, media.LangId)
		if err != nil {
			log.Println(err)
			return
		}
		/*if task.Status == "error" {
			s.telegramService.Send(telegram.ChanVideo, fmt.Sprintf("Error converting video %s", task.SourcePath))
		}

		*/

	}

	_, _ = w.Write([]byte("OK"))
}
