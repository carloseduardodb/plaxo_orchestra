package analyzer

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type BoundedContext struct {
	Domain      string `json:"domain"`
	Context     string `json:"context"`
	Description string `json:"description"`
}

type ProjectAnalysis struct {
	BoundedContexts []BoundedContext `json:"bounded_contexts"`
}

func AnalyzeProjectRequirements(input string) ([]BoundedContext, error) {
	prompt := fmt.Sprintf(`
Analise esta requisição de software e identifique os bounded contexts necessários:
"%s"

Retorne APENAS um JSON válido no formato:
{
  "bounded_contexts": [
    {"domain": "user", "context": "authentication", "description": "Login e autenticação"},
    {"domain": "user", "context": "profile", "description": "Perfis de usuário"}
  ]
}

Regras:
- Crie bounded contexts específicos e granulares
- Evite contextos muito amplos
- Use nomes em inglês
- Máximo 3-4 contextos por domínio
- Seja específico para o que foi pedido
`, input)

	cmd := exec.Command("q", "chat", "--message", prompt)
	output, err := cmd.Output()
	if err != nil {
		return getDefaultContexts(input), nil // Fallback
	}

	// Extrai JSON da resposta
	response := string(output)
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}") + 1
	
	if jsonStart == -1 || jsonEnd <= jsonStart {
		return getDefaultContexts(input), nil
	}
	
	jsonStr := response[jsonStart:jsonEnd]
	
	var analysis ProjectAnalysis
	if err := json.Unmarshal([]byte(jsonStr), &analysis); err != nil {
		return getDefaultContexts(input), nil
	}
	
	return analysis.BoundedContexts, nil
}

func getDefaultContexts(input string) []BoundedContext {
	input = strings.ToLower(input)
	var contexts []BoundedContext
	
	// Detecta padrões comuns e cria contextos mínimos
	if containsAny(input, []string{"usuário", "user", "login", "auth"}) {
		contexts = append(contexts, BoundedContext{
			Domain: "user", Context: "management", 
			Description: "Gestão de usuários e autenticação",
		})
	}
	
	if containsAny(input, []string{"produto", "product", "item", "catálogo"}) {
		contexts = append(contexts, BoundedContext{
			Domain: "product", Context: "catalog", 
			Description: "Catálogo de produtos",
		})
	}
	
	if containsAny(input, []string{"pedido", "order", "compra", "venda"}) {
		contexts = append(contexts, BoundedContext{
			Domain: "order", Context: "processing", 
			Description: "Processamento de pedidos",
		})
	}
	
	if containsAny(input, []string{"pagamento", "payment", "cobrança"}) {
		contexts = append(contexts, BoundedContext{
			Domain: "payment", Context: "processing", 
			Description: "Processamento de pagamentos",
		})
	}
	
	// Se não detectou nada, cria um contexto genérico
	if len(contexts) == 0 {
		contexts = append(contexts, BoundedContext{
			Domain: "core", Context: "business", 
			Description: "Lógica de negócio principal",
		})
	}
	
	return contexts
}

func containsAny(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}
