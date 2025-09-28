# GLCli

The `glcli` program allows Gitlab variable and environment management, i.e., it synchronizes JSON files and Gitlab data. It can do this the other way around with the `-export` option.

It uses the flat file format for the project ID and its groupe ID (`.gitlab.id` and `.gitlab.gid` files) and the JSON format for environments (`.gitlab-envs.json` file) and variables (`.gitlab-vars.json` file).

It requires the Gitlab URL and a valid token for identification. It is not possible to pass the token directly to the application; you can only specify a file containing this token for security reasons.

The application also has a read-only mode, in which only read calls are made: `-dryrun`

The application can also get the project ID from a project export data (`.gitlab-projects.json` file)

## Installation

```
go install github.com/didier13150/gitlabcli@latest
```

## Application Usage

```
❯ ./glcli -help
Usage: ./glcli [options]
  -all-projects
        Export all projects, not only projects where I'm a membership.
  -debug
        Enable debug mode
  -delete
        Delete Gitlab var if not present in var file.
  -dryrun
        Run in dry-run mode (read only).
  -envfile string
        File which contains envs. (default ".gitlab-envs.json")
  -export
        Export current variables in var file.
  -export-projects
        Export current projects in project file.
  -full-projects-data
        Requesting full data about projects.
  -gid string
        Gitlab group identifiant.
  -gidfile string
        Gitlab group identifiant file. (default ".gitlab.gid")
  -groupvarfile string
        File which contains group vars. (default ".gitlab-groupvars.json")
  -id string
        Gitlab project identifiant.
  -idfile string
        Gitlab project identifiant file. (default ".gitlab.id")
  -projectfile string
        File which contains projects. (default "$HOME/.gitlab-projects.json"))
  -remote string
        Git remote name.
  -token string
        File which contains token to access Gitlab API. (default "$HOME/.gitlab.token")
  -url string
        Gitlab URL. (default "https://gitlab.com")
  -varfile string
        File which contains vars. (default ".gitlab-vars.json")
  -verbose
        Make application more talkative.
```

Debug mode exports the environment and variable tables to the `debug.txt` file.

## File Descriptions

* **Project ID** file (`.gitlab.id` file). This file contains a single-line, without spaces, corresponding to the project ID.

    ```
    51
    ```
* **Group ID** file (`.gitlab.gid` file). This file contains a single-line, without spaces, corresponding to the project GID.

    ```
    69
    ```
* **Environment** file

    ```
    [
      {
        "id": 10,
        "name": "review",
        "state": "available",
        "external_url": "https://app.example.com",
        "description": "Production environment"
      }
    ]
    ```

    | Key          | Description             | Value Type      | Default Value | Notes                   |
    | ------------ | ----------------------- | --------------- | ------------- | ----------------------- |
    | id           | Environment ID          | integer         |               | not editable, read-only |
    | name         | Environment name        | non-null string |               | required                |
    | state        | State: started/stopped  | non-null string |               | not editable, read-only |
    | external_url | Check URL               | nullable string | _null_        | optional for creation   |
    | description  | Environment description | nullable string | _null_        | optional for creation   |

* **variables** file

    ```
    [
      {
        "key": "DEBUG_ENABLED",
        "value": "1",
        "description": null,
        "environment_scope": "*",
        "raw": true,
        "hidden": false,
        "protected": false,
        "masked": false
      }
    ]
    ```
    | Key               | Description                                                         | Value Type      | Default Value | Notes                 |
    | ----------------- | ------------------------------------------------------------------- | --------------- | ------------- | --------------------- |
    | key               | Variable key (unique name per environment)                          | non-null string |               | required              |
    | value             | Variable value                                                      | non-null string |               | required              |
    | description       | Environment description                                             | nullable string | _null_        | optional for creation |
    | environment_scope | Variable scope                                                      | non-null string | __*__         | required              |
    | raw               | Flag indicating that the variable is uninterpretable                | boolean         | false         | required              |
    | hidden            | Flag indicating that the variable should be hidden in the *job* log | boolean         | false         | required              |
    | protected         | Flag indicating that the variable is a protected variable           | boolean         | false         | required              |
    | masked            | Flag indicating that the variable is a masked variable              | boolean         | false         | required              |

    * hidden: Hidden from job logs and can never be revealed in pipelines once the variable is saved.
    * protected: Export the variable to pipelines running only on protected branches and tags.
    * masked: Hidden from job logs, but the value can be revealed in pipelines.

* **Project** file, obtained with the `-export-projects` option. This file can be shared across all projects to avoid creating `.gitlab.id` files in all local repositories.

    ```
    [
      {
        "id": 2,
        "name": "GLCli",
        "description": null,
        "path": "glcli",
        "name_with_namespace": "Sources / GLCli",
        "path_with_namespace": "sources/glcli",
        "ssh_url_to_repo": "git@gitlab.tartarefr.eu:sources/glcli.git",
        "http_url_to_repo": "https://gitlab.tartarefr.eu/sources/glcli.git",
        "web_url": "https://gitlab.tartarefr.eu/sources/glcli",
        "visibility": "public"
      },
      {
        "id": 1,
        "name": "Gitlab Community Edition",
        "description": null,
        "path": "gitlab-ce",
        "name_with_namespace": "ARM 64 version 8 / Gitlab Community Edition",
        "path_with_namespace": "arm64v8/gitlab-ce",
        "ssh_url_to_repo": "git@gitlab.tartarefr.eu:arm64v8/gitlab-ce.git",
        "http_url_to_repo": "https://gitlab.tartarefr.eu/arm64v8/gitlab-ce.git",
        "web_url": "https://gitlab.tartarefr.eu/arm64v8/gitlab-ce",
        "visibility": "public"
      }
    ]
    ```

## Runtime

The application can use environment variables to simplify command-line options.

| Variable             | valeur par défaut           |
| -------------------- | --------------------------- |
| GLCLI_GITLAB_URL     | https://gitlab.com          |
| GLCLI_TOKEN_FILE     | $HOME/.gitlab.token         |
| GLCLI_PROJECT_FILE   | $HOME/.gitlab.projects.json |
| GLCLI_VAR_FILE       | .gitlab-vars.json           |
| GLCLI_GROUP_VAR_FILE | .gitlab-groupvars.json      |
| GLCLI_ENV_FILE       | .gitlab-envs.json           |
| GLCLI_ID_FILE        | .gitlab.id                  |
| GLCLI_GROUP_ID_FILE  | .gitlab.gid                 |
| GLCLI_DEBUG_FILE     | debug.txt                   |

Before using the application, you must first enter the project ID in the `.gitlab.id` file or using an export of projects.

### Export

Exports existing environments and variables from Gitlab into files. This creates the `.gitlab-envs.json` and `.gitlab-vars.json` files. If these files already exist, they will be overwritten.

[ Gitlab ] ---> [ Files ]

```
❯ ./glcli -export
```

### Import

Imports environments and variables from the `.gitlab-envs.json` and `.gitlab-vars.json` files into Gitlab. By default, Gitlab variables not present in the files are not deleted.

[ Files ] ---> [ Gitlab ]

```
❯ ./glcli
```

To delete extra variables, add the `-delete` option.

## Examples

### Starting with a project without an environment or variables.

1. Since there are no variable or environment, we can simply create the files (there's nothing to export).

    ```
    echo -n '[]' > .gitlab-envs.json
    echo -n '[]' > .gitlab-vars.json
    ```
2. Modify the files to specify the variables and environments. Add an environment named **production** and three **variables**, one of which is present in all environments, but overridden for the production environment.

    * `.gitlab-envs.json` file 
        
        ``` 
        [ 
          { 
            "id": 0, 
            "name": "production", 
            "state": "available", 
            "external_url": null, 
            "description": "Production environment" 
          } 
        ] 
        ``` 
    * `.gitlab-vars.json` file 
        
        ``` 
        [ 
          { 
            "key": "DEBUG_ENABLED", 
            "value": "1", 
            "description": null, 
            "environment_scope": "*", 
            "raw": true, 
            "hidden": false, 
            "protected": false, 
            "masked": false 
          }, 
          { 
            "key": "DEBUG_ENABLED", 
            "value": "0", 
            "description": null, 
            "environment_scope": "production", 
            "raw": true, 
            "hidden": false, 
            "protected": false, 
            "masked": false 
          }, 
          { 
            "key": "VAR_PREFIX", 
            "value": "GLCLI", 
            "description": null, 
            "environment_scope": "*", 
            "raw": true, 
            "hidden": false, 
            "protected": false, 
            "masked": false 
          }, 
          { 
            "key": "GLCLI_VAR_LOCK_PREFIX", 
            "value": "/var/lock", 
            "description": "Prefix for lock file", 
            "environment_scope": "*", 
            "raw": true, 
            "hidden": false, 
            "protected": false, 
            "masked": false 
          } 
        ] 
        ```
3. Import our declarations into gitlab 

    ``` 
    ❯ ./glcli 
    2025/08/02 13:07:43 Fetching envs from gitlab with URL https://gitlab.tartarefr.eu 
    2025/08/02 13:07:44 Fetching vars from gitlab with URL https://gitlab.tartarefr.eu 
    2025/08/02 13:07:44 Env {0 production available <nil> Production environment} should be added 
    2025/08/02 13:07:44 Use URL https://gitlab.tartarefr.eu/api/v4/projects/52/environments to insert env 
    2025/08/02 13:07:44 Insert env production (0) 
    2025/08/02 13:07:44 No env to update 
    2025/08/02 13:07:44 No env to delete 
    2025/08/02 13:07:44 Var {DEBUG_ENABLED 1 <nil> * true false false false} should be added 
    2025/08/02 13:07:44 Var {DEBUG_ENABLED 0 <nil> production true false false false} should be added 
    2025/08/02 13:07:44 Var {VAR_PREFIX GLCLI <nil> * true false false false} should be added 
    2025/08/02 13:07:44 Var {GLCLI_VAR_LOCK_PREFIX /var/lock Prefix for lock file * true false false false} should be added 
    2025/08/02 13:07:44 Use URL https://gitlab.tartarefr.eu/api/v4/projects/52/variables to insert var 
    2025/08/02 13:07:44 Insert var DEBUG_ENABLED in * env 
    2025/08/02 13:07:44 Use URL https://gitlab.tartarefr.eu/api/v4/projects/52/variables to insert var 
    2025/08/02 13:07:44 Insert var DEBUG_ENABLED in production env 
    2025/08/02 13:07:44 Use URL https://gitlab.tartarefr.eu/api/v4/projects/52/variables to insert var 
    2025/08/02 13:07:44 Insert var VAR_PREFIX in * env 
    2025/08/02 13:07:44 Use URL https://gitlab.tartarefr.eu/api/v4/projects/52/variables to insert var
    2025/08/02 13:07:44 Insert var GLCLI_VAR_LOCK_PREFIX in * env
    2025/08/02 13:07:45 No var to update
    2025/08/02 13:07:45 No var to delete
    2025/08/02 13:07:45 Exit
    ```
4. Export from Gitlab to update the environment's **id** field.

    ``` 
    ❯ ./glcli -export 
    2025/08/02 13:08:18 Export requested 
    2025/08/02 13:08:18 Fetching envs from gitlab with URL https://gitlab.tartarefr.eu 
    2025/08/02 13:08:19 Fetching vars from gitlab with URL https://gitlab.tartarefr.eu 
    2025/08/02 13:08:19 Export current Gitlab vars to .gitlab-vars.json file 
    2025/08/02 13:08:19 Export current Gitlab envs to .gitlab-envs.json file 
    2025/08/02 13:08:19 Exit now because export is done 
    ```
5. Deleting the variable **GLCLI_VAR_LOCK_PREFIX** in the `.gitlab-vars.json` file 

    ``` 
    [ 
      { 
        "key": "DEBUG_ENABLED", 
        "value": "1", 
        "description": null, 
        "environment_scope": "*", 
        "raw": true, 
        "hidden": false, 
        "protected": false, 
        "masked": false 
      }, 
      { 
        "key": "DEBUG_ENABLED", 
        "value": "0", 
        "description": null, 
        "environment_scope": "production", 
        "raw": true, 
        "hidden": false, 
        "protected": false, 
        "masked": false 
      }, 
      { 
        "key": "VAR_PREFIX", 
        "value": "GLCLI", 
        "description": null, 
        "environment_scope": "*", 
        "raw": true, 
        "hidden": false,
        "protected": false,
        "masked": false
      }
    ]
    ```
6. Synchronization between files and GitLab. The application recognizes the deletion but cannot find the flag allowing the deletion operation in the command line.

    ``` 
    ❯ ./glcli 
    2025/08/02 13:21:14 Fetching envs from gitlab with URL https://gitlab.tartarefr.eu 
    2025/08/02 13:21:14 Fetching vars from gitlab with URL https://gitlab.tartarefr.eu 
    2025/08/02 13:21:14 No env to insert 
    2025/08/02 13:21:14 No env to update 
    2025/08/02 13:21:14 No env to delete 
    2025/08/02 13:21:14 Var {GLCLI_VAR_LOCK_PREFIX /var/lock Prefix for lock file * true false false false} should be deleted 
    2025/08/02 13:21:14 No var to insert 
    2025/08/02 13:21:14 No var to update 
    2025/08/02 13:21:14 1 var(s) may be deleted, but delete flag in command line is not set 
    2025/08/02 13:21:14 Exit 
    ```
7. Sync between files and gitlab with delete option 

    ``` 
    ❯ ./glcli -delete 
    2025/08/02 13:22:32 Fetching envs from gitlab with URL https://gitlab.tartarefr.eu 
    2025/08/02 13:22:33 Fetching vars from gitlab with URL https://gitlab.tartarefr.eu 
    2025/08/02 13:22:33 No env to insert 
    2025/08/02 13:22:33 No env to update 
    2025/08/02 13:22:33 No env to delete 
    2025/08/02 13:22:33 Var {GLCLI_VAR_LOCK_PREFIX /var/lock Prefix for lock file * true false false false} should be deleted 
    2025/08/02 13:22:33 No var to insert 
    2025/08/02 13:22:33 No var to update 
    2025/08/02 13:22:33 Use URL https://gitlab.tartarefr.eu/api/v4/projects/52/variables/GLCLI_VAR_LOCK_PREFIX?filter[environment_scope]=* to delete var 
    2025/08/02 13:22:33 Delete var GLCLI_VAR_LOCK_PREFIX in * env 
    2025/08/02 13:22:33 Exit 
    ```
