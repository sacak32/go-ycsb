#!/bin/bash

for i in {1..5}; do ./bin/go-ycsb run mydb -P workloads/workloada & done
