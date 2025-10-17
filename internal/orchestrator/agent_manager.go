package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"gopkg.in/yaml.v2"
)

type AgentConfig struct {
	Name             string            `yaml:"name"`
	Domain           string            `yaml:"domain"`
	Complexity       int               `yaml:"complexity"`
	FilesCount       int               `yaml:"files_count"`
	TechStack        []string          `yaml:"tech_stack"`
	Responsibilities []string          `yaml:"responsibilities"`
	Context          AgentContext      `yaml:"context"`
	Commands         map[string]string `yaml:"commands"`
}

type AgentContext struct {
	Path  string `yaml:"path"`
	Files int    `yaml:"files"`
}

type OrchestraConfig struct {
	AppName      string              `yaml:"app_name"`
	Complexity   string              `yaml:"complexity"`
	TechStack    []string            `yaml:"tech_stack"`
	TotalDomains int                 `yaml:"total_domains"`
	Agents       map[string][]string `yaml:"agents"`
	Orchestration map[string]string  `yaml:"orchestration"`
}

type AgentManager struct {
	rootPath        string
	orchestraConfig *OrchestraConfig
	agents          map[string]*AgentConfig
}

func NewAgentManager(rootPath string) *AgentManager {
	return &AgentManager{
		rootPath: rootPath,
		agents:   make(map[string]*AgentConfig),
	}
}

func (am *AgentManager) LoadConfiguration() error {
	configPath := filepath.Join(am.rootPath, "orchestra.yaml")
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("configuração do orchestra não encontrada. Execute 'plaxo spread' primeiro")
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("erro lendo configuração: %v", err)
	}
	
	am.orchestraConfig = &OrchestraConfig{}
	if err := yaml.Unmarshal(data, am.orchestraConfig); err != nil {
		return fmt.Errorf("erro parseando configuração: %v", err)
	}
	
	// Carregar configurações dos agentes
	return am.loadAgentConfigs()
}

func (am *AgentManager) loadAgentConfigs() error {
	for domain, paths := range am.orchestraConfig.Agents {
		for _, agentPath := range paths {
			configFile := filepath.Join(agentPath, "agent.yaml")
			
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				continue
			}
			
			data, err := os.ReadFile(configFile)
			if err != nil {
				continue
			}
			
			config := &AgentConfig{}
			if err := yaml.Unmarshal(data, config); err != nil {
				continue
			}
			
			am.agents[domain] = config
		}
	}
	
	return nil
}

func (am *AgentManager) ListAgents() {
	fmt.Println("🤖 Agentes Distribuídos:")
	fmt.Println(strings.Repeat("=", 50))
	
	if am.orchestraConfig == nil {
		fmt.Println("❌ Nenhuma configuração carregada")
		return
	}
	
	fmt.Printf("📱 Aplicação: %s\n", am.orchestraConfig.AppName)
	fmt.Printf("📊 Complexidade: %s\n", am.orchestraConfig.Complexity)
	fmt.Printf("🎯 Total de domínios: %d\n", am.orchestraConfig.TotalDomains)
	fmt.Println()
	
	for domain, config := range am.agents {
		fmt.Printf("🤖 %s (%s)\n", config.Name, domain)
		fmt.Printf("   📁 Caminho: %s\n", config.Context.Path)
		fmt.Printf("   📄 Arquivos: %d\n", config.Context.Files)
		fmt.Printf("   🎯 Responsabilidades: %d\n", len(config.Responsibilities))
		
		fmt.Printf("   💻 Comandos disponíveis:\n")
		for cmd, desc := range config.Commands {
			fmt.Printf("     • %s: %s\n", cmd, desc)
		}
		fmt.Println()
	}
}

func (am *AgentManager) ExecuteAgentCommand(domain, command, input string) error {
	agent, exists := am.agents[domain]
	if !exists {
		return fmt.Errorf("agente '%s' não encontrado", domain)
	}
	
	cmdDesc, exists := agent.Commands[command]
	if !exists {
		return fmt.Errorf("comando '%s' não disponível para agente '%s'", command, domain)
	}
	
	fmt.Printf("🤖 Executando: %s.%s\n", domain, command)
	fmt.Printf("📋 Descrição: %s\n", cmdDesc)
	fmt.Printf("🎯 Contexto: %s (%d arquivos)\n", agent.Context.Path, agent.Context.Files)
	fmt.Println(strings.Repeat("─", 50))
	
	// Construir prompt contextualizado
	contextualPrompt := am.buildContextualPrompt(agent, command, input)
	
	// Aqui integraria com o enhanced orchestrator para executar
	fmt.Printf("🧠 Prompt contextualizado:\n%s\n", contextualPrompt)
	
	return nil
}

func (am *AgentManager) buildContextualPrompt(agent *AgentConfig, command, input string) string {
	prompt := fmt.Sprintf(`Você é um agente especializado no domínio '%s'.

CONTEXTO DO AGENTE:
- Nome: %s
- Domínio: %s
- Caminho: %s
- Arquivos: %d
- Tech Stack: %v
- Complexidade: %d

RESPONSABILIDADES:
%s

COMANDO SOLICITADO: %s
DESCRIÇÃO: %s

ENTRADA DO USUÁRIO: %s

Por favor, execute a tarefa considerando:
1. O contexto específico do domínio %s
2. Os arquivos localizados em %s
3. As tecnologias utilizadas: %v
4. As responsabilidades do agente

Forneça uma resposta detalhada e específica para este domínio.`,
		agent.Domain,
		agent.Name, agent.Domain, agent.Context.Path, agent.Context.Files, agent.TechStack, agent.Complexity,
		strings.Join(agent.Responsibilities, "\n- "),
		command, agent.Commands[command],
		input,
		agent.Domain, agent.Context.Path, agent.TechStack)
	
	return prompt
}

func (am *AgentManager) GetAvailableCommands(domain string) ([]string, error) {
	agent, exists := am.agents[domain]
	if !exists {
		return nil, fmt.Errorf("agente '%s' não encontrado", domain)
	}
	
	var commands []string
	for cmd := range agent.Commands {
		commands = append(commands, cmd)
	}
	
	return commands, nil
}

func (am *AgentManager) GetDomains() []string {
	var domains []string
	for domain := range am.agents {
		domains = append(domains, domain)
	}
	return domains
}

func (am *AgentManager) ExecuteOrchestrationCommand(command string) error {
	if am.orchestraConfig == nil {
		return fmt.Errorf("configuração não carregada")
	}
	
	cmdDesc, exists := am.orchestraConfig.Orchestration[command]
	if !exists {
		return fmt.Errorf("comando de orquestração '%s' não encontrado", command)
	}
	
	fmt.Printf("🎼 Executando orquestração: %s\n", command)
	fmt.Printf("📋 Descrição: %s\n", cmdDesc)
	fmt.Println(strings.Repeat("─", 50))
	
	// Executar comando em todos os agentes relevantes
	switch command {
	case "analyze_all":
		return am.executeOnAllAgents("analyze", "Analisar código completo")
	case "refactor_all":
		return am.executeOnAllAgents("refactor", "Refatorar seguindo melhores práticas")
	case "test_all":
		return am.executeOnAllAgents("test", "Executar testes completos")
	case "deploy_all":
		return am.executeOnAllAgents("document", "Preparar documentação para deploy")
	default:
		return fmt.Errorf("comando de orquestração não implementado: %s", command)
	}
}

func (am *AgentManager) executeOnAllAgents(command, description string) error {
	fmt.Printf("🚀 Executando '%s' em todos os agentes...\n\n", command)
	
	for domain := range am.agents {
		fmt.Printf("🤖 Processando domínio: %s\n", domain)
		if err := am.ExecuteAgentCommand(domain, command, description); err != nil {
			fmt.Printf("⚠️  Erro no domínio %s: %v\n", domain, err)
		}
		fmt.Println()
	}
	
	return nil
}
