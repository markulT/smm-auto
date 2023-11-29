package scheduler

import (
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

func TestMemoryLeak(t *testing.T) {

	// Your test logic here
	// For example, run the scheduler for some time
	schedulerTask := &SchedulerTask{}
	go schedulerTask.FetchAndProcessPosts()

	// Allow some time for the scheduler to run
	time.Sleep(5 * time.Second)

	f, err := os.Create("memprofile")
	if err != nil {
		t.Fatal(err)
	}
	pprof.WriteHeapProfile(f)
	f.Close()
}
