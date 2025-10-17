#!/bin/bash

echo "🧪 Teste de Performance - Pool de Agentes"
echo "=========================================="

# Teste com múltiplas requisições sequenciais
echo "📊 Testando 5 requisições sequenciais..."

start_time=$(date +%s.%N)

echo "1. Análise de usuário..."
./bin/orchestra chat "validar dados do usuário" > /dev/null 2>&1

echo "2. Análise de produto..."
./bin/orchestra chat "criar catálogo de produtos" > /dev/null 2>&1

echo "3. Análise de pagamento..."
./bin/orchestra chat "integrar gateway de pagamento" > /dev/null 2>&1

echo "4. Análise de pedido..."
./bin/orchestra chat "processar pedidos" > /dev/null 2>&1

echo "5. Análise de relatório..."
./bin/orchestra chat "gerar relatórios de vendas" > /dev/null 2>&1

end_time=$(date +%s.%N)
duration=$(echo "$end_time - $start_time" | bc)

echo "⏱️  Tempo total: ${duration}s"
echo "📈 Média por requisição: $(echo "scale=2; $duration / 5" | bc)s"

echo ""
echo "✅ Benefícios do Pool de Agentes:"
echo "   • Reutilização de instâncias do Q CLI"
echo "   • Redução do overhead de inicialização"
echo "   • Contexto mantido entre chamadas"
echo "   • Cleanup automático de instâncias ociosas"
