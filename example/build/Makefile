gobuild:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making build<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	./script/build/gobuild.sh hello 999 bin/hello

version:export APP_REGION=huabei
version:export APP_ZONE=ali-3
version:
	./bin/hello --version

run:export MOCRO_MODE=dev
run:export APP_REGION=huabei
run:export APP_ZONE=ali-3
run:export APP_MODE=dev
run:
	./bin/hello --config=config.toml

governor:
	curl http://127.0.0.1:19090/routes
