# [Simple solution to the challange proposed by Francisco Zanfranceschi on twitter](https://twitter.com/zanfranceschi/status/1501583683685425159)
# Application Flow
![application flow](/flow.png?raw=true)

# Running the apps
- clone this repo
- add a .env file to the producer and receiver directories following the .env.example file
- run the command "go run cmd/cmd.go" inside each directory
- send a post request to localhost:8000/vote with a body containing a json with paredao_id and emparedado_id(both ints)
