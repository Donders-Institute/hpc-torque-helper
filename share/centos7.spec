# build with the following command:
# rpmbuild -bb
%define debug_package %{nil}

Name:       torque-helper
Version:    0.1
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
mkdir -p %{buildroot}/%{_bindir}
install -m 755 bin/%{name} %{buildroot}/%{_bindir}/%{name}
install -m 644 share/%{name}.service /run/systemd/system/%{name}.service

%files
%{_bindir}/%{name}

%changelog
* Wed Sep 19 2018 Hong Lee <h.lee@donders.ru.nl> - 0.1-1
- implemented the first interface for "checkjob --xml".
