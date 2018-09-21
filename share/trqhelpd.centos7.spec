# build with the following command:
# rpmbuild -bb
%define debug_package %{nil}

Name:       torque-helper
Version:    0.3
Release:    1%{?dist}
Summary:    A helper server for Torque/Moab
License:    FIXME
URL: https://github.com/Donders-Institute/%{name}
Source0: https://github.com/Donders-Institute/%{name}/archive/%{version}.tar.gz

BuildArch: x86_64

%description
A helper server for retrieving torque/moab job information with leveraged privilege.

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
install -m 644 share/trqhelpd.service %{buildroot}/usr/lib/systemd/system/trqhelpd.service
install -m 644 share/trqhelpd.env %{buildroot}/etc/sysconfig/trqhelpd

%files
%{_sbindir}/trqhelpd
%{_bindir}/cluster-qstat
/usr/lib/systemd/system/trqhelpd.service
/etc/sysconfig/trqhelpd

%changelog
* Fri Sep 21 2018 Hong Lee <h.lee@donders.ru.nl> - 0.3-1
- added cluster-qstat, a demo for client CLI program
- improved the client-server protocol to handle multiple commands
* Thu Sep 20 2018 Hong Lee <h.lee@donders.ru.nl> - 0.2-1
- introduced environment file in /etc/sysconfig
- added more commands to the service
* Wed Sep 19 2018 Hong Lee <h.lee@donders.ru.nl> - 0.1-1
- implemented the first interface for "checkjob --xml".
