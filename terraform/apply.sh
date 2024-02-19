#!/usr/bin/env bash

echo "Executing terraform/apply.sh"
echo "Execuring tflocal init"
tflocal init
echo "Executing tflocal apply -auto-approve"
tflocal apply -auto-approve
