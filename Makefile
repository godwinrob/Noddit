# Noddit Deployment Makefile
# Manual deployment commands for GCP e2-micro VM

# VM Configuration - override via environment or .env.deploy
VM_NAME    ?= noddit-vm
VM_ZONE    ?= us-central1-a
VM_USER    ?= deploy
PROJECT_DIR = /home/$(VM_USER)/noddit

# SSH via gcloud
SSH = gcloud compute ssh $(VM_NAME) --zone=$(VM_ZONE)
SCP = gcloud compute scp --zone=$(VM_ZONE)

.PHONY: help deploy build up down restart logs status ps ssh sync health backup

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

deploy: sync build up health ## Full deployment: sync, build, start, health check

sync: ## Sync project files to VM
	$(SCP) --recurse \
		backend-go $(VM_USER)@$(VM_NAME):$(PROJECT_DIR)/
	$(SCP) --recurse \
		noddit-next $(VM_USER)@$(VM_NAME):$(PROJECT_DIR)/
	$(SCP) \
		docker-compose.production.yml \
		Makefile \
		$(VM_USER)@$(VM_NAME):$(PROJECT_DIR)/
	$(SSH) --command="cd $(PROJECT_DIR) && cp docker-compose.production.yml docker-compose.yml"

build: ## Build containers on VM
	$(SSH) --command="cd $(PROJECT_DIR) && docker compose build --no-cache"

up: ## Start containers on VM
	$(SSH) --command="cd $(PROJECT_DIR) && docker compose up -d"

down: ## Stop containers on VM
	$(SSH) --command="cd $(PROJECT_DIR) && docker compose down"

restart: ## Restart containers on VM
	$(SSH) --command="cd $(PROJECT_DIR) && docker compose restart"

logs: ## Tail container logs
	$(SSH) --command="cd $(PROJECT_DIR) && docker compose logs --tail=100 -f"

logs-backend: ## Tail backend logs
	$(SSH) --command="cd $(PROJECT_DIR) && docker compose logs --tail=100 -f backend"

logs-frontend: ## Tail frontend logs
	$(SSH) --command="cd $(PROJECT_DIR) && docker compose logs --tail=100 -f frontend"

logs-db: ## Tail database logs
	$(SSH) --command="cd $(PROJECT_DIR) && docker compose logs --tail=100 -f postgres"

status: ## Show container status and resource usage
	$(SSH) --command="cd $(PROJECT_DIR) && docker compose ps && echo '' && docker stats --no-stream"

ps: ## Show container status
	$(SSH) --command="cd $(PROJECT_DIR) && docker compose ps"

health: ## Run health checks against VM
	@echo "Checking backend..."
	@$(SSH) --command="curl -sf --max-time 10 http://localhost:8080/api/public/subnoddits > /dev/null && echo 'Backend: OK' || echo 'Backend: FAILED'"
	@echo "Checking frontend..."
	@$(SSH) --command="curl -sf --max-time 10 http://localhost:8081/ > /dev/null && echo 'Frontend: OK' || echo 'Frontend: FAILED'"

ssh: ## SSH into the VM
	$(SSH)

backup: ## Run database backup on VM
	$(SSH) --command="$(PROJECT_DIR)/backup-db.sh"

backup-download: ## Download latest backup to local machine
	$(SSH) --command="ls -t $(HOME)/backups/noddit_*.sql.gz | head -1" | xargs -I{} \
		$(SCP) $(VM_USER)@$(VM_NAME):{} ./

clean: ## Remove old Docker images on VM
	$(SSH) --command="docker image prune -af --filter 'until=24h'"

vm-info: ## Show VM memory and disk usage
	$(SSH) --command="echo '=== Memory ===' && free -h && echo '' && echo '=== Disk ===' && df -h /"
