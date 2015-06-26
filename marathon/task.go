package marathon

type Task struct {
	Id        string `json:"id"`
	Ports     []int  `json:"ports"`
	Host      string `json:"host"`
	stagedAt  string `json:"stagedAt"`
	startedAt string `json:"startedAt"`
	version   string `json:"version"`
}
type TaskList []Task

type TaskStatus string

const (
	TaskStaging  TaskStatus = "TASK_STAGING"
	TaskStarting            = "TASK_STARTING"
	TaskRunning             = "TASK_RUNNING"
	TaskFinished            = "TASK_FINISHED"
	TaskFailed              = "TASK_FAILED"
	TaskKilled              = "TASK_KILLED"
	TaskLost                = "TASK_LOST"
)

func (a TaskList) Len() int           { return len(a) }
func (a TaskList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TaskList) Less(i, j int) bool { return a[i].Id < a[j].Id }
