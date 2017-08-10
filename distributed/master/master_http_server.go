package master

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/chrislusf/gleam/distributed/master/ui"
	"github.com/chrislusf/gleam/pb"
	"github.com/gorilla/mux"
	"github.com/hashicorp/golang-lru"
)

func (ms *MasterServer) uiStatusHandler(w http.ResponseWriter, r *http.Request) {
	infos := make(map[string]interface{})
	infos["Version"] = 0.01
	args := struct {
		Version   string
		Topology  interface{}
		StartTime time.Time
		Logs      *lru.Cache
	}{
		"0.01",
		ms.Topology,
		ms.startTime,
		ms.statusCache,
	}
	ui.MasterStatusTpl.Execute(w, args)
}

func (ms *MasterServer) jobStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobId, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		log.Printf("Failed to parse job id %s", vars["id"])
		return
	}
	status, ok := ms.statusCache.Get(uint32(jobId))
	if !ok {
		log.Printf("Failed to find job status for %d", jobId)
		return
	}

	infos := make(map[string]interface{})
	infos["Version"] = 0.01
	args := struct {
		Version   string
		Topology  interface{}
		Status    interface{}
		Svg       string
		StartTime time.Time
		Logs      *lru.Cache
	}{
		"0.01",
		ms.Topology,
		status,
		ui.GenSvg(status.(*pb.FlowExecutionStatus)),
		ms.startTime,
		ms.statusCache,
	}
	ui.JobStatusTpl.Execute(w, args)
}
