package orchestrator

type Spec struct {
	workingDir string
}

type Watcher struct {
	workingDir string
	orch       *Orchestrator
}

func NewSpec(workingDir string) *Spec {
	return &Spec{
		workingDir: workingDir,
	}
}

func NewWatcher(workingDir string) *Watcher {
	return &Watcher{
		workingDir: workingDir,
		orch:       New(workingDir),
	}
}

func (s *Spec) Generate() error {
	// Implementação básica para gerar especificação
	return nil
}

func (w *Watcher) Start() error {
	return w.orch.StartWatchMode()
}
