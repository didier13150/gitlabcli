# GLCli

Le programme `glcli` permet une gestion des variables Gitlab, c'est à dire qu'il synchronise les fichiers json et les données du gitlab. Il peut le faire dans l'autre sens avec l'option `-export`.

Il utilise le format de fichier plat pour l'id du projet (fichier `.gitlab.id`) et le format JSON pour les environnements (fichier `.gitlab.env.json`) et les variables (fichier `.gitlab.var.json`).

Il a besoin de l'url du gitlab, ainsi que d'un token valide pour l'identification. Il n'est pas possible de passer directement le token à l'application, on ne peut que spécifier un fichier contenant ce token pour des raisons de sécurité.

L'application possède également un mode lecture seule, dans lequel seul les appels de lecture sont effectués: `-dryrun`

L'application peut également obtenir l'identifiant du projet à partir d'un export des projets (fichier `.gitlab.projects.json`)

## Usage de l'application

```
❯ ./glcli -help
Usage: ./glcli [--id <Project ID>] [--varfile <VAR FILE>] [--envfile <ENV FILE>] [--projectfile <PROJECT FILE>] [--token <TOKEN FILE>] [--dryrun] [--export] [--export-projects] [--delete]
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
  -id string
        Gitlab project identifiant.
  -projectfile string
        File which contains projects. (default "$HOME/.gitlab-projects.json")
  -token string
        File which contains token to access Gitlab API. (default "$HOME/.gitlab.token")
  -url string
        Gitlab URL. (default "https://gitlab.com")
  -varfile string
        File which contains vars. (default ".gitlab-vars.json")
  -verbose
        Make application more talkative.
```

Le mode debug exporte les tableaux des environnements et des variables dans le fichier `debug.txt`

Pour obtenir automatiquement l'identifiant du projet, il faut exporter les données concernant les projets avec l'option `-export-projects`

## Description des fichiers

* Fichier concernant l'**identifiant du projet** (fichier `.gitlab.id`). Ce fichier contient un nombre sur une ligne unique, sans espace, correspondant à l'id du projet.

    ```
    51
    ```
* Fichier concernant **les environnements**

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
    
    | Clé          | Description                    | Type de valeur                           | Valeur par défaut | Remarques                        |
    | ------------ | ------------------------------ | ---------------------------------------- | ----------------- | -------------------------------- |
    | id           | Identifiant de l'environnement | nombre entier                            |                   | non modifiable, en lecture seule |
    | name         | Nom de l'environnement         | chaîne de caractères non nulle           |                   | obligatoire                      |
    | state        | État: démarré/stoppé           | chaîne de caractères non nulle           |                   | non modifiable, en lecture seule |
    | external_url | URL de vérification            | chaîne de caractères qui peut être nulle | _null_            | facultatif pour la création      |
    | description  | Description de l'environnement | chaîne de caractères qui peut être nulle | _null_            | facultatif pour la création      |

* Fichier concernant **les variables**

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
    
    | Clé               | Description                                                                   | Type de valeur                           | Valeur par défaut | Remarques                        |
    | ----------------- | ----------------------------------------------------------------------------- | ---------------------------------------- | ----------------- | -------------------------------- |
    | key               | Clé de la variable (nom unique par environnement)                             | chaîne de caractères non nulle           |                   | obligatoire                      |
    | value             | Valeur de la variable                                                         | chaîne de caractères non nulle           |                   | obligatoire                      |
    | description       | Description de l'environnement                                                | chaîne de caractères qui peut être nulle | _null_            | facultatif pour la création      |
    | environment_scope | Portée de la variable                                                         | chaîne de caractères non nulle           | __*__             | obligatoire                      |
    | raw               | Drapeau indiquant que la variable est une variable non interprétable          | boolean                                  | false             | obligatoire                      |
    | hidden            | Drapeau indiquant que la variable doit être cachée dans le journal des *jobs* | boolean                                  | false             | obligatoire                      |
    | protected         | Drapeau indiquant que la variable est une variable protégée                   | boolean                                  | false             | obligatoire                      |
    | masked            | Drapeau indiquant que la variable est une variable masquée                    | boolean                                  | false             | obligatoire                      |
    

    * hidden: Masqué dans les journaux des *jobs* et ne peut jamais être révélé dans les pipelines une fois la variable enregistrée.
    * protected: Exporter la variable vers les pipelines exécutés uniquement sur des branches et des *tags* protégés.
    * masked: Masqué dans les journaux des *jobs*, mais la valeur peut être révélée dans les pipelines.

* Fichier concernant les projets, obtenu avec l'option `-export-projects`. Ce fichier peut être mutualisé pour tous les projets afin de s'affranchir la création de fichier `.gitlab.id` dans tous les dépôts locaux. 

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
## Utilisation

L'application peut utiliser des variables d'environnement afin de simplifier les options de la ligne de commande.

| Variable           | valeur par défaut           |
| ------------------ | --------------------------- |
| GLCLI_GITLAB_URL   | https://gitlab.com          |
| GLCLI_TOKEN_FILE   | $HOME/.gitlab.token         |
| GLCLI_PROJECT_FILE | $HOME/.gitlab.projects.json |
| GLCLI_VAR_FILE     | .gitlab-vars.json           |
| GLCLI_ENV_FILE     | .gitlab-envs.json           |
| GLCLI_ID_FILE      | .gitlab.id                  |
| GLCLI_DEBUG_FILE   | debug.txt                   |

Avant d'utiliser l'application, on doit d'abord inscrire l'identifiant du projet dans le fichier `.gitlab.id`. On peut se servir du script `get-projects-id.sh` afin d'obtenir une correspondance entre tous les projets et leur identifiant.

```
GITLAB_DOMAIN=gitlab.tartarefr.eu GITLAB_PRIV_TOKEN_FILE=./.${GITLAB_DOMAIN}.token ./get-projects-id.sh 
2 : sources/glcli
1 : arm64v8/gitlab-ce
```

### Export

Exporte les environnements et les variables existants depuis gitlab dans des fichiers. Cette action créé les fichiers `.gitlab.env.json` et `.gitlab.var.json`. si ces fichiers existent déjà, ils seront écrasés.

[ Gitlab ] ---> [ Fichiers ]

```
❯ ./glcli -export
```

### Import

Importe les environnements et les variables depuis les fichiers `.gitlab.env.json` et `.gitlab.var.json` dans gitlab. Par défaut les variables de gitlab non présentes dans les fichiers ne sont pas supprimées.

[ Fichiers ] ---> [ Gitlab ]

```
❯ ./glcli
```

Pour supprimer les variables surnumémaires il faut ajouter l'option `-delete`


## Exemples


### À partir d'un projet sans environnement ni variable.

1. Comme il n'y a ni variables ni environnement, on peut simplement créer les fichiers (il n'y a rien à exporter)

    ```
    echo -n '[]' > .gitlab.env.json
    echo -n '[]' > .gitlab.var.json
    ```
2. Modification des fichiers afin de spécifier les variables et environnements. Ajout d'un environnement nommé **production** et de trois **variables** dont une présente sur tous les environnement, mais surchargée pour l'environnement de production.
    
    * Fichier `.gitlab.env.json`
        
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
    * Fichier `.gitlab.var.json`
        
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
3. Import de nos déclarations dans gitlab

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
4. Export depuis gitlab afin de mettre à jour le champs **id** de l'environnement.

    ```
    ❯ ./glcli -export
    2025/08/02 13:08:18 Export requested
    2025/08/02 13:08:18 Fetching envs from gitlab with URL https://gitlab.tartarefr.eu
    2025/08/02 13:08:19 Fetching vars from gitlab with URL https://gitlab.tartarefr.eu
    2025/08/02 13:08:19 Export current Gitlab vars to .gitlab-vars.json file
    2025/08/02 13:08:19 Export current Gitlab envs to .gitlab-envs.json file
    2025/08/02 13:08:19 Exit now because export is done
    ```
5. Suppresion de la variable **GLCLI_VAR_LOCK_PREFIX** dans le fichier `.gitlab.var.json`

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
6. Synchronisation entre les fichiers et le gitlab. L'application reconnait bien la suppression mais ne trouve pas le drapeau permettant l'opération de suppression dans le ligne de commande.

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
7. Synchronisation entre les fichiers et le gitlab avec l'option de suppression

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
