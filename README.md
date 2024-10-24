<p align="center">
  <picture>
    <source srcset="https://github.com/philipp-mlr/al-id-maestro/blob/main/public/img/logo.png?raw=true" width="25%" height="25%">
    <img src="https://github.com/philipp-mlr/al-id-maestro/blob/main/public/img/logo.png?raw=true" width="25%" height="25%">
  </picture>
</p>

<!-- omit in toc -->

# AL-ID Maestro
The ultimate tool for managing object IDs in Business Central

## Docker

You need to map a volume to the /app/data directory. Inside this directory, the application will
- look for the config.yml
- create a database
- store repositories

Use the following command to run the container locally:

```
docker run -d -p 8080:8080 -v "C:/path:/app/data" --name al-id-maestro ghcr.io/philipp-mlr/al-id-maestro:main 
```

## Config

The container needs a file named config.yml mounted to the /app/data directory.
For more information about the configuration options, checkout the [config.yml example file](https://github.com/philipp-mlr/al-id-maestro/blob/main/data/config.yml.example).


### ID Ranges
You have to define ID ranges for the following object types:
- Table
- TableExtension
- Page
- PageExtension
- Report
- ReportExtension
- Enum
- EnumExtension
- PermissionSet
- PermissionSetExtension
- Codeunit
- Query
- XMLPort
- Menusuite

You may define multiple ranges for the same object. They must not overlap.

Works
```  yaml
idRanges:
  - objectType: Codeunit
    from: 50000
    to:   60000
  - objectType: Codeunit
    from: 60000
    to:   70000
```
Doesn't work because overlapping id ranges
```  yaml
idRanges:
  - objectType: Codeunit
    from: 50000
    to:   60000
  - objectType: Codeunit
    from: 59000
    to:   70000
```

