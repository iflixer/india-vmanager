package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"videomanager/database"
	"videomanager/helper"
)

type Job struct {
	Video  *database.Video
	Media  *database.Media
	Format *database.Format
	Error  string
}

// VideoGetJob creates the video and returns it to converter
func (s *Service) VideoGetJob(w http.ResponseWriter, r *http.Request) {
	nodeName := r.URL.Query().Get("nodeName")
	cpuQty := helper.StrToInt(r.URL.Query().Get("cpuQty"))
	version := r.URL.Query().Get("version")

	converter := database.Converter{Name: nodeName, CpuQty: cpuQty, Version: version}

	converter.Register(s.dbService)

	converter.Load(s.dbService)

	if !converter.Active {
		w.Write([]byte("you are blocked"))
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
		Video:  video,
		Media:  media,
		Format: format,
	}
	res, _ := json.Marshal(job)
	w.Write(res)

}

// VideoProgress for converters to report tasks progreee
func (s *Service) VideoProgress(_ http.ResponseWriter, r *http.Request) {
	// log.Println("handler progress")
	params := mux.Vars(r)
	pass := helper.StrToInt(params["pass"])
	videoId := helper.StrToInt(params["videoId"])
	size := helper.StrToInt(params["size"])
	speed := params["speed"]
	timeMs := helper.StrToInt(params["timeMs"])

	/*taskProgress := &TaskProgress{
		Pass:    pass,
		Speed:   r.URL.Query().Get("speed"),
		Frame:   r.URL.Query().Get("frame"),
		Bitrate: r.URL.Query().Get("bitrate"),
		Size:    r.URL.Query().Get("size"),
		TimeMs:  r.URL.Query().Get("timeMs"),
	}*/

	log.Printf("progress video %d: pass %d size: %d speed: %s timeMs: %d \n", videoId, pass, size, speed, timeMs)

	// save progress to DB

	/*reqDump, _ := httputil.DumpRequest(r, true)
	// body, _ := io.ReadAll(r.Body)
	log.Printf("XXXXXXXX: %s", reqDump)
	log.Println("\n-----XXXXXXX")*/
}

// VideoDone for converters to report tasks
func (s *Service) VideoDone(w http.ResponseWriter, r *http.Request) {
	log.Println("handler task report")
	params := mux.Vars(r)
	videoId := helper.StrToInt(params["videoId"])
	converterId := helper.StrToInt(params["converterId"])
	size := helper.StrToInt(params["size"])

	log.Printf("task report: videoId:%d converterId: %d size: %d \n", videoId, converterId, size)

	// update video
	video, err := database.VideoDone(s.dbService, videoId, size)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// update media
	err = database.MediaReadyToPlay(s.dbService, video.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// TODO: update lang in post

	/*if task.Status == "error" {
		s.telegramService.Send(telegram.ChanVideo, fmt.Sprintf("Error converting video %s", task.SourcePath))
	}

	*/
	_, _ = w.Write([]byte("OK"))
}
