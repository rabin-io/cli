#
# spec file for package step-cli
#

%undefine _missing_build_ids_terminate_build
%define   debug_package %{nil}

Name:           step-cli
Version:        0.14.6
Release:        0
Summary:        Swiss-army knife for day-to-day production identity operations
License:        Apache-2.0
Group:          Productivity/Networking/Security
URL:            https://smallstep.com/cli/
Source0:        https://github.com/smallstep/cli/releases/download/v%{version}/step_linux_%{version}_amd64.tar.gz
Source1:        https://raw.githubusercontent.com/smallstep/cli/v%{version}/autocomplete/bash_autocomplete
Source2:        https://raw.githubusercontent.com/smallstep/cli/v%{version}/autocomplete/zsh_autocomplete
                
%description
The command-line interface for all things smallstep & a swiss-army knife for
day-to-day production identity operations

step is an Open Source command-line tool for developers, operators, and
security professionals to configure, operate, and automate the smallstep
toolchain and open standard identity technologies.

%prep
%autosetup -n step_%{version}

%build
%global _missing_build_ids_terminate_build 0

%install
install -Dm755 bin/step %{buildroot}/%{_bindir}/step
install -Dm644 %{SOURCE1} "%{buildroot}%{_datadir}/bash-completion/completions/step"
install -Dm644 %{SOURCE2} "%{buildroot}%{_datadir}/zsh/site-functions/step"

%files
#% doc README.md
%{_bindir}/step
%dir %{_datadir}/bash-completion/completions
%{_datadir}/bash-completion/completions/step
%dir %{_datadir}/zsh
%dir %{_datadir}/zsh/site-functions
%{_datadir}/zsh/site-functions/step

%changelog
