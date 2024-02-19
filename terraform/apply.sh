#!/usr/bin/env bash

echo "Executing terraform/apply.sh"
echo "Executing tflocal init"
tflocal init
echo "Executing tflocal apply -auto-approve"
tflocal apply -auto-approve
