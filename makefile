project_name := xpeto
root_file := ./cmd/root.go
dist_dir := dist
version := 0.1.0

.PHONY: build upload install clean

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

test:
	clear && go test -v -count=1 $(pkg)

clean:
	rm -rf $(dist_dir)
