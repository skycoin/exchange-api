#!/bin/sh

available_tags=""

if test -n "$CRYPTOPIA_API_SECRET"
then
    available_tags="$available_tags cryptopia_integration_test"
fi

if test -n "$C2CX_TEST_KEY"
then
    available_tags="$available_tags c2cx_integration_test"
fi

echo $available_tags
