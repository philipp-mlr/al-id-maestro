<p align="center">
  <picture>
    <source srcset="https://github.com/philipp-mlr/al-id-maestro/blob/main/public/img/logo.png?raw=true" width="25%" height="25%">
    <img src="https://github.com/philipp-mlr/al-id-maestro/blob/main/public/img/logo.png?raw=true" width="25%" height="25%">
  </picture>
</p>

<!-- omit in toc -->

# al-id-maestro
Streamline object ID management for Business Central app development, ensuring consistency and avoiding collisions in your AL development workflow.
Whether you're working alone or part of a large team, **al-object-maestro** helps you maintain accurate and unique object IDs across all features, branches, and repositories.

## Docker

You need to map a volume to the `/app/data` directory. Inside this directory, the application will
- look for the `config.yml`
- create a database
- store repositories

Use the following command to run the container locally:

```
docker run -d -p 8080:8080 -v "C:/path:/app/data" --name al-id-maestro ghcr.io/philipp-mlr/al-id-maestro:main 
```

## Config

The container needs a file named `config.yml` mounted to the `/app/data` directory.
For more information about the configuration options, checkout the [config.yml example file](https://github.com/philipp-mlr/al-id-maestro/blob/main/data/config.yml.example).

### Repositories
You may define multiple repositories like so:

```
repositories:
  - name: repo-1
    url: https://github.com/myorg/cool-repo
    authToken: ghp_ABC123
    remoteName: origin
    excludeBranches:
      - release/
  - name: repo-2
    url: https://github.com/myuser/awesome-apps
    authToken: ghp_ABC123
    remoteName: origin
    excludeBranches:
      - somebranch/
```

####  Settings

| Setting         	| Description                                                             	| State    	|
|-----------------	|-------------------------------------------------------------------------	|----------	|
| name            	| A friendly name for the repository                                      	| required 	|
| authToken       	| Github auth token                                                       	| required 	|
| remoteName      	| Remote name                                                             	| required 	|
| excludeBranches 	| Array of branch patterns which get ignored during scan. Omit * asterisk 	| optional 	|

### ID Ranges
You have to define ID ranges for the following object types:

#### Settings

| Object type            	| State    	|
|------------------------	|----------	|
| Page                   	| required 	|
| PageExtension          	| required 	|
| Table                  	| required 	|
| TableExtension         	| required 	|
| Enum                   	| required 	|
| EnumExtension          	| required 	|
| Report                 	| required 	|
| ReportExtension        	| required 	|
| PermissionSet          	| required 	|
| PermisisonSetExtension 	| required 	|
| Codeunit               	| required 	|
| XMLPort                	| required 	|
| MenuSuite              	| required 	|

You may define multiple ranges for the same object. They must not overlap.

This will work:
```  yaml
idRanges:
  - objectType: Codeunit
    from: 50000
    to:   60000
  - objectType: Codeunit
    from: 60000
    to:   70000
```
This won't work: overlapping ID ranges!
```  yaml
idRanges:
  - objectType: Codeunit
    from: 50000
    to:   60000
  - objectType: Codeunit
    from: 59000
    to:   70000
```

## Found a bug?
Please open an issue.


