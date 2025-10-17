package stream

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type StreamHandler struct {
	onProgress func(string)
	onComplete func(string)
	onError    func(error)
}

type StreamResult struct {
	Content string
	Error   error
}

func NewStreamHandler() *StreamHandler {
	return &StreamHandler{
		onProgress: func(s string) { fmt.Print(s) },
		onComplete: func(s string) {},
		onError:    func(e error) { fmt.Printf("❌ Erro: %v\n", e) },
	}
}

func (sh *StreamHandler) SetProgressCallback(fn func(string)) {
	sh.onProgress = fn
}

func (sh *StreamHandler) SetCompleteCallback(fn func(string)) {
	sh.onComplete = fn
}

func (sh *StreamHandler) SetErrorCallback(fn func(error)) {
	sh.onError = fn
}

func (sh *StreamHandler) ExecuteWithStream(ctx context.Context, request string) *StreamResult {
	result := &StreamResult{}
	
	// Show initial progress
	sh.onProgress("🧠 Analisando requisição...\n")
	time.Sleep(200 * time.Millisecond)
	
	sh.onProgress("🎯 Iniciando processamento...\n")
	time.Sleep(200 * time.Millisecond)
	
	// Execute Q CLI with streaming
	cmd := exec.CommandContext(ctx, "q", "chat", "--no-interactive", request)
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		result.Error = err
		sh.onError(err)
		return result
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		result.Error = err
		sh.onError(err)
		return result
	}
	
	if err := cmd.Start(); err != nil {
		result.Error = err
		sh.onError(err)
		return result
	}
	
	// Stream output in real-time
	var content strings.Builder
	
	// Read stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			content.WriteString(line + "\n")
			sh.onProgress(line + "\n")
		}
	}()
	
	// Read stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "Error") || strings.Contains(line, "error") {
				sh.onProgress("⚠️  " + line + "\n")
			}
		}
	}()
	
	// Wait for completion
	if err := cmd.Wait(); err != nil {
		result.Error = err
		sh.onError(err)
		return result
	}
	
	result.Content = content.String()
	sh.onComplete(result.Content)
	
	return result
}

func (sh *StreamHandler) ExecuteWorkflowWithStream(ctx context.Context, steps []string, input string) *StreamResult {
	result := &StreamResult{}
	var fullContent strings.Builder
	
	sh.onProgress(fmt.Sprintf("📋 Executando workflow com %d etapas...\n\n", len(steps)))
	
	for i, step := range steps {
		sh.onProgress(fmt.Sprintf("🔄 Etapa %d/%d: %s\n", i+1, len(steps), step))
		sh.onProgress(strings.Repeat("─", 50) + "\n")
		
		// Simulate processing time for each step
		stepResult := sh.ExecuteWithStream(ctx, fmt.Sprintf("%s - %s", input, step))
		if stepResult.Error != nil {
			result.Error = stepResult.Error
			return result
		}
		
		fullContent.WriteString(fmt.Sprintf("=== %s ===\n%s\n\n", step, stepResult.Content))
		
		if i < len(steps)-1 {
			sh.onProgress("\n✅ Etapa concluída. Próxima etapa...\n\n")
			time.Sleep(500 * time.Millisecond)
		}
	}
	
	result.Content = fullContent.String()
	sh.onProgress("\n🎉 Workflow concluído com sucesso!\n")
	
	return result
}

// StreamingProgressBar shows a simple progress indicator
func ShowProgressBar(current, total int, description string) {
	percentage := float64(current) / float64(total) * 100
	filled := int(percentage / 5) // 20 chars max
	
	bar := "["
	for i := 0; i < 20; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	bar += "]"
	
	fmt.Printf("\r%s %.1f%% %s", bar, percentage, description)
	if current == total {
		fmt.Println()
	}
}

// StreamingSpinner shows a spinning indicator
type Spinner struct {
	chars   []string
	current int
	active  bool
}

func NewSpinner() *Spinner {
	return &Spinner{
		chars: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

func (s *Spinner) Start(message string) {
	s.active = true
	go func() {
		for s.active {
			fmt.Printf("\r%s %s", s.chars[s.current], message)
			s.current = (s.current + 1) % len(s.chars)
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Print("\r")
	}()
}

func (s *Spinner) Stop() {
	s.active = false
	time.Sleep(150 * time.Millisecond) // Allow goroutine to finish
}
