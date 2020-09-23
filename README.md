## Feedback Generator for the Interviews

## Development is in progress

### Build and embed static resources
statik -src=./public -include=*.jpg,*.txt,*.html,*.css,*.js

### Build server
make build

### Build client
make build-client

### Run server
./bin/feedback-generator

### Run client
#### Create Feedback
./bin/feedback-client c

#### Delete feedback
./bin/feedback-client d -d [request-id]

#### Read feedback
./bin/feedback-client r -r [request-id]

#### Generate Feedback from the created request
./bin/feedback-client g -g [request-id]
