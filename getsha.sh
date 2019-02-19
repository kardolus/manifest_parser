#!/usr/bin/env bash
# Usage: ./getsha <url_archive>

wget $1 -O output.tar.gz > /dev/null 2>&1 && shasum -a 256 output.tar.gz | cut -d ' ' -f 1 && rm output.tar.gz
