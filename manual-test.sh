#!/usr/bin/env bash
green="\033[32m"
red="\033[31m"
bold="\033[1m"
reset="\033[0m"

curr_dir=$PWD

log_success() {
  echo -e "${green}$@${reset}"
}

log_err() {
  echo -e "${red}$@${reset}"
}

log_bold() {
  echo -e "${bold}$@${reset}"
}

assert_eq() {
  if [ "$2" != "$3" ]; then
    echo -e "  ${red}$1: $2 != $3${reset}"
    exit -1
  fi
}

assert_not_null() {
  if [ "$2" == "" ]; then
    echo -e "  ${red}$1: $2 is empty${reset}"
    exit -1
  fi
}

ensure_tools() {
  set +e
  for tool_check in "${@}"
  do
    # log_bold "  ..testing tool $tool_check"
    $tool_check &> /dev/null
    assert_eq "$tool_check exit code" $? 0
  done
  set -e
}

ensure_tools "curl --version" "jq --version"

# test: unlimited requests
resp1=`curl localhost:8080/pingLimited | jq -r '.message'`
resp2=`curl localhost:8080/healthLimited | jq -r '.health'`
sleep 1
resp3=`curl localhost:8080/pingLimited | jq -r '.message'`
sleep 1
resp4=`curl localhost:8080/healthLimited | jq -r '.health'`
sleep 2
resp5=`curl localhost:8080/pingLimited | jq -r '.message'`
resp6=`curl localhost:8080/healthLimited | jq -r '.error'`
echo $resp1 $resp2 $resp3 $resp4 $resp5 $resp6
assert_eq "resp1 message" "$resp1" "pong"
assert_eq "resp2 health"  "$resp2" "good"
assert_eq "resp3 message" "$resp3" "pong"
assert_eq "resp4 health"  "$resp4" "good"
assert_eq "resp5 message" "$resp5" "pong"
assert_eq "resp6 error"   "$resp6" "Rate limit breached due to rule: {5 10s}"

sleep 8
resp7=`curl localhost:8080/pingLimited | jq -r '.message'`
resp8=`curl localhost:8080/healthLimited | jq -r '.error'`
echo $resp7 $resp8
#assert_eq "resp7 message" "$resp7" "pong"
assert_eq "resp8 error"   "$resp8" "Rate limit breached due to rule: {7 20s}"

log_success "Success!"