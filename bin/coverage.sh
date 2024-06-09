#!/bin/bash

bin/run_tests.sh

go tool cover -html=coverage.txt
