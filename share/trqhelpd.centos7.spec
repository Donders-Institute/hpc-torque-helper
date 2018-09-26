# build with the following command:
# rpmbuild -bb
%define debug_package %{nil}

Name:       torque-helper
Version:    %{getenv:VERSION}
Release:    1%{?dist}
Summary:    A helper package for retrieving job/system information from Torque/Moab systems
License:    FIXME
URL: https://github.com/Donders-Institute/%{name}
Source0: https://github.com/Donders-Institute/%{name}/archive/%{version}.tar.gz

BuildArch: x86_64

%description
A helper server for retrieving torque/moab job information with leveraged privilege.

%package server-srv
Summary: the server component of the %{name} for the pbs_server node
%description server
The server interfacing the torque/mom systems running on the pbs_server node to deliver job/system information to the client.

%package server-mom
Summary: the server component of the %{name} for the pbs_mom node
%description server
The server interfacing the pbs_mom node to deliver job information to the client.

%package client
Summary: the client component of the %{name}
%description client
A set of client CLI tools to interact with the server for retrieving job/system information. 

%prep
%setup -q

%build
make

%install
mkdir -p %{buildroot}/%{_sbindir}
mkdir -p %{buildroot}/%{_bindir}
mkdir -p %{buildroot}/usr/lib/systemd/system
mkdir -p %{buildroot}/etc/sysconfig
install -m 755 bin/trqhelpd %{buildroot}/%{_sbindir}/trqhelpd
install -m 755 bin/cluster-qstat %{buildroot}/%{_bindir}/cluster-qstat
install -m 755 bin/cluster-config %{buildroot}/%{_bindir}/cluster-config
install -m 644 share/trqhelpd.service %{buildroot}/usr/lib/systemd/system/trqhelpd.service

%install server-mom
install -m 644 share/trqhelpd.mom.env %{buildroot}/etc/sysconfig/trqhelpd

%install server-srv
install -m 644 share/trqhelpd.srv.env %{buildroot}/etc/sysconfig/trqhelpd

%files server-srv
%{_sbindir}/trqhelpd
/usr/lib/systemd/system/trqhelpd.service
/etc/sysconfig/trqhelpd

%files server-mom
%{_sbindir}/trqhelpd
/usr/lib/systemd/system/trqhelpd.service
/etc/sysconfig/trqhelpd

%files client
%{_bindir}/cluster-qstat
%{_bindir}/cluster-config

%changelog
* Tue Sep 25 2018 Hong Lee <h.lee@donders.ru.nl> - 0.4-1
- added more commands to the server
- added max. connections to the server, default is 100 and changeable via the argument
- split server and client into two RPM packages
* Fri Sep 21 2018 Hong Lee <h.lee@donders.ru.nl> - 0.3-1
- added cluster-qstat, a demo for client CLI program
- improved the client-server protocol to handle multiple commands
* Thu Sep 20 2018 Hong Lee <h.lee@donders.ru.nl> - 0.2-1
- introduced environment file in /etc/sysconfig
- added more commands to the service
* Wed Sep 19 2018 Hong Lee <h.lee@donders.ru.nl> - 0.1-1
- implemented the first interface for "checkjob --xml".
