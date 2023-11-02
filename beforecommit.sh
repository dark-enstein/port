#!/usr/bin/env bash
echo init
#PURE FUNCTIONS DEFINITION
function intro() {
    echo "
    ██     ██  ██████  ██████
    ██     ██ ██    ██ ██   ██
    ██  █  ██ ██    ██ ██████
    ██ ███ ██ ██    ██ ██   ██
     ███ ███   ██████  ██████

    Copyright ©2023 Port Inc®. All Rights Reserved.\n"
}
function demarc() {
    printf -- '-%.0s' $(seq 100); echo ""
}

# ACTUAL LOGICC
intro
demarc
echo "Generating a secret secure .env.example from current .env for commit"
# delete env file if present
if [ -f .env.example ]; then
  cp .env.example .env.about.deleted
  : > .env.example
else
  touch .env.example
fi

#process .env into .env.example
while IFS= read -r on; do
  echo "$on" | cut -d "=" -f1 | xargs printf "%s=\"*************\"\n" >> .env.example
done < .env
rm -f .env.about.deleted
demarc