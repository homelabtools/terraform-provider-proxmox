#!/bin/sh
set -eu
URL="${1?Must provide URL as 1st argument}"
STATUS="${2?Must provide HTTP status code as 2nd argument}"
DELAY=1
MAX=300

for i in $(seq 1 $MAX); do
    echo "Waiting for '$URL' to return status '$STATUS', attempt #$i..."
    status="$(curl -kLso /dev/null -w '%{http_code}' "$URL")"
    if [ "$status" -eq "$STATUS" ]; then
        echo "SUCCESS: URL '$URL' returned status '$STATUS' after trying $i time(s) and $(( DELAY * i )) seconds"
        exit 0
    fi
    sleep $DELAY
done
echo "ERROR: URL '$URL' never returned status '$STATUS' after trying $i times for $(( DELAY * i )) seconds"
exit 1