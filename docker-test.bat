docker volume create --name=gjg-test-task_pgdata
docker-compose -f docker-compose.tests.yml up ^
  --build ^
  --remove-orphans ^
  --abort-on-container-exit ^
  --exit-code-from tests