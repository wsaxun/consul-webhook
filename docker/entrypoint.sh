#!/bin/sh

run_pre() {
    echo "run pre"
}

run_post() {
    echo "run post"
    rm -rf /tmp/*
}


run_in_online() {
    echo "run in online"
    run_pre
    mv docker/online.toml config.toml
    run_post
}

run_in_beta() {
    echo "run in beta"
    run_pre
    mv docker/beta.toml config.toml
    run_post
}

run_in_dev() {
    echo "run in develop"
    run_pre
    mv docker/dev.toml config.toml
    run_post
}


case ${RUN_ENV:=dev} in
    online)
        echo "online"
        run_in_online
        ;;
    beta)
        echo "beta"
        run_in_beta
        ;;
    dev)
        echo "develop"
        run_in_dev
        ;;
    *)
        echo "default"
        run_in_dev
        ;;
esac

exec "$@"
