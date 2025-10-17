# 🕷️ Plaxo Orchestra - Agent Spread Mode

## 🎯 Conceito

O **Agent Spread** é um modo revolucionário que analisa automaticamente uma aplicação existente e **distribui agentes especializados** por toda sua estrutura, criando uma rede inteligente de assistentes focados em domínios específicos.

## 🚀 Como Funciona

### 1. **Análise Automática da Aplicação**
```bash
orchestra spread
```

O sistema escaneia:
- 📁 **Estrutura de diretórios**
- 📚 **Tech stack** (Python, JavaScript, Go, etc.)
- 🎯 **Domínios funcionais** (auth, products, api, etc.)
- 📊 **Complexidade** (simple, medium, complex)
- 📄 **Arquivos de código** por domínio

### 2. **Detecção Inteligente de Domínios**

O analisador identifica automaticamente:

| Domínio | Padrões Detectados |
|---------|-------------------|
| **auth** | auth/, authentication/, login/, users/ |
| **products** | products/, catalog/, items/ |
| **api** | api/, routes/, controllers/, handlers/ |
| **models** | models/, entities/, schemas/, database/ |
| **services** | services/, business/, logic/, core/ |
| **tests** | tests/, test/, spec/, __tests__/ |
| **config** | config/, settings/, env/ |
| **deploy** | deploy/, k8s/, docker/, infra/ |

### 3. **Distribuição Automática de Agentes**

Para cada domínio detectado, cria:
```
domain_name/
├── agents/
│   └── agent.yaml    # Configuração do agente especializado
└── [arquivos existentes]
```

### 4. **Configuração Central**
```yaml
# orchestra.yaml (raiz do projeto)
app_name: minha_app
complexity: medium
tech_stack: [Python, FastAPI, PostgreSQL]
total_domains: 4

agents:
  auth:
    - /path/to/auth/agents
  products:
    - /path/to/products/agents
  api:
    - /path/to/api/agents

orchestration:
  analyze_all: "Analisar todos os domínios"
  refactor_all: "Refatorar aplicação completa"
  test_all: "Executar todos os testes"
```

## 🎮 Uso Prático

### **Passo 1: Análise e Distribuição**
```bash
cd meu_projeto
orchestra spread

# Saída:
🔍 Analisando estrutura da aplicação...
📚 Tech Stack detectado: [Python FastAPI PostgreSQL]
📊 Complexidade: medium

🤖 Planejando distribuição de agentes...
  📍 auth: /projeto/auth/agents (5 arquivos)
  📍 products: /projeto/products/agents (8 arquivos)
  📍 api: /projeto/api/agents (12 arquivos)

❓ Deseja distribuir os agentes? (s/N): s

🚀 Distribuindo agentes pela aplicação...
✅ Agente auth criado
✅ Agente products criado  
✅ Agente api criado
🎉 Distribuição concluída!
```

### **Passo 2: Gerenciamento de Agentes**
```bash
orchestra agents

# Interface interativa:
agents> list                    # Lista todos os agentes
agents> domains                 # Mostra domínios disponíveis
agents> auth.analyze           # Analisa domínio auth
agents> products.refactor      # Refatora produtos
agents> orchestrate test_all   # Testa toda aplicação
agents> quit
```

### **Passo 3: Comandos Especializados**

Cada agente oferece comandos específicos:

```bash
# Comandos por domínio
agents> auth.analyze "verificar segurança"
agents> auth.refactor "melhorar validação"
agents> auth.test "criar testes unitários"
agents> auth.document "gerar documentação"

# Comandos de orquestração
agents> orchestrate analyze_all    # Analisa todos os domínios
agents> orchestrate refactor_all   # Refatora aplicação completa
agents> orchestrate test_all       # Executa todos os testes
agents> orchestrate deploy_all     # Prepara deploy
```

## 🏗️ Estrutura Criada

### **Antes do Agent Spread**
```
meu_projeto/
├── auth/
│   ├── models.py
│   └── routes.py
├── products/
│   ├── models.py
│   └── services.py
└── main.py
```

### **Depois do Agent Spread**
```
meu_projeto/
├── orchestra.yaml              # ← Configuração central
├── auth/
│   ├── agents/
│   │   └── agent.yaml         # ← Agente especializado
│   ├── models.py
│   └── routes.py
├── products/
│   ├── agents/
│   │   └── agent.yaml         # ← Agente especializado
│   ├── models.py
│   └── services.py
└── main.py
```

## 🤖 Configuração dos Agentes

Cada agente possui configuração especializada:

```yaml
# auth/agents/agent.yaml
name: auth_agent
domain: auth
complexity: 3
files_count: 5
tech_stack: [Python, FastAPI, SQLAlchemy]

responsibilities:
  - Análise de código do domínio auth
  - Refatoração e otimização
  - Testes e validação
  - Documentação técnica

context:
  path: /projeto/auth
  files: 5

commands:
  analyze: "Analisar código do domínio auth"
  refactor: "Refatorar código seguindo melhores práticas"
  test: "Criar/executar testes para o domínio"
  document: "Gerar documentação técnica"
```

## 🎯 Casos de Uso Ideais

### ✅ **Perfeito para Agent Spread**
- 🏢 **Aplicações empresariais** com múltiplos domínios
- 🔧 **Refatoração de código legado**
- 📚 **Projetos com documentação defasada**
- 🧪 **Aplicações sem testes adequados**
- 🚀 **Preparação para deploy/migração**

### 🎯 **Exemplos Práticos**
```bash
# E-commerce
domains: auth, products, cart, payment, orders

# SaaS Platform  
domains: auth, billing, analytics, api, admin

# Microserviços
domains: user-service, product-service, order-service

# Aplicação Monolítica
domains: models, views, controllers, services, tests
```

## 📊 Benefícios

### **1. Especialização Automática**
- Cada agente conhece profundamente seu domínio
- Contexto específico para melhores respostas
- Comandos otimizados por área

### **2. Escalabilidade**
- Adiciona agentes conforme aplicação cresce
- Distribui carga de trabalho
- Paralelização natural de tarefas

### **3. Manutenibilidade**
- Organização clara por domínios
- Facilita refatoração incremental
- Histórico de mudanças por área

### **4. Produtividade**
- Comandos específicos por contexto
- Orquestração de tarefas complexas
- Automação de workflows

## 🔧 Comandos Disponíveis

### **Análise e Distribuição**
```bash
orchestra spread              # Analisa e distribui agentes
```

### **Gerenciamento**
```bash
orchestra agents             # Interface de gerenciamento
orchestra agents list       # Lista agentes (futuro)
orchestra agents status     # Status dos agentes (futuro)
```

### **Integração**
```bash
orchestra interactive        # Modo interativo com agentes
orchestra insights          # Métricas incluindo agentes
```

## 🚀 Exemplo Completo

### **1. Aplicação FastAPI**
```bash
# Estrutura inicial
fastapi_app/
├── main.py
├── auth/
│   ├── models.py
│   └── routes.py
├── products/
│   ├── models.py
│   └── routes.py
└── requirements.txt

# Executar Agent Spread
cd fastapi_app
orchestra spread

# Resultado: 2 agentes criados (auth, products)
```

### **2. Usando os Agentes**
```bash
orchestra agents

agents> list
# Mostra: auth_agent, products_agent

agents> auth.analyze "verificar vulnerabilidades de segurança"
# Analisa especificamente o código de autenticação

agents> products.refactor "otimizar queries do banco"
# Refatora apenas o domínio de produtos

agents> orchestrate test_all
# Executa testes em todos os domínios
```

## 🎉 Resultado

O **Agent Spread** transforma uma aplicação monolítica em uma **rede inteligente de agentes especializados**, onde cada parte do sistema tem seu próprio assistente dedicado, resultando em:

- 🎯 **Análises mais precisas** por domínio
- ⚡ **Refatorações mais seguras** e contextualizadas  
- 🧪 **Testes mais abrangentes** por área
- 📚 **Documentação mais detalhada** e específica
- 🚀 **Deploy mais confiável** com validação por domínio

**Antes**: "Analise toda a aplicação" (genérico, superficial)
**Agora**: "auth.analyze segurança + products.optimize performance" (específico, profundo)
