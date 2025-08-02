# GLVars

The `glvars` program allows Gitlab variable management, i.e., it synchronizes JSON files and Gitlab data. It can do this the other way around with the `-export` option.

It uses the flat file format for the project ID (`.gitlab.id` file) and the JSON format for environments (`.gitlab.env.json` file) and variables (`.gitlab.var.json` file).

It requires the Gitlab URL and a valid token for identification. It is not possible to pass the token directly to the application; you can only specify a file containing this token for security reasons.

The application also has a read-only mode, in which only read calls are made: `-dryrun`

## Application Usage

```
❯ ./glvars -help
Usage: ./glvars [--id <Poject ID>] [--varfile <VAR FILE>] [--envfile <ENV FILE>] [--token <TOKEN FILE>] [--dryrun] [--export] [--delete]
  -debug
        Enable debug mode
  -delete
        Delete Gitlab var if not present in var file. Default is false.
  -dryrun
        Run in dry-run mode (read only). Default is false.
  -envfile string
        File which contains envs. Default is '.gitlab-envs.json' in the current directory. (default ".gitlab-envs.json")
  -export
        Export current variables in var file. Default is false.
  -id string
        Gitlab project identifiant. Default is to read it from '.gitlab.id' file in the current directory.
  -token string
        File which contains token to access Gitlab API. Default is '/home/didier/.gitlab.tartarefr.eu.token' (default "/home/didier/.gitlab.tartarefr.eu.token")
  -url string
        Gitlab URL. Default is 'https://gitlab.tartarefr.eu' (default "https://gitlab.tartarefr.eu")
  -varfile string
        File which contains vars. Default is '.gitlab-vars.json' in the current directory. (default ".gitlab-vars.json")
  -verbose
        Make application more talkative. Default is false.
```

Debug mode exports the environment and variable tables to the `debug.txt` file.

## File Descriptions

* **Project ID** file (`.gitlab.id` file). This file contains a single-line, without spaces, corresponding to the project ID.

    ```
    51
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

## Runtime

The application can use environment variables to simplify command-line options.

| Variable          | Default value       |
| ----------------- | ------------------- |
| GLVARS_GITLAB_URL | https://gitlab.com  |
| GLVARS_TOKEN_FILE | $HOME/.gitlab.token |
| GLVARS_VAR_FILE   | .gitlab-vars.json   |
| GLVARS_ENV_FILE   | .gitlab-envs.json   |
| GLVARS_ID_FILE    | .gitlab.id          |
| GLVARS_DEBUG_FILE | debug.txt           |

Before using the application, you must first enter the project ID in the `.gitlab.id` file. You can use the `get-projects-id.sh` script to obtain a mapping between all projects and their IDs.

```
GITLAB_DOMAIN=gitlab.tartarefr.eu GITLAB_PRIV_TOKEN_FILE=./.${GITLAB_DOMAIN}.token ./get-projects-id.sh
2: sources/glvars
1: arm64v8/gitlab-ce
```

### Export

Exports existing environments and variables from Gitlab into files. This creates the `.gitlab.env.json` and `.gitlab.var.json` files. If these files already exist, they will be overwritten.

[ Gitlab ] ---> [ Files ]

```
❯ ./glvars -export
```

### Import

Imports environments and variables from the `.gitlab.env.json` and `.gitlab.var.json` files into Gitlab. By default, Gitlab variables not present in the files are not deleted.

[ Files ] ---> [ Gitlab ]

```
❯ ./glvars
```

To delete extra variables, add the `-delete` option.

## Examples

### Starting with a project without an environment or variables.

1. Since there are no variable or environment, we can simply create the files (there's nothing to export).

    ```
    echo -n '[]' > .gitlab.env.json
    echo -n '[]' > .gitlab.var.json
    ```
2. Modify the files to specify the variables and environments. Add an environment named **production** and three **variables**, one of which is present in all environments, but overridden for the production environment.

    * `.gitlab.env.json` file 
        
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
    * `.gitlab.var.json` file 
        
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
            "value": "GLVARS", 
            "description": null, 
            "environment_scope": "*", 
            "raw": true, 
            "hidden": false, 
            "protected": false, 
            "masked": false 
          }, 
          { 
            "key": "GLVARS_VAR_LOCK_PREFIX", 
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
    ❯ ./glvars 
    2025/08/02 13:07:43 Fetching envs from gitlab with URL https://gitlab.tartarefr.eu 
    2025/08/02 13:07:44 Fetching vars from gitlab with URL https://gitlab.tartarefr.eu 
    2025/08/02 13:07:44 Env {0 production available <nil> Production environment} should be added 
    2025/08/02 13:07:44 Use URL https://gitlab.tartarefr.eu/api/v4/projects/52/environments to insert env 
    2025/08/02 13:07:44 Insert env production (0) 
    2025/08/02 13:07:44 No env to update 
    2025/08/02 13:07:44 No env to delete 
    2025/08/02 13:07:44 Var {DEBUG_ENABLED 1 <nil> * true false false false} should be added 
    2025/08/02 13:07:44 Var {DEBUG_ENABLED 0 <nil> production true false false false} should be added 
    2025/08/02 13:07:44 Var {VAR_PREFIX GLVARS <nil> * true false false false} should be added 
    2025/08/02 13:07:44 Var {GLVARS_VAR_LOCK_PREFIX /var/lock Prefix for lock file * true false false false} should be added 
    2025/08/02 13:07:44 Use URL https://gitlab.tartarefr.eu/api/v4/projects/52/variables to insert var 
    2025/08/02 13:07:44 Insert var DEBUG_ENABLED in * env 
    2025/08/02 13:07:44 Use URL https://gitlab.tartarefr.eu/api/v4/projects/52/variables to insert var 
    2025/08/02 13:07:44 Insert var DEBUG_ENABLED in production env 
    2025/08/02 13:07:44 Use URL https://gitlab.tartarefr.eu/api/v4/projects/52/variables to insert var 
    2025/08/02 13:07:44 Insert var VAR_PREFIX in * env 
    2025/08/02 13:07:44 Use URL https://gitlab.tartarefr.eu/api/v4/projects/52/variables to insert var
    2025/08/02 13:07:44 Insert var GLVARS_VAR_LOCK_PREFIX in * env
    2025/08/02 13:07:45 No var to update
    2025/08/02 13:07:45 No var to delete
    2025/08/02 13:07:45 Exit
    ```
4. Export from Gitlab to update the environment's **id** field.

    ``` 
    ❯ ./glvars -export 
    2025/08/02 13:08:18 Export requested 
    2025/08/02 13:08:18 Fetching envs from gitlab with URL https://gitlab.tartarefr.eu 
    2025/08/02 13:08:19 Fetching vars from gitlab with URL https://gitlab.tartarefr.eu 
    2025/08/02 13:08:19 Export current Gitlab vars to .gitlab-vars.json file 
    2025/08/02 13:08:19 Export current Gitlab envs to .gitlab-envs.json file 
    2025/08/02 13:08:19 Exit now because export is done 
    ```
5. Deleting the variable **GLVARS_VAR_LOCK_PREFIX** in the `.gitlab.var.json` file 

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
        "value": "GLVARS", 
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
    ❯ ./glvars 
    2025/08/02 13:21:14 Fetching envs from gitlab with URL https://gitlab.tartarefr.eu 
    2025/08/02 13:21:14 Fetching vars from gitlab with URL https://gitlab.tartarefr.eu 
    2025/08/02 13:21:14 No env to insert 
    2025/08/02 13:21:14 No env to update 
    2025/08/02 13:21:14 No env to delete 
    2025/08/02 13:21:14 Var {GLVARS_VAR_LOCK_PREFIX /var/lock Prefix for lock file * true false false false} should be deleted 
    2025/08/02 13:21:14 No var to insert 
    2025/08/02 13:21:14 No var to update 
    2025/08/02 13:21:14 1 var(s) may be deleted, but delete flag in command line is not set 
    2025/08/02 13:21:14 Exit 
    ```
7. Sync between files and gitlab with delete option 

    ``` 
    ❯ ./glvars -delete 
    2025/08/02 13:22:32 Fetching envs from gitlab with URL https://gitlab.tartarefr.eu 
    2025/08/02 13:22:33 Fetching vars from gitlab with URL https://gitlab.tartarefr.eu 
    2025/08/02 13:22:33 No env to insert 
    2025/08/02 13:22:33 No env to update 
    2025/08/02 13:22:33 No env to delete 
    2025/08/02 13:22:33 Var {GLVARS_VAR_LOCK_PREFIX /var/lock Prefix for lock file * true false false false} should be deleted 
    2025/08/02 13:22:33 No var to insert 
    2025/08/02 13:22:33 No var to update 
    2025/08/02 13:22:33 Use URL https://gitlab.tartarefr.eu/api/v4/projects/52/variables/GLVARS_VAR_LOCK_PREFIX?filter[environment_scope]=* to delete var 
    2025/08/02 13:22:33 Delete var GLVARS_VAR_LOCK_PREFIX in * env 
    2025/08/02 13:22:33 Exit 
    ```
