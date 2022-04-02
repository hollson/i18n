all: help

## build@编译程序。
.PHONY:build
build:
	@cd ./cmd && go build -o ../bin/goi18n
	@echo " ./bin/goi18n"
	@echo "\033[31m ☕️ 编译完成\033[0m";


## clean@清理编译、日志和缓存等数据。
.PHONY:clean
clean:
	@rm -rf ./logs;
	@rm -rf ./log;
	@rm -rf ./debug;
	@rm -rf ./tmp;
	@rm -rf ./temp;
	@rm -rf ./model;
	@rm -rf ./bin;
	@rm -rf ./_gen;
	@echo "\033[31m ✅  清理完毕\033[0m";


## install@安装应用。
.PHONY:install
install:build
	@mv ./bin/goi18n ${GOPATH}/bin
	@echo "\033[31m ✅  安装完毕\033[0m";


## commit <msg>@提交Git(格式:make commit msg=备注内容,msg为可选参数)。
.PHONY:commit
message:=$(if $(msg),$(msg),"Rebuilded at $$(date '+%Y/%m/%d %H:%M:%S')")
commit:
	@echo "\033[0;34mPush to remote...\033[0m"
	@git add .
	@git commit -m $(message)
	@echo "\033[0;31m 💿 Commit完毕\033[0m"


## push <msg>@提交并推送到Git仓库(格式:make push msg=备注内容,msg为可选参数)。
.PHONY:push
push:commit
	@git push #origin master
	@echo "\033[0;31m ⬆️ Push完毕\033[0m"


### test@执行测试。
#.PHONY:test
#test: build
#	@./bin/dbcoder -driver=pgsql -host=localhost -port=5432 -user=postgres -auth=123456 -dbname=deeplink -gorm=true -package=hello


## help@查看make帮助。
.PHONY:help
help:Makefile
	@echo "Usage:\n  make [command]"
	@echo
	@echo "Available Commands:"
	@sed -n "s/^##//p" $< | column -t -s '@' |grep --color=auto "^[[:space:]][a-z]\+[[:space:]]"
	@echo
	@echo "For more to see https://makefiletutorial.com/"
