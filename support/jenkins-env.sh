#!/bin/bash 

# If you run jirarelease from Jenkins, this can help parse the component name and versions.

set -vux

JOB_NAME=theservice-release
CURRENT_RELEASE_VERSION=1.1
DEVELOPMENT_VERSION=1.2-SNAPSHOT

COMPONENT_NAME=$(perl -e 'if ( $ARGV[0] =~ /(.*)$ARGV[1]$/ ) { print $1 } else { exit(1) }' -- $JOB_NAME -release)
if [ $? != 0 ]; then
	echo "JOB_NAME does not contain the suffix -release.  Cannot determine Jira component name. Exiting."
	exit 1
fi

NEXT_VERSION=$(perl -e 'if ( $ARGV[0] =~ /(.*)$ARGV[1]$/ ) { print $1 } else { exit(1) }' -- $DEVELOPMENT_VERSION -SNAPSHOT)
if [ $? != 0 ]; then
	echo "DEVELOPMENT_VERSION does not contain the suffix -SNAPSHOT  Cannot determine Jira component name. Exiting."
	exit 1
fi

echo ./jirarelease-linux-amd64 \
	-jira-base-url http://localhost:8080 \
	-jira-username admin \
	-jira-password admin123 \
	-project-key BP \
	-component-name $COMPONENT_NAME \
	-release-version-name $CURRENT_RELEASE_VERSION \
	-next-version-name $NEXT_VERSION


