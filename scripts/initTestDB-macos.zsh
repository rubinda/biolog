#!/bin/zsh
# Ustvari testno podatkovno bazo na Postgres. Uporabljeno na macOS 10.13.3.
# Za delovanje potrebuje PostgreSQL bin v PATH, prav tako pa veljavno datoteko .pgpass
# 
# Skripta sprejme 3 opcijske parametre, katerih privzete vrednosti so spodaj
# Primer delovanja: ./initTestDB-macosh.zsh <test_db> <db> <user>
#   test_db ... ime testne podatkovne baze
#   db      ... ime baze od aplikacije biolog
#   user    ... uporabnik, ki ima dostop do baze db
#
# @author David Rubin
TEST_DB=${1:-biolog_development}
DBNAME=${2:-biolog}
USER=${3:-biolog}


# Check if a database with the same name already exists and drop it
DB_EXISTS=`psql -U $USER -tc "SELECT 1 FROM pg_database WHERE datname = '$TEST_DB'"`
if [[ $DB_EXISTS ]]; then
    dropdb -U $USER $TEST_DB
fi
# Create a new database with the selected name
createdb -U $USER $TEST_DB
# Dump the original database schema into a file
pg_dump -U $USER  --schema-only --format c $DBNAME > biolog.dump
# Restore the previously dumped schema into the newly created database
pg_restore -U postgres --clean --schema-only -d $TEST_DB biolog.dump
# Repair public schema permissions
psql -U postgres -d $TEST_DB -c "GRANT ALL ON SCHEMA public TO public"
# Delete the dump file (no longer needed)
rm ./biolog.dump
# Get the directory where the script is currently at
DIR=`dirname $0`
# Load some test data from './sample-data.sql' into the new test database
psql -U $USER -d $TEST_DB -f $DIR/sample-data.sql
