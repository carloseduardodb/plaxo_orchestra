#!/bin/bash

echo "🕷️  Demonstração: Plaxo Orchestra Agent Spread"
echo "=============================================="
echo

echo "📁 Estrutura da aplicação de exemplo:"
echo "example_app/"
echo "├── main.py"
echo "├── requirements.txt"
echo "├── auth/"
echo "│   ├── models.py"
echo "│   └── routes.py"
echo "└── products/"
echo "    ├── models.py"
echo "    └── routes.py"
echo

echo "🔍 1. Analisando aplicação e distribuindo agentes..."
cd example_app
../bin/orchestra spread

echo
echo "🤖 2. Gerenciando agentes distribuídos..."
echo "   (Execute: ../bin/orchestra agents)"

echo
echo "📋 3. Comandos disponíveis após distribuição:"
echo "   • orchestra agents           - Gerenciar agentes"
echo "   • orchestra interactive      - Modo interativo"
echo "   • orchestra insights         - Ver estatísticas"

echo
echo "🎯 Exemplo de uso dos agentes:"
echo "   agents> list                 - Listar agentes"
echo "   agents> auth.analyze         - Analisar domínio auth"
echo "   agents> products.refactor    - Refatorar produtos"
echo "   agents> orchestrate test_all - Testar tudo"

echo
echo "✅ Demonstração preparada!"
echo "   Execute os comandos acima para ver o Agent Spread em ação!"
