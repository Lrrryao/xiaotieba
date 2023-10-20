openxiaotieba:
powershell 
cd "D:\GO\xiaotieba"
#复制后`powershell`命令自动执行，再按一次enter执行 `cd "D:\GO\xiaotieba"`


postgres:
docker run --name postgres2 -p 5432:5432 -e POSTGRES_PASSWORD=sjy8873353 -d postgres

createdb:
	docker exec -it postgres2 createdb --username=postgres --owner=postgres postgres3
#最后的postgres3是新建的数据库名称，前面的参数是postgres3建库的信息
#即在名为“postgres2”容器里创建一个新数据库，名为postgres3，中间的参数用于验证信息

dropdb:
	docker exec -it postgres dropdb postgres3

up:
	docker run -v "D:/GO/xiaotieba/db/migration:/migrations" --network host migrate/migrate     -path=/migrations/ -database postgres://postgres:sjy8873353@localhost:5432/postgres3?sslmode=disable up
                                                              #使用network因为go程序、migration和postgres在不同容器/环境 
    #迁移到最高版本


	
down:
	docker run -v "D:/GO/xiaotieba/db/migration:/migrations" --network host migrate/migrate     -path=/migrations/ -database postgres://postgres:sjy8873353@localhost:5432/postgres3?sslmode=disable down
    #`这个语句一直提示有语法错误，以后在看看吧已经解决，把table名从user改成users就好了，原因未知`
    #迁移到最低版本
	

generate:
	docker run --rm --platform linux/arm64 -v "D:\GO\xiaotieba:/src" -w /src sqlc/sqlc generate

redis:
docker run --name xiaotieba -p 6379:6379 -d redis:latest

ping_redis:
docker exec -it xiaotieba redis-cli ping