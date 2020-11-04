#!/bin/bash

# read argument
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -b|--branch) branch="$2"; shift ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

if [ "$branch" = "" ]; then
  echo "Unknown parameter branch"
  exit 1;
fi

echo "Deploy to $branch"

# copy to remote server
if [ "$branch" = "dev" ]; then
  zip -r dist_metric.zip dist/
  scp dist_metric.zip lionnix-dev-main:/home/lionnix
  rm dist_metric.zip
elif [ "$branch" = "master" ]; then
  zip -r dist_metric.zip dist/
  scp dist_metric.zip lionnix-prod-main:/home/lionnix
  rm dist_metric.zip
else
  echo "Invalid branch"
  exit 1;
fi

# ssh remote server
if [ "$branch" = "dev" ]; then
  ssh lionnix-dev-main 'mkdir -p /usr/lionnix-dashboard/metric; mv /home/lionnix/dist_metric.zip /usr/lionnix-dashboard/metric; cd /usr/lionnix-dashboard/metric; zip -r dist.bak dist/; rm -rf dist; unzip dist_metric.zip; sudo systemctl reload nginx'

elif [ "$branch" = "master" ]; then
  ssh lionnix-prod-main 'mkdir -p /usr/lionnix-dashboard/metric; mv /home/lionnix/dist_metric.zip /usr/lionnix-dashboard/metric; cd /usr/lionnix-dashboard/metric; zip -r dist.bak dist/; rm -rf dist; unzip dist_metric.zip; sudo systemctl reload nginx'
else
  echo "Invalid branch"
  exit 1;
fi

exit