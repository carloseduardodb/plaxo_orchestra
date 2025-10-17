package orchestrator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plaxo-orchestra/internal/agent"
	"plaxo-orchestra/internal/analyzer"
	"plaxo-orchestra/internal/detector"
	"plaxo-orchestra/internal/pool"
	"strings"
	"time"
)

type Orchestrator struct {
	workingDir string
	agents     map[string]*agent.Agent
	agentPool  *pool.AgentPool
}

func New(workingDir string) *Orchestrator {
	return &Orchestrator{
		workingDir: workingDir,
		agents:     make(map[string]*agent.Agent),
		agentPool:  pool.NewAgentPool(),
	}
}

func (o *Orchestrator) Process(input string) error {
	projectInfo := detector.DetectProject(o.workingDir)
	
	switch projectInfo.Type {
	case detector.SingleAgent:
		return o.handleSingleAgent(input)
	case detector.MultiAgent:
		return o.handleMultiAgent(input, projectInfo.Domains)
	case detector.NewProject:
		if detector.IsComplexSoftwareRequest(input) {
			return o.createNewProject(input)
		}
		return o.handleSingleAgent(input)
	}
	
	return nil
}

func (o *Orchestrator) handleSingleAgent(input string) error {
	fmt.Println("🤖 Modo agente livre")
	
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, "q", "chat", input)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

func (o *Orchestrator) handleMultiAgent(input string, domains []string) error {
	fmt.Println("🎼 Modo multi-agente ativo")
	fmt.Printf("📁 Domínios: %s\n", strings.Join(domains, ", "))
	
	// Carrega agentes dos domínios
	for _, domain := range domains {
		if _, exists := o.agents[domain]; !exists {
			agent := agent.NewAgent(domain, o.workingDir, o.agentPool)
			if err := agent.LoadInstructions(); err != nil {
				fmt.Printf("⚠️  Erro carregando agente %s: %v\n", domain, err)
				continue
			}
			o.agents[domain] = agent
		}
	}
	
	// Analisa se precisa coordenação entre agentes
	if o.needsCoordination(input) {
		return o.coordinateAgents(input, domains)
	}
	
	// Determina qual agente deve responder
	targetAgent := o.selectAgent(input, domains)
	if targetAgent != "" {
		fmt.Printf("🎯 Delegando para: %s\n", targetAgent)
		result, err := o.agents[targetAgent].Execute(input)
		if err != nil {
			return err
		}
		fmt.Println(result)
	} else {
		fmt.Println("🤔 Delegando para todos os agentes relevantes")
		return o.coordinateAgents(input, domains)
	}
	
	return nil
}

func (o *Orchestrator) needsCoordination(input string) bool {
	keywords := []string{"integrar", "conectar", "comunicar", "sincronizar", "coordenar", "funcionar", "implementar sistema"}
	input = strings.ToLower(input)
	
	for _, keyword := range keywords {
		if strings.Contains(input, keyword) {
			return true
		}
	}
	return false
}

func (o *Orchestrator) coordinateAgents(input string, domains []string) error {
	fmt.Println("🔗 Coordenando múltiplos agentes...")
	
	// Fase 1: Análise - cada agente analisa o que precisa fazer
	analysisResults := make(map[string]string)
	for _, domain := range domains {
		if agent, exists := o.agents[domain]; exists {
			analysisPrompt := fmt.Sprintf(`
Analise esta requisição do ponto de vista do seu domínio (%s):
"%s"

Responda APENAS:
1. O que você precisa implementar/modificar no seu domínio
2. Que contratos/interfaces você precisa de outros domínios
3. Que contratos/interfaces você pode fornecer para outros domínios
`, domain, input)
			
			result, err := agent.Execute(analysisPrompt)
			if err != nil {
				fmt.Printf("❌ Erro na análise do %s: %v\n", domain, err)
				continue
			}
			analysisResults[domain] = result
			fmt.Printf("📋 %s analisou suas responsabilidades\n", domain)
		}
	}
	
	// Fase 2: Coordenação - compartilha análises entre agentes
	fmt.Println("🤝 Compartilhando análises entre agentes...")
	for _, domain := range domains {
		if agent, exists := o.agents[domain]; exists {
			coordinationPrompt := fmt.Sprintf(`
Baseado nas análises de todos os domínios, implemente sua parte:

Requisição original: "%s"

Análises dos outros domínios:
%s

Agora IMPLEMENTE concretamente sua parte, considerando as interfaces necessárias.
`, input, o.formatAnalysisResults(analysisResults, domain))
			
			result, err := agent.Execute(coordinationPrompt)
			if err != nil {
				fmt.Printf("❌ Erro na implementação do %s: %v\n", domain, err)
				continue
			}
			fmt.Printf("✅ %s implementou sua parte\n", domain)
			fmt.Println(result)
			fmt.Println("---")
		}
	}
	
	// Fase 3: Validação - verifica se tudo está integrado
	fmt.Println("🔍 Validando integração...")
	return o.validateIntegration(input, domains)
}

func (o *Orchestrator) formatAnalysisResults(results map[string]string, excludeDomain string) string {
	var formatted strings.Builder
	for domain, analysis := range results {
		if domain != excludeDomain {
			formatted.WriteString(fmt.Sprintf("\n=== %s ===\n%s\n", domain, analysis))
		}
	}
	return formatted.String()
}

func (o *Orchestrator) validateIntegration(input string, domains []string) error {
	validationPrompt := fmt.Sprintf(`
Valide se a implementação está completa e integrada:
Requisição: "%s"
Domínios implementados: %s

Verifique:
1. Todas as funcionalidades foram implementadas?
2. As integrações entre domínios estão corretas?
3. Há algum erro ou inconsistência?
4. O sistema está funcional?
`, input, strings.Join(domains, ", "))

	cmd := exec.Command("q", "chat", validationPrompt)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("⚠️  Erro na validação: %v\n", err)
		return nil
	}
	
	fmt.Println("🔍 Validação final:")
	fmt.Println(string(output))
	return nil
}

func (o *Orchestrator) createNewProject(input string) error {
	fmt.Println("🏗️  Gerando projeto com Amazon Q CLI...")
	
	// Primeiro, deixa o Amazon Q CLI gerar a estrutura do projeto
	cmd := exec.Command("q", "chat", input)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	
	fmt.Println("\n🔍 Analisando estrutura gerada para configurar multi-agentes...")
	
	// Analisa a estrutura criada pelo Q CLI
	boundedContexts, err := o.detectGeneratedStructure()
	if err != nil {
		return err
	}
	
	if len(boundedContexts) == 0 {
		fmt.Println("📁 Estrutura simples detectada, mantendo modo agente livre")
		return nil
	}
	
	fmt.Printf("📋 Configurando %d agentes especializados:\n", len(boundedContexts))
	
	// Configura agentes nas pastas já criadas pelo Q CLI
	var allContexts []string
	for _, bc := range boundedContexts {
		contextPath := fmt.Sprintf("%s/%s", bc.Domain, bc.Context)
		if err := o.setupAgentInExistingStructure(bc.Domain, bc.Context, bc.Description); err != nil {
			fmt.Printf("  ⚠️  Erro configurando %s: %v\n", contextPath, err)
			continue
		}
		allContexts = append(allContexts, contextPath)
		fmt.Printf("  ✅ %s: %s\n", contextPath, bc.Description)
	}
	
	fmt.Println("🎼 Multi-agentes configurados na estrutura existente")
	return nil
}

func (o *Orchestrator) detectGeneratedStructure() ([]analyzer.BoundedContext, error) {
	var contexts []analyzer.BoundedContext
	
	// Procura por pastas que parecem bounded contexts
	entries, err := os.ReadDir(o.workingDir)
	if err != nil {
		return contexts, err
	}
	
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		
		domainPath := filepath.Join(o.workingDir, entry.Name())
		
		// Verifica se é um domínio com sub-contextos
		subEntries, err := os.ReadDir(domainPath)
		if err != nil {
			continue
		}
		
		for _, subEntry := range subEntries {
			if subEntry.IsDir() && !strings.HasPrefix(subEntry.Name(), ".") {
				// Verifica se parece um bounded context (tem pastas como domain, application, etc)
				contextPath := filepath.Join(domainPath, subEntry.Name())
				if o.looksLikeBoundedContext(contextPath) {
					contexts = append(contexts, analyzer.BoundedContext{
						Domain:      entry.Name(),
						Context:     subEntry.Name(),
						Description: fmt.Sprintf("Especialista em %s/%s", entry.Name(), subEntry.Name()),
					})
				}
			}
		}
		
		// Se não tem sub-contextos, mas parece um domínio
		if len(subEntries) > 0 && o.looksLikeBoundedContext(domainPath) {
			contexts = append(contexts, analyzer.BoundedContext{
				Domain:      entry.Name(),
				Context:     "main",
				Description: fmt.Sprintf("Especialista em %s", entry.Name()),
			})
		}
	}
	
	return contexts, nil
}

func (o *Orchestrator) looksLikeBoundedContext(path string) bool {
	// Verifica se tem estrutura típica de bounded context
	indicators := []string{"domain", "application", "infrastructure", "src", "controllers", "services"}
	
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	
	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			for _, indicator := range indicators {
				if entry.Name() == indicator {
					count++
					break
				}
			}
		}
	}
	
	return count >= 2 // Se tem pelo menos 2 indicadores, provavelmente é um bounded context
}

func (o *Orchestrator) setupAgentInExistingStructure(domain, context, description string) error {
	var agentPath string
	
	if context == "main" {
		// Agente no nível do domínio
		agentPath = filepath.Join(o.workingDir, domain, "agents")
	} else {
		// Agente no bounded context
		agentPath = filepath.Join(o.workingDir, domain, context, "agents")
	}
	
	if err := os.MkdirAll(agentPath, 0755); err != nil {
		return err
	}
	
	instructions := fmt.Sprintf("Especialista em %s/%s: %s\n\nResponsável por implementar e manter funcionalidades específicas deste bounded context.", domain, context, description)
	instructionsFile := filepath.Join(agentPath, "instructions.txt")
	
	return os.WriteFile(instructionsFile, []byte(instructions), 0644)
}

func (o *Orchestrator) createBoundedContextStructure(domain, context, description string) error {
	contextPath := filepath.Join(o.workingDir, domain, context)
	agentPath := filepath.Join(contextPath, "agents")
	
	if err := os.MkdirAll(agentPath, 0755); err != nil {
		return err
	}
	
	instructions := fmt.Sprintf("Especialista em %s/%s: %s", domain, context, description)
	instructionsFile := filepath.Join(agentPath, "instructions.txt")
	
	return os.WriteFile(instructionsFile, []byte(instructions), 0644)
}

func (o *Orchestrator) selectAgent(input string, domains []string) string {
	input = strings.ToLower(input)
	
	bestMatch := ""
	maxScore := 0
	
	// Busca por palavras-chave nos nomes dos domínios/contextos
	for _, domain := range domains {
		score := 0
		parts := strings.Split(domain, "/")
		
		// Pontua por correspondência com partes do nome
		for _, part := range parts {
			if strings.Contains(input, part) {
				score += 2 // Correspondência exata vale mais
			}
			// Correspondência parcial
			if len(part) > 3 && strings.Contains(input, part[:3]) {
				score += 1
			}
		}
		
		if score > maxScore {
			maxScore = score
			bestMatch = domain
		}
	}
	
	return bestMatch
}
