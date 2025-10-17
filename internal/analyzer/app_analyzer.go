package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type AppStructure struct {
	RootPath    string
	Domains     map[string]*Domain
	TechStack   []string
	Complexity  string
	AgentPlan   map[string][]string
}

type Domain struct {
	Name        string
	Path        string
	Files       []string
	SubDomains  map[string]*Domain
	Complexity  int
	AgentNeeded bool
}

type AppAnalyzer struct {
	rootPath string
}

func NewAppAnalyzer(rootPath string) *AppAnalyzer {
	return &AppAnalyzer{rootPath: rootPath}
}

func (aa *AppAnalyzer) AnalyzeApplication() (*AppStructure, error) {
	fmt.Println("🔍 Analisando estrutura da aplicação...")
	
	structure := &AppStructure{
		RootPath:  aa.rootPath,
		Domains:   make(map[string]*Domain),
		AgentPlan: make(map[string][]string),
	}
	
	// Detectar tech stack
	structure.TechStack = aa.detectTechStack()
	fmt.Printf("📚 Tech Stack detectado: %v\n", structure.TechStack)
	
	// Analisar estrutura de diretórios
	if err := aa.analyzeDomains(structure); err != nil {
		return nil, err
	}
	
	// Calcular complexidade
	structure.Complexity = aa.calculateComplexity(structure)
	fmt.Printf("📊 Complexidade: %s\n", structure.Complexity)
	
	// Planejar distribuição de agentes
	aa.planAgentDistribution(structure)
	
	return structure, nil
}

func (aa *AppAnalyzer) detectTechStack() []string {
	var stack []string
	
	// Detectar linguagens e frameworks
	patterns := map[string][]string{
		"Python":     {"*.py", "requirements.txt", "setup.py", "pyproject.toml"},
		"JavaScript": {"*.js", "package.json", "*.ts", "*.jsx", "*.tsx"},
		"Go":         {"*.go", "go.mod", "go.sum"},
		"Java":       {"*.java", "pom.xml", "build.gradle"},
		"PHP":        {"*.php", "composer.json"},
		"Ruby":       {"*.rb", "Gemfile"},
		"C#":         {"*.cs", "*.csproj", "*.sln"},
		"Rust":       {"*.rs", "Cargo.toml"},
		"FastAPI":    {"main.py", "app.py", "**/routers/**"},
		"Django":     {"manage.py", "settings.py", "**/models.py"},
		"Flask":      {"app.py", "**/templates/**"},
		"React":      {"src/App.js", "src/App.tsx", "public/index.html"},
		"Vue":        {"src/App.vue", "vue.config.js"},
		"Angular":    {"angular.json", "src/app/app.module.ts"},
		"Docker":     {"Dockerfile", "docker-compose.yml"},
		"Kubernetes": {"*.yaml", "*.yml", "**/k8s/**"},
	}
	
	for tech, patterns := range patterns {
		if aa.hasPatterns(patterns) {
			stack = append(stack, tech)
		}
	}
	
	return stack
}

func (aa *AppAnalyzer) hasPatterns(patterns []string) bool {
	for _, pattern := range patterns {
		matches, _ := filepath.Glob(filepath.Join(aa.rootPath, pattern))
		if len(matches) > 0 {
			return true
		}
		
		// Busca recursiva para padrões com **
		if strings.Contains(pattern, "**") {
			found := false
			filepath.Walk(aa.rootPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				relPath, _ := filepath.Rel(aa.rootPath, path)
				if matched, _ := filepath.Match(strings.Replace(pattern, "**", "*", -1), relPath); matched {
					found = true
					return filepath.SkipDir
				}
				return nil
			})
			if found {
				return true
			}
		}
	}
	return false
}

func (aa *AppAnalyzer) analyzeDomains(structure *AppStructure) error {
	// Padrões de domínios comuns
	domainPatterns := map[string][]string{
		"auth":     {"auth", "authentication", "login", "users", "accounts"},
		"api":      {"api", "routes", "controllers", "handlers", "endpoints"},
		"models":   {"models", "entities", "schemas", "database", "db"},
		"services": {"services", "business", "logic", "core"},
		"utils":    {"utils", "helpers", "common", "shared", "lib"},
		"config":   {"config", "settings", "env", "configuration"},
		"tests":    {"tests", "test", "spec", "__tests__", "testing"},
		"docs":     {"docs", "documentation", "readme"},
		"frontend": {"frontend", "ui", "web", "client", "public", "static"},
		"backend":  {"backend", "server", "api"},
		"data":     {"data", "migrations", "seeds", "fixtures"},
		"deploy":   {"deploy", "deployment", "infra", "infrastructure", "k8s", "docker"},
		"products": {"products", "catalog", "items"},
		"orders":   {"orders", "cart", "checkout"},
		"payment":  {"payment", "billing", "transactions"},
	}
	
	return filepath.Walk(aa.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return nil
		}
		
		relPath, _ := filepath.Rel(aa.rootPath, path)
		if relPath == "." || strings.HasPrefix(relPath, ".") {
			return nil
		}
		
		dirName := strings.ToLower(info.Name())
		
		for domain, patterns := range domainPatterns {
			for _, pattern := range patterns {
				if strings.Contains(dirName, pattern) || dirName == pattern {
					if structure.Domains[domain] == nil {
						structure.Domains[domain] = &Domain{
							Name:       domain,
							Path:       path,
							Files:      []string{},
							SubDomains: make(map[string]*Domain),
						}
					}
					
					// Contar arquivos no domínio
					files := aa.countFiles(path)
					structure.Domains[domain].Files = append(structure.Domains[domain].Files, files...)
					structure.Domains[domain].Complexity += len(files)
					
					// Sempre marcar como necessário se há arquivos
					if len(files) > 0 {
						structure.Domains[domain].AgentNeeded = true
					}
					
					break
				}
			}
		}
		
		return nil
	})
}

func (aa *AppAnalyzer) countFiles(dirPath string) []string {
	var files []string
	
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		
		ext := strings.ToLower(filepath.Ext(path))
		codeExts := []string{".py", ".js", ".ts", ".go", ".java", ".php", ".rb", ".cs", ".rs", ".jsx", ".tsx", ".vue"}
		
		for _, codeExt := range codeExts {
			if ext == codeExt {
				files = append(files, path)
				break
			}
		}
		
		return nil
	})
	
	return files
}

func (aa *AppAnalyzer) calculateComplexity(structure *AppStructure) string {
	totalFiles := 0
	totalDomains := len(structure.Domains)
	
	for _, domain := range structure.Domains {
		totalFiles += len(domain.Files)
	}
	
	if totalFiles < 10 && totalDomains < 3 {
		return "simple"
	} else if totalFiles < 50 && totalDomains < 8 {
		return "medium"
	} else {
		return "complex"
	}
}

func (aa *AppAnalyzer) planAgentDistribution(structure *AppStructure) {
	fmt.Println("\n🤖 Planejando distribuição de agentes...")
	
	for domainName, domain := range structure.Domains {
		// Sempre criar agente se há arquivos de código
		if len(domain.Files) > 0 {
			agentPath := filepath.Join(domain.Path, "agents")
			structure.AgentPlan[domainName] = []string{agentPath}
			
			fmt.Printf("  📍 %s: %s (%d arquivos)\n", domainName, agentPath, len(domain.Files))
		}
	}
	
	// Agente principal se complexidade alta ou múltiplos domínios
	if structure.Complexity == "complex" || len(structure.Domains) > 2 {
		mainAgentPath := filepath.Join(structure.RootPath, "orchestra_agents")
		structure.AgentPlan["orchestrator"] = []string{mainAgentPath}
		fmt.Printf("  🎼 orchestrator: %s (coordenação geral)\n", mainAgentPath)
	}
}

func (aa *AppAnalyzer) DeployAgents(structure *AppStructure) error {
	fmt.Println("\n🚀 Distribuindo agentes pela aplicação...")
	
	for domain, paths := range structure.AgentPlan {
		for _, agentPath := range paths {
			if err := aa.createAgentStructure(domain, agentPath, structure); err != nil {
				return fmt.Errorf("erro criando agente %s: %v", domain, err)
			}
			fmt.Printf("✅ Agente %s criado em: %s\n", domain, agentPath)
		}
	}
	
	// Criar arquivo de configuração central
	if err := aa.createOrchestraConfig(structure); err != nil {
		return fmt.Errorf("erro criando configuração: %v", err)
	}
	
	fmt.Println("🎉 Distribuição de agentes concluída!")
	return nil
}

func (aa *AppAnalyzer) createAgentStructure(domain, agentPath string, structure *AppStructure) error {
	// Criar diretório do agente
	if err := os.MkdirAll(agentPath, 0755); err != nil {
		return err
	}
	
	// Criar arquivo de configuração do agente
	configContent := aa.generateAgentConfig(domain, structure)
	configFile := filepath.Join(agentPath, "agent.yaml")
	
	return os.WriteFile(configFile, []byte(configContent), 0644)
}

func (aa *AppAnalyzer) generateAgentConfig(domain string, structure *AppStructure) string {
	domainInfo := structure.Domains[domain]
	
	config := fmt.Sprintf(`# Agente especializado para domínio: %s
name: %s_agent
domain: %s
complexity: %d
files_count: %d
tech_stack: %v

# Responsabilidades
responsibilities:
  - Análise de código do domínio %s
  - Refatoração e otimização
  - Testes e validação
  - Documentação técnica

# Contexto do domínio
context:
  path: %s
  files: %d
  
# Comandos especializados
commands:
  analyze: "Analisar código do domínio %s"
  refactor: "Refatorar código seguindo melhores práticas"
  test: "Criar/executar testes para o domínio"
  document: "Gerar documentação técnica"
`, 
		domain, domain, domain, 
		domainInfo.Complexity, len(domainInfo.Files), structure.TechStack,
		domain, domainInfo.Path, len(domainInfo.Files), domain)
	
	return config
}

func (aa *AppAnalyzer) createOrchestraConfig(structure *AppStructure) error {
	configPath := filepath.Join(structure.RootPath, "orchestra.yaml")
	
	config := fmt.Sprintf(`# Configuração do Plaxo Orchestra
app_name: %s
complexity: %s
tech_stack: %v
total_domains: %d

# Agentes distribuídos
agents:
`, filepath.Base(structure.RootPath), structure.Complexity, structure.TechStack, len(structure.Domains))
	
	for domain, paths := range structure.AgentPlan {
		config += fmt.Sprintf("  %s:\n", domain)
		for _, path := range paths {
			config += fmt.Sprintf("    - %s\n", path)
		}
	}
	
	config += `
# Comandos de orquestração
orchestration:
  analyze_all: "Analisar todos os domínios"
  refactor_all: "Refatorar aplicação completa"
  test_all: "Executar todos os testes"
  deploy_all: "Preparar deploy da aplicação"
`
	
	return os.WriteFile(configPath, []byte(config), 0644)
}
