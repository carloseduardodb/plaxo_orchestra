package orchestrator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (o *Orchestrator) StartWatchMode() error {
	fmt.Println("👁️ Modo watch ativo - monitorando mudanças...")
	
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	lastModTimes := make(map[string]time.Time)
	
	for {
		select {
		case <-ticker.C:
			if err := o.checkForChanges(lastModTimes); err != nil {
				log.Printf("Erro no watch: %v", err)
			}
		}
	}
}

func (o *Orchestrator) checkForChanges(lastModTimes map[string]time.Time) error {
	return filepath.Walk(o.workingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Ignora pastas de agentes e arquivos temporários
		if strings.Contains(path, "/agents/") || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Só monitora arquivos de código
		if !o.isCodeFile(path) {
			return nil
		}

		lastMod, exists := lastModTimes[path]
		if !exists {
			lastModTimes[path] = info.ModTime()
			return nil
		}

		if info.ModTime().After(lastMod) {
			lastModTimes[path] = info.ModTime()
			go o.handleFileChange(path)
		}

		return nil
	})
}

func (o *Orchestrator) isCodeFile(path string) bool {
	extensions := []string{".go", ".js", ".py", ".java", ".ts", ".rs", ".cpp", ".c", ".rb"}
	for _, ext := range extensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}

func (o *Orchestrator) handleFileChange(filePath string) {
	fmt.Printf("📝 Arquivo modificado: %s\n", filePath)
	
	// Determina qual agente deve revisar
	domain := o.getDomainFromPath(filePath)
	if domain == "" {
		return
	}

	if agent, exists := o.agents[domain]; exists {
		prompt := fmt.Sprintf("Arquivo %s foi modificado. Revise se há problemas ou melhorias necessárias.", filePath)
		
		result, err := agent.Execute(prompt)
		if err != nil {
			fmt.Printf("❌ Erro na revisão automática: %v\n", err)
			return
		}
		
		if strings.Contains(strings.ToLower(result), "problema") || 
		   strings.Contains(strings.ToLower(result), "erro") {
			fmt.Printf("⚠️ Revisão automática encontrou issues em %s:\n%s\n", filePath, result)
		}
	}
}

func (o *Orchestrator) getDomainFromPath(filePath string) string {
	rel, err := filepath.Rel(o.workingDir, filePath)
	if err != nil {
		return ""
	}
	
	parts := strings.Split(rel, string(filepath.Separator))
	if len(parts) >= 2 {
		return fmt.Sprintf("%s/%s", parts[0], parts[1])
	} else if len(parts) >= 1 {
		return parts[0]
	}
	
	return ""
}