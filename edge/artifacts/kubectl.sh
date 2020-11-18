#!/bin/bash
${BASH_SOURCE%/*}/kubectl --kubeconfig=${BASH_SOURCE%/*}/admin.conf $@
