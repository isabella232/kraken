Jira release tool that works with Component Versions Jira Add-on.

Jira API: https://docs.atlassian.com/jira/REST/latest/

Component Versions API:  http://componentversions.denizoguz.com/

Build
-----

     make

Run
---

     Usage of ./jirarelease-darwin-amd64:
       -component-name="": JIRA project component name.  For example, rest-server.  Required.
       -jira-base-url="http://localhost:8080": JIRA base REST URL.  Required.
       -jira-password="": JIRA admin password.  Required.
       -jira-username="": JIRA admin user.  Required.
       -next-version-name="": JIRA next version name. For example, 1.2.  Optional.
       -project-key="": JIRA project key.  For example, PLAT.  Required.
       -release-version-name="": JIRA release version name. For example, 1.1.  Required.
       -version=false: Print version and exit.

A mapping is defined as an entry returned by Component Versions
get-mappings that bears a give project-id, component-id, and
version-id.

The following invocation will get or create a mapping for version
2.1 of the component component-9 marked as released with today's
date, and get or create a mapping for the same component and mark
it as unreleased version 2.2.

     $ ./jirarelease-darwin-amd64 \
	-jira-base-url http://localhost:8080 \
	-jira-username admin \
	-jira-password admin123 \
	-project-key BP \
	-release-version-name 2.1 \
	-component-name component-9 \
	-next-version-name 2.2

Idempotency
-----------

The tool can be run more than once with the same arguments, though
that would be an odd use case.  

In that case the mappings will exist for both the release and next
versions of the component.  The release version will not be mutated
in any way.  The next version will again be marked as unreleased.
