#!/bin/bash

go test -cover
go test -coverprofile=covprofile
go tool cover -html=covprofile -o coverage.html