#!/bin/bash

rpmbuild --undefine=_disable_source_fetch -bb build/rpm/centos7.spec
