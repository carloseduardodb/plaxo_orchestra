package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"plaxo-orchestra/internal/agent"
	"plaxo-orchestra/internal/analyzer"
	"plaxo-orchestra/internal/detector"
	"plaxo-orchestra/internal/intelligence"
	"plaxo-orchestra/internal/pool"
	"strings"
)

type SmartOrchestrator struct {
	workingDir   string
	agents       map[string]*agent.Agent
	semantic     *intelligence.SemanticAnalyzer
	coordinator  *intelligence.Coordinator
	learning     *intelligence.LearningSystem
	agentPool    *pool.AgentPool
}

func NewSmart(workingDir string) *SmartOrchestrator {
	return &SmartOrchestrator{
		workingDir:  workingDir,
		agents:      make(map[string]*agent.Agent),
		semantic:    intelligence.NewSemanticAnalyzer(),
		coordinator: intelligence.NewCoordinator(),
		learning:    intelligence.NewLearningSystem(workingDir),
		agentPool:   pool.NewAgentPool(),
	}
}

func (o *SmartOrchestrator) Process(input string) error {
	fmt.Println("🧠 Analisando requisição com IA...")
	
	// Análise semântica da requisição
	analysis, err := o.semantic.AnalyzeIntent(input)
	if err != nil {
		fmt.Printf("⚠️  Erro na análise semântica, usando modo básico: %v\n", err)
		return o.processBasic(input)
	}

	fmt.Printf("🎯 Intent: %s | Complexidade: %s | Domínios: %s\n", 
		analysis.Intent, analysis.Complexity, strings.Join(analysis.Domains, ", "))

	projectInfo := detector.DetectProject(o.workingDir)
	
	switch projectInfo.Type {
	case detector.SingleAgent:
		return o.handleSingleAgent(input)
	case detector.MultiAgent:
		return o.handleSmartMultiAgent(input, projectInfo.Domains, analysis)
	case detector.NewProject:
		if analysis.Complexity != "simple" {
			return o.createSmartProject(input, analysis)
		}
		return o.handleSingleAgent(input)
	}
	
	return nil
}

func (o *SmartOrchestrator) handleSmartMultiAgent(input string, domains []string, analysis *intelligence.SemanticResult) error {
	fmt.Println("🎼 Modo multi-agente inteligente ativo")
	
	// Carrega agentes
	for _, domain := range domains {
		if _, exists := o.agents[domain]; !exists {
			agent := agent.NewAgent(domain, o.workingDir)
			if err := agent.LoadInstructions(); err != nil {
				fmt.Printf("⚠️  Erro carregando agente %s: %v\n", domain, err)
				continue
			}
			o.agents[domain] = agent
		}
	}

	// Verifica se precisa coordenação baseado na análise semântica
	if analysis.Complexity == "complex" || analysis.Intent == "integrate" {
		return o.executeSmartWorkflow(input, domains, analysis)
	}

	// Seleção inteligente de agente
	selectedAgent := o.selectSmartAgent(input, domains, analysis)
	if selectedAgent != "" {
		fmt.Printf("🎯 Agente selecionado: %s (IA)\n", selectedAgent)
		
		// Registra decisão para aprendizado
		context := map[string]string{
			"intent":     analysis.Intent,
			"complexity": analysis.Complexity,
			"domains":    strings.Join(analysis.Domains, ","),
		}
		o.learning.RecordDecision(input, selectedAgent, context)
		
		result, err := o.agents[selectedAgent].Execute(input)
		if err != nil {
			o.learning.RecordFeedback(input, false, fmt.Sprintf("Erro: %v", err))
			return err
		}
		
		fmt.Println(result)
		o.learning.RecordFeedback(input, true, "Execução bem-sucedida")
		return nil
	}

	// Fallback para coordenação
	return o.executeSmartWorkflow(input, domains, analysis)
}

func (o *SmartOrchestrator) selectSmartAgent(input string, domains []string, analysis *intelligence.SemanticResult) string {
	// Primeiro, tenta usar aprendizado histórico
	bestFromHistory := o.learning.GetBestAgentForInput(input, domains)
	if bestFromHistory != "" {
		fmt.Printf("📚 Usando aprendizado histórico: %s\n", bestFromHistory)
		return bestFromHistory
	}

	// Usa análise semântica para seleção
	bestAgent := ""
	maxScore := 0.0

	for _, domain := range domains {
		score := o.semantic.CalculateSimilarity(input, domain)
		
		// Bonus por correspondência com domínios detectados
		for _, detectedDomain := range analysis.Domains {
			if strings.Contains(domain, detectedDomain) {
				score += 0.3
			}
		}

		// Bonus por correspondência com entidades
		for _, entity := range analysis.Entities {
			if strings.Contains(domain, entity) {
				score += 0.2
			}
		}

		if score > maxScore {
			maxScore = score
			bestAgent = domain
		}
	}

	if maxScore > 0.4 {
		return bestAgent
	}

	return ""
}

func (o *SmartOrchestrator) executeSmartWorkflow(input string, domains []string, analysis *intelligence.SemanticResult) error {
	fmt.Println("🔗 Executando workflow inteligente...")

	// Planeja workflow baseado na análise semântica
	workflow, err := o.coordinator.PlanWorkflow(input, domains)
	if err != nil {
		fmt.Printf("⚠️  Erro no planejamento, usando coordenação básica: %v\n", err)
		return o.coordinateBasic(input, domains)
	}

	fmt.Printf("📋 Workflow planejado com %d etapas\n", len(workflow.Steps))

	// Executa workflow
	agentExecutor := func(agentName, prompt string) (string, error) {
		if agent, exists := o.agents[agentName]; exists {
			return agent.Execute(prompt)
		}
		return "", fmt.Errorf("agente %s não encontrado", agentName)
	}

	return o.coordinator.ExecuteWorkflow(workflow, agentExecutor)
}

func (o *SmartOrchestrator) createSmartProject(input string, analysis *intelligence.SemanticResult) error {
	fmt.Println("🏗️  Criando projeto inteligente...")
	
	// Gera projeto com Amazon Q CLI usando pool
	response, err := o.agentPool.Execute("project_creator", input)
	if err != nil {
		return err
	}
	fmt.Print(response)

	fmt.Println("\n🧠 Analisando estrutura com IA...")

	// Detecta bounded contexts usando análise semântica
	boundedContexts, err := o.detectSmartStructure(analysis)
	if err != nil {
		return err
	}

	if len(boundedContexts) == 0 {
		fmt.Println("📁 Estrutura simples detectada")
		return nil
	}

	fmt.Printf("🎯 Configurando %d agentes especializados:\n", len(boundedContexts))

	// Configura agentes com instruções inteligentes
	for _, bc := range boundedContexts {
		if err := o.setupSmartAgent(bc, analysis); err != nil {
			fmt.Printf("  ⚠️  Erro configurando %s/%s: %v\n", bc.Domain, bc.Context, err)
			continue
		}
		fmt.Printf("  ✅ %s/%s: %s\n", bc.Domain, bc.Context, bc.Description)
	}

	fmt.Println("🎼 Multi-agentes inteligentes configurados!")
	return nil
}

func (o *SmartOrchestrator) detectSmartStructure(analysis *intelligence.SemanticResult) ([]analyzer.BoundedContext, error) {
	var contexts []analyzer.BoundedContext

	// Usa análise semântica para detectar bounded contexts mais precisos
	entries, err := os.ReadDir(o.workingDir)
	if err != nil {
		return contexts, err
	}

	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		domainPath := filepath.Join(o.workingDir, entry.Name())
		
		// Verifica se o diretório corresponde aos domínios detectados
		relevantDomain := false
		for _, detectedDomain := range analysis.Domains {
			if strings.Contains(entry.Name(), detectedDomain) || 
			   o.semantic.CalculateSimilarity(entry.Name(), detectedDomain) > 0.5 {
				relevantDomain = true
				break
			}
		}

		if relevantDomain || o.looksLikeBoundedContext(domainPath) {
			// Detecta sub-contextos
			subEntries, err := os.ReadDir(domainPath)
			if err != nil {
				continue
			}

			hasSubContexts := false
			for _, subEntry := range subEntries {
				if subEntry.IsDir() && !strings.HasPrefix(subEntry.Name(), ".") {
					contextPath := filepath.Join(domainPath, subEntry.Name())
					if o.looksLikeBoundedContext(contextPath) {
						contexts = append(contexts, analyzer.BoundedContext{
							Domain:      entry.Name(),
							Context:     subEntry.Name(),
							Description: o.generateSmartDescription(entry.Name(), subEntry.Name(), analysis),
						})
						hasSubContexts = true
					}
				}
			}

			// Se não tem sub-contextos, mas é relevante
			if !hasSubContexts && relevantDomain {
				contexts = append(contexts, analyzer.BoundedContext{
					Domain:      entry.Name(),
					Context:     "main",
					Description: o.generateSmartDescription(entry.Name(), "main", analysis),
				})
			}
		}
	}

	return contexts, nil
}

func (o *SmartOrchestrator) generateSmartDescription(domain, context string, analysis *intelligence.SemanticResult) string {
	// Gera descrição baseada na análise semântica
	description := fmt.Sprintf("Especialista em %s", domain)
	
	if context != "main" {
		description += fmt.Sprintf("/%s", context)
	}

	// Adiciona contexto baseado nas entidades detectadas
	if len(analysis.Entities) > 0 {
		description += fmt.Sprintf(" - Foca em: %s", strings.Join(analysis.Entities, ", "))
	}

	return description
}

func (o *SmartOrchestrator) setupSmartAgent(bc analyzer.BoundedContext, analysis *intelligence.SemanticResult) error {
	var agentPath string
	
	if bc.Context == "main" {
		agentPath = filepath.Join(o.workingDir, bc.Domain, "agents")
	} else {
		agentPath = filepath.Join(o.workingDir, bc.Domain, bc.Context, "agents")
	}
	
	if err := os.MkdirAll(agentPath, 0755); err != nil {
		return err
	}

	// Gera instruções inteligentes baseadas na análise
	instructions := o.generateSmartInstructions(bc, analysis)
	instructionsFile := filepath.Join(agentPath, "instructions.txt")
	
	return os.WriteFile(instructionsFile, []byte(instructions), 0644)
}

func (o *SmartOrchestrator) generateSmartInstructions(bc analyzer.BoundedContext, analysis *intelligence.SemanticResult) string {
	instructions := fmt.Sprintf(`%s

RESPONSABILIDADES:
- Implementar funcionalidades específicas de %s/%s
- Manter consistência com outros bounded contexts
- Seguir padrões de arquitetura limpa

CONTEXTO DA REQUISIÇÃO:
- Intent: %s
- Entidades relevantes: %s
- Complexidade: %s

DIRETRIZES:
- Sempre considere integrações com outros domínios
- Implemente testes quando apropriado
- Documente APIs e contratos
- Siga princípios SOLID e DDD
`, bc.Description, bc.Domain, bc.Context, 
   analysis.Intent, strings.Join(analysis.Entities, ", "), analysis.Complexity)

	return instructions
}

func (o *SmartOrchestrator) looksLikeBoundedContext(path string) bool {
	indicators := []string{"domain", "application", "infrastructure", "src", "controllers", "services", "models", "handlers"}
	
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
	
	return count >= 2
}

func (o *SmartOrchestrator) handleSingleAgent(input string) error {
	fmt.Println("🤖 Modo agente livre")
	
	response, err := o.agentPool.Execute("free_agent", input)
	if err != nil {
		return err
	}
	fmt.Print(response)
	
	return nil
}

func (o *SmartOrchestrator) executeAgent(agentPath, input string) (string, error) {
	agentID := strings.ReplaceAll(agentPath, "/", "_")
	
	// Contexto específico do agente
	contextPrompt := fmt.Sprintf(`Você é um agente especializado em: %s
Diretório de trabalho: %s
Contexto do agente: %s

Requisição: %s`, agentPath, o.workingDir, agentPath, input)
	
	return o.agentPool.Execute(agentID, contextPrompt)
}

func (o *SmartOrchestrator) coordinateBasic(input string, domains []string) error {
	// Fallback para coordenação básica
	fmt.Println("🔗 Coordenação básica entre agentes...")
	
	for _, domain := range domains {
		if agent, exists := o.agents[domain]; exists {
			fmt.Printf("▶️  Executando: %s\n", domain)
			result, err := agent.Execute(input)
			if err != nil {
				fmt.Printf("❌ Erro em %s: %v\n", domain, err)
				continue
			}
			fmt.Printf("✅ %s concluído\n", domain)
			fmt.Println(result)
			fmt.Println("---")
		}
	}
	
	return nil
}

func (o *SmartOrchestrator) processBasic(input string) error {
	// Fallback para modo básico
	orchestrator := New(o.workingDir)
	return orchestrator.Process(input)
}

func (o *SmartOrchestrator) ShowInsights() {
	fmt.Println(o.learning.GetInsights())
}
