#!/bin/bash
# Create an .env file for staging
echo VIBER_KEY=$VIBER_KEY_STAG >> .env
echo CALLBACK_URL=$CALLBACK_URL_STAG >> .env
echo DB_CONNECTION=$DB_CONNECTION_STAG >> .env
