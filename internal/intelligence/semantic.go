package intelligence

import (
	"encoding/json"
	"fmt"
	"plaxo-orchestra/internal/pool"
	"strings"
)

type SemanticAnalyzer struct {
	cache     map[string]SemanticResult
	agentPool *pool.AgentPool
}

type SemanticResult struct {
	Intent     string            `json:"intent"`
	Entities   []string          `json:"entities"`
	Domains    []string          `json:"domains"`
	Complexity string            `json:"complexity"`
	Keywords   map[string]float64 `json:"keywords"`
}

func NewSemanticAnalyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{
		cache:     make(map[string]SemanticResult),
		agentPool: pool.NewAgentPool(),
	}
}

func (s *SemanticAnalyzer) AnalyzeIntent(input string) (*SemanticResult, error) {
	// Verifica cache primeiro
	if cached, exists := s.cache[input]; exists {
		return &cached, nil
	}

	// Usa Amazon Q CLI para análise semântica
	prompt := fmt.Sprintf(`
Analise semanticamente esta requisição e retorne APENAS um JSON válido:

Requisição: "%s"

Formato esperado:
{
  "intent": "create|modify|query|debug|integrate",
  "entities": ["entidade1", "entidade2"],
  "domains": ["dominio_provavel1", "dominio_provavel2"],
  "complexity": "simple|medium|complex",
  "keywords": {"palavra1": 0.9, "palavra2": 0.7}
}

Regras:
- intent: ação principal (create, modify, query, debug, integrate)
- entities: substantivos importantes (user, product, order, etc)
- domains: domínios técnicos prováveis (user, catalog, payment, etc)
- complexity: simple (1 domínio), medium (2-3), complex (4+)
- keywords: palavras-chave com peso de relevância (0.0-1.0)
`, input)

	output, err := s.agentPool.Execute("semantic_analyzer", prompt)
	if err != nil {
		return nil, fmt.Errorf("erro na análise semântica: %v", err)
	}

	// Extrai JSON da resposta
	jsonStr := s.extractJSON(output)
	if jsonStr == "" {
		return s.fallbackAnalysis(input), nil
	}

	var result SemanticResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return s.fallbackAnalysis(input), nil
	}

	// Cache o resultado
	s.cache[input] = result
	return &result, nil
}

func (s *SemanticAnalyzer) extractJSON(text string) string {
	// Procura por JSON na resposta
	start := strings.Index(text, "{")
	if start == -1 {
		return ""
	}

	end := strings.LastIndex(text, "}")
	if end == -1 || end <= start {
		return ""
	}

	return text[start : end+1]
}

func (s *SemanticAnalyzer) fallbackAnalysis(input string) *SemanticResult {
	input = strings.ToLower(input)
	
	// Análise básica por palavras-chave
	result := &SemanticResult{
		Intent:     "create",
		Entities:   []string{},
		Domains:    []string{},
		Complexity: "simple",
		Keywords:   make(map[string]float64),
	}

	// Detecta intent
	if strings.Contains(input, "criar") || strings.Contains(input, "novo") {
		result.Intent = "create"
	} else if strings.Contains(input, "modificar") || strings.Contains(input, "alterar") {
		result.Intent = "modify"
	} else if strings.Contains(input, "como") || strings.Contains(input, "o que") {
		result.Intent = "query"
	} else if strings.Contains(input, "erro") || strings.Contains(input, "bug") {
		result.Intent = "debug"
	} else if strings.Contains(input, "integrar") || strings.Contains(input, "conectar") {
		result.Intent = "integrate"
	}

	// Detecta entidades comuns
	entities := map[string]string{
		"usuário": "user", "usuario": "user", "user": "user",
		"produto": "product", "item": "product",
		"pedido": "order", "compra": "order",
		"pagamento": "payment", "checkout": "payment",
		"entrega": "delivery", "envio": "delivery",
	}

	for word, entity := range entities {
		if strings.Contains(input, word) {
			result.Entities = append(result.Entities, entity)
			result.Keywords[word] = 0.8
		}
	}

	// Detecta domínios
	domains := map[string]string{
		"e-commerce": "catalog", "loja": "catalog",
		"biblioteca": "book", "livro": "book",
		"delivery": "restaurant", "restaurante": "restaurant",
		"usuário": "user", "cliente": "customer",
	}

	for word, domain := range domains {
		if strings.Contains(input, word) {
			result.Domains = append(result.Domains, domain)
		}
	}

	// Determina complexidade
	if len(result.Domains) > 3 {
		result.Complexity = "complex"
	} else if len(result.Domains) > 1 {
		result.Complexity = "medium"
	}

	return result
}

func (s *SemanticAnalyzer) CalculateSimilarity(input string, agentDomain string) float64 {
	analysis, err := s.AnalyzeIntent(input)
	if err != nil {
		return s.basicSimilarity(input, agentDomain)
	}

	score := 0.0
	
	// Pontuação por domínios detectados
	for _, domain := range analysis.Domains {
		if strings.Contains(agentDomain, domain) {
			score += 0.4
		}
	}

	// Pontuação por entidades
	for _, entity := range analysis.Entities {
		if strings.Contains(agentDomain, entity) {
			score += 0.3
		}
	}

	// Pontuação por palavras-chave
	for keyword, weight := range analysis.Keywords {
		if strings.Contains(agentDomain, keyword) {
			score += weight * 0.3
		}
	}

	return score
}

func (s *SemanticAnalyzer) basicSimilarity(input, agentDomain string) float64 {
	input = strings.ToLower(input)
	agentDomain = strings.ToLower(agentDomain)
	
	words := strings.Fields(input)
	matches := 0
	
	for _, word := range words {
		if len(word) > 3 && strings.Contains(agentDomain, word) {
			matches++
		}
	}
	
	if len(words) == 0 {
		return 0.0
	}
	
	return float64(matches) / float64(len(words))
}
