package pool

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

type Connection struct {
	ID        string
	CreatedAt time.Time
	LastUsed  time.Time
	InUse     bool
}

type ConnectionPool struct {
	connections chan *Connection
	active      map[string]*Connection
	mutex       sync.RWMutex
	maxSize     int
	timeout     time.Duration
}

type AsyncProcessor struct {
	pool     *ConnectionPool
	jobQueue chan Job
	workers  int
}

type Job struct {
	ID       string
	Request  string
	Response chan JobResult
	Context  context.Context
}

type JobResult struct {
	Data  interface{}
	Error error
}

func NewConnectionPool(maxSize int, timeout time.Duration) *ConnectionPool {
	pool := &ConnectionPool{
		connections: make(chan *Connection, maxSize),
		active:      make(map[string]*Connection),
		maxSize:     maxSize,
		timeout:     timeout,
	}
	
	// Initialize connections
	for i := 0; i < maxSize; i++ {
		conn := &Connection{
			ID:       fmt.Sprintf("conn_%d", i),
			InUse:    false,
			LastUsed: time.Now(),
		}
		pool.connections <- conn
	}
	
	return pool
}

func NewAsyncProcessor(pool *ConnectionPool, workers int) *AsyncProcessor {
	processor := &AsyncProcessor{
		pool:     pool,
		jobQueue: make(chan Job, workers*2),
		workers:  workers,
	}
	
	for i := 0; i < workers; i++ {
		go processor.worker()
	}
	
	return processor
}

func (p *ConnectionPool) Get(ctx context.Context) (*Connection, error) {
	select {
	case conn := <-p.connections:
		p.mutex.Lock()
		conn.InUse = true
		conn.LastUsed = time.Now()
		p.active[conn.ID] = conn
		p.mutex.Unlock()
		return conn, nil
	case <-time.After(p.timeout):
		return nil, fmt.Errorf("connection timeout")
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (p *ConnectionPool) Release(conn *Connection) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	conn.InUse = false
	delete(p.active, conn.ID)
	
	select {
	case p.connections <- conn:
	default:
		// Pool full, discard connection
	}
}

func (ap *AsyncProcessor) Submit(ctx context.Context, request string) <-chan JobResult {
	job := Job{
		ID:       fmt.Sprintf("job_%d", time.Now().UnixNano()),
		Request:  request,
		Response: make(chan JobResult, 1),
		Context:  ctx,
	}
	
	select {
	case ap.jobQueue <- job:
	case <-ctx.Done():
		job.Response <- JobResult{Error: ctx.Err()}
	}
	
	return job.Response
}

func (ap *AsyncProcessor) worker() {
	for job := range ap.jobQueue {
		result := ap.processJob(job)
		select {
		case job.Response <- result:
		case <-job.Context.Done():
		}
	}
}

func (ap *AsyncProcessor) processJob(job Job) JobResult {
	conn, err := ap.pool.Get(job.Context)
	if err != nil {
		return JobResult{Error: err}
	}
	defer ap.pool.Release(conn)
	
	// Execute Amazon Q CLI
	ctx, cancel := context.WithTimeout(job.Context, 45*time.Second)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, "q", "chat", "--no-interactive", job.Request)
	output, err := cmd.Output()
	if err != nil {
		return JobResult{Error: fmt.Errorf("Q CLI error: %v", err)}
	}
	
	return JobResult{
		Data: string(output),
	}
}
