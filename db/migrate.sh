set -euo pipefail

export PGHOST=localhost
export PGUSER=postgres
export PGDATABASE=pursuemail
if [ "`uname -s`" != "Linux" ]; then
    # For Mac OS X
    export PGUSER=$USER
fi

# Run migrations
for file in $*; do
    psql -f "$file"
done
