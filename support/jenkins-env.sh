#!/bin/bash 

# If you run kraken from Jenkins, this can help parse the component name and versions.

set -vux

function cleanup {
   rm -f prog.pl
}
trap cleanup EXIT

JOB_NAME=theservice-release
CURRENT_RELEASE_VERSION=1.1
DEVELOPMENT_VERSION=1.2-SNAPSHOT

cat <<"F" > prog.pl
if ( $ARGV[0] =~ /(.*)$ARGV[1]$/ ) { print $1 } else { exit(1) }
F

COMPONENT_NAME=$(perl -- prog.pl $JOB_NAME -release)
if [ $? != 0 ]; then
	echo "JOB_NAME does not contain the suffix -release.  Cannot determine Jira component name. Exiting."
	exit 1
fi

NEXT_VERSION=$(perl -- prog.pl $DEVELOPMENT_VERSION -SNAPSHOT)
if [ $? != 0 ]; then
	echo "DEVELOPMENT_VERSION does not contain the suffix -SNAPSHOT  Cannot determine component next-version name. Exiting."
	exit 1
fi

echo ./kraken-linux-amd64 \
	-jira-base-url http://localhost:8080 \
	-jira-username admin \
	-jira-password admin123 \
	-project-key BP \
	-component-name $COMPONENT_NAME \
	-release-version-name $CURRENT_RELEASE_VERSION \
	-next-version-name $NEXT_VERSION

