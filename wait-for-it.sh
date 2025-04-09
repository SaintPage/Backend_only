#!/usr/bin/env bash
# wait-for-it.sh: espera a que un host y puerto específicos estén disponibles.
# Uso: wait-for-it.sh <host:port> [-- command]

set -e

TIMEOUT=30
HOSTPORT=$1
shift

IFS=":" read -r HOST PORT <<< "$HOSTPORT"

echo "Esperando a que $HOST:$PORT esté disponible (timeout ${TIMEOUT}s)..."

for ((i=0;i<TIMEOUT;i++)); do
  if nc -z "$HOST" "$PORT"; then
    echo "$HOST:$PORT está listo."
    exec "$@"
    exit 0
  fi
  sleep 1
done

echo "Timeout: $HOST:$PORT no está disponible tras $TIMEOUT segundos."
exit 1