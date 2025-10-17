package orchestrator

import (
	"fmt"
	"os"
	"os/exec"
)

func (o *Orchestrator) InitFromSpec(specFile string) error {
	if _, err := os.Stat(specFile); os.IsNotExist(err) {
		return fmt.Errorf("arquivo de especificação não encontrado: %s", specFile)
	}

	content, err := os.ReadFile(specFile)
	if err != nil {
		return err
	}

	prompt := fmt.Sprintf(`
Baseado nesta especificação, crie um projeto completo:

%s

Gere toda a estrutura de pastas, código e documentação necessária.
`, string(content))

	fmt.Println("🏗️ Gerando projeto baseado na especificação...")
	
	cmd := exec.Command("q", "chat", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Após gerar, configura agentes automaticamente
	return o.setupAgentsFromGenerated()
}

func (o *Orchestrator) setupAgentsFromGenerated() error {
	fmt.Println("\n🤖 Configurando agentes especializados...")
	
	contexts, err := o.detectGeneratedStructure()
	if err != nil {
		return err
	}

	for _, bc := range contexts {
		contextPath := fmt.Sprintf("%s/%s", bc.Domain, bc.Context)
		if err := o.setupAgentInExistingStructure(bc.Domain, bc.Context, bc.Description); err != nil {
			fmt.Printf("⚠️ Erro configurando %s: %v\n", contextPath, err)
			continue
		}
		fmt.Printf("✅ Agente configurado: %s\n", contextPath)
	}

	return nil
}