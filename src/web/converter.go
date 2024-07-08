package web

// HandlerTask for converters to get tasks
/*func (s *Service) HandlerTask(w http.ResponseWriter, r *http.Request) {
	// converter asks for task
	cpuQty, _ := strconv.Atoi(r.URL.Query().Get("cpu"))
	nodeId := r.URL.Query().Get("nodeId")
	version := r.URL.Query().Get("version")
	files := r.URL.Query().Get("files")
	_ = files

	if nodeId == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error in HandlerTask: empty nodeId")
		return
	}

	//  register converter and update LastSeen
	now := time.Now()
	converter := &database.Converter{
		NodeId:   nodeId,
		CpuQty:   cpuQty,
		LastSeen: &now,
		Attached: &now,
		Version:  version,
	}
	if err := converter.Register(s.dbService); err != nil {
		fmt.Println(err)
	}
	converter.Version = version
	if err := converter.UpdateLastSeen(s.dbService); err != nil {
		fmt.Println(err)
	}
	if err := converter.UpdateVersion(s.dbService); err != nil {
		fmt.Println(err)
	}

	task := &database.Task{
		ConverterId: converter.ID,
	}

	if converter.Active {
		// search task for this converter

		if err := task.Lock(s.dbService); err != nil {
			task = &database.Task{}
			fmt.Println(err)
		}

		if task.ID > 0 {
			task.StartTime = &now
			_ = task.Save(s.dbService)
			// get format
			format := database.Format{
				ID: task.FormatId,
			}
			if err := format.Load(s.dbService); err != nil {
				fmt.Println(err)
			}
			if task.Pass1 == "" {
				task.Pass1 = format.Pass1
			}
			if task.Pass2 == "" {
				task.Pass2 = format.Pass2
			}
			// update converter with task id
			if err := converter.UpdateTaskId(s.dbService, task.ID); err != nil {
				fmt.Println(err)
			}
			fmt.Printf("task: %d\n", task.ID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")

	taskJson, _ := json.Marshal(task)

	_, _ = w.Write(taskJson)

}*/
