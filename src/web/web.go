package web

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"videomanager/database"
	"videomanager/telegram"
)

type metrics struct {
	tasks      *prometheus.CounterVec
	converters *prometheus.GaugeVec
}

type TaskProgress struct {
	Pass    int
	Speed   string
	Frame   string
	Bitrate string
	Size    string
	TimeMs  string
}

// Service stores all the Cdn servers synced with DB
type Service struct {
	mu              sync.RWMutex
	metrics         *metrics
	dbService       *database.Service
	telegramService *telegram.Service
	tasksProgress   map[int]*TaskProgress
	tasksProgressMu sync.RWMutex
}

func NewService(dbService *database.Service, promRegistry *prometheus.Registry, telegramService *telegram.Service) (s *Service, err error) {
	m := &metrics{
		tasks: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "india_vmanager_tasks",
			Help: "The number of tasks",
		}, []string{"status"}),
		converters: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "india_vmanager_converters",
			Help: "The network usage by cdn",
		}, []string{"status"}),
	}
	promRegistry.MustRegister(m.tasks, m.converters)

	s = &Service{
		metrics:         m,
		dbService:       dbService,
		telegramService: telegramService,
		tasksProgress:   make(map[int]*TaskProgress),
	}

	/*processCount := 0
	s.dbService.DB.Raw("select count(*) from flix_video_task where status='process' or status=''").Scan(&processCount)
	log.Println("tasks in queue:", processCount)*/
	return
}
