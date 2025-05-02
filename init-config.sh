#!/bin/sh

envsubst < /app/config.template.yaml > /app/config.yaml

exec ./main