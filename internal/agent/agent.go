package agent

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Agent struct {
	Name         string
	Domain       string
	Instructions string
	Memory       []string
	WorkingDir   string
}

func NewAgent(domain, workingDir string) *Agent {
	// Suporte para bounded contexts (ex: "user/profile")
	parts := strings.Split(domain, "/")
	name := strings.Join(parts, "-") + "-agent"
	
	return &Agent{
		Name:       name,
		Domain:     domain,
		WorkingDir: workingDir,
		Memory:     make([]string, 0),
	}
}

func (a *Agent) LoadInstructions() error {
	// Suporte para bounded contexts
	parts := strings.Split(a.Domain, "/")
	var instructionsPath string
	
	if len(parts) == 1 {
		// Domínio simples: user/agents/instructions.txt
		instructionsPath = filepath.Join(a.WorkingDir, parts[0], "agents", "instructions.txt")
	} else {
		// Bounded context: user/profile/agents/instructions.txt
		instructionsPath = filepath.Join(a.WorkingDir, parts[0], parts[1], "agents", "instructions.txt")
	}
	
	if _, err := os.Stat(instructionsPath); os.IsNotExist(err) {
		return fmt.Errorf("instructions not found for domain %s", a.Domain)
	}
	
	content, err := os.ReadFile(instructionsPath)
	if err != nil {
		return err
	}
	
	a.Instructions = string(content)
	return nil
}

func (a *Agent) SaveMemory(entry string) error {
	a.Memory = append(a.Memory, entry)
	
	// Suporte para bounded contexts
	parts := strings.Split(a.Domain, "/")
	var memoryPath string
	
	if len(parts) == 1 {
		memoryPath = filepath.Join(a.WorkingDir, parts[0], "agents", "memory.txt")
	} else {
		memoryPath = filepath.Join(a.WorkingDir, parts[0], parts[1], "agents", "memory.txt")
	}
	
	os.MkdirAll(filepath.Dir(memoryPath), 0755)
	
	f, err := os.OpenFile(memoryPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	
	_, err = f.WriteString(entry + "\n")
	return err
}

func (a *Agent) Execute(task string) (string, error) {
	context := fmt.Sprintf(`
Domain: %s
Instructions: %s
Recent Memory: %s
Task: %s
`, a.Domain, a.Instructions, strings.Join(a.Memory[max(0, len(a.Memory)-5):], "\n"), task)

	cmd := exec.Command("q", "chat", context)
	output, err := cmd.Output()
	
	if err == nil {
		a.SaveMemory(fmt.Sprintf("Task: %s | Result: %s", task, string(output)))
	}
	
	return string(output), err
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
