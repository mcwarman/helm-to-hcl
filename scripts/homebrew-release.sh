#!/usr/bin/env bash

set -euo pipefail

version=${GITHUB_REF##*/}
branch=bump-helm-to-hcl-${version}

rm -rf ./dist/homebrew-tap

git clone "https://bot:${GH_PUBLIC_REPO_TOKEN}@github.com/mcwarman/homebrew-tap.git" ./dist/homebrew-tap

cd ./dist/homebrew-tap

git config user.name  "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

git checkout -b "${branch}"

cp ../helm-to-hcl.rb ./Formula/helm-to-hcl.rb

git commit -a -m "helm-to-hcl: Bump to ${version}"

git push -u origin "${branch}"

curl \
   -d "{\"title\":\"helm-to-hcl: Bump to ${version}\",\"base\":\"main\", \"head\":\"${branch}\"}"  \
   -H "Authorization: token ${GH_PUBLIC_REPO_TOKEN}" \
  https://api.github.com/repos/mcwarman/homebrew-tap/pulls
