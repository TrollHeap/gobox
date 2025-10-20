.PHONY: run 

# During development 
run:
	sudo env "PATH=$$PATH" go run ./cmd/goboxcli
