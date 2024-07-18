#!/usr/bin/env bats

@test "accepted and mutated" {
  run kwctl run annotated-policy.wasm -r test_data/pod.json
  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request accepted
  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*true') -ne 0 ]

  # request mutated
  [ $(expr "$output" : '.*patch":".*') -ne 0 ]
}
