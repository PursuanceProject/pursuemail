set -euo pipefail

export PGUSER=postgres
if [ "`uname -s`" != "Linux" ]; then
    # For Mac OS X
    export PGUSER=$USER
fi

source ../.env
psql -d postgres -c "CREATE USER pursuemail WITH PASSWORD '$PGPASSWORD';" || true
psql -d postgres -f sql/pre.sql
export PGHOST=localhost
export PGUSER=pursuemail
export PGDATABASE=pursuemail

# More initialization
for file in sql/init*.sql; do
    psql -f "$file"
done
# Create tables
for file in sql/table*.sql; do
    psql -f "$file"
done

/bin/bash migrate.sh sql/migration*.sql
