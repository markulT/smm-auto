#!/bin/bash
echo "Waiting for mongo to start..."
sleep 10

echo "Setting up replica set..."
mongo --host mongo --eval 'rs.initiate({_id: "rs0", members: [{_id: 0, host: "localhost:27017"}]})'

echo "Replica set configured!"