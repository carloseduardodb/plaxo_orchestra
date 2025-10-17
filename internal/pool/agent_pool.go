package pool

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

type AgentInstance struct {
	ID       string
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
	// Use direct execution instead of persistent process
	instance := &AgentInstance{
		ID:       agentID,
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
	
	// Execute Q CLI directly with longer timeout for initialization
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, "q", "chat", input)
	
	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("Q CLI timeout after 120 seconds")
		}
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("Q CLI error (exit %d): %s", exitError.ExitCode(), string(exitError.Stderr))
		}
		return "", fmt.Errorf("Q CLI error: %v", err)
	}
	
	return string(output), nil
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
				delete(p.instances, id)
			}
		}
		p.mutex.Unlock()
	}
}


