#!/bin/bash

kubectl -n order-tracker  create configmap proto-descriptor --from-file=./descriptors/proto.pb --save-config=true  -o yaml > ./deployment/proto-descriptor-configmap.yml
kubectl -n order-tracker apply -f ./deployment/proto-descriptor-configmap.yml