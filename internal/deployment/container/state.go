package container

import "github.com/docker/docker/api/types"

type State struct {
	Status     string        `json:"status"` // created
	Running    bool          `json:"running"`
	Paused     bool          `json:"paused"`
	Restarting bool          `json:"restarting"`
	OOMKilled  bool          `json:"oom_killed"`
	Dead       bool          `json:"dead"`
	Pid        int           `json:"pid"`
	ExitCode   int           `json:"exit_code"`
	Error      string        `json:"error"`
	StartedAt  string        `json:"started"`
	FinishedAt string        `json:"finished_at"`
	Health     *types.Health `json:",omitempty"`
}

type ContainerState string

const (
	ContainerRunning ContainerState = "running"
	ContainerStarted ContainerState = "started"
	ContainerCreated ContainerState = "created"
)

// fromDockerStateToKstate : map a Docker state to a Krane container State
func fromDockerStateToKstate(state types.ContainerState) State {
	return State{
		Status:     state.Status,
		Running:    state.Running,
		Paused:     state.Paused,
		Restarting: state.Restarting,
		OOMKilled:  state.OOMKilled,
		Dead:       state.Dead,
		Pid:        state.Pid,
		ExitCode:   state.ExitCode,
		Error:      state.Error,
		StartedAt:  state.StartedAt,
		FinishedAt: state.FinishedAt,
		Health:     state.Health,
	}
}
