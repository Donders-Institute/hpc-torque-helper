#!/bin/bash

rpmbuild --undefine=_disable_source_fetch -bb build/trqhelpd.centos7.spec
