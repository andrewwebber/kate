go install github.com/andrewwebber/kate
go install github.com/andrewwebber/kate/contrib/analyze-local-images
cp ./bin/kate ./src/github.com/andrewwebber/kate/container
cp ./bin/analyze-local-images ./src/github.com/andrewwebber/kate/container
docker build -t andrewwebber/kate:next ./src/github.com/andrewwebber/kate/container
docker push andrewwebber/kate:next
rm ./src/github.com/andrewwebber/kate/container/kate
rm ./src/github.com/andrewwebber/kate/container/analyze-local-images
