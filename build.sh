#!/bin/bash

GOPATH=$HOME/Projects/go rpmbuild --undefine=_disable_source_fetch -bb build/trqhelpd.centos7.spec
