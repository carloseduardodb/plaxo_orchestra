package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"plaxo-orchestra/internal/orchestrator"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: plaxo <comando> [argumentos]")
		fmt.Println("Comandos:")
		fmt.Println("  chat \"<mensagem>\"    - Executa comando único inteligente")
		fmt.Println("  interactive          - Modo interativo com IA avançada")
		fmt.Println("  insights             - Insights avançados do sistema")
		fmt.Println("  metrics              - Métricas de performance")
		fmt.Println("  spec                 - Gera especificação do projeto")
		fmt.Println("  watch                - Monitora mudanças no projeto")
		os.Exit(1)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Erro obtendo diretório atual: %v\n", err)
		os.Exit(1)
	}

	// Usa o orquestrador aprimorado com IA
	enhancedOrch := orchestrator.NewEnhancedOrchestrator(workingDir)
	ctx := context.Background()

	switch os.Args[1] {
	case "chat":
		if len(os.Args) < 3 {
			fmt.Println("Uso: plaxo chat \"<mensagem>\"")
			os.Exit(1)
		}
		
		message := strings.Join(os.Args[2:], " ")
		
		// Timeout context for requests
		ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
		
		if err := enhancedOrch.ProcessWithIntelligence(ctx, message); err != nil {
			fmt.Printf("Erro: %v\n", err)
			os.Exit(1)
		}

	case "interactive":
		runEnhancedInteractive(enhancedOrch)

	case "insights":
		showAdvancedInsights(enhancedOrch)

	case "metrics":
		showMetrics(enhancedOrch)

	case "spec":
		spec := orchestrator.NewSpec(workingDir)
		if err := spec.Generate(); err != nil {
			fmt.Printf("Erro gerando especificação: %v\n", err)
			os.Exit(1)
		}

	case "watch":
		watcher := orchestrator.NewWatcher(workingDir)
		if err := watcher.Start(); err != nil {
			fmt.Printf("Erro iniciando watcher: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Comando desconhecido: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runEnhancedInteractive(orch *orchestrator.EnhancedOrchestrator) {
	fmt.Println("🧠 Plaxo Orchestra v2.0 - Modo Interativo com Streaming")
	fmt.Println("Comandos especiais:")
	fmt.Println("  'quit' - sair")
	fmt.Println("  'insights' - estatísticas de aprendizado")
	fmt.Println("  'metrics' - métricas de performance")
	fmt.Println("  'cache clear' - limpar cache")
	fmt.Println("  'stream on/off' - ativar/desativar streaming")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	streamingEnabled := true
	
	for {
		if streamingEnabled {
			fmt.Print("plaxo🧠📡> ")
		} else {
			fmt.Print("plaxo🧠> ")
		}
		
		if !scanner.Scan() {
			break
		}
		
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		
		switch input {
		case "quit":
			fmt.Println("👋 Até logo!")
			return
		case "insights":
			showAdvancedInsights(orch)
			continue
		case "metrics":
			showMetrics(orch)
			continue
		case "cache clear":
			fmt.Println("🗑️  Cache limpo")
			continue
		case "stream on":
			streamingEnabled = true
			fmt.Println("📡 Streaming ativado - você verá o progresso em tempo real")
			continue
		case "stream off":
			streamingEnabled = false
			fmt.Println("📴 Streaming desativado - aguarde resposta completa")
			continue
		}
		
		// Process with timeout and streaming feedback
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		
		start := time.Now()
		
		if streamingEnabled {
			fmt.Println("🚀 Processando com streaming ativo...")
		}
		
		if err := orch.ProcessWithIntelligence(ctx, input); err != nil {
			fmt.Printf("❌ Erro: %v\n", err)
		}
		
		duration := time.Since(start)
		fmt.Printf("⏱️  Processado em %v\n", duration)
		
		cancel()
		fmt.Println()
	}
}

func showAdvancedInsights(orch *orchestrator.EnhancedOrchestrator) {
	fmt.Println("📊 Insights Avançados do Sistema:")
	fmt.Println(strings.Repeat("=", 50))
	
	insights := orch.GetAdvancedInsights()
	
	// Learning insights
	if totalDecisions, ok := insights["total_decisions"].(int); ok && totalDecisions > 0 {
		fmt.Printf("🎯 Decisões Totais: %d\n", totalDecisions)
		
		if successRate, ok := insights["success_rate"].(float64); ok {
			fmt.Printf("📈 Taxa de Sucesso: %.1f%%\n", successRate*100)
		}
		
		if avgRating, ok := insights["average_rating"].(float64); ok {
			fmt.Printf("⭐ Avaliação Média: %.1f/5.0\n", avgRating)
		}
		
		if maturity, ok := insights["learning_maturity"].(string); ok {
			fmt.Printf("🧠 Maturidade do Sistema: %s\n", maturity)
		}
	}
	
	// Performance metrics
	if cacheHitRate, ok := insights["cache_hit_rate"].(float64); ok {
		fmt.Printf("🚀 Taxa de Cache Hit: %.1f%%\n", cacheHitRate*100)
	}
	
	// Temporal patterns
	if patterns, ok := insights["temporal_patterns"].(int); ok {
		fmt.Printf("📅 Padrões Temporais: %d identificados\n", patterns)
	}
	
	fmt.Println()
}

func showMetrics(orch *orchestrator.EnhancedOrchestrator) {
	fmt.Println("📈 Métricas de Performance:")
	fmt.Println(strings.Repeat("=", 40))
	
	insights := orch.GetAdvancedInsights()
	
	// Show counters
	if counters, ok := insights["metrics_counters"].(map[string]int64); ok {
		fmt.Println("🔢 Contadores:")
		for name, value := range counters {
			fmt.Printf("  %s: %d\n", name, value)
		}
	}
	
	// Show percentiles
	if percentiles, ok := insights["metrics_percentiles"].(map[string]map[string]float64); ok {
		fmt.Println("\n⏱️  Latências (segundos):")
		for operation, stats := range percentiles {
			fmt.Printf("  %s:\n", operation)
			for metric, value := range stats {
				fmt.Printf("    %s: %.3f\n", metric, value)
			}
		}
	}
	
	fmt.Println()
}
