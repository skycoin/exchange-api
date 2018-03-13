#!/bin/sh

# every source file ending in .go
find ./exchange -path '.*\.go' | \

while read src_path
do
  # dirname turns ./foo/bar/baz into ./foo/bar
  dirname=$(dirname $src_path)

  # grabbing the package line from within the source file
  pkgname=$(grep '^package' < $src_path |\
            sed -e 's/package\s\(\s\)*\(.*\)/\2/g')

  echo "$pkgname:$dirname"
done | \

# peeling off the module names
sed -e 's/.*:\(.*\)/\1/g' | \

# unique-ifying the modules
sort | \
uniq
