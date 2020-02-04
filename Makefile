MOD=github.com/mana-sys/adhesive
PREFIX=/usr/local
.PHONY: clean cli test

cli:
	go build -o bin/adhesive $(MOD)/cmd

test:
	go test $(MOD)/...

install: cli
	cp bin/adhesive $(DESTDIR)$(PREFIX)/bin/
	chmod +x $(DESTDIR)$(PREFIX)/bin/adhesive

uninstall:
	rm $(DESTDIR)$(PREFIX)/bin/adhesive

clean:
	rm -rf bin/
