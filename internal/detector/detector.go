package detector

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ProjectType int

const (
	SingleAgent ProjectType = iota
	MultiAgent
	NewProject
)

type ProjectInfo struct {
	Type    ProjectType
	Domains []string
	Root    string
}

func DetectProject(workingDir string) *ProjectInfo {
	// Verifica se é um projeto multi-agente existente
	if isMultiAgentProject(workingDir) {
		domains := findDomains(workingDir)
		return &ProjectInfo{
			Type:    MultiAgent,
			Domains: domains,
			Root:    workingDir,
		}
	}
	
	// Verifica se precisa criar novo projeto
	if needsNewProject(workingDir) {
		return &ProjectInfo{
			Type: NewProject,
			Root: workingDir,
		}
	}
	
	return &ProjectInfo{
		Type: SingleAgent,
		Root: workingDir,
	}
}

func isMultiAgentProject(dir string) bool {
	// Procura por estrutura de domínios com agentes
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			agentPath := filepath.Join(dir, entry.Name(), "agents", "instructions.txt")
			if _, err := os.Stat(agentPath); err == nil {
				return true
			}
		}
	}
	
	return false
}

func findDomains(dir string) []string {
	var domains []string
	
	entries, err := os.ReadDir(dir)
	if err != nil {
		return domains
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			// Verifica se tem agentes no nível do domínio
			agentPath := filepath.Join(dir, entry.Name(), "agents", "instructions.txt")
			if _, err := os.Stat(agentPath); err == nil {
				domains = append(domains, entry.Name())
				continue
			}
			
			// Verifica bounded contexts dentro do domínio
			domainPath := filepath.Join(dir, entry.Name())
			subEntries, err := os.ReadDir(domainPath)
			if err != nil {
				continue
			}
			
			for _, subEntry := range subEntries {
				if subEntry.IsDir() {
					subAgentPath := filepath.Join(domainPath, subEntry.Name(), "agents", "instructions.txt")
					if _, err := os.Stat(subAgentPath); err == nil {
						domains = append(domains, entry.Name()+"/"+subEntry.Name())
					}
				}
			}
		}
	}
	
	return domains
}

func needsNewProject(dir string) bool {
	// Verifica se o diretório está vazio ou tem poucos arquivos
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	
	// Se tem menos de 3 arquivos, provavelmente precisa de novo projeto
	return len(entries) < 3
}

func IsComplexSoftwareRequest(input string) bool {
	prompt := fmt.Sprintf(`
Analise esta requisição e responda APENAS "SIM" ou "NAO":
"%s"

É um software complexo que precisa de múltiplos domínios/módulos coordenados?
Considere complexo se envolve:
- Múltiplas entidades de negócio
- Diferentes responsabilidades
- Integrações entre módulos
- Arquitetura com vários componentes

Responda apenas: SIM ou NAO
`, input)

	cmd := exec.Command("q", "chat", "--message", prompt)
	output, err := cmd.Output()
	if err != nil {
		// Fallback: se Q CLI falhar, usa heurística simples
		return len(strings.Fields(input)) > 5
	}

	response := strings.ToUpper(strings.TrimSpace(string(output)))
	return strings.Contains(response, "SIM")
}
