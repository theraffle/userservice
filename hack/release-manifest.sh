#!/bin/bash

#
# Copyright 2022 The Raffle Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

if [ "$#" -eq 0 ]; then
  VERSION=latest
else
  VERSION=$1
  REGISTRY=$2
fi

BASEDIR=$(dirname "$0")
PROJECT_DIR="$BASEDIR/.."
CONFIG_DIR="$PROJECT_DIR/kubernetes-manifests"

RELEASE_MANIFEST="$CONFIG_DIR/release.yaml"

TARGETS=("$CONFIG_DIR/service.yaml" "$CONFIG_DIR/deployment.yaml")
function append_target(){
  local TARGET="$1"

  if [ "${TARGET: -5}" == ".yaml" ]; then
    cat "$TARGET" >> "$RELEASE_MANIFEST"
    echo "---" >> "$RELEASE_MANIFEST"
  else
    for f in "$TARGET"/*; do
      append_target "$f"
    done
  fi
}

rm -rf "$RELEASE_MANIFEST"

touch "$RELEASE_MANIFEST"

for target in "${TARGETS[@]}"; do
  append_target "$target"
done

sed -i '' "s/changjjjjjjjj\/raffle-user-service:latest/$REGISTRY\/raffle-user-service:$VERSION/g" "$RELEASE_MANIFEST"
