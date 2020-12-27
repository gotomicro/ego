gobuild:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making build<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	./scripts/build/gobuild.sh hello bin/hello

version:export EGO_REGION=huabei
version:export EGO_ZONE=ali-3
version:
	./bin/hello --version

run:export EGO_DEBUG=true
run:export EGO_LOG_ADD_APP=true
run:export EGO_REGION=huabei
run:export EGO_ZONE=ali-3
run:export EGO_MODE=dev
run:
	./bin/hello --config=config.toml

governor:
	curl http://127.0.0.1:19090/routes
