#!/bin/bash

# Get id and path with namespace from all projects on gitlab

GITLAB_PRIV_TOKEN=${GITLAB_PRIV_TOKEN:-}
GITLAB_DOMAIN=gitlab.tartarefr.eu
GITLAB_API_URL=${GITLAB_API_URL:-https://${GITLAB_DOMAIN}/api/v4}
PER_PAGE=${PER_PAGE:-50}
GITLAB_PRIV_TOKEN_FILE=${HOME}/.${GITLAB_DOMAIN}.token

[ -z "${GITLAB_PRIV_TOKEN}" ] && [ -f ${GITLAB_PRIV_TOKEN_FILE} ] && GITLAB_PRIV_TOKEN=$(cat ${GITLAB_PRIV_TOKEN_FILE} | tr -d '\n')
[ -z "${GITLAB_PRIV_TOKEN}" ] && echo -e "No token found !\nRetry with GITLAB_PRIV_TOKEN=\"<TOKEN>\" $0 or save it in ${GITLAB_PRIV_TOKEN_FILE} file" && exit 1

# Get the number of page (${PER_PAGE} entries by page)
nb=$(curl -kfsSL --head "${GITLAB_API_URL}/projects?private_token=${GITLAB_PRIV_TOKEN}&per_page=${PER_PAGE}" | grep -i "x-total-pages" | cut -d':' -f2 | tr -d ' "\r')

for page in $(seq 1 ${nb})
do
    curl -k -fsSL "${GITLAB_API_URL}/projects?private_token=${GITLAB_PRIV_TOKEN}&per_page=${PER_PAGE}&page=${page}" | jq -r '.[] | "\(.id) : \(.path_with_namespace)"' | tr -d '"'
done
