# uncomment these for local debugging:
#BLINK_USER_API_KEY=<your-blink-api-key>
#E2E_CTRL_ADDR=<your controller addr (e.g https://app.dev.blinkops.com)>

# helper vars for constructing blink cmds
addr=${E2E_CTRL_ADDR#https://} # blink cli expects a url without schema
flags="--api-key $BLINK_USER_API_KEY --hostname $addr --workspace 55555555-5555-5555-5555-555555555555"

# remove this line to test the script locally without docker
alias blink="docker run --rm -v $(pwd)/test/backwards_compatible:/inputs blinkops/blink-cli"
# change to "." to test the script locally without docker
inputs_dir="/inputs" #inputs_dir="."

# fails the script in case a command has failed
assert() { # $1: test case description
  if [ $? -ne 0 ]
  then
    echo
    echo "$1 failed" >&2
    exit 1
  fi
}

# tests begin here
echo
echo "create playbooks:\n"
automation_id=$(blink automations create -f $inputs_dir/cli_automation_case_1.yaml -p "test pack" $flags)
assert "create playbook"
echo $automation_id
automation_id2=$(blink automations create -f $inputs_dir/cli_automation_case_2.yaml -p "test pack" $flags)
assert "create playbook"
echo $automation_id2

echo
echo "execute playbook:\n"
execution_id=$(blink automations execute --id $automation_id $flags)
assert "execute playbook"
echo $execution_id

execution_id=$(echo $execution_id | egrep '[0-9a-f]{8}-([0-9a-f]{4}-){3}[0-9a-f]{12}' -o)
assert "egrep uuid from execution"
echo $execution_id

echo
echo "get execution:\n"
blink executions get --id $execution_id $flags
assert "get automation"

echo
echo "update automation:\n"
blink automations update -f $inputs_dir/cli_automation_case_1_update.yaml $flags
assert "update automation"

echo
echo "delete automations:\n"
blink automations delete --id $automation_id $flags
assert "delete automation"

echo
blink automations delete --id $automation_id2 $flags
assert "delete automation"

unalias blink

