tentool_bin := ./gobin/tentool

gobin:
	mkdir gobin

.PHONY: tentool
tentool:
	go build -o $(tentool_bin) "./"

.PHONY: init
init: tentool
	$(tentool_bin) stats init
	$(tentool_bin) logs init

.PHONY: yadisk
yadisk: tentool
	$(tentool_bin) stats yadisk
	$(tentool_bin) logs yadisk

.PHONY: download
download: yadisk
	$(tentool_bin) stats download
	$(tentool_bin) logs update
	$(tentool_bin) logs download
