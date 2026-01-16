package controller

// Message types.
type estimationMsg struct {
	total     int
	paths     int
	fileStats map[string]int
	err       error
}

type upcomingMsg struct {
	count int
}

type startMutationMsg struct {
	id     int
	thread int
	kind   interface{}
	path   string
}

type completedMutationMsg struct {
	id     int
	kind   interface{}
	status string
}

type concurrencyMsg struct {
	threads    int
	shardIndex int
	shards     int
}

// List item types.
type fileItem struct {
	path  string
	count int
}

func (f fileItem) FilterValue() string {
	return f.path
}
