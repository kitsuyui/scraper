#!/usr/bin/env bash
outdir=$(mktemp -d)

for pkg in '' 'scraper' 'commands' ; do
  go test \
    -covermode=atomic \
    -coverprofile="$outdir"/"$pkg".out \
    ./"$pkg" \
  > /dev/null
done
cat - - <<<'mode: atomic' <(tail -n +2 -q "$outdir"/*.out) > coverage.out
rm -rf "$outdir"
