rm config.properties
rm dump.db

rm main
go build main.go
sudo chmod +s main
./main -init