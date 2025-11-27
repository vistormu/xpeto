project_name := xpeto
root_file := ./cmd/root.go
dist_dir := dist
version := 0.0.7

OS_ARCHES := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: build upload install clean graph

build: clean
	mkdir -p $(dist_dir)
	@for target in $(OS_ARCHES); do \
	  os=$${target%/*}; \
	  arch=$${target#*/}; \
	  case $$os in \
	    windows) ext=".exe" ;; \
	    *)       ext="" ;; \
	  esac; \
	  for tag in "" headless; do \
	    if [ "$$tag" = "headless" ]; then \
	      suffix="_headless"; \
	      tagflag="-tags=headless"; \
	    else \
	      suffix=""; \
	      tagflag=""; \
	    fi; \
	    echo "Building $(project_name) for $$os/$$arch$${suffix}"; \
	    GOOS=$$os GOARCH=$$arch go build $$tagflag -o $(dist_dir)/$(project_name)$$ext $(root_file); \
	    binname="$(project_name)_$${os}_$${arch}$${suffix}"; \
	    tarname="$(project_name)_$(version)_$${os}_$${arch}$${suffix}.tar.gz"; \
	    tar -czf $(dist_dir)/$${tarname} -C $(dist_dir) $(project_name)$$ext; \
	    mv $(dist_dir)/$(project_name)$$ext $(dist_dir)/$${binname}$$ext; \
	  done; \
	done

upload:
	git add .
	git commit -m $(commit)
	git push

release: build
	git add .
	git commit -m "release v$(version)"
	git tag -a v$(version) -m "release v$(version)"
	git push origin main
	git push origin --tags
	gh release create v$(version) $(dist_dir)/* --title "release v$(version)" --notes "release v$(version)"

install:
	@echo "Detecting host OS and architecture..."
	@uname_os=$$(uname -s | tr '[:upper:]' '[:lower:]'); \
	 uname_arch=$$(uname -m); \
	 case $$uname_arch in \
	    x86_64) arch="amd64" ;; \
	    arm64|aarch64) arch="arm64" ;; \
	    *) echo "Unsupported architecture: $$uname_arch"; exit 1 ;; \
	 esac; \
	 case $$uname_os in \
	    linux)   os="linux" ;; \
	    darwin)  os="darwin" ;; \
	    msys*|mingw*|cygwin*) os="windows" ;; \
	    *) echo "Unsupported OS: $$uname_os"; exit 1 ;; \
	 esac; \
	 suffix=""; \
	 if [ "$(HEADLESS)" = "1" ]; then suffix="_headless"; fi; \
	 bin="$(project_name)_$${os}_$${arch}$${suffix}"; \
	 echo "Installing $$bin..."; \
	 if [ ! -f "$(dist_dir)/$$bin" ]; then echo "Binary not found: $(dist_dir)/$$bin"; exit 1; fi; \
	 if [ "$$os" = "windows" ]; then \
	    install_path="/usr/bin/$(project_name).exe"; \
	    cp "$(dist_dir)/$$bin" "$$install_path"; \
	    echo "Installed to $$install_path"; \
	 else \
	    install_path="/usr/local/bin/$(project_name)"; \
	    sudo cp "$(dist_dir)/$$bin" "$$install_path"; \
	    sudo chmod +x "$$install_path"; \
	    echo "Installed to $$install_path"; \
	 fi
	

clean:
	rm -rf $(dist_dir)
