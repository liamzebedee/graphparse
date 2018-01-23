# Press [Enter] to recompile.

go install
~/go/bin/graphparse & export PID=$!
read && kill $PID
echo Waiting...
wait $PID
echo Restarting...
./run.sh
