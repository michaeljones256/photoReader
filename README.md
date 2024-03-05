# Golang Photo Appliction
This application will scan a local file system for photos and show data results to the user

Planned features:
- Dispaly a heat map of most active photography months, years, and days

## Running this application
Create a `.env` file in root directory of this project and set the variable `PHOTOS_PATH` to the location where the photos are located e.g. `/home/photos/camera1`
`go run .`

## Testing
To run benchmark testing for the various collection algorithms use the command `go test -bench . -count x` where x is the amount of times each benchmark is run sequentially