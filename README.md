# Plaxo Orchestra 🎼 (Em testes)

Orquestrador inteligente de agentes IA que potencializa o Amazon Q CLI com:

- **Análise Semântica**: Entende intenções usando IA
- **Pool de Agentes**: Instâncias persistentes para performance otimizada
- **Multi-Agente**: Coordena agentes especializados automaticamente
- **Sistema de Aprendizado**: Melhora decisões com base no histórico

## 🚀 Funcionalidades Principais

### 🕷️ **Agent Spread** (NOVO!)

- **Análise Automática**: Escaneia aplicação e detecta domínios
- **Distribuição Inteligente**: Cria agentes especializados por área
- **Comandos Específicos**: Cada agente conhece seu domínio profundamente
- **Orquestração Global**: Coordena todos os agentes automaticamente

### 📡 **Streaming em Tempo Real**

- **Feedback Imediato**: Vê o progresso em tempo real
- **Zero Timeout**: Elimina percepção de travamento
- **Progress Tracking**: Barra de progresso e indicadores visuais
- **Controle Flexível**: Liga/desliga streaming no modo interativo

### ⚡ Pool de Agentes Persistentes

- **Zero Cold Start**: Reutiliza instâncias do Amazon Q CLI
- **80-90% mais rápido**: Elimina overhead de inicialização
- **Contexto Preservado**: Mantém estado entre chamadas
- **Cleanup Automático**: Remove instâncias ociosas (10min)

### 🧠 Análise Semântica Inteligente

- **Detecção de Intenções**: create, modify, query, debug, integrate
- **Extração de Entidades**: Identifica substantivos importantes
- **Análise de Domínios**: user, catalog, payment, order, etc.
- **Classificação de Complexidade**: simple, medium, complex

### 🎯 Detecção Automática de Projetos

- **Projeto Novo**: Cria estrutura multi-agente automaticamente
- **Agente Único**: Repassa diretamente para Amazon Q CLI
- **Multi-Agente**: Coordena agentes especializados por domínio

### 📊 Sistema de Observabilidade

- **Métricas de Performance**: Tempo de resposta, cache hits
- **Insights de Aprendizado**: Taxa de sucesso, padrões de uso
- **Cache Inteligente**: Otimiza análises repetidas

## 📦 Instalação

```bash
# Clone e compile
git clone <repo>
cd plaxo_orchestra
make build
make install
```

## 🎮 Uso

### Comando Único

```bash
orchestra chat "criar sistema de e-commerce completo"
```

### Modo Interativo

```bash
orchestra interactive
# plaxo🧠> criar API de usuários
# plaxo🧠> insights
# plaxo🧠> quit
```

### Agent Spread - Distribuição Automática

```bash
# Analisar aplicação e distribuir agentes
orchestra spread

# Gerenciar agentes distribuídos
orchestra agents
# agents> list                    # Lista agentes
# agents> auth.analyze           # Analisa domínio auth
# agents> products.refactor      # Refatora produtos
# agents> orchestrate test_all   # Testa tudo
# agents> quit
```

### Modo Interativo com Streaming

```bash
orchestra interactive
# plaxo🧠📡> criar API de usuários    # 📡 indica streaming ativo
# plaxo🧠📡> stream off               # Desativa streaming  
# plaxo🧠> stream on                  # Ativa streaming
# plaxo🧠📡> insights
# plaxo🧠📡> quit
```

### Comandos Disponíveis

```bash
orchestra chat "mensagem"    # Executa comando único
orchestra interactive        # Modo interativo inteligente
orchestra spread            # Analisa e distribui agentes
orchestra agents            # Gerencia agentes distribuídos
orchestra insights          # Estatísticas de aprendizado
orchestra metrics           # Métricas de performance
orchestra spec              # Gera especificação do projeto
orchestra watch             # Monitora mudanças no projeto
```

## 🏗️ Como Funciona

### 1. Detecção Automática

```
Diretório vazio → Modo Agente Único → Amazon Q CLI direto
Projeto existente → Análise de domínios → Multi-Agente
```

### 2. Análise Semântica

```
Input: "validar dados do usuário no cadastro"
↓
🧠 Análise com IA:
- Intent: modify
- Entities: [usuário, dados, cadastro]
- Domains: [user]
- Complexity: simple
↓
🎯 Agente selecionado: user/registration
```

### 3. Pool de Agentes

```
Primeira chamada: Cria instância Q CLI (~3s)
Próximas chamadas: Reutiliza instância (~0.1s)
Após 10min inativo: Remove automaticamente
```

### 4. Coordenação Multi-Agente

```
Input complexo → Planeja workflow → Executa em ordem → Compartilha contexto
```

## 📁 Estrutura de Projeto Multi-Agente

Quando detecta um projeto complexo, cria automaticamente:

```
projeto/
├── user/
│   ├── registration/agents/    # Cadastro de usuários
│   └── authentication/agents/  # Autenticação
├── catalog/
│   ├── products/agents/        # Gestão de produtos
│   └── categories/agents/      # Categorias
├── order/
│   ├── cart/agents/           # Carrinho de compras
│   └── checkout/agents/       # Finalização
└── payment/
    └── gateway/agents/        # Gateway de pagamento
```

## 🎯 Exemplos Práticos

### Análise Automática

```bash
orchestra chat "integrar pagamento com carrinho"
# 🧠 Analisando requisição com IA...
# 🎯 Intent: integrate | Complexidade: complex
# 🔗 Executando workflow inteligente...
# 📋 Workflow planejado com 3 etapas
```

### Performance Otimizada

```bash
# Primeira chamada (startup)
time orchestra chat "teste 1"  # ~3.2s

# Próximas chamadas (pool)
time orchestra chat "teste 2"  # ~0.3s
time orchestra chat "teste 3"  # ~0.3s
```

### Insights do Sistema

```bash
orchestra insights
# 📊 Insights do Sistema:
# 📈 Taxa de Sucesso: 95.2% (20/21 decisões)
# 🤖 Agentes mais utilizados:
#   1. user/registration: 8 usos (38.1%)
#   2. catalog/products: 6 usos (28.6%)
# ⚡ Performance:
#   • Tempo médio: 0.8s
#   • Cache hits: 73%
```

## 🔧 Arquitetura Técnica

### Pool de Agentes

```
┌─────────────────────────────────────────┐
│           Agent Pool Manager            │
├─────────────────────────────────────────┤
│ ┌─────────┐ ┌─────────┐ ┌─────────┐    │
│ │semantic │ │project  │ │free     │    │
│ │analyzer │ │creator  │ │agent    │    │
│ └─────────┘ └─────────┘ └─────────┘    │
├─────────────────────────────────────────┤
│ • Reutilização inteligente              │
│ • Contexto preservado                   │
│ • Cleanup automático                    │
└─────────────────────────────────────────┘
```

### Fluxo de Processamento

```
Input → Análise Semântica → Detecção de Projeto → Seleção de Agente → Execução → Aprendizado
```

## 📈 Benefícios

- **Performance**: 80-90% redução no tempo de resposta
- **Inteligência**: Seleção automática do melhor agente
- **Escalabilidade**: Suporte a projetos complexos
- **Aprendizado**: Melhora contínua das decisões
- **Simplicidade**: Zero configuração necessária

## 🧪 Teste de Performance

```bash
# Execute o benchmark
./test_performance.sh

# Resultado esperado:
# ⏱️  Tempo total: 4.2s
# 📈 Média por requisição: 0.84s
# ✅ 65% mais rápido que versão anterior
```

O Plaxo Orchestra transforma o Amazon Q CLI em um sistema inteligente e performático! 🎼✨
