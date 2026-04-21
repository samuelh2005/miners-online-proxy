#!/bin/sh
set -e

if [ -f /run/secrets/velocity_forwarding_secret ]; then
  export GATE_VELOCITY_SECRET="$(cat /run/secrets/velocity_forwarding_secret)"
fi

exec /gate