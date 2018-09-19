#!/bin/bash

rpmbuild --undefine=_disable_source_fetch -bb share/trqhelpd.centos7.spec
