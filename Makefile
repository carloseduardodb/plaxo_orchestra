BINARY_NAME=orchestra
BUILD_DIR=bin

.PHONY: build clean install test update benchmark demo

build:
	@echo "🔨 Building Plaxo Orchestra v2.0..."
	@mkdir -p $(BUILD_DIR)
	@go mod tidy
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/plaxo
	@echo "✅ Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

install: build
	@echo "📦 Installing Plaxo Orchestra globally..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✅ Installed! Use 'orchestra' command anywhere"

clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean

test:
	@echo "🧪 Running tests..."
	@go test ./internal/...

benchmark:
	@echo "⚡ Running benchmarks..."
	@go test -bench=. ./internal/...

dev: build
	@echo "🚀 Running in development mode..."
	@./$(BUILD_DIR)/$(BINARY_NAME) interactive

demo: build
	@echo "🎬 Demonstração das melhorias v2.0..."
	@echo "Cache distribuído, aprendizado avançado e coordenação inteligente"
	@./$(BUILD_DIR)/$(BINARY_NAME) chat "criar sistema de e-commerce com pagamentos"

update: build
	@echo "🔄 Updating global installation..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✅ Global installation updated!"
