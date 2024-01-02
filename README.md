# Golang Photo Parser Appliction

## Running 
Create a `.env` file in root directory of this project and set the variable `PHOTOS_PATH` to the location your photos are located e.g. `/home/photos/camera1`

## Testing
To run benchmark testing for the various collection algorithms use the command `go test -bench . -count x` where x is the amount of times each benchmark is run sequentially