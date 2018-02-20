# Press [Enter] to recompile.
trap ctrl_c INT

function ctrl_c() {
    kill $PID
    exit
}

go install
~/go/bin/graphparse & export PID=$!
read && kill $PID
echo Waiting...
wait $PID
echo Restarting...
./run.sh
