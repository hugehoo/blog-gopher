go run . &
sleep 4 && command

curl -X POST "http://localhost:8080/latest"

sleep 4 && command

# PID 가져오기
PID=$(lsof -ti tcp:8080)

# 가져온 PID로 프로세스 종료하기
if [ ! -z "$PID" ]; then
    kill -9 $PID
    echo "Process $PID killed"
else
    echo "No process found on port 8080"
fi
echo "애플리케이션이 종료되었습니다."