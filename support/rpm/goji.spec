Name:		  goji
Version:	0.1.2
Release:	1%{?dist}.tumblr
Summary:	Simple event consumer and config templater for service discovery with Marathon

Group:		  System Environment/Daemons
License:	  Apache 2.0
URL:		    https://github.com/byxorna/goji
Source0:	  https://github.com/byxorna/goji/archive/%{version}.tar.gz
Source1:	  goji.sysconfig
Source2:	  config.json
BuildRoot:	%(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)

BuildRequires:	golang >= 1.3.3
%if 0%{?rhel} >= 7
BuildRequires: pkgconfig(systemd)
BuildRequires: systemd-units
%endif

%description
goji is a server that registers with a Marathon instance, consumes events, and emits templated configs containing information about running tasks for a set of apps that you care about.


%prep
%setup -q


%build
mkdir -p ./_build/src/github.com/byxorna
ln -s $(pwd) ./_build/src/github.com/byxorna/goji
export GOPATH=$(pwd)/_build:$(pwd)/vendor:%{gopath}
go build


%install
rm -rf %{buildroot}
install -d -m 755 %{buildroot}%{_bindir}
install    -m 755 %{name}-%{version} %{buildroot}%{_bindir}/goji

install -d -m 755 %{buildroot}%{_sysconfdir}/sysconfig
install    -m 644 %{S:1} %{buildroot}%{_sysconfdir}/sysconfig/goji

install -d -m 755 %{buildroot}%{_sysconfdir}/goji
install    -m 644 example/haproxy.tmpl %{buildroot}%{_sysconfdir}/goji/haproxy.tmpl
install    -m 644 example/nginx.tmpl %{buildroot}%{_sysconfdir}/goji/nginx.tmpl
install    -m 644 %{S:2} %{buildroot}%{_sysconfdir}/goji/config.json

%if 0%{?rhel} < 7
install -d -m 755 %{buildroot}%{_initrddir}
install    -m 755 support/goji.init %{buildroot}%{_initrddir}/goji
%else
install -d -m 755 %{buildroot}%{_unitdir}
install    -m 755 support/goji.service %{buildroot}%{_unitdir}/goji.service
%endif

# drop the example configs in /usr/share
install -d -m 755 %{buildroot}%{_defaultdocdir}/%{name}-%{version}/example
cp -r example/* %{buildroot}%{_defaultdocdir}/%{name}-%{version}/example/


%clean
rm -rf %{buildroot}

%post
%if 0%{?rhel} >= 7
%systemd_post goji.service
%endif

%preun
%if 0%{?rhel} >= 7
%systemd_preun goji.service
%endif

%postun
%if 0%{?rhel} >= 7
%systemd_postun goji.service
%endif

%files
%defattr(-,root,root,-)
%doc
%{_bindir}/goji
%config(noreplace) %{_sysconfdir}/sysconfig/goji
%config(noreplace) %{_sysconfdir}/goji/haproxy.tmpl
%config(noreplace) %{_sysconfdir}/goji/nginx.tmpl
%config(noreplace) %{_sysconfdir}/goji/config.json
%if 0%{?rhel} < 7
%{_initrddir}/goji
%else
%{_unitdir}/goji.service
%endif
%{_defaultdocdir}/%{name}-%{version}/example/



%changelog
* Mon May 11 2015 Gabe Conradi <gabe@tumblr.com> 0.1.2-1.tumblr
- goji to 0.1.2 (gabe@tumblr.com)

* Fri Mar 13 2015 Gabe Conradi <gabe@tumblr.com> 0.1.1-3.tumblr
- add goji config sample (gabe@tumblr.com)
- patching build (gabe@tumblr.com)
- add source (gabe@tumblr.com)

* Fri Mar 13 2015 Gabe Conradi <gabe@tumblr.com> 0.1.1-2.tumblr
- new package built with tito


