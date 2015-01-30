Kraken:  A Jira release tool that works with Component Versions Jira Add-on.

Jira API: https://docs.atlassian.com/jira/REST/latest/

Component Versions API:  http://componentversions.denizoguz.com/

Build
-----

     make

Run
---

     Usage of ./kraken-darwin-amd64:
       -component-name="": JIRA project component name.  For example, rest-server.  Required.
       -jira-base-url="http://localhost:8080": JIRA base REST URL.  Required.
       -jira-password="": JIRA admin password.  Required.
       -jira-username="": JIRA admin user.  Required.
       -next-version-name="": JIRA next version name. For example, 1.2.  Optional.
       -project-key="": JIRA project key.  For example, PLAT.  Required.
       -release-version-name="": JIRA release version name. For example, 1.1.  Required.
       -version=false: Print version and exit.

A mapping is defined as an entry returned by Component Versions
get-mappings that bears a given project-id, component-id, and
version-id.

The following invocation will get or create a mapping for version
2.1 of the component component-9 marked as released with today's
date, and get or create a mapping for the same component and mark
it as unreleased version 2.2.

     $ ./kraken-darwin-amd64 \
	-jira-base-url http://localhost:8080 \
	-jira-username admin \
	-jira-password admin123 \
	-project-key BP \
	-release-version-name 2.1 \
	-component-name component-9 \
	-next-version-name 2.2

Idempotency
-----------

kraken is idempotent for two successive runs with the same arguments.

Consider a kraken run from yesterday.  If kraken is run today
with the same arguments, the release mapping will exist and will
be marked as released with yesterday's date.  Because the mapping
is already marked as released, kraken will not attempt to
update the mapping's release date with today's date and will therefore
produce the same result as yesterday's run for the release mapping.
As concerns the next-mapping, kraken returns the mapping if
it exists and creates if it does not exist.
