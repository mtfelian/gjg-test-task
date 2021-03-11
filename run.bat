docker volume create --name=gjg-test-task_pgdata
docker start gjg-test-task_postgres_1
go build && .\gjg-test-task.exe --port=4000
docker stop gjg-test-task_postgres_1