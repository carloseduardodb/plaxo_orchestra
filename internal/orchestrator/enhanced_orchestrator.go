package orchestrator

import (
	"context"
	"fmt"
	"plaxo-orchestra/internal/cache"
	"plaxo-orchestra/internal/learning"
	"plaxo-orchestra/internal/observability"
	"plaxo-orchestra/internal/pool"
	"sync"
	"time"
)

type EnhancedOrchestrator struct {
	*Orchestrator
	cache       *cache.DistributedCache
	learning    *learning.AdvancedLearning
	observer    *observability.Observer
	processor   *pool.AsyncProcessor
	circuitBreaker *CircuitBreaker
}

type CircuitBreaker struct {
	failures    int
	lastFailure time.Time
	state       CircuitState
	threshold   int
	timeout     time.Duration
	mutex       sync.RWMutex
}

type CircuitState int

const (
	Closed CircuitState = iota
	Open
	HalfOpen
)

type WorkflowStep struct {
	Agent        string
	Dependencies []string
	Parallel     bool
	Context      map[string]interface{}
}

func NewEnhancedOrchestrator(workingDir string) *EnhancedOrchestrator {
	connectionPool := pool.NewConnectionPool(10, 60*time.Second)
	
	return &EnhancedOrchestrator{
		Orchestrator:   New(workingDir),
		cache:          cache.NewDistributedCache(),
		learning:       learning.NewAdvancedLearning(),
		observer:       observability.NewObserver(),
		processor:      pool.NewAsyncProcessor(connectionPool, 5),
		circuitBreaker: NewCircuitBreaker(5, 1*time.Minute),
	}
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		threshold: threshold,
		timeout:   timeout,
		state:     Closed,
	}
}

func (eo *EnhancedOrchestrator) ProcessWithIntelligence(ctx context.Context, input string) error {
	span := eo.observer.StartSpan("process_request", map[string]string{
		"input_length": fmt.Sprintf("%d", len(input)),
	})
	defer func() {
		eo.observer.FinishSpan(span, true, nil)
	}()
	
	// Check cache first
	cacheKey := eo.cache.GenerateKey(input)
	if cached, found := eo.cache.Get(ctx, cacheKey); found {
		fmt.Println("🚀 Resposta do cache")
		fmt.Println(cached)
		return nil
	}
	
	// Get proactive suggestions
	suggestions := eo.learning.GetProactiveSuggestions(input)
	if len(suggestions) > 0 {
		fmt.Println("💡 Sugestões baseadas em aprendizado:")
		for _, suggestion := range suggestions {
			fmt.Printf("  • %s\n", suggestion)
		}
	}
	
	// Check circuit breaker
	if !eo.circuitBreaker.CanExecute() {
		return fmt.Errorf("circuit breaker is open, service temporarily unavailable")
	}
	
	// Plan intelligent workflow
	workflow, err := eo.planIntelligentWorkflow(ctx, input)
	if err != nil {
		eo.circuitBreaker.RecordFailure()
		return err
	}
	
	// Execute workflow with parallelization
	result, err := eo.executeWorkflow(ctx, workflow, input)
	if err != nil {
		eo.circuitBreaker.RecordFailure()
		eo.learning.RecordFeedback(span.SpanID, input, "workflow", false, 1)
		return err
	}
	
	// Cache successful result
	eo.cache.Set(ctx, cacheKey, result, 10*time.Minute)
	eo.circuitBreaker.RecordSuccess()
	eo.learning.RecordFeedback(span.SpanID, input, "workflow", true, 5)
	
	fmt.Println(result)
	return nil
}

func (eo *EnhancedOrchestrator) planIntelligentWorkflow(ctx context.Context, input string) ([]WorkflowStep, error) {
	span := eo.observer.StartSpan("plan_workflow", map[string]string{
		"input": input,
	})
	defer eo.observer.FinishSpan(span, true, nil)
	
	// Analyze dependencies and create execution plan
	workflow := []WorkflowStep{}
	
	// Simple workflow planning (would be more sophisticated in production)
	if eo.needsCoordination(input) {
		// Multi-agent coordination workflow
		workflow = append(workflow, WorkflowStep{
			Agent:        "analyzer",
			Dependencies: []string{},
			Parallel:     false,
			Context:      map[string]interface{}{"phase": "analysis"},
		})
		
		workflow = append(workflow, WorkflowStep{
			Agent:        "coordinator",
			Dependencies: []string{"analyzer"},
			Parallel:     false,
			Context:      map[string]interface{}{"phase": "coordination"},
		})
		
		workflow = append(workflow, WorkflowStep{
			Agent:        "validator",
			Dependencies: []string{"coordinator"},
			Parallel:     false,
			Context:      map[string]interface{}{"phase": "validation"},
		})
	} else {
		// Single agent workflow
		workflow = append(workflow, WorkflowStep{
			Agent:        "single",
			Dependencies: []string{},
			Parallel:     false,
			Context:      map[string]interface{}{"phase": "execution"},
		})
	}
	
	fmt.Printf("📋 Workflow planejado com %d etapas\n", len(workflow))
	return workflow, nil
}

func (eo *EnhancedOrchestrator) executeWorkflow(ctx context.Context, workflow []WorkflowStep, input string) (string, error) {
	span := eo.observer.StartSpan("execute_workflow", map[string]string{
		"steps": fmt.Sprintf("%d", len(workflow)),
	})
	defer eo.observer.FinishSpan(span, true, nil)
	
	results := make(map[string]string)
	executed := make(map[string]bool)
	
	for len(executed) < len(workflow) {
		// Find steps that can be executed (dependencies satisfied)
		var readySteps []WorkflowStep
		
		for _, step := range workflow {
			if executed[step.Agent] {
				continue
			}
			
			canExecute := true
			for _, dep := range step.Dependencies {
				if !executed[dep] {
					canExecute = false
					break
				}
			}
			
			if canExecute {
				readySteps = append(readySteps, step)
			}
		}
		
		if len(readySteps) == 0 {
			return "", fmt.Errorf("workflow deadlock detected")
		}
		
		// Execute ready steps (potentially in parallel)
		if len(readySteps) > 1 {
			eo.executeStepsInParallel(ctx, readySteps, input, results)
		} else {
			result, err := eo.executeStep(ctx, readySteps[0], input, results)
			if err != nil {
				return "", err
			}
			results[readySteps[0].Agent] = result
		}
		
		// Mark steps as executed
		for _, step := range readySteps {
			executed[step.Agent] = true
		}
	}
	
	// Combine results
	finalResult := ""
	for agent, result := range results {
		finalResult += fmt.Sprintf("=== %s ===\n%s\n\n", agent, result)
	}
	
	return finalResult, nil
}

func (eo *EnhancedOrchestrator) executeStepsInParallel(ctx context.Context, steps []WorkflowStep, input string, results map[string]string) {
	var wg sync.WaitGroup
	resultsChan := make(chan struct {
		agent  string
		result string
		err    error
	}, len(steps))
	
	for _, step := range steps {
		wg.Add(1)
		go func(s WorkflowStep) {
			defer wg.Done()
			
			result, err := eo.executeStep(ctx, s, input, results)
			resultsChan <- struct {
				agent  string
				result string
				err    error
			}{s.Agent, result, err}
		}(step)
	}
	
	wg.Wait()
	close(resultsChan)
	
	// Collect results
	for res := range resultsChan {
		if res.err == nil {
			results[res.agent] = res.result
		}
	}
}

func (eo *EnhancedOrchestrator) executeStep(ctx context.Context, step WorkflowStep, input string, previousResults map[string]string) (string, error) {
	span := eo.observer.StartSpan("execute_step", map[string]string{
		"agent": step.Agent,
		"phase": fmt.Sprintf("%v", step.Context["phase"]),
	})
	defer eo.observer.FinishSpan(span, true, nil)
	
	// Build context-aware prompt
	prompt := eo.buildContextualPrompt(input, step, previousResults)
	
	// Execute using async processor
	resultChan := eo.processor.Submit(ctx, prompt)
	
	select {
	case result := <-resultChan:
		if result.Error != nil {
			return "", result.Error
		}
		return fmt.Sprintf("%v", result.Data), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (eo *EnhancedOrchestrator) buildContextualPrompt(input string, step WorkflowStep, previousResults map[string]string) string {
	prompt := fmt.Sprintf("Input: %s\n\nAgent: %s\nPhase: %v\n", input, step.Agent, step.Context["phase"])
	
	if len(previousResults) > 0 {
		prompt += "\nPrevious Results:\n"
		for agent, result := range previousResults {
			prompt += fmt.Sprintf("- %s: %s\n", agent, result)
		}
	}
	
	return prompt
}

func (cb *CircuitBreaker) CanExecute() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	
	switch cb.state {
	case Closed:
		return true
	case Open:
		if time.Since(cb.lastFailure) > cb.timeout {
			cb.state = HalfOpen
			return true
		}
		return false
	case HalfOpen:
		return true
	}
	return false
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	cb.failures = 0
	cb.state = Closed
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	cb.failures++
	cb.lastFailure = time.Now()
	
	if cb.failures >= cb.threshold {
		cb.state = Open
	}
}

func (eo *EnhancedOrchestrator) GetAdvancedInsights() map[string]interface{} {
	insights := eo.learning.GetInsights()
	metrics := eo.observer.GetMetrics()
	
	// Combine insights
	combined := make(map[string]interface{})
	for k, v := range insights {
		combined[k] = v
	}
	for k, v := range metrics {
		combined[fmt.Sprintf("metrics_%s", k)] = v
	}
	
	// Add cache statistics
	hits, misses := eo.cache.GetStats()
	combined["cache_hit_rate"] = float64(hits) / float64(hits+misses)
	
	return combined
}
