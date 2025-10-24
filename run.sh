rm ./main
go build main.go
docker build -t godash . --progress=plain --no-cache
docker stop godash
docker rm godash
docker run --name godash --restart=always -p 8868:80 -d godash
echo y|docker system prune