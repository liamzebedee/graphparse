# Press [Enter] to recompile.

go install
~/go/bin/graphparse -api & export PID=$!
read && kill $PID
echo Waiting...
wait $PID
echo Restarting...
./run.sh
