#!/bin/bash
# This script will run any necessary post commands to the API, then save any necessary IDs to environment variables for use in the test scripts

# Variables
HOST="http://localhost:8080"
OUTPUT="/dev/null" # Change to /dev/stdout for debugging, or a file for logging

# Make sure the API is running, but both Dev and Engineer resources are empty
if [[ $(curl ${HOST}/dev | jq '. | length') -ne 0 ]]; then
    echo "Dev resources are not empty. Please clear the resources before running this script."
    exit 1
fi
if [[ $(curl ${HOST}/engineers | jq '. | length') -ne 0 ]]; then
    echo "Engineer resources are not empty. Please clear the resources before running this script."
    exit 1
fi

# Post Dev and Engineer resources
# creates a new Engineer resource named ryan
curl ${HOST}/engineers \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"name": "Ryan","email": "ryan@ferrets.com"}' > ${OUTPUT} 2>> ${OUTPUT}
# creates a new Engineer resource named zach
curl ${HOST}/engineers \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"name": "zach", "email": "zach@bengal.com"}' \
    >> ${OUTPUT} 2>> ${OUTPUT}
# creates a new Engineer resource named bob
curl ${HOST}/engineers \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"name": "bob", "email": "bob@bob.com"}' \
    >> ${OUTPUT} 2>> ${OUTPUT}
# creates a new Dev resource named dev_ferrets
curl ${HOST}/dev \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"name": "dev_ferrets"}' \
    >> ${OUTPUT} 2>> ${OUTPUT}
# creates a new Dev resource named dev_bengal
curl ${HOST}/dev \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"name": "dev_bengal"}' \
    >> ${OUTPUT} 2>> ${OUTPUT}

# Save ID to environment variables
export ENGINEER_ID=$(curl ${HOST}/engineers | jq -r '.[] | select(.name == "Ryan") | .id')
export DEV_ID=$(curl ${HOST}/dev | jq -r '.[] | select(.name == "dev_ferrets") | .id')

# Add Engineer to Dev resource
curl ${HOST}/dev/${DEV_ID} \
    --include \
    --header "Content-Type: application/json" \
    --request "PUT" \
    --data '{"name":"dev_ferrets", "engineers":[{"id":"'"${ENGINEER_ID}"'", "name":"Ryan", "email":"ryan@ferrets.com"}]}' \
    >> ${OUTPUT} 2>> ${OUTPUT}


# Get all Dev and Engineer resources for verification
echo "Dev Resources:"
curl ${HOST}/dev
echo "Engineer Resources:"
curl ${HOST}/engineers

# Run the tests
go install .
TF_ACC=1 go test . -v