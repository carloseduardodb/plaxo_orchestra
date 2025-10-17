package intelligence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LearningSystem struct {
	historyFile string
	decisions   []Decision
}

type Decision struct {
	Timestamp   time.Time         `json:"timestamp"`
	Input       string            `json:"input"`
	SelectedAgent string          `json:"selected_agent"`
	Success     bool              `json:"success"`
	Feedback    string            `json:"feedback"`
	Context     map[string]string `json:"context"`
}

func NewLearningSystem(workingDir string) *LearningSystem {
	historyFile := filepath.Join(workingDir, ".plaxo", "learning_history.json")
	
	ls := &LearningSystem{
		historyFile: historyFile,
		decisions:   []Decision{},
	}
	
	ls.loadHistory()
	return ls
}

func (ls *LearningSystem) loadHistory() {
	if _, err := os.Stat(ls.historyFile); os.IsNotExist(err) {
		return
	}

	data, err := os.ReadFile(ls.historyFile)
	if err != nil {
		return
	}

	json.Unmarshal(data, &ls.decisions)
}

func (ls *LearningSystem) saveHistory() error {
	dir := filepath.Dir(ls.historyFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(ls.decisions, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ls.historyFile, data, 0644)
}

func (ls *LearningSystem) RecordDecision(input, selectedAgent string, context map[string]string) {
	decision := Decision{
		Timestamp:     time.Now(),
		Input:         input,
		SelectedAgent: selectedAgent,
		Success:       true, // Assume sucesso inicialmente
		Context:       context,
	}

	ls.decisions = append(ls.decisions, decision)
	ls.saveHistory()
}

func (ls *LearningSystem) RecordFeedback(input string, success bool, feedback string) {
	// Encontra a decisão mais recente para este input
	for i := len(ls.decisions) - 1; i >= 0; i-- {
		if ls.decisions[i].Input == input {
			ls.decisions[i].Success = success
			ls.decisions[i].Feedback = feedback
			ls.saveHistory()
			break
		}
	}
}

func (ls *LearningSystem) GetBestAgentForInput(input string, availableAgents []string) string {
	// Analisa histórico para encontrar padrões
	scores := make(map[string]float64)
	
	for _, agent := range availableAgents {
		scores[agent] = ls.calculateAgentScore(input, agent)
	}

	// Encontra o agente com maior score
	bestAgent := ""
	maxScore := 0.0
	
	for agent, score := range scores {
		if score > maxScore {
			maxScore = score
			bestAgent = agent
		}
	}

	return bestAgent
}

func (ls *LearningSystem) calculateAgentScore(input, agent string) float64 {
	score := 0.0
	totalDecisions := 0
	
	for _, decision := range ls.decisions {
		if decision.SelectedAgent == agent {
			totalDecisions++
			
			// Pontuação por similaridade de input
			similarity := ls.calculateInputSimilarity(input, decision.Input)
			
			if similarity > 0.5 {
				if decision.Success {
					score += similarity * 1.0
				} else {
					score -= similarity * 0.5 // Penaliza falhas
				}
			}
		}
	}

	// Normaliza pela quantidade de decisões
	if totalDecisions > 0 {
		score = score / float64(totalDecisions)
	}

	return score
}

func (ls *LearningSystem) calculateInputSimilarity(input1, input2 string) float64 {
	words1 := ls.extractKeywords(input1)
	words2 := ls.extractKeywords(input2)
	
	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	matches := 0
	for word1 := range words1 {
		if _, exists := words2[word1]; exists {
			matches++
		}
	}

	// Jaccard similarity
	union := len(words1) + len(words2) - matches
	if union == 0 {
		return 0.0
	}

	return float64(matches) / float64(union)
}

func (ls *LearningSystem) extractKeywords(text string) map[string]bool {
	// Palavras irrelevantes
	stopWords := map[string]bool{
		"o": true, "a": true, "os": true, "as": true,
		"um": true, "uma": true, "de": true, "do": true,
		"da": true, "em": true, "no": true, "na": true,
		"para": true, "com": true, "por": true, "que": true,
		"como": true, "quando": true, "onde": true,
	}

	words := make(map[string]bool)
	
	// Extrai palavras significativas (> 3 caracteres, não stop words)
	for _, word := range ls.tokenize(text) {
		if len(word) > 3 && !stopWords[word] {
			words[word] = true
		}
	}

	return words
}

func (ls *LearningSystem) tokenize(text string) []string {
	// Tokenização simples
	var words []string
	var current string
	
	for _, char := range text {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			current += string(char)
		} else {
			if current != "" {
				words = append(words, current)
				current = ""
			}
		}
	}
	
	if current != "" {
		words = append(words, current)
	}
	
	return words
}

func (ls *LearningSystem) GetInsights() string {
	if len(ls.decisions) == 0 {
		return "Nenhum histórico de decisões disponível."
	}

	// Estatísticas básicas
	totalDecisions := len(ls.decisions)
	successCount := 0
	agentUsage := make(map[string]int)
	
	for _, decision := range ls.decisions {
		if decision.Success {
			successCount++
		}
		agentUsage[decision.SelectedAgent]++
	}

	successRate := float64(successCount) / float64(totalDecisions) * 100

	insights := fmt.Sprintf(`
📊 Insights do Sistema de Aprendizado:

📈 Taxa de Sucesso: %.1f%% (%d/%d decisões)
📅 Período: %s até %s

🤖 Agentes mais utilizados:
`, successRate, successCount, totalDecisions,
		ls.decisions[0].Timestamp.Format("02/01/2006"),
		ls.decisions[len(ls.decisions)-1].Timestamp.Format("02/01/2006"))

	// Top 3 agentes mais usados
	type agentStat struct {
		name  string
		count int
	}
	
	var stats []agentStat
	for agent, count := range agentUsage {
		stats = append(stats, agentStat{agent, count})
	}
	
	// Ordena por uso
	for i := 0; i < len(stats)-1; i++ {
		for j := i + 1; j < len(stats); j++ {
			if stats[j].count > stats[i].count {
				stats[i], stats[j] = stats[j], stats[i]
			}
		}
	}

	for i, stat := range stats {
		if i >= 3 {
			break
		}
		percentage := float64(stat.count) / float64(totalDecisions) * 100
		insights += fmt.Sprintf("  %d. %s: %d usos (%.1f%%)\n", i+1, stat.name, stat.count, percentage)
	}

	return insights
}
