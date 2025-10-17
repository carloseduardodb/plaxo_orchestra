package learning

import (
	"fmt"
	"math"
	"sort"
	"time"
)

type DecisionFeedback struct {
	DecisionID string    `json:"decision_id"`
	Success    bool      `json:"success"`
	Timestamp  time.Time `json:"timestamp"`
	Context    string    `json:"context"`
	AgentUsed  string    `json:"agent_used"`
	UserRating int       `json:"user_rating"` // 1-5
}

type TemporalPattern struct {
	Hour        int     `json:"hour"`
	DayOfWeek   int     `json:"day_of_week"`
	Frequency   int     `json:"frequency"`
	SuccessRate float64 `json:"success_rate"`
}

type AdvancedLearning struct {
	feedback        []DecisionFeedback `json:"feedback"`
	patterns        []TemporalPattern  `json:"patterns"`
	agentPerformance map[string]*AgentMetrics `json:"agent_performance"`
}

type AgentMetrics struct {
	TotalUses     int     `json:"total_uses"`
	SuccessCount  int     `json:"success_count"`
	AvgRating     float64 `json:"avg_rating"`
	LastUsed      time.Time `json:"last_used"`
	Confidence    float64 `json:"confidence"`
}

func NewAdvancedLearning() *AdvancedLearning {
	return &AdvancedLearning{
		feedback:         make([]DecisionFeedback, 0),
		patterns:         make([]TemporalPattern, 0),
		agentPerformance: make(map[string]*AgentMetrics),
	}
}

func (al *AdvancedLearning) RecordFeedback(decisionID, context, agent string, success bool, rating int) {
	feedback := DecisionFeedback{
		DecisionID: decisionID,
		Success:    success,
		Timestamp:  time.Now(),
		Context:    context,
		AgentUsed:  agent,
		UserRating: rating,
	}
	
	al.feedback = append(al.feedback, feedback)
	al.updateAgentMetrics(agent, success, rating)
	al.analyzeTemporalPatterns()
}

func (al *AdvancedLearning) updateAgentMetrics(agent string, success bool, rating int) {
	if al.agentPerformance[agent] == nil {
		al.agentPerformance[agent] = &AgentMetrics{}
	}
	
	metrics := al.agentPerformance[agent]
	metrics.TotalUses++
	metrics.LastUsed = time.Now()
	
	if success {
		metrics.SuccessCount++
	}
	
	// Update average rating
	oldAvg := metrics.AvgRating
	metrics.AvgRating = (oldAvg*float64(metrics.TotalUses-1) + float64(rating)) / float64(metrics.TotalUses)
	
	// Calculate confidence based on usage and success rate
	successRate := float64(metrics.SuccessCount) / float64(metrics.TotalUses)
	usageWeight := math.Min(float64(metrics.TotalUses)/10.0, 1.0)
	metrics.Confidence = successRate * usageWeight
}

func (al *AdvancedLearning) analyzeTemporalPatterns() {
	patternMap := make(map[string]*TemporalPattern)
	
	for _, fb := range al.feedback {
		hour := fb.Timestamp.Hour()
		dayOfWeek := int(fb.Timestamp.Weekday())
		key := fmt.Sprintf("%d_%d", hour, dayOfWeek)
		
		if patternMap[key] == nil {
			patternMap[key] = &TemporalPattern{
				Hour:      hour,
				DayOfWeek: dayOfWeek,
			}
		}
		
		pattern := patternMap[key]
		pattern.Frequency++
		if fb.Success {
			pattern.SuccessRate = (pattern.SuccessRate*float64(pattern.Frequency-1) + 1.0) / float64(pattern.Frequency)
		} else {
			pattern.SuccessRate = (pattern.SuccessRate * float64(pattern.Frequency-1)) / float64(pattern.Frequency)
		}
	}
	
	al.patterns = make([]TemporalPattern, 0, len(patternMap))
	for _, pattern := range patternMap {
		al.patterns = append(al.patterns, *pattern)
	}
}

func (al *AdvancedLearning) GetProactiveSuggestions(context string) []string {
	now := time.Now()
	currentHour := now.Hour()
	currentDay := int(now.Weekday())
	
	suggestions := make([]string, 0)
	
	// Find patterns for current time
	for _, pattern := range al.patterns {
		if pattern.Hour == currentHour && pattern.DayOfWeek == currentDay && pattern.SuccessRate > 0.8 {
			suggestions = append(suggestions, fmt.Sprintf("Based on patterns, this is a good time for %s operations", context))
		}
	}
	
	// Suggest best performing agents
	type agentScore struct {
		name  string
		score float64
	}
	
	scores := make([]agentScore, 0)
	for agent, metrics := range al.agentPerformance {
		score := metrics.Confidence * metrics.AvgRating / 5.0
		scores = append(scores, agentScore{agent, score})
	}
	
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})
	
	if len(scores) > 0 {
		suggestions = append(suggestions, fmt.Sprintf("Consider using agent '%s' (confidence: %.2f)", scores[0].name, scores[0].score))
	}
	
	return suggestions
}

func (al *AdvancedLearning) GetInsights() map[string]interface{} {
	totalFeedback := len(al.feedback)
	successCount := 0
	totalRating := 0.0
	
	for _, fb := range al.feedback {
		if fb.Success {
			successCount++
		}
		totalRating += float64(fb.UserRating)
	}
	
	insights := map[string]interface{}{
		"total_decisions":    totalFeedback,
		"success_rate":       float64(successCount) / float64(totalFeedback),
		"average_rating":     totalRating / float64(totalFeedback),
		"agent_performance":  al.agentPerformance,
		"temporal_patterns":  len(al.patterns),
		"learning_maturity":  al.calculateMaturity(),
	}
	
	return insights
}

func (al *AdvancedLearning) calculateMaturity() string {
	if len(al.feedback) < 10 {
		return "nascent"
	} else if len(al.feedback) < 50 {
		return "developing"
	} else if len(al.feedback) < 200 {
		return "mature"
	}
	return "expert"
}
