#!/bin/bash

mockgen -source=internal/handler/computer_management.go -destination=internal/mocks/computer_management_service.go -package=mocks
