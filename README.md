Jira release tool that works with Component Versions Jira Add-on.

Build
=====

     make

Run
===

     Usage of ./jirarelease-darwin-amd64:
       -component-name="": JIRA project component name.  For example, rest-server.  Required.
       -jira-base-url="http://localhost:8080": JIRA base REST URL.  Required.
       -jira-password="": JIRA admin password.  Required.
       -jira-username="": JIRA admin user.  Required.
       -next-version-name="": JIRA next version name. For example, 1.2.  Optional.
       -project-key="": JIRA project key.  For example, PLAT.  Required.
       -release-version-name="": JIRA release version name. For example, 1.1.  Required.
       -version=false: Print version and exit.

     ./jirarelease-darwin-amd64 -jira-base-url http://localhost:8080 -jira-username admin -jira-password admin123 -project-key BP -release-version-name 2.1 -component-name component-9 -next-version-name 2.2
