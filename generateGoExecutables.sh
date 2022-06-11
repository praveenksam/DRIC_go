#Windows files
##Generate Windows executables
echo "\nBuilding Windows executables for x64 architecture..."
GOOS=windows GOARCH=amd64 go build -o bin/windows/x64/serviceListDownload.exe main.go
echo "Building Windows executables for arm64 architecture..."
GOOS=windows GOARCH=arm64 go build -o bin/windows/arm64/serviceListDownload.exe main.go
##Windows copy files
echo "Copying static files (README, config and logo) to relevant folders for Windows binaries..."
cp ./README.md ./bin/windows/x64
cp ./README.md ./bin/windows/arm64
cp ./config.json ./bin/windows/x64
cp ./config.json ./bin/windows/arm64
cp "./DRIC Full Logo.png" ./bin/windows/x64
cp "./DRIC Full Logo.png" ./bin/windows/arm64
echo "Completed Windows x64 and ARM build!"


#Mac/Darwin files
##Generate Mac/Darwin executables
echo "\nBuilding Mac executables for x64 architecture..."
GOOS=darwin GOARCH=amd64 go build -o bin/mac/x64/serviceListDownload main.go
echo "Building Mac executables for arm64 architecture..."
GOOS=darwin GOARCH=arm64 go build -o bin/mac/M1/serviceListDownload main.go
##Mac/Darwin copy files
echo "Copying static files (README, config and logo) to relevant folders for Mac binaries..."
cp ./README.md ./bin/mac/x64
cp ./README.md ./bin/mac/M1
cp ./config.json ./bin/mac/x64
cp ./config.json ./bin/mac/M1
cp "./DRIC Full Logo.png" ./bin/mac/x64
cp "./DRIC Full Logo.png" ./bin/mac/M1
echo "Completed Mac x64 and ARM build!"


#Linux files
##Generate Linux executables
echo "\nBuilding Linux executables for x64 architecture..."
GOOS=linux GOARCH=amd64 go build -o bin/linux/x64/serviceListDownload main.go
echo "Building Linux executables for arm64 architecture..."
GOOS=linux GOARCH=arm64 go build -o bin/linux/arm64/serviceListDownload main.go
##Linux copy files
echo "Copying static files (README, config and logo) to relevant folders for Linux binaries..."
cp ./README.md ./bin/linux/x64
cp ./README.md ./bin/linux/arm64
cp ./config.json ./bin/linux/x64
cp ./config.json ./bin/linux/arm64
cp "./DRIC Full Logo.png" ./bin/linux/x64
cp "./DRIC Full Logo.png" ./bin/linux/arm64
echo "Completed Linux x64 and ARM build!"