#!/bin/bash
# Create an .env file
echo VIBER_KEY=$VIBER_KEY >>.env
echo CALLBACK_URL=$CALLBACK_URL >>.env
echo DB_CONNECTION=$DB_CONNECTION >>.env
echo DATASTORE_USER_ANSWER_LOG_TABLE=$DATASTORE_USER_ANSWER_LOG_TABLE >>.env
echo DATASTORE_USERS_TABLE=$DATASTORE_USERS_TABLE >>.env
echo ELASTIC_HOSTS=$ELASTIC_HOSTS >>.env
echo ELASTIC_BASIC_AUTH_USER=$ELASTIC_BASIC_AUTH_USER >>.env
echo ELASTIC_BASIC_AUTH_PASSWORD=$ELASTIC_BASIC_AUTH_PASSWORD >>.env
