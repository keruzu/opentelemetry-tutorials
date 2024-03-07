Name: otelcol-snmptrap
Version: 1.0.0
Release: 1
License: GPL
Summary: OTel Collector based around the use case of SNMP traps
BuildRequires: systemd
%description
OpenTelemetry Collector with modules built to support receiving and processing
SNMP traps.

%build
cd ${RPM_SOURCE_DIR}
make build


%install
cd ${RPM_SOURCE_DIR}
mkdir -p %{buildroot}%{_sysconfdir}/systemd/system
install -m 750 %{name}.service %{buildroot}%{_sysconfdir}/systemd/system


%files
%defattr (-,root,root)
%config /containers/%{name}/%{name}.conf
%{_sysconfdir}/systemd/system/%{name}.service

