#!/bin/bash

go test -coverprofile=covprofile
go tool cover -html=covprofile -o coverage.html