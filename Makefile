MOD=github.com/mana-sys/adhesive

.PHONY: clean cli test

cli:
	go build -o bin/adhesive $(MOD)/cmd

test:
	go test $(MOD)/...


clean:
	rm -rf bin/
