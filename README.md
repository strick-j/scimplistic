# SCIMPLISTIC
Scimplistic is a simplistic Go webapp to manage tasks, I built this tool to demonstrate the integration avaialble via SCIM for CyberArk Privileged Access Manager.

## Install and Configure
### Automated 
via Script: `bash install.sh`

### Manual
1. `go get github.com/strick-j/go-form-webserver`
2. Change dir to respective folder and run go-build
3. `./go-form-webserver`
4. Open localhost:8080

### Setup
Prior to using the web server you must configure your connection to the CyberArk SCIM endpoint provided by CyberArk Identity. Use the settings icon to configure you SCIM URL and Auth Token.

### Notes / Warning
This web app is in developement and should not be used in a production environment. Currently, the application lacks proper security measures and interacts directly with your PAM environment. Use at your own risk, ideally in a development environment.

## Currently Working Capabilities:
- Users
  - Retrieving all User data (GET https://scimurl/scim/users) via the Users tab
  - Adding Users (POST https://scimurl/scim/users)
  - Deleteing Users (DELETE https://scimurl/scim/users)
- Users
  - Retrieving all Safe data (GET https://scimurl/scim/containers) via the Safes tab
  - Adding Safes (POST https://scimurl/scim/containers)
  - Deleteing Safes (DELETE https://scimurl/scim/containers)
- Groups
  - Retrieving all Group data (GET https://scimurl/scim/groups) via the Groups tab
  - Adding Groups (POST https://scimurl/scim/groups)
  - Deleteing Groups (DELETE https://scimurl/scim/groups)

## Planned/Todo:
- Build out update functions to allow for update of users, groups, safes
- Add login/logout functionality
- Place check box / confirmation infront of "Delete" functions
- Create review function to examine all of a users access based on direct Safe Access and Group based access
- Add database to track changes over time and allow for more complex queries (e.g. overall access)

## Screenshots
Example Users Page:
![Users Page](https://github.com/strick-j/blob/master/screenshots/users.png)

Add User Form:
![Add User Form](https://github.com/strick-j/blob/master/screenshots/adduserform.png)

Example User Info:
![User Info](https://github.com/strick-j/blob/master/screenshots/userinfo.png)

Example Groups Page:
![Groups](https://github.com/strick-j/blob/master/screenshots/groups.png)


## Acknowledgements
Thanks to @thewhitetulip for the guides / book they created and the sample Tasks application [Tasks](https://github.com/thewhitetulip/Tasks)