package pool

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

type AgentInstance struct {
	ID       string
	Cmd      *exec.Cmd
	Stdin    io.WriteCloser
	Stdout   io.ReadCloser
	Scanner  *bufio.Scanner
	LastUsed time.Time
	InUse    bool
	Context  string
}

type AgentPool struct {
	instances map[string]*AgentInstance
	mutex     sync.RWMutex
	maxIdle   time.Duration
}

func NewAgentPool() *AgentPool {
	pool := &AgentPool{
		instances: make(map[string]*AgentInstance),
		maxIdle:   10 * time.Minute,
	}
	
	// Cleanup routine
	go pool.cleanup()
	
	return pool
}

func (p *AgentPool) GetOrCreate(agentID string) (*AgentInstance, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	if instance, exists := p.instances[agentID]; exists {
		if !instance.InUse {
			instance.InUse = true
			instance.LastUsed = time.Now()
			return instance, nil
		}
	}
	
	// Create new instance
	instance, err := p.createInstance(agentID)
	if err != nil {
		return nil, err
	}
	
	p.instances[agentID] = instance
	return instance, nil
}

func (p *AgentPool) createInstance(agentID string) (*AgentInstance, error) {
	cmd := exec.Command("q", "chat")
	
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	
	scanner := bufio.NewScanner(stdout)
	
	instance := &AgentInstance{
		ID:       agentID,
		Cmd:      cmd,
		Stdin:    stdin,
		Stdout:   stdout,
		Scanner:  scanner,
		LastUsed: time.Now(),
		InUse:    true,
	}
	
	return instance, nil
}

func (p *AgentPool) Execute(agentID, input string) (string, error) {
	instance, err := p.GetOrCreate(agentID)
	if err != nil {
		return "", err
	}
	defer p.Release(instance)
	
	// Send input
	if _, err := instance.Stdin.Write([]byte(input + "\n")); err != nil {
		return "", err
	}
	
	// Read response
	var response string
	timeout := time.After(30 * time.Second)
	done := make(chan bool)
	
	go func() {
		for instance.Scanner.Scan() {
			line := instance.Scanner.Text()
			response += line + "\n"
			if isResponseComplete(line) {
				done <- true
				return
			}
		}
	}()
	
	select {
	case <-done:
		return response, nil
	case <-timeout:
		return "", fmt.Errorf("timeout waiting for response")
	}
}

func (p *AgentPool) Release(instance *AgentInstance) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	instance.InUse = false
	instance.LastUsed = time.Now()
}

func (p *AgentPool) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		p.mutex.Lock()
		for id, instance := range p.instances {
			if !instance.InUse && time.Since(instance.LastUsed) > p.maxIdle {
				instance.Cmd.Process.Kill()
				delete(p.instances, id)
			}
		}
		p.mutex.Unlock()
	}
}

func isResponseComplete(line string) bool {
	// Simple heuristic - adjust based on Q CLI output patterns
	return line == "" || len(line) == 0
}
