# SCIMPLISTIC
Scimplistic is a simplistic Go webapp to manage tasks, I built this tool to demonstrate the integration avaialble via SCIM for CyberArk Privileged Access Manager.

## Install and Configure

### Manual
1. `git clone https://github.com/strick-j/scimplistic`
2. Change dir to respective folder and run go-build
3. Rename "example-settings.json" to "settings.json"
3. `./scimplistic`
4. Open localhost:8080

### Setup
Prior to using the web server you must configure your connection to the CyberArk SCIM endpoint provided by CyberArk Identity. Use the settings icon to configure you SCIM URL and Auth Token.

### Notes / Warning
This web app is in developement and should not be used in a production environment. Currently, the application lacks proper security measures and interacts directly with your PAM environment. Use at your own risk, ideally in a development environment.

## Currently Working Capabilities:
- Users
  - Retrieving all User data (GET https://scimurl/scim/v2/users) via the Users tab
  - Adding Users (POST https://scimurl/scim/v2/users)
  - Deleteing Users (DELETE https://scimurl/scim/v2/users)
- Users
  - Retrieving all Safe data (GET https://scimurl/scim/v2/containers) via the Safes tab
  - Adding Safes (POST https://scimurl/scim/v2/containers)
  - Deleteing Safes (DELETE https://scimurl/scim/v2/containers)
- Groups
  - Retrieving all Group data (GET https://scimurl/scim/v2/groups) via the Groups tab
  - Adding Groups (POST https://scimurl/scim/v2/groups)
  - Deleteing Groups (DELETE https://scimurl/scim/v2/groups)
- Accounts
  - Retrieving all Account data (GET https://scimurl/scim/v2/privilegeddata) via the Accounts tab
 

## Planned/Todo:
- Build out update functions to allow for update of users, groups, safes
- Add login/logout functionality
- Create review function to examine all of a users access based on direct Safe Access and Group based access
- Add database to track changes over time and allow for more complex queries (e.g. overall access)
- Add notifications for action completion status (Success / Failure / etc...)
- Add support for https

## Screenshots
Example Users Page:
![Users Page](https://github.com/strick-j/scimplistic/blob/main/screenshots/users.png)

Add User Form:
![Add User Form](https://github.com/strick-j/scimplistic/blob/main/screenshots/adduserform.png)

Example User Info:
![User Info](https://github.com/strick-j/scimplistic/blob/main/screenshots/userinfo.png)

Example Groups Page:
![Groups](https://github.com/strick-j/scimplistic/blob/main/screenshots/groups.png)


## Acknowledgements
Thanks to @thewhitetulip for the guides / book they created and the sample Tasks application [Tasks](https://github.com/thewhitetulip/Tasks)