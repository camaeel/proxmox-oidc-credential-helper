all: clean test build

clean:
	rm -rf bin || true
build:
	go build -o bin/proxmox-oidc-credential-helper github.com/proxmox-oidc-credential-helper/proxmox-oidc-credential-helper/cmd/proxmox-oidc-credential-helper

test:
	go test ./...

lint:
	golangci-lint run ./... -v

vet:
	go vet ./...

trivy:
	trivy fs .
