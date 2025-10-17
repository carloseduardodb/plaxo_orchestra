package intelligence

import (
	"fmt"
	"plaxo-orchestra/internal/pool"
	"strings"
)

type Coordinator struct {
	semantic  *SemanticAnalyzer
	memory    map[string]WorkflowMemory
	agentPool *pool.AgentPool
}

type WorkflowMemory struct {
	Steps     []WorkflowStep `json:"steps"`
	Context   string         `json:"context"`
	Completed []string       `json:"completed"`
}

type WorkflowStep struct {
	Agent       string            `json:"agent"`
	Action      string            `json:"action"`
	Dependencies []string         `json:"dependencies"`
	Outputs     map[string]string `json:"outputs"`
	Status      string            `json:"status"`
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		semantic:  NewSemanticAnalyzer(),
		memory:    make(map[string]WorkflowMemory),
		agentPool: pool.NewAgentPool(),
	}
}

func (c *Coordinator) PlanWorkflow(input string, availableAgents []string) (*WorkflowMemory, error) {
	analysis, err := c.semantic.AnalyzeIntent(input)
	if err != nil {
		return nil, err
	}

	// Usa Amazon Q CLI para planejar workflow
	prompt := fmt.Sprintf(`
Crie um plano de execução para esta requisição:

Requisição: "%s"
Análise semântica: %+v
Agentes disponíveis: %s

Retorne um plano estruturado indicando:
1. Quais agentes devem ser executados
2. Em que ordem (dependências)
3. Que informações cada agente precisa
4. Que outputs cada agente deve gerar

Formato:
AGENTE: nome_do_agente
AÇÃO: o que deve fazer
DEPENDE: agentes que devem executar antes (ou "nenhum")
SAÍDA: que informação deve gerar

Exemplo:
AGENTE: user
AÇÃO: criar estrutura de autenticação
DEPENDE: nenhum
SAÍDA: interfaces de autenticação

AGENTE: catalog
AÇÃO: implementar listagem de produtos
DEPENDE: user
SAÍDA: API de produtos com autenticação
`, input, analysis, strings.Join(availableAgents, ", "))

	output, err := c.agentPool.Execute("workflow_planner", prompt)
	if err != nil {
		return c.createSimpleWorkflow(input, availableAgents), nil
	}

	workflow := c.parseWorkflowPlan(string(output), availableAgents)
	c.memory[input] = *workflow
	
	return workflow, nil
}

func (c *Coordinator) parseWorkflowPlan(plan string, availableAgents []string) *WorkflowMemory {
	workflow := &WorkflowMemory{
		Steps:     []WorkflowStep{},
		Context:   plan,
		Completed: []string{},
	}

	lines := strings.Split(plan, "\n")
	var currentStep *WorkflowStep

	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		if strings.HasPrefix(line, "AGENTE:") {
			if currentStep != nil {
				workflow.Steps = append(workflow.Steps, *currentStep)
			}
			agent := strings.TrimSpace(strings.TrimPrefix(line, "AGENTE:"))
			if c.isValidAgent(agent, availableAgents) {
				currentStep = &WorkflowStep{
					Agent:        agent,
					Dependencies: []string{},
					Outputs:      make(map[string]string),
					Status:       "pending",
				}
			}
		} else if currentStep != nil {
			if strings.HasPrefix(line, "AÇÃO:") {
				currentStep.Action = strings.TrimSpace(strings.TrimPrefix(line, "AÇÃO:"))
			} else if strings.HasPrefix(line, "DEPENDE:") {
				deps := strings.TrimSpace(strings.TrimPrefix(line, "DEPENDE:"))
				if deps != "nenhum" && deps != "" {
					currentStep.Dependencies = strings.Split(deps, ",")
					for i, dep := range currentStep.Dependencies {
						currentStep.Dependencies[i] = strings.TrimSpace(dep)
					}
				}
			} else if strings.HasPrefix(line, "SAÍDA:") {
				output := strings.TrimSpace(strings.TrimPrefix(line, "SAÍDA:"))
				currentStep.Outputs["main"] = output
			}
		}
	}

	if currentStep != nil {
		workflow.Steps = append(workflow.Steps, *currentStep)
	}

	return workflow
}

func (c *Coordinator) isValidAgent(agent string, availableAgents []string) bool {
	for _, available := range availableAgents {
		if strings.Contains(available, agent) || strings.Contains(agent, available) {
			return true
		}
	}
	return false
}

func (c *Coordinator) createSimpleWorkflow(input string, availableAgents []string) *WorkflowMemory {
	// Fallback: cria workflow simples baseado em similaridade
	workflow := &WorkflowMemory{
		Steps:     []WorkflowStep{},
		Context:   input,
		Completed: []string{},
	}

	// Ordena agentes por relevância
	scores := make(map[string]float64)
	for _, agent := range availableAgents {
		scores[agent] = c.semantic.CalculateSimilarity(input, agent)
	}

	// Adiciona agentes com score > 0.3
	for agent, score := range scores {
		if score > 0.3 {
			step := WorkflowStep{
				Agent:        agent,
				Action:       fmt.Sprintf("Implementar funcionalidades relacionadas a: %s", input),
				Dependencies: []string{},
				Outputs:      map[string]string{"main": "Implementação completa"},
				Status:       "pending",
			}
			workflow.Steps = append(workflow.Steps, step)
		}
	}

	return workflow
}

func (c *Coordinator) ExecuteWorkflow(workflow *WorkflowMemory, agentExecutor func(string, string) (string, error)) error {
	fmt.Printf("🎯 Executando workflow com %d etapas\n", len(workflow.Steps))

	for len(workflow.Completed) < len(workflow.Steps) {
		executed := false
		
		for i, step := range workflow.Steps {
			if step.Status == "completed" {
				continue
			}

			// Verifica se dependências foram completadas
			if c.dependenciesCompleted(step.Dependencies, workflow.Completed) {
				fmt.Printf("▶️  Executando: %s\n", step.Agent)
				
				// Prepara contexto com outputs das dependências
				context := c.buildContextForStep(step, workflow)
				prompt := fmt.Sprintf("%s\n\nContexto das etapas anteriores:\n%s", step.Action, context)
				
				result, err := agentExecutor(step.Agent, prompt)
				if err != nil {
					fmt.Printf("❌ Erro em %s: %v\n", step.Agent, err)
					workflow.Steps[i].Status = "failed"
					continue
				}

				workflow.Steps[i].Status = "completed"
				workflow.Steps[i].Outputs["result"] = result
				workflow.Completed = append(workflow.Completed, step.Agent)
				
				fmt.Printf("✅ %s concluído\n", step.Agent)
				executed = true
			}
		}

		if !executed {
			// Deadlock ou erro - força execução dos pendentes
			for i, step := range workflow.Steps {
				if step.Status == "pending" {
					fmt.Printf("⚠️  Forçando execução: %s\n", step.Agent)
					result, _ := agentExecutor(step.Agent, step.Action)
					workflow.Steps[i].Status = "completed"
					workflow.Steps[i].Outputs["result"] = result
					workflow.Completed = append(workflow.Completed, step.Agent)
				}
			}
			break
		}
	}

	fmt.Println("🎉 Workflow concluído!")
	return nil
}

func (c *Coordinator) dependenciesCompleted(dependencies []string, completed []string) bool {
	for _, dep := range dependencies {
		found := false
		for _, comp := range completed {
			if strings.Contains(comp, dep) || strings.Contains(dep, comp) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (c *Coordinator) buildContextForStep(step WorkflowStep, workflow *WorkflowMemory) string {
	var context strings.Builder
	
	for _, dep := range step.Dependencies {
		for _, completedStep := range workflow.Steps {
			if completedStep.Status == "completed" && 
			   (strings.Contains(completedStep.Agent, dep) || strings.Contains(dep, completedStep.Agent)) {
				context.WriteString(fmt.Sprintf("\n=== Output de %s ===\n", completedStep.Agent))
				if result, exists := completedStep.Outputs["result"]; exists {
					context.WriteString(result)
				}
				context.WriteString("\n")
			}
		}
	}
	
	return context.String()
}
