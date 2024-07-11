#!/bin/bash

function handler () {
EVENT_DATA=$1

# Copy contents of GAMConfig to /tmp/resources so it's writable
mkdir -p /tmp/resources
mkdir -p /tmp/resources/GAMWork
cp -r /opt/GAMConfig /tmp/resources

# Set GAMCFGDIR env variable
export GAMCFGDIR="/tmp/resources/GAMConfig"

# Set the path to the GAM binary and jq binary
gamBin="/opt/gamadv-xtd3/gam"
jqBin="/opt/jq_linux"

event="$EVENT_DATA"
orgName=$(echo "$event" | $jqBin -r ".org")

# Run the select command
if $gamBin select "$orgName" save > /dev/null 2>&1; then # CSM
    # If successful, run the other command

    IFS=$'\n' 
    readarray -t commands <<< "$(echo "$event" | $jqBin -r ".commands[]")"
    unset IFS

    for cmd in "${commands[@]}"; do
        output=$($gamBin $cmd 2>&1)
        commandObject=$($jqBin -n \
            --arg command "$cmd" \
            --arg output "$output" \
            '{command: $command, output: $output}')

        commandObjects+=("$commandObject")
    done

    commandsJson=$($jqBin -s '.' <<< "${commandObjects[@]}")

    # Final JSON output for success
    json=$($jqBin -n \
            --arg org "$orgName" \
            --argjson commands "$commandsJson" \
            '{org: $org, commands: $commands}')

else

    # Final JSON output for error
    json=$($jqBin -n \
            --arg org "$orgName" \
            '{org: $org, output: "error selecting org"}')

fi

  echo "$json"

}