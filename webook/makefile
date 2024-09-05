.PHONY:mock
mock:
	#@mockgen -source=D:\code\GOcode\basic-go\webook\internal\service\code.go -package=svcmocks -destination=D:\code\GOcode\basic-go\webook\internal\service\mocks\code.mock.go
	#@mockgen -source=D:\code\GOcode\basic-go\webook\internal\repository\user.go -package=repomocks -destination=D:\code\GOcode\basic-go\webook\internal\repository\mocks\user.mock.go
	#@mockgen -source=D:\code\GOcode\basic-go\webook\internal\repository\code.go -package=repomocks -destination=D:\code\GOcode\basic-go\webook\internal\repository\mocks\code.mock.go
	#@mockgen -source=D:\code\GOcode\basic-go\webook\internal\repository\dao\user.go -package=cachemocks -destination=D:\code\GOcode\basic-go\webook\internal\repository\dao\mocks\user.mock.go
	#@mockgen -source=D:\code\GOcode\basic-go\webook\internal\repository\cache\user.go -package=cachemocks -destination=D:\code\GOcode\basic-go\webook\internal\repository\cache\mocks\user.mock.go
	#@mockgen -package=cmdablemocks -destination=D:\code\GOcode\basic-go\webook\internal\repository\cache\redismocks\cmdableuser.mock.go   github.com/redis/go-redis/v9 Cmdable
	@mockgen  -source=D:\code\GOcode\basic-go\webook\internal\service\sms\types.go -package=smsmocks -destination=D:\code\GOcode\basic-go\webook\internal\service\sms\mocks\types.mock.go
	@go mod tidy
