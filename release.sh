#!/bin/bash -x

function get_script_dir() {
    ## resolve the base directory of this executable
    local SOURCE=$1
    while [ -h "$SOURCE" ]; do
        # resolve $SOURCE until the file is no longer a symlink
        DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
        SOURCE="$(readlink "$SOURCE")"

        # if $SOURCE was a relative symlink,
        # we need to resolve it relative to the path
        # where the symlink file was located

        [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE"
    done

    echo "$( cd -P "$( dirname "$SOURCE" )" && pwd )"
}

function new_release_post_data() {
    t=1
    p=2
    cat <<EOF
{
    "tag_name": "${tag}",
    "tag_commitish": "master",
    "name": "${tag}",
    "body": "Release ${tag}",
    "draft": false,
    "prerelease": $2
}
EOF
}

if [ $# -ne 3 ]; then
    echo "$0 <release_tag> <is_prerelease> <github_username>"
    exit 1
fi

tag=$1
pre=$2
gh_uname=$3

GH_API="https://api.github.com"
GH_REPO="$GH_API/repos/Donders-Institute/torque-helper"
GH_RELE="$GH_REPO/releases"
GH_TAG="$GH_REPO/releases/tags/$tag"
GH_REPO_ASSET_PREFIX="https://uploads.github.com/repos/Donders-Institute/torque-helper/releases"

# check if version tag already exists
response=$(curl -X GET $GH_TAG 2>/dev/null)
eval $(echo "$response" | grep -m 1 "id.:" | grep -w id | tr : = | tr -cd '[[:alnum:]]=')
if [ "$id" ]; then
    read -p "release tag already exists: ${tag}, continue? y/[n]: " cnt
    if [ "${cnt,,}" != "y" ]; then
        exit 1
    fi
fi

# make sure the go command is available
which go > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "golang is required for building RPM."
    exit 1
fi 

# create a new tag with current master branch
# if the $id of the release is not available.
if [ ! "$id" ]; then
    response=$(curl -u ${gh_uname} -X POST --data "$(new_release_post_data ${tag} ${pre})" $GH_RELE)
    eval $(echo "$response" | grep -m 1 "id.:" | grep -w id | tr : = | tr -cd '[[:alnum:]]=')
    [ "$id" ] || { echo "release tag not created successfully: ${tag}"; exit 1; }
fi

# copy over id to rid (release id)
rid=$id

mydir=$( get_script_dir $0 )
path_spec=${mydir}/share/trqhelpd.centos7.spec

## replace the release version in 
out=$( VERSION=${tag} rpmbuild --undefine=_disable_source_fetch -bb ${path_spec} )
if [ $? -ne 0 ]; then
    echo "rpm build failure"
    exit 1
fi

## parse the RPM build output to get paths of output RPMs
rpms=( $( echo "${out}" | egrep -o -e 'Wrote:.*\.rpm' | sed 's/Wrote: //g' ) )

## upload RPMs as release assets 
if [ ${#rpms[@]} -gt 0 ]; then
    upload="y"
    read -p "upload ${#rpms[@]} RPMs as release assets?, continue? [y]/n: " upload
    for rpm in ${rpms[@]}; do
        echo ${rpm}
        if [ "${upload,,}" == "y" ]; then
            echo "uploading ${rpm} ..."
            fname=$( basename $rpm )
            # check if the asset with the same name already exists
            id=""
            eval $(echo "$response" | grep -C1 "name.:.\+${fname}" | grep -m 1 "id.:" | grep -w id | tr : = | tr -cd '[[:alnum:]]=')
            if [ "$id" != "" ]; then
                # delete existing asset
                echo "deleting asset: ${id} ..."
                curl -u ${gh_uname} -X DELETE "${GH_RELE}/assets/${asset_id}"
            fi
            # post new asset
            GH_ASSET="${GH_REPO_ASSET_PREFIX}/${rid}/assets?name=$(basename $rpm)"
            response=$( curl -u ${gh_uname} --data-binary @${rpm} \
                             -H "Content-Type: application/octet-stream" $GH_ASSET )
        fi
    done
fi
